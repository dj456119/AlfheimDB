/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 18:05:16
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-16 11:34:47
 */
package store

import (
	"github.com/AlfheimDB/config"
	"github.com/sirupsen/logrus"
)

var ADBStore AlfheimdbStore

type AlfheimdbStore interface {
	Set(string, string) string
	Get(string) string
	Incr(string) (string, error)
	Del(string) string
	Snapshot() ([]byte, error)
	LoadSnapshot(data []byte) error
}

func Init() {
	switch config.Config.StoreEngine {
	case "syncmap":
		ADBStore = NewSyncMemStoreDatabase()
	case "map":
		ADBStore = NewMemStoreDatabase()
	default:
		logrus.Fatal("Unknow store engine")
	}

}
