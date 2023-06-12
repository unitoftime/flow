package cod

import (
	"encoding"
	"encoding/binary"
)

// TODOs:
// 1. Measure varint encoding vs fixed int encoding

type Marshaler = encoding.BinaryMarshaler

var order = binary.LittleEndian

// Append values
func AppendUint8(bs []byte, v uint8) []byte {
	return append(bs, v)
}

func AppendUint16(bs []byte, v uint16) []byte {
	return order.AppendUint16(bs, v)
}

func AppendUint32(bs []byte, v uint32) []byte {
	return order.AppendUint32(bs, v)
}

func AppendUint64(bs []byte, v uint64) []byte {
	return order.AppendUint64(bs, v)
}

func AppendUvarint(bs []byte, v uint64) []byte {
	return binary.AppendUvarint(bs, v)
}
func AppendVarint(bs []byte, v int64) []byte {
	return binary.AppendVarint(bs, v)
}

// Read values
func Uint8(bs []byte) uint8 {
	return bs[0]
}
func Uint16(bs []byte) (uint16, []byte) {
	return order.Uint16(bs), bs[2:]
}
func Uint32(bs []byte) uint32 {
	return order.Uint32(bs)
}
func Uint64(bs []byte) uint64 {
	return order.Uint64(bs)
}

// Standard Codecs
var CodecUint8 = NewCodec(encodeUint8, decodeUint8)
func encodeUint8(buf *Buffer, v uint8) {
	buf.WriteUint8(v)
}
func decodeUint8(buf *Buffer) (uint8, error) {
	return buf.ReadUint8(), nil
}

var CodecUint16 = NewCodec(encodeUint16, decodeUint16)
func encodeUint16(buf *Buffer, v uint16) {
	buf.WriteUint16(v)
}
func decodeUint16(buf *Buffer) (uint16, error) {
	return buf.ReadUint16()
}

var CodecUint32 = NewCodec(encodeUint32, decodeUint32)
func encodeUint32(buf *Buffer, v uint32) {
	buf.WriteUint32(v)
}
func decodeUint32(buf *Buffer) (uint32, error) {
	return buf.ReadUint32()
}

var CodecUint64 = NewCodec(encodeUint64, decodeUint64)
func encodeUint64(buf *Buffer, v uint64) {
	buf.WriteUint64(v)
}
func decodeUint64(buf *Buffer) (uint64, error) {
	return buf.ReadUint64()
}

// More complex types

// Simply just converting to byte slices
var CodecString = NewCodec(encodeString, decodeString)
func encodeString(buf *Buffer, v string) {
	buf.WriteString(v)
}
func decodeString(buf *Buffer) (string, error) {
	return buf.ReadString()
}
