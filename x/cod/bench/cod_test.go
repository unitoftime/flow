package bench

import (
	"testing"
	"math/rand"

	"github.com/unitoftime/flow/cod"
)

var musCodec = cod.NewCodec(encodeCod, decodeCod)
func encodeCod(buf *cod.Buffer, p MUSA) {
	buf.WriteString(p.Name)
	buf.WriteInt64(p.BirthDay)
	buf.WriteString(p.Phone)
	// buf.WriteVarint(int64(p.Siblings))
	buf.WriteBool(p.Spouse)
	buf.WriteInt32(p.Siblings)
	buf.WriteFloat64(p.Money)
}
func decodeCod(buf *cod.Buffer) (p MUSA, err error) {
	p.Name, err = buf.ReadString()
	if err != nil { return }
	p.BirthDay, err = buf.ReadInt64()
	if err != nil { return }
	p.Phone, err = buf.ReadString()
	if err != nil { return }
	p.Spouse = buf.ReadBool()
	// s, err := buf.ReadVarint()
	// if err != nil { return }
	// p.Siblings = int32(s)
	p.Siblings, err = buf.ReadInt32()
	if err != nil { return }
	p.Money = buf.ReadFloat64()

	return
}

func Benchmark_cod_Marshal(b *testing.B) {
	data := generateMUS()
	b.ReportAllocs()
	var serialSize int
	var buf *cod.Buffer
	var o MUSA
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		o = data[rand.Intn(len(data))]
		buf = cod.NewBuffer(34)

		encodeCod(buf, o)

		serialSize += len(buf.Bytes())
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

func Benchmark_cod_Unmarshal(b *testing.B) {
	b.StopTimer()
	data := generateMUS()
	ser := make([]*cod.Buffer, len(data))
	var serialSize int
	for i, d := range data {
		buf := cod.NewBuffer(0)
		encodeCod(buf, d)
		ser[i] = buf
		serialSize += len(ser[i].Bytes())
	}
	b.ReportMetric(float64(serialSize)/float64(len(data)), "B/serial")
	b.ReportAllocs()
	b.StartTimer()

	for i := 0; i < b.N; i++ {
		n := rand.Intn(len(ser))
		ser[n].Seek(0) // Just have to do this b/c we are rereading from buffers
		o, err := decodeCod(ser[n])
		if err != nil {
			b.Fatalf("mus failed to unmarshal: %s (%d)", err, n)
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
