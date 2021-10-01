package errors

import (
	"fmt"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRooted(t *testing.T) {
	e1 := NewErrWithRoot(io.EOF)
	e2 := e1.Wrap(io.ErrClosedPipe)

	assert.True(t, Is(e2, io.EOF))
	assert.True(t, Is(e2, io.ErrClosedPipe))

	e3 := NewErrWithRoot(e1.Wrap(io.ErrClosedPipe).Wrap(io.ErrUnexpectedEOF).Wrap(io.ErrShortBuffer))
	assert.True(t, Is(e3, io.ErrShortBuffer))
	assert.ErrorIs(t, e3.Root(), io.ErrShortBuffer)

	assert.ErrorIs(t, Root(e3), io.ErrShortBuffer)

	// stdlib errors
	e4 := fmt.Errorf("err3: %w", io.ErrClosedPipe)
	e5 := fmt.Errorf("err4: %w", e4)

	assert.ErrorIs(t, Root(e4), io.ErrClosedPipe)
	assert.ErrorIs(t, Root(e5), io.ErrClosedPipe)
}
