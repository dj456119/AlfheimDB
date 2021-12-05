/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-30 22:08:26
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-05 00:13:46
 */
package store

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

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
			switch err {
			case badger.ErrKeyNotFound:
				result = "nil"
				return nil
			default:
				logrus.Fatal("badgerDB get error, ", err)
			}
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
		var resultInt64 int64
		item, err := txn.Get([]byte(key))
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				resultInt64 = 0
			default:
				logrus.Fatal("badgerDB incr get error, ", err)
			}
		}

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
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				return nil
			default:
				logrus.Fatal("badgerDB delete error, ", err)
			}
		}
		return err
	})
	return err
}

func (bDB *BadgerDBStore) Keys(prefix string) ([]string, error) {
	result := []string{}
	if prefix == "*" {
		prefix = ""
	}
	bDB.DB.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		prefix := []byte(prefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()
			k := item.Key()
			result = append(result, string(k))
		}
		return nil
	})
	return result, nil
}

func (bDB *BadgerDBStore) SetEx(key string, value string, timeout int64) error {
	err := bDB.DB.Update(func(txn *badger.Txn) error {
		e := badger.NewEntry([]byte(key), []byte(value)).WithTTL(time.Duration(timeout * (int64(time.Millisecond))))
		err := txn.SetEntry(e)

		return err
	})
	switch err {
	case badger.ErrKeyNotFound:
		return err
	default:
		logrus.Fatal("badgerDB set ex error, ", err)
	}
	return err
}
func (bDB *BadgerDBStore) TTL(key string) (string, error) {
	var result uint64
	err := bDB.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		result = item.ExpiresAt()
		return nil
	})
	return fmt.Sprintf("%d", result), err
}

func (bDB *BadgerDBStore) SetNx(key string, value string) (int, error) {
	result := 0
	err := bDB.DB.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				return err
			default:
				logrus.Fatal("badgerDB incr get error, ", err)
			}
		}

		err = txn.Set([]byte(key), []byte(value))
		if err != nil {
			logrus.Fatal("badgerDB incr set value error, ", err)
		}
		result = 1
		return nil
	})
	return result, err
}

func (bDB *BadgerDBStore) Snapshot() ([]byte, error) {
	return nil, nil
}
func (bDB *BadgerDBStore) LoadSnapshot(data []byte) error {
	return nil
}
