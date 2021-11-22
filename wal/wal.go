/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-22 11:33:51
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-22 12:58:30
 */
package wal

import (
	"path/filepath"

	"github.com/AlfheimDB/config"
	raftbadger "github.com/BBVA/raft-badger"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
)

func NewWAL(baseDir string) (raft.LogStore, error) {
	switch config.Config.WALEngine {
	case "badger":
		return raftbadger.NewBadgerStore(filepath.Join(baseDir, "logs.dat"))
	case "alfheimdbwal":
		return NewAlheimDBWALRaftEngine(baseDir), nil
	default:
		logrus.Fatal("Unknow wal engine: ", config.Config.WALEngine)
	}
	return nil, nil
}
