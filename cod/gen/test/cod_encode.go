package main

import (
	"github.com/unitoftime/flow/cod"
)

func (t MyUnion) EncodeCod(buf *cod.Buffer) {

	// codUnion := cod.Union(t)
	// rawVal := codUnion.GetRawValue()
	rawVal := t.Get()
	if rawVal == nil {
		buf.WriteUint8(0) // Zero tag indicates nil
		return
	}

	switch sv := rawVal.(type) {

	case Id:
		buf.WriteUint8(1)
		sv.EncodeCod(buf)
		// <no value>

	case SpecialMap:
		buf.WriteUint8(2)
		sv.EncodeCod(buf)
		// <no value>

	default:
		panic("unknown type placed in union")
	}

}

func (t *MyUnion) DecodeCod(buf *cod.Buffer) error {
	var err error

	// codUnion := cod.Union(*t)

	tagVal := buf.ReadUint8()

	switch tagVal {
	case 0: // Zero tag indicates nil
		return nil

	case 1:
		var decoded Id
		err = decoded.DecodeCod(buf)
		if err != nil {
			return err
		}
		// codUnion.PutRawValue(decoded)
		t.Set(decoded)

	case 2:
		var decoded SpecialMap
		err = decoded.DecodeCod(buf)
		if err != nil {
			return err
		}
		// codUnion.PutRawValue(decoded)
		t.Set(decoded)

	default:
		panic("unknown type placed in union")
	}
	return err

	return err
}

func (t SpecialMap) EncodeCod(buf *cod.Buffer) {

	{
		value0 := map[string][]uint8(t)

		{
			buf.WriteUint64(uint64(len(value0)))

			for k1, v1 := range value0 {

				buf.WriteString(k1)
				{
					buf.WriteUint64(uint64(len(v1)))
					for i2 := range v1 {

						buf.WriteUint8(v1[i2])
					}
				}
			}

		}

	}
}

func (t *SpecialMap) DecodeCod(buf *cod.Buffer) error {
	var err error

	{
		var value0 map[string][]uint8

		{
			length, err := buf.ReadUint64()
			if err != nil {
				return err
			}

			if value0 == nil {
				value0 = make(map[string][]uint8)
			}

			for i1 := 0; i1 < int(length); i1++ {
				var key1 string
				var val1 []uint8

				key1, err = buf.ReadString()
				{
					length, err := buf.ReadUint64()
					if err != nil {
						return err
					}
					for i2 := 0; i2 < int(length); i2++ {
						var value2 uint8

						value2 = buf.ReadUint8()
						if err != nil {
							return err
						}

						val1 = append(val1, value2)
					}
				}
				if err != nil {
					return err
				}

				value0[key1] = val1
			}
		}
		*t = SpecialMap(value0)
	}

	return err
}

func (t Id) EncodeCod(buf *cod.Buffer) {

	buf.WriteUint16(t.Val)
}

func (t *Id) DecodeCod(buf *cod.Buffer) error {
	var err error

	t.Val, err = buf.ReadUint16()

	return err
}

func (t Person) EncodeCod(buf *cod.Buffer) {

	buf.WriteString(t.Name)
	buf.WriteUint8(t.Age)
	t.Id.EncodeCod(buf)
	for i0 := range t.Array {

		buf.WriteUint16(t.Array[i0])
	}
	{
		buf.WriteUint64(uint64(len(t.Slice)))
		for i0 := range t.Slice {

			buf.WriteUint32(t.Slice[i0])
		}
	}
	{
		buf.WriteUint64(uint64(len(t.DoubleSlice)))
		for i0 := range t.DoubleSlice {

			{
				buf.WriteUint64(uint64(len(t.DoubleSlice[i0])))
				for i1 := range t.DoubleSlice[i0] {

					buf.WriteUint8(t.DoubleSlice[i0][i1])
				}
			}
		}
	}
	{
		buf.WriteUint64(uint64(len(t.Map)))

		for k0, v0 := range t.Map {

			buf.WriteString(k0)
			{
				buf.WriteUint64(uint64(len(v0)))
				for i1 := range v0 {

					buf.WriteUint64(v0[i1])
				}
			}
		}

	}
	{
		buf.WriteUint64(uint64(len(t.MultiMap)))

		for k0, v0 := range t.MultiMap {

			buf.WriteString(k0)
			{
				buf.WriteUint64(uint64(len(v0)))

				for k1, v1 := range v0 {

					buf.WriteUint32(k1)
					{
						buf.WriteUint64(uint64(len(v1)))
						for i2 := range v1 {

							buf.WriteUint8(v1[i2])
						}
					}
				}

			}
		}

	}
	t.MyUnion.EncodeCod(buf)
}

