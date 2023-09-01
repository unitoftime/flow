package rpc

import (
	"fmt"
	"testing"

	"github.com/unitoftime/flow/net"
)

type Req struct {
	Val int32
}

type Res struct {
	Val int32
}

type BasicService interface {
	HandleThing(Req) (Res, error)
	SendThing(Req) (error)
}

type testService struct {
	// put data here
}

func (t *testService) HandleThing(req Req) (Res, error) {
	return Res{
		req.Val+2,
	}, nil
}

func (t *testService) SendThing(req Req) (error) {
	return Res{
		req.Val+2,
	}, nil
}

// type testServiceClient struct {
// 	// put data here
// }

// func (t *testServiceClient) CallHandleThing(req Req) (Res, error) {
// 	call := NewCall[Req, Res](t.service)
// 	req, err := call.Make(Req{9999})
// 	if err != nil { panic(err) }
// 	fmt.Println(req)
// 	return 
// }

// func (t *testServiceClient) CallSendThing(req Req) (error) {
// 	return nil
// }

func TestRpc(t *testing.T) {
	serviceDef := NewServiceDef(new(BasicService))

	serviceHandler := &testService{}
	server := serviceDef.NewServer(sock, serviceHandler)

	client := serviceDef.NewClient(sock,


	NewServiceDef(new(BasicService))

	reqSerdes := net.NewUnion(Req{})
	resSerdes := net.NewUnion(Res{})
	service := NewService(reqSerdes, resSerdes)

	// Client Side
	// client := service.Client()
	// rpc.Register(client, reqSerdes, resSerdes)
	// res, err := rpc.Call(Req{5})
	// res, err := rpc.CallAsync(Req{5})

	call := NewCall[Req, Res](service)
	req, err := call.Make(Req{9999})
	if err != nil { panic(err) }
	fmt.Println(req)

	// Server side
	svc := testService{}
	Register(service, svc.HandleThing)
	rpcResp, err := service.HandleRequest(req)
	if err != nil { panic(err) }

	// ClientSide
	resp, err := call.Unmake(rpcResp)
	if err != nil { panic(err) }
	fmt.Println(rpcResp.Id)
	fmt.Println(resp)


	// dat, err := reqSerdes.Serialize(Req{555})
	// if err != nil { panic(err) }

	// resp, err := service.Handle(RpcRequest{
	// 	Identifier: 111,
	// 	Data: dat,
	// })
	// if err != nil { panic(err) }
	// fmt.Println(resp)

	// respDat, err := resSerdes.Deserialize(resp.Data)
	// if err != nil { panic(err) }
	// fmt.Println(respDat)

	// client := rpc.NewClient(sock, serdes)
	// rpc.Register(client, reqSerdes, resSerdes)
	// res, err := rpc.Call(Req{5})

	// res, err := rpc.CallAsync(Req{5})
}


//--------------------------------------------------------------------------------
// Attempt 1

// import (
// 	"fmt"
// 	"testing"

// 	"github.com/unitoftime/flow/net"
// )

// type Req struct {
// 	Val int32
// }

// type Res struct {
// 	Val int32
// }

// func handler(req Req) (Res, error) {
// 	return Res{
// 		req.Val,
// 	}, nil
// }

// func TestRpc(t *testing.T) {
// 	reqSerdes := net.NewUnion(Req{})
// 	resSerdes := net.NewUnion(Res{})
// 	service := NewService(reqSerdes, resSerdes)

// 	// Client Side
// 	// client := service.Client()
// 	// rpc.Register(client, reqSerdes, resSerdes)
// 	// res, err := rpc.Call(Req{5})
// 	// res, err := rpc.CallAsync(Req{5})

// 	call := NewCall[Req, Res](service)
// 	req, err := call.Make(Req{9999})
// 	if err != nil { panic(err) }
// 	fmt.Println(req)

// 	// Server side
// 	Register(service, handler)
// 	rpcResp, err := service.HandleRequest(req)
// 	if err != nil { panic(err) }

// 	// ClientSide
// 	resp, err := call.Unmake(rpcResp)
// 	if err != nil { panic(err) }
// 	fmt.Println(rpcResp.Id)
// 	fmt.Println(resp)


// 	// dat, err := reqSerdes.Serialize(Req{555})
// 	// if err != nil { panic(err) }

// 	// resp, err := service.Handle(RpcRequest{
// 	// 	Identifier: 111,
// 	// 	Data: dat,
// 	// })
// 	// if err != nil { panic(err) }
// 	// fmt.Println(resp)

// 	// respDat, err := resSerdes.Deserialize(resp.Data)
// 	// if err != nil { panic(err) }
// 	// fmt.Println(respDat)

// 	// client := rpc.NewClient(sock, serdes)
// 	// rpc.Register(client, reqSerdes, resSerdes)
// 	// res, err := rpc.Call(Req{5})

// 	// res, err := rpc.CallAsync(Req{5})
// }
