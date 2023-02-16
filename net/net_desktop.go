// +build !js

package net

import (
	"time"
	"context"
	"net/http"
	"crypto/tls"

	"nhooyr.io/websocket"
)

func dialWs(ctx context.Context, url string, tlsConfig *tls.Config) (*websocket.Conn, error) {
	wsConn, _, err := websocket.Dial(ctx, url, &websocket.DialOptions{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: tlsConfig,
			},
		},
	})
	return wsConn, err
}

const redialHackDur = 1 * time.Second
