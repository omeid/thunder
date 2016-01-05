package thunder

import (
	"fmt"
	"reflect"

	"github.com/boltdb/bolt"
	"github.com/omeid/thunder/codec"
)

type Bucket struct {
	bucket *bolt.Bucket

	kc codec.Codec
	vc codec.Codec

	err error
}

func (b *Bucket) Err() error {
	return b.err
}

func (b *Bucket) Bucket(name interface{}) *Bucket {

	if b.err != nil {
		return &Bucket{nil, b.kc, b.vc, b.err}
	}

	n, err := b.kc.Marshaler(name).MarshalBinary()
	if err != nil {
		return &Bucket{nil, b.kc, b.vc, err}
	}

	bucket := b.bucket.Bucket(n)
	if bucket == nil {
		err = bolt.ErrBucketNotFound
	}

	return &Bucket{bucket, b.kc, b.vc, err}
}

func (b *Bucket) Cursor() *Cursor {
	return &Cursor{b.bucket.Cursor(), b.kc, b.vc, b.err}
}

func (b *Bucket) Put(key interface{}, value interface{}) error {
	if b.err != nil {
		return b.err
	}
	k, err := b.kc.Marshaler(key).MarshalBinary()
	if err != nil {
		return err
	}
	v, err := b.vc.Marshaler(value).MarshalBinary()
	if err != nil {
		return err
	}
	return b.bucket.Put(k, v)
}

func (b *Bucket) Get(key interface{}, value interface{}) error {
	if b.err != nil {
		return b.err
	}

	k, err := b.kc.Marshaler(key).MarshalBinary()
	if err != nil {
		return err
	}
	v := b.bucket.Get(k)

	return b.vc.Unmarshaler(value).UnmarshalBinary(v)
}

func (b *Bucket) All(kvm interface{}) error {
	return b.Where(func(interface{}) bool { return true }, kvm)
}

func (b *Bucket) Where(match func(interface{}) bool, kvm interface{}) error {

	mv := reflect.Indirect(reflect.ValueOf(kvm))

	if mv.Kind() != reflect.Map {
		return fmt.Errorf("Thunder: Bucket.All Expects a Map")
	}

	if mv.IsNil() {
		return fmt.Errorf("Thunder: Bucket.All expects a map, nil given.")
	}

	mvt := mv.Type()

	kt := mvt.Key()
	vt := mvt.Elem()

	return b.bucket.ForEach(func(k, v []byte) error {

		key, value := reflect.New(kt), reflect.New(vt)
		err := unmarshal(b.kc, b.vc, k, v, key.Interface(), value.Interface())
		if err != nil {
			return err
		}
		if kt.Kind() != reflect.Ptr {
			key = reflect.Indirect(key)
		}
		if vt.Kind() != reflect.Ptr {
			value = reflect.Indirect(value)
		}

		if match(value.Interface()) {
			mv.SetMapIndex(key, value)
		}
		return nil
	})
}
