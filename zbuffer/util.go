package zbuffer

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/RexLetRock/scriptcache/libs/zcount"
	"github.com/RexLetRock/zlib/zbench"
	"github.com/sirupsen/logrus"
)

var countAll zcount.Counter

func Bench() {
	zbuffer := BufferCreate()
	logrus.Warnf("\n\n==== ZBUFFER ===\n")
	fmt.Printf("Buffer size: %T, %d\n", zbuffer, unsafe.Sizeof(*zbuffer))

	zbench.Run(5_000_000, 12, func(i, j int) {
		zbuffer.Write([]byte("How Are You|||"))
	})

	zbench.Run(50_000_000, 12, func(i, j int) {
		zbuffer.Write([]byte("How Are You|||"))
	})

	time.Sleep(time.Second)
	log("Hello countall ", countAll.Value())
}

func log(args ...interface{}) {
	logrus.Warn(args...)
}

func logf(args ...interface{}) {
	fmt.Print(args...)
}
