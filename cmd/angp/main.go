package main

import (
	"image"

	"angp/internal/canvas"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	fynecanvas "fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("angp")
	w.SetMaster()
	w.Resize(fyne.NewSize(600, 400))

	c, _ := canvas.New(600, 400)

	var drawing bool
	var img *image.RGBA

	img = c.Render()

	raster := fynecanvas.NewRaster(func(width, height int) image.Image {
		return img
	})

	refresh := func() {
		img = c.Render()
		fynecanvas.Refresh(raster)
	}

	drawArea := newDrawable(raster, func(pos fyne.Position, widgetSize fyne.Size) {
		// scale from widget coordinates to canvas pixel coordinates
		scaleX := float64(c.Width) / float64(widgetSize.Width)
		scaleY := float64(c.Height) / float64(widgetSize.Height)
		p := canvas.Point{
			X: int(float64(pos.X) * scaleX),
			Y: int(float64(pos.Y) * scaleY),
		}
		if !drawing {
			drawing = true
			c.BeginStroke(p)
		} else {
			c.AddPoint(p)
		}
		refresh()
	}, func() {
		drawing = false
	})

	toolbar := buildToolbar(c, refresh, w)

	content := container.NewBorder(nil, nil, nil, toolbar, drawArea)
	w.SetContent(content)
	w.ShowAndRun()
}

func buildToolbar(c *canvas.Canvas, refresh func(), w fyne.Window) *fyne.Container {
	penBtn := widget.NewButtonWithIcon("Pen", theme.DocumentCreateIcon(), func() {
		c.Tool = canvas.ToolPen
	})
	eraserBtn := widget.NewButtonWithIcon("Eraser", theme.DeleteIcon(), func() {
		c.Tool = canvas.ToolEraser
	})

	sizeLabel := widget.NewLabel("Size: M")
	sizeS := widget.NewButton("S", func() {
		c.Size = 1
		sizeLabel.SetText("Size: S")
	})
	sizeM := widget.NewButton("M", func() {
		c.Size = 2
		sizeLabel.SetText("Size: M")
	})
	sizeL := widget.NewButton("L", func() {
		c.Size = 5
		sizeLabel.SetText("Size: L")
	})

	clearBtn := widget.NewButtonWithIcon("Clear", theme.ContentClearIcon(), func() {
		c.Clear()
		refresh()
	})

	saveBtn := widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), func() {
		d := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
			if err != nil || writer == nil {
				return
			}
			defer func() { _ = writer.Close() }()
			_ = c.Save(writer)
		}, w)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))
		d.SetFileName("drawing.png")
		d.Show()
	})

	loadBtn := widget.NewButtonWithIcon("Load", theme.FolderOpenIcon(), func() {
		d := dialog.NewFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err != nil || reader == nil {
				return
			}
			defer func() { _ = reader.Close() }()
			loaded, loadErr := canvas.Load(reader)
			if loadErr != nil {
				return
			}
			applyLoaded(c, loaded)
			refresh()
		}, w)
		d.SetFilter(storage.NewExtensionFileFilter([]string{".png"}))
		d.Show()
	})

	sizes := container.NewHBox(sizeS, sizeM, sizeL)

	return container.NewVBox(
		penBtn,
		eraserBtn,
		widget.NewSeparator(),
		sizeLabel,
		sizes,
		widget.NewSeparator(),
		clearBtn,
		widget.NewSeparator(),
		saveBtn,
		loadBtn,
	)
}

func applyLoaded(c *canvas.Canvas, img *image.RGBA) {
	c.Clear()
	b := img.Bounds()
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			r, g, bl, _ := img.At(x, y).RGBA()
			if r == 0xffff && g == 0xffff && bl == 0xffff {
				continue
			}
			c.Color = img.RGBAAt(x, y)
			c.Size = 1
			c.BeginStroke(canvas.Point{X: x, Y: y})
		}
	}
}

// drawable wraps a raster and captures mouse drag/tap for drawing.
type drawable struct {
	widget.BaseWidget
	raster    *fynecanvas.Raster
	onDrag    func(fyne.Position, fyne.Size)
	onEndDrag func()
}

func newDrawable(raster *fynecanvas.Raster, onDrag func(fyne.Position, fyne.Size), onEnd func()) *drawable {
	d := &drawable{raster: raster, onDrag: onDrag, onEndDrag: onEnd}
	d.ExtendBaseWidget(d)
	return d
}

func (d *drawable) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(d.raster)
}

func (d *drawable) Dragged(ev *fyne.DragEvent) {
	d.onDrag(ev.Position, d.Size())
}

func (d *drawable) DragEnd() {
	d.onEndDrag()
}

func (d *drawable) Tapped(ev *fyne.PointEvent) {
	d.onDrag(ev.Position, d.Size())
	d.onEndDrag()
}

var _ fyne.Draggable = (*drawable)(nil)
var _ fyne.Tappable = (*drawable)(nil)
