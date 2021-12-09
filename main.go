/*
 * @Descripttion:Main
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-06 18:48:23
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-09 18:11:57
 */

package main

import (
	_ "net/http/pprof"

	"github.com/AlfheimDB/config"
	"github.com/AlfheimDB/httpserver"
	"github.com/AlfheimDB/log"
	"github.com/AlfheimDB/raft"
	"github.com/AlfheimDB/resp"
	"github.com/AlfheimDB/store"
	"github.com/sirupsen/logrus"
)

//Module init
func init() {
	//Init config module
	config.Init()

	//Init log module
	log.Init()

	//Init store module
	store.Init()

	//Init redcon server
	go resp.Init()

	//Init raft server
	go raft.Init()

	//Init test http server
	go httpserver.Init()
}

//main
func main() {
	logrus.Info("AlfheimDB is start")
	select {}
}
