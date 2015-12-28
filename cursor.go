package thunder

import (
	"github.com/boltdb/bolt"
	"github.com/omeid/thunder/codec"
)

func unmarshal(kc, vc codec.Codec, k, v []byte, key, value interface{}) error {
	if k == nil && v == nil {
		return ErrKeyValueNotFound
	}
	err := kc.Unmarshaler(key).UnmarshalBinary(k)
	if err != nil {
		return err
	}

	err = vc.Unmarshaler(value).UnmarshalBinary(v)
	return err
}

type Cursor struct {
	cursor *bolt.Cursor

	kc codec.Codec
	vc codec.Codec
}

func (cur *Cursor) Bucket() *Bucket {
	return &Bucket{cur.cursor.Bucket(), cur.kc, cur.vc}
}

func (cur *Cursor) First(key interface{}, value interface{}) error {
	k, v := cur.cursor.First()
	return unmarshal(cur.kc, cur.vc, k, v, key, value)
}

func (cur *Cursor) Last(key interface{}, value interface{}) error {
	k, v := cur.cursor.Last()
	return unmarshal(cur.kc, cur.vc, k, v, key, value)
}

func (cur *Cursor) Next(key interface{}, value interface{}) error {
	k, v := cur.cursor.Next()
	return unmarshal(cur.kc, cur.vc, k, v, key, value)
}

func (cur *Cursor) Prev(key interface{}, value interface{}) error {
	k, v := cur.cursor.Prev()
	return unmarshal(cur.kc, cur.vc, k, v, key, value)
}

func (cur *Cursor) Seek(seek interface{}, key interface{}, value interface{}) error {
	s, err := cur.kc.Marshaler(seek).MarshalBinary()
	if err != nil {
		return err
	}

	k, v := cur.cursor.Seek(s)
	return unmarshal(cur.kc, cur.vc, k, v, key, value)
}

func New(db *bolt.DB, KeyCodec, ValueCodec codec.Codec) *DB {
	return &DB{db: db, kc: KeyCodec, vc: ValueCodec}
}
