package errors

// Traceable knows how return a runtime stack trace captured within an error
type Traceable interface {
	StackTrace()
}

// WithStack compose an error with a stack trace
func WithStack(err error) Traceable {
	return &stacked{
		err:   err,
		stack: nil,
	}
}

type stacked struct {
	err   error
	stack []string
}

func (s *stacked) StackTrace() {
	// TODO
}
