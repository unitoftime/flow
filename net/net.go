package net

import (
	"fmt"
	"errors"
	"time"
	// "math/rand"
	"net"
	"net/url"
	"net/http"
	"crypto/tls"
	"sync"
	"sync/atomic"
)

// TODO - Ensure sent messages remain under this
// Calculation: 1460 Byte = 1500 Byte - 20 Byte IP Header - 20 Byte TCP Header
// const MaxMsgSize = 1460 // bytes
const MaxRecvMsgSize = 4 * 1024 // 8 KB // TODO - this is arbitrary

var ErrSerdes = errors.New("serdes errror")
var ErrNetwork = errors.New("network error")
var ErrDisconnected = errors.New("socket disconnected")
var ErrClosed = errors.New("socket closed") // Indicates that the socket is closed. Currently if you get this error then it means the socket will never receive or send again!

// This is the interface used to marshal and unmarshal messages over the network.
// Hardest part is when you have interfaces in messages, you'll likely need custom serializers for that
type Serdes interface {
	Marshal(v any) ([]byte, error)
	Unmarshal(dat []byte) (any, error)
}

type Listener interface {
	// Accept waits for and returns the next connection to the listener.
	Accept() (Socket, error)

	// Close closes the listener.
	// Any blocked Accept operations will be unblocked and return errors.
	Close() error

	// Addr returns the listener's network address.
	Addr() net.Addr
}

type Transport interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

type Socket interface {
	// TODO - SetReadDeadline and SetWriteDeadline could be nice to have!

	Send(any) error
	Recv() (any, error)
	Close() error

	Connected() bool
	Closed() bool
	Wait() // Wait for the connection to stabalize
}

type SocketListener struct {
	listener net.Listener
	serdes Serdes
}
func (l *SocketListener) Accept() (Socket, error) {
	c, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}

	framedConn := NewFrameConn(c)
	return newAcceptedSocket(framedConn, l.serdes), nil
}
func (l *SocketListener) Close() error {
	return l.listener.Close()
}
func (l *SocketListener) Addr() net.Addr {
	return l.listener.Addr()
}

// --------------------------------------------------------------------------------
// - Config
// --------------------------------------------------------------------------------

// For dialing a socket
type DialConfig struct {
	Url string   // Note: We only use the [scheme]://[host] portion of this
	Serdes Serdes
	TlsConfig *tls.Config

	// These are generated based on the upper config
	scheme string
	host string
}

// Returns a created socket which may not be connected, but will be actively trying to connect
func (c *DialConfig) Dial() Socket {
	// Parse the config
	u, err := url.Parse(c.Url)
	if err != nil {
		// TODO - wrap this up in the creation of the dialconfig
		panic(fmt.Sprintf("URL Parsing Error:", err))
	}
	c.scheme = u.Scheme
	c.host = u.Host

	sock := newDialSocket(c)

	go sock.continuallyRedial()

	return sock
}


func (c *DialConfig) dialTransport() (Transport, error) {
	// Handle websockets
	if c.scheme == "ws" || c.scheme == "wss" {
		return dialWebsocket(c)
	} else if c.scheme == "tcp" {
		return dialTcpSocket(c)
	} else if c.scheme == "webrtc" {
		return dialWebRtc(c)
	}

	return nil, fmt.Errorf("Failed to Dial, unknown scheme")
}

// For listening for sockets
type ListenConfig struct {
	Url string   // Note: We only use the [scheme]://[host] portion of this
	Serdes Serdes
	TlsConfig *tls.Config

	HttpServer *http.Server // TODO - For Websockets only, maybe split up? - Note we have to wrap their Handler with our own handler!
	OriginPatterns []string
}

