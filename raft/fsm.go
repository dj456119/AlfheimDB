/*
 * @Descripttion: AlfheimDB raft fsm core implement raft's interface.
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 11:19:46
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-30 21:18:19
 */

package raft

import (
	"errors"
	"io"
	"io/ioutil"
	"strings"

	"github.com/AlfheimDB/store"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/redcon"
)

//Fsm apply response struct
type FsmResponse struct {
	Data  string
	Error error
}

type AlfheimRaftFSMImpl struct {
}

func NewAlfheimRaftFSM() raft.BatchingFSM {
	return new(AlfheimRaftFSMImpl)
}

//Single apply
func (aFsm *AlfheimRaftFSMImpl) Apply(l *raft.Log) interface{} {
	cmd, err := redcon.Parse(l.Data)
	if err != nil {
		return err.Error
	}
	cmdLow := strings.ToLower(string(cmd.Args[0]))
	switch string(cmdLow) {
	case "test":
		return FsmResponse{Data: "test ok"}
	case "set":
		err := store.ADBStore.Set(string(cmd.Args[1]), string(cmd.Args[2]))
		return FsmResponse{Error: err}
	case "del":
		err := store.ADBStore.Del(string(cmd.Args[1]))
		return FsmResponse{Error: err}
	case "incr":
		result, err := store.ADBStore.Incr(string(cmd.Args[1]))
		return FsmResponse{Data: result, Error: err}
	default:
		return FsmResponse{Error: errors.New("Unknow command, " + cmdLow)}
	}

}

//Batch apply, just run apply multiple times
func (aFsm *AlfheimRaftFSMImpl) ApplyBatch(logs []*raft.Log) []interface{} {
	result := make([]interface{}, len(logs))
	for i, l := range logs {
		result[i] = aFsm.Apply(l)
	}
	return result
}

//Snapshot
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

//Load snapshot
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
