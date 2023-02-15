package net

import (
	"errors"
	"sync/atomic"
	"time"
	"net"
	"net/http"
	"crypto/tls"
	"context"

	"nhooyr.io/websocket"
)

type wsPipe struct {
	conn net.Conn
	cancel context.CancelFunc
}
func newWsPipe(wsConn *websocket.Conn) *wsPipe {
	ctx, cancel := context.WithCancel(context.Background())
	conn := websocket.NetConn(ctx, wsConn, websocket.MessageBinary)

	pipe := &wsPipe{
		conn: conn,
		cancel: cancel,
	}
	return pipe
}

func (t *wsPipe) Read(b []byte) (int, error) {
	return t.conn.Read(b)
}

func (t *wsPipe) Write(b []byte) (int, error) {
	return t.conn.Write(b)
}

func (t *wsPipe) Close() error {
	defer t.cancel()
	return t.conn.Close()
}

// Returns a connected socket or fails with an error
func dialWebsocket(c *DialConfig) (*wsPipe, error) {
	ctx, _ := context.WithTimeout(context.Background(), 5 * time.Second)
	conn, err := dialWs(ctx, c.Url, c.TlsConfig)
	if err != nil {
		return nil, err
	}

	return newWsPipe(conn), nil
}

// --------------------------------------------------------------------------------
// - Listener
// --------------------------------------------------------------------------------
type WebsocketListener struct {
	httpServer http.Server
	originPatterns []string
	addr net.Addr
	serdes Serdes
	closed atomic.Bool
	pendingAccepts chan Socket // TODO - should this get buffered?
	pendingAcceptErrors chan error // TODO - should this get buffered?
}
func newWebsocketListener(c *ListenConfig) (*WebsocketListener, error) {
	// TODO - Is tcp always correct here?
	listener, err := tls.Listen("tcp", c.host, c.TlsConfig)
	if err != nil {
		panic(err)
	}

	wsl := &WebsocketListener{
		serdes: c.Serdes,
		addr: listener.Addr(),
		pendingAccepts: make(chan Socket),
		pendingAcceptErrors: make(chan error),
		originPatterns: c.OriginPatterns,
	}

	httpServer := c.HttpServer
	httpServer.Handler = wsl

	go func() {
		for {
			err := httpServer.Serve(listener)
			// ErrServerClosed is returned when shutdown or close is called
			if errors.Is(err, http.ErrServerClosed) {
				return // Just close if the server is shutdown or closed
			} else if wsl.closed.Load() {
				return // Else if closed then just exit
			}

			// TODO - Passing serve errors back through the accept channel. This might be a slightly leaky abstraction. Because these are server errors not really accept errors.
			wsl.pendingAcceptErrors <- err

			time.Sleep(1 * time.Second)
		}
	}()

	return wsl, nil
}

func (l *WebsocketListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	conn, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		OriginPatterns: l.originPatterns,
	})
	if err != nil {
		// Return as an accept error
		l.pendingAcceptErrors <- err
		return
	}

	// Build the socket and push to channel
	pipe := newWsPipe(conn)
	sock := newAcceptedSocket(pipe, l.serdes)
	l.pendingAccepts <- sock
}

func (l *WebsocketListener) Accept() (Socket, error) {
	select{
	case sock := <-l.pendingAccepts:
		return sock, nil
	case err := <-l.pendingAcceptErrors:
		return nil, err
	}
}
func (l *WebsocketListener) Close() error {
	l.closed.Store(true)
	close(l.pendingAccepts)
	close(l.pendingAcceptErrors)

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	return l.httpServer.Shutdown(ctx)
}
func (l *WebsocketListener) Addr() net.Addr {
	return l.addr
}
