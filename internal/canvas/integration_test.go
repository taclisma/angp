package canvas_test

import (
	"bytes"
	"strings"
	"testing"

	"angp/internal/canvas"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// #5 + #6 + #7 — draw with pen + erase + save + load + verify pixels
func TestIntegration_DrawSaveLoad(t *testing.T) {
	c := newCanvas(t, 100, 100)
	c.Size = 1

	// draw black stroke
	c.BeginStroke(canvas.Point{X: 50, Y: 50})
	c.AddPoint(canvas.Point{X: 51, Y: 50})
	c.AddPoint(canvas.Point{X: 52, Y: 50})

	// erase middle point
	c.Tool = canvas.ToolEraser
	c.BeginStroke(canvas.Point{X: 51, Y: 50})

	// save & load
	var buf bytes.Buffer
	require.NoError(t, c.Save(&buf))
	img, err := canvas.Load(&buf)
	require.NoError(t, err)

	// drawn pixel is black
	r, g, b, _ := img.At(50, 50).RGBA()
	assert.Equal(t, uint32(0), r)
	assert.Equal(t, uint32(0), g)
	assert.Equal(t, uint32(0), b)

	// erased pixel is white
	r2, g2, b2, _ := img.At(51, 50).RGBA()
	assert.Equal(t, uint32(0xffff), r2)
	assert.Equal(t, uint32(0xffff), g2)
	assert.Equal(t, uint32(0xffff), b2)
}

// #5 — draw then clear then save = white image
func TestIntegration_ClearThenSave(t *testing.T) {
	c := newCanvas(t, 50, 50)
	c.BeginStroke(canvas.Point{X: 25, Y: 25})
	c.Clear()

	var buf bytes.Buffer
	require.NoError(t, c.Save(&buf))
	img, err := canvas.Load(&buf)
	require.NoError(t, err)

	r, g, b, _ := img.At(25, 25).RGBA()
	assert.Equal(t, uint32(0xffff), r)
	assert.Equal(t, uint32(0xffff), g)
	assert.Equal(t, uint32(0xffff), b)
}

// #8 — load invalid PNG, canvas preserves its previous state
func TestIntegration_LoadInvalid_PreservesState(t *testing.T) {
	c := newCanvas(t, 100, 100)
	c.Size = 1
	c.BeginStroke(canvas.Point{X: 50, Y: 50})

	// attempt to load garbage
	_, err := canvas.Load(strings.NewReader("not a png"))
	require.Error(t, err)

	// canvas still has its stroke, render and verify pixel is black
	img := c.Render()
	r, g, b, _ := img.At(50, 50).RGBA()
	assert.Equal(t, uint32(0), r)
	assert.Equal(t, uint32(0), g)
	assert.Equal(t, uint32(0), b)
}

// #14 — save to bad writer fails, canvas preserves state
func TestIntegration_SaveBadWriter_PreservesState(t *testing.T) {
	c := newCanvas(t, 100, 100)
	c.Size = 1
	c.BeginStroke(canvas.Point{X: 50, Y: 50})

	// attempt to save to failing writer
	err := c.Save(&failWriter{})
	require.Error(t, err)

	// canvas still works — save to a valid buffer and reload
	var buf bytes.Buffer
	require.NoError(t, c.Save(&buf))
	img, loadErr := canvas.Load(&buf)
	require.NoError(t, loadErr)

	// pixel at (50,50) is still black
	r, g, b, _ := img.At(50, 50).RGBA()
	assert.Equal(t, uint32(0), r)
	assert.Equal(t, uint32(0), g)
	assert.Equal(t, uint32(0), b)
}

// #15 — load non-PNG file, canvas preserves state
func TestIntegration_LoadNonPNG_PreservesState(t *testing.T) {
	c := newCanvas(t, 100, 100)
	c.Size = 1
	c.BeginStroke(canvas.Point{X: 50, Y: 50})

	// attempt to load JPEG data
	_, err := canvas.Load(strings.NewReader("\xff\xd8\xff\xe0"))
	require.Error(t, err)

	// canvas still has its stroke
	img := c.Render()
	r, g, b, _ := img.At(50, 50).RGBA()
	assert.Equal(t, uint32(0), r)
	assert.Equal(t, uint32(0), g)
	assert.Equal(t, uint32(0), b)
}
