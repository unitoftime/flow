package net

import (
	"net"
	"sync"
	"time"
	"errors"
	"strings"
	// "context"

	"github.com/rs/zerolog/log"

	"github.com/pion/webrtc/v3"
)

// TODO - Investigate Detaching the datachannel: https://github.com/pion/webrtc/tree/master/examples/data-channels-detach

type RtcSdpMsg struct {
	// SessionDescription webrtc.SessionDescription
	Type webrtc.SDPType
	SDP string
}

type RtcCandidateMsg struct {
	CandidateInit webrtc.ICECandidateInit
}

type RtcUpgradeSerdes struct {
	union *UnionBuilder
}
func NewRtcUpgradeSerdes() *RtcUpgradeSerdes {
	return &RtcUpgradeSerdes{NewUnion(RtcSdpMsg{}, RtcCandidateMsg{})}
}
func (s *RtcUpgradeSerdes) Marshal(v any) ([]byte, error) {
	return s.union.Serialize(v)
}
func (s *RtcUpgradeSerdes) Unmarshal(dat []byte) (any, error) {
		return s.union.Deserialize(dat)
}

type WebRtcListener struct {
	listener Listener
	serdes Serdes
	pendingAccepts chan *Socket // TODO - should this get buffered?
	pendingAcceptErrors chan error // TODO - should this get buffered?
}

func newWebRtcListener(c *Config) (*WebRtcListener, error) {
	websocketConfig := &Config{
		Url: strings.Replace(c.Url, "webrtc", "wss", 1), // TODO! - Not super clean. we want just the url with the schema replaced to be a wss schema
		Serdes: NewRtcUpgradeSerdes(),
		TlsConfig: c.TlsConfig,
		HttpServer: c.HttpServer,
		OriginPatterns: c.OriginPatterns,
	}
	wsl, err := websocketConfig.Listen()
	if err != nil {
		return nil, err
	}

	rtcListener := &WebRtcListener{
		listener: wsl,
		serdes: c.Serdes,
		pendingAccepts: make(chan *Socket),
		pendingAcceptErrors: make(chan error),
	}

	// TODO - some way to cancel?
	go func() {
		for {
			wsConn, err := rtcListener.listener.Accept()
			if err != nil {
				rtcListener.pendingAcceptErrors <- err
				continue
			}

			// Try and negotiate a webrtc connection
			go rtcListener.attemptWebRtcNegotiation(wsConn)
		}
	}()

	return rtcListener, nil
}

