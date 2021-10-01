package errors

import (
	"fmt"
	"io"
	"testing"

	stderrors "errors"

	"github.com/stretchr/testify/assert"
)

const str = "test error"

func TestWrap(t *testing.T) {
	t.Parallel()

	e := New(str)
	assert.EqualValues(t, e, e.Wrap(nil))

	assert.True(t, Is(e, e))

	assert.False(t, stderrors.Is(e, nil))
	assert.True(t, stderrors.Is(nil, nil))

	assert.False(t, Is(e, nil))
	assert.True(t, Is(nil, nil))

	withIs := e.(interface {
		Is(error) bool
		Error() string
	})

	assert.True(t, withIs.Is(e))
	assert.True(t, withIs.Is(withIs))
	assert.False(t, withIs.Is(nil))

	w1 := e.Wrap(io.ErrUnexpectedEOF)

	assert.Contains(t, w1.Error(), io.ErrUnexpectedEOF.Error())
	assert.Contains(t, w1.Error(), str)

	assert.True(t, Is(w1, io.ErrUnexpectedEOF))

	w2 := w1.Wrap(io.EOF)

	assert.Contains(t, w1.Error(), io.ErrUnexpectedEOF.Error())
	assert.Contains(t, w1.Error(), io.EOF.Error())
	assert.Contains(t, w1.Error(), str)

	u1 := Unwrap(w1)

	assert.NotContains(t, u1.Error(), str)
	assert.Contains(t, w1.Error(), io.ErrUnexpectedEOF.Error())

	assert.True(t, Is(w2, io.EOF))
	assert.True(t, Is(w2, io.ErrUnexpectedEOF))
	assert.False(t, Is(w2, io.ErrClosedPipe))

	var tg Wrappable
	_ = As(w2, &tg)

	assert.ErrorIs(t, tg, io.ErrUnexpectedEOF)

	u2 := Unwrap(w2)
	u3 := Unwrap(u2)

	assert.True(t, Is(u3, io.EOF))
	assert.False(t, Is(u3, io.ErrUnexpectedEOF))

	assert.True(t, As(w2, &tg))
	assert.EqualValues(t, w2, tg)

	assert.True(t, Is(w2, w2))

	assert.True(t, Is(w2, e))

	w3 := w2.Wrap(fmt.Errorf("message: %w", io.ErrClosedPipe))
	assert.True(t, Is(w3, io.ErrClosedPipe))

	w4 := New(str).Errorf("message: %w", io.ErrClosedPipe)
	assert.True(t, Is(w4, io.ErrClosedPipe))
}

func TestIsNested(t *testing.T) {
	t.Parallel()

	e1 := New(str).Wrap(fmt.Errorf("err: %w", io.EOF))
	assert.True(t, Is(e1, io.EOF))
	assert.False(t, Is(e1, io.ErrClosedPipe))

	e2 := e1.Wrap(io.ErrClosedPipe)
	assert.True(t, Is(e2, io.EOF))
	assert.True(t, Is(e2, io.ErrClosedPipe))

	e3 := NewErr(fmt.Errorf("err: %w", io.EOF))
	assert.True(t, Is(e3, io.EOF))
	assert.False(t, Is(e3, io.ErrClosedPipe))

	e4 := NewErr(NewErr(fmt.Errorf("err: %w", io.EOF)))
	assert.True(t, Is(e4, io.EOF))
	assert.False(t, Is(e4, io.ErrClosedPipe))

	e5 := New("message")
	e6 := New(str).Wrap(e5)
	assert.False(t, Is(e5, e6))

	assert.False(t, Is(e5, New("message")))        // this is a new value, even if it has the same message
	assert.False(t, Is(e5, fmt.Errorf("message"))) // same here
}

func TestIsEdge(t *testing.T) {
	t.Parallel()

	// edge cases
	nest2 := &wrapped{
		err:   io.EOF,
		cause: io.ErrClosedPipe,
	}
	nest1 := &wrapped{
		err: nest2,
	}

	nest3 := &wrapped{
		err:   io.EOF,
		cause: io.ErrShortBuffer,
	}

	e7 := &wrapped{
		err: nest1,
	}

	assert.True(t, Is(e7, io.EOF))
	assert.True(t, Is(e7, io.ErrClosedPipe))
	assert.True(t, Is(e7, io.EOF))
	assert.True(t, Is(e7, nest1))
	assert.True(t, Is(e7, nest2))

	assert.True(t, Is(e7, nest3))                                                              // we prioritize the topmost match, even if the inner cause don't match
	assert.True(t, Is(&wrapped{err: io.EOF}, nest3))                                           //same here
	assert.False(t, Is(&wrapped{err: io.EOF}, &wrapped{err: io.ErrClosedPipe, cause: io.EOF})) // this does not match
	assert.True(t, Is(&wrapped{err: io.ErrClosedPipe, cause: io.EOF}, io.EOF))                 // this matches the cause

	assert.True(t, Is(&wrapped{err: io.ErrClosedPipe, cause: &wrapped{err: io.EOF}}, &wrapped{err: io.EOF})) // this matches the cause
	assert.False(t, Is(&wrapped{err: io.ErrClosedPipe, cause: io.EOF}, &wrapped{err: io.EOF}))               // this does not match ???

	assert.True(t, Is(&wrapped{err: io.EOF, cause: io.ErrClosedPipe}, &wrapped{err: io.EOF})) // but this is a match

	// with type composition
	type expanded struct {
		*wrapped
	}

	e8 := expanded{wrapped: &wrapped{err: io.EOF}}
	assert.True(t, Is(e8, io.EOF))
	assert.True(t, Is(e8, &wrapped{err: io.EOF}))
}

func TestAsEdge(t *testing.T) {
	e1 := &wrapped{
		err: &wrapped{
			err: io.EOF,
		},
		cause: io.ErrClosedPipe,
	}
	target := &wrapped{
		err: &wrapped{
			err: io.EOF,
		},
	}
	assert.True(t, As(e1, &target))

	// with type composition
	type expanded struct {
		*wrapped
	}

	type inner struct {
		*wrapped
	}

	w1 := &inner{wrapped: &wrapped{err: io.ErrClosedPipe}}
	w2 := &wrapped{err: io.EOF, cause: w1}
	e2 := &expanded{wrapped: w2}

	t0 := io.ErrUnexpectedEOF // *errors.errorsString
	assert.True(t, e2.As(&t0))
	assert.EqualValues(t, io.EOF, t0)

	assert.True(t, As(e2, &t0))
	assert.EqualValues(t, io.EOF, t0)

	var t1 inner
	assert.True(t, As(e2, &t1))
	assert.EqualValues(t, *w1, t1)
	/*
		assert.True(t, As(e2, &t1))
		assert.EqualValues(t, *w1, t1)
	*/
}
