package codec

import "encoding"

type Codec interface {
	Marshaler(v interface{}) encoding.BinaryMarshaler
	Unmarshaler(v interface{}) encoding.BinaryUnmarshaler
}
