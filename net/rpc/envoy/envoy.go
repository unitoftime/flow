package envoy

import (
	"fmt"
	"time"
	"math/rand"
	"reflect"

	"github.com/unitoftime/flow/net"
	// "github.com/unitoftime/flow/net/rpc"
)

// Needs
// 1. Bidirectional RPCs
// 2. Gotta have fire-and-forget style RPCs (ie just send it and don't block, like a msg)
// 3. Easy setup and management

type ServiceDefinition struct {
	Requests, Responses *net.UnionBuilder
}

// TODO - I think I'd prefer this to be based on method name and not based on input argument type
// TODO - reordering the definition, or switching between a message and an RPC will break api compatibility
func NewServiceDef(def any) ServiceDefinition {
	ty := reflect.TypeOf(def).Elem()
	fmt.Println(ty)
	numMethod := ty.NumMethod()
	fmt.Println(numMethod)

	requests := make([]any, 0)
	responses := make([]any, 0)
	for i := 0; i < numMethod; i++ {
		method := ty.Method(i)
		fmt.Println(method)

		numInputs := method.Type.NumIn()
		// for j := 0; j < numInputs; j++ {
		// 	in := method.Type.In(j)
		// 	fmt.Println(in)
		// }

		if numInputs != 1 { panic("We only support methods of form: func (req) (resp, error) or func (req) error") }
		reqType := method.Type.In(0)
		reqStruct := reflect.New(reqType).Elem().Interface()
		requests = append(requests, reqStruct)

		numOutputs := method.Type.NumOut()
		// for j := 0; j < numOutputs; j++ {
		// 	out := method.Type.Out(j)
		// 	fmt.Println(out)
		// }

		if numOutputs == 1 {
			// Rpc doesn't expect a response
		} else if numOutputs == 2 {
			// Rpc expects a response
			respType := method.Type.Out(0)
			respStruct := reflect.New(respType).Elem().Interface()
			responses = append(responses, respStruct)
		} else {
			panic("We only support methods of form: func (req) (resp, error) or func (req) error")
		}

		// TODO - check last argument should be an error
	}

	return ServiceDefinition{
		Requests: net.NewUnion(requests...),
		Responses: net.NewUnion(responses...),
	}
}

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

// TODO - requests and responses from interface definition
func NewClient(sock net.Socket, serviceDef, clientDef ServiceDefinition) *Client {
	client := &Client{
		sock: sock,

		serviceDef: serviceDef,
		clientDef: clientDef,

		handlers: make(map[reflect.Type]HandlerFunc),
		messageHandlers: make(map[reflect.Type]MessageHandlerFunc),
		activeCalls: make(map[uint32]chan Response),
	}

	client.start() // This doesn't block
	return client
}

type Client struct {
	sock net.Socket

	serviceDef, clientDef ServiceDefinition

	handlers map[reflect.Type]HandlerFunc
	messageHandlers map[reflect.Type]MessageHandlerFunc
	activeCalls map[uint32]chan Response
}

func (c *Client) Close() {
	// TODO - what to do about active calls and handlers?
	c.sock.Close()
}

func (c *Client) start() {
	dat := make([]byte, 8 * 1024) // TODO - hardcoded
	go func() {
		for {
			if c.sock.Closed() {
				return // sockets can never redial
			}

			err := c.sock.Read(dat)
			if err != nil {
				fmt.Println("ERROR: ", err)
				// TODO!!!!! - this might cause the for loop to spin if we are trying to reconnect, for example
				time.Sleep(100 * time.Millisecond)
				continue
			}

			msg, err := rpcSerdes.Unmarshal(dat)
			if err != nil {
				fmt.Println("ERROR: ", err)
				continue
			}

			// If the message was empty, just continue to the next one
			if msg == nil { continue }

			switch typedMsg := msg.(type) {
			case Request:
				resp, err := c.HandleRequest(typedMsg)
				if err != nil {
					fmt.Println("ERROR: ", err)
				}

				respDat, err := rpcSerdes.Marshal(resp)
				if err != nil {
					fmt.Println("ERROR: ", err)
				}

				err = c.sock.Write(respDat)
				if err != nil {
					fmt.Println("ERROR: ", err)
				}
			case Response:
				err := c.HandleResponse(typedMsg)
				if err != nil {
					fmt.Println("ERROR: ", err)
				}
			case Message:
				err := c.HandleMessage(typedMsg)
				if err != nil {
					fmt.Println("ERROR: ", err)
				}
			default:
				fmt.Printf("Unknown message type: %T\n", typedMsg)
			}
		}
	}()
}


type MessageHandlerFunc func(req any) error
type HandlerFunc func(req any) ([]byte, error)

func (c *Client) HandleResponse(rpcResp Response) error {
	callChan, ok := c.activeCalls[rpcResp.Id]
	if !ok {
		return fmt.Errorf("Disassociated RpcResponse")
	}

	// Send the response to the appropriate call
	callChan <-rpcResp

	// Cleanup
	close(callChan)
	c.activeCalls[rpcResp.Id] = nil
	delete(c.activeCalls, rpcResp.Id)
	return nil
}

