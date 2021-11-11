/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 12:15:28
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-11 21:36:34
 */
package raft

import (
	"github.com/hashicorp/raft"
)

type AlfheimRaftSnapshot struct {
	dataBytes []byte
}

func (s *AlfheimRaftSnapshot) Persist(sink raft.SnapshotSink) error {
	return sink.Close()
}

func (s *AlfheimRaftSnapshot) Release() {
}
