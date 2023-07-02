package envoy

import (
	"fmt"
	"time"
	"errors"
	"math/rand"
	"reflect"
	"sync"

	"github.com/unitoftime/flow/net"
	"github.com/unitoftime/flow/cod"
)

// Has
// 1. Bidirectional RPCs
// 2. Fire-and-forget style RPCs (ie just send it and don't block, like a msg)
// 3. Easy setup and management

// Wants
// 4. Different reliability levels
// 5. Automatic retries of messages?
// 6. Message batching?
// TODO: Maybe eventually migrate over to the auto generated serialization code

var ErrTimeout = errors.New("timeout reached")

const (
	wireTypeRequest uint8 = iota
	wireTypeResponse
	wireTypeMessage
)

// // Internal serialization
// var rpcSerdes serdes
// func init() {
// 	rpcSerdes = serdes{
// 		// union: net.NewUnion(
// 		// 	Request{},
// 		// 	Response{},
// 		// 	Message{},
// 		// ),
// 	}
// }

// type serdes struct {
// 	// union *net.UnionBuilder
// }
// func (s *serdes) Marshal(v any) ([]byte, error) {
// 	// return s.union.Serialize(v)

// 	buf := cod.NewBuffer(1024)

// 	switch msgDat := m.Data.(type) {
// 	case Request:
// 		buf.WriteUint8(headerRequest)
// 		// encodeRequest(buf, msgDat)
// 		msgDat.EncodeCod(buf)
// 	case Response:
// 		buf.WriteUint8(headerResponse)
// 		msgDat.EncodeCod(buf)
// 	case Message:
// 		buf.WriteUint8(headerMessage)
// 		msgDat.EncodeCod(buf)
// 	default:
// 		panic(fmt.Errorf("Envoy internal: unknown type error: %T\n", msgDat))
// 	}

// 	return buf.Bytes(), nil
// }
// func (s *serdes) Unmarshal(dat []byte) (any, error) {
// 	// return s.union.Deserialize(dat)

// 	header := buf.ReadUint8()

// 	var err error

// 	switch header {
// 	case headerRequest:
// 		req := Request{}
// 		err = req.DecodeCod(buf)
// 		// err = decodeRequest(buf)
// 	case headerResponse:
// 		res := Response{}
// 		err = res.DecodeCod(buf)
// 		// m.Data, err = decodeResponse(buf)
// 	case headerMessage:
// 		msg := Message{}
// 		err = msg.DecodeCod(buf)
// 		// m.Data, err = decodeMessage(buf)
// 	default:
// 		err = fmt.Errorf("unknown header id: %d", header)
// 	}
// 	return err
// }


// Internal messages
type Request struct {
	Id uint32 // Tracks the request Id
	Data []byte
}

type Response struct {
	Id uint32 // Tracks the request Id
	Data []byte
}

type Message struct {
	Data []byte
}

// Internal interfaces
type rpcClient interface {
	doRpc(any, time.Duration) (any, error)
	doMsg(any) error
}

type MsgDefinition interface {
	MsgType() any
}

type RpcDefinition interface {
	ReqType() any
	RespType() any
}

type clientSetter interface {
	setClient(rpcClient)
}

type RpcHandler interface {
	Handler() (reflect.Type, HandlerFunc)
}

type MsgHandler interface {
	Handler() (reflect.Type, MessageHandlerFunc)
}

// Message Definition
type MsgDef[A any] struct {
	handler MessageHandlerFunc
	client rpcClient
}

func (d *MsgDef[A]) setClient(client rpcClient) {
	d.client = client
}

func (d MsgDef[A]) Handler() (reflect.Type, MessageHandlerFunc) {
	var a A
	return reflect.TypeOf(a), d.handler
}

func (d *MsgDef[A]) Register(handler func(A)) {
	d.handler = makeMsgHandler(handler)
}

func (d MsgDef[A]) MsgType() any {
	var a A
	return a
}

