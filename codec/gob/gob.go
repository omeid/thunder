package gob

import (
	"bytes"
	"encoding"
	"encoding/gob"

	"github.com/omeid/thunder/codec"
)

func Codec() codec.Codec {
	return gobCodec{}
}

type marshaler struct {
	value interface{}
}

func (m marshaler) MarshalBinary() ([]byte, error) {
	var buff *bytes.Buffer
	err := gob.NewEncoder(buff).Encode(m.value)
	return buff.Bytes(), err
}

type unmarshaler struct {
	value interface{}
}

func (m unmarshaler) UnmarshalBinary(data []byte) error {
	return gob.NewDecoder(bytes.NewReader(data)).Decode(m.value)
}

type gobCodec struct{}

func (c gobCodec) Marshaler(v interface{}) encoding.BinaryMarshaler {
	return marshaler{v}
}

func (c gobCodec) Unmarshaler(v interface{}) encoding.BinaryUnmarshaler {
	return unmarshaler{v}
}
