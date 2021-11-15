/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 18:43:19
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-15 13:57:39
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

func (memStore *SyncMemStoreDatabase) Set(key, value string) string {
	memStore.StringStore.Store(key, value)
	return "ok"
}

func (memStore *SyncMemStoreDatabase) Get(key string) string {
	v, b := memStore.StringStore.Load(key)
	if !b {
		return "nil"
	}
	return v.(string)
}

func (memStore *SyncMemStoreDatabase) Incr(key string) (string, error) {

	if value, ok := memStore.StringStore.Load(key); ok {
		int64value, err := strconv.ParseInt(value.(string), 10, 64)
		if err != nil {
			return "parse err", err
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

func (memStore *SyncMemStoreDatabase) Del(key string) string {
	memStore.StringStore.Delete(key)
	return "ok"
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
