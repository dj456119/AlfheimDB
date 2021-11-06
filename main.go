/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-06 18:48:23
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-06 18:59:11
 */

package main

import (
	"github.com/AlfheimDB/config"
	"github.com/AlfheimDB/log"
	"github.com/sirupsen/logrus"
)

func init() {
	config.Init()
	log.Init()
}

func main() {
	logrus.Info("AlfheimDB is start")
}
