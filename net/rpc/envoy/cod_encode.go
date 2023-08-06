package envoy

import (
	"github.com/unitoftime/cod/backend"
)

func (t ServerMsg) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ServerMsg) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}

func (t ServerReq) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ServerReq) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}

func (t ServerResp) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ServerResp) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}

func (t ClientReq) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ClientReq) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}

func (t ClientResp) EncodeCod(bs []byte) []byte {

	bs = backend.WriteVarUint16(bs, (t.Val))

	return bs
}

func (t *ClientResp) DecodeCod(bs []byte) (int, error) {
	var err error
	var n int
	var nOff int

	{
		var decoded uint16
		decoded, nOff, err = backend.ReadVarUint16(bs[n:])
		if err != nil {
			return 0, err
		}
		n += nOff
		t.Val = (decoded)
	}

	return n, err
}
