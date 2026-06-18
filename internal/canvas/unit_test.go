package canvas_test

import (
	"bytes"
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"angp/internal/canvas"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnit_NewCanvas(t *testing.T) {
	c := newCanvas(t, 100, 100)
	assert.Equal(t, 100, c.Width)
	assert.Equal(t, 100, c.Height)
	assert.Equal(t, canvas.ToolPen, c.Tool)
	assert.Empty(t, c.Strokes)
}

func TestUnit_NewCanvas_InvalidSize(t *testing.T) {
	tests := []struct {
		name string
		w, h int
	}{
		{"zero width", 0, 100},
		{"negative width", -1, 100},
		{"zero height", 100, 0},
		{"negative height", 100, -1},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := canvas.New(tt.w, tt.h)
			require.Error(t, err)
			assert.ErrorIs(t, err, canvas.ErrInvalidSize)
		})
	}
}

// #10 — single click with pen creates a point
func TestUnit_SingleClick(t *testing.T) {
	c := newCanvas(t, 50, 50)
	c.BeginStroke(canvas.Point{X: 25, Y: 25})

	require.Len(t, c.Strokes, 1)
	assert.Len(t, c.Strokes[0].Points, 1)
	assert.Equal(t, canvas.Point{X: 25, Y: 25}, c.Strokes[0].Points[0])
}

// #11 — points outside bounds are clamped
func TestUnit_Clamp(t *testing.T) {
	tests := []struct {
		name   string
		input  canvas.Point
		expect canvas.Point
	}{
		{"within bounds", canvas.Point{X: 50, Y: 50}, canvas.Point{X: 50, Y: 50}},
		{"negative X", canvas.Point{X: -5, Y: 50}, canvas.Point{X: 0, Y: 50}},
		{"negative Y", canvas.Point{X: 50, Y: -5}, canvas.Point{X: 50, Y: 0}},
		{"X over width", canvas.Point{X: 200, Y: 50}, canvas.Point{X: 99, Y: 50}},
		{"Y over height", canvas.Point{X: 50, Y: 200}, canvas.Point{X: 50, Y: 99}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCanvas(t, 100, 100)
			c.BeginStroke(tt.input)
			got := c.Strokes[0].Points[0]
			assert.Equal(t, tt.expect, got)
		})
	}
}

// pen draws with pen color, eraser draws with background
func TestUnit_PenVsEraser(t *testing.T) {
	black := color.RGBA{A: 255}
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}

	tests := []struct {
		name      string
		tool      canvas.Tool
		wantColor color.RGBA
	}{
		{"pen draws black", canvas.ToolPen, black},
		{"eraser draws background", canvas.ToolEraser, white},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := newCanvas(t, 50, 50)
			c.Tool = tt.tool
			c.BeginStroke(canvas.Point{X: 25, Y: 25})
			assert.Equal(t, tt.wantColor, c.Strokes[0].Color)
		})
	}
}

// stroke accumulates multiple points
func TestUnit_Stroke(t *testing.T) {
	c := newCanvas(t, 100, 100)
	c.BeginStroke(canvas.Point{X: 10, Y: 20})
	c.AddPoint(canvas.Point{X: 30, Y: 40})

	require.Len(t, c.Strokes, 1)
	assert.Len(t, c.Strokes[0].Points, 2)
	assert.Equal(t, canvas.Point{X: 30, Y: 40}, c.Strokes[0].Points[1])
}

// #5 — clear removes all strokes
func TestUnit_Clear(t *testing.T) {
	c := newCanvas(t, 50, 50)
	c.BeginStroke(canvas.Point{X: 5, Y: 5})
	c.BeginStroke(canvas.Point{X: 10, Y: 10})

	c.Clear()
	assert.Empty(t, c.Strokes)
}

// #13 — clear on empty canvas is safe
func TestUnit_Clear_AlreadyEmpty(t *testing.T) {
	c := newCanvas(t, 50, 50)
	c.Clear()
	assert.Empty(t, c.Strokes)
}

// #6 — save produces a valid PNG file
func TestUnit_Save(t *testing.T) {
	dir := t.TempDir()
	c := newCanvas(t, 64, 64)
	c.BeginStroke(canvas.Point{X: 32, Y: 32})

	path := filepath.Join(dir, "out.png")
	f, err := os.Create(path) //nolint:gosec // test file uses t.TempDir()
	require.NoError(t, err)
	defer func() { _ = f.Close() }()

	require.NoError(t, c.Save(f))

	info, err := os.Stat(path)
	require.NoError(t, err)
	assert.Greater(t, info.Size(), int64(0))
}

// #12 — save empty canvas produces valid white PNG
func TestUnit_Save_EmptyCanvas(t *testing.T) {
	c := newCanvas(t, 30, 30)

	var buf bytes.Buffer
	require.NoError(t, c.Save(&buf))

	img, err := canvas.Load(&buf)
	require.NoError(t, err)
	assert.Equal(t, 30, img.Bounds().Dx())

	r, g, b, _ := img.At(15, 15).RGBA()
	assert.Equal(t, uint32(0xffff), r)
	assert.Equal(t, uint32(0xffff), g)
	assert.Equal(t, uint32(0xffff), b)
}

// #14 — save to bad writer returns error, no panic
func TestUnit_Save_BadWriter(t *testing.T) {
	c := newCanvas(t, 10, 10)
	c.BeginStroke(canvas.Point{X: 5, Y: 5})

	err := c.Save(&failWriter{})
	assert.Error(t, err)
}

// #7 — load valid PNG returns correct dimensions
func TestUnit_Load(t *testing.T) {
	c := newCanvas(t, 40, 40)
	c.BeginStroke(canvas.Point{X: 20, Y: 20})

	var buf bytes.Buffer
	require.NoError(t, c.Save(&buf))

	img, err := canvas.Load(&buf)
	require.NoError(t, err)
	assert.Equal(t, 40, img.Bounds().Dx())
	assert.Equal(t, 40, img.Bounds().Dy())
}

// #8 + #15 — load invalid/non-PNG data returns error, no panic
func TestUnit_Load_Invalid(t *testing.T) {
	tests := []struct {
		name string
		data string
	}{
		{"empty", ""},
		{"garbage", "not a png"},
		{"partial png header", "\x89PNG"},
		{"jpeg header", "\xff\xd8\xff\xe0"},
		{"gif header", "GIF89a"},
		{"plain text", "hello world"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := canvas.Load(strings.NewReader(tt.data))
			assert.Error(t, err)
		})
	}
}
