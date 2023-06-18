package cod

import (
	"testing"

	// "bytes"
	// "encoding/binary"

	// crunch "github.com/superwhiskers/crunch/v3"
	// "github.com/zhuangsirui/binpacker"
)

//--------------------------------------------------------------------------------
// 1. Basic Tests
//--------------------------------------------------------------------------------
type testData struct {
	A, B, C uint16
}
func newData() testData {
	return testData{
		A: 1,
		B: 2,
		C: 3,
	}
}

func encodeTestData(bs []byte, v testData) []byte {
	bs = AppendUint16(bs, v.A)
	bs = AppendUint16(bs, v.B)
	bs = AppendUint16(bs, v.C)
	return bs
}
func encodeTestDataVarint(bs []byte, v testData) []byte {
	bs = AppendUvarint(bs, uint64(v.A))
	bs = AppendUvarint(bs, uint64(v.B))
	bs = AppendUvarint(bs, uint64(v.C))
	return bs
}

func decodeTestData(bs []byte) (testData, error) {
	res := testData{}
	res.A, bs = Uint16(bs)
	res.B, bs = Uint16(bs)
	res.C, bs = Uint16(bs)
	return res, nil

	// idx := 0
	// n1 := 0
	// res := testData{}
	// res.A, n1 = Uint16(bs[idx:])
	// idx += n1
	// res.B, n1 = Uint16(bs[idx:])
	// idx += n1
	// res.C, n1 = Uint16(bs[idx:])
	// return res, nil
}

func encodeTestDataBuffer(buf *Buffer, d testData) {
	buf.WriteUint16(d.A).WriteUint16(d.B).WriteUint16(d.C)
}

func decodeTestDataBuffer(buf *Buffer) (testData, error) {
	res := testData{}
	var err error
	res.A, err = buf.ReadUint16()
	if err != nil { panic(err) }
	res.B, _ = buf.ReadUint16()
	if err != nil { panic(err) }
	res.C, _ = buf.ReadUint16()
	if err != nil { panic(err) }
	return res, nil
}

func TestBasic(t *testing.T) {
	d := newData()
	bs := make([]byte, 0, 1024)
	bs = bs[:0]
	bs = encodeTestData(bs, d)
	t.Log("Serialized: ", bs)


	// n := 0
	// n1 := 0
	// res := testData{}
	// res.A, n1 = Uint16(bs[n:])
	// n += n1
	// res.B, n1 = Uint16(bs[n:])
	// n += n1
	// res.C, n1 = Uint16(bs[n:])
	res, _ := decodeTestData(bs)
	t.Log("Deserializ: ", res)
}


// func BenchmarkBasic(b *testing.B) {
// 	d := newData()
// 	bs := make([]byte, 0, 1024)

// 	b.ResetTimer()

// 	for n := 0; n < b.N; n++ {
// 		bs = bs[:0]
// 		bs = encodeTestData(bs, d)

// 		_, err := decodeTestData(bs)
// 		if err != nil {
// 			panic(err)
// 		}
// 	}

// 	// b.Log(bs)
// }

// cpu: Intel(R) Core(TM) i7-8700K CPU @ 3.70GHz
// BenchmarkBasic2
// BenchmarkBasic2-12       	1000000000	         0.9900 ns/op	       0 B/op	       0 allocs/op
// BenchmarkBasic
// BenchmarkBasic-12        	1000000000	         0.9247 ns/op	       0 B/op	       0 allocs/op
// BenchmarkCrunch
// BenchmarkCrunch-12       	139578205	         8.614 ns/op	       0 B/op	       0 allocs/op
// BenchmarkBinPacker
// BenchmarkBinPacker-12    	21341906	        54.96 ns/op	       6 B/op	       3 allocs/op

// func BenchmarkCrunch(b *testing.B) {
// 	d := newData()
// 	buf := crunch.NewBuffer()
// 	buf.Grow(2 * 3)

// 	b.ResetTimer()

// 	for n := 0; n < b.N; n++ {
// 		buf.SeekByte(0x00, false)
// 		buf.WriteU16LENext([]uint16{d.A})
// 		buf.WriteU16LENext([]uint16{d.B})
// 		buf.WriteU16LENext([]uint16{d.C})
// 	}

// 	// b.Log(buf.Bytes())
// }

// func BenchmarkBinPacker(b *testing.B) {
// 	d := newData()
// 	buffer := new(bytes.Buffer)
// 	buffer.Grow(1024)
// 	packer := binpacker.NewPacker(binary.BigEndian, buffer)

// 	b.ResetTimer()

// 	for n := 0; n < b.N; n++ {
// 		buffer.Reset()
// 		packer.PushUint16(d.A)
// 		packer.PushUint16(d.B)
// 		packer.PushUint16(d.C)
// 	}

// 	// b.Log(buffer.Bytes())
// }

