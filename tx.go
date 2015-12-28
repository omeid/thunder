package thunder

import (
	"github.com/boltdb/bolt"
	"github.com/omeid/thunder/codec"
)

type Tx struct {
	tx *bolt.Tx

	kc codec.Codec
	vc codec.Codec
}

func (tx *Tx) Bucket(name interface{}) (*Bucket, error) {

	n, err := tx.kc.Marshaler(name).MarshalBinary()
	if err != nil {
		return nil, err
	}

	bucket := tx.tx.Bucket(n)
	if bucket == nil {
		return nil, bolt.ErrBucketNotFound
	}

	return &Bucket{bucket, tx.kc, tx.vc}, nil
}

func (tx *Tx) Cursor() *Cursor {
	return &Cursor{tx.tx.Cursor(), tx.kc, tx.vc}
}
