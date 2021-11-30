/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 18:00:19
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-30 22:34:28
 */

package raft

import (
	"bytes"
	"fmt"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/AlfheimDB/config"
	"github.com/AlfheimDB/log"
	"github.com/AlfheimDB/wal"
	raftbadger "github.com/BBVA/raft-badger"

	"github.com/hashicorp/raft"
	"github.com/sirupsen/logrus"
)

type AlfheimRaftServer struct {
	RaftId  string
	MyIP    string
	MyPort  string
	RaftDir string

	RaftFsm     raft.BatchingFSM
	Raft        *raft.Raft
	RaftCluster []string
}

var RaftServer *AlfheimRaftServer

func Init() {
	logrus.Info("Init raft server ")
	initRaft(config.Config.RaftAddr, config.Config.RaftDir, config.Config.RaftId)

}

func initRaft(address string, raftDir string, raftId string) {
	RaftServer = new(AlfheimRaftServer)
	RaftServer.RaftCluster = config.Config.RaftCluster

	ip, port, err := net.SplitHostPort(config.Config.RaftAddr)
	if err != nil {
		logrus.Fatal("Unknow ip and port", config.Config.RaftAddr)
	}
	logrus.Info("Raft address Ip: ", ip, " ,port: ", port)

	RaftServer.MyIP = ip
	RaftServer.MyPort = port
	RaftServer.RaftId = raftId
	RaftServer.RaftDir = raftDir

	//raft config
	raftConfig := raft.DefaultConfig()
	raftConfig.LocalID = raft.ServerID(raftId)
	raftConfig.BatchApplyCh = true
	raftConfig.MaxAppendEntries = config.Config.RaftMaxAppendEntris
	raftConfig.TrailingLogs = config.Config.RaftTrailingLogs
	raftConfig.LogLevel = config.Config.LogLevel
	raftConfig.LogOutput = log.LogWriter

	//If umask is not 0, need chmod
	shell(fmt.Sprintf("mkdir %s", config.Config.BaseDir))

	//Init log db
	ldb, err := wal.NewWAL(config.Config.BaseDir)
	if err != nil {
		logrus.Fatal("Init log db error, ", err)
	}

	//Init stable db
	sdb, err := raftbadger.NewBadgerStore(filepath.Join(config.Config.BaseDir, "stable.dat"))
	if err != nil {
		logrus.Fatal("Init stable db error, ", err)
	}

	//Init file snapshot store
	fss, err := raft.NewFileSnapshotStore(config.Config.BaseDir, 1, os.Stderr)
	if err != nil {
		logrus.Fatal("Init snapshot dir error, ", err)
	}

	//Create raft fsm
	RaftServer.RaftFsm = NewAlfheimRaftFSM()

	addr, err := net.ResolveTCPAddr("tcp", config.Config.RaftAddr)
	if err != nil {
		logrus.Fatal("Raft addr resolve error, ", err)
	}

	//Use default net transport in raft lib
	transport, err := raft.NewTCPTransport(config.Config.RaftAddr, addr, 100, 5*time.Second, os.Stderr)
	if err != nil {
		logrus.Fatal("Raft addr create error, ", err)
	}

	//Create raft instance
	raftIns, err := raft.NewRaft(raftConfig, RaftServer.RaftFsm, ldb, sdb, fss, transport)
	if err != nil {
		logrus.Fatal("Init raft instance error, ", err)
	}
	RaftServer.Raft = raftIns
	RaftServer.Bootstrap()

}

func shell(command string) {
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd := exec.Command("bash", "-c", command)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		logrus.Info("err cmd :", command)
		logrus.Info("stdout: ", stdout.String())
		logrus.Info("stderr: ", stderr.String())
		logrus.Info(err)
	}
	logrus.Info("exec cmd done:", command)
	logrus.Info("stdout: ", stdout.String())
	logrus.Info("stderr: ", stderr.String())
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
		logrus.Fatal("Bootstrap raft cluster error, ", err)
	}
	logrus.Info("Bootstrap done")
}
