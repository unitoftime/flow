package net

import (
	"fmt"
	"time"
	"math/rand"
	"errors"
	"io"

	"sync"
	"sync/atomic"
)

// --------------------------------------------------------------------------------
// - Transport based sockets
// --------------------------------------------------------------------------------
type PipeSocket struct {
	dialConfig *DialConfig // The config used for redialing, if nil, then the socket cant reconnect

	encoder Serdes        // The encoder to use for serialization
	decoder Serdes        // The decoder to use for serialization
	pipe Pipe         // The underlying network connection to send and receive on

	// Note: sendMut I think is needed now that I'm using goframe
	sendMut sync.Mutex    // The mutex for multiple threads writing at the same time
	recvMut sync.Mutex    // The mutex for multiple threads reading at the same time
	recvBuf []byte        // The buffer that reads are buffered into

	closed atomic.Bool    // Used to indicate that the user has requested to close this ClientConn
	connected atomic.Bool // Used to indicate that the underlying connection is still active

	redialTimer *time.Timer // Tracks the redial timer
}

func newGlobalSocket() *PipeSocket {
	sock := PipeSocket{
		recvBuf: make([]byte, MaxRecvMsgSize),
	}
	return &sock
}

// Creates a socket spawned by a dialer (as opposed to a listener). These sockets can reconnect.
func newDialSocket(c *DialConfig) *PipeSocket {
	sock := newGlobalSocket()
	sock.dialConfig = c
	sock.encoder = c.Encoder
	sock.decoder = c.Decoder
	return sock
}

// Creates a socket spawned by a listener (as opposed to a dialer). These sockets can't reconnect, the dialer-side socket must reconnect by redialing the listener and getting re-accepted.
func newAcceptedSocket(pipe Pipe, encoder Serdes, decoder Serdes) *PipeSocket {
	sock := newGlobalSocket()
	sock.encoder = encoder
	sock.decoder = decoder

	sock.connectTransport(pipe)

	return sock
}

func (s *PipeSocket) connectTransport(pipe Pipe) {
	if s.connected.Load() {
		panic("Error: This shouldn't happen")
		// return // Skip as we are already connected
	}

	// TODO - close old transport?
	// TODO - ensure that we aren't already connected?
	s.pipe = pipe
	s.connected.Store(true)
}

func (s *PipeSocket) disconnectTransport() error {
	// We have already disconnected the transport
	if !s.connected.Load() {
		return nil
	}
	s.connected.Store(false)

	// Automatically close the socket if it has disconnected and isn't configured to reconnect
	if s.dialConfig == nil {
		s.Close()
	}

	if s.pipe != nil {
		return s.pipe.Close()
	}

	return nil
}

func (s *PipeSocket) Connected() bool {
	return s.connected.Load()
}

func (s *PipeSocket) Closed() bool {
	return s.closed.Load()
}

func (s *PipeSocket) Close() error {
	s.closed.Store(true)
	if s.redialTimer != nil {
		s.redialTimer.Stop()
	}

	s.disconnectTransport()

	return nil
}

// Sends the message through the connection
func (s *PipeSocket) Send(msg any) error {
	if s.Closed() {
		return ErrClosed
	}

	if !s.Connected() {
		return ErrDisconnected
	}

	ser, err := s.encoder.Marshal(msg)
	if err != nil {
		return err
	}

	s.sendMut.Lock()
	defer s.sendMut.Unlock()

	_, err = s.pipe.Write(ser)
	if err != nil {
		s.disconnectTransport()
		err = fmt.Errorf("%w: %s", ErrNetwork, err)
		return err
	}
	return nil
}

// Reads the next message (blocking) on the connection
func (s *PipeSocket) Recv() (any, error) {
	if s.Closed() {
		return nil, ErrClosed
	}

	if !s.Connected() {
		return nil, ErrDisconnected
	}

	s.recvMut.Lock()
	defer s.recvMut.Unlock()

	n, err := s.pipe.Read(s.recvBuf)
	if err != nil {
		s.disconnectTransport()
		// TODO - use new go 1.20 errors.Join() function
		if errors.Is(err, io.EOF) {
			return nil, err
		}
		err = fmt.Errorf("%w: %s", ErrNetwork, err)
		return nil, err
	}
	if n <= 0 { return nil, nil } // There was no message, and no error (likely a keepalive)

	// Note: slice off based on how many bytes we read
	msg, err := s.decoder.Unmarshal(s.recvBuf[:n])
	if err != nil {
		err = fmt.Errorf("%w: %s", ErrSerdes, err)
		return nil, err
	}
	return msg, nil
}

func (s *PipeSocket) Write(dat []byte) (int, error) {
	if s.Closed() {
		return 0, ErrClosed
	}

	if !s.Connected() {
		return 0, ErrDisconnected
	}
	s.sendMut.Lock()
	defer s.sendMut.Unlock()

	// TODO - error if n not big enough
	n, err := s.pipe.Write(dat)
	if err != nil {
		s.disconnectTransport()
		err = fmt.Errorf("%w: %s", ErrNetwork, err)
		return 0, err
	}
	return n, nil
}

func (s *PipeSocket) Read(buf []byte) (int, error) {
	if s.Closed() {
		return 0, ErrClosed
	}

	if !s.Connected() {
		return 0, ErrDisconnected
	}

	s.recvMut.Lock()
	defer s.recvMut.Unlock()

	// TODO - error if buffer wasn't big enough
	n, err := s.pipe.Read(buf)
	if err != nil {
		s.disconnectTransport()
		// TODO - use new go 1.20 errors.Join() function
		if errors.Is(err, io.EOF) {
			return 0, err
		}
		err = fmt.Errorf("%w: %s", ErrNetwork, err)
		return 0, err
	}
	return n, nil
	// if n <= 0 { return nil } // There was no message, and no error (likely a keepalive)
	// return nil
}

