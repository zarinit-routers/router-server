package storage

import (
	"fmt"
	"sync"

	"github.com/charmbracelet/log"
	"github.com/dgraph-io/badger/v4"
	"github.com/spf13/viper"
)

var logger log.Logger
var dbMutex sync.Mutex

func getDBPath() string {
	return viper.GetString("storage.key-value.path")
}

func init() {
	logger = *log.WithPrefix("Storage")
}
func mustConnect() *badger.DB {
	db, err := connect()

	if err != nil {
		logger.Fatal("Failed get connection to key-value storage", "error", err, "suggestion", "check if the key-value store is available at the start of a program with storage.Check()")
	}
	return db
}
func connect() (*badger.DB, error) {
	path := getDBPath()
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, fmt.Errorf("failed open connection to storage file %q: %s", path, err)
	}
	return db, nil
}

// Tries to connect to the key-value store and returns an error if it fails.
// Use this function to check if the key-value store is available at the start of a program.
func Check() error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer closeDB(db)
	return nil
}
func GetString(key string) string {
	return getString(key)
}
func SetString(key, value string) error {
	return setString(key, value)
}

func getString(key string) string {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	db := mustConnect()
	defer closeDB(db)

	txn := db.NewTransaction(false)
	item, err := txn.Get([]byte(key))
	if err != nil {
		logger.Warn("Failed get value from key-value store", "key", key, "error", err.Error())
		return ""
	}
	value, err := item.ValueCopy(nil)
	if err != nil {
		logger.Warn("Failed get value from key-value store", "key", key, "error", err.Error())
		return ""
	}
	if err := txn.Commit(); err != nil {
		logger.Warn("Failed commit transaction for getting value", "key", key, "error", err.Error())
	}
	return string(value)
}

func setString(key, value string) error {
	dbMutex.Lock()
	defer dbMutex.Unlock()
	db := mustConnect()
	defer closeDB(db)
	txn := db.NewTransaction(true)
	err := txn.Set([]byte(key), []byte(value))
	if err != nil {
		return err
	}
	return txn.Commit()
}
func closeDB(db *badger.DB) {
	if err := db.Close(); err != nil {
		log.Error("Error while closing key-value storage", "error", err)
	}

}
