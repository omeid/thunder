package thunder

import (
	"github.com/boltdb/bolt"
	"github.com/omeid/thunder/codec"
)

type Bucket struct {
	bucket *bolt.Bucket

	kc codec.Codec
	vc codec.Codec
}

func (b *Bucket) Bucket(name interface{}) (*Bucket, error) {

	n, err := b.kc.Marshaler(name).MarshalBinary()
	if err != nil {
		return nil, err
	}

	bucket := b.bucket.Bucket(n)
	if bucket == nil {
		return nil, bolt.ErrBucketNotFound
	}

	return &Bucket{bucket, b.kc, b.vc}, nil
}

func (b *Bucket) Cursor() *Cursor {
	return &Cursor{b.bucket.Cursor(), b.kc, b.vc}
}
