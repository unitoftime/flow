package envoy

import (
	"fmt"
	"time"
	"testing"

	"github.com/unitoftime/flow/net"
	// "github.com/unitoftime/flow/net/rpc"
)

type TestServiceApi interface {
	DoThing(ServerReq) (ServerResp, error)
	HandleMsg(ServerMsg) error
}

type TestServiceApiStruct struct {
	DoThing RpcDef[ServerReq, ServerResp]
	HandleMsg MsgDef[ServerMsg]
}

type TestServiceClientApiStruct struct {
	ClientDoThing RpcDef[ClientReq, ClientResp]
}

type TestServiceClientApi interface {
	ClientDoThing(ClientReq) (ClientResp, error)
}

type ServerMsg struct {
	Val uint16
}
type ServerReq struct {
	Val uint16
}
type ServerResp struct {
	Val uint16
}

type TestService struct {
}
func (s *TestService) HandleMsg(r ServerMsg) error {
	fmt.Println("HandleMsg: ", r)
	return nil
}

func (s *TestService) DoThing(r ServerReq) (ServerResp, error) {
	fmt.Println("DoThing: ", r)
	return ServerResp{
		r.Val + 1,
	}, nil
}

type ClientReq struct {
	Val uint16
}
type ClientResp struct {
	Val uint16
}

type TestServiceClient struct {
}
func (s *TestServiceClient) ClientDoThing(r ClientReq) (ClientResp, error) {
	fmt.Println("ClientDoThing: ", r)
	return ClientResp{
		r.Val + 100,
	}, nil
}

func TestMain(t *testing.T) {
	// testServiceDef := TestServiceApiStruct{}
	// testServiceClientDef := TestServiceClientApiStruct{}

	// {
	// 	s1 := DefineService(testServiceDef)
	// 	s2 := DefineService(testServiceClientDef)

	// 	fmt.Println(s1)
	// 	fmt.Println(s2)

	// 	ts := &TestService{}
	// 	testServiceDef.DoThing.Register(ts.DoThing)
	// 	// Do some call to create service

	// 	// Do some calls to create calls and message
	// 	// testServiceClientDef.ClientDoThing.Get()
	// 	// you need to like combine these two structs somehow then create a client based on that by passing in a socket. preregister all of your handlers, then on the outside you can just lookup the calls and msg executors or whatever
	// }

	// tsDoThingCall := testServiceDef.DoThing.Call()

	// testService := DefineService(TestServiceApiStruct{})
	// fmt.Println(testService)

	// Full data would be something like:
	// rpc1 := RpcDef[ServerReq, ServerResp]()
	// rpc2 := MsgDef[ServerMsg]()

	// // iface := new(TestServiceApi)
	// // Add(def, iface.DoThing)
	// // Add(def, iface.HandleMsg)
	// rpc1.Register(client, ts.DoThing)
	// rpc2.Register(client, ts.HandleMsg)

	// doThingCall := rpc1.Call(client)
	// handleMessage := rpc2.Call(client)

	// definitions:
	// interfaceDef := NewInterfaceDef(new(TestServiceApi), new(TestServiceClientApi))

	// interfaceDef := NewInterfaceDef[TestServiceApi, TestServiceClientApi]()
	// interfaceDef := NewInterfaceDef2[TestServiceApiStruct, TestServiceClientApiStruct]()

	// Server
	go func() {
		listenConfig := net.ListenConfig{
			Url: "tcp://localhost:8000",
		}
		listener, err := listenConfig.Listen()
		if err != nil { panic(err) }

		for {
			sock, err := listener.Accept()
			if err != nil { panic(err) }

			ts := &TestService{}

			interfaceDef := NewInterfaceDef[TestServiceApiStruct, TestServiceClientApiStruct]()
			interfaceDef.Service.DoThing.Register(ts.DoThing)
			interfaceDef.Service.HandleMsg.Register(ts.HandleMsg)

			client := interfaceDef.NewServer()
			client.Connect(sock)

			// Register(client, ts.DoThing)
			// RegisterMessage(client, ts.HandleMsg)

			// call := client.Call.ClientDoThing.Get()
			call := NewCall2(client, client.Call.ClientDoThing)

			// call := interfaceDef.Client.ClientDoThing.Get(client)
			// call := NewCall[ClientReq, ClientResp](client)

			// var iface TestServiceClientApi
			// call := GetCall(client, func() {iface.ClientDoThing)
			fmt.Println(client)

			resp, err := call.Do(ClientReq{1})
			if err != nil { panic(err) }
			fmt.Println("ClientDoThingResp: ", resp)
		}
	}()

	dialConfig := net.DialConfig{
		Url: "tcp://localhost:8000",
	}
	sock := dialConfig.Dial()

	ts := &TestServiceClient{}
	interfaceDef := NewInterfaceDef[TestServiceApiStruct, TestServiceClientApiStruct]()
	interfaceDef.Client.ClientDoThing.Register(ts.ClientDoThing)

	client := interfaceDef.NewClient()
	client.Connect(sock)

	// call := client.Call.DoThing.Get()
	call := NewCall2(client, client.Call.DoThing)

	// call := interfaceDef.Service.DoThing.Get(client)
	// Register(client, ts.ClientDoThing)
	// call := NewCall[ServerReq, ServerResp](client)
	// msgCall := NewMessage[ServerMsg](client)
	msgCall := NewMessage2(client, client.Call.HandleMsg)
	fmt.Println(client)

	time.Sleep(1 * time.Second)

	msgCall.Send(ServerMsg{9})

	// resp := client.MakeRequest(Req{5})
	resp, err := call.Do(ServerReq{5})
	if err != nil { panic(err) }
	fmt.Println("Resp: ", resp)

	time.Sleep(1 * time.Second)
}