func (c *ListenConfig) Listen() (Listener, error) {
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
	} else if u.Scheme == "webrtc" {
		listener, err := newWebRtcListener(c)
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

// --------------------------------------------------------------------------------
// - Transport based sockets
// --------------------------------------------------------------------------------
type TransportSocket struct {
	dialConfig *DialConfig

	encoder Serdes        // The encoder to use for serialization
	conn Transport         // The underlying network connection to send and receive on

	// Note: sendMut I think is needed now that I'm using goframe
	sendMut sync.Mutex    // The mutex for multiple threads writing at the same time
	recvMut sync.Mutex    // The mutex for multiple threads reading at the same time
	recvBuf []byte        // The buffer that reads are buffered into

	closed atomic.Bool    // Used to indicate that the user has requested to close this ClientConn
	connected atomic.Bool // Used to indicate that the underlying connection is still active

	// Packetloss float64    // This is the probability that the packet will be lossed for every send/recv
	// MinDelay time.Duration // This is the min delay added to every packet sent or recved
	// MaxDelay time.Duration // This is the max delay added to every packet sent or recved
	// sendDelayErr, recvDelayErr chan error
	// recvDelayMsg chan any
	// recvThreadCount int

}

func dialTcpSocket(c *DialConfig) (Transport, error) {
	conn, err := net.Dial("tcp", c.host)
	if err != nil {
		return nil, err
	}

	// Create a Framed connection and set it to our connection
	framedConn := NewFrameConn(conn)
	return framedConn, nil
}

func newGlobalSocket() *TransportSocket {
	sock := TransportSocket{
		recvBuf: make([]byte, MaxRecvMsgSize),

		// recvDelayMsg: make(chan any, 10),
		// recvDelayErr: make(chan error, 10),
	}
	return &sock
}

func newDialSocket(c *DialConfig) *TransportSocket {
	sock := newGlobalSocket()
	sock.dialConfig = c
	sock.encoder = c.Serdes
	return sock
}

func newAcceptedSocket(conn Transport, encoder Serdes) *TransportSocket {
	sock := newGlobalSocket()
	sock.encoder = encoder

	sock.connectTransport(conn)

	return sock
}

func (s *TransportSocket) connectTransport(transport Transport) {
	if s.connected.Load() {
		panic("Error: This shouldn't happen")
		// return // Skip as we are already connected
	}

	// TODO - close old transport?
	// TODO - ensure that we aren't already connected?
	s.conn = transport
	s.connected.Store(true)
}

func (s *TransportSocket) disconnectTransport() error {
	// We have already disconnected the transport
	if !s.connected.Load() {
		return nil
	}

	var err error
	if s.conn != nil {
		err = s.conn.Close()
	}

	s.conn = nil
	s.connected.Store(false)

	return err
}

func (s *TransportSocket) Connected() bool {
	return s.connected.Load()
}

func (s *TransportSocket) Closed() bool {
	return s.closed.Load()
}

func (s *TransportSocket) Close() error {
	s.disconnectTransport()

	s.closed.Store(true)

	return nil
}

// Sends the message through the connection
func (s *TransportSocket) Send(msg any) error {
	if s.Closed() {
		return ErrClosed
	}

	if !s.Connected() {
		return ErrDisconnected
	}

	return s.send(msg)

	// TODO - add back in some wrapper class I think
	// if s.MaxDelay <= 0 {
	// 	return s.send(msg)
	// }

	// // Else send with delay
	// go func() {
	// 	r := rand.Float64()
	// 	delay := time.Duration(1_000_000_000 * r * ((s.MaxDelay-s.MinDelay).Seconds())) + s.MinDelay
	// 	// fmt.Println("SendDelay: ", delay)
	// 	time.Sleep(delay)
	// 	err := s.send(msg)
	// 	if err != nil {
	// 		s.sendDelayErr <- err
	// 	}
	// }()

	// select {
	// case err := <-s.sendDelayErr:
	// 	return err
	// default:
	// 	return nil
	// }
}

func (s *TransportSocket) Recv() (any, error) {
	if s.Closed() {
		return nil, ErrClosed
	}

	if !s.Connected() {
		return nil, ErrDisconnected
	}

	return s.recv()

	// TODO - fix this
	// if s.MaxDelay <= 0 {
	// 	return s.recv()
	// }

	// for {
	// 	if s.recvThreadCount > 100 {
	// 		break
	// 	}
	// 	s.recvThreadCount++ // TODO - not thread safe
	// 	go func() {
	// 		msg, err := s.recv()

	// 		r := rand.Float64()
	// 		delay := time.Duration(1_000_000_000 * r * ((s.MaxDelay-s.MinDelay).Seconds())) + s.MinDelay
	// 		fmt.Println("RecvDelay: ", delay)
	// 		time.Sleep(delay)

	// 		s.recvThreadCount--
	// 		if err != nil {
	// 			s.recvDelayErr <- err
	// 		} else {
	// 			fmt.Println("Recv: ", msg, err)
	// 			s.recvDelayMsg <- msg
	// 		}
	// 	}()
	// }

	// select {
	// case err := <-s.recvDelayErr:
	// 	return nil, err
	// default:
	// 	msg := <-s.recvDelayMsg
	// 	fmt.Println("RETURNING")
	// 	return msg, nil
	// }
}

func (s *TransportSocket) send(msg any) error {
	// TODO - I'd prefer this to never be nil!
	if s.conn == nil {
		return ErrNetwork
	}

	// if rand.Float64() < s.Packetloss {
	// 	return nil
	// }

	ser, err := s.encoder.Marshal(msg)
	if err != nil {
		return err
	}

	s.sendMut.Lock()
	defer s.sendMut.Unlock()

	fmt.Println("AttemptedSend")
	_, err = s.conn.Write(ser)
	if err != nil {
		s.disconnectTransport()
		err = fmt.Errorf("%w: %s", ErrNetwork, err)
		return err
	}
	return nil
}

// Reads the next message (blocking) on the connection
func (s *TransportSocket) recv() (any, error) {
	// TODO - I'd prefer this to never be nil!
	if s.conn == nil {
		return nil, ErrNetwork
	}

	s.recvMut.Lock()
	defer s.recvMut.Unlock()

	n, err := s.conn.Read(s.recvBuf)
	if err != nil {
		s.disconnectTransport()
		err = fmt.Errorf("%w: %s", ErrNetwork, err)
		return nil, err
	}
	if n <= 0 { return nil, nil } // There was no message, and no error (likely a keepalive)

	// if rand.Float64() < s.Packetloss {
	// 	return nil, nil
	// }

	// Note: slice off based on how many bytes we read
	msg, err := s.encoder.Unmarshal(s.recvBuf[:n])
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrSerdes, err)
		return nil, err
	}
	return msg, nil
}

