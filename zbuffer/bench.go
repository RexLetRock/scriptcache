package zbuffer

import (
	"strings"
	"time"
	"unsafe"

	"github.com/RexLetRock/scriptcache/libs/zcount"
	"github.com/RexLetRock/zlib/zbench"
)

const cRun = 100_000_000
const cCpu = 12
const cMsg = "How are you today ?"
const cSplit = "|||"

var countAll zcount.Counter

func Bench() {
	handle := func(data []byte) {
		a := strings.Split(string(data), cSplit)
		countAll.Add(int64(len(a) - 1))
	}

	zbuffer := ZBufferCreate(handle)
	warn("==== ZBUFFER ===\n")
	warn("WRITE ---msg---> BUFFER ---msg---> READER < " + cMsg + " >")
	warnf("Buffer size: %T, %d\n", zbuffer, unsafe.Sizeof(*zbuffer))

	zbench.Run(cRun, cCpu, func(i, j int) {
		zbuffer.Write([]byte(cMsg + cSplit))
	})

	zbench.Run(cRun, cCpu, func(i, j int) {
		zbuffer.Write([]byte(cMsg + cSplit))
	})

	zbench.Run(cRun, cCpu, func(i, j int) {
		zbuffer.Write([]byte(cMsg + cSplit))
	})

	zbench.Run(cRun, cCpu, func(i, j int) {
		zbuffer.Write([]byte(cMsg + cSplit))
	})

	zbench.Run(cRun, cCpu, func(i, j int) {
		zbuffer.Write([]byte(cMsg + cSplit))
	})

	time.Sleep(time.Second)
	warnf("CountAll %v \n", Commaize(int64(countAll.Value())))
}

func ExampleSimple() {
	handle := func(data []byte) {
		a := strings.Split(string(data), cSplit)
		countAll.Add(int64(len(a) - 1))
	}
	zbuffer := ZBufferCreate(handle)
	zbuffer.Write([]byte("How are you today ?|||"))
	time.Sleep(2 * time.Second)
	warnf("CountAll %v \n", Commaize(int64(countAll.Value())))
}
