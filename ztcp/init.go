package ztcp

import (
	"strconv"
	"time"

	"github.com/RexLetRock/scriptcache/ztcp/ztcpclient"
	"github.com/RexLetRock/scriptcache/ztcp/ztcpserver"
	"github.com/RexLetRock/zlib/zbench"
	"github.com/sirupsen/logrus"

	cmap "github.com/orcaman/concurrent-map/v2"
)

const Address = "127.0.0.1:9000"

func Bench() {
	logrus.Warnf("==== ZTCP ===\n")
	go ztcpserver.ServerStart(Address)
	time.Sleep(time.Second)
	ztcpclient.ClientStart(Address)
	time.Sleep(10 * time.Second)
}

func BenchMap() {
	logrus.Warnf("==== ZMAP ====")
	cmap := cmap.New[int]()
	cdata := [30_000_000]string{}
	zbench.Run(30_000_000, 12, func(i, thread int) {
		cdata[i] = strconv.Itoa(i)
	})
	zbench.Run(5_000_000, 12, func(i, thread int) {
		cmap.Set(cdata[i], i)
	})
}
