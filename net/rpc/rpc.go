package rpc

import (
	"fmt"
	"math/rand"
	"reflect"
	"time"

	"github.com/unitoftime/flow/net"
)

// Needs
// 1. Bidirectional RPCs
// 2. Gotta have fire-and-forget style RPCs (ie just send it and don't block, like a msg)

type RpcRequest struct {
	Id uint32 // Tracks the request Id
	Data []byte
}

type RpcResponse struct {
	Id uint32 // Tracks the request Id
	Data []byte
}

type Service struct {
	reqSerdes *net.UnionBuilder
	resSerdes *net.UnionBuilder

	handlers map[reflect.Type]HandlerFunc

	activeCalls map[uint32]chan RpcResponse
}
func NewService(reqSerdes, resSerdes *net.UnionBuilder) *Service {
	return &Service{
		reqSerdes: reqSerdes,
		resSerdes: resSerdes,

		handlers: make(map[reflect.Type]HandlerFunc),
		// clients: make(map[reflect.Type]any),

		activeCalls: make(map[uint32]chan RpcResponse),
	}
}

func (s *Service) HandleResponse(rpcResp RpcResponse) error {
	callChan, ok := s.activeCalls[rpcResp.Id]
	if !ok {
		return fmt.Errorf("Disassociated RpcResponse")
	}

	// Send the response to the appropriate call
	callChan <-rpcResp

	// Cleanup
	close(callChan)
	s.activeCalls[rpcResp.Id] = nil
	delete(s.activeCalls, rpcResp.Id)
	return nil
}

func (s *Service) HandleRequest(rpcReq RpcRequest) (RpcResponse, error) {
	rpcResp := RpcResponse{
		Id: rpcReq.Id,
	}

	reqVal, err := s.reqSerdes.Deserialize(rpcReq.Data)
	if err != nil { return rpcResp, err }
	reqValType := reflect.TypeOf(reqVal)

	handler, ok := s.handlers[reqValType]
	if !ok {
		return rpcResp, fmt.Errorf("Handler not set for type: %T", reqVal)
	}

	data, err := handler(reqVal)
	rpcResp.Data = data

	return rpcResp, err
}

type HandlerFunc func(req any) ([]byte, error)

func Register[Req any, Resp any](service *Service, handler func(Req) (Resp, error)) {
	var reqVal Req
	reqValType := reflect.TypeOf(reqVal)
	_, exists := service.handlers[reqValType]
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

		data, err := service.resSerdes.Serialize(res)
		if err != nil {
			return nil, err
		}

		return data, nil
	}

	// Store the handler function
	service.handlers[reqValType] = generalHandlerFunc
}

func NewCall[Req, Resp any](service *Service) Call[Req, Resp] {
	rngSrc := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(rngSrc)

	return Call[Req, Resp]{
		service: service,
		rng: rng,
	}
}
type Call[Req, Resp any] struct {
	service *Service
	rng *rand.Rand
}
func (c *Call[Req, Resp]) Request() any {
	var ret Req
	return ret
}
func (c *Call[Req, Resp]) Response() any {
	var ret Resp
	return ret
}

// // TODO - ideally
// func (c *Call[Req, Resp]) Do(req Req) (Resp, error) {
// 	var resp Resp
// 	rpcReq, err := c.Make(req)
// 	if err != nil { return resp, err }

// 	// TODO!!! - check if this ID is already being used, if it is, then use a different one

// 	// Make a channel to wait for a response on this Id
// 	respChan := make(chan RpcResponse)
// 	c.service.activeCalls[rpcReq.Id] = respChan

// 	// Send over socket
// 	err = sock.Send(rpcReq)
// 	if err != nil { return resp, err }

// 	// TODO - set some timeout too
// 	rpcResp := <-respChan
// 	return Unmake(rpcResp)
// }

func (c *Call[Req, Resp]) Make(req Req) (RpcRequest, error) {
	dat, err := c.service.reqSerdes.Serialize(req)

	return RpcRequest{
		Id: c.rng.Uint32(),
		Data: dat,
	}, err
}

func (c *Call[Req, Resp]) Unmake(rpcResp RpcResponse) (Resp, error) {
	anyResp, err := c.service.resSerdes.Deserialize(rpcResp.Data)
	var resp Resp
	if err != nil { return resp, err }
	resp, ok := anyResp.(Resp)
	if !ok { panic("Mismatched type!") }
	return resp, err
}

