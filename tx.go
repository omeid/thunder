package thunder

import (
	"github.com/boltdb/bolt"
	"github.com/omeid/thunder/codec"
)

type Tx struct {
	tx *bolt.Tx

	kc codec.Codec
	vc codec.Codec

	err error
}

func (tx *Tx) Err() error {
	return tx.err
}

func (tx *Tx) Commit() error {
	return tx.tx.Commit()
}

func (tx *Tx) Rollback() error {
	return tx.tx.Rollback()
}

func (tx *Tx) CreateBucketIfNotExists(name interface{}) *Bucket {
	if tx.err != nil {
		return &Bucket{nil, tx.kc, tx.vc, tx.err}
	}

	n, err := tx.kc.Marshaler(name).MarshalBinary()
	if err != nil {
		return &Bucket{nil, tx.kc, tx.vc, err}
	}

	bucket, err := tx.tx.CreateBucketIfNotExists(n)
	return &Bucket{
		bucket: bucket,
		kc:     tx.kc,
		vc:     tx.vc,
		err:    err,
	}

}

func (tx *Tx) Bucket(name interface{}) *Bucket {

	if tx.err != nil {
		return &Bucket{nil, tx.kc, tx.vc, tx.err}
	}

	n, err := tx.kc.Marshaler(name).MarshalBinary()

	bucket := tx.tx.Bucket(n)
	if bucket == nil {
		err = bolt.ErrBucketNotFound
	}

	return &Bucket{
		bucket: bucket,
		kc:     tx.kc,
		vc:     tx.vc,
		err:    err,
	}
}

func (tx *Tx) Cursor() *Cursor {
	return &Cursor{
		cursor: tx.tx.Cursor(),
		kc:     tx.kc,
		vc:     tx.vc,
		err:    tx.err,
	}
}
