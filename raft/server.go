/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 18:00:19
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-13 03:06:13
 */

package raft

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/AlfheimDB/config"
	"github.com/hashicorp/raft"
	boltdb "github.com/hashicorp/raft-boltdb"
	"github.com/sirupsen/logrus"
)

type AlfheimRaftServer struct {
	RaftId      string
	MyIP        string
	MyPort      string
	RaftDir     string
	RaftFsm     raft.BatchingFSM
	Raft        *raft.Raft
	RaftCluster []string
}

var RaftServer *AlfheimRaftServer

func Init() {
	logrus.Info("init raft server ")
	initRaft(config.Config.RaftAddr, config.Config.RaftDir, config.Config.RaftId)

}

func initRaft(address string, raftDir string, raftId string) {
	RaftServer = new(AlfheimRaftServer)
	RaftServer.RaftCluster = config.Config.RaftCluster
	ip, port, err := net.SplitHostPort(config.Config.RaftAddr)
	if err != nil {
		logrus.Fatal("Unknow ip and port", config.Config.RaftAddr)
	}
	logrus.Info("Bind address Ip: ", ip, " ,port: ", port)
	RaftServer.MyIP = ip
	RaftServer.MyPort = port
	RaftServer.RaftId = raftId
	RaftServer.RaftDir = raftDir
	// sock, err := net.Listen("tcp", config.Config.RaftAddr)
	// if err != nil {
	// 	logrus.Fatal("Listen port error", err)
	// }

	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(raftId)
	raftConfig.BatchApplyCh = true
	raftConfig.MaxAppendEntries = 1000
	baseDir := filepath.Join(raftDir, raftId)
	//ldb := raft.NewInmemStore()
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

	fsm := AlfheimRaftFSMImpl{Counter: 0, RWLock: new(sync.RWMutex)}
	RaftServer.RaftFsm = &fsm

	//tm := transport.New(raft.ServerAddress(address), []grpc.DialOption{grpc.WithInsecure()})

	addr, err := net.ResolveTCPAddr("tcp", config.Config.RaftAddr)
	if err != nil {
		logrus.Fatal("Raft addr resolve error", err)
	}
	transport, err := raft.NewTCPTransport(config.Config.RaftAddr, addr, 2, 5*time.Second, os.Stderr)
	if err != nil {
		logrus.Fatal("Raft addr create error", err)
	}

	raftIns, err := raft.NewRaft(raftConfig, &fsm, ldb, sdb, fss, transport)
	if err != nil {
		logrus.Fatal("Init raft instance error", err)
	}
	RaftServer.Raft = raftIns
	RaftServer.Bootstrap()
	// grpcServer := grpc.NewServer()
	// tm.Register(grpcServer)
	// leaderhealth.Setup(raftIns, grpcServer, []string{"Example"})
	// raftadmin.Register(grpcServer, raftIns)
	// reflection.Register(grpcServer)
	// logrus.Info("raft init success")
	// if err := grpcServer.Serve(sock); err != nil {
	// 	logrus.Fatal("Grpc serve sock error, ", err)
	// }
}

func (aServer *AlfheimRaftServer) Bootstrap() {
	servers := aServer.Raft.GetConfiguration().Configuration().Servers
	if len(servers) > 0 {
		logrus.Info("Not first startup, don't need bootstrap")
		return
	}
	logrus.Info("First start, need bootstrap")
	var configuration raft.Configuration
	for _, peerInfo := range aServer.RaftCluster {
		peer := strings.Split(peerInfo, "/")
		id := peer[1]
		addr := peer[0]
		server := raft.Server{
			ID:      raft.ServerID(id),
			Address: raft.ServerAddress(addr),
		}
		configuration.Servers = append(configuration.Servers, server)
	}
	raftFuture := RaftServer.Raft.BootstrapCluster(configuration)
	if err := raftFuture.Error(); err != nil {
		logrus.Fatal("Bootstrap raft cluster error", err)
	}
	logrus.Info("Bootstrap done")
}
