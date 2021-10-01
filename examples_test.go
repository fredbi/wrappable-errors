package errors

import (
	"fmt"
	"io"
)

func ExampleWrappable_error_type() {
	// MyErrorType represents some abstract error type, e.g. to capture all
	// errors returned by some package.
	type MyErrorType struct {
		Wrappable
	}

	newMyErr := func(msg string) *MyErrorType {
		return &MyErrorType{Wrappable: New(msg)}
	}

	var (
		// ErrMyErr1 is assumed to be an immutable var, e.g. for a package
		ErrMyErr1 = newMyErr("err1")
	)

	// raise an error with some inner cause
	inner := io.EOF
	err := ErrMyErr1.Wrap(inner)

	fmt.Printf("err: %v\n", err)
	fmt.Printf("is io.EOF? %t\n", Is(err, inner))

	// Output:
	// err: err1: EOF
	// is io.EOF? true
}

func ExampleIs_sentinel_value() {
	inner := io.EOF

	var (
		ErrFirst  = New("err1")
		ErrSecond = NewErr(inner)
	)

	// raise an error with some inner cause
	err := ErrFirst.Wrap(ErrSecond)

	fmt.Printf("err: %v\n", err)
	fmt.Printf("is io.EOF? %t\n", Is(err, inner))

	// Output:
	// err: err1: EOF
	// is io.EOF? true
}

func ExampleIs_with_wrapped() {
	inner := io.EOF

	var (
		ErrFirst  = New("err1")
		ErrSecond = NewErr(inner)
	)

	// raise an error with some inner cause
	err := ErrFirst.Wrap(ErrSecond)

	fmt.Printf("err: %v\n", err)
	fmt.Printf("is ErrSecond? %t\n", Is(err, ErrSecond))

	err2 := ErrSecond.Wrap(ErrFirst)

	fmt.Printf("err2: %v\n", err2)
	fmt.Printf("is EOF? %t\n", Is(err, inner))

	// Output:
	// err: err1: EOF
	// is ErrSecond? true
	// err2: EOF: err1
	// is EOF? true
}

func ExampleWrappable_Err() {
	inner := io.EOF

	var (
		err1 = NewErr(inner).Errorf("inside")
		err2 = New("outside").Wrap(inner)
	)

	// display the wrapped error at the top level
	fmt.Printf("err1: %v\n", err1.Err())
	fmt.Printf("err2: %v\n", err2.Err())

	// Output:
	// err1: EOF
	// err2: outside
}

type testError string

func (e testError) Error() string {
	return string(e)
}

func ExampleAs_custom_type() {
	inner := io.EOF

	var (
		ErrFirst  = New("err1")
		ErrSecond = NewErr(inner)
		ErrThird  = testError("err3") // specific error type
	)

	// raise an error with some inner causes
	err := ErrFirst.Wrap(ErrSecond).Wrap(ErrThird)

	fmt.Printf("err: %v\n", err)

	var target testError
	isTestError := As(err, &target)
	fmt.Printf("as testError? %t, target: %v\n", isTestError, target)

	// Output:
	// err: err1: EOF: err3
	// as testError? true, target: err3
}

func ExampleRootable_wrapped() {
	inner := io.EOF

	var (
		ErrFirst = NewWithRoot("err1")
	)

	// raise an error with some inner cause
	err := ErrFirst.Wrap(io.EOF)

	fmt.Printf("err: %v\n", err)
	fmt.Printf("is io.EOF? %t\n", Is(err, inner))

	fmt.Printf("root cause: %v\n", err.(Rootable).Root())

	// Output:
	// err: err1: EOF
	// is io.EOF? true
	// root cause: EOF
}
