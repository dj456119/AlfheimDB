/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 18:05:16
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-02 21:48:53
 */
package store

import (
	"github.com/AlfheimDB/config"
	"github.com/sirupsen/logrus"
)

var ADBStore AlfheimdbStore

type AlfheimdbStore interface {
	Set(key string, value string) error
	Get(key string) (string, error)
	Incr(key string) (string, error)
	Del(key string) error
	Keys(prefix string) ([]string, error)
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
