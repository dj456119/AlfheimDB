/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-10-16 17:50:03
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-11 18:05:06
 */
package config

import "github.com/jinzhu/configor"

type GTConfig struct {
	LogLevel string `default:"info"`
	RaftDir  string `default:"data/"`
	RaftId   string
	RaftAddr string `default:"localhost:50011"`
}

var Config GTConfig

func Init() {
	configor.Load(&Config, "config.yaml")
}
