/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-30 22:08:26
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-30 22:40:44
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
	return bDB.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), []byte(value))
	})

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
			return err
		}
		var resultInt64 int64
		item.Value(func(val []byte) error {
			resultInt64, err = strconv.ParseInt(string(val), 10, 64)
			if err != nil {
				return err
			}
			return nil
		})
		resultInt64 = resultInt64 + 1

		result = fmt.Sprintf("%d", resultInt64)
		return txn.Set([]byte(key), []byte(result))
	})
	return result, err
}
func (bDB *BadgerDBStore) Del(key string) error {
	return bDB.DB.Update(func(txn *badger.Txn) error {
		err := txn.Delete([]byte(key))
		return err
	})
}
func (bDB *BadgerDBStore) Snapshot() ([]byte, error) {
	return nil, nil
}
func (bDB *BadgerDBStore) LoadSnapshot(data []byte) error {
	return nil
}
