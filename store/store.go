/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 18:05:16
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-15 13:57:35
 */
package store

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
	ADBStore = NewSyncMemStoreDatabase()
}
