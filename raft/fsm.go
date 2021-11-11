/*
 * @Descripttion: raft
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 11:19:46
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-11 21:36:39
 */

package raft

import (
	"fmt"
	"io"

	"github.com/hashicorp/raft"
)

type AlfheimRaftFSMImpl struct {
	DataMap map[string]string
}

func (aFsm *AlfheimRaftFSMImpl) Apply(l *raft.Log) interface{} {
	fmt.Println(string(l.Data))
	return nil
}

func (aFsm *AlfheimRaftFSMImpl) Snapshot() (raft.FSMSnapshot, error) {

	return nil, nil
}

func (aFsm *AlfheimRaftFSMImpl) Restore(r io.ReadCloser) error {

	return nil
}
