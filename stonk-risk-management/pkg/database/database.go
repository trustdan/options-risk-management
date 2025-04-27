package database

import (
	"encoding/json"
	"os"
	"time"

	"github.com/dgraph-io/badger/v3"
)

// DB encapsulates the badger database
type DB struct {
	db       *badger.DB
	gcTicker *time.Ticker
	stopGC   chan struct{}
}

// New creates a new database instance
func New(dbPath string) (*DB, error) {
	// Create directory if it doesn't exist
	if err := os.MkdirAll(dbPath, 0755); err != nil {
		return nil, err
	}

	options := badger.DefaultOptions(dbPath)
	options.Logger = nil // Disable Badger's default logger

	// Optimize for storage efficiency
	options.ValueLogFileSize = 10 * 1024 * 1024 // 10MB instead of default 1GB
	options.NumVersionsToKeep = 1               // Only keep the latest version of each key
	options.CompactL0OnClose = true             // Compact level 0 files on close

	db, err := badger.Open(options)
	if err != nil {
		return nil, err
	}

	// Create DB instance
	dbInstance := &DB{
		db:     db,
		stopGC: make(chan struct{}),
	}

	// Start background garbage collection
	dbInstance.startGC()

	return dbInstance, nil
}

// startGC starts a background goroutine for periodic garbage collection
func (d *DB) startGC() {
	d.gcTicker = time.NewTicker(30 * time.Minute) // Run GC every 30 minutes
	go func() {
		for {
			select {
			case <-d.gcTicker.C:
				// Run garbage collection with a more aggressive threshold (0.5)
				err := d.RunGC()
				if err != nil && err != badger.ErrNoRewrite {
					// Log error but continue (badger.ErrNoRewrite is normal when no garbage to collect)
					// We don't have a proper logger here, so we'll just print to stdout
					// In a production app, you might want to use a proper logger
					println("Badger GC error:", err.Error())
				}
			case <-d.stopGC:
				return
			}
		}
	}()
}

// Close closes the database
func (d *DB) Close() error {
	// Stop the GC goroutine
	if d.gcTicker != nil {
		d.gcTicker.Stop()
		close(d.stopGC)
	}
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

// GetKeysWithPrefix retrieves all keys matching a given prefix
func (d *DB) GetKeysWithPrefix(prefix string) ([]string, error) {
	var keys []string

	err := d.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false // We only need keys
		it := txn.NewIterator(opts)
		defer it.Close()

		prefixBytes := []byte(prefix)
		for it.Seek(prefixBytes); it.ValidForPrefix(prefixBytes); it.Next() {
			item := it.Item()
			keys = append(keys, string(item.Key()))
		}
		return nil
	})

	return keys, err
}

// RunGC runs the garbage collector to free up space
func (d *DB) RunGC() error {
	return d.db.RunValueLogGC(0.5)
}
