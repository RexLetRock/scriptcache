package ztcp

import (
	"time"

	"github.com/RexLetRock/scriptcache/ztcp/ztcpclient"
	"github.com/RexLetRock/scriptcache/ztcp/ztcpserver"
	"github.com/sirupsen/logrus"
)

const Address = "127.0.0.1:9000"

func Bench() {
	logrus.Warnf("\n\n==== ZTCP ===\n")
	// go zevio.MainEvio(Address)
	go ztcpserver.ServerStart(Address)
	time.Sleep(time.Second)
	ztcpclient.ClientStart(Address)
}
