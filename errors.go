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

// Unwrap nested error.
//
// This method is only provided for this package to nicely supersede standard lib errors:
// it just calls Unwrap() from the standard library.
func Unwrap(err error) error {
	return errors.Unwrap(err)
}
