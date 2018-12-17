package criticalsection

import (
	"runtime"
	"sync"
	"time"
)

const maxSleepCycle = 6

// A CriticalSection is a kind of lock like mutex. But it doesn't block
// first locked section again.
//
// A CriticalSection must not be copied after first use.
type CriticalSection struct {
	mu sync.Mutex
	v  int32
	id uint64
	sc Section
}

// Lock locks cs.
// If the lock is already in use different goroutine, the different
// goroutine blocks until the CriticalSection is available.
func (cs *CriticalSection) Lock() {
	id := getGID()
	i := uint(0)
	for {
		cs.mu.Lock()
		if cs.v == 0 || cs.id == id {
			cs.v++
			cs.id = id
			cs.mu.Unlock()
			break
		}
		cs.mu.Unlock()
		runtime.Gosched()
		time.Sleep(time.Duration(1<<i) * 1 * time.Millisecond)
		i++
		if i >= maxSleepCycle {
			i = 0
		}
	}
}

// Unlock unlocks cs.
// It panics if cs is not locked on entry to Unlock.
func (cs *CriticalSection) Unlock() {
	cs.mu.Lock()
	if cs.v < 0 {
		cs.mu.Unlock()
		panic(ErrNotLocked)
	}
	cs.v--
	if cs.v == 0 {
		cs.id = 0
		cs.sc = 0
	}
	cs.mu.Unlock()
}

// LockSection locks cs by given section.
// If the lock is already in use different section, the different
// section blocks until the CriticalSection is available.
func (cs *CriticalSection) LockSection(sc Section) {
	i := uint(0)
	for {
		cs.mu.Lock()
		if cs.v == 0 || cs.sc == sc {
			cs.v++
			cs.sc = sc
			cs.mu.Unlock()
			break
		}
		cs.mu.Unlock()
		runtime.Gosched()
		time.Sleep(time.Duration(1<<i) * 1 * time.Millisecond)
		i++
		if i >= maxSleepCycle {
			i = 0
		}
	}
}