func (s *TransportSocket) Wait() {
	for {
		if s.connected.Load() {
			return
		}
		time.Sleep(100 * time.Millisecond)
	}
}

func (s *TransportSocket) continuallyRedial() {
	attempt := 1
	sleepDur := 100 * time.Millisecond // TODO - Tweakable?
	maxSleep := 10 * time.Second // TODO - Tweakable?
	for {
		if s.closed.Load() { return } // If socket is closed, then never reconnect

		fmt.Println("Redial Loop")
		if s.connected.Load() {
			fmt.Println("Already Connected")
			// If socket is already connected, then just sleep
			time.Sleep(sleepDur) // TODO - I feel like I'd prefer this to be some better sync mechanism, but I'm not sure what to use
			continue
		}

		fmt.Println("Attempting Redial")
		trans, err := s.dialConfig.dialTransport()
		if err != nil {
			fmt.Printf("Socket Reconnect attempt %d - Waiting %s\n", attempt, sleepDur)
			fmt.Println(err)
			attempt++
			time.Sleep(sleepDur)
			sleepDur = 2 * sleepDur // TODO - Tweakable?
			if sleepDur > maxSleep {
				sleepDur = maxSleep
			}
			continue
		}

		fmt.Println("Socket Reconnected")
		s.connectTransport(trans)
	}
}
