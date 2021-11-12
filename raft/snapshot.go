/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 12:15:28
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-12 22:37:19
 */
package raft

import (
	"github.com/hashicorp/raft"
)

type AlfheimRaftSnapshot struct {
	SnapshotBytes []byte
}

func (s *AlfheimRaftSnapshot) Persist(sink raft.SnapshotSink) error {
	sink.Write(s.SnapshotBytes)
	return sink.Close()
}

func (s *AlfheimRaftSnapshot) Release() {
}
