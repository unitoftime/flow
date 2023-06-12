package cod

import (
	"errors"
	"math"
	"encoding/binary"

	binreflection "github.com/unitoftime/binary"
)

var ErrUvarintCorrupted = errors.New("failed to read uvarint is corrupted")
var ErrVarintCorrupted = errors.New("failed to read varint is corrupted")

type Buffer struct {
	buf []byte
	off int
}
func NewBuffer(length int) *Buffer {
	return &Buffer{
		buf: make([]byte, 0, length),
		// endian: binary.LittleEndian, // TODO: Made things a bit slower. retest at scale
	}
}

func NewBufferFrom(dat []byte) *Buffer {
	return &Buffer{
		buf: dat,
		// endian: binary.LittleEndian, // TODO: Made things a bit slower. retest at scale
	}
}

// // TODO: Remove I just used this for testing benchmarks
// func (b *Buffer) Realloc(length int) {
// 	b.off = 0
// 	b.buf = make([]byte, 0, length)
// }

func (b *Buffer) Seek(offset int) {
	b.off = offset
}

func (b *Buffer) Reset() {
	b.off = 0
	b.buf = b.buf[:0]
}

// Retursn all of the bytes in the buffer
func (b *Buffer) Bytes() []byte {
	return b.buf
}

// Write Uint values
func (b *Buffer) WriteUint8(v uint8) *Buffer {
	b.buf = append(b.buf, v)
	return b
}
func (b *Buffer) WriteRawUint16(v uint16) *Buffer {
	b.buf = order.AppendUint16(b.buf, v)
	return b
}
func (b *Buffer) WriteRawUint32(v uint32) *Buffer {
	b.buf = order.AppendUint32(b.buf, v)
	return b
}
func (b *Buffer) WriteRawUint64(v uint64) *Buffer {
	b.buf = order.AppendUint64(b.buf, v)
	return b
}

// Write Int values
func (b *Buffer) WriteRawInt8(v int8) *Buffer {
	return b.WriteUint8(uint8(v))
}
func (b *Buffer) WriteRawInt16(v int16) *Buffer {
	return b.WriteRawUint16(uint16(v))
}
func (b *Buffer) WriteRawInt32(v int32) *Buffer {
	return b.WriteRawUint32(uint32(v))
}
func (b *Buffer) WriteRawInt64(v int64) *Buffer {
	return b.WriteRawUint64(uint64(v))
}

// Read Uint values
func (b *Buffer) ReadUint8() uint8 {
	ret := b.buf[b.off]
	b.off += 1
	return ret
}
func (b *Buffer) ReadRawUint16() uint16 {
	ret := order.Uint16(b.buf[b.off:])
	b.off += 2
	return ret
}
func (b *Buffer) ReadRawUint32() uint32 {
	ret := order.Uint32(b.buf[b.off:])
	b.off += 4
	return ret
}
func (b *Buffer) ReadRawUint64() uint64 {
	ret := order.Uint64(b.buf[b.off:])
	b.off += 8
	return ret
}

// Read Int values
func (b *Buffer) ReadInt8() int8 {
	return int8(b.ReadUint8())
}
func (b *Buffer) ReadRawInt16() int16 {
	return int16(b.ReadRawUint16())
}
func (b *Buffer) ReadRawInt32() int32 {
	return int32(b.ReadRawUint32())
}
func (b *Buffer) ReadRawInt64() int64 {
	return int64(b.ReadRawUint64())
}

// Floats
func (b *Buffer) WriteFloat32(v float32) *Buffer {
	return b.WriteRawUint32(math.Float32bits(v))
}
func (b *Buffer) WriteFloat64(v float64) *Buffer {
	return b.WriteRawUint64(math.Float64bits(v))
}

func (b *Buffer) ReadFloat32() float32 {
	ret := math.Float32frombits(b.ReadRawUint32())
	return ret
}
func (b *Buffer) ReadFloat64() float64 {
	ret := math.Float64frombits(b.ReadRawUint64())
	return ret
}

