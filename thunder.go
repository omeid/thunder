package thunder

import (
	"errors"
	"os"
	"time"

	"github.com/boltdb/bolt"
	"github.com/omeid/thunder/codec"
)

var (
	ErrKeyValueNotFound = errors.New("Key Value not found.")
)

type Options struct {
	Timeout    time.Duration
	NoGrowSync bool
	ReadOnly   bool
	MmapFlags  int
	KeyCodec   codec.Codec
	ValueCodec codec.Codec
}

func Open(path string, mode os.FileMode, options *Options) (*DB, error) {
	db, err := bolt.Open(path, mode, &bolt.Options{
		Timeout:    options.Timeout,
		NoGrowSync: options.NoGrowSync,
		ReadOnly:   options.ReadOnly,
		MmapFlags:  options.MmapFlags,
	})
	return &DB{db, options.KeyCodec, options.ValueCodec}, err
}
