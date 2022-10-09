package ztcp

import (
	"time"

	"github.com/RexLetRock/scriptcache/zevio"
	"github.com/RexLetRock/scriptcache/ztcp/ztcpclientv2"
	"github.com/sirupsen/logrus"
)

const Address = "127.0.0.1:9000"

func Bench() {
	logrus.Warnf("\n\n==== ZTCP ===\n")
	go zevio.MainEvio(Address)
	time.Sleep(2 * time.Second)
	ztcpclientv2.ClientStart(Address)
}