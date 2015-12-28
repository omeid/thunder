package json

import (
	"encoding"
	"encoding/json"

	"github.com/omeid/thunder/codec"
)

func Codec() codec.Codec {
	return jsonCodec{}
}

type marshaler struct {
	value interface{}
}

func (m marshaler) MarshalBinary() ([]byte, error) {
	return json.Marshal(m.value)
}

type unmarshaler struct {
	data  []byte
	value interface{}
}

func (m unmarshaler) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(m.data, m.value)
}

type jsonCodec struct{}

func (c jsonCodec) Marshaler(v interface{}) encoding.BinaryMarshaler {
	return marshaler{v}
}

func (c jsonCodec) Unmarshaler(v interface{}) encoding.BinaryUnmarshaler {
	return unmarshaler{}
}
