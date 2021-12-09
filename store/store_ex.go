/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-12-04 13:55:21
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-06 16:53:25
 */
package store

import (
	"time"

	"github.com/sirupsen/logrus"
)

type DataEx struct {
	Callback func(key string) error
	Des      DataExStore
	IsLeader bool
}

type DataExStore interface {
	Set(key string, dev DataExValue) error
	Get(key string) (*DataExValue, error)
	GetExKeys(nowtime time.Duration) error
}

type DataExValue struct {
	Key       string
	When      time.Duration
	Timeout   time.Duration
	TimeoutAt time.Duration
}

func NewDataEx(callback func(key string) error) *DataEx {
	return &DataEx{
		Callback: callback,
	}
}

type Void struct{}

func (de *DataEx) SetEx(key string, timeout time.Duration) {
	nowTime := time.Duration(time.Now().UnixNano()) / 1e6
	dev := DataExValue{Key: key, When: nowTime, Timeout: timeout, TimeoutAt: nowTime + timeout}
	err := de.Des.Set(key, dev)
	if err != nil {
		logrus.Fatal("Set expire time error, ", err)
	}
}

func (de *DataEx) TTL(key string) (time.Duration, error) {
	dev, err := de.Des.Get(key)
	if err != nil {
		return -1, err
	}
	nowTime := time.Duration(time.Now().UnixNano()) / 1e6
	return dev.TimeoutAt - nowTime, nil
}

// func (de *DataEx) Loop() {
// 	for {
// 		if de.IsLeader {

// 		}
// 	}
// }
