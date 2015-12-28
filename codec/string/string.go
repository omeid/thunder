package strings

import (
	"encoding"
	"fmt"
	"reflect"

	"github.com/omeid/thunder/codec"
)

func Codec() codec.Codec {
	return stringCodec{}
}

type marshaler struct {
	value interface{}
}

func (m marshaler) MarshalBinary() ([]byte, error) {

	var err error
	bytes, ok := m.value.([]byte)
	if !ok {
		err = fmt.Errorf("Strings Codec: Expected type string but got %s", reflect.TypeOf(m.value))
	}

	return bytes, err
}

type unmarshaler struct {
	data  []byte
	value interface{}
}

func (m unmarshaler) UnmarshalBinary(data []byte) error {

	rv := reflect.ValueOf(m.value)
	if rv.IsNil() || rv.Type() != reflect.TypeOf((*string)(nil)) {
		return fmt.Errorf("Strings Codec: Expects Non-nil String Pointer.")
	}

	m.value = string(m.data)
	return nil
}

type stringCodec struct{}

func (c stringCodec) Marshaler(v interface{}) encoding.BinaryMarshaler {
	return marshaler{v}
}

func (c stringCodec) Unmarshaler(v interface{}) encoding.BinaryUnmarshaler {
	return unmarshaler{}
}
