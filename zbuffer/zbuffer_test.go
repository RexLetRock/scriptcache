package zbuffer

import (
	"strings"
	"testing"
)

var testData = []byte("How are you today|||")

func BenchmarkWrite(b *testing.B) {
	handle := func(data []byte) {
		a := strings.Split(string(data), cSplit)
		countAll.Add(int64(len(a) - 1))
	}

	zbuffer := ZBufferCreate(handle)
	for i := 0; i < b.N; i++ {
		zbuffer.Write(testData)
	}
}
