package errors

import (
	"errors"
)

/* TODOs(fred)
I'd like to:
- optionally add a stack trace
- nice json unmarshalling
- Formatter
*/

// Wrappable is a wrappable error.
//
// The particularity of this error type is that is adds a Wrap(error) Wrappable method to add new errors to a stack of errors.
type Wrappable interface {
	error

	// Wrap an error into another one
	Wrap(error) Wrappable

	// Unwrap the inner error, like standard lib errors.Unwrap()
	Unwrap() error

	// Errorf is like fmt.Errorf, but wraps the newly created error into the current one
	Errorf(string, ...interface{}) Wrappable

	// Err returns the topmost error in the stack (head)
	Err() error
}

// Rootable is a Wrappable that knows how to yield its head and tail
type Rootable interface {
	Wrappable

	// Root returns the deepest error in the case (tail, i.e. "root cause")
	Root() error
}

// We expose here the same interface as the stdlib errors package.
//
// This is mostly to avoid importing both and managing package aliases: wrapping
// standard lib calls allow for a single import ( "github.com/.../errors") clause.

// New wrappable error from a string
func New(msg string) Wrappable {
	return &wrapped{err: errors.New(msg)}
}

// NewErr wrappable error from another error
func NewErr(err error) Wrappable {
	return &wrapped{err: err}
}

// NewWithRoot wrappable & rootable error from a string
func NewWithRoot(msg string) Rootable {
	return &wrapped{err: errors.New(msg)}
}

// NewErrWithRoot wrappable & rootable error from another error
func NewErrWithRoot(err error) Rootable {
	return &wrapped{err: err}
}

// Is behaves like errors.Is from the standard library
//
// This method is only provided for this package to nicely supersede standard lib errors:
// it just calls Is() from the standard library.
func Is(err, target error) bool {
	return errors.Is(err, target)
}

// As behaves like errors.As from the standard library.
//
// This method is only provided for this package to nicely supersede standard lib errors:
// it just calls As() from the standard library.
func As(err error, target interface{}) bool {
	return errors.As(err, target)
	/*
		// standard lib case: types are comparable
		if errors.As(err, target) {
			return true
		}

			if target == nil {
				panic("errors: target cannot be nil")
			}

		// wrapped case: check topmost (might itself be a stack)
		if wrap, ok := err.(interface{ Err() error }); ok && As(wrap.Err(), target) {
			return true
		}

		// wrapped case: check inner error stack
		if unwrap, ok := err.(interface {
			error
			Unwrap() error
		}); ok {
			return As(unwrap.Unwrap(), target)
		}

		return false
	*/
}

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

// Unwrap nested error.
//
// This method is only provided for this package to nicely supersede standard lib errors:
// it just calls Unwrap() from the standard library.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
