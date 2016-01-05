# Thunder 
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](https://godoc.org/github.com/omeid/thunder) ![Project status](https://img.shields.io/badge/status-experimental-red.svg?style=flat-square)

Thunder is a high productivity wrapper for [BoltDB](https://github.com/boltdb/bolt).
Thunder allows you to do develop rapidly and take care of performance when it becomes and issue by simply falling back to using BoltDB directly.

Thunder is unopinionated about serialization format and allows you to bring your own [Codec](https://godoc.org/omeid/thunder/codec") by implementing a simple interface:

```go
type Codec interface {
  Marshaler(v interface{}) encoding.BinaryMarshaler
  Unmarshaler(v interface{}) encoding.BinaryUnmarshaler
}
```

Thunder comes with `json`, `string`, `gob`, and `Protocol Buffers` codecs for free.


