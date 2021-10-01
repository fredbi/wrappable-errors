package errors

import (
	"errors"
	"fmt"
	"reflect"
)

var _ Wrappable = &wrapped{}

// wrapped produces a stack of errors. It implements the Wrappable interface.
//
//
// This allows to check typed errors and Wrap them, which still remains impractical as of go1.13+ and is
// not the goal of github.com/pkg/errors either.
//
// wrapped is assumed to remain immutable and all methods produce shallow clones of the error.
type wrapped struct {
	err   error
	cause error
}

// Error implements the error interface, with plain formatting: all nested errors are printed, separated by a ":".
func (e wrapped) Error() string {
	if e.cause == nil {
		return e.err.Error()
	}

	return fmt.Sprintf("%v: %v", e.err, e.cause)
}

// Errorf wraps a nested error built from the extra message.
//
// This is a shorthand for Wrap(fmt.Errorf(format, args...)).
func (e wrapped) Errorf(format string, args ...interface{}) Wrappable {
	return e.Wrap(fmt.Errorf(format, args...))
}

// Wrap another error. Returns a shallow clone.
//
// Notice that you may wrap inner errors using fmt.Errorf() and it will still work:
// e.g. WithMessage("err: %w", err) will recognize the usual
// errors.Is() and errors.As() methods from the standard library.
//
// More generally error stacking supports any other stacking mechanism on underlying errors
// equipped with the standard Unwrap() error method.
func (e *wrapped) Wrap(err error) Wrappable {
	if err == nil {
		return e
	}

	if e.cause == nil {
		return &wrapped{
			err:   e.err,
			cause: err,
		}
	}

	wrapper, ok := e.cause.(Wrappable)
	if ok {
		// stack err at the tail of the cause
		return &wrapped{
			err:   e.err,
			cause: wrapper.Wrap(err),
		}
	}

	unwrapper, ok := e.cause.(interface{ Unwrap() error })
	if ok {
		// destructure the cause and wrap err at the tail
		cause := &wrapped{
			err:   e.cause,
			cause: unwrapper.Unwrap(),
		}
		return &wrapped{
			err:   e.err,
			cause: cause.Wrap(err),
		}
	}

	return &wrapped{
		err: e.err,
		cause: &wrapped{
			err:   e.cause,
			cause: err,
		},
	}
}

// Unwrap implements errors.Unwrap
func (e wrapped) Unwrap() error {
	if e.cause != nil {
		return e.cause
	}

	return e.err
}

func (e wrapped) Err() error {
	return e.err
}

// Is implements errors.Is
func (e *wrapped) Is(err error) bool {
	if e == err {
		return true
	}

	if err == nil {
		return false
	}

	if errors.Is(e.err, err) {
		return true
	}

	// special case for another wrapped error
	errable, ok := err.(*wrapped)
	if ok && errors.Is(e.err, errable.err) {
		return true
	}

	return errors.Is(e.cause, err)
}

var errorType = reflect.TypeOf((*error)(nil)).Elem()

// As implements errors.As
func (e *wrapped) As(target interface{}) bool {
	fmt.Printf("e: %T => As(%T)\n", e, target)
	fmt.Printf("e.err: %T => As(%T)\n", e.err, target)
	fmt.Printf("e.cause: %T => As(%T)\n", e.cause, target)

	u := reflect.TypeOf(e)
	v := reflect.TypeOf(target)
	val := reflect.ValueOf(target)
	if v.Kind() != reflect.Ptr || val.IsNil() {
		panic("wrappable-errors: target must be a non-nil pointer")
	}

	targetType := v.Elem()
	if targetType.Kind() != reflect.Interface && !targetType.Implements(errorType) {
		panic("wrappable-errors: *target must be interface or implement error")
	}

	if u.AssignableTo(targetType) {
		val.Elem().Set(reflect.ValueOf(e))

		return true
	}

	if e.err != nil {
		if errors.As(e.err, target) {
			return true
		}
	}

	if e.cause != nil {
		return errors.As(e.cause, target)
	}

	return false
}
