/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 23:33:20
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-13 02:12:01
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
	future := raft.RaftServer.Raft.Apply([]byte("heelo"), 10*time.Second)
	err := future.Error()
	if err != nil {
		w.Write([]byte(err.Error()))
	} else {
		w.Write([]byte(future.Response().(string)))
	}
}

func Init() {
	logrus.Info("Http Test Server is start")
	http.HandleFunc("/hello", HelloServer)
	err := http.ListenAndServe(config.Config.HttpServerAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
