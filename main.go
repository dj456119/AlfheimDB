/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-06 18:48:23
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-13 02:50:24
 */

package main

import (
	"github.com/AlfheimDB/config"
	"github.com/AlfheimDB/httpserver"
	"github.com/AlfheimDB/log"
	"github.com/AlfheimDB/raft"
	"github.com/AlfheimDB/resp"
	"github.com/sirupsen/logrus"
)

func init() {
	config.Init()
	log.Init()

	go resp.Init()
	raft.Init()
	httpserver.Init()
}

func main() {
	logrus.Info("AlfheimDB is start")
}
