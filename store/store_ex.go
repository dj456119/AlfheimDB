/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-12-04 13:55:21
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-04 14:15:22
 */
package store

import (
	"sync"
	"time"
)

type DataEx struct {
	Pool     *sync.Map
	loopChan chan *DataExValue
}

type DataExValue struct {
	Key      string
	When     time.Duration
	Timeout  time.Duration
	Callback func(key string) error
}

func NewDataEx() *DataEx {
	return &DataEx{
		Pool:     new(sync.Map),
		loopChan: make(chan *DataExValue),
	}
}

type Void struct{}

func (de *DataEx) SetEx(key string, timeout time.Duration, callback func(key string) error) {
	dev := DataExValue{Key: key, When: time.Duration(time.Now().UnixNano()) / 1e6, Timeout: timeout, Callback: callback}
	de.Pool.Store(key, Void{})
	time.AfterFunc(dev.Timeout, func() {
		dev.Callback(key)
	})

}
