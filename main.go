/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-06 18:48:23
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-12 21:44:40
 */

package main

import (
	"github.com/AlfheimDB/config"
	"github.com/AlfheimDB/httpserver"
	"github.com/AlfheimDB/log"
	"github.com/AlfheimDB/raft"
	"github.com/sirupsen/logrus"
)

func init() {
	config.Init()
	log.Init()
	go httpserver.Init()
	raft.Init()

}

func main() {
	logrus.Info("AlfheimDB is start")
}
