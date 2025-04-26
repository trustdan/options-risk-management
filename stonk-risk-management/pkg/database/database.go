package database

import (
	"encoding/json"
	"os"

	"github.com/dgraph-io/badger/v3"
)

// DB encapsulates the badger database
type DB struct {
	db *badger.DB
}

// New creates a new database instance
func New(dbPath string) (*DB, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, err
	}

	options := badger.DefaultOptions(dbPath)
	options.Logger = nil // Disable Badger's default logger

	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	return &DB{db: db}, nil
}

// Close closes the database
func (d *DB) Close() error {
	return d.db.Close()
}

// Put stores a value in the database
func (d *DB) Put(key string, value interface{}) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	})
}

// Get retrieves a value from the database
func (d *DB) Get(key string, value interface{}) error {
	var data []byte
	err := d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}

		data, err = item.ValueCopy(nil)
		return err
	})

	if err != nil {
		return err
	}

	return json.Unmarshal(data, value)
}

// Delete removes a key from the database
func (d *DB) Delete(key string) error {
	return d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	})
}

// GetAllWithPrefix retrieves all keys with a given prefix
func (d *DB) GetAllWithPrefix(prefix string) ([][]byte, error) {
	var values [][]byte

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		it := txn.NewIterator(opts)
		defer it.Close()

		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}
			values = append(values, val)
		}
		return nil
	})

	return values, err
}

// RunGC runs the garbage collector to free up space
func (d *DB) RunGC() error {
	return d.db.RunValueLogGC(0.5)
}