// type Client[ struct {
// 	reqSerdes *net.UnionBuilder
// 	resSerdes *net.UnionBuilder
// }

// func Client(service *Service) *Client {
// 	return &Client{
// 		reqSerdes: service.reqSerdes,
// 		resSerdes: service.resSerdes,
// 	}
// }

// func (c *Client) NewRequest(val any) (RpcRequest, error) {
// 	dat, err := c.reqSerdes.Serialize(val)
// 	return RpcRequest{
// 		Id: math.RandUint16(),
// 		Data: dat,
// 	}, err
// }

	

// func (s *Service) NewRequest(val any) (RpcRequest, error) {
// 	dat, err := s.reqSerdes.Serialize(val)
// 	return RpcRequest{
// 		Id: math.RandUint16(),
// 		Data: dat,
// 	}, err
// }

// func (s *Service) NewResponse(id uint16, val any) (RpcResponse, error) {
// 	dat, err := s.resSerdes.Serialize(val)
// 	return RpcResponse{
// 		Id: id,
// 		Data: dat,
// 	}, err
// }

	// dat, err := reqSerdes.Serialize(Req{555})
	// if err != nil { panic(err) }

	// if err != nil { panic(err) }
	// fmt.Println(resp)

	// respDat, err := resSerdes.Deserialize(resp.Data)
	// if err != nil { panic(err) }
	// fmt.Println(respDat)


//--------------------------------------------------------------------------------
// Attempt 1

// package rpc

// import (
// 	"fmt"
// 	"math/rand"
// 	"reflect"
// 	"time"

// 	"github.com/unitoftime/flow/net"
// )

// // type Service struct {
// // 	sock net.Socket
// // }

// // type Client struct {
// // 	sock net.Socket
// // }

// type RpcRequest struct {
// 	Id uint32 // Tracks the request Id
// 	Data []byte
// }

// type RpcResponse struct {
// 	Id uint32 // Tracks the request Id
// 	Data []byte
// }

// type Service struct {
// 	reqSerdes *net.UnionBuilder
// 	resSerdes *net.UnionBuilder

// 	handlers map[reflect.Type]HandlerFunc

// 	activeCalls map[uint32]chan RpcResponse
// }
// func NewService(reqSerdes, resSerdes *net.UnionBuilder) *Service {
// 	return &Service{
// 		reqSerdes: reqSerdes,
// 		resSerdes: resSerdes,

// 		handlers: make(map[reflect.Type]HandlerFunc),
// 		// clients: make(map[reflect.Type]any),

// 		activeCalls: make(map[uint32]chan RpcResponse),
// 	}
// }

// func (s *Service) HandleResponse(rpcResp RpcResponse) error {
// 	callChan, ok := s.activeCalls[rpcResp.Id]
// 	if !ok {
// 		return fmt.Errorf("Disassociated RpcResponse")
// 	}

// 	// Send the response to the appropriate call
// 	callChan <-rpcResp

// 	// Cleanup
// 	close(callChan)
// 	s.activeCalls[rpcResp.Id] = nil
// 	delete(s.activeCalls, rpcResp.Id)
// 	return nil
// }

// func (s *Service) HandleRequest(rpcReq RpcRequest) (RpcResponse, error) {
// 	rpcResp := RpcResponse{
// 		Id: rpcReq.Id,
// 	}

// 	reqVal, err := s.reqSerdes.Deserialize(rpcReq.Data)
// 	if err != nil { return rpcResp, err }
// 	reqValType := reflect.TypeOf(reqVal)

// 	handler, ok := s.handlers[reqValType]
// 	if !ok {
// 		return rpcResp, fmt.Errorf("Handler not set for type: %T", reqVal)
// 	}

// 	data, err := handler(reqVal)
// 	rpcResp.Data = data

// 	return rpcResp, err
// }

// type HandlerFunc func(req any) ([]byte, error)

// func Register[Req any, Resp any](service *Service, handler func(Req) (Resp, error)) {
// 	var reqVal Req
// 	reqValType := reflect.TypeOf(reqVal)
// 	_, exists := service.handlers[reqValType]
// 	if exists {
// 		panic("Cant reregister the same handler type")
// 	}

