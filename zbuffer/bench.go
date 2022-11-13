package zbuffer

import (
	"fmt"
	"strings"
	"time"
	"unsafe"

	"github.com/RexLetRock/scriptcache/libs/zcount"
	"github.com/RexLetRock/zlib/zbench"
)

var countAll zcount.Counter

func Bench() {
	handle := func(data []byte) {
		a := strings.Split(string(data), "|||")
		countAll.Add(int64(len(a) - 1))
	}

	zbuffer := BufferCreate(handle)
	log("\n\n==== ZBUFFER ===\n")
	log("WRITE ---msg---> BUFFER ---msg---> READER < How Are You Today >")
	fmt.Printf("Buffer size: %T, %d\n", zbuffer, unsafe.Sizeof(*zbuffer))

	zbench.Run(5_000_000, 12, func(i, j int) {
		zbuffer.Write([]byte("How Are You Today|||"))
	})

	zbench.Run(50_000_000, 12, func(i, j int) {
		zbuffer.Write([]byte("How Are You Today|||"))
	})

	zbench.Run(100_000_000, 12, func(i, j int) {
		zbuffer.Write([]byte("How Are You Today|||"))
	})

	time.Sleep(time.Second)
	log("CountAll ", countAll.Value())
}
