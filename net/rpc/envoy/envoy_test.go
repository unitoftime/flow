package envoy

import (
	"fmt"
	"time"
	"testing"

	"github.com/unitoftime/flow/net"
	// "github.com/unitoftime/flow/net/rpc"
)

type ServerService interface {
	DoThing(ServerReq) (ServerResp, error)
	HandleMsg(ServerMsg) error
}

type ClientService interface {
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
func HandleMsg(r ServerMsg) error {
	fmt.Println("HandleMsg: ", r)
	return nil
}

func DoThing(r ServerReq) (ServerResp, error) {
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
func ClientDoThing(r ClientReq) (ClientResp, error) {
	fmt.Println("ClientDoThing: ", r)
	return ClientResp{
		r.Val + 100,
	}, nil
}

func TestMain(t *testing.T) {
	// definitions:
	serverDef := NewServiceDef(new(ServerService))
	clientDef := NewServiceDef(new(ClientService))

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

			client := NewClient(sock, serverDef, clientDef)

			Register(client, DoThing)
			RegisterMessage(client, HandleMsg)
			call := NewCall[ClientReq, ClientResp](client)
			fmt.Println(client)

			resp, err := call.Do(ClientReq{1})
			if err != nil { panic(err) }
			fmt.Println("Resp: ", resp)
		}
	}()

	dialConfig := net.DialConfig{
		Url: "tcp://localhost:8000",
	}
	sock := dialConfig.Dial()

	client := NewClient(sock, clientDef, serverDef)
	Register(client, ClientDoThing)
	call := NewCall[ServerReq, ServerResp](client)
	msgCall := NewMessage[ServerMsg](client)
	fmt.Println(client)

	time.Sleep(1 * time.Second)

	msgCall.Send(ServerMsg{9})

	// resp := client.MakeRequest(Req{5})
	resp, err := call.Do(ServerReq{5})
	if err != nil { panic(err) }
	fmt.Println("Resp: ", resp)

	time.Sleep(1 * time.Second)
}
