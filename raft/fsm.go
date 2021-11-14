/*
 * @Descripttion: raft
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 11:19:46
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-14 22:25:57
 */

package raft

import (
	"io"
	"io/ioutil"
	"strings"

	"github.com/AlfheimDB/store"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/redcon"
)

type AlfheimRaftFSMImpl struct {
}

func NewAlfheimRaftFSM() raft.BatchingFSM {
	fsm := new(AlfheimRaftFSMImpl)
	return fsm
}

func (aFsm *AlfheimRaftFSMImpl) Apply(l *raft.Log) interface{} {
	cmd, err := redcon.Parse(l.Data)
	if err != nil {
		return err.Error
	}
	cmdLow := strings.ToLower(string(cmd.Args[0]))
	switch string(cmdLow) {
	case "test":
		return "test ok"
	case "set":
		result := store.ADBStore.Set(string(cmd.Args[1]), string(cmd.Args[2]))
		return result
	case "del":
		return store.ADBStore.Del(string(cmd.Args[1]))
	case "incr":
		result, err := store.ADBStore.Incr(string(cmd.Args[1]))
		return []interface{}{result, err}
	default:
		return "unknow cmd"
	}

}

func (aFsm *AlfheimRaftFSMImpl) ApplyBatch(logs []*raft.Log) []interface{} {
	result := make([]interface{}, len(logs))
	for i, l := range logs {
		result[i] = aFsm.Apply(l)
	}
	return result
}

func (aFsm *AlfheimRaftFSMImpl) Snapshot() (raft.FSMSnapshot, error) {
	logrus.Info("Start create snapshot")
	snapshot := new(AlfheimRaftSnapshot)
	buff, err := store.ADBStore.Snapshot()
	if err != nil {
		logrus.Fatal("Snapshot create error, ", err)
	}
	snapshot.SnapshotBytes = buff
	return snapshot, err
}

func (aFsm *AlfheimRaftFSMImpl) Restore(r io.ReadCloser) error {
	logrus.Info("Start load snapshot")
	buff, err := ioutil.ReadAll(r)
	if err != nil {
		logrus.Fatal("Snapshot load error, ", err)
	}
	err = store.ADBStore.LoadSnapshot(buff)
	if err != nil {
		logrus.Fatal("Snapshot load error, ", err)
	}
	return err
}
