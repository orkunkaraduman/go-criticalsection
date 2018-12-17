package criticalsection

import (
	"sync/atomic"
)

// A Section is an identifier of sections.
type Section uint64

var sec = uint64(0)

// NewSection creates a new Section.
func NewSection() Section {
	return Section(atomic.AddUint64(&sec, 1))
}