func (d MsgDef[A]) Send(msg A) error {
	return d.client.doMsg(msg)
}

// RPC Definition
type RpcDef[Req, Resp any] struct {
	handler HandlerFunc
	client rpcClient
}

func (d *RpcDef[Req, Resp]) setClient(client rpcClient) {
	d.client = client
}

func (d RpcDef[Req, Resp]) ReqType() any {
	var req Req
	return req
}

func (d RpcDef[Req, Resp]) RespType() any {
	var resp Resp
	return resp
}

func (d *RpcDef[Req, Resp]) Register(handler func(Req) Resp) {
	d.handler = makeRpcHandler(handler)
}

func (d RpcDef[Req, Resp]) Call(req Req) (Resp, error) {
	var resp Resp
	anyResp, err := d.client.doRpc(req, 15 * time.Second) // TODO: configurable
	if err != nil { return resp, err }
	resp, ok := anyResp.(Resp)
	if !ok { panic("Mismatched type!") }
	return resp, nil
}

func (d RpcDef[Req, Resp]) Handler() (reflect.Type, HandlerFunc) {
	var req Req
	return reflect.TypeOf(req), d.handler
}

// Interface Definition
type InterfaceDef[S, C any] struct {
	Service S
	Client C
	serviceApi serviceDef
	clientApi serviceDef
}

func NewInterfaceDef[S, C any]() InterfaceDef[S, C] {
	var serviceApi S
	var clientApi C
	return InterfaceDef[S, C]{
		serviceApi: makeServiceDef(serviceApi),
		clientApi: makeServiceDef(clientApi),
	}
}

type Config struct {
	MaxRecvPacketSize int // Default 65K
}

// Returns the client side of the interface
func (d InterfaceDef[S, C]) NewClient(config *Config) *Client[C, S] {
	// Note: The C and S are reversed because we call the service and serve the client
	client := newClient[C, S](d.clientApi, d.serviceApi)

	if config != nil {
		client.maxRecvPacketSize = config.MaxRecvPacketSize
	}

	return client
}

// Returns the server side of the interface
func (d InterfaceDef[S, C]) NewServer(config *Config) *Client[S, C] {
	client := newClient[S, C](d.serviceApi, d.clientApi)

	if config != nil {
		client.maxRecvPacketSize = config.MaxRecvPacketSize
	}

	return client
}


// TODO - I should make an interface to better capture the fact that this is just for serialization
// TODO - could I make servicedef generic on the interface type. Then when I register handlers I just pass in a struct which implements the interface?

// TODO - I think I'd prefer this to be based on method name and not based on input argument type
// TODO - reordering the definition, or switching between a message and an RPC will break api compatibility
type serviceDef struct {
	Requests, Responses *cod.UnionBuilder
}

func makeServiceDef(def any) serviceDef {
	ty := reflect.TypeOf(def)
	numField := ty.NumField()

	requests := make([]any, 0)
	responses := make([]any, 0)
	for i := 0; i < numField; i++ {
		field := ty.Field(i)

		fieldAny := reflect.New(field.Type).Elem().Interface()

		switch rpcDef := fieldAny.(type) {
		case RpcDefinition:
			reqStruct := rpcDef.ReqType()
			requests = append(requests, reqStruct)

			respStruct := rpcDef.RespType()
			responses = append(responses, respStruct)

		case MsgDefinition:
			msgStruct := rpcDef.MsgType()
			requests = append(requests, msgStruct)

		default:
			panic("Error: Fields must all either be RpcDef or MsgDef")
		}
	}

	return serviceDef{
		Requests: cod.NewUnion(requests...),
		Responses: cod.NewUnion(responses...),
	}
}

type Client[S, C any] struct {
	sock net.Socket

	Handler S
	Call C

	serviceDef, clientDef serviceDef

	handlers map[reflect.Type]HandlerFunc
	messageHandlers map[reflect.Type]MessageHandlerFunc

	reqLock sync.Mutex
	activeCalls map[uint32]chan any

	rng *rand.Rand
	maxRecvPacketSize int
}

