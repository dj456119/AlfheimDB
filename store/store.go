/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 18:05:16
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-09 18:17:38
 */
package store

import (
	"github.com/AlfheimDB/config"
	"github.com/sirupsen/logrus"
)

var ADBStore AlfheimdbStore

type AlfheimdbStore interface {
	//String store
	Set(key string, value string, nowTime int64) error
	Get(key string) (*string, error)
	Incr(key string, nowTime int64) (int64, error)
	Del(key string) error
	Keys(prefix string) ([]string, error)
	SetEx(key string, value string, nowTime int64, time int64) (string, error)
	TTL(key string) (int64, error)
	SetNx(key string, value string, nowTime int64) (int, error)
	Expire(key string, nowTime, timeout int64) (int, error)
	Snapshot() ([]byte, error)
	LoadSnapshot(data []byte) error
}

func Init() {
	switch config.Config.StoreEngine {
	case "map":
		ADBStore = NewMemStoreDatabase()
	case "badger":
		ADBStore = NewBadgerDBStore(config.Config.BaseDir)
	default:
		logrus.Fatal("Unknow store engine")
	}

}
