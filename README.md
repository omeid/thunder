# Thunder 
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/omeid/thunder) ![Project status](https://img.shields.io/badge/status-experimental-red.svg?style=flat-square)

Thunder is a high productivity wrapper for [BoltDB](https://github.com/boltdb/bolt).
The idea behind thunder is to allow users to rapidly develop but allow them to fallback to using BoltDB directly by maintaining an almost functionally identical API yet hide some of the complexity involved around using a serializaiton format, filtering, batch insert and reads.


Thunder stays unopinionated about serialization format used by allowing users to bring their own [Codec](https://godoc.org/omeid/thunder/codec") by implementing a simple interface:

```go
type Codec interface {
  Marshaler(v interface{}) encoding.BinaryMarshaler
  Unmarshaler(v interface{}) encoding.BinaryUnmarshaler
}
```

Thunder comes with `json`, `string`, `gob`, and `Protocol Buffers` codecs for free. You may mix and match any combination of Codecs for your keys and values.