func (t *Person) DecodeCod(buf *cod.Buffer) error {
	var err error

	t.Name, err = buf.ReadString()
	t.Age = buf.ReadUint8()
	t.Id.DecodeCod(buf)
	for i0 := range t.Array {

		t.Array[i0], err = buf.ReadUint16()

		if err != nil {
			return err
		}
	}
	{
		length, err := buf.ReadUint64()
		if err != nil {
			return err
		}
		for i0 := 0; i0 < int(length); i0++ {
			var value0 uint32

			value0, err = buf.ReadUint32()
			if err != nil {
				return err
			}

			t.Slice = append(t.Slice, value0)
		}
	}
	{
		length, err := buf.ReadUint64()
		if err != nil {
			return err
		}
		for i0 := 0; i0 < int(length); i0++ {
			var value0 []uint8

			{
				length, err := buf.ReadUint64()
				if err != nil {
					return err
				}
				for i1 := 0; i1 < int(length); i1++ {
					var value1 uint8

					value1 = buf.ReadUint8()
					if err != nil {
						return err
					}

					value0 = append(value0, value1)
				}
			}
			if err != nil {
				return err
			}

			t.DoubleSlice = append(t.DoubleSlice, value0)
		}
	}
	{
		length, err := buf.ReadUint64()
		if err != nil {
			return err
		}

		if t.Map == nil {
			t.Map = make(map[string][]uint64)
		}

		for i0 := 0; i0 < int(length); i0++ {
			var key0 string
			var val0 []uint64

			key0, err = buf.ReadString()
			{
				length, err := buf.ReadUint64()
				if err != nil {
					return err
				}
				for i1 := 0; i1 < int(length); i1++ {
					var value1 uint64

					value1, err = buf.ReadUint64()
					if err != nil {
						return err
					}

					val0 = append(val0, value1)
				}
			}
			if err != nil {
				return err
			}

			t.Map[key0] = val0
		}
	}
	{
		length, err := buf.ReadUint64()
		if err != nil {
			return err
		}

		if t.MultiMap == nil {
			t.MultiMap = make(map[string]map[uint32][]uint8)
		}

		for i0 := 0; i0 < int(length); i0++ {
			var key0 string
			var val0 map[uint32][]uint8

			key0, err = buf.ReadString()
			{
				length, err := buf.ReadUint64()
				if err != nil {
					return err
				}

				if val0 == nil {
					val0 = make(map[uint32][]uint8)
				}

				for i1 := 0; i1 < int(length); i1++ {
					var key1 uint32
					var val1 []uint8

					key1, err = buf.ReadUint32()
					{
						length, err := buf.ReadUint64()
						if err != nil {
							return err
						}
						for i2 := 0; i2 < int(length); i2++ {
							var value2 uint8

							value2 = buf.ReadUint8()
							if err != nil {
								return err
							}

							val1 = append(val1, value2)
						}
					}
					if err != nil {
						return err
					}

					val0[key1] = val1
				}
			}
			if err != nil {
				return err
			}

			t.MultiMap[key0] = val0
		}
	}
	t.MyUnion.DecodeCod(buf)

	return err
}

func (t MyUnion) Get() any {
	codUnion := cod.Union(t)
	rawVal := codUnion.GetRawValue()
	return rawVal

	// switch rawVal.(type) {
	// <no value>
	// default:
	//    panic("unknown type placed in union")
	// }
}

func (t *MyUnion) Set(v any) {
	codUnion := cod.Union(*t)
	codUnion.PutRawValue(v)
	*t = MyUnion(codUnion)

	// switch tagVal {
	// case 0: // Zero tag indicates nil
	//    return nil

	// <no value>
	// default:
	//    panic("unknown type placed in union")
	// }
	// return err
}

func NewMyUnion(v any) MyUnion {
	var ret MyUnion
	ret.Set(v)
	return ret
}