// func (s *PipeSocket) Wait() {
// 	for {
// 		if s.connected.Load() {
// 			return
// 		}
// 		// fmt.Println("PipeSocket.Wait()")
// 		time.Sleep(1 * time.Nanosecond)
// 	}
// }

func (s *PipeSocket) redial() {
	if s.dialConfig == nil { return } // If socket cant dial, then skip
	if s.Closed() { return } // If socket is closed, then never reconnect

	go func() {
		// TODO - I'd like this to be more on-demand
		// Trigger the next redial attempt
		defer func() {
			s.redialTimer = time.AfterFunc(1 * time.Second, s.redial)
		}()

		if s.connected.Load() {
			return
		}

		trans, err := s.dialConfig.dialPipe()
		if err != nil {
			return
		}

		// fmt.Println("Socket Reconnected")
		s.connectTransport(trans)
	}()
}

// func (s *PipeSocket) continuallyRedial() {
// 	attempt := 1
// 	const sleepBase = 100 * time.Millisecond // TODO - Tweakable?
// 	const maxSleep = 10 * time.Second // TODO - Tweakable?

// 	sleepDur := sleepBase
// 	for {
// 		if s.Closed() { return } // If socket is closed, then never reconnect

// 		if s.Connected() {
// 			// If socket is already connected, then just sleep
// 			// TODO - for some weird reason, if I dont have a sleep here my game has frame stutters on firefox. I think there's something in here which is making go "give control back" to JS. which makes frames not complete. I haven't been able to pin down exactly where the issue is:
// 			// Relevant: https://github.com/golang/go/issues/27894
// 			// Note: Maybe wasm threads fixes this?
// 			// time.Sleep(redialHack())
// 			// time.Sleep(1 * time.Nanosecond)
// 			time.Sleep(redialHackDur)
// 			continue
// 		}

// 		trans, err := s.dialConfig.dialPipe()
// 		if err != nil {
// 			fmt.Printf("Socket Reconnect attempt %d - Waiting %s. Error: %s\n", attempt, sleepDur, err)
// 			time.Sleep(sleepDur)

// 			// TODO - Tweakable Math?
// 			sleepDur = 2 * sleepDur
// 			if sleepDur > maxSleep {
// 				sleepDur = maxSleep
// 			}

// 			attempt++
// 			continue
// 		}
// 		s.connectTransport(trans)
// 	}
// }

// --------------------------------------------------------------------------------
// - Packetloss code
// --------------------------------------------------------------------------------
type SimSocket struct {
	Socket

	Packetloss float64    // This is the probability that the packet will be lossed for every send/recv
	MinDelay time.Duration // This is the min delay added to every packet sent or recved
	MaxDelay time.Duration // This is the max delay added to every packet sent or recved
	sendDelayErr, recvDelayErr chan error
	recvDelayMsg chan any
	// recvThreadCount int
}

func NewSimSocket(s Socket) *SimSocket {
	return &SimSocket{
		Socket: s,

		sendDelayErr: make(chan error, 10),
		recvDelayMsg: make(chan any, 10),
		recvDelayErr: make(chan error, 10),
	}
}


func (s *SimSocket) Write(buf []byte) (int, error) {
	if rand.Float64() < s.Packetloss {
		// fmt.Println("SEND DROPPING PACKET")
		return 0, nil
	}

	if s.MaxDelay <= 0 {
		return s.Socket.Write(buf)
	}

	// Else send with delay
	go func() {
		r := rand.Float64()
		delay := time.Duration(1_000_000_000 * r * ((s.MaxDelay-s.MinDelay).Seconds())) + s.MinDelay
		// fmt.Println("SendDelay: ", delay)
		time.Sleep(delay)
		_, err := s.Socket.Write(buf)
		if err != nil {
			s.sendDelayErr <- err
		}
	}()

	select {
	case err := <-s.sendDelayErr:
		return 0, err
	default:
		return len(buf), nil // TODO - meh, how to get n out?
	}
}

func (s *SimSocket) Read(buf []byte) (int, error) {
	if rand.Float64() < s.Packetloss {
		// fmt.Println("RECV DROPPING PACKET")
		return 0, nil
	}

	return s.Socket.Read(buf)
}


func (s *SimSocket) Send(msg any) error {
	if rand.Float64() < s.Packetloss {
		// fmt.Println("SEND DROPPING PACKET")
		return nil
	}

	if s.MaxDelay <= 0 {
		return s.Socket.Send(msg)
	}

	// Else send with delay
	go func() {
		r := rand.Float64()
		delay := time.Duration(1_000_000_000 * r * ((s.MaxDelay-s.MinDelay).Seconds())) + s.MinDelay
		// fmt.Println("SendDelay: ", delay)
		time.Sleep(delay)
		err := s.Socket.Send(msg)
		if err != nil {
			s.sendDelayErr <- err
		}
	}()

	select {
	case err := <-s.sendDelayErr:
		return err
	default:
		return nil
	}
}

func (s *SimSocket) Recv() (any, error) {
	if rand.Float64() < s.Packetloss {
		// fmt.Println("RECV DROPPING PACKET")
		return nil, nil
	}

	return s.Socket.Recv()
}

// Code to add delayto recv
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
