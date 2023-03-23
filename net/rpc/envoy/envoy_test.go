package envoy

import (
	"fmt"
	"time"
	"testing"

	"github.com/unitoftime/flow/net"
)

// Define Services
type TestServiceApiStruct struct {
	DoThing RpcDef[ServerReq, ServerResp]
	HandleMsg MsgDef[ServerMsg]
}

type TestServiceClientApiStruct struct {
	ClientDoThing RpcDef[ClientReq, ClientResp]
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
func (s *TestService) HandleMsg(r ServerMsg) {
	fmt.Println("HandleMsg: ", r)
}

func (s *TestService) DoThing(r ServerReq) ServerResp {
	fmt.Println("DoThing: ", r)
	return ServerResp{
		r.Val + 1,
	}
}

type ClientReq struct {
	Val uint16
}
type ClientResp struct {
	Val uint16
}

type TestServiceClient struct {
}
func (s *TestServiceClient) ClientDoThing(r ClientReq) ClientResp {
	fmt.Println("ClientDoThing: ", r)
	return ClientResp{
		r.Val + 100,
	}
}

func TestMain(t *testing.T) {
	interfaceDef := NewInterfaceDef[TestServiceApiStruct, TestServiceClientApiStruct]()

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

			client := interfaceDef.NewServer()

			ts := &TestService{}
			client.Handler.DoThing.Register(ts.DoThing)
			client.Handler.HandleMsg.Register(ts.HandleMsg)
			// call := NewCall(client, client.Call.ClientDoThing)

			client.Connect(sock)

			resp, err := client.Call.ClientDoThing.Call(ClientReq{1})
			// resp, err := call.Do(ClientReq{1})
			if err != nil { panic(err) }
			fmt.Println("ClientDoThingResp: ", resp)
		}
	}()

	dialConfig := net.DialConfig{
		Url: "tcp://localhost:8000",
	}
	sock := dialConfig.Dial()

	ts := &TestServiceClient{}

	client := interfaceDef.NewClient()
	client.Handler.ClientDoThing.Register(ts.ClientDoThing)
	client.Connect(sock)

	// call := NewCall(client, client.Call.DoThing)
	msgCall := NewMessage(client, client.Call.HandleMsg)

	fmt.Println(client)

	time.Sleep(1 * time.Second)

	err := msgCall.Send(ServerMsg{9})
	if err != nil { panic(err) }

	resp, err := client.Call.DoThing.Call(ServerReq{5})
	// resp := client.MakeRequest(Req{5})
	// resp, err := call.Do(ServerReq{5})
	if err != nil { panic(err) }
	fmt.Println("Resp: ", resp)

	time.Sleep(1 * time.Second)
}
