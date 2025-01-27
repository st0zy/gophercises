package database

import (
	"errors"
	"time"

	"github.com/boltdb/bolt"
	"github.com/st0zy/gophercises/urlshort/mapping"
)

var DBConnectionError = errors.New("failed to open connection to the database")

func OpenDB(path string) (*bolt.DB, error) {
	db, err := bolt.Open(path, 0600, &bolt.Options{
		Timeout: time.Second * 1,
	})
	if err != nil {
		return nil, DBConnectionError
	}
	LoadInitialData(db)
	return db, nil
}

func LoadInitialData(db *bolt.DB) error {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("url_mapping"))
		if err != nil {
			return err
		}
		err = bucket.Put([]byte("/db-tesst"), []byte("https://linkedin.com"))
		return err

	})

	return err
}

func GetAllMappings(db *bolt.DB) []mapping.PathMapping {

	var allMappings []mapping.PathMapping
	db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("url_mapping"))
		cursor := bucket.Cursor()
		for k, v := cursor.First(); k != nil; k, v = cursor.Next() {
			allMappings = append(allMappings, mapping.PathMapping{
				Path:         string(k),
				RedirectPath: string(v),
			})
		}

		return nil
	})

	return allMappings
}
