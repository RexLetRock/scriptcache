package main

import (
	"scriptcache/tcp"
)

const connStr = "developer:password@tcp(127.0.0.1:4000)/imsystem?parseTime=true"
const connHost = "0.0.0.0:8888"

func main() {
	tcp.ServerStartViaOptions(connStr, connHost)
	select {}
}
