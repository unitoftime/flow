package bench

import (
	"github.com/unitoftime/flow/cod"
)

func (t MUSA) EncodeCod(buf *cod.Buffer) {

	buf.WriteString(t.Name)
	buf.WriteInt64(t.BirthDay)
	buf.WriteString(t.Phone)
	buf.WriteInt32(t.Siblings)
	buf.WriteBool(t.Spouse)
	buf.WriteFloat64(t.Money)
}

func (t *MUSA) DecodeCod(buf *cod.Buffer) error {
	var err error

	t.Name, err = buf.ReadString()
	t.BirthDay, err = buf.ReadInt64()
	t.Phone, err = buf.ReadString()
	t.Siblings, err = buf.ReadInt32()
	t.Spouse = buf.ReadBool()
	t.Money = buf.ReadFloat64()

	return err
}
