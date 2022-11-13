package zbufferv2

import (
	"strings"
	"time"
	"unsafe"

	"github.com/RexLetRock/scriptcache/libs/zcount"
	"github.com/RexLetRock/zlib/zbench"
)

const cRun = 10_000_000
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
	log("\n\n==== ZBUFFER ===\n")
	log("WRITE ---msg---> BUFFER ---msg---> READER < " + cMsg + " >")
	logf("Buffer size: %T, %d\n", zbuffer, unsafe.Sizeof(*zbuffer))

	zbench.Run(cRun, cCpu, func(i, j int) {
		zbuffer.Write([]byte(cMsg + cSplit))
	})

	time.Sleep(5 * time.Second)
	logf("CountAll %v \n", countAll.Value())
}
