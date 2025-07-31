package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Paint App")

	// Create components
	paint := NewPaintWidget()
	paint.Resize(fyne.NewSize(20, 20))
	refreshBtn := widget.NewButton("Refresh", func() {
		paint.Clear()
	})
	exportBtn := widget.NewButton("Export PNG", func() {
		err := paint.ExportToPNG(w, "drawing.png")
		if err != nil {
			panic(err)
		}
	})

	// Create a container with proper layout
	content := container.NewBorder(
		nil,                                      // Top
		container.NewVBox(refreshBtn, exportBtn), // Bottom
		nil,                                      // Left
		nil,                                      // Right
		paint,                                    // Center
	)

	// Set window content and size
	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 500)) // Window size
	w.SetFixedSize(true)
	w.ShowAndRun()
}
