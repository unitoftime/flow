package main

import (
	"testing"

	"bytes"
)

// func codWriteUint8(buf *cod.Buffer, t uint8) {
// 	buf.WriteUint8(t)
// }

// func codReadUint8(buf *cod.Buffer, t *uint8) error {
// 	value, err := buf.WriteUint8(t)
// 	if err != nil { return err }
// 	*t = value
// }

// func TestBasicType(t *testing.T) {
// 	buf := new(bytes.Buffer)

// 	executeBasicMarshal(buf, "t.Age", "uint8")
// 	executeBasicUnmarshal(buf, "t.Age", "uint8")
// 	t.Log(string(buf.Bytes()))
// }

// func TestStructType(t *testing.T) {
// 	buf := new(bytes.Buffer)

// 	executeBasicMarshal(buf, "t.Age", "uint8")
// 	executeBasicUnmarshal(buf, "t.Age", "uint8")
// 	t.Log(string(buf.Bytes()))
// }
