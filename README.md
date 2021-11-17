<!--
 * @Descripttion: 
 * @version: 
 * @Author: cm.d
 * @Date: 2021-11-12 09:51:16
 * @LastEditors: cm.d
 * @LastEditTime: 2021-11-17 18:08:52
-->

# AlfheimDB

A linearizability distributed database by raft and wisckey, which supports redis client.  

# Build

This project build by mage, you will need install it at first:

```` shell
go get github.com/magefile/mage
````

Execute "mage build" in project dir:

```` shell
mage build
````

# Test Case

The single node test:

```` shell
mage initDataDir
mage test single
````

You can also startup a test cluster, the magefile include a cluster test case:

```` shell
mage initDataDir
mage test id1
mage test id2
mage test id3
````

# Dependencies

+ Go 1.16  
+ mage
+ [raft](https://github.com/hashicorp/raft)
+ [badger](https://github.com/dgraph-io/badger)
+ [redcon](https://github.com/tidwall/redcon)

# Command Support

String

+ Set
+ Get
+ Incr
+ Del

# Benchmarks

## Macbook pro 13, 2020(M1)  

### Single node test case

Set/Incr and Get:

```` shell
./redis-benchmark -p 6379 -t set,get -n 1000000 -q  -c 512
SET: 114850.12 requests per second, p50=2.727 msec                    
GET: 161524.80 requests per second, p50=1.447 msec 
````

### Three node test case

Set/Incr and Get:

```` shell
./redis-benchmark -p 6379 -t set,get -n 500000 -q  -c 512
SET: 66952.33 requests per second, p50=6.279 msec                    
GET: 161917.09 requests per second, p50=1.447 msec 
````

## Centos 7 8c8g

### Single node test case

```` shell
./redis-benchmark -p 6379 -t set,get -n 500000 -q  -c 512
SET: 40038.44 requests per second, p50=10.431 msec
GET: 100684.66 requests per second, p50=2.423 msec 
````  


# References

Raft: ["Raft: In Search of an Understandable Consensus Algorithm"](https://raft.github.io/raft.pdf)  

Wisckey: ["WiscKey: Separating Keys from Values in SSD-conscious Storage"](https://www.usenix.org/system/files/conference/fast16/fast16-papers-lu.pdf)(FAST2016)

# Todo List

+ [ ] WASM script
+ [x] mage build
+ [ ] High performance WAL log
+ [ ] Set support
+ [ ] Map support
+ [ ] ZSet support
  