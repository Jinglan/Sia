// Package blockdb provides read/write access to an on-disk block database.
// blockdb uses Bolt as its underlying database, but this is subject to change.
package blockdb

import (
	"errors"

	"github.com/boltdb/bolt"

	"github.com/NebulousLabs/Sia/encoding"
	"github.com/NebulousLabs/Sia/types"
)

var (
	ErrUnknownBlock = errors.New("unknown block")
)

// A DB is a block database. Right now it is equivalent to a []Block.
type DB interface {
	Block(types.BlockHeight) (types.Block, error)
	AddBlock(types.Block) error
	RemoveBlock() error
	Height() (types.BlockHeight, error)
	Close() error
}

// boltDB implements the DB interface. It is a Bolt database of Blocks, ordered
// by BlockHeight.
type boltDB struct {
	*bolt.DB
}

func (db *boltDB) Block(height types.BlockHeight) (types.Block, error) {
	key := encoding.EncUint64(uint64(height))
	var block types.Block
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("chain"))
		value := b.Get(key)
		if value == nil {
			return ErrUnknownBlock
		}
		// TODO: move outside the tx?
		// NOTE: value is not valid outside the tx
		return encoding.Unmarshal(value, &block)
	})
	return block, err
}

func (db *boltDB) AddBlock(block types.Block) error {
	value := encoding.Marshal(block)
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("chain"))
		key := encoding.EncUint64(uint64(b.Stats().KeyN))
		return b.Put(key, value)
	})
}

func (db *boltDB) RemoveBlock() error {
	return db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("chain"))
		key := encoding.EncUint64(uint64(b.Stats().KeyN - 1))
		return b.Delete(key)
	})
}

func (db *boltDB) Height() (types.BlockHeight, error) {
	var height types.BlockHeight
	err := db.View(func(tx *bolt.Tx) error {
		height = types.BlockHeight(tx.Bucket([]byte("chain")).Stats().KeyN)
		return nil
	})
	return height, err
}

// Open returns a database ready for use. If the database file does not exist,
// it will be created. Only one view of a given database file should be open at
// any one time.
func Open(filename string) (DB, error) {
	db, err := bolt.Open(filename, 0600, nil)
	if err != nil {
		return nil, err
	}
	// create buckets
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("chain"))
		if err != nil {
			return err
		}
		//_, err = tx.CreateBucketIfNotExists([]byte("utxos"))
		return err
	})
	return &boltDB{db}, nil
}
