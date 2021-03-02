package blabdb

import (
	"encoding/json"

	"math/rand"
	"strconv"
	"strings"

	"github.com/dgraph-io/badger/v3"
)

// Collection interface for basic crud actions inside collection
type Collection interface {
	Insert(item interface{}) (string, error)
	Update(key string, item interface{}) (string, error)
	GetByKey(key string, writeTo interface{}) error
	GetAllKeys() ([]string, error)
	DeleteByKey(key string) error
}

type blabDbCollection struct {
	collectioName string
	db            *blabDb
}

func useCollection(collectionName string, db *blabDb) Collection {
	return &blabDbCollection{
		collectioName: collectionName,
		db:            db,
	}
}

func concatCollectionAndKey(collectionName string, key string) string {
	trimmedCollection := strings.ReplaceAll(collectionName, "/", "")
	trimmedKey := strings.ReplaceAll(key, "/", "")
	return trimmedCollection + "/" + trimmedKey
}

func splitCollectionAndKey(entry string) (string, string) {
	sep := strings.Split(entry, "/")

	if len(sep) > 2 {
		panic("Data corrupted")
	}

	return sep[0], sep[1]
}

func createKey(db *blabDb, collection string) string {
	largest := 0

	err := db.badgerdb.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(collection)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			_, kk := splitCollectionAndKey(string(k))
			ik, err := strconv.Atoi(kk)
			if err != nil {
				return err
			}
			if ik > largest {
				largest = ik
			}
		}
		return nil
	})

	if err != nil {

		panic("Data corrupted")
	}

	min := 1234
	max := 3456
	randomNum := rand.Intn(max-min) + min

	return strconv.Itoa(largest + randomNum)

}

func (col *blabDbCollection) Insert(item interface{}) (string, error) {
	id := concatCollectionAndKey(col.collectioName, createKey(col.db, col.collectioName))
	payload, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	err = col.db.badgerdb.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte(id), []byte(payload))
		return err
	})

	if err != nil {
		return "", err
	}

	_, key := splitCollectionAndKey(id)

	return key, nil
}
func (col *blabDbCollection) Update(key string, item interface{}) (string, error) {
	id := concatCollectionAndKey(col.collectioName, key)

	// Check if found
	err := col.GetByKey(key, make(map[string]string))

	if err != nil {
		return "", err
	}

	payload, err := json.Marshal(item)
	if err != nil {
		return "", err
	}
	err = col.db.badgerdb.Update(func(txn *badger.Txn) error {
		err = txn.Set([]byte(id), []byte(payload))
		return err
	})

	return key, nil
}
func (col *blabDbCollection) GetByKey(key string, writeTo interface{}) error {
	id := concatCollectionAndKey(col.collectioName, key)
	var value []byte
	err := col.db.badgerdb.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(id))
		if err != nil {
			// Not found
			return err
		}
		err = item.Value(func(val []byte) error {
			clone := append([]byte{}, val...)
			value = clone
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	},
	)

	if err != nil {
		return err
	}

	json.Unmarshal(value, writeTo)

	return nil
}
func (col *blabDbCollection) DeleteByKey(key string) error {
	id := concatCollectionAndKey(col.collectioName, key)

	err := col.db.badgerdb.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(id))
		return err
	})

	if err != nil {
		return err
	}

	return nil
}

func (col *blabDbCollection) GetAllKeys() ([]string, error) {
	keys := []string{}
	err := col.db.badgerdb.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions

		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()

		prefix := []byte(col.collectioName)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()

			// Check if still exists
			// Badgerdb returns deleted keys for some time
			_, err := txn.Get(k)
			if err == nil {
				_, kk := splitCollectionAndKey(string(k))
				keys = append(keys, string(kk))
			}

		}
		return nil
	})

	if err != nil {

		return nil, err
	}

	return keys, nil
}
