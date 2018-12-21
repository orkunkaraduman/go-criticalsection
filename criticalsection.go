package criticalsection

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

// A CriticalSection is a kind of lock like mutex. But it doesn't block
// first locked goroutine/section again.
//
// A CriticalSection must not be copied after first use.
type CriticalSection struct {
	mu sync.Mutex
	c  *chan struct{}
	v  int32
	id uint64
	sc uint64
}

// Lock locks cs.
// If the lock is already in use different goroutine, the different
// goroutine blocks until the CriticalSection is available.
func (cs *CriticalSection) Lock() {
	/*id := getGID()
	for {
		cs.mu.Lock()
		if cs.c == nil {
			cs.c = make(chan struct{}, 1)
		}
		if cs.v == 0 || cs.id == id {
			cs.v++
			cs.id = id
			cs.mu.Unlock()
			break
		}
		cs.mu.Unlock()
		<-cs.c
	}*/
}

// Unlock unlocks cs.
// It panics if cs is not locked on entry to Unlock.
func (cs *CriticalSection) Unlock() {
	if atomic.AddInt32(&cs.v, -1) < 0 {
		atomic.AddInt32(&cs.v, 1)
		panic(ErrNotLocked)
	}
	select {
	case *cs.c <- struct{}{}:
	default:
	}
}

// LockSection locks cs by given section. LockSection is faster than Lock().
// If the lock is already in use different section, the different
// section blocks until the CriticalSection is available.
func (cs *CriticalSection) LockSection(sc Section) {
	for {
		c := make(chan struct{}, 1)
		if !atomic.CompareAndSwapPointer(
			(*unsafe.Pointer)(unsafe.Pointer(&cs.c)),
			nil,
			(unsafe.Pointer)(unsafe.Pointer(&c))) {
			close(c)
			c = nil
		}
		if atomic.AddInt32(&cs.v, 1) == 1 {
			cs.id = 0
			cs.sc = uint64(sc)
			break
		}
		if atomic.CompareAndSwapUint64(&cs.sc, uint64(sc), uint64(sc)) {
			break
		}
		atomic.AddInt32(&cs.v, -1)
		<-*cs.c
	}
}
