// +build js,wasm

package net

import (
	"context"
	"crypto/tls"

	"nhooyr.io/websocket"
)

func dialWs(ctx context.Context, url string, tlsConfig *tls.Config) (*websocket.Conn, error) {
	wsConn, _, err := websocket.Dial(ctx, url, nil)
	return wsConn, err
}