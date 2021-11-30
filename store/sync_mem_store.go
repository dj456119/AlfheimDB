/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 18:43:19
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-30 21:18:55
 */

package store

import (
	"encoding/json"
	"strconv"
	"sync"
)

type SyncMemStoreDatabase struct {
	StringStore *sync.Map
}

func NewSyncMemStoreDatabase() *SyncMemStoreDatabase {
	return &SyncMemStoreDatabase{StringStore: new(sync.Map)}
}

func (memStore *SyncMemStoreDatabase) Set(key, value string) error {
	memStore.StringStore.Store(key, value)
	return nil
}

func (memStore *SyncMemStoreDatabase) Get(key string) (string, error) {
	v, b := memStore.StringStore.Load(key)
	if !b {
		return "nil", nil
	}
	return v.(string), nil
}

func (memStore *SyncMemStoreDatabase) Incr(key string) (string, error) {
	if value, ok := memStore.StringStore.Load(key); ok {
		int64value, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return "", err
		}
		int64value = int64value + 1
		result := strconv.FormatInt(int64value, 10)
		memStore.StringStore.Store(key, result)

		return result, nil
	} else {
		memStore.StringStore.Store(key, "1")
		return "1", nil
	}
}

func (memStore *SyncMemStoreDatabase) Del(key string) error {
	memStore.StringStore.Delete(key)
	return nil
}

func (memStore *SyncMemStoreDatabase) Snapshot() ([]byte, error) {

	stringstoreMap := make(map[string]string)
	memStore.StringStore.Range(func(k interface{}, v interface{}) bool {
		stringstoreMap[k.(string)] = v.(string)
		return true
	})
	buff, err := json.Marshal(stringstoreMap)
	return buff, err
}

func (memStore *SyncMemStoreDatabase) LoadSnapshot(data []byte) error {
	newSyncMap := new(sync.Map)
	stringstoreMap := make(map[string]string)
	err := json.Unmarshal(data, &stringstoreMap)
	for k, v := range stringstoreMap {
		newSyncMap.Store(k, v)
	}
	memStore.StringStore = newSyncMap
	return err
}
