package main

import (
	"time"

	"github.com/RexLetRock/scriptcache/zbufferv3"
)

const cWait = 1 * time.Second

func main() {
	// ztcp.Bench()
	// zbufferv2.Bench()
	zbufferv3.Bench()
	time.Sleep(cWait)
}