//--------------------------------------------------------------------------------
// 2. Serialize Structs
//--------------------------------------------------------------------------------
func BenchmarkStructManual(b *testing.B) {
	d := newData()
	bs := make([]byte, 0, 1024)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		bs = bs[:0]
		// bs = encodeTestDataVarint(bs, d)
		bs = encodeTestData(bs, d)

		res, err := decodeTestData(bs)
		if err != nil { panic(err) }
		if d != res { panic("dne") }
	}
}

// func BenchmarkStructPacker(b *testing.B) {
// 	d := newData()
// 	buffer := new(bytes.Buffer)
// 	buffer.Grow(1024)
// 	packer := NewPacker(buffer)

// 	b.ResetTimer()

// 	for n := 0; n < b.N; n++ {
// 		buffer.Reset()
// 		packer.PushUint16(d.A)
// 		packer.PushUint16(d.B)
// 		packer.PushUint16(d.C)
// 	}

// 	// b.Log(buffer.Bytes())
// }

func BenchmarkStructBuffer(b *testing.B) {
	d := newData()
	buffer := NewBuffer(1024)

	// codec := NewCodec(encodeTestDataBuffer, decodeTestDataBuffer)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buffer.Reset()
		// buffer.PushUint16(d.A).PushUint16(d.B).PushUint16(d.C)
		// res, err := decodeTestDataBuffer(buffer)
		// if err != nil { panic(err) }
		// if d != res { panic("dne") }

		// codec.Write(buffer, d)
		// res, err := codec.Read(buffer)
		encodeTestDataBuffer(buffer, d)
		res, err := decodeTestDataBuffer(buffer)
		if err != nil { panic(err) }
		if d != res { panic("dne") }
	}

	// b.Log(buffer.Bytes())
}

// func BenchmarkStructCod(b *testing.B) {
// 	d := newData()
// 	bs := make([]byte, 0, 1024)

// 	codec := NewCodec(encodeTestData, decodeTestData)

// 	b.ResetTimer()

// 	for n := 0; n < b.N; n++ {
// 		bs = bs[:0]
// 		bs = codec.Append(bs, d)

// 		res, err := codec.Get(bs)
// 		if err != nil { panic(err) }
// 		if d != res { panic("dne") }
// 	}
// }


//--------------------------------------------------------------------------------
// 3. Serialize Slices of Struct
//--------------------------------------------------------------------------------
func BenchmarkSliceManual(b *testing.B) {
	count := 1000
	d := make([]testData, count)
	for i := range d {
		d[i] = newData()
	}

	bs := make([]byte, 0)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		bs = bs[:0]

		bs = AppendUint16(bs, uint16(len(d)))
		for i := range d {
			bs = encodeTestData(bs, d[i])
		}

		res := make([]testData, count)
		var l uint16
		l, bs = Uint16(bs)
		length := int(l)
		for i := 0; i < length; i++ {
			val, err := decodeTestData(bs)
			if err != nil { panic(err) }

			res[i] = val
		}

		if len(res) != len(d) { panic("DNE") }
		for i := range res {
			if res[i] != d[i] { panic("DNE") }
		}
	}

	// b.Log(bs)
}

func BenchmarkSliceCod(b *testing.B) {
	count := 1000
	d := make([]testData, count)
	for i := range d {
		d[i] = newData()
	}

	buffer := NewBuffer(0)
	codec := NewCodec(encodeTestDataBuffer, decodeTestDataBuffer)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buffer.Reset()
		codec.WriteSlice(buffer, d)

		res, err := codec.ReadSlice(buffer)
		if err != nil { panic(err) }

		if len(res) != len(d) { panic("DNE") }
		for i := range res {
			if res[i] != d[i] { panic("DNE") }
		}
	}

	// b.Log(bs)
}

// func BenchmarkSliceCodOLD(b *testing.B) {
// 	count := 1000
// 	d := make([]testData, count)
// 	for i := range d {
// 		d[i] = newData()
// 	}
// 	codec := NewCodec2(encodeTestData, decodeTestData)

// 	bs := make([]byte, 0, 1024)

// 	b.ResetTimer()

// 	for n := 0; n < b.N; n++ {
// 		bs = bs[:0]
// 		bs = codec.AppendSlice(bs, d)
// 	}

// 	// b.Log(bs)
// }


//--------------------------------------------------------------------------------
// 4. Hard Struct
//--------------------------------------------------------------------------------
type pos struct {
	X, Y float64
}
var posCodec = NewCodec(
	func(buf *Buffer, p pos) {
		buf.WriteFloat64(p.X)
		buf.WriteFloat64(p.Y)
	},
	func(buf *Buffer) (p pos, err error) {
		p.X = buf.ReadFloat64()
		p.Y = buf.ReadFloat64()
		return
	},
)

