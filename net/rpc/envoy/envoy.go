package envoy

import (
	"fmt"
	"time"
	"errors"
	"math/rand"
	"reflect"
	"sync"

	"github.com/unitoftime/flow/net"
)

// Has
// 1. Bidirectional RPCs
// 2. Fire-and-forget style RPCs (ie just send it and don't block, like a msg)
// 3. Easy setup and management

// Wants
// 4. Different reliability levels
// 5. Automatic retries of messages?
// 6. Message batching?

var ErrTimeout = errors.New("timeout reached")
var ErrDisconnected = errors.New("socket disconnected")

// Internal serialization
var rpcSerdes serdes
func init() {
	rpcSerdes = serdes{
		union: net.NewUnion(
			Request{},
			Response{},
			Message{},
		),
	}
}

type serdes struct {
	union *net.UnionBuilder
}
func (s *serdes) Marshal(v any) ([]byte, error) {
	return s.union.Serialize(v)
}
func (s *serdes) Unmarshal(dat []byte) (any, error) {
	return s.union.Deserialize(dat)
}

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
	anyResp, err := d.client.doRpc(req, 5 * time.Second)
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

// Returns the client side of the interface
func (d InterfaceDef[S, C]) NewClient() *Client[C, S] {
	// Note: The C and S are reversed because we call the service and serve the client
	client := newClient[C, S](d.clientApi, d.serviceApi)

	return client
}

// Returns the server side of the interface
func (d InterfaceDef[S, C]) NewServer() *Client[S, C] {
	client := newClient[S, C](d.serviceApi, d.clientApi)

	return client
}


// TODO - I should make an interface to better capture the fact that this is just for serialization
// TODO - could I make servicedef generic on the interface type. Then when I register handlers I just pass in a struct which implements the interface?

// TODO - I think I'd prefer this to be based on method name and not based on input argument type
// TODO - reordering the definition, or switching between a message and an RPC will break api compatibility
type serviceDef struct {
	Requests, Responses *net.UnionBuilder
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
		Requests: net.NewUnion(requests...),
		Responses: net.NewUnion(responses...),
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
	activeCalls map[uint32]chan Response

	rng *rand.Rand
}

func newClient[S, C any](serviceDef, clientDef serviceDef) *Client[S, C] {
	rngSrc := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rngSrc) // TODO - push this up to the client?

	client := &Client[S, C]{
		serviceDef: serviceDef,
		clientDef: clientDef,

		handlers: make(map[reflect.Type]HandlerFunc),
		messageHandlers: make(map[reflect.Type]MessageHandlerFunc),
		activeCalls: make(map[uint32]chan Response),

		rng: rng,
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

func (c *Client[S, C]) start() {
	dat := make([]byte, 64 * 1024) // TODO!!!! - hardcoded to max webrtc packet size
	go func() {
		for {
			if c.sock.Closed() {
				return // sockets can never redial
			}

			n, err := c.sock.Read(dat)
			if err != nil {
				// TODO - I need a better way to wait if the socket is disconnected. If I remove the sleep then we will spin when disconnected
				time.Sleep(100 * time.Millisecond)
				continue
			}
			if n == 0 { continue } // Empty message


			msg, err := rpcSerdes.Unmarshal(dat)
			if err != nil {
				fmt.Println("Envoy: data unmarshal error:", err)
				continue
			}

			// If the message was empty, just continue to the next one
			if msg == nil { continue }

			switch typedMsg := msg.(type) {
			case Request:
				go func() {
					err := c.handleRequest(typedMsg)
					if err != nil {
						fmt.Println("Envoy.Request error:", err)
					}
				}()
			case Response:
				// TODO - this one may not need to be run in a goroutine. b/c it will exit quickly enough?
				go func() {
					err := c.handleResponse(typedMsg)
					if err != nil {
						fmt.Println("Envoy.Response error:", err)
					}
				}()
			case Message:
				go func() {
					err := c.handleMessage(typedMsg)
					if err != nil {
						fmt.Println("Envoy.Message error:", err)
					}
				}()
			default:
				fmt.Printf("Envoy: unknown type error: %T\n", typedMsg)
			}
		}
	}()
}

func (c *Client[S, C]) handleResponse(rpcResp Response) error {
	c.reqLock.Lock()
	defer c.reqLock.Unlock()

	callChan, ok := c.activeCalls[rpcResp.Id]
	if !ok {
		return fmt.Errorf("disassociated response")
	}

	// Send the response to the appropriate call
	callChan <-rpcResp

	return nil
}

func (c *Client[S, C]) handleRequest(rpcReq Request) error {
	reqVal, err := c.serviceDef.Requests.Deserialize(rpcReq.Data)
	if err != nil { return err }
	reqValType := reflect.TypeOf(reqVal)

	handler, ok := c.handlers[reqValType]
	if !ok {
		return fmt.Errorf("RPC Handler not set for type: %T", reqVal)
	}

	anyResp := handler(reqVal)

	data, err := c.serviceDef.Responses.Serialize(anyResp)
	if err != nil {
		return err
	}

	rpcResp := Response{
		Id: rpcReq.Id,
		Data: data,
	}

	respDat, err := rpcSerdes.Marshal(rpcResp)
	if err != nil {
		return err
	}

	// TODO - check that n is correct?
	_, err = c.sock.Write(respDat)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client[S, C]) handleMessage(msg Message) error {
	msgVal, err := c.serviceDef.Requests.Deserialize(msg.Data)
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
	rpcReq, respChan, err := c.makeRequest(req)
	if err != nil { return nil, err }

	defer c.cleanupResponseChannel(rpcReq.Id)

	// Send over socket
	reqDat, err := rpcSerdes.Marshal(rpcReq)
	if err != nil { return nil, err }

	// TODO - retry sending? Or push to a queue to be batch sent?
	// TODO - check that n is correct?
	_, err = c.sock.Write(reqDat)
	if err != nil {
		// TODO - snuffing underlying error because if we coudln't send it means we are trying to reconnect.
		return nil, ErrDisconnected
	}

	select {
	case rpcResp := <-respChan:
		return c.clientDef.Responses.Deserialize(rpcResp.Data)
	case <-time.After(timeout):
		return nil, ErrTimeout
	}
}

func (c *Client[S, C]) doMsg(msg any) error {
	dat, err := c.clientDef.Requests.Serialize(msg)
	if err != nil { return err }

	rpcMsg := Message{
		Data: dat,
	}

	// Send over socket
	msgDat, err := rpcSerdes.Marshal(rpcMsg)
	if err != nil { return err }

	// TODO - check that n is correct?
	_, err = c.sock.Write(msgDat)
	if err != nil {
		// TODO - snuffing underlying error because if we coudln't send it means we are trying to reconnect.
		return ErrDisconnected
	}

	return nil
}

func (c *Client[S, C]) makeRequest(req any) (Request, chan Response, error) {
	dat, err := c.clientDef.Requests.Serialize(req)
	if err != nil { return Request{}, nil, err }

	c.reqLock.Lock()
	defer c.reqLock.Unlock()

	respChan := make(chan Response)

	// TODO - instead of a map could I use some other data structure? is it important that the Id is hard to guess? Should I use some other crypto-random-rng generator?
	var rngId uint32
	for {
		rngId = c.rng.Uint32()
		_, ok := c.activeCalls[rngId]
		if ok { continue } // call Id already exists, look for another

		c.activeCalls[rngId] = respChan
		break
	}

	request := Request{
		Id: rngId,
		Data: dat,
	}
	return request, respChan, nil
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
