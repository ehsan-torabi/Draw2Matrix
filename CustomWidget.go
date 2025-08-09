// Package main provides custom widget implementations for the Draw2Matrix application
package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
)

// PrevPos stores the previous mouse position for drawing continuous lines
var PrevPos fyne.Position = fyne.NewPos(0, 0)

// PaintWidget represents a custom widget for drawing
// It extends the base widget and maintains a collection of lines
type PaintWidget struct {
	widget.BaseWidget
	lines []*canvas.Line // Collection of lines drawn on the widget
}

// CreateRenderer implements the Widget interface, creating a new renderer for the paint widget
func (p *PaintWidget) CreateRenderer() fyne.WidgetRenderer {
	// Create background rectangle
	rect := canvas.NewRectangle(color.RGBA{255, 255, 255, 255})
	rect.StrokeColor = color.RGBA{100, 100, 100, 255}
	rect.StrokeWidth = 1
	rect.Resize(fyne.NewSize(20, 20))

	return &paintRenderer{
		widget:  p,
		rect:    rect,
		lines:   p.lines,
		objects: append([]fyne.CanvasObject{rect}, p.getLineObjects()...),
	}
}

// getLineObjects converts the lines array to canvas objects
func (p *PaintWidget) getLineObjects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(p.lines))
	for i, line := range p.lines {
		objects[i] = line
	}
	return objects
}

// MouseDown handles mouse button press events
func (p *PaintWidget) MouseDown(ev *desktop.MouseEvent) {
	PrevPos = ev.Position
}

// MouseMoved handles mouse movement events
// Creates new line segments when the primary button is pressed
func (p *PaintWidget) MouseMoved(ev *desktop.MouseEvent) {
	if ev.Button == desktop.MouseButtonPrimary {
		newPos := ev.Position
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 8
		line.Position1 = PrevPos
		line.Position2 = newPos
		p.lines = append(p.lines, line)
		PrevPos = ev.Position
		p.Refresh()
	}
}

// MouseIn handles mouse enter events
func (p *PaintWidget) MouseIn(ev *desktop.MouseEvent) {}

// MouseOut handles mouse leave events
func (p *PaintWidget) MouseOut() {}

// MouseUp handles mouse button release events
func (p *PaintWidget) MouseUp(ev *desktop.MouseEvent) {}

// PrintMatrix outputs the current drawing as a binary matrix
// If flat is true, outputs as a flattened array
func (p *PaintWidget) PrintMatrix(w fyne.Window, flat bool) {
	img := captureAndProcessImage(w, p)
	mat := image2BinaryMatrix(img)
	fmt.Println()
	if flat {
		fmt.Println(ToFlattenMatrix(mat))
	} else {
		for _, line := range mat {
			fmt.Println(line)
		}
	}
}

// GetMatrix returns the current drawing as a binary matrix
func (p *PaintWidget) GetMatrix(w fyne.Window) [][]int8 {
	img := captureAndProcessImage(w, p)
	return image2BinaryMatrix(img)
}

// ExportToPNG saves the current drawing as a PNG file
func (p *PaintWidget) ExportToPNG(w fyne.Window, filename string) error {
	// Create output directory if it doesn't exist
	wd, _ := os.Getwd()
	err := os.Mkdir(filepath.Join(wd, "output"), os.ModePerm)
	if err != nil && !os.IsExist(err) {
		return err
	}
	
	// Create output file
	file, err := os.Create(filepath.Join("output", filename))
	if err != nil {
		return err
	}
	defer file.Close()

	// Process and save image
	result := captureAndProcessImage(w, p)
	return png.Encode(file, result)
}

// Clear removes all drawn lines from the widget
func (p *PaintWidget) Clear() {
	p.lines = []*canvas.Line{}
	p.Refresh()
}

// paintRenderer implements the fyne.WidgetRenderer interface
type paintRenderer struct {
	widget  *PaintWidget
	rect    *canvas.Rectangle
	lines   []*canvas.Line
	objects []fyne.CanvasObject
}

// Layout implements WidgetRenderer interface
func (r *paintRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

// MinSize implements WidgetRenderer interface
func (r *paintRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

// Refresh implements WidgetRenderer interface
func (r *paintRenderer) Refresh() {
	r.objects = append([]fyne.CanvasObject{r.rect}, r.widget.getLineObjects()...)
	canvas.Refresh(r.widget)
}

// Objects implements WidgetRenderer interface
func (r *paintRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

// Destroy implements WidgetRenderer interface
func (r *paintRenderer) Destroy() {}

// NewPaintWidget creates and initializes a new paint widget
func NewPaintWidget() *PaintWidget {
	p := &PaintWidget{
		lines: make([]*canvas.Line, 0),
	}
	p.ExtendBaseWidget(p)
	return p
}
