# wrappable-errors

Wrappable errors allow programs to handle nested errors.

It is largely inspired from `github.com/pkg/errors`, and may be used as a standin for the standard library's `errors`

The main difference with the standard libray is the new method `Wrap(error)`
which makes it easier than `fmt.Errorf("...%w")` to reason with static error values
and error types.

I wanted to improve errors in order for a package to expose a _class_ of errors, with
some static values, and remove the need to check for errors by their string representation.

Since `go1.13` the standard library supports extended errors, which know how to `Unwrap() error`,
`Is(error) bool` and `As(interface{}) bool`.

We build on top of this major improvement with a compatible error type.

## Usage

### Sentinel errors
```go
import (
    "fmt"
    "io"
    "github.com/fredbi/wrappable-errors"
)

var (
	// ErrMyErr1 is assumed to be an immutable var, e.g. for a package
	ErrMyErr1 = errors.New("err1")
)

func meetError() error {
    return ErrMyErr1.Wrap(io.EOF)
}

func main() {
    err := meetError()

    fmt.Printf("is ErrMyErr1: %t", Is(err, ErrMyErr1))
}
```

### Custom error classes

#### Defining a simple class of errors
```go
import (
    "github.com/fredbi/wrappable-errors"
)

// MyErrorType represents some abstract error type, e.g. to capture all
// errors returned by some package.
type MyErrorType struct {
	Wrappable
}

newMyErr := func(msg string) *MyErrorType {
	return &MyErrorType{Wrappable: errors.New(msg)}
}

var (
	// ErrMyErr1 is assumed to be an immutable var, e.g. for a package
	ErrMyErr1 = newMyErr("err1")
)
```

#### Enriched errors
```go
import (
    "github.com/fredbi/wrappable-errors"
)

// MyErrorType represents some abstract error type, e.g. to capture all
// errors returned by some package.
type MyErrorType struct {
	Wrappable

    Code int `json:"code"`
    Message string `json:"message"`
}

newMyErr := func(msg string) *MyErrorType {
	return &MyErrorType{Wrappable: New(msg)}
}

var (
	// ErrMyErr1 is assumed to be an immutable var, e.g. for a package
	ErrMyErr1 = newMyErr("err1")
)
```
