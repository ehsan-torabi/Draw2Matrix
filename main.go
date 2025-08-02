//go:generate fyne bundle -o data.go Icon.png
package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"os"
	"strconv"
)

func main() {
	a := app.New()
	w := a.NewWindow("Draw2Matrix")

	flatMatrix := true
	matlabSaveFormat := true
	// Create components
	paint := NewPaintWidget()
	paint.Resize(fyne.NewSize(20, 20))
	refreshBtn := widget.NewButton("Clear", func() {
		paint.Clear()
	})

	printBtn := widget.NewButton("Save File", func() {
		paint.PrintMatrix(w, flatMatrix)
	})
	flatMatrixCheck := widget.NewCheck("Flat Matric", func(b bool) {
		flatMatrix = b
	})
	matlabSaveCheck := widget.NewCheck("Matlab Save Matrix", func(b bool) {
		matlabSaveFormat = b
	})
	matlabSaveCheck.Checked = true
	flatMatrixCheck.Checked = true
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter Label")

	exportBtn := widget.NewButton("Export PNG", func() {
		filename := "draw.png"
		if input.Text != "" {
			filename = input.Text + ".png"
		}
		err := paint.ExportToPNG(w, filename)
		if err != nil {
			fmt.Printf("Export error: %s", err)
		}
	})

	rowInput := widget.NewEntry()
	rowInput.SetPlaceHolder("Row")
	rowInput.SetText(strconv.Itoa(MatrixRowNum))
	rowInput.Validator = func(s string) error {
		val, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("please enter valid number")
		}
		MatrixRowNum = val + 1
		return nil
	}
	colInput := widget.NewEntry()
	colInput.SetPlaceHolder("Column")
	colInput.SetText(strconv.Itoa(MatrixColNum))
	colInput.Validator = func(s string) error {
		val, err := strconv.Atoi(s)
		if err != nil {
			return fmt.Errorf("please enter valid number")
		}
		MatrixColNum = val + 1
		return nil
	}

	rowColContainer := container.NewGridWithColumns(2, rowInput, colInput)
	submitBtn := widget.NewButton("Add", func() {})
	labelContainer := container.NewAdaptiveGrid(3, input, container.NewVBox(flatMatrixCheck, matlabSaveCheck), submitBtn)
	bottomContainer := container.NewVBox(refreshBtn, exportBtn, labelContainer, printBtn, rowColContainer)

	// Create a container with proper layout
	content := container.NewBorder(
		nil,             // Top
		bottomContainer, // Bottom
		nil,             // Left
		nil,             // Right
		paint,           // Center
	)

	// Set window content and size
	w.SetContent(content)
	w.Resize(fyne.NewSize(600, 600)) // Window size
	w.SetFixedSize(true)
	// This is because first time process gives incorrect results
	a.Lifecycle().SetOnStarted(func() {
		oldStdOut := os.Stdout
		os.Stdout = nil
		paint.PrintMatrix(w, flatMatrix)
		fmt.Println(matlabSaveFormat)
		os.Stdout = oldStdOut
	})
	w.ShowAndRun()

}
