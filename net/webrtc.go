package net

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
	"strings"
	"context"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/pion/webrtc/v3"
)

// TODO - Investigate Detaching the datachannel: https://github.com/pion/webrtc/tree/master/examples/data-channels-detach

type RtcSdpMsg struct {
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
	pendingAccepts chan Socket // TODO - should this get buffered?
	pendingAcceptErrors chan error // TODO - should this get buffered?
	closed atomic.Bool
}

func newWebRtcListener(c *ListenConfig) (*WebRtcListener, error) {
	websocketConfig := &ListenConfig{
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
		pendingAccepts: make(chan Socket),
		pendingAcceptErrors: make(chan error),
	}

	go func() {
		for {
			wsConn, err := rtcListener.listener.Accept()
			if err != nil {
				rtcListener.pendingAcceptErrors <- err
				continue
			}

			if rtcListener.closed.Load() {
				return // If closed then just exit
			}

			// Try and negotiate a webrtc connection for the websocket connection
			go rtcListener.attemptWebRtcNegotiation(wsConn)
		}
	}()

	return rtcListener, nil
}

func (l *WebRtcListener) Accept() (Socket, error) {
	select{
	case sock := <-l.pendingAccepts:
		return sock, nil
	case err := <-l.pendingAcceptErrors:
		return nil, err
	}
}
func (l *WebRtcListener) Close() error {
	l.closed.Store(true)
	close(l.pendingAccepts)
	close(l.pendingAcceptErrors)

	return l.listener.Close()
}
func (l *WebRtcListener) Addr() net.Addr {
	return l.listener.Addr()
}

func (l *WebRtcListener) attemptWebRtcNegotiation(wSock Socket) {
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
		l.pendingAcceptErrors <- err
		return
	}

	// When an ICE candidate is available send to the other Pion instance
	// the other Pion instance will add this candidate by calling AddICECandidate
	peerConnection.OnICECandidate(func(c *webrtc.ICECandidate) {
		if c == nil {
			return // Do nothing because the ice candidate was nil for some reason
		}

		candidatesMux.Lock()
		defer candidatesMux.Unlock()

		desc := peerConnection.RemoteDescription()
		if desc == nil {
			pendingCandidates = append(pendingCandidates, c)
		} else {
			candidateMsg := RtcCandidateMsg{c.ToJSON()}
			err := wSock.Send(candidateMsg)
			if err != nil {
				l.pendingAcceptErrors <- fmt.Errorf("OnIceCandidate Send - Possible websocket disconnect: %w", err)
				return
			}
		}
	})

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Print("Listener: Peer Connection State has changed: ", s.String())

		// if s == webrtc.PeerConnectionStateClosed {
		// 	// This means the webrtc was closed by one side. Just close it on the other side
		// 	// Note: because this is the listen side. I don't think we actually need to close this
		// }

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
		conn := NewRtcConn(peerConnection, wSock)
		conn.dataChannel = d

		sock := newAcceptedSocket(conn, l.serdes)
		// Register channel opening handling
		d.OnOpen(func() {
			l.pendingAccepts <- sock
		})

		// Register channel opening handling
		d.OnClose(func() {
			log.Print("Listener: Data channel was closed!!")
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
		anyMsg, err := wSock.Recv()

		// Exit if the socket is closed
		if wSock.Closed() {
			break
		}

		if err != nil {
			// log.Warn().Err(err).Msg("attemptWebRtcNegotiation")
			continue
		}
		if anyMsg == nil { continue }

		switch msg := anyMsg.(type) {
		case RtcCandidateMsg:
			err := peerConnection.AddICECandidate(msg.CandidateInit)
			if err != nil {
				l.pendingAcceptErrors <- fmt.Errorf("RtcCandidateMsg Recv - Failed to add candidate: %w", err)
				return
			}

		case RtcSdpMsg:
			sdp := webrtc.SessionDescription{}
			sdp.Type = msg.Type
			sdp.SDP = msg.SDP

			err := peerConnection.SetRemoteDescription(sdp)
			if err != nil {
				l.pendingAcceptErrors <- fmt.Errorf("RtcSdpMsg Recv - Failed to set remote description: %w", err)
				return
			}

			// Create an answer to send to the other process
			answer, err := peerConnection.CreateAnswer(nil)
			if err != nil {
				l.pendingAcceptErrors <- fmt.Errorf("RtcSdpMsg Recv - Failed to create answer: %w", err)
				return
			}

			answerMessage := RtcSdpMsg{ answer.Type, answer.SDP }
			err = wSock.Send(answerMessage)
			if err != nil {
				l.pendingAcceptErrors <- fmt.Errorf("RtcSdpMsg Recv - Failed to send SDP answer: %w", err)
				return
			}

			// Sets the LocalDescription, and starts our UDP listeners
			err = peerConnection.SetLocalDescription(answer)
			if err != nil {
				l.pendingAcceptErrors <- fmt.Errorf("RtcSdpMsg Recv - Failed to set local SDP: %w", err)
				return
			}

			candidatesMux.Lock()
			for _, c := range pendingCandidates {
				candidateMsg := RtcCandidateMsg{c.ToJSON()}
				err := wSock.Send(candidateMsg)
				if err != nil {
					l.pendingAcceptErrors <- fmt.Errorf("RtcSdpMsg Recv - Failed to send RtcCandidate: %w", err)
					return
				}
			}
			candidatesMux.Unlock()
		}
	}
}

type RtcConn struct {
	peerConn *webrtc.PeerConnection
	dataChannel *webrtc.DataChannel
	websocket Socket
	readChan chan []byte
	errorChan chan error
}
func NewRtcConn(peer *webrtc.PeerConnection, websocket Socket) *RtcConn {
	return &RtcConn{
		peerConn: peer,
		websocket: websocket,
		readChan: make(chan []byte, 100), //TODO! - Sizing
		errorChan: make(chan error), //TODO! - Sizing
	}
}

