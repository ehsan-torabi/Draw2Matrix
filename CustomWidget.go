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
)

var PrevPose fyne.Position = fyne.NewPos(0, 0)

type PaintWidget struct {
	widget.BaseWidget
	lines []*canvas.Line
}

func (p *PaintWidget) CreateRenderer() fyne.WidgetRenderer {
	rect := canvas.NewRectangle(color.RGBA{255, 255, 255, 255})
	rect.StrokeWidth = 5
	rect.Resize(fyne.NewSize(20, 20)) // Set a default size

	return &paintRenderer{
		widget:  p,
		rect:    rect,
		lines:   p.lines,
		objects: append([]fyne.CanvasObject{rect}, p.getLineObjects()...),
	}
}

func (p *PaintWidget) getLineObjects() []fyne.CanvasObject {
	objects := make([]fyne.CanvasObject, len(p.lines))
	for i, line := range p.lines {
		objects[i] = line
	}
	return objects
}

func (p *PaintWidget) MouseDown(ev *desktop.MouseEvent) {
	PrevPose = ev.Position
}

func (p *PaintWidget) MouseMoved(ev *desktop.MouseEvent) {
	if ev.Button == desktop.MouseButtonPrimary {
		newPos := ev.Position
		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 8
		line.Position1 = PrevPose
		line.Position2 = newPos
		p.lines = append(p.lines, line)
		PrevPose = ev.Position
		p.Refresh()
	}
}
func (p *PaintWidget) MouseIn(ev *desktop.MouseEvent) {
}
func (p *PaintWidget) MouseOut() {}

func (p *PaintWidget) MouseUp(ev *desktop.MouseEvent) {
}

func (p *PaintWidget) PrintMatrix(w fyne.Window, flat bool) {
	img := captureAndProcessImage(w, p)
	mat := image2BinaryMatrix(img)
	fmt.Println()
	if flat {
		fmt.Println(ConvertToFlattenMatric(mat))
	} else {
		for _, line := range mat {
			fmt.Println(line)
		}
	}

}

func (p *PaintWidget) ExportToPNG(w fyne.Window, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	result := captureAndProcessImage(w, p)
	defer file.Close()
	return png.Encode(file, result)
}

func (p *PaintWidget) Clear() {
	p.lines = []*canvas.Line{}
	p.Refresh()
}

type paintRenderer struct {
	widget  *PaintWidget
	rect    *canvas.Rectangle
	lines   []*canvas.Line
	objects []fyne.CanvasObject
}

func (r *paintRenderer) Layout(size fyne.Size) {
	r.rect.Resize(size)
}

func (r *paintRenderer) MinSize() fyne.Size {
	return fyne.NewSize(100, 100)
}

func (r *paintRenderer) Refresh() {
	r.objects = append([]fyne.CanvasObject{r.rect}, r.widget.getLineObjects()...)
	canvas.Refresh(r.widget)
}

func (r *paintRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *paintRenderer) Destroy() {}

func NewPaintWidget() *PaintWidget {
	p := &PaintWidget{
		lines: make([]*canvas.Line, 0),
	}
	p.ExtendBaseWidget(p)
	return p
}
