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

	v := reflect.Indirect(reflect.ValueOf(m.value))
	if v.Kind() != reflect.String {
		err = fmt.Errorf("Strings Codec: Expected type string but got %s", reflect.TypeOf(m.value))
	}

	vs := v.Interface().(string)
	if vs == "" {
		err = fmt.Errorf("Empty Key.")
	}
	return []byte(vs), err
}

type unmarshaler struct {
	value interface{}
}

func (m unmarshaler) UnmarshalBinary(data []byte) error {

	rv := reflect.ValueOf(m.value)
	if rv.Kind() != reflect.Ptr || rv.Elem().Kind() != reflect.String {
		return fmt.Errorf("Strings Codec: Expects Non-nil String Pointer.")
	}
	rv.Elem().SetString(string(data))
	return nil
}

type stringCodec struct{}

func (c stringCodec) Marshaler(v interface{}) encoding.BinaryMarshaler {
	return marshaler{v}
}

func (c stringCodec) Unmarshaler(v interface{}) encoding.BinaryUnmarshaler {
	return unmarshaler{v}
}
