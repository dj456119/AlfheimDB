/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 18:43:19
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-02 22:48:28
 */

package store

import (
	"encoding/json"
	"strconv"
	"strings"
	"sync"
)

type MemStoreDatabase struct {
	RWMutex     *sync.RWMutex
	StringStore map[string]string `json:"stringstore"`
}

func NewMemStoreDatabase() *MemStoreDatabase {
	return &MemStoreDatabase{RWMutex: new(sync.RWMutex), StringStore: make(map[string]string)}
}

func (memStore *MemStoreDatabase) Set(key, value string) error {
	memStore.RWMutex.Lock()
	memStore.StringStore[key] = value
	memStore.RWMutex.Unlock()
	return nil
}

func (memStore *MemStoreDatabase) Get(key string) (string, error) {
	var result string
	memStore.RWMutex.RLock()
	result = memStore.StringStore[key]
	memStore.RWMutex.RUnlock()
	return result, nil
}

func (memStore *MemStoreDatabase) Incr(key string) (string, error) {
	memStore.RWMutex.Lock()
	if value, ok := memStore.StringStore[key]; ok {
		int64value, err := strconv.ParseInt(value, 10, 64)
		if err != nil {
			memStore.RWMutex.Unlock()
			return "", err
		}
		int64value = int64value + 1
		result := strconv.FormatInt(int64value, 10)
		memStore.StringStore[key] = result
		memStore.RWMutex.Unlock()
		return result, nil
	} else {
		memStore.StringStore[key] = "1"
		memStore.RWMutex.Unlock()
		return "1", nil
	}
}

func (memStore *MemStoreDatabase) Del(key string) error {
	memStore.RWMutex.Lock()
	delete(memStore.StringStore, key)
	memStore.RWMutex.Unlock()
	return nil
}

func (memStore *MemStoreDatabase) Snapshot() ([]byte, error) {
	memStore.RWMutex.RLock()
	buff, err := json.Marshal(memStore.StringStore)
	memStore.RWMutex.RUnlock()
	return buff, err
}

func (memStore *MemStoreDatabase) LoadSnapshot(data []byte) error {
	stringstore := make(map[string]string)
	memStore.RWMutex.Lock()
	err := json.Unmarshal(data, &stringstore)
	memStore.RWMutex.Unlock()
	return err
}

func (memStore *MemStoreDatabase) Keys(prefix string) ([]string, error) {
	memStore.RWMutex.RLock()
	defer memStore.RWMutex.RUnlock()
	result := []string{}
	for k, v := range memStore.StringStore {
		if strings.HasPrefix(k, prefix) {
			result = append(result, v)
		}
	}
	return result, nil
}
