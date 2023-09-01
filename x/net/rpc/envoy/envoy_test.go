package envoy

import (
	"fmt"
	"time"
	"testing"

	"github.com/unitoftime/flow/net"

	"github.com/unitoftime/cod/backend"
)

//go:generate cod

// Define Services
type TestServiceApiStruct struct {
	DoThing RpcDef[ServerReq, ServerResp]
	HandleMsg MsgDef[ServerMsg]
}

type TestServiceClientApiStruct struct {
	ClientDoThing RpcDef[ClientReq, ClientResp]
}


//cod:struct
type ServerMsg struct {
	Val uint16
}
//cod:struct
type ServerReq struct {
	Val uint16
}
//cod:struct
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

//cod:struct
type ClientReq struct {
	Val uint16
}
//cod:struct
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

			client := interfaceDef.NewServer(nil)

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

	client := interfaceDef.NewClient(nil)
	client.Handler.ClientDoThing.Register(ts.ClientDoThing)
	client.Connect(sock)

	// call := NewCall(client, client.Call.DoThing)
	// msgCall := NewMessage(client, client.Call.HandleMsg)

	fmt.Println(client)

	time.Sleep(1 * time.Second)

	err := client.Call.HandleMsg.Send(ServerMsg{9})
	// err := msgCall.Send(ServerMsg{9})
	if err != nil { panic(err) }

	resp, err := client.Call.DoThing.Call(ServerReq{5})
	// resp := client.MakeRequest(Req{5})
	// resp, err := call.Do(ServerReq{5})
	if err != nil { panic(err) }
	fmt.Println("Resp: ", resp)

	time.Sleep(1 * time.Second)
}


// --------------------------------------------------------------------------------
// - Cod generated for tests
// Note: weird setup, normally these would be in another file but because the structs are defined in the test file we have to copy them into here
// --------------------------------------------------------------------------------

func (t ServerMsg) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ServerMsg) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}

func (t ServerReq) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ServerReq) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}

func (t ServerResp) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ServerResp) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}

func (t ClientReq) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ClientReq) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}

func (t ClientResp) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ClientResp) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}
