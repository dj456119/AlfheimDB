/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-30 22:08:26
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-02 21:48:24
 */
package store

import (
	"fmt"
	"path/filepath"
	"strconv"

	badger "github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
)

type BadgerDBStore struct {
	DB *badger.DB
}

func NewBadgerDBStore(basedir string) *BadgerDBStore {
	bDBStore := new(BadgerDBStore)
	DB, err := badger.Open(badger.DefaultOptions(filepath.Join(basedir, "/badger")))
	if err != nil {
		logrus.Fatal("Open badger db error, ", err)
	}
	bDBStore.DB = DB
	return bDBStore
}

func (bDB *BadgerDBStore) Set(key string, value string) error {
	err := bDB.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})
	if err != nil {
		logrus.Fatal("badgerDB set error, ", err)
	}
	return nil
}
func (bDB *BadgerDBStore) Get(key string) (string, error) {
	result := ""
	err := bDB.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			result = string(val)
			return nil
		})

	})
	return result, err
}

func (bDB *BadgerDBStore) Incr(key string) (string, error) {
	var result string
	err := bDB.DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			logrus.Fatal("badgerDB incr get error, ", err)
		}
		var resultInt64 int64
		item.Value(func(val []byte) error {
			resultInt64, err = strconv.ParseInt(string(val), 10, 64)
			if err != nil {
				logrus.Fatal("badgerDB incr get value error, ", err)
			}
			return nil
		})
		resultInt64 = resultInt64 + 1

		result = fmt.Sprintf("%d", resultInt64)
		err = txn.Set([]byte(key), []byte(result))
		if err != nil {
			logrus.Fatal("badgerDB incr set value error, ", err)
		}
		return nil
	})

	return result, err
}
func (bDB *BadgerDBStore) Del(key string) error {
	err := bDB.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
	if err != nil {
		logrus.Fatal("badgerDB delete error, ", err)
	}
	return nil
}

func (bDB *BadgerDBStore) Keys(prefix string) ([]string, error) {

}

func (bDB *BadgerDBStore) Snapshot() ([]byte, error) {
	return nil, nil
}
func (bDB *BadgerDBStore) LoadSnapshot(data []byte) error {
	return nil
}
