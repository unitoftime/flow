package net

import (
	"fmt"
	"reflect"
)

type MuxHandler func(any) error

type Mux struct {
	handlers map[reflect.Type]MuxHandler
}
func NewMux() *Mux {
	return &Mux{
		handlers: make(map[reflect.Type]MuxHandler),
	}
}

func Register[A any](mux *Mux, handler func(A) error) {
	var msgVal A
	msgValType := reflect.TypeOf(msgVal)
	_, exists := mux.handlers[msgValType]
	if exists {
		panic("Cant reregister the same handler type")
	}

	// Create a handler function
	generalHandlerFunc := func(anyMsg any) error {
		msg, ok := anyMsg.(A)
		if !ok {
			panic(fmt.Errorf("Mismatched request types: %T, %T", anyMsg, msgVal))
		}

		return handler(msg)
	}

	// Store the handler function
	mux.handlers[msgValType] = generalHandlerFunc
}

func (m *Mux) Handle(msg any) (bool, error) {
	msgType := reflect.TypeOf(msg)
	handler, ok := m.handlers[msgType]
	if !ok {
		return false, nil
	}

	err := handler(msg)
	return true, err
}
