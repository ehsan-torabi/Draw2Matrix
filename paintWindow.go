package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func NewPaintWindow(a fyne.App, paintObject *PaintWidget) fyne.Window {
	paintWindow := a.NewWindow("Paint")
	paintWindow.SetContent(container.NewPadded(paintObject))
	paintWindow.Resize(fyne.NewSize(400, 500))
	paintWindow.SetFixedSize(true)
	return paintWindow
}
