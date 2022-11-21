package net

import (
	"fmt"
	"errors"
	"time"
	"net"
	"net/url"
	"net/http"
	"crypto/tls"
	"sync"
	"sync/atomic"
	"context"

	"nhooyr.io/websocket"
)

// TODO - Ensure sent messages remain under this
// Calculation: 1460 Byte = 1500 Byte - 20 Byte IP Header - 20 Byte TCP Header
const MaxMsgSize = 1460 // bytes
const MaxRecvMsgSize = 4 * 1024 // 4 KB // TODO - this is arbitrary

var ErrSerdes = errors.New("serdes errror")
var ErrNetwork = errors.New("network error")

// This is the interface used to marshal and unmarshal messages over the network.
// Hardest part is when you have interfaces in messages, you'll likely need custom serializers for that
type Serdes interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(dat []byte) (any, error)
}

// type Socket interface {
// 	net.Conn
// 	// net.PacketConn???
// 	Send(any) error
// 	Recv() (any, error)
// }

type Listener interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (*Socket, error)

	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error

	// Addr returns the listener's network address.
	Addr() net.Addr
}

type SocketListener struct {
	listener net.Listener
	serdes Serdes
}
func (l *SocketListener) Accept() (*Socket, error) {
	c, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	framedConn := NewFrameConn(c)
	return newConnectedSocket(framedConn, l.serdes), nil
}
func (l *SocketListener) Close() error {
	return l.listener.Close()
}
func (l *SocketListener) Addr() net.Addr {
	return l.listener.Addr()
}

// TODO - (When I migrate to TCP) TCP will send 0 byte messages to indicate closes, websockets sends them without closing
type WebsocketListener struct {
	httpServer http.Server
	originPatterns []string
	addr net.Addr
	serdes Serdes
	pendingAccepts chan *Socket // TODO - should this get buffered?
	pendingAcceptErrors chan error // TODO - should this get buffered?
}
func newWebsocketListener(c *Config) (*WebsocketListener, error) {
	u, err := url.Parse(c.Url)
	if err != nil {
		return nil, err
	}

	// TODO - is TCP always correct?
	listener, err := tls.Listen("tcp", u.Host, c.TlsConfig)
	if err != nil {
		panic(err)
	}

	wsl := &WebsocketListener{
		serdes: c.Serdes,
		addr: listener.Addr(),
		pendingAccepts: make(chan *Socket),
		pendingAcceptErrors: make(chan error),
		originPatterns: c.OriginPatterns,
	}

	httpServer := c.HttpServer
	httpServer.Handler = wsl

	go func() {
		for {
			err := httpServer.Serve(listener)
			// TODO - what happens if this continually fails, how do we notify back?
			// ErrServerClosed is returned when shutdown or close is called
			fmt.Println("Serving Error:", err)

			if errors.Is(err, http.ErrServerClosed) {
				return // Just close if the server is shutdown or closed
			}

			time.Sleep(1 * time.Second)
		}
	}()

	return wsl, nil
}

func (l *WebsocketListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: l.originPatterns,
	})
	if err != nil {
		// Return as an accept error
		l.pendingAcceptErrors <- err
		return
	}

	// Build the socket and push to channel
	ctx := context.Background()
	conn := websocket.NetConn(ctx, c, websocket.MessageBinary)
	sock := newConnectedSocket(conn, l.serdes)
	l.pendingAccepts <- sock
}

func (l *WebsocketListener) Accept() (*Socket, error) {
	select{
	case sock := <-l.pendingAccepts:
		return sock, nil
	case err := <-l.pendingAcceptErrors:
		return nil, err
	}
}
func (l *WebsocketListener) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	return l.httpServer.Shutdown(ctx)
}
func (l *WebsocketListener) Addr() net.Addr {
	return l.addr
}

// --------------------------------------------------------------------------------
// - TCP
// --------------------------------------------------------------------------------

// TODO - split up dialers and listeners?
type Config struct {
	Url string   // Note: We only use the [scheme]://[host] portion of this
	Serdes Serdes
	TlsConfig *tls.Config
	ReconnectHandler func(*Socket) error // TODO - I think for listeners, this should be when we try to re-listen (ie if listening failed). TODO this should change to be less of a goroutine and more of a (reconnect when an API is called and we are currently disconnected)

	HttpServer *http.Server // TODO - For Websockets only, maybe split up? - Note we have to wrap their Handler with our own handler!
	OriginPatterns []string
}

func (c *Config) Listen() (Listener, error) {
	u, err := url.Parse(c.Url)
	if err != nil {
		return nil, err
	}

	if u.Scheme == "tcp" || u.Scheme == "tcp4" || u.Scheme == "tcp6" || u.Scheme == "unix" || u.Scheme == "unixpacket" {
		listener, err := net.Listen(u.Scheme, u.Host)
		if err != nil {
			return nil, err
		}
		sockListener := &SocketListener{
			listener: listener,
			serdes: c.Serdes,
		}
		return sockListener, nil
	} else if u.Scheme == "wss" {
		listener, err := newWebsocketListener(c)
		if err != nil {
			return nil, err
		}
		return listener, nil
	} else if u.Scheme == "ws" {
		panic("Not implemented yet")
	} else {
		panic("Unsupported network")
	}
}

func (c *Config) Dial() (*Socket, error) {
	sock, err := newSocket(c.Url, c.Serdes, c.TlsConfig)
	if err != nil {
		return nil, err
	}

	// This automatically stops trying to reconnect when Close is called
	go ReconnectLoop(sock, c.ReconnectHandler)

	return sock, nil
}

