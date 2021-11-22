/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-22 11:39:04
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-22 15:31:11
 */
package wal

import (
	"github.com/AlfheimDB/config"
	"github.com/AlfheimDB/pb"
	alfheimdbwal "github.com/dj456119/AlfheimDB-WAL"
	"github.com/golang/protobuf/ptypes"
	"github.com/hashicorp/raft"
	"google.golang.org/protobuf/proto"
)

type AlheimDBWALRaftEngine struct {
	WAL       *alfheimdbwal.AlfheimDBWAL
	WriteBuff []byte
}

func NewAlheimDBWALRaftEngine(basedir string) raft.LogStore {
	engine := new(AlheimDBWALRaftEngine)
	engine.WAL = alfheimdbwal.NewWAL(basedir)
	engine.WriteBuff = make([]byte, 1024*100)
	return engine
}

// FirstIndex returns the first index written. 0 for no entries.
func (ae AlheimDBWALRaftEngine) FirstIndex() (uint64, error) {
	if ae.WAL.MinIndex == -1 {
		return 0, nil
	}
	return uint64(ae.WAL.MinIndex), nil
}

// LastIndex returns the last index written. 0 for no entries.
func (ae AlheimDBWALRaftEngine) LastIndex() (uint64, error) {
	return uint64(ae.WAL.MaxIndex), nil
}

// GetLog gets a log entry at a given index.
func (ae AlheimDBWALRaftEngine) GetLog(index uint64, log *raft.Log) error {
	buff := ae.WAL.GetLog(int64(index))
	return Decode(buff, log)
}

// StoreLog stores a log entry.
func (ae AlheimDBWALRaftEngine) StoreLog(log *raft.Log) error {
	buff, err := Encode(log)
	if err != nil {
		return err
	}
	lItem := alfheimdbwal.NewLogItemBuff(int64(log.Index), buff, ae.WriteBuff, config.Config.IsBigEndian)
	ae.WAL.WriteLog(lItem, ae.WriteBuff[:8+8+len(buff)])
	return nil
}

// StoreLogs stores multiple log entries.
func (ae AlheimDBWALRaftEngine) StoreLogs(logs []*raft.Log) error {
	var length int64 = 0
	for _, log := range logs {
		length = length + int64(len(log.Data))
	}
	pos := 0
	lItems := make([]*alfheimdbwal.LogItem, len(logs))
	for i, log := range logs {
		buff, err := Encode(log)
		if err != nil {
			return err
		}
		lItem := alfheimdbwal.NewLogItemBuff(int64(log.Index), buff, ae.WriteBuff[pos:], config.Config.IsBigEndian)
		lItems[i] = lItem
		pos = pos + 8 + 8 + len(buff)
	}
	ae.WAL.BatchWriteLog(lItems, ae.WriteBuff[:pos])
	return nil
}

// DeleteRange deletes a range of log entries. The range is inclusive.
func (ae AlheimDBWALRaftEngine) DeleteRange(min, max uint64) error {
	ae.WAL.TruncateLog(ae.WAL.MinIndex, ae.WAL.MaxIndex)
	return nil
}

func Decode(buff []byte, log *raft.Log) error {
	raftLogPB := pb.RaftLog{}
	err := proto.Unmarshal(buff, &raftLogPB)
	if err != nil {
		return err
	}
	log.AppendedAt, err = ptypes.Timestamp(raftLogPB.AppendedAt)

	if err != nil {
		return nil
	}
	log.Data = raftLogPB.Data
	log.Extensions = raftLogPB.Extensions
	log.Index = raftLogPB.Index
	log.Term = raftLogPB.Term
	log.Type = raft.LogType(raftLogPB.Type)
	return nil
}

func Encode(log *raft.Log) ([]byte, error) {
	logTimeStampProto, err := ptypes.TimestampProto(log.AppendedAt)
	if err != nil {
		return nil, err
	}
	raftLogPB := pb.RaftLog{
		Index:      log.Index,
		Term:       log.Term,
		Data:       log.Data,
		Extensions: log.Extensions,
		Type:       uint32(log.Type),
		AppendedAt: logTimeStampProto,
	}
	return proto.Marshal(&raftLogPB)
}