func (l *WebRtcListener) Accept() (*Socket, error) {
	select{
	case sock := <-l.pendingAccepts:
		return sock, nil
	case err := <-l.pendingAcceptErrors:
		return nil, err
	}
}
func (l *WebRtcListener) Close() error {
	return l.listener.Close()
}
func (l *WebRtcListener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *WebRtcListener) attemptWebRtcNegotiation(wsConn *Socket) {
	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)
	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			// {
			// 	URLs: []string{"stun:stun.l.google.com:19302"},
			// },
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	// TODO - when Should I close?
	// defer func() {
	// 	if err := peerConnection.Close(); err != nil {
	// 		fmt.Printf("cannot close peerConnection: %v\n", err)
	// 	}
	// }()

	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else {
			candidateMsg := RtcCandidateMsg{c.ToJSON()}
			onICECandidateErr := wsConn.Send(candidateMsg)
			if onICECandidateErr != nil {
				panic(onICECandidateErr)
			}
		}
	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		// log.Debug().Msg("Peer Connection State has changed: ", s.String())
		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			log.Error().Msg("Peer Connection has gone to failed")
			// TODO - Do some cancellation
		}
	})

	// Register data channel creation handling
	peerConnection.OnDataChannel(func(d *webrtc.DataChannel) {
		conn := NewRtcConn(peerConnection, wsConn)
		conn.dataChannel = d

		// Register channel opening handling
		d.OnOpen(func() {
			sock := newConnectedSocket(conn, l.serdes)
			l.pendingAccepts <- sock
		})

		// Register text message handling
		d.OnMessage(func(msg webrtc.DataChannelMessage) {
			// log.Print("Server: Received Msg from DataChannel", len(msg.Data))
			if msg.IsString {
				log.Warn().Msg("DataChannel OnMessage: Received string message, skipping")
				return
			}
			conn.readChan <- msg.Data
		})
	})

	for {
		anyMsg, err := wsConn.Recv()

		if errors.Is(err, ErrNetwork) {
			// Handle errors where we should stop (ie connection closed or something)
			log.Warn().Err(err).Msg("attemptWebRtcNegotiation: NetworkErr")
			return
		} else if errors.Is(err, ErrSerdes) {
			// Handle errors where we should continue (ie serialization)
			log.Error().Err(err).Msg("attemptWebRtcNegotiation:  SerdesErr")
			continue
		}
		if anyMsg == nil { continue }

		switch msg := anyMsg.(type) {
		case RtcCandidateMsg:
			candidateErr := peerConnection.AddICECandidate(msg.CandidateInit)
			if candidateErr != nil {
				panic(candidateErr)
			}

		case RtcSdpMsg:
			sdp := webrtc.SessionDescription{}
			sdp.Type = msg.Type
			sdp.SDP = msg.SDP

			if err := peerConnection.SetRemoteDescription(sdp); err != nil {
				panic(err)
			}

			// Create an answer to send to the other process
			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				panic(err)
			}

			answerMessage := RtcSdpMsg{ answer.Type, answer.SDP }
			err = wsConn.Send(answerMessage)
			if err != nil {
				panic(err)
			}

			// Sets the LocalDescription, and starts our UDP listeners
			err = peerConnection.SetLocalDescription(answer)
			if err != nil {
				panic(err)
			}

			candidatesMux.Lock()
			for _, c := range pendingCandidates {
				candidateMsg := RtcCandidateMsg{c.ToJSON()}
				onICECandidateErr := wsConn.Send(candidateMsg)
				if onICECandidateErr != nil {
					panic(onICECandidateErr)
				}
			}
			candidatesMux.Unlock()
		}
	}
}

type RtcConn struct {
	peerConn *webrtc.PeerConnection
	dataChannel *webrtc.DataChannel
	websocket *Socket
	readChan chan []byte
}
func NewRtcConn(peer *webrtc.PeerConnection, websocket *Socket) *RtcConn {
	return &RtcConn{
		peerConn: peer,
		websocket: websocket,
		readChan: make(chan []byte, 100), //TODO! - Sizing
	}
}

func (c *RtcConn) Read(b []byte) (int, error) {
	dat := <- c.readChan
	if len(dat) > len(b) {
		panic("Read Buffer is too small") // TODO - Fix
	}
	copy(b, dat)
	return len(dat), nil
}