func (c *Client) HandleRequest(rpcReq Request) (Response, error) {
	rpcResp := Response{
		Id: rpcReq.Id,
	}

	reqVal, err := c.serviceDef.Requests.Deserialize(rpcReq.Data)
	if err != nil { return rpcResp, err }
	reqValType := reflect.TypeOf(reqVal)

	handler, ok := c.handlers[reqValType]
	if !ok {
		return rpcResp, fmt.Errorf("RPC Handler not set for type: %T", reqVal)
	}

	data, err := handler(reqVal)
	rpcResp.Data = data

	return rpcResp, err
}

func (c *Client) HandleMessage(msg Message) error {
	msgVal, err := c.serviceDef.Requests.Deserialize(msg.Data)
	if err != nil { return err }
	msgValType := reflect.TypeOf(msgVal)

	handler, ok := c.messageHandlers[msgValType]
	if !ok {
		return fmt.Errorf("Message Handler not set for type: %T", msgVal)
	}

	return handler(msgVal)
}

func RegisterMessage[M any](client *Client, handler func(M) error) {
	var msgVal M
	msgValType := reflect.TypeOf(msgVal)
	_, exists := client.messageHandlers[msgValType]
	if exists {
		panic("Cant reregister the same handler type")
	}

	// Create a handler function
	generalHandlerFunc := func(anyMsg any) error {
		msg, ok := anyMsg.(M)
		if !ok {
			panic(fmt.Errorf("Mismatched message types: %T, %T", msgVal, msg))
		}

		err := handler(msg)
		return err
	}

	// Store the handler function
	client.messageHandlers[msgValType] = generalHandlerFunc
}

func Register[Req any, Resp any](client *Client, handler func(Req) (Resp, error)) {
	var reqVal Req
	reqValType := reflect.TypeOf(reqVal)
	_, exists := client.handlers[reqValType]
	if exists {
		panic("Cant reregister the same handler type")
	}

	// Create a handler function
	generalHandlerFunc := func(anyReq any) ([]byte, error) {
		req, ok := anyReq.(Req)
		if !ok {
			panic(fmt.Errorf("Mismatched request types: %T, %T", reqVal, req))
		}

		res, err := handler(req)
		if err != nil {
			return nil, err
		}

		data, err := client.serviceDef.Responses.Serialize(res)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	// Store the handler function
	client.handlers[reqValType] = generalHandlerFunc
}



// Client - making requests
// func (c *Client) MakeRequest(req any) (any, error) {
// 	dat, err := c.reqSerdes.Serialize(req)
// 	if err != nil { return err }

// 	reqDat, err := rpcSerdes.Marshal(Request{
// 		Id: 0, // TODO
// 		Data: dat,
// 	})
// 	if err != nil { return err }

// 	err = c.sock.Write(reqDat)
// 	return err
// }

func NewCall[Req, Resp any](client *Client) *Call[Req, Resp] {
	rngSrc := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rngSrc) // TODO - push this up to the client?

	return &Call[Req, Resp]{
		client: client,
		rng: rng,
	}
}
type Call[Req, Resp any] struct {
	client *Client
	rng *rand.Rand
}

func (c *Call[Req, Resp]) Do(req Req) (Resp, error) {
	var resp Resp
	rpcReq, err := c.Make(req)
	if err != nil { return resp, err }

	// TODO!!! - check if this ID is already being used, if it is, then use a different one

	// Make a channel to wait for a response on this Id
	respChan := make(chan Response)
	c.client.activeCalls[rpcReq.Id] = respChan

	// Send over socket
	reqDat, err := rpcSerdes.Marshal(rpcReq)
	if err != nil { return resp, err }

	err = c.client.sock.Write(reqDat)
	if err != nil { return resp, err }

	// TODO!!! - set some timeout too
	rpcResp := <-respChan
	return c.Unmake(rpcResp)
}

func (c *Call[Req, Resp]) Make(req Req) (Request, error) {
	dat, err := c.client.clientDef.Requests.Serialize(req)

	return Request{
		Id: c.rng.Uint32(),
		Data: dat,
	}, err
}

func (c *Call[Req, Resp]) Unmake(rpcResp Response) (Resp, error) {
	anyResp, err := c.client.clientDef.Responses.Deserialize(rpcResp.Data)
	var resp Resp
	if err != nil { return resp, err }
	resp, ok := anyResp.(Resp)
	if !ok { panic("Mismatched type!") }
	return resp, err
}

func NewMessage[A any](client *Client) *Msg[A] {
	return &Msg[A]{
		client: client,
	}
}
type Msg[A any] struct {
	client *Client
}
func (m *Msg[A]) Send(req A) error {
	rpcMsg, err := m.Make(req)
	if err != nil { return err }

	// Send over socket
	reqDat, err := rpcSerdes.Marshal(rpcMsg)
	if err != nil { return err }

	err = m.client.sock.Write(reqDat)
	if err != nil { return err }

	return nil
}

func (m *Msg[A]) Make(req A) (Message, error) {
	dat, err := m.client.clientDef.Requests.Serialize(req)

	return Message{
		Data: dat,
	}, err
}
