package tcp

import (
	"bufio"
	"io"
	"io/ioutil"
	"net"

	"github.com/RexLetRock/zlib/zbench"
)

const NCpu = 30
const NRun = 3_000_000

var conns = [NCpu]net.Conn{}

func ClientStart() {
	for i := 0; i < NCpu; i++ {
		conns[i], _ = net.Dial("tcp", "127.0.0.1:8888")
	}

	a := []byte("How are you today :D \n")

	pipeReader, pipeWriter := io.Pipe()
	reader := bufio.NewReader(pipeReader)
	go func() {
		for {
			data, _ := ioutil.ReadAll(reader)
			conns[0].Write(data)
		}
	}()

	zbench.Run(NRun, NCpu, func(i, thread int) {
		// conns[thread].Write(a)
		pipeWriter.Write(a)
	})
}