func newClient[S, C any](serviceDef, clientDef serviceDef) *Client[S, C] {
	rngSrc := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rngSrc) // TODO - push this up to the client?

	client := &Client[S, C]{
		serviceDef: serviceDef,
		clientDef: clientDef,

		handlers: make(map[reflect.Type]HandlerFunc),
		messageHandlers: make(map[reflect.Type]MessageHandlerFunc),
		activeCalls: make(map[uint32]chan any),

		rng: rng,
		maxRecvPacketSize: 1000 * 65, // TODO!!!! - hardcoded: 65K Bytes is approximately the largest UDP/WebRTC packet you can send.
	}

	return client
}

func (c *Client[S, C]) Connect(sock net.Socket) {
	c.registerHandlers(c.Handler)
	c.registerCallers(&c.Call)

	c.sock = sock
	c.start() // This doesn't block
}

func (c *Client[S, C]) Close() error {
	// TODO - what to do about active calls and handlers?
	return c.sock.Close()
}

// TODO - get rid of this
func (c *Client[S, C]) Closed() bool {
	return c.sock.Closed()
}

var sendBufPool = sync.Pool{
	New: func() any {
		// The Pool's New function should generally only return pointer
		// types, since a pointer can be put into the return interface
		// value without an allocation:
		return cod.NewBuffer(8 * 1024) // TODO!!!! - hardcoded
	},
}

func (c *Client[S, C]) start() {
	var receiveBufPool = sync.Pool{
		New: func() any {
			// The Pool's New function should generally only return pointer
			// types, since a pointer can be put into the return interface
			// value without an allocation:
			return make([]byte, c.maxRecvPacketSize)
		},
	}

	go func() {
		for {
			if c.sock.Closed() {
				return // sockets can never redial
			}

			dat := receiveBufPool.Get().([]byte)
			n, err := c.sock.Read(dat)
			if err != nil {
				receiveBufPool.Put(dat)
				// TODO - I need a better way to wait if the socket is disconnected. If I remove the sleep then we will spin when disconnected
				time.Sleep(100 * time.Millisecond)
				continue
			}

			// if n > 30000 {
			// 	fmt.Println("Envoy.Recv: ", n)
			// }

			if n == 0 {
				// fmt.Println("Envoy.Recv: empty! ", n)
				receiveBufPool.Put(dat)
				continue
			} // Empty message

			if n >= len(dat) {
				// TODO: need to handle this case by dynamically resizing the receive buffer (within reason)
				// Full message
				fmt.Println("Envoy error: message was too big: ", n)
				receiveBufPool.Put(dat)
				continue
			}

			buf := cod.NewBufferFrom(dat)
			wireType := buf.ReadUint8()
			// if err != nil {
			// 	fmt.Println("Envoy: data unmarshal error:", err)
			// 	continue
			// }

			switch wireType {
			case wireTypeRequest:
				go func() {
					err := c.handleRequest(buf)
					if err != nil {
						fmt.Println("Envoy.Request error:", err)
					}
					receiveBufPool.Put(dat)
				}()
			case wireTypeResponse:
				// TODO - this one may not need to be run in a goroutine. b/c it will exit quickly enough?
				go func() {
					err := c.handleResponse(buf)
					if err != nil {
						fmt.Println("Envoy.Response error:", err)
					}
					receiveBufPool.Put(dat)
				}()
			case wireTypeMessage:
				go func() {
					err := c.handleMessage(buf)
					if err != nil {
						fmt.Println("Envoy.Message error:", err)
					}
					receiveBufPool.Put(dat)
				}()
			default:
				fmt.Printf("Envoy: unknown wire type error: %d\n", wireType)
				receiveBufPool.Put(dat)
			}
		}
	}()
}

