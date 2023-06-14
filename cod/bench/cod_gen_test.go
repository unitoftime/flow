package bench

import (
	"testing"
	"math/rand"

	"github.com/unitoftime/flow/cod"
)

//go:generate go run ../gen

func Benchmark_cod_gen_Marshal(b *testing.B) {
	data := generateMUS()
	b.ReportAllocs()
	var serialSize int
	var buf *cod.Buffer
	var o MUSA
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		o = data[rand.Intn(len(data))]
		buf = cod.NewBuffer(34)

		o.EncodeCod(buf)

		serialSize += len(buf.Bytes())
	}
	b.ReportMetric(float64(serialSize)/float64(b.N), "B/serial")
}

func Benchmark_cod_gen_Unmarshal(b *testing.B) {
	b.StopTimer()
	data := generateMUS()
	ser := make([]*cod.Buffer, len(data))
	var serialSize int
	for i, d := range data {
		buf := cod.NewBuffer(0)
		d.EncodeCod(buf)

		ser[i] = buf
		serialSize += len(ser[i].Bytes())
	}
	b.ReportMetric(float64(serialSize)/float64(len(data)), "B/serial")
	b.ReportAllocs()
	b.StartTimer()

	o := MUSA{}
	for i := 0; i < b.N; i++ {
		n := rand.Intn(len(ser))
		ser[n].Seek(0) // Just have to do this b/c we are rereading from buffers

		err := o.DecodeCod(ser[n])
		// o, err := decodeCod(ser[n])
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