type testUnion struct {
	A *uint16
	B *uint16
}
// var testUnionCodec = NewCodec()
func encodeTestUnion(buf *Buffer, p testUnion) {
	if p.A != nil {
		buf.WriteUint8(0) // Header
		buf.WriteUint16(*p.A)
	} else if p.B != nil {
		buf.WriteUint8(1) // Header
		buf.WriteUint16(*p.B)
	} else {
		panic("both items nil in union")
	}
}
func decodeTestUnion(buf *Buffer) (p testUnion, err error) {
	header := buf.ReadUint8()
	if header == 0 {
		v, _ := buf.ReadUint16()
		p.A = &v
	} else if header == 1 {
		v, _ := buf.ReadUint16()
		p.B = &v
	}
	return
}


// type testTable struct {
// 	thingA *uint16
// 	thingB *uint16
// }

type hardData struct {
	Dat []uint32
	Str []string
	Something testData
	Mapping map[string]pos
	// Anything any
	Union testUnion
}
func newHardData() hardData {
	a := uint16(5)

	return hardData{
		Dat: []uint32{10, 20, 30},
		Str: []string{"hello world", "ok"},
		Something: testData{
			A: 1,
			B: 2,
			C: 3,
		},
		Mapping: map[string]pos{
			"a": pos{1, 2},
			"b": pos{3, 4},
		},
		// anything: testData{
		// 	A: 10000,
		// 	B: 1000,
		// 	C: 100,
		// },
		Union: testUnion{
			A: &a,
		},

	}
}

var codecTestData = NewCodec(encodeTestDataBuffer, decodeTestDataBuffer)
var mapCodec = NewMapCodec(CodecString, posCodec)
func encodeHard(buf *Buffer, d hardData) {
	CodecUint32.WriteSlice(buf, d.Dat)
	CodecString.WriteSlice(buf, d.Str)
	codecTestData.Write(buf, d.Something)
	mapCodec.Write(buf, d.Mapping)

	// buf.WriteAny(d.Anything)

	// testUnionCodec.Write(buf, d.Union)
	encodeTestUnion(buf, d.Union)
}

func decodeHard(buf *Buffer) (hardData, error) {
	res := hardData{}
	var err error
	// res.Dat, err = CodecUint32.ReadSlice(buf)
	res.Dat, err = CodecUint32.ReadSlice(buf)
	if err != nil { panic(err) }

	res.Str, err = CodecString.ReadSlice(buf)
	if err != nil { panic(err) }

	res.Something, err = codecTestData.Read(buf)
	if err != nil { panic(err) }

	res.Mapping, err = mapCodec.Read(buf)
	if err != nil { panic(err) }

	// td := testData{}
	// err = buf.ReadAny(&td)
	// if err != nil { panic(err) }
	// res.Anything = td

	res.Union, err = decodeTestUnion(buf)
	if err != nil { panic(err) }

	return res, err
}

func BenchmarkStructOfStruct(b *testing.B) {
	d := newHardData()
	res := hardData{}
	var err error

	buffer := NewBuffer(0)
	serialSize := 0

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buffer.Reset()

		encodeHard(buffer, d)
		res, err = decodeHard(buffer)
		if err != nil { panic(err) }
		serialSize += len(buffer.Bytes())
	}

	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
	// b.Log(buffer.Bytes())
	b.Log(res)
}


func BenchmarkStructOfStructAnyEncoding(b *testing.B) {
	d := newHardData()
	res := hardData{}
	var err error
	serialSize := 0

	buffer := NewBuffer(0)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buffer.Reset()

		buffer.WriteAny(d)

		res = hardData{}
		err = buffer.ReadAny(&res)
		if err != nil { panic(err) }

		serialSize += len(buffer.Bytes())
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
	// b.Log(buffer.Bytes())
	b.Log(res)
}

//--------------------------------------------------------------------------------
func (d hardData) EncodeCod(buf *Buffer) {
	encodeHard(buf, d)
}

func (d *hardData) DecodeCod(buf *Buffer) error {
	hd, err := decodeHard(buf)
	if err != nil { return err }
	*d = hd
	return nil
}

func BenchmarkInterfaceImpl(b *testing.B) {
	d := newHardData()
	res := hardData{}
	serialSize := 0

	buffer := NewBuffer(0)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		buffer.Reset()
		buffer.WriteCod(d)

		err := buffer.ReadCod(&res)
		if err != nil { panic(err) }
		serialSize += len(buffer.Bytes())
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
	// b.Log(buffer.Bytes())
	b.Log(res)
}

// type lookup struct {
// 	enc map[reflect.Type]func(*Buffer)
// 	dec map[reflect.Type]func(*Buffer) error
// }

// func BenchmarkInterfaceMapLookup(b *testing.B) {
// 	encMap := make(map[reflect.Type]func(*Buffer))
// 	decMap := make(map[reflect.Type]func(*Buffer) error)
// 	l := lookup{
// 		enc: encMap,
// 		dec: decMap,
// 	}

// 	d := newHardData()

// 	buffer := NewBuffer(0)

// 	b.ResetTimer()

// 	for n := 0; n < b.N; n++ {
// 		buffer.Reset()
// 		buffer.WriteCod(d)

// 		res := hardData{}
// 		buffer.ReadCod(&res)
// 	}
// }
