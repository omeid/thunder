package thunder

import (
	"github.com/boltdb/bolt"
	"github.com/omeid/thunder/codec"
)

type DB struct {
	db *bolt.DB

	kc codec.Codec
	vc codec.Codec
}

func (b *DB) View(fn func(*Tx) error) error {
	return b.db.View(func(tx *bolt.Tx) error {
		return fn(&Tx{tx, b.kc, b.vc})
	})
}

func (b *DB) Update(fn func(*Tx) error) error {
	return b.db.Update(func(tx *bolt.Tx) error {
		return fn(&Tx{tx, b.kc, b.vc})
	})
}

func (db *DB) Sync() error {
	return db.db.Sync()
}
