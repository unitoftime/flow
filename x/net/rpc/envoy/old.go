package envoy

// Old stuff that I just want to keep around for reference
//--------------------------------------------------------------------------------


//--------------------------------------------------------------------------------

// TODO - could I make servicedef generic on the interface type. Then when I register handlers I just pass in a struct which implements the interface?

// TODO - I think I'd prefer this to be based on method name and not based on input argument type
// TODO - reordering the definition, or switching between a message and an RPC will break api compatibility
// func NewServiceDef(def any) serviceDef {
// 	ty := reflect.TypeOf(def).Elem()
// 	fmt.Println(ty)
// 	numMethod := ty.NumMethod()
// 	fmt.Println(numMethod)

// 	requests := make([]any, 0)
// 	responses := make([]any, 0)
// 	for i := 0; i < numMethod; i++ {
// 		method := ty.Method(i)
// 		fmt.Println(method)

// 		numInputs := method.Type.NumIn()
// 		// for j := 0; j < numInputs; j++ {
// 		// 	in := method.Type.In(j)
// 		// 	fmt.Println(in)
// 		// }

// 		if numInputs != 1 { panic("We only support methods of form: func (req) (resp, error) or func (req) error") }
// 		reqType := method.Type.In(0)
// 		reqStruct := reflect.New(reqType).Elem().Interface()
// 		requests = append(requests, reqStruct)

// 		numOutputs := method.Type.NumOut()
// 		// for j := 0; j < numOutputs; j++ {
// 		// 	out := method.Type.Out(j)
// 		// 	fmt.Println(out)
// 		// }

// 		if numOutputs == 1 {
// 			// Rpc doesn't expect a response
// 		} else if numOutputs == 2 {
// 			// Rpc expects a response
// 			respType := method.Type.Out(0)
// 			respStruct := reflect.New(respType).Elem().Interface()
// 			responses = append(responses, respStruct)
// 		} else {
// 			panic("We only support methods of form: func (req) (resp, error) or func (req) error")
// 		}

// 		// TODO - check last argument should be an error
// 	}

// 	return ServiceDefinition{
// 		Requests: net.NewUnion(requests...),
// 		Responses: net.NewUnion(responses...),
// 	}
// }

//--------------------------------------------------------------------------------

// func RegisterMessage[S, C, M any](client *Client[S, C], handler func(M) error) {
// 	var msgVal M
// 	msgValType := reflect.TypeOf(msgVal)
// 	_, exists := client.messageHandlers[msgValType]
// 	if exists {
// 		panic("Cant reregister the same handler type")
// 	}

// 	// Create a handler function
// 	generalHandlerFunc := makeMsgHandler(handler)

// 	// Store the handler function
// 	client.messageHandlers[msgValType] = generalHandlerFunc
// }

// func Register[S, C, Req, Resp any](client *Client[S, C], handler func(Req) (Resp, error)) {
// 	var reqVal Req
// 	reqValType := reflect.TypeOf(reqVal)
// 	_, exists := client.handlers[reqValType]
// 	if exists {
// 		panic("Cant reregister the same handler type")
// 	}

// 	// Create a handler function
// 	generalHandlerFunc := makeRpcHandler(handler)

// 	// Store the handler function
// 	client.handlers[reqValType] = generalHandlerFunc
// }

// --------------------------------------------------------------------------------
// Calls
// // Client - making requests
// // func (c *Client) MakeRequest(req any) (any, error) {
// // 	dat, err := c.reqSerdes.Serialize(req)
// // 	if err != nil { return err }

// // 	reqDat, err := rpcSerdes.Marshal(Request{
// // 		Id: 0, // TODO
// // 		Data: dat,
// // 	})
// // 	if err != nil { return err }

// // 	err = c.sock.Write(reqDat)
// // 	return err
// // }

// // func Rpc[Req, Resp any](f func(Req) (Resp, error))

// // func GetCall[Req, Resp any](client *Client, _ func(Req) (Resp, error)) *Call[Req, Resp] {
// // 	return NewCall[Req, Resp](client)
// // }

// func NewCall[S, C, Req, Resp any](client *Client[S, C], rpc RpcDef[Req, Resp]) *Call[S, C, Req, Resp] {
// 	return &Call[S, C, Req, Resp]{
// 		client: client,
// 		timeout: 5 * time.Second,
// 	}
// }

// // func NewCall[S, C, Req, Resp any](client *Client[S, C]) *Call[S, C, Req, Resp] {
// // 	return &Call[S, C, Req, Resp]{
// // 		client: client,
// // 		timeout: 5 * time.Second,
// // 	}
// // }
// type Call[S, C, Req, Resp any] struct {
// 	client *Client[S, C]
// 	timeout time.Duration
// }


// func (c *Call[S, C, Req, Resp]) Do(req Req) (Resp, error) {
// 	var resp Resp
// 	rpcReq, err := c.client.MakeRequest(req)
// 	if err != nil { return resp, err }

// 	// TODO!!! - check if this ID is already being used, if it is, then use a different one

// 	// Send over socket
// 	reqDat, err := rpcSerdes.Marshal(rpcReq)
// 	if err != nil { return resp, err }

// 	// Make a channel to wait for a response on this Id
// 	// TODO - you need to clean this up on any error
// 	respChan := make(chan Response)
// 	c.client.activeCalls[rpcReq.Id] = respChan

// 	// TODO - retry sending? Or push to a queue to be batch sent?
// 	err = c.client.sock.Write(reqDat)
// 	if err != nil { return resp, err }

// 	select {
// 	case rpcResp := <-respChan:
// 		anyResp, err := c.client.UnmakeResponse(rpcResp)
// 		if err != nil { return resp, err }
// 		resp, ok := anyResp.(Resp)
// 		if !ok { panic("Mismatched type!") }
// 		return resp, nil
// 	case <-time.After(c.timeout):
// 		// TODO - I need to cleanup the channel here
// 		return resp, ErrTimeout
// 	}
// }


//--------------------------------------------------------------------------------
// Message
// func NewMessage[S, C, A any](client *Client[S, C], rpc MsgDef[A]) *Msg[S, C, A] {
// 	return &Msg[S, C, A]{
// 		client: client,
// 	}
// }

// func NewMessage[S, C, A any](client *Client[S, C]) *Msg[S, C, A] {
// 	return &Msg[S, C, A]{
// 		client: client,
// 	}
// }
// type Msg[S, C, A any] struct {
// 	client *Client[S, C]
// }
// func (m *Msg[S, C, A]) Send(req A) error {
// 	rpcMsg, err := m.Make(req)
// 	if err != nil { return err }

// 	// Send over socket
// 	reqDat, err := rpcSerdes.Marshal(rpcMsg)
// 	if err != nil { return err }

// 	err = m.client.sock.Write(reqDat)
// 	if err != nil { return err }

// 	return nil
// }

// func (m *Msg[S, C, A]) Make(req A) (Message, error) {
// 	dat, err := m.client.clientDef.Requests.Serialize(req)

// 	return Message{
// 		Data: dat,
// 	}, err
// }
