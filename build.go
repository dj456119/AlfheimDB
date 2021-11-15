//+build mage

/*
 * @Descripttion: AlfheimDB's makefile by mage
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-15 20:18:19
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-15 22:11:22
 */
package main

import (
	"os"

	"github.com/magefile/mage/sh"
	"github.com/sirupsen/logrus"
)

func Build() {
	logrus.Info("Build AlfheimDB without cgo")
	env := make(map[string]string)
	env["CGO_ENABLED"] = "0"
	err := sh.Run("go", "build")
	if err != nil {
		logrus.Fatal("Build AlfheimDB error, ", err)
	}
	logrus.Info("Build AlfheimDB ok")
}

func InitDataDir() {
	logrus.Info("Clean data dir")
	sh.Run("rm", "-rf", "data/")
	logrus.Info("Create new data dir")
	sh.Run("mkdir", "data/")
	logrus.Info("Create data dir done")
}

func testSingle() {
	createIDDir("id1")
	logrus.Info("Start id1 test single")
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "./AlfheimDB", "--httpserver_addr=localhost:12345", "--raft_addr=localhost:40000", "--raft_id=id1", "--raft_cluster=localhost:40000/id1", "--respserver_addr=localhost:6379")
	if err != nil {
		logrus.Fatal("AlfheimDB start error, ", err)
	}
	logrus.Fatal("AlfheimDB start ok")
}

func Test(testType string) {
	switch testType {
	case "single":
		testSingle()
	case "id1":
		testStartID1()
	case "id2":
		testStartID2()
	case "id3":
		testStartID3()
	default:
		logrus.Fatal("Test error, unknow test arg ", testType)
	}
}

func testStartID1() {
	createIDDir("id1")
	logrus.Info("Start id1 test node, all 3 nodes in cluster")
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "./AlfheimDB", "--httpserver_addr=localhost:12345", "--raft_addr=localhost:40000", "--raft_id=id1", "--raft_cluster=localhost:40000/id1,localhost:40001/id2,localhost:40002/id3", "--respserver_addr=localhost:6379")
	if err != nil {
		logrus.Fatal("AlfheimDB start error, ", err)
	}
	logrus.Fatal("AlfheimDB start ok")
}

func createIDDir(id string) {
	logrus.Info("Create ", id, " dir")
	err := sh.Run("mkdir", "data/"+id)
	if err != nil {
		logrus.Warn("Create ", id, " dir error, ", err)
		return
	}
	logrus.Info("Create ", id, " dir ok")
}

func testStartID2() {
	createIDDir("id2")
	logrus.Info("Start id2 test node, all 3 nodes in cluster")
	_, err := sh.Exec(nil, os.Stdout, os.Stdin, "./AlfheimDB", "--httpserver_addr=localhost:12346", "--raft_addr=localhost:40001", "--raft_id=id2", "--raft_cluster=localhost:40000/id1,localhost:40001/id2,localhost:40002/id3", "--respserver_addr=localhost:6380")
	if err != nil {
		logrus.Fatal("AlfheimDB start error, ", err)
	}
	logrus.Fatal("AlfheimDB start ok")
}

func testStartID3() {
	createIDDir("id3")
	logrus.Info("Start id3 test node, all 3 nodes in cluster")
	_, err := sh.Exec(nil, os.Stdout, os.Stdin, "./AlfheimDB", "--httpserver_addr=localhost:12347", "--raft_addr=localhost:40002", "--raft_id=id3", "--raft_cluster=localhost:40000/id1,localhost:40001/id2,localhost:40002/id3", "--respserver_addr=localhost:6381")
	if err != nil {
		logrus.Fatal("AlfheimDB start error, ", err)
	}
	logrus.Fatal("AlfheimDB start ok")
}
