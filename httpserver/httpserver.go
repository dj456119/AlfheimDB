/*
 * @Descripttion:Http test server.
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-11 23:33:20
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-04 13:36:09
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

const (
	HTTP_TEST_APPLY_TIMEOUT = 10 * time.Second
)

//The http test server only supprot test command to test raft-fsm.
func TestServer(w http.ResponseWriter, req *http.Request) {
	future := raft.RaftServer.Raft.Apply([]byte("test"), HTTP_TEST_APPLY_TIMEOUT)
	err := future.Error()
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	} else {
		respInter := future.Response()
		if respInter == nil {
			w.Write([]byte("Raft fsm error"))
			return
		}
		resp := respInter.(raft.FsmResponse)
		if resp.Error != nil {
			w.Write([]byte(resp.Error.Error()))
			return
		}
		w.Write([]byte(resp.Data.(string)))
	}
}

func PingServer(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("pong"))
}

func Init() {
	logrus.Info("Http Test Server is start")
	http.HandleFunc("/test", TestServer)
	http.HandleFunc("/ping", PingServer)
	err := http.ListenAndServe(config.Config.HttpServerAddr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
