// Package errors provide a wrappable error type.
//
// It is intended to use with sentinel errors or type assertions on errors,
// e.g. using errors.Is().
//
// It generalizes the concept of wrapping errors already present with fmt.Errorf("... %w", err),
// but allows to wrap typed errors with a Wrap(err error) method.
//
// As its simplest, this package may be used to derive error values or types and proceed with type or value assertion on
// sentinel errors using Wrap() and Is() or As().
//
// Runtime stack trace capture is provided as an optional addon (using WithStackTrace()).
//
// To capture the root cause of an error stack (i.e. the deepest error in the stack), one can use the Root() method.
package errors
