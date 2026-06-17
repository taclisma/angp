// Package canvas provides domain logic for a drawing canvas with pen and eraser.
package canvas

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"io"
)

type Tool int

const (
	ToolPen Tool = iota
	ToolEraser
)

type Point struct {
	X, Y int
}

type Stroke struct {
	Points []Point
	Color  color.RGBA
	Size   int
}

type Canvas struct {
	Width, Height int
	Strokes       []Stroke
	Tool          Tool
	Color         color.RGBA
	Size          int
	Background    color.RGBA
}

var ErrInvalidSize = errors.New("dimensions must be positive")

// New creates a canvas. Width and height must be positive.
func New(width, height int) (*Canvas, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("w=%d h=%d: %w", width, height, ErrInvalidSize)
	}
	return &Canvas{
		Width:      width,
		Height:     height,
		Tool:       ToolPen,
		Color:      color.RGBA{A: 255},                          // black
		Size:       2,
		Background: color.RGBA{R: 255, G: 255, B: 255, A: 255}, // white
	}, nil
}

// BeginStroke starts a new stroke at point p.
func (c *Canvas) BeginStroke(p Point) {
	p = c.clamp(p)
	col := c.Color
	if c.Tool == ToolEraser {
		col = c.Background
	}
	c.Strokes = append(c.Strokes, Stroke{
		Points: []Point{p},
		Color:  col,
		Size:   c.Size,
	})
}

// AddPoint appends a point to the current stroke.
func (c *Canvas) AddPoint(p Point) {
	p = c.clamp(p)
	if len(c.Strokes) == 0 {
		c.BeginStroke(p)
		return
	}
	last := &c.Strokes[len(c.Strokes)-1]
	last.Points = append(last.Points, p)
}

// Clear removes all strokes.
func (c *Canvas) Clear() {
	c.Strokes = nil
}

// Render draws all strokes onto a new image.
func (c *Canvas) Render() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, c.Width, c.Height))
	// fill background
	for y := range c.Height {
		for x := range c.Width {
			img.SetRGBA(x, y, c.Background)
		}
	}
	// draw strokes
	for _, s := range c.Strokes {
		for _, p := range s.Points {
			c.drawDot(img, p, s.Color, s.Size)
		}
	}
	return img
}

// Save writes the rendered canvas as PNG.
func (c *Canvas) Save(w io.Writer) error {
	if err := png.Encode(w, c.Render()); err != nil {
		return fmt.Errorf("encoding png: %w", err)
	}
	return nil
}

// Load decodes a PNG from the reader.
func Load(r io.Reader) (*image.RGBA, error) {
	img, err := png.Decode(r)
	if err != nil {
		return nil, fmt.Errorf("decoding png: %w", err)
	}
	b := img.Bounds()
	rgba := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			rgba.Set(x, y, img.At(x, y))
		}
	}
	return rgba, nil
}

func (c *Canvas) clamp(p Point) Point {
	p.X = max(0, min(p.X, c.Width-1))
	p.Y = max(0, min(p.Y, c.Height-1))
	return p
}

func (c *Canvas) drawDot(img *image.RGBA, p Point, col color.RGBA, size int) {
	half := size / 2
	for dy := -half; dy <= half; dy++ {
		for dx := -half; dx <= half; dx++ {
			px, py := p.X+dx, p.Y+dy
			if px >= 0 && px < c.Width && py >= 0 && py < c.Height {
				img.SetRGBA(px, py, col)
			}
		}
	}
}
