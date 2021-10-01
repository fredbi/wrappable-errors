# wrappable-errors

Wrappable errors allow programs to handle errors as stacks.

The main difference with the standard libray error is the new method `Wrap(error)`
which makes it easier than `fmt.Errorf("...%w")` to reason with static error values
and error types.

I wanted to improve errors in order for a package to expose a _class_ of errors, with
some static values, and remove the need to check for errors by their string representation.

Since `go1.13` the standard library supports extended errors, which know how to `Unwrap() error`,
`Is(error) bool` and `As(interface{}) bool`.

We build on top of this major improvement with a compatible error type.
