package criticalsection

import "errors"

// Library specific errors.
var (
	ErrNotLocked = errors.New("not locked")
)
