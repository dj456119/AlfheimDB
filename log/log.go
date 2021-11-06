/*
 * @Descripttion:Log文件
 * @version:
 * @Author: cm.d
 * @Date: 2021-10-16 16:13:36
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-06 18:57:33
 */

package log

import (
	"os"

	"github.com/AlfheimDB/config"
	"github.com/sirupsen/logrus"
)

func Init() {
	//	logFile := "./runtime-log/runtime.log"
	logrus.SetFormatter(&logrus.TextFormatter{})
	logrus.SetOutput(os.Stdout)
	// writer, _ := rotatelogs.New(
	// 	logFile+".%Y%m%d",
	// 	rotatelogs.WithLinkName(logFile),
	// 	rotatelogs.WithMaxAge(time.Duration(72)*time.Hour),
	// )
	// logrus.SetOutput(writer)
	logLevel := config.Config.LogLevel
	switch logLevel {
	case "info":
		logrus.SetLevel(logrus.InfoLevel)
	case "debug":
		logrus.SetLevel(logrus.DebugLevel)
	case "error":
		logrus.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logrus.SetLevel(logrus.FatalLevel)
	default:
		logrus.SetLevel(logrus.InfoLevel)
	}

	logrus.Info("Log module is start, the log level is ", logLevel)
}
