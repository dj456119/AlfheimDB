/*
 * @Descripttion: AlfheimdDB store engine by badgerdb
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-30 22:08:26
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-17 14:38:06
 */
package store

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/AlfheimDB/pb"
	badger "github.com/dgraph-io/badger/v3"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type BadgerDBStore struct {
	DB *badger.DB
}

type BadgerDBValue struct {
	Value     []byte
	TimeStamp int64
	Timeout   int64
	buff      []byte
}

func NewBadgerDBValue(value []byte, timeStamp, timeout int64) (*BadgerDBValue, error) {
	bdsv := BadgerDBValue{
		Value:     value,
		TimeStamp: timeStamp,
		Timeout:   timeout,
	}
	err := bdsv.Encode()
	if err != nil {
		return nil, err
	}
	return &bdsv, nil
}

func (bdsv *BadgerDBValue) Encode() error {
	pbBdsv := pb.BadgerStringValue{
		Value:     bdsv.Value,
		TimeStamp: int64(bdsv.TimeStamp),
		Timeout:   int64(bdsv.Timeout),
	}
	var err error
	bdsv.buff, err = proto.Marshal(&pbBdsv)
	return err
}

func (bdsv *BadgerDBValue) Decode(buff []byte) error {
	pbBdsv := pb.BadgerStringValue{}
	err := proto.Unmarshal(buff, &pbBdsv)
	if err != nil {
		return err
	}
	bdsv.TimeStamp = pbBdsv.TimeStamp
	bdsv.Timeout = pbBdsv.Timeout
	bdsv.Value = pbBdsv.Value
	bdsv.buff = buff
	return nil
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

func (bDB *BadgerDBStore) Set(key string, value string, nowTime int64) error {
	stringValue, err := NewBadgerDBValue([]byte(value), nowTime/1e6, -1)
	if err != nil {
		return err
	}
	err = bDB.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), stringValue.buff)
	})
	if err != nil {
		logrus.Fatal("badgerDB set error, ", err)
	}
	return nil
}
func (bDB *BadgerDBStore) Get(key string) (*string, error) {
	var result *BadgerDBValue
	err := bDB.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			bdsv := BadgerDBValue{}
			result = &bdsv
			return bdsv.Decode(val)
		})

	})
	if err != nil {
		switch err {
		case badger.ErrKeyNotFound:
			return nil, nil
		default:
			logrus.Fatal("badgerDB get error, ", err)
		}
	}
	if result.Timeout != -1 && time.Now().UnixNano()/1e6 > result.TimeStamp+result.Timeout {
		return nil, nil
	}
	resultString := string(result.Value)
	return &resultString, err
}

func (bDB *BadgerDBStore) Incr(key string, nowTime int64) (int64, error) {
	var resultInt64 int64
	var bdsv *BadgerDBValue
	err := bDB.DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				resultInt64 = 0
			default:
				logrus.Fatal("badgerDB incr get error, ", err)
			}
		} else {
			item.Value(func(val []byte) error {
				bdsv = &BadgerDBValue{}
				err := bdsv.Decode(val)
				if err != nil {
					logrus.Fatal("badgerDB get value error, ", err)
				}

				if bdsv.Timeout != -1 && nowTime/1e6 > bdsv.TimeStamp+bdsv.Timeout {
					return nil
				}

				resultInt64, err = strconv.ParseInt(string(bdsv.Value), 10, 64)
				if err != nil {
					logrus.Fatal("badgerDB incr get value error, ", err)
				}
				return nil
			})
		}

		resultInt64 = resultInt64 + 1
		result := fmt.Sprintf("%d", resultInt64)
		if bdsv == nil {
			bdsv, err = NewBadgerDBValue([]byte(result), nowTime/1e6, -1)
			if err != nil {
				logrus.Fatal("Incr unknow err, ", err)
			}
		} else {
			bdsv.Value = []byte(result)
			err = bdsv.Encode()
			if err != nil {
				logrus.Fatal("Incr unknow err, ", err)
			}
		}

		err = txn.Set([]byte(key), bdsv.buff)
		if err != nil {
			logrus.Fatal("badgerDB incr set value error, ", err)
		}
		return nil
	})

	return resultInt64, err
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
			value, err := bDB.Get(string(k))
			if err != nil {
				return err
			}
			if value != nil {
				result = append(result, string(k))
			}
		}
		return nil
	})
	return result, nil
}

func (bDB *BadgerDBStore) SetEx(key string, value string, nowTime, timeout int64) (string, error) {
	stringValue, err := NewBadgerDBValue([]byte(value), nowTime/1e6, timeout)
	if err != nil {
		return "", err
	}
	err = bDB.DB.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), stringValue.buff)
	})
	if err != nil {
		logrus.Fatal("badgerDB set error, ", err)
	}
	return "ok", nil
}

func (bDB *BadgerDBStore) TTL(key string) (int64, error) {
	var result *BadgerDBValue
	err := bDB.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			bdsv := BadgerDBValue{}
			result = &bdsv
			return bdsv.Decode(val)
		})

	})
	if err != nil {
		switch err {
		case badger.ErrKeyNotFound:
			return -1, nil
		default:
			logrus.Fatal("badgerDB get error, ", err)
		}
	}

	nowTime := time.Now().UnixNano() / 1e6
	exTime := result.TimeStamp + result.Timeout
	if result.Timeout != -1 && nowTime > exTime {
		return -1, nil
	}
	return exTime - nowTime, err
}

func (bDB *BadgerDBStore) SetNx(key string, value string, nowTime int64) (int, error) {
	result := 0
	stringValue, err := NewBadgerDBValue([]byte(value), nowTime/1e6, -1)
	if err != nil {
		return 0, err
	}
	err = bDB.DB.Update(func(txn *badger.Txn) error {
		_, err := txn.Get([]byte(key))
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
				err = txn.Set([]byte(key), stringValue.buff)
				if err != nil {
					logrus.Fatal("badgerDB incr set value error, ", err)
				}
				result = 1
				return nil
			default:
				logrus.Fatal("badgerDB incr get error, ", err)
			}
		}
		return nil
	})
	return result, err
}

func (bDB *BadgerDBStore) Expire(key string, nowTime, timeout int64) (int, error) {
	result := 0
	err := bDB.DB.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			switch err {
			case badger.ErrKeyNotFound:
			default:
				logrus.Fatal("badgerDB incr get error, ", err)
			}
		} else {
			var bdsv *BadgerDBValue
			err := item.Value(func(val []byte) error {
				bdsv = &BadgerDBValue{}
				return bdsv.Decode(val)
			})
			if err != nil {
				logrus.Fatal("Expire: read value error, ", err)
			}
			bdsv.TimeStamp = nowTime / 1e6
			bdsv.Timeout = timeout
			err = bdsv.Encode()
			if err != nil {
				logrus.Fatal("Expire: encode value error, ", err)
			}
			err = txn.Set(item.Key(), bdsv.buff)
			if err != nil {
				logrus.Fatal("Expire: write value error, ", err)
			}
			result = 1
		}
		return nil
	})
	return result, err
}

//No need snapshot
func (bDB *BadgerDBStore) Snapshot() ([]byte, error) {
	return nil, nil
}

//No need load snapshot
func (bDB *BadgerDBStore) LoadSnapshot(data []byte) error {
	return nil
}
