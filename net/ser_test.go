package net

import (
	"testing"
	"math"
	"bytes"
	"encoding"
	"encoding/binary"


	unitbinary "github.com/unitoftime/binary"
)

// This one has a custom marshaller
type SerializeTest2 struct {
	X, Y, Z float64
	A uint32
	// Other []int32
}
func newData2() SerializeTest2 {
	return SerializeTest2{
		X: 1.1,
		Y: 2.1,
		Z: 3.1,
		A: 4,
		// Other: []int32{1, 2, 3, 4, 5},
	}
}

func (s *SerializeTest2) MarshalBinary() (data []byte, err error) {
	ser := make([]byte, 28)
	binary.LittleEndian.PutUint64(ser[0:], math.Float64bits(s.X))
	binary.LittleEndian.PutUint64(ser[8:], math.Float64bits(s.Y))
	binary.LittleEndian.PutUint64(ser[16:], math.Float64bits(s.Z))
	binary.LittleEndian.PutUint32(ser[24:], s.A)

	return ser, nil
}

func (s *SerializeTest2) UnmarshalBinary(data []byte) error {
	s.X = math.Float64frombits(binary.LittleEndian.Uint64(data[0:]))
	s.Y = math.Float64frombits(binary.LittleEndian.Uint64(data[8:]))
	s.Z = math.Float64frombits(binary.LittleEndian.Uint64(data[16:]))
	s.A = binary.LittleEndian.Uint32(data[24:])

	return nil
}

func serializeWhole2(s SerializeTest2) []byte {
	ser := new(bytes.Buffer)

	err := binary.Write(ser, binary.LittleEndian, &s)
	if err != nil { panic(err) }
	return ser.Bytes()
}

func deserializeWhole2(data []byte) SerializeTest2 {
	des := SerializeTest2{}

	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &des)
	if err != nil { panic(err) }

	return des
}

func serializeWholeCustomBinary2(s SerializeTest2) []byte {
	ser, err := unitbinary.Marshal(&s)
	if err != nil { panic(err) }
	return ser
}

func deserializeWholeCustomBinary2(data []byte) SerializeTest2 {
	des := SerializeTest2{}
	err := unitbinary.Unmarshal(data, &des)
	if err != nil { panic(err) }

	return des
}

func serializeWholeInterfaceCall(s encoding.BinaryMarshaler) []byte {
	ser, err := s.MarshalBinary()
	if err != nil { panic(err) }

	return ser
}

func deserializeWholeInterfaceCall(data []byte) SerializeTest2 {
	des := SerializeTest2{}
	err := des.UnmarshalBinary(data)
	if err != nil { panic(err) }

	return des
}

// --------------------------------------------------------------------------------

type SerializeTest struct {
	X, Y, Z float64
	A uint32
	// Other []int32
}
func newData() SerializeTest {
	return SerializeTest{
		X: 1.1,
		Y: 2.1,
		Z: 3.1,
		A: 4,
		// Other: []int32{1, 2, 3, 4, 5},
	}
}
func serializeWholeCustomBinary(s SerializeTest) []byte {
	ser, err := unitbinary.Marshal(s)
	if err != nil { panic(err) }
	return ser
}

func deserializeWholeCustomBinary(data []byte) SerializeTest {
	des := SerializeTest{}
	err := unitbinary.Unmarshal(data, &des)
	if err != nil { panic(err) }

	return des
}

func serializeWhole(s SerializeTest) []byte {
	ser := new(bytes.Buffer)

	err := binary.Write(ser, binary.LittleEndian, s)
	if err != nil { panic(err) }
	return ser.Bytes()
}

func deserializeWhole(data []byte) SerializeTest {
	des := SerializeTest{}

	buf := bytes.NewReader(data)
	err := binary.Read(buf, binary.LittleEndian, &des)
	if err != nil { panic(err) }

	return des
}

func serializeManual(s SerializeTest) []byte {
	ser := make([]byte, 28)

	binary.LittleEndian.PutUint64(ser[0:], math.Float64bits(s.X))
	binary.LittleEndian.PutUint64(ser[8:], math.Float64bits(s.Y))
	binary.LittleEndian.PutUint64(ser[16:], math.Float64bits(s.Z))
	binary.LittleEndian.PutUint32(ser[24:], s.A)

	return ser
}

func deserializeManual(data []byte) SerializeTest {
	des := SerializeTest{}

	des.X = math.Float64frombits(binary.LittleEndian.Uint64(data[0:]))
	des.Y = math.Float64frombits(binary.LittleEndian.Uint64(data[8:]))
	des.Z = math.Float64frombits(binary.LittleEndian.Uint64(data[16:]))
	des.A = binary.LittleEndian.Uint32(data[24:])

	return des
}


func BenchmarkSerializeUnion(b *testing.B) {
	d := newData()
	union := NewUnion(SerializeTest{})

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		union.Serialize(d)
	}
}

func BenchmarkDeserializeUnion(b *testing.B) {
	d := newData()
	union := NewUnion(SerializeTest{})
	ser, err := union.Serialize(d)
	if err != nil { panic(err) }

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		union.Deserialize(ser)
	}
}

func BenchmarkSerializeWhole(b *testing.B) {
	d := newData()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		serializeWhole(d)
	}
}

func BenchmarkDeserializeWhole(b *testing.B) {
	d := newData()

	ser := serializeWhole(d)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		deserializeWhole(ser)
	}
}

func BenchmarkSerializeWholeCustomBinary(b *testing.B) {
	d := newData()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		serializeWholeCustomBinary(d)
	}
}

func BenchmarkDeserializeWholeCustomBinary(b *testing.B) {
	d := newData()

	ser := serializeWholeCustomBinary(d)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		deserializeWholeCustomBinary(ser)
	}
}


func BenchmarkSerializeManual(b *testing.B) {
	d := newData()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		serializeManual(d)
	}
}

func BenchmarkDeserializeManual(b *testing.B) {
	d := newData()

	ser := serializeManual(d)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		deserializeManual(ser)
	}
}

// --------------------------------------------------------------------------------
func BenchmarkSerializeWhole2(b *testing.B) {
	d := newData2()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		serializeWhole2(d)
	}
}

func BenchmarkDeserializeWhole2(b *testing.B) {
	d := newData2()

	ser := serializeWhole2(d)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		deserializeWhole2(ser)
	}
}

func BenchmarkSerializeWholeCustomBinary2(b *testing.B) {
	d := newData2()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		serializeWholeCustomBinary2(d)
	}
}

func BenchmarkDeserializeWholeCustomBinary2(b *testing.B) {
	d := newData2()

	ser := serializeWholeCustomBinary2(d)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		deserializeWholeCustomBinary2(ser)
	}
}

func BenchmarkSerializeWholeInterfaceCall(b *testing.B) {
	d := newData2()

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		serializeWholeInterfaceCall(&d)
	}
}

func BenchmarkDeserializeWholeInterfaceCall(b *testing.B) {
	d := newData2()

	ser := serializeWholeInterfaceCall(&d)

	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		deserializeWholeInterfaceCall(ser)
	}
}
