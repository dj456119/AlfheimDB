/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-10-16 17:50:03
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-12 22:16:12
 */
package config

import (
	"flag"
	"strings"

	"github.com/jinzhu/configor"
	"github.com/sirupsen/logrus"
)

type GTConfig struct {
	LogLevel       string `default:"info"`
	RaftDir        string `default:"data/"`
	RaftId         string
	RaftAddr       string `default:"localhost:50011"`
	HttpServerAddr string `default:"localhost:12345"`
	RaftCluster    []string
}

var Config GTConfig

func Init() {
	configor.Load(&Config, "config.yaml")
	myAddr := flag.String("raft_addr", Config.RaftAddr, "TCP host+port for this node")
	raftId := flag.String("raft_id", Config.RaftId, "Node id used by Raft")
	raftCluster := flag.String("raft_cluster", "", "Raft cluster list")
	httpserverAddr := flag.String("httpserver_addr", Config.HttpServerAddr, "Http test server addr")
	flag.Parse()
	if *raftCluster == "" {
		logrus.Fatal("Raft cluster is empty!")
	}
	Config.RaftCluster = strings.Split(*raftCluster, ",")
	Config.RaftAddr = *myAddr
	Config.RaftId = *raftId
	Config.HttpServerAddr = *httpserverAddr
	logrus.Info("Init config ok, ", Config)
}