func (c *Client[S, C]) handleResponse(buf *cod.Buffer) error {
	respId, err := buf.ReadUint32()
	if err != nil { return err } // TODO: Should this error inform the doRpc call?

	// TODO: codify this to remove the alloc
	respVal, err := c.clientDef.Responses.Deserialize(buf)
	if err != nil { return err } // TODO: Should this error inform the doRpc call?

	c.reqLock.Lock()
	defer c.reqLock.Unlock()

	callChan, ok := c.activeCalls[respId]
	if !ok {
		return fmt.Errorf("disassociated response")
	}

	// TODO: Should the channel disassociate if the receive side has closed?
	// Send the response to the appropriate call
	callChan <-respVal

	return nil
}

// func (c *Client[S, C]) handleRequest(rpcReq Request) error {
func (c *Client[S, C]) handleRequest(buf *cod.Buffer) error {
	reqId, err := buf.ReadUint32()
	if err != nil { return err }

	// TODO: codify this to remove the alloc
	reqVal, err := c.serviceDef.Requests.Deserialize(buf)
	if err != nil { return err }
	reqValType := reflect.TypeOf(reqVal)

	handler, ok := c.handlers[reqValType]
	if !ok {
		return fmt.Errorf("RPC Handler not set for type: %T", reqVal)
	}

	anyResp := handler(reqVal)

	// Build response
	sendBuf := sendBufPool.Get().(*cod.Buffer)
	sendBuf.Reset()
	defer sendBufPool.Put(sendBuf)

	sendBuf.WriteUint8(wireTypeResponse)
	sendBuf.WriteUint32(reqId)

	// TODO: codify this to remove the alloc
	err = c.serviceDef.Responses.Serialize(sendBuf, anyResp)
	if err != nil {
		return err
	}

	// TODO - check that n is correct?
	_, err = c.sock.Write(sendBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (c *Client[S, C]) handleMessage(buf *cod.Buffer) error {
	// TODO: codify this to remove the alloc
	msgVal, err := c.serviceDef.Requests.Deserialize(buf)
	if err != nil {
		return err
	}
	msgValType := reflect.TypeOf(msgVal)

	handler, ok := c.messageHandlers[msgValType]
	if !ok {
		return fmt.Errorf("Message Handler not set for type: %T", msgVal)
	}

	handler(msgVal)
	return nil
}

// TODO - Note: caller must be passed in as a pointer
func (client *Client[S, C]) registerCallers(caller any) {
	ty := reflect.TypeOf(caller)
	val := reflect.ValueOf(caller)
	numField := ty.Elem().NumField()
	for i := 0; i < numField; i++ {

		field := val.Elem().Field(i).Addr()

		fieldAny := field.Interface()

		switch rpcDef := fieldAny.(type) {
		case clientSetter:
			rpcDef.setClient(client)
		default:
			panic("Error: Must be a clientSetter")
		}
	}
}

func (client *Client[S, C]) registerHandlers(service any) {
	ty := reflect.TypeOf(service)
	val := reflect.ValueOf(service)
	numField := ty.NumField()
	for i := 0; i < numField; i++ {

		field := val.Field(i)

		fieldAny := field.Interface()

		switch rpcHandler := fieldAny.(type) {
		case RpcHandler:
			reqType, handler := rpcHandler.Handler()
			if handler == nil { panic("All Handlers must be defined!") }
			client.registerRpc(reqType, handler)
		case MsgHandler:
			msgType, handler := rpcHandler.Handler()
			if handler == nil { panic("All Handlers must be defined!") }
			client.registerMsg(msgType, handler)
		default:
			panic("ERROR") // TODO - must be an RpcDef or a message Def
		}
	}
}

func (client *Client[S, C]) registerRpc(reqValType reflect.Type, handler HandlerFunc) {
	if handler == nil { panic("Handler must not be nil!") }
	_, exists := client.handlers[reqValType]
	if exists {
		panic("Cant reregister the same handler type")
	}

	client.handlers[reqValType] = handler
}

func (client *Client[S, C]) registerMsg(msgValType reflect.Type, handler MessageHandlerFunc) {

	_, exists := client.messageHandlers[msgValType]
	if exists {
		panic("Cant reregister the same handler type")
	}

	client.messageHandlers[msgValType] = handler
}

type HandlerFunc func(req any) any
func makeRpcHandler[Req, Resp any](handler func(Req) Resp) HandlerFunc {
	return func(anyReq any) any {
		req, ok := anyReq.(Req)
		if !ok {
			panic(fmt.Errorf("Mismatched request types: %T", req))
		}

		res := handler(req)

		return res
	}
}

type MessageHandlerFunc func(req any)
func makeMsgHandler[M any](handler func(M)) MessageHandlerFunc {
	return func(anyMsg any) {
		msg, ok := anyMsg.(M)
		if !ok {
			panic(fmt.Errorf("Mismatched message types: %T", anyMsg))
		}

		handler(msg)
	}
}

func (c *Client[S, C]) doRpc(req any, timeout time.Duration) (any, error) {
	sendBuf := sendBufPool.Get().(*cod.Buffer)
	sendBuf.Reset()

	// Build RPC
	reqId, respChan, err := c.encodeRequest(sendBuf, req)
	defer c.cleanupResponseChannel(reqId) // Even if we fail, we definitely created a channel. so we need to make sure we clean that up. This code is coupled to encodeRequest
	if err != nil { return nil, err }

	// Send over socket
	// TODO - retry sending? Or push to a queue to be batch sent?
	// TODO - check that n is correct?
	_, err = c.sock.Write(sendBuf.Bytes())
	sendBufPool.Put(sendBuf)
	if err != nil {
		return nil, err
	}

	select {
	case rpcResp := <-respChan:
		return rpcResp, nil
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

func (c *Client[S, C]) doMsg(msg any) error {
	sendBuf := sendBufPool.Get().(*cod.Buffer)
	sendBuf.Reset()
	defer sendBufPool.Put(sendBuf)

	sendBuf.WriteUint8(wireTypeMessage)

	err := c.clientDef.Requests.Serialize(sendBuf, msg)
	if err != nil { return err }

	// println("Envoy.doMsg: ", len(sendBuf.Bytes()))

	// Send over socket
	// TODO - check that n is correct?
	_, err = c.sock.Write(sendBuf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

// Encode a request into the buffer
// Make a channel for the request
func (c *Client[S, C]) encodeRequest(buf *cod.Buffer, req any) (uint32, chan any, error) {
	reqId, respChan := c.generateReqId()

	buf.WriteUint8(wireTypeRequest)
	buf.WriteUint32(reqId)
	err := c.clientDef.Requests.Serialize(buf, req)
	if err != nil { return reqId, respChan, err } // TODO: can i get rid of this serialize failure?

	return reqId, respChan, nil
}

// Generate a request Id and create a channel for it
func (c *Client[S, C]) generateReqId() (uint32, chan any) {
	c.reqLock.Lock()
	defer c.reqLock.Unlock()

	respChan := make(chan any)

	// TODO - instead of a map could I use some other data structure? is it important that the Id is hard to guess? Should I use some other crypto-random-rng generator?
	var rngId uint32
	for {
		rngId = c.rng.Uint32()
		_, ok := c.activeCalls[rngId]
		if ok { continue } // call Id already exists, look for another

		c.activeCalls[rngId] = respChan
		break
	}
	return rngId, respChan
}

func (c *Client[S, C]) cleanupResponseChannel(id uint32) {
	c.reqLock.Lock()
	defer c.reqLock.Unlock()

	channel, ok := c.activeCalls[id]
	if !ok {
		// This is weird, the channel didn't exist but should have. Probably not worth panicking on
		fmt.Println("Envoy: Tried to cleanup channel that was already cleaned up")
		return
	}

	close(channel)
	c.activeCalls[id] = nil
	delete(c.activeCalls, id)
}

