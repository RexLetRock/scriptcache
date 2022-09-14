package main

import (
	"github.com/RexLetRock/scriptcache/snowflake"
	"github.com/RexLetRock/zlib/zbench"
	"github.com/sirupsen/logrus"
)

const Address = "127.0.0.1:9000"
const NCpu = 1

func main() {
	// go zgnet.MainGnet()

	// go zevio.MainEvio(Address)
	// go ztcpserver.ServerStartViaOptions(Address)

	// time.Sleep(2 * time.Second)
	// ztcpclient.ClientStart(Address)
	snow, _ := snowflake.NewSnowflake(1)
	id := uint64(6671919771961558991)
	for i := 1; i < 100; i++ {
		id = snow.NextIdWithSeq(id)
		logrus.Warnf("Nextid %v %v \n", id, id&0xFFF)
	}

	var snows [NCpu]*snowflake.Snowflake
	for i := 0; i < NCpu; i++ {
		snows[i], _ = snowflake.NewSnowflake(uint64(i))
	}

	zbench.Run(10_000, NCpu, func(i, thread int) {
		id = snows[thread].NextIdWithSeqViaCache(id)
	})
}