// --------------------------------------------------------------------------------
// - Websockets
// --------------------------------------------------------------------------------

// Continually attempts to reconnect to the proxy if disconnected. If connected, receives data and sends over the networkChannel
// TODO - It might be nice to not have a reconnect loop and just handle reconnects automatically
func ReconnectLoop(sock *Socket, handler func(*Socket) error) {
	for {
		if sock.Closed.Load() { break } // Exit if the ClientConn has been closed

		err := sock.Dial()
		if err != nil {
			// log.Warn().Err(err).Msg("Client Websocket Dial Failed")
			time.Sleep(5 * time.Second) // TODO - Probably want some random value so everyone isn't reconnecting simultaneously
			continue
		}

		// Start the handler
		err = handler(sock)
		if err != nil {
			// log.Warn().Err(err).Msg("ReconnectLoop handler finished")

			// TODO - Is this a good idea?
			// Try to close the connection one last time
			sock.conn.Close()

			// Set connected to false, because we just closed it
			sock.Connected.Store(false)
		}
		// log.Print("Looping!")
	}

	// Final attempt to cleanup the connection
	sock.Connected.Store(false)
	sock.conn.Close()
	// log.Print("Exiting ClientConn.ReconnectLoop")
}

// This is a wrapper for the client websocket connection
type Socket struct {
	url string            // The URL to connect to
	scheme string         // The scheme of the parsed URL
	host string           // The host of the parsed URL
	tlsConfig *tls.Config // This is the config for tls, nil if not using tls

	encoder Serdes        // The encoder to use for serialization
	conn net.Conn         // The underlying network connection to send and receive on

	// Note: sendMut I think is needed now that I'm using goframe
	sendMut sync.Mutex    // The mutex for multiple threads writing at the same time
	recvMut sync.Mutex    // The mutex for multiple threads reading at the same time
	recvBuf []byte        // The buffer that reads are buffered into

	Closed atomic.Bool    // Used to indicate that the user has requested to close this ClientConn
	Connected atomic.Bool // Used to indicate that the underlying connection is still active
}

// TODO - Combine NewSocket and NewConnectedSocket
func newSocket(network string, encoder Serdes, tlsConfig *tls.Config) (*Socket, error) {
	u, err := url.Parse(network)
	if err != nil {
		return nil, err
	}

	sock := Socket{
		scheme: u.Scheme,
		host: u.Host,
		url: network,
		tlsConfig: tlsConfig,
		encoder: encoder,
		recvBuf: make([]byte, MaxRecvMsgSize),
	}
	return &sock, nil
}

func newConnectedSocket(conn net.Conn, encoder Serdes) *Socket {
	sock := Socket{
		// Create a Framed connection and set it to our connection
		// conn: NewFrameConn(conn),
		conn: conn, // TODO - need to frame when doing TCP and not frame when doing WS. How to handle? Maybe move server abstraction over or something?
		encoder: encoder,
		recvBuf: make([]byte, MaxRecvMsgSize),
	}
	return &sock
}

func (s *Socket) Dial() error {
	// log.Print("Dialing", s.url)
	// Handle websockets
	if s.scheme == "ws" || s.scheme == "wss" {
		// ctx := context.Background()
		// wsConn, _, err := websocket.Dial(ctx, s.url, nil)
		ctx := context.Background()
		wsConn, err := dialWs(ctx, s.url, s.tlsConfig)


		// log.Println("Connection Response:", resp)
		if err != nil { return err }

		// Note: This connection is automagically framed by websockets
		s.conn = websocket.NetConn(ctx, wsConn, websocket.MessageBinary)
		s.Connected.Store(true)
		return nil
	} else if s.scheme == "tcp" {
		conn, err := net.Dial("tcp", s.host)
		if err != nil { return err }

		// Create a Framed connection and set it to our connection
		s.conn = NewFrameConn(conn)
		s.Connected.Store(true)
		return nil
	}

	return fmt.Errorf("Failed to Dial, unknown ClientConn")
}

func (s *Socket) Close() error {
	s.Connected.Store(false)
	s.Closed.Store(true)
	if s.conn != nil {
		err := s.conn.Close()
		return err
	}

	return nil
}

// Sends the message through the connection
func (s *Socket) Send(msg any) error {
	if s.conn == nil {
		return fmt.Errorf("Socket Closed")
	}

	ser, err := s.encoder.Marshal(msg)
	if err != nil {
		return err
	}

	s.sendMut.Lock()
	defer s.sendMut.Unlock()

	// log.Println("ClientSendUpdate:", len(ser))
	_, err = s.conn.Write(ser)
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrNetwork, err)
		return err
	}
	return nil
}
// Reads the next message (blocking) on the connection
func (s *Socket) Recv() (any, error) {
	if s.conn == nil {
		return nil, fmt.Errorf("Socket Closed")
	}

	s.recvMut.Lock()
	defer s.recvMut.Unlock()

	n, err := s.conn.Read(s.recvBuf)
	if err != nil {
		// log.Warn().Err(err).Msg("Failed to receive")
		err = fmt.Errorf("%w: %s", ErrNetwork, err)
		return nil, err
	}
	if n <= 0 { return nil, nil } // There was no message, and no error (likely a keepalive)

	// Note: slice off based on how many bytes we read
	msg, err := s.encoder.Unmarshal(s.recvBuf[:n])
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrSerdes, err)
		return nil, err
	}
	return msg, nil
}
