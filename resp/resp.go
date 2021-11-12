/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 01:06:51
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-13 02:20:04
 */
package resp

import (
	"strings"
	"time"

	"github.com/AlfheimDB/config"
	"github.com/AlfheimDB/raft"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/redcon"
)

func Init() {
	logrus.Info("Start resp server, ", config.Config.RespServerAddr)
	err := redcon.ListenAndServe(config.Config.RespServerAddr, CommandExec, Accept, Close)
	if err != nil {
		logrus.Fatal("Resp server start error, ", err)
	}
	logrus.Info("Start resp server done")
}

var count int

func CommandExec(conn redcon.Conn, cmd redcon.Command) {
	switch strings.ToLower(string(cmd.Args[0])) {
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	case "ping":
		conn.WriteString("PONG")
	case "quit":
		conn.WriteString("OK")
		conn.Close()
	case "test":
		raft.RaftServer.Raft.Apply([]byte("incr"), 3*time.Second)
		conn.WriteString("ok")
	case "incr":
		count++
		//	fmt.Println("aaa ", count)
		future := raft.RaftServer.Raft.Apply([]byte("incr"), 30*time.Second)
		err := future.Error()
		if err != nil {
			conn.WriteError(err.Error())
			return
		}
		conn.WriteString(future.Response().(string))
	case "set":
		if len(cmd.Args) != 3 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}

		conn.WriteString("OK")
	case "get":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
	case "del":
		if len(cmd.Args) != 2 {
			conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
			return
		}
	}
}

func Accept(conn redcon.Conn) bool {
	return true
}

func Close(conn redcon.Conn, err error) {
}
