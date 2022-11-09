package ztcp

import (
	"github.com/RexLetRock/scriptcache/zlibs/zbuffer"
	"github.com/sirupsen/logrus"
)

const Address = "127.0.0.1:9000"

func Bench() {
	logrus.Warnf("\n\n==== ZTCP ===\n")

	zbuffer.Bench()
	// go ztcpserver.ServerStart(Address)
	// time.Sleep(time.Second)
	// ztcpclient.ClientStart(Address)
}
