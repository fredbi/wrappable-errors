package errors_test

import (
	"io"
	"testing"

	errors "github.com/fredbi/wrappable-errors"
	"github.com/stretchr/testify/assert"
)

type (
	// MyErrorType represents a class of errors returned by the current package.
	MyErrorType struct {
		errors.Wrappable
	}

	OtherErrors struct {
		error
	}
)

// Wrap overrides the underlying Wrappable to maintain a strongly type wrapping
func (e MyErrorType) Wrap(err error) *MyErrorType {
	return &MyErrorType{Wrappable: e.Wrappable.Wrap(err)}
}

// NewMyErrorType builds a new MyErrorType
func NewMyErrorType(msg string) *MyErrorType {
	return &MyErrorType{Wrappable: errors.New(msg)}
}

var (
	// ErrPkg1 represents some package level error. Its value won't change
	ErrPkg1 = NewMyErrorType("err1")

	// ErrPkg2 represents some other package level error. Its value won't change
	ErrPkg2 = NewMyErrorType("err2")

	// ErrOther represents another class of errors. This one is not wrappable.
	ErrOther = OtherErrors{error: errors.New("other")}
)

// TestCustomType illustrates how to use typed errors
func TestCustomType(t *testing.T) {
	e := &MyErrorType{} // template value

	inner := io.EOF

	err1 := ErrPkg1.Wrap(inner) // err1 <- EOF
	err2 := err1.Wrap(ErrPkg2)  // err1 <- EOF <- err2

	// err1 and err2 are both of class MyErrorType
	assert.IsType(t, e, ErrPkg1)
	assert.IsType(t, e, ErrPkg2)

	t.Logf("err1.Error(): %v\n", err1)
	t.Logf("err2.Error(): %v\n", err2)

	// errors.Is() works to identify the inner EOF in the stack
	assert.Truef(t, errors.Is(err2, inner), "err2: %v, inner: %v", err2, inner)
	assert.Truef(t, errors.Is(err2, ErrPkg2), "err2: %v, ErrPkg2: %v", err2, ErrPkg2)

	// different error trees / TODO(fred)

	// errors.As() works to extract a typed error
	err3 := err2.Wrap(ErrOther) // err1 <- EOF <- err2 <- other
	t.Logf("err3.Error(): %v\n", err3)
	var target OtherErrors
	assert.Truef(t, errors.As(err3, &target), "err3: %v, target: %v", err3, target)
}
