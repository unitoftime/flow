package phy2

import (
	"github.com/unitoftime/cod/backend"
)

// func (t Pos) EncodeCod(bs []byte) []byte {

// 	{
// 		value0 := Vec2(t)

// 		bs = value0.EncodeCod(bs)

// 	}
// 	return bs
// }

// func (t *Pos) DecodeCod(bs []byte) (int, error) {
// 	var err error
// 	var n int
// 	var nOff int

// 	{
// 		var value0 Vec2

// 		{
// 			var decoded Vec2
// 			nOff, err = decoded.DecodeCod(bs[n:])
// 			if err != nil {
// 				return 0, err
// 			}
// 			n += nOff
// 			value0 = decoded
// 		}

// 		*t = Pos(value0)
// 	}

// 	return n, err
// }

func (t Vel) EncodeCod(bs []byte) []byte {

	{
		value0 := Vec2(t)

		bs = value0.EncodeCod(bs)

	}
	return bs
}

func (t *Vel) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var value0 Vec2

		{
			var decoded Vec2
			nOff, err = decoded.DecodeCod(bs[n:])
			if err != nil {
				return 0, err
			}
			n += nOff
			value0 = decoded
		}

		*t = Vel(value0)
	}

	return n, err
}

func (t Scale) EncodeCod(bs []byte) []byte {

	{
		value0 := Vec2(t)

		bs = value0.EncodeCod(bs)

	}
	return bs
}

func (t *Scale) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var value0 Vec2

		{
			var decoded Vec2
			nOff, err = decoded.DecodeCod(bs[n:])
			if err != nil {
				return 0, err
			}
			n += nOff
			value0 = decoded
		}

		*t = Scale(value0)
	}

	return n, err
}

func (t Rotation) EncodeCod(bs []byte) []byte {

	{
		value0 := float64(t)

		bs = backend.WriteFloat64(bs, (value0))

	}
	return bs
}

func (t *Rotation) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var value0 float64

		{
			var decoded float64
			decoded, nOff, err = backend.ReadFloat64(bs[n:])
			if err != nil {
				return 0, err
			}
			n += nOff
			value0 = (decoded)
		}

		*t = Rotation(value0)
	}

	return n, err
}

func (t Rigidbody) EncodeCod(bs []byte) []byte {

	bs = backend.WriteFloat64(bs, (t.Mass))

	bs = t.Velocity.EncodeCod(bs)
	return bs
}

func (t *Rigidbody) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded float64
		decoded, nOff, err = backend.ReadFloat64(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Mass = (decoded)
	}

	{
		var decoded Vec2
		nOff, err = decoded.DecodeCod(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Velocity = decoded
	}

	return n, err
}

func (t Vec2) EncodeCod(bs []byte) []byte {

	bs = backend.WriteFloat64(bs, (t.X))

	bs = backend.WriteFloat64(bs, (t.Y))

	return bs
}

func (t *Vec2) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded float64
		decoded, nOff, err = backend.ReadFloat64(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.X = (decoded)
	}

	{
		var decoded float64
		decoded, nOff, err = backend.ReadFloat64(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Y = (decoded)
	}

	return n, err
}
