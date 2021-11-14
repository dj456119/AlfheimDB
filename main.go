/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-06 18:48:23
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-14 12:46:00
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

func init() {
	config.Init()
	log.Init()
	store.Init()
	go resp.Init()
	go raft.Init()
	httpserver.Init()
}

func main() {
	logrus.Info("AlfheimDB is start")
}
