/*
 * @Descripttion: AlfheimDB raft fsm core implement raft's interface.
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 11:19:46
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-09 17:56:20
 */

package raft

import (
	"errors"
	"io"
	"io/ioutil"
	"strconv"

	"github.com/AlfheimDB/pb"
	"github.com/AlfheimDB/store"
	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

//Fsm apply response struct
type FsmResponse struct {
	Data  interface{}
	Error error
}

//AlfheimDB raft fsm which implents raft batch-fsm interface
type AlfheimRaftFSMImpl struct {
}

func NewAlfheimRaftFSM() raft.BatchingFSM {
	return new(AlfheimRaftFSMImpl)
}

//The fsm's apply log
type AlfheimRaftFSMLog struct {
	Cmd       string
	Args      [][]byte
	TimeStamp int64
	Buff      []byte
}

//The bytes in apply log convert to AlfheimRaftFSMLog
func (fLog *AlfheimRaftFSMLog) Decode() error {
	pbLog := pb.FsmLog{}
	err := proto.Unmarshal(fLog.Buff, &pbLog)
	if err != nil {
		return err
	}
	fLog.Cmd = pbLog.Cmd
	fLog.Args = pbLog.CmdArgs
	fLog.TimeStamp = pbLog.TimeStamp
	return nil
}

//AlfheimRaftFSMLog convert to bytes
func (fLog *AlfheimRaftFSMLog) Encode() error {
	pbLog := pb.FsmLog{
		Cmd:       fLog.Cmd,
		CmdArgs:   fLog.Args,
		TimeStamp: fLog.TimeStamp,
	}
	var err error
	fLog.Buff, err = proto.Marshal(&pbLog)
	if err != nil {
		return err
	}
	return nil
}

//Single apply
func (aFsm *AlfheimRaftFSMImpl) Apply(l *raft.Log) interface{} {
	fLog := AlfheimRaftFSMLog{Buff: l.Data}
	err := fLog.Decode()
	if err != nil {
		return FsmResponse{Error: err}
	}
	switch fLog.Cmd {
	case "test":
		return FsmResponse{Data: "test ok"}
	case "set":
		err := store.ADBStore.Set(string(fLog.Args[1]), string(fLog.Args[2]), fLog.TimeStamp)
		if err != nil {
			return FsmResponse{Error: err}
		}
		return FsmResponse{Data: "ok"}
	case "del":
		err := store.ADBStore.Del(string(fLog.Args[1]))
		if err != nil {
			return FsmResponse{Error: err}
		}
		return FsmResponse{Data: "ok"}
	case "incr":
		result, err := store.ADBStore.Incr(string(fLog.Args[1]), fLog.TimeStamp)
		return FsmResponse{Data: result, Error: err}
	case "setnx":
		result, err := store.ADBStore.SetNx(string(fLog.Args[1]), string(fLog.Args[2]), fLog.TimeStamp)
		return FsmResponse{Data: result, Error: err}
	case "setex":
		timeout, _ := strconv.ParseInt(string(fLog.Args[3]), 10, 64)
		result, err := store.ADBStore.SetEx(string(fLog.Args[1]), string(fLog.Args[2]), fLog.TimeStamp, timeout)
		return FsmResponse{Data: result, Error: err}
	default:
		return FsmResponse{Error: errors.New("Unknow command, " + fLog.Cmd)}
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
