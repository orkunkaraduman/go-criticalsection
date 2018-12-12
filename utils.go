package criticalsection

import (
	"bytes"
	"runtime"
	"strconv"
)

func getGID() uint64 {
	var stack [64]byte
	b := stack[:runtime.Stack(stack[:], false)]
	b = bytes.TrimPrefix(b, []byte("goroutine "))
	b = b[:bytes.IndexByte(b, ' ')]
	n, e := strconv.ParseUint(string(b), 10, 64)
	if e != nil {
		panic(e)
	}
	return n
}