// Variable width and floats
func (b *Buffer) WriteUint16(v uint16) *Buffer {
	b.buf = binary.AppendUvarint(b.buf, uint64(v))
	return b
}
func (b *Buffer) WriteUint32(v uint32) *Buffer {
	b.buf = binary.AppendUvarint(b.buf, uint64(v))
	return b
}
func (b *Buffer) WriteUint64(v uint64) *Buffer {
	b.buf = binary.AppendUvarint(b.buf, v)
	return b
}
func (b *Buffer) ReadUint16() (uint16, error) {
	v, err := b.ReadUint64()
	return uint16(v), err
}
func (b *Buffer) ReadUint32() (uint32, error) {
	v, err := b.ReadUint64()
	return uint32(v), err
}
func (b *Buffer) ReadUint64() (uint64, error) {
	ret, n := binary.Uvarint(b.buf[b.off:])
	if n <= 0 {
		return 0, ErrUvarintCorrupted
	}
	b.off += n
	return ret, nil
}


func (b *Buffer) WriteInt16(v int16) *Buffer {
	b.buf = binary.AppendVarint(b.buf, int64(v))
	return b
}
func (b *Buffer) WriteInt32(v int32) *Buffer {
	b.buf = binary.AppendVarint(b.buf, int64(v))
	return b
}
func (b *Buffer) WriteInt64(v int64) *Buffer {
	b.buf = binary.AppendVarint(b.buf, v)
	return b
}

func (b *Buffer) ReadInt16() (int16, error) {
	v, err := b.ReadInt64()
	return int16(v), err
}
func (b *Buffer) ReadInt32() (int32, error) {
	v, err := b.ReadInt64()
	return int32(v), err
}
func (b *Buffer) ReadInt64() (int64, error) {
	ret, n := binary.Varint(b.buf[b.off:])
	if n <= 0 {
		return 0, ErrVarintCorrupted
	}
	b.off += n
	return ret, nil
}

// --------------------------------------------------------------------------------
// - Complex types
// --------------------------------------------------------------------------------
// This copies the data
func (b *Buffer) WriteByteSlice(v []byte) *Buffer {
	b.WriteInt64(int64(len(v)))
	b.buf = append(b.buf, v...)
	return b
}

// This does not copy the data
func (b *Buffer) ReadByteSlice() ([]byte, error) {
	v, err := b.ReadInt64()
	if err != nil { return nil, err }
	if v < 0 { return nil, ErrVarintCorrupted }

	ret := b.buf[b.off:b.off+int(v)]
	b.off += int(v)
	return ret, nil
}

func (b *Buffer) WriteString(v string) *Buffer {
	return b.WriteByteSlice([]byte(v))
}
func (b *Buffer) ReadString() (string, error) {
	dat, err := b.ReadByteSlice()
	if err != nil { return "", err }

	return string(dat), nil
}

func (b *Buffer) WriteCod(v Er) *Buffer {
	v.EncodeCod(b)
	return b
}
func (b *Buffer) ReadCod(v Er) error {
	return v.DecodeCod(b)
}

func (b *Buffer) WriteAny(v any) *Buffer {
	coder, ok := v.(Er)
	if ok {
		return b.WriteCod(coder)
	}

	dat, err := binreflection.Marshal(v)
	if err != nil {
		panic(err)
	}

	return b.WriteByteSlice(dat)
}
// Must be pointer to underlying type
func (b *Buffer) ReadAny(p any) error {
	coder, ok := p.(Er)
	if ok {
		return b.ReadCod(coder)
	}

	dat, err := b.ReadByteSlice()
	if err != nil { return err }

	return binreflection.Unmarshal(dat, p)
}

// The cod.Er interface
type Er interface {
	EncodeCod(*Buffer) // Encode data to the buffer
	DecodeCod(*Buffer) error // Decode data from the buffer
}
