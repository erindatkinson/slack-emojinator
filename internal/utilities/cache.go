package utilities

import (
	"github.com/dgraph-io/badger/v4"
)

var DB *badger.DB

func SetupDB(fPath string) error {
	var err error
	DB, err = badger.Open(
		badger.DefaultOptions(fPath).WithLoggingLevel(badger.ERROR))
	if err != nil {
		return err
	}
	return nil
}

func CloseDB() {
	DB.Close()
}

func MarkDownloaded(key []byte) error {
	return DB.Update(func(txn *badger.Txn) error {
		return txn.Set(
			[]byte(key),
			[]byte("true"))
	})
}

func CheckDownloaded(key []byte) error {
	return DB.View(func(txn *badger.Txn) error {
		_, err := txn.Get(key)
		if err != nil {
			return err
		}

		return nil
	})
}
