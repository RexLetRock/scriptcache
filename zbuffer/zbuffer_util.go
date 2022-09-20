package zbuffer

import (
	"fmt"
	"unsafe"

	"github.com/RexLetRock/zlib/zbench"
	"github.com/RexLetRock/zlib/zgoid"
	"github.com/sirupsen/logrus"
)

func (s *ZBuffer) getGID() (id uint16, gid uint16) {
	gid = uint16(zgoid.Get())
	if gid >= CMaxCpu {
		id = gid % CMaxCpu
	} else {
		id = gid
	}
	return
}

func Bench() {
	zbuffer := ZBufferCreate()
	logrus.Warnf("==== ZBUFFER ===\n")
	fmt.Printf("ZBuffer size: %T, %d\n", zbuffer, unsafe.Sizeof(*zbuffer))

	zbench.Run(50_000_000, 12, func(i, thread int) {
		zbuffer.Write([]byte("Hello How Are You Today|||"))
	})
}
