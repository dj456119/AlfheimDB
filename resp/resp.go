/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 01:06:51
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-14 15:53:42
 */
package resp

import (
	"strings"
	"time"

	"github.com/AlfheimDB/config"
	"github.com/AlfheimDB/raft"
	"github.com/AlfheimDB/store"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/redcon"
)

const (
	RESP_TEST_TIMEOUT = 3 * time.Second
	RESP_SET_TIMEOUT  = 10 * time.Second
	RESP_GET_TIMEOUT  = 10 * time.Second
	RESP_INCR_TIMEOUT = 10 * time.Second
)

func Init() {
	logrus.Info("Start resp server, ", config.Config.RespServerAddr)
	err := redcon.ListenAndServe(config.Config.RespServerAddr, CommandExec, Accept, Close)
	if err != nil {
		logrus.Fatal("Resp server start error, ", err)
	}
	logrus.Info("Start resp server done")
}

func CommandExec(conn redcon.Conn, cmd redcon.Command) {
	switch strings.ToLower(string(cmd.Args[0])) {
	default:
		conn.WriteError("ERR unknown command '" + string(cmd.Args[0]) + "'")
	case "ping":
		pingCommand(conn, cmd)
	case "quit":
		quitCommand(conn, cmd)
	case "test":
		testCommand(conn, cmd)
	case "incr":
		incrCommand(conn, cmd)
	case "set":
		setCommand(conn, cmd)
	case "get":
		getCommand(conn, cmd)
	case "del":
		delCommand(conn, cmd)
	}
}

func Accept(conn redcon.Conn) bool {
	logrus.Debug("Accept client access, address: ", conn.NetConn().RemoteAddr())
	return true
}

func Close(conn redcon.Conn, err error) {
	logrus.Debug("Close client, address: ", conn.NetConn().RemoteAddr())
}

func testCommand(conn redcon.Conn, cmd redcon.Command) {
	future := raft.RaftServer.Raft.Apply(cmd.Raw, RESP_TEST_TIMEOUT)
	err := future.Error()
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteString(future.Response().(string))
}

func incrCommand(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	future := raft.RaftServer.Raft.Apply(cmd.Raw, RESP_INCR_TIMEOUT)
	err := future.Error()
	if err != nil {
		conn.WriteError(err.Error())
		return
	}

	result := future.Response().([]interface{})
	if result[1] != nil {
		conn.WriteError(result[1].(error).Error())
		return
	}
	conn.WriteString(result[0].(string))
}

func pingCommand(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("pong")
}

func quitCommand(conn redcon.Conn, cmd redcon.Command) {
	conn.WriteString("quit ok")
	conn.Close()
}

func setCommand(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 3 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	future := raft.RaftServer.Raft.Apply(cmd.Raw, RESP_SET_TIMEOUT)
	err := future.Error()
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteString(future.Response().(string))
}

func getCommand(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	result := store.ADBStore.Get(string(cmd.Args[1]))
	conn.WriteString(result)
}

func delCommand(conn redcon.Conn, cmd redcon.Command) {
	if len(cmd.Args) < 2 {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	future := raft.RaftServer.Raft.Apply(cmd.Raw, RESP_SET_TIMEOUT)
	err := future.Error()
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	conn.WriteString(future.Response().(string))
}
