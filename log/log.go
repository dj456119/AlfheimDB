/*
 * @Descripttion:The log
 * @version:
 * @Author: cm.d
 * @Date: 2021-10-16 16:13:36
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-14 20:24:00
 */

package log

import (
	"io"
	"os"
	"time"

	"github.com/AlfheimDB/config"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var LogWriter io.Writer

func Init() {

	logrus.SetFormatter(&logrus.TextFormatter{})
	if config.Config.LogType == "file" {
		logFile := "./runtime-log/runtime.log"
		writer, _ := rotatelogs.New(
			logFile+".%Y%m%d",
			rotatelogs.WithLinkName(logFile),
			rotatelogs.WithMaxAge(time.Duration(72)*time.Hour),
		)
		LogWriter = writer
	} else {
		LogWriter = os.Stdout
	}
	logrus.SetOutput(LogWriter)
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

	logrus.Info("Log module is start, the log type is ", config.Config.LogType, ", log level is ", logLevel)
}
