/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-10-16 17:50:03
 * @LastEditors: cm.d
 * @LastEditTime: 2021-10-16 17:59:20
 */
package config

import "github.com/jinzhu/configor"

type GTConfig struct {
	LogLevel string `default:"info"`
}

var Config GTConfig

func Init() {
	configor.Load(&Config, "config.yaml")
}