func (c *RtcConn) Read(b []byte) (int, error) {
	select {
	case err := <-c.errorChan:
		return 0, err // There was some error

		case dat := <- c.readChan:
		if len(dat) > len(b) {
			panic("Read Buffer is too small") // TODO - Fix
		}
		copy(b, dat)
		return len(dat), nil
	}
}

func (c *RtcConn)	Write(b []byte) (int, error) {
	err := c.dataChannel.Send(b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}

func (c *RtcConn) Close() error {
	err1 := c.dataChannel.Close()
	err2 := c.peerConn.Close()
	err3 := c.websocket.Close()
	// var err3 error

	close(c.readChan)
	close(c.errorChan)

	if err1 != nil || err2 != nil || err3 != nil {
		return fmt.Errorf("RtcConn Close Error: datachannel: %s peerconn: %s websocket: %s", err1, err2, err3)
	}
	return nil
}

func dialWebRtc(c *DialConfig) (Pipe, error) {
	ctx, _ := context.WithTimeout(context.Background(), 15 * time.Second)

	websocketConfig := &DialConfig{
		Url: strings.Replace(c.Url, "webrtc", "wss", 1), // TODO! - Not super clean. we want just the url with the schema replaced to be a wss schema
		Serdes: NewRtcUpgradeSerdes(),
		TlsConfig: c.TlsConfig,
	}

	wSock := websocketConfig.Dial()
	for {
		select {
		case <- ctx.Done():
			// The context is over. Exit because we couldn't connect
			wSock.Close()
			return nil, fmt.Errorf("Websocket dial failed")
		default:
		}

		// Check if we've connected and if we have then break
		if wSock.Connected() { break }
		time.Sleep(100 * time.Millisecond)
	}

	wSock.Wait()

	// Offer WebRtc Upgrade
	fmt.Println("offerWebRtcUpgrade")
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
		return nil, err
	}
	fmt.Println("newpeerconn")

	conn := NewRtcConn(peerConnection, wSock)
	connFinish := make(chan bool)
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
			err := conn.websocket.Send(candidateMsg)
			if err != nil {
				conn.errorChan <- err
				return
			}
		}
	})

	fmt.Println("startlisten")
	go func() {
		for {
			anyMsg, err := conn.websocket.Recv()

			// Exit if the socket was closed
			if conn.websocket.Closed() {
				return
			}
			if err != nil {
				log.Warn().Err(err).Msg("dialWebRtc")
				// Because this is an inner goroutine. If there are any issues at all we just want to give up and restart the entire connection process
				return
			}
			if anyMsg == nil { continue }

			switch msg := anyMsg.(type) {
			case RtcCandidateMsg:
				err := peerConnection.AddICECandidate(msg.CandidateInit)
				if err != nil {
					conn.errorChan <- err
					return
				}

			case RtcSdpMsg:
				sdp := webrtc.SessionDescription{}
				sdp.Type = msg.Type
				sdp.SDP = msg.SDP

				err := peerConnection.SetRemoteDescription(sdp)
				if err != nil {
					conn.errorChan <- err
					return
				}

				candidatesMux.Lock()
				defer candidatesMux.Unlock()

				for _, c := range pendingCandidates {
					candidateMsg := RtcCandidateMsg{c.ToJSON()}
					err := conn.websocket.Send(candidateMsg)
					if err != nil {
						conn.errorChan <- err
						return
					}
				}
			}
		}
	}()

	// Create a datachannel with label 'data'
	dataChannel, err := peerConnection.CreateDataChannel("data", nil)
	if err != nil {
		return nil, err
	}

	// Set the handler for Peer connection state
	// This will notify you when the peer has connected/disconnected
	peerConnection.OnConnectionStateChange(func(s webrtc.PeerConnectionState) {
		log.Print("Peer Connection State has changed: ", s.String())

		// if s == webrtc.PeerConnectionStateClosed {
		// 	// This means the webrtc was closed by one side. Just close it on the other side
		// 	sock.Close()
		// }

		if s == webrtc.PeerConnectionStateFailed {
			// Wait until PeerConnection has had no network activity for 30 seconds or another failure. It may be reconnected using an ICE Restart.
			// Use webrtc.PeerConnectionStateDisconnected if you are interested in detecting faster timeout.
			// Note that the PeerConnection may come back from PeerConnectionStateDisconnected.
			log.Error().Msg("Peer Connection has gone to failed")

			conn.errorChan <- fmt.Errorf("Peer Connection has gone to failed")
		} else if s == webrtc.PeerConnectionStateDisconnected {
			conn.errorChan <- fmt.Errorf("Peer Connection has gone to disconnected")
		}
	})

	// Register channel opening handling
	dataChannel.OnOpen(func() {
		conn.dataChannel = dataChannel
		// sock.conn = conn
		connFinish <- true
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
		return nil, err
	}
	fmt.Println("CreateOffer")

	// Sets the LocalDescription, and starts our UDP listeners
	// Note: this will start the gathering of ICE candidates
	err = peerConnection.SetLocalDescription(offer)
	if err != nil {
		return nil, err
	}
	fmt.Println("SetLocalDesc")

	offerMessage := RtcSdpMsg{ offer.Type, offer.SDP }
	err = conn.websocket.Send(offerMessage)
	if err != nil {
		return nil, err
	}
	fmt.Println("websocket send")

	// Wait until the webrtc connection is finished getting setup
	select {
	case err := <-conn.errorChan:
		return nil, err // There was an error in setup
	case <-connFinish:
		// Socket finished getting setup
		return conn, nil
	}
}
