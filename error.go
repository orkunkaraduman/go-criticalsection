package criticalsection

// Library specific errors.
var (
	ErrPanic     = NewError("panic")
	ErrNotLocked = NewError("not locked")
)

// Error struct is base of library specific errors.
type Error struct {
	s string
}

// NewError returns a new Error.
func NewError(text string) *Error {
	return &Error{text}
}

// Error implements error interface.
func (e *Error) Error() string {
	return e.s
}
