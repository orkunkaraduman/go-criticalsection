package criticalsection

import (
	"sync"
)

// A CriticalSection is a kind of lock like mutex. But it doesn't block
// first locked goroutine/section again.
//
// A CriticalSection must not be copied after first use.
type CriticalSection struct {
	mu sync.Mutex
	c  chan struct{}
	v  int32
	id uint64
	sc Section
}

// Lock locks cs.
// If the lock is already in use different goroutine, the different
// goroutine blocks until the CriticalSection is available.
func (cs *CriticalSection) Lock() {
	id := getGID()
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
	}
}

// Unlock unlocks cs.
// It panics if cs is not locked on entry to Unlock.
func (cs *CriticalSection) Unlock() {
	cs.mu.Lock()
	if cs.c == nil {
		cs.c = make(chan struct{}, 1)
	}
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
	select {
	case cs.c <- struct{}{}:
	default:
	}
}

// LockSection locks cs by given section. LockSection is faster than Lock().
// If the lock is already in use different section, the different
// section blocks until the CriticalSection is available.
func (cs *CriticalSection) LockSection(sc Section) {
	for {
		cs.mu.Lock()
		if cs.c == nil {
			cs.c = make(chan struct{}, 1)
		}
		if cs.v == 0 || cs.sc == sc {
			cs.v++
			cs.sc = sc
			cs.mu.Unlock()
			break
		}
		cs.mu.Unlock()
		<-cs.c
	}
}
