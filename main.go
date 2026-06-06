package main

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	angp := app.New()

	mainWin := angp.NewWindow("angp")
	mainWin.SetMaster() // set as main window

	mainWin.Resize(fyne.NewSize(400, 400))

	fmt.Println(fyne.CurrentDevice())
	mainWin.SetContent(toolBar(angp))

	mainWin.ShowAndRun()
	tidyUp()
}

func toolBar(app fyne.App) *fyne.Container {
	message := widget.NewLabel("Welcome")

	button := widget.NewButton("Update", func() {
		formatted := time.Now().Format("Time: 03:04:05")
		message.SetText(formatted)
	})
	colorPicker := widget.NewButton("color", func() {
		wColor := app.NewWindow("Color picker")
		wColor.SetContent(widget.NewLabel("TBD"))
		wColor.Show()
	})

	toolbar := container.NewVBox(message, button, colorPicker)
	return toolbar
}

func tidyUp() {
	fmt.Println("Exited")
}