// 	// Create a handler function
// 	generalHandlerFunc := func(anyReq any) ([]byte, error) {
// 		req, ok := anyReq.(Req)
// 		if !ok {
// 			panic(fmt.Errorf("Mismatched request types: %T, %T", reqVal, req))
// 		}

// 		res, err := handler(req)
// 		if err != nil {
// 			return nil, err
// 		}

// 		data, err := service.resSerdes.Serialize(res)
// 		if err != nil {
// 			return nil, err
// 		}

// 		return data, nil
// 	}

// 	// Store the handler function
// 	service.handlers[reqValType] = generalHandlerFunc
// }

// func NewCall[Req, Resp any](service *Service) Call[Req, Resp] {
// 	rngSrc := rand.NewSource(time.Now().UnixNano())
// 	rng := rand.New(rngSrc)

// 	return Call[Req, Resp]{
// 		service: service,
// 		rng: rng,
// 	}
// }
// type Call[Req, Resp any] struct {
// 	service *Service
// 	rng *rand.Rand
// }
// func (c *Call[Req, Resp]) Request() any {
// 	var ret Req
// 	return ret
// }
// func (c *Call[Req, Resp]) Response() any {
// 	var ret Resp
// 	return ret
// }

// // // TODO - ideally
// // func (c *Call[Req, Resp]) Do(req Req) (Resp, error) {
// // 	var resp Resp
// // 	rpcReq, err := c.Make(req)
// // 	if err != nil { return resp, err }

// // 	// TODO!!! - check if this ID is already being used, if it is, then use a different one

// // 	// Make a channel to wait for a response on this Id
// // 	respChan := make(chan RpcResponse)
// // 	c.service.activeCalls[rpcReq.Id] = respChan

// // 	// Send over socket
// // 	err = sock.Send(rpcReq)
// // 	if err != nil { return resp, err }

// // 	// TODO - set some timeout too
// // 	rpcResp := <-respChan
// // 	return Unmake(rpcResp)
// // }

// func (c *Call[Req, Resp]) Make(req Req) (RpcRequest, error) {
// 	dat, err := c.service.reqSerdes.Serialize(req)

// 	return RpcRequest{
// 		Id: c.rng.Uint32(),
// 		Data: dat,
// 	}, err
// }

// func (c *Call[Req, Resp]) Unmake(rpcResp RpcResponse) (Resp, error) {
// 	anyResp, err := c.service.resSerdes.Deserialize(rpcResp.Data)
// 	var resp Resp
// 	if err != nil { return resp, err }
// 	resp, ok := anyResp.(Resp)
// 	if !ok { panic("Mismatched type!") }
// 	return resp, err
// }

// // type Client[ struct {
// // 	reqSerdes *net.UnionBuilder
// // 	resSerdes *net.UnionBuilder
// // }

// // func Client(service *Service) *Client {
// // 	return &Client{
// // 		reqSerdes: service.reqSerdes,
// // 		resSerdes: service.resSerdes,
// // 	}
// // }

// // func (c *Client) NewRequest(val any) (RpcRequest, error) {
// // 	dat, err := c.reqSerdes.Serialize(val)
// // 	return RpcRequest{
// // 		Id: math.RandUint16(),
// // 		Data: dat,
// // 	}, err
// // }

	

// // func (s *Service) NewRequest(val any) (RpcRequest, error) {
// // 	dat, err := s.reqSerdes.Serialize(val)
// // 	return RpcRequest{
// // 		Id: math.RandUint16(),
// // 		Data: dat,
// // 	}, err
// // }

// // func (s *Service) NewResponse(id uint16, val any) (RpcResponse, error) {
// // 	dat, err := s.resSerdes.Serialize(val)
// // 	return RpcResponse{
// // 		Id: id,
// // 		Data: dat,
// // 	}, err
// // }

// 	// dat, err := reqSerdes.Serialize(Req{555})
// 	// if err != nil { panic(err) }

// 	// if err != nil { panic(err) }
// 	// fmt.Println(resp)

// 	// respDat, err := resSerdes.Deserialize(resp.Data)
// 	// if err != nil { panic(err) }
// 	// fmt.Println(respDat)
