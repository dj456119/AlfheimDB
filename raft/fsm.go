/*
 * @Descripttion: raft
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 11:19:46
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-13 02:32:04
 */

package raft

import (
	"encoding/binary"
	"fmt"
	"io"
	"sync"

	"github.com/hashicorp/raft"
)

type AlfheimRaftFSMImpl struct {
	Counter uint64
	RWLock  *sync.RWMutex
}

func (aFsm *AlfheimRaftFSMImpl) Apply(l *raft.Log) interface{} {
	//aFsm.RWLock.Lock()
	aFsm.Counter = aFsm.Counter + 1
	if aFsm.Counter%1000 == 0 {
		fmt.Println("nobatch", aFsm.Counter)
	}
	//aFsm.RWLock.Unlock()
	return fmt.Sprintf("%d", aFsm.Counter)
}

func (aFsm *AlfheimRaftFSMImpl) ApplyBatch(logs []*raft.Log) []interface{} {
	result := make([]interface{}, len(logs))
	for i := range logs {
		aFsm.Counter = aFsm.Counter + 1
		if aFsm.Counter%1000 == 0 {
			fmt.Println("batch", aFsm.Counter, len(logs))
		}

		result[i] = fmt.Sprintf("%d", aFsm.Counter)

	}
	return result
}

func (aFsm *AlfheimRaftFSMImpl) Snapshot() (raft.FSMSnapshot, error) {
	b := make([]byte, 8)
	snapshot := new(AlfheimRaftSnapshot)
	aFsm.RWLock.Lock()
	binary.BigEndian.PutUint64(b, aFsm.Counter)
	aFsm.RWLock.Unlock()
	snapshot.SnapshotBytes = b
	return snapshot, nil
}

func (aFsm *AlfheimRaftFSMImpl) Restore(r io.ReadCloser) error {
	b := make([]byte, 8)
	for {
		count, err := r.Read(b)
		if err != nil {
			return err
		}
		if count != 8 {
			continue
		}
		if count == 8 {
			break
		}
	}
	aFsm.Counter = uint64(binary.BigEndian.Uint64(b))
	return nil
}
