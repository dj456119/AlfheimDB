/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-10-16 17:50:03
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-22 11:36:01
 */
package config

import (
	"flag"
	"log"
	"strings"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

type GTConfig struct {
	LogLevel             string `default:"info"`
	RaftDir              string `default:"data/"`
	RaftId               string
	RaftAddr             string `default:"localhost:50011"`
	HttpServerAddr       string `default:"localhost:12345"`
	RaftCluster          []string
	RespServerAddr       string
	RaftSnapshotInterval int    `default:"1"`
	RaftMaxAppendEntris  int    `default:"1000"`
	RaftTrailingLogs     uint64 `default:"1024000"`
	LogType              string `default:"stdout"`
	StoreEngine          string `default:"syncmap"`
	WALEngine            string `default:"badger"`
}

var Config GTConfig

func Init() {
	configor.Load(&Config, "config.yaml")

	myAddr := flag.String("raft_addr", Config.RaftAddr, "TCP host+port for this node")
	raftId := flag.String("raft_id", Config.RaftId, "Node id used by Raft")
	raftCluster := flag.String("raft_cluster", "", "Raft cluster list")
	httpserverAddr := flag.String("httpserver_addr", Config.HttpServerAddr, "Http test server addr")
	respserverAddr := flag.String("respserver_addr", Config.RespServerAddr, "Resp server addr")
	logtype := flag.String("logtype", Config.LogType, "Log type, file or stdout")
	flag.Parse()

	if *raftId == "" {
		log.Fatal("Raft id is empty")
	}

	if *raftCluster != "" {
		Config.RaftCluster = strings.Split(*raftCluster, ",")
	}

	if Config.RaftCluster == nil || len(Config.RaftCluster) == 0 {
		log.Fatal("Raft cluster is empty")
	}

	Config.RaftAddr = *myAddr
	Config.RaftId = *raftId
	Config.LogType = *logtype
	if *httpserverAddr == "" {
		logrus.Info("No need http server")
	}
	Config.HttpServerAddr = *httpserverAddr

	if *respserverAddr == "" {
		logrus.Fatal("Resp server addr is empty!")
	}

	Config.RespServerAddr = *respserverAddr
	logrus.Info("Init config ok, ", Config)
}
