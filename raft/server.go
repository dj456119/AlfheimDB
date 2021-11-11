/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 18:00:19
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-11 21:38:32
 */

package raft

import (
	"net"
	"os"
	"path/filepath"

	"github.com/AlfheimDB/config"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/sirupsen/logrus"
)

var RaftServer *AlfheimRaftServer

type AlfheimRaftServer struct {
	RaftId  string
	MyIP    string
	MyPort  string
	RaftDir string
	RaftFsm raft.FSM
}

func Init() {

}

func New(address string, raftDir string, raftId string) *AlfheimRaftServer {
	RaftServer = new(AlfheimRaftServer)
	ip, port, err := net.SplitHostPort(config.Config.RaftAddr)
	if err != nil {
		logrus.Fatal("Unknow ip and port", config.Config.RaftAddr)
	}
	logrus.Info("Bind address Ip: ", ip, " ,port: ", port)
	RaftServer.MyIP = ip
	RaftServer.MyPort = port
	RaftServer.RaftId = raftId
	RaftServer.RaftDir = raftDir
	sock, err := net.Listen("tcp", config.Config.RaftAddr)
	if err != nil {
		logrus.Fatal("Listen port error", err)
	}
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(raftId)
	baseDir := filepath.Join(raftDir, raftId)
	ldb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "logs.dat"))
	if err != nil {
		logrus.Fatal("Init log db error", err)
	}
	sdb, err := boltdb.NewBoltStore(filepath.Join(baseDir, "stable.dat"))
	if err != nil {
		logrus.Fatal("Init stable db error", err)
	}
	fss, err := raft.NewFileSnapshotStore(baseDir, 1, os.Stderr)
	if err != nil {
		logrus.Fatal("Init snapshot dir error", err)
	}

	fsm := AlfheimRaftFSMImpl{}
	RaftServer.RaftFsm = &fsm
	raftIns, err := raft.NewRaft(raftConfig, &fsm, ldb, sdb, fss, nil)
	if err != nil {
		logrus.Fatal("Init raft instance error", err)
	}
	cfg := raft.Configuration{
		Servers: []raft.Server{
			{
				Suffrage: raft.Voter,
				ID:       raft.ServerID(raftId),
				Address:  raft.ServerAddress(address),
			},
		},
	}
	raftFuture := raftIns.BootstrapCluster(cfg)
	if err := raftFuture.Error(); err != nil {
		logrus.Fatal("Bootstrap raft cluster error")
	}

}
