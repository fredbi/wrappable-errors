package errors

var _ Rootable = &wrapped{}

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
