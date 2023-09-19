package errors

import "errors"

var _ Rootable = &wrapped{}

// Root cause of the error: returns the deepest wrapped error in the chain
func Root(err error) error {
	if rootable, ok := err.(Rootable); ok {
		return rootable.Root()
	}

	last := err
	next := errors.Unwrap(err)

	for next != nil {
		if rootable, ok := next.(Rootable); ok {
			return rootable.Root()
		}

		last = next

		if unwrapped, ok := next.(interface{ Unwrap() error }); ok {
			next = unwrapped.Unwrap()
		} else {
			next = nil
		}
	}

	return last
}

// NewWithRoot wrappable & rootable error from a string
func NewWithRoot(msg string) Rootable {
	return &wrapped{err: errors.New(msg)}
}

// NewErrWithRoot wrappable & rootable error from another error
func NewErrWithRoot(err error) Rootable {
	return &wrapped{err: err}
}

// Root returns the root cause of a wrapped error
func (e wrapped) Root() error {
	last := e.err
	next := e.Unwrap()

	for next != nil {
		if errable, ok := next.(interface{ Err() error }); ok {
			last = errable.Err()
		} else {
			last = next
		}

		if unwrapped, ok := next.(interface{ Unwrap() error }); ok {
			next = unwrapped.Unwrap()
		} else {
			next = nil
		}
	}

	return last
}
