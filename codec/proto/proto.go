package proto

import (
	"encoding"
	"fmt"

	"github.com/golang/protobuf/proto"

	"github.com/omeid/thunder/codec"
)

func Codec() codec.Codec {
	return protoCodec{}
}

type marshaler struct {
	value interface{}
}

func (m marshaler) MarshalBinary() ([]byte, error) {
	bm, ok := m.value.(proto.Message)
	if !ok {
		return nil, fmt.Errorf("Proto Encoder Marshaler expects a proto.Message")
	}
	return proto.Marshal(bm)
}

type unmarshaler struct {
	value interface{}
}

func (m unmarshaler) UnmarshalBinary(data []byte) error {
	bm, ok := m.value.(proto.Message)
	if !ok {
		return fmt.Errorf("Proto Encoder Marshaler expects a proto.Message")
	}
	return proto.Unmarshal(data, bm)
}

type protoCodec struct{}

func (c protoCodec) Marshaler(v interface{}) encoding.BinaryMarshaler {
	return marshaler{v}
}

func (c protoCodec) Unmarshaler(v interface{}) encoding.BinaryUnmarshaler {
	return unmarshaler{v}
}
