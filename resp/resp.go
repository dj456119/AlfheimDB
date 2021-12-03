/*
 * @Descripttion:
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 01:06:51
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-02 22:47:16
 */
package resp

import (
	"errors"
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
		execCommand(1, conn, cmd, func(cmd redcon.Command) raft.FsmResponse {
			return raft.FsmResponse{Error: errors.New("ERR unknown command '" + string(cmd.Args[0]) + "'")}
		})
	case "ping":
		execCommand(1, conn, cmd, func(cmd redcon.Command) raft.FsmResponse {
			return raft.FsmResponse{Data: "pong"}
		})
	case "quit":
		execCommand(1, conn, cmd, func(cmd redcon.Command) raft.FsmResponse {
			return raft.FsmResponse{Data: "quit ok"}
		})
		conn.Close()
	case "test":
		execCommandByFsm(1, conn, cmd, RESP_TEST_TIMEOUT)
	case "incr":
		execCommandByFsm(2, conn, cmd, RESP_INCR_TIMEOUT)
	case "set":
		execCommandByFsm(3, conn, cmd, RESP_SET_TIMEOUT)
	case "get":
		execCommand(2, conn, cmd, func(cmd redcon.Command) raft.FsmResponse {
			data, err := store.ADBStore.Get(string(cmd.Args[1]))
			return raft.FsmResponse{Data: data, Error: err}
		})
	case "del":
		execCommandByFsm(2, conn, cmd, RESP_GET_TIMEOUT)
	}
}

func Accept(conn redcon.Conn) bool {
	logrus.Debug("Accept client access, address: ", conn.NetConn().RemoteAddr())
	return true
}

func Close(conn redcon.Conn, err error) {
	logrus.Debug("Close client, address: ", conn.NetConn().RemoteAddr())
}

func execCommand(argsLength int, conn redcon.Conn, cmd redcon.Command, exec func(cmd redcon.Command) raft.FsmResponse) {
	if len(cmd.Args) < argsLength {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	response := exec(cmd)
	if response.Error != nil {
		conn.WriteError(response.Error.Error())
		return
	}
	//conn.WriteString(response.Data)
	conn.WriteAny(response.Data)
}

func execCommandByFsm(argsLength int, conn redcon.Conn, cmd redcon.Command, timeout time.Duration) {
	if len(cmd.Args) < argsLength {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	future := raft.RaftServer.Raft.Apply(cmd.Raw, timeout)
	err := future.Error()
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	respInterface := future.Response()
	if respInterface == nil {
		conn.WriteError("Unknow error")
		return
	}
	response := respInterface.(raft.FsmResponse)
	if response.Error != nil {
		conn.WriteError(response.Error.Error())
		return
	}
	conn.WriteAny(response.Data)
	//conn.WriteString(response.Data)
}
