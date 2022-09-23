package zbuff

import (
	"fmt"
	"time"
	"unsafe"

	"github.com/RexLetRock/zlib/zbench"
	"github.com/RexLetRock/zlib/zgoid"
	"github.com/sirupsen/logrus"
)

const CMaxChan = 100000
const CMaxCpu = 1000

type Buff struct {
	Chan [CMaxCpu]chan []byte
}

func BuffCreate() *Buff {
	s := &Buff{}
	return s
}

func ChanCreateMulti() {
	var Chan [CMaxCpu]chan []byte
	for i := 0; i < CMaxCpu; i++ {
		Chan[i] = make(chan []byte, CMaxChan)
	}
}

func (s *Buff) Write(data []byte) {
	id, _ := getGID()
	select {
	case s.Chan[id] <- data:
	default:
	}
}

func getGID() (id uint16, gid uint16) {
	gid = uint16(zgoid.Get())
	if gid >= CMaxCpu {
		id = gid % CMaxCpu
	} else {
		id = gid
	}
	return
}

func Bench() {
	zbuffer := BuffCreate()
	logrus.Warnf("==== ZBUFFER ===\n")
	fmt.Printf("Buffer size: %T, %d\n", zbuffer, unsafe.Sizeof(*zbuffer))

	zbench.Run(100_000_000, 12, func(i, thread int) {
		zbuffer.Write([]byte("Hello How Are You Today|||"))
	})

	time.Sleep(5 * time.Second)
}
