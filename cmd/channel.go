package cmd

import (
	"bufio"
	"os"
	"time"
)

const NChannel = 12
const NCpuname = 10000

var (
	NRun   = 10_000_000
	NCpu   = 12
	ACount [NCpuname]int

	timeStart = time.Now().Unix()
	timeNow   = int64(0)
	stotal    = 0
)

func benchChannel(threadName int) {
	// Prebuffer channel
	c := make(chan string, 1000)

	// Producer
	go func() {
		for {
			c <- "How is the weather like today ? hope you okie"
		}
	}()

	// Consumer
	for r := range c {
		_ = r
		ACount[threadName] += 1
	}
}

func MultiChannel() {
	// Bench
	for i := 0; i <= NCpu; i++ {
		go benchChannel(i)
	}

	// Result
	go showResult()
	bufio.NewScanner(os.Stdin).Scan()
}

// type Conns struct {
// 	c chan (string)
// 	conn
// }

// type ScriptCacheClient struct {
// }