func (c *RtcConn)	Write(b []byte) (int, error) {
	err := c.dataChannel.Send(b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *RtcConn) Close() error {
	err := c.dataChannel.Close()
	if err != nil {
		log.Error().Err(err).Msg("RtcConn: Error Closing WebRtc DataChannel")
	}
	err2 := c.peerConn.Close()
	if err2 != nil {
		log.Error().Err(err2).Msg("RtcConn: Error Closing WebRtc Peer Connection")
	}
	err3 := c.websocket.Close()
	if err3 != nil {
		log.Error().Err(err3).Msg("RtcConn: Error Closing Websocket Connection")
	}

	// TODO! - I'm not sure the best way to do this. Maybe wrap these if not nil?
	if err != nil { return err }
	if err2 != nil { return err2 }
	if err3 != nil { return err3 }

	return nil
}
// TODO - Rethink -> How do these affect the webrtc connection
func (c *RtcConn) LocalAddr() net.Addr {
	return c.websocket.conn.LocalAddr()
}
func (c *RtcConn) RemoteAddr() net.Addr {
	return c.websocket.conn.RemoteAddr()
}
func (c *RtcConn) SetDeadline(t time.Time) error {
	return c.websocket.conn.SetDeadline(t)
}
func (c *RtcConn) SetReadDeadline(t time.Time) error {
	return c.websocket.conn.SetReadDeadline(t)
}
func (c *RtcConn) SetWriteDeadline(t time.Time) error {
	return c.websocket.conn.SetWriteDeadline(t)
}

func dialWebRtc(sock *Socket) error {
	websocketConfig := &Config{
		Url: strings.Replace(sock.url, "webrtc", "wss", 1), // TODO! - Not super clean. we want just the url with the schema replaced to be a wss schema
		Serdes: NewRtcUpgradeSerdes(),
		TlsConfig: sock.tlsConfig,
		ReconnectHandler: func(wSock *Socket) error {
			// TODO! - On reconnects I kind of want to offer an Rtc Upgrade, but ONLY if I haven't already upgraded
			for {
				err := offerWebRtcUpgrade(wSock, sock)
				if err != nil {
					log.Error().Err(err).Msg("Client Reconnect Error")
				}
				time.Sleep(1 * time.Second)
			}
		},
	}

	_, err := websocketConfig.Dial()
	if err != nil {
		return err
	}

	// TODO! - Wait for sock.conn to get set (ie the webrtc upgrade completes)
	for {
		if sock.conn != nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	return nil
}

func offerWebRtcUpgrade(wSock *Socket, sock *Socket) error {
	var candidatesMux sync.Mutex
	pendingCandidates := make([]*webrtc.ICECandidate, 0)

	config := webrtc.Configuration{
		ICEServers: []webrtc.ICEServer{
			// {
			// 	URLs: []string{"stun:stun.l.google.com:19302"},
			// },
		},
	}

	peerConnection, err := webrtc.NewPeerConnection(config)
	if err != nil {
		panic(err)
	}
	// TODO - when should I close the peerConnection?
	// defer func() {
	// 	if cErr := peerConnection.Close(); cErr != nil {
	// 		fmt.Printf("cannot close peerConnection: %v\n", cErr)
	// 	}
	// }()

	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else {
			candidateMsg := RtcCandidateMsg{c.ToJSON()}
			onICECandidateErr := wSock.Send(candidateMsg)
			if onICECandidateErr != nil {
				panic(onICECandidateErr)
			}
		}
	})

	go func() {
		for {
			anyMsg, err := wSock.Recv()

			if errors.Is(err, ErrNetwork) {
				log.Warn().Err(err).Msg("dialWebRtc: NetworkErr")
				return
			} else if errors.Is(err, ErrSerdes) {
				log.Error().Err(err).Msg("dialWebRtc:  SerdesErr")
				continue
			}
			if anyMsg == nil { continue }

			switch msg := anyMsg.(type) {
			case RtcCandidateMsg:
				candidateErr := peerConnection.AddICECandidate(msg.CandidateInit)
				if candidateErr != nil {
					panic(candidateErr)
				}

			case RtcSdpMsg:
				sdp := webrtc.SessionDescription{}
				sdp.Type = msg.Type
				sdp.SDP = msg.SDP

				if sdpErr := peerConnection.SetRemoteDescription(sdp); sdpErr != nil {
					panic(sdpErr)
				}

				candidatesMux.Lock()
				defer candidatesMux.Unlock()

				for _, c := range pendingCandidates {
					candidateMsg := RtcCandidateMsg{c.ToJSON()}
					onICECandidateErr := wSock.Send(candidateMsg)
					if onICECandidateErr != nil {
						panic(onICECandidateErr)
					}
				}
			}
		}
	}()

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		panic(err)
	}

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Print("Peer Connection State has changed: ", s.String())

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			log.Error().Msg("Peer Connection has gone to failed")
			// TODO - do some cancellation
		}
	})

	conn := NewRtcConn(peerConnection, wSock)

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		conn.dataChannel = dataChannel
		sock.conn = conn
	})

	// Register text message handling
	dataChannel.OnMessage(func(msg webrtc.DataChannelMessage) {
			// log.Print("Client: Received Msg from DataChannel", len(msg.Data))
			if msg.IsString {
				log.Print("DataChannel OnMessage: Received string message, skipping")
				return
			}
			conn.readChan <- msg.Data
	})

	// Create an offer to send to the other process
	offer, err := peerConnection.CreateOffer(nil)
	if err != nil {
		panic(err)
	}

	// Sets the LocalDescription, and starts our UDP listeners
	// Note: this will start the gathering of ICE candidates
	if err = peerConnection.SetLocalDescription(offer); err != nil {
		panic(err)
	}

	offerMessage := RtcSdpMsg{ offer.Type, offer.SDP }
	err = wSock.Send(offerMessage)
	if err != nil {
		panic(err)
	}

	select{} // TODO! - Do Cancellation: I guess we don't want to loop we just want to wait here until we need to reconnect? Not sure how to exit
}
