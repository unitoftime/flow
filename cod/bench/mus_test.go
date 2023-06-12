package bench

import (
	"time"
	"testing"
	// "math"
	// "bytes"
	// "encoding"
	// "encoding/binary"

	// "time"
	"math/rand"

	// "github.com/mus-format/mus-go"
	"github.com/mus-format/mus-go/ord"
	"github.com/mus-format/mus-go/raw"
	"github.com/mus-format/mus-go/unsafe"
	"github.com/mus-format/mus-go/varint"

	// unitbinary "github.com/unitoftime/binary"
)

type MUSA struct {
	Name     string
	BirthDay int64
	Phone    string
	Siblings int32
	Spouse   bool
	Money    float64
}

func randString(n int) string {
	return "abcdefg"
}
func generateMUS() []MUSA {
	a := make([]MUSA, 0, 1000)
	for i := 0; i < 1000; i++ {
		a = append(a, MUSA{
			Name:     randString(16),
			BirthDay: time.Now().UnixNano(),
			Phone:    randString(10),
			Siblings: rand.Int31n(5),
			Spouse:   rand.Intn(2) == 1,
			Money:    rand.Float64(),
		})
	}
	return a
}

func MarshalMUS(v MUSA) (buf []byte) {
	n := ord.SizeString(v.Name)
	n += raw.SizeInt64(v.BirthDay)
	n += ord.SizeString(v.Phone)
	n += varint.SizeInt32(v.Siblings)
	n += ord.SizeBool(v.Spouse)
	n += raw.SizeFloat64(v.Money)
	buf = make([]byte, n)
	n = ord.MarshalString(v.Name, buf)
	n += raw.MarshalInt64(v.BirthDay, buf[n:])
	n += ord.MarshalString(v.Phone, buf[n:])
	n += varint.MarshalInt32(v.Siblings, buf[n:])
	n += ord.MarshalBool(v.Spouse, buf[n:])
	raw.MarshalFloat64(v.Money, buf[n:])
	return
}

func UnmarshalMUS(bs []byte) (v MUSA, n int, err error) {
	v.Name, n, err = ord.UnmarshalString(bs)
	if err != nil {
		return
	}
	var n1 int
	v.BirthDay, n1, err = raw.UnmarshalInt64(bs[n:])
	n += n1
	if err != nil {
		return
	}
	v.Phone, n1, err = ord.UnmarshalString(bs[n:])
	n += n1
	if err != nil {
		return
	}
	v.Siblings, n1, err = varint.UnmarshalInt32(bs[n:])
	n += n1
	if err != nil {
		return
	}
	v.Spouse, n1, err = ord.UnmarshalBool(bs[n:])
	n += n1
	if err != nil {
		return
	}
	v.Money, n1, err = raw.UnmarshalFloat64(bs[n:])
	n += n1
	return
}

func MarshalMUSUnsafe(v MUSA) (buf []byte) {
	n := unsafe.SizeString(v.Name)
	n += unsafe.SizeInt64(v.BirthDay)
	n += unsafe.SizeString(v.Phone)
	n += unsafe.SizeInt32(v.Siblings)
	n += unsafe.SizeBool(v.Spouse)
	n += unsafe.SizeFloat64(v.Money)
	buf = make([]byte, n)
	n = unsafe.MarshalString(v.Name, buf)
	n += unsafe.MarshalInt64(v.BirthDay, buf[n:])
	n += unsafe.MarshalString(v.Phone, buf[n:])
	n += unsafe.MarshalInt32(v.Siblings, buf[n:])
	n += unsafe.MarshalBool(v.Spouse, buf[n:])
	unsafe.MarshalFloat64(v.Money, buf[n:])
	return
}

func UnmarshalMUSUnsafe(bs []byte) (v MUSA, n int, err error) {
	v.Name, n, err = unsafe.UnmarshalString(bs)
	if err != nil {
		return
	}
	var n1 int
	v.BirthDay, n1, err = unsafe.UnmarshalInt64(bs[n:])
	n += n1
	if err != nil {
		return
	}
	v.Phone, n1, err = unsafe.UnmarshalString(bs[n:])
	n += n1
	if err != nil {
		return
	}
	v.Siblings, n1, err = unsafe.UnmarshalInt32(bs[n:])
	n += n1
	if err != nil {
		return
	}
	v.Spouse, n1, err = unsafe.UnmarshalBool(bs[n:])
	n += n1
	if err != nil {
		return
	}
	v.Money, n1, err = unsafe.UnmarshalFloat64(bs[n:])
	n += n1
	return
}

func Benchmark_MUS_Marshal(b *testing.B) {
	// data := generateMUS()
	data := generateMUS()
	b.ReportAllocs()
	b.ResetTimer()
	var serialSize int
	var buf []byte
	var o MUSA
	for i := 0; i < b.N; i++ {
		o = data[rand.Intn(len(data))]
		buf = MarshalMUS(o)
		serialSize += len(buf)
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

var validate = ""

func Benchmark_MUS_Unmarshal(b *testing.B) {
	b.StopTimer()
	data := generateMUS()
	ser := make([][]byte, len(data))
	var serialSize int
	for i, d := range data {
		ser[i] = MarshalMUS(d)
		serialSize += len(ser[i])
	}
	b.ReportMetric(float64(serialSize)/float64(len(data)), "B/serial")
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		n := rand.Intn(len(ser))
		o, _, err := UnmarshalMUS(ser[n])
		if err != nil {
			b.Fatalf("mus failed to unmarshal: %s (%s)", err, ser[n])
		}
		// Validate unmarshalled data.
		if validate != "" {
			i := data[n]
			correct := o.Name == i.Name && o.Phone == i.Phone && o.Siblings == i.Siblings && o.Spouse == i.Spouse && o.Money == i.Money && o.BirthDay == i.BirthDay //&& cmpTags(o.Tags, i.Tags) && cmpAliases(o.Aliases, i.Aliases)
			if !correct {
				b.Fatalf("unmarshaled object differed:\n%v\n%v", i, o)
			}
		}
	}
}

func Benchmark_MUSUnsafe_Marshal(b *testing.B) {
	data := generateMUS()
	b.ReportAllocs()
	b.ResetTimer()
	var serialSize int
	var buf []byte
	var o MUSA
	for i := 0; i < b.N; i++ {
		o = data[rand.Intn(len(data))]
		buf = MarshalMUSUnsafe(o)
		serialSize += len(buf)
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

func Benchmark_MUSUnsafe_Unmarshal(b *testing.B) {
	b.StopTimer()
	data := generateMUS()
	ser := make([][]byte, len(data))
	var serialSize int
	for i, d := range data {
		ser[i] = MarshalMUSUnsafe(d)
		serialSize += len(ser[i])
	}
	b.ReportMetric(float64(serialSize)/float64(len(data)), "B/serial")
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		n := rand.Intn(len(ser))
		o, _, err := UnmarshalMUSUnsafe(ser[n])
		if err != nil {
			b.Fatalf("mus failed to unmarshal: %s (%s)", err, ser[n])
		}
		// Validate unmarshalled data.
		if validate != "" {
			i := data[n]
			correct := o.Name == i.Name && o.Phone == i.Phone && o.Siblings == i.Siblings && o.Spouse == i.Spouse && o.Money == i.Money && o.BirthDay == i.BirthDay //&& cmpTags(o.Tags, i.Tags) && cmpAliases(o.Aliases, i.Aliases)
			if !correct {
				b.Fatalf("unmarshaled object differed:\n%v\n%v", i, o)
			}
		}
	}
}

