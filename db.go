package blabdb

import (
	badger "github.com/dgraph-io/badger/v3"
)

var globaldb *blabDb

// BlabDb tenrface for db
type BlabDb interface {
	Collection(name string) Collection
	GetAllCollections() ([]string, error)
	Close() error
}

type blabDb struct {
	badgerdb *badger.DB
}

// Open creates a new embedded dtabase or connects if one exists
func Open(path string) (BlabDb, error) {
	if globaldb != nil {
		return globaldb, nil
	}
	bdb, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}

	globaldb = &blabDb{
		badgerdb: bdb,
	}

	return globaldb, nil
}

func (db *blabDb) Collection(name string) Collection {
	return useCollection(name, db)
}

func (db *blabDb) GetAllCollections() ([]string, error) {
	collections := []string{}
	err := db.badgerdb.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			// Check if exists
			_, err := txn.Get([]byte(k))
			if err == nil {
				collection, _ := splitCollectionAndKey(string(k))
				collections = append(collections, collection)
			}

		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	return collections, nil
}

func (db *blabDb) Close() error {
	err := db.badgerdb.Close()
	globaldb = nil
	return err
}
