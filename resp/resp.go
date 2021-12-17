/*
 * @Descripttion: Redcon server, it implents redis protocol and support complicate.
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-13 01:06:51
 * @LastEditors: cm.d
 * @LastEditTime: 2021-12-17 14:16:29
 */
package resp

import (
	"errors"
	"reflect"
	"strconv"
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

//Init redcon server
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
		execCommandByFsm(1, nil, conn, cmd, RESP_TEST_TIMEOUT)
	case "incr":
		execCommandByFsm(2, nil, conn, cmd, RESP_INCR_TIMEOUT)
	case "set":
		execCommandByFsm(3, nil, conn, cmd, RESP_SET_TIMEOUT)
	case "setex":
		execCommandByFsm(4, func(cmd redcon.Command) error {
			_, err := strconv.ParseInt(string(cmd.Args[3]), 10, 64)
			return err
		}, conn, cmd, RESP_SET_TIMEOUT)
	case "setnx":
		execCommandByFsm(3, nil, conn, cmd, RESP_SET_TIMEOUT)
	case "ttl":
		execCommand(2, conn, cmd, func(cmd redcon.Command) raft.FsmResponse {
			result, err := store.ADBStore.TTL(string(cmd.Args[1]))
			return raft.FsmResponse{Data: result, Error: err}
		})
	case "keys":
		execCommand(2, conn, cmd, func(cmd redcon.Command) raft.FsmResponse {
			result, err := store.ADBStore.Keys(string(cmd.Args[1]))
			return raft.FsmResponse{Data: result, Error: err}
		})
	case "get":
		execCommand(2, conn, cmd, func(cmd redcon.Command) raft.FsmResponse {
			data, err := store.ADBStore.Get(string(cmd.Args[1]))
			return raft.FsmResponse{Data: data, Error: err}
		})
	case "del":
		execCommandByFsm(2, nil, conn, cmd, RESP_GET_TIMEOUT)
	}
}

func Accept(conn redcon.Conn) bool {
	logrus.Debug("Accept client access, address: ", conn.NetConn().RemoteAddr())
	return true
}

func Close(conn redcon.Conn, err error) {
	logrus.Debug("Close client, address: ", conn.NetConn().RemoteAddr())
}

//Exec command
func execCommand(argsLength int, conn redcon.Conn, cmd redcon.Command, exec func(cmd redcon.Command) raft.FsmResponse) {
	if len(cmd.Args) < argsLength {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	response := exec(cmd)
	WriteResponse(response, conn)
}

//Exec command by fsm
func execCommandByFsm(argsLength int, validArgs func(cmd redcon.Command) error, conn redcon.Conn, cmd redcon.Command, timeout time.Duration) {
	if len(cmd.Args) < argsLength {
		conn.WriteError("ERR wrong number of arguments for '" + string(cmd.Args[0]) + "' command")
		return
	}
	if validArgs != nil {
		err := validArgs(cmd)
		if err != nil {
			conn.WriteError(err.Error())
			return
		}
	}
	fLog := raft.AlfheimRaftFSMLog{
		TimeStamp: time.Now().UnixNano(),
	}
	fLog.Cmd = strings.ToLower(string(cmd.Args[0]))
	fLog.Args = cmd.Args
	err := fLog.Encode()
	if err != nil {
		conn.WriteError("Cmd encode error, " + err.Error())
	}
	future := raft.RaftServer.Raft.Apply(fLog.Buff, timeout)
	err = future.Error()
	if err != nil {
		conn.WriteError(err.Error())
		return
	}
	respInterface := future.Response()
	if respInterface == nil {
		conn.WriteError("RespInterface is nil, unknow error")
		return
	}
	response := respInterface.(raft.FsmResponse)
	WriteResponse(response, conn)
}

func WriteResponse(response raft.FsmResponse, conn redcon.Conn) {
	if response.Error != nil {
		conn.WriteError(response.Error.Error())
		return
	}
	switch response.Data.(type) {
	case int:
		conn.WriteInt(response.Data.(int))
	case int64:
		conn.WriteInt64(response.Data.(int64))
	case *string:
		if response.Data == nil {
			conn.WriteNull()
		} else {
			t := reflect.ValueOf(response.Data)
			if t.IsNil() {
				conn.WriteNull()
			} else {
				conn.WriteString(*response.Data.(*string))
			}
		}
	default:
		conn.WriteAny(response.Data)
	}
}
