//+build mage

/*
 * @Descripttion: AlfheimDB's makefile by mage
 * @version:
 * @Author: cm.d
 * @Date: 2021-11-15 20:18:19
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-17 17:29:46
 */
package main

import (
	"log"
	"os"

	"github.com/magefile/mage/sh"
	"github.com/sirupsen/logrus"
)

func Build(system string) {
	env := make(map[string]string)
	env["CGO_ENABLED"] = "0"
	switch system {
	case "linux":
		env["GOARCH"] = "amd64"
		env["GOOS"] = "linux"
	case "macos":
		env["GOARCH"] = "amd64"
		env["GOOS"] = "darwin"
	default:
		logrus.Fatal("Unknow platform, ", system)
	}
	logrus.Info("Build AlfheimDB without cgo, platform: ", system)

	err := sh.RunWith(env, "go", "build")
	if err != nil {
		logrus.Fatal("Build AlfheimDB error, ", err)
	}
	logrus.Info("Build AlfheimDB ok")
}

func Package() {
	logrus.Info("Clean old target dir")
	sh.Run("rm", "-rf", "target")
	logrus.Info("Create new target dir")
	err := sh.Run("mkdir", "target")
	if err != nil {
		log.Fatal("Create new target dir error, ", err)
	}
	logrus.Info("Packaging...")
	err = sh.Run("cp", "AlfheimDB", "target/")
	if err != nil {
		logrus.Fatal("Copy AlftheimDB error, ", err)
	}
	err = sh.Run("cp", "config.yaml", "target/")
	if err != nil {
		logrus.Fatal("Copy config.yaml error, ", err)
	}
	err = sh.Run("mkdir", "target/data")
	if err != nil {
		logrus.Fatal("Mkdir raft data dir error, ", err)
	}
	err = sh.Run("mkdir", "target/runtime-log")
	if err != nil {
		logrus.Fatal("Mkdir log dir error, ", err)
	}
	err = sh.Run("chmod", "+x", "target/AlfheimDB")
	if err != nil {
		logrus.Fatal("Chmod AlfheimDB error, ", err)
	}
	logrus.Info("Package done")
}

func InitDataDir() {
	logrus.Info("Clean data dir")
	sh.Run("rm", "-rf", "data/")
	logrus.Info("Create new data dir")
	sh.Run("mkdir", "data/")
	logrus.Info("Create data dir done")
}

func testSingle() {
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
	logrus.Info("Start id1 test node, all 3 nodes in cluster")
	_, err := sh.Exec(nil, os.Stdout, os.Stderr, "./AlfheimDB", "--httpserver_addr=localhost:12345", "--raft_addr=localhost:40000", "--raft_id=id1", "--raft_cluster=localhost:40000/id1,localhost:40001/id2,localhost:40002/id3", "--respserver_addr=localhost:6379")
	if err != nil {
		logrus.Fatal("AlfheimDB start error, ", err)
	}
	logrus.Fatal("AlfheimDB start ok")
}

func testStartID2() {
	logrus.Info("Start id2 test node, all 3 nodes in cluster")
	_, err := sh.Exec(nil, os.Stdout, os.Stdin, "./AlfheimDB", "--httpserver_addr=localhost:12346", "--raft_addr=localhost:40001", "--raft_id=id2", "--raft_cluster=localhost:40000/id1,localhost:40001/id2,localhost:40002/id3", "--respserver_addr=localhost:6380")
	if err != nil {
		logrus.Fatal("AlfheimDB start error, ", err)
	}
	logrus.Fatal("AlfheimDB start ok")
}

func testStartID3() {
	logrus.Info("Start id3 test node, all 3 nodes in cluster")
	_, err := sh.Exec(nil, os.Stdout, os.Stdin, "./AlfheimDB", "--httpserver_addr=localhost:12347", "--raft_addr=localhost:40002", "--raft_id=id3", "--raft_cluster=localhost:40000/id1,localhost:40001/id2,localhost:40002/id3", "--respserver_addr=localhost:6381")
	if err != nil {
		logrus.Fatal("AlfheimDB start error, ", err)
	}
	logrus.Fatal("AlfheimDB start ok")
}
