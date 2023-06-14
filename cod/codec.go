package cod

//--------------------------------------------------------------------------------
type Encoder[T any] func(*Buffer, T)
type Decoder[T any] func(*Buffer) (T, error)

type Codec[T any] struct {
	enc Encoder[T]
	dec Decoder[T]
}

func NewCodec[T any](enc Encoder[T], dec Decoder[T]) *Codec[T] {
	return &Codec[T]{
		enc: enc,
		dec: dec,
	}
}

func (c *Codec[T]) Write(buf *Buffer, v T) *Buffer {
	c.enc(buf, v)
	return buf
}

func (c *Codec[T]) Read(buf *Buffer) (T, error) {
	return c.dec(buf)
}


func (c *Codec[T]) WriteSlice(buf *Buffer, v []T) *Buffer {
	buf.WriteInt64(int64(len(v)))
	for i := range v {
		c.Write(buf, v[i])
	}
	return buf
}

func (c *Codec[T]) ReadSlice(buf *Buffer) ([]T, error) {
	l, err := buf.ReadInt64()
	if err != nil { return nil, err }
	if l < 0 { return nil, ErrVarintCorrupted }
	length := int(l)

	ret := make([]T, length)
	for i := 0; i < length; i++ {
		val, err := c.Read(buf)
		if err != nil {
			return nil, err
		}
		ret[i] = val
	}
	return ret, nil
}

type MapCodec[K comparable,V any] struct {
	keyCodec *Codec[K]
	valCodec *Codec[V]
}

func NewMapCodec[K comparable,V any](keyCodec *Codec[K], valCodec *Codec[V]) *MapCodec[K,V] {
	return &MapCodec[K,V]{
		keyCodec: keyCodec,
		valCodec: valCodec,
	}
}

func (c *MapCodec[K,V]) Write(buf *Buffer, m map[K]V) *Buffer {
	buf.WriteInt64(int64(len(m)))
	for k,v := range m {
		c.keyCodec.Write(buf, k)
		c.valCodec.Write(buf, v)
	}
	return buf
}

func (c *MapCodec[K,V]) Read(buf *Buffer) (map[K]V, error) {
	l, err := buf.ReadInt64()
	if err != nil { return nil, err }
	if l < 0 { return nil, ErrVarintCorrupted }
	length := int(l)

	ret := make(map[K]V, length) // TODO: Not sure how useful initial capacity is
	for i := 0; i < length; i++ {
		key, err := c.keyCodec.Read(buf)
		if err != nil { return nil, nil }
		val, err := c.valCodec.Read(buf)
		if err != nil { return nil, nil }

		ret[key] = val
	}
	return ret, nil
}
