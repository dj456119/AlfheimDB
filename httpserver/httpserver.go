/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 23:33:20
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-12 13:31:59
 */
package httpserver

import (
	"log"
	"net/http"
	"time"

	"github.com/AlfheimDB/config"
	"github.com/AlfheimDB/raft"
	"github.com/sirupsen/logrus"
)

func HelloServer(w http.ResponseWriter, req *http.Request) {
	logrus.Debug(raft.RaftServer)
	logrus.Debug(raft.RaftServer.Raft)
	raft.RaftServer.Raft.Apply([]byte("heelo"), 10*time.Second)
}

func Init() {
	logrus.Info("Http Test Server is start")
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(config.Config.HttpServerAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
