package canvas_test

import (
	"os"
	"testing"

	"angp/internal/canvas"

	"github.com/stretchr/testify/require"
)

// newCanvas creates a canvas or fails the test.
func newCanvas(t *testing.T, w, h int) *canvas.Canvas {
	t.Helper()
	c, err := canvas.New(w, h)
	require.NoError(t, err)
	return c
}

// failWriter always returns a permission error.
type failWriter struct{}

func (f *failWriter) Write([]byte) (int, error) {
	return 0, os.ErrPermission
}
