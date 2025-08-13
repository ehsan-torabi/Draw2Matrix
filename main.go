//go:generate fyne bundle -o data.go Icon.png

// Package main implements a drawing application that converts drawings to matrices
// The application supports both standard CSV format and MATLAB compatible format
// for saving the drawn patterns and their corresponding labels.
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"strconv"
)

// Options stores the global application settings
var Options struct {
	FlatMatrix         bool // Whether to flatten the matrix when saving
	MatlabSaveFormat   bool // Whether to save in MATLAB compatible format
	MatrixCol          int  // Number of columns in the output matrix
	MatrixRow          int  // Number of rows in the output matrix
	SettingsSaved      bool // Whether settings have been saved and locked
	oneHotEncodingSave bool // Whether to save target to one-hot-encoding format
}

var (
	mainApp     = app.New()
	Application struct {
		mainWindow  fyne.Window
		paintWindow fyne.Window
		paintObject *PaintWidget
	}
)

// main initializes and runs the Draw2Matrix application
func main() {
	// Initialize application and main window
	window := mainApp.NewWindow("Draw2Matrix")
	Application.mainWindow = window

	mainApp.Settings().SetTheme(theme.LightTheme())

	// Initialize default application options
	Options.FlatMatrix = false       // Default to flat matrix output
	Options.MatlabSaveFormat = false // Default to MATLAB format
	Options.MatrixRow = 20
	Options.MatrixCol = 20

	// Initialize UI components
	paint := NewPaintWidget()
	paintWindow := NewPaintWindow(mainApp, paint)
	Application.paintWindow = paintWindow
	Application.paintObject = paint
	refreshBtn.Importance = widget.HighImportance
	savePath.SetPlaceHolder("Directory path For save file")
	dataFileEntry.SetPlaceHolder("Data file name")
	dataFileEntry.SetText("data")

	targetFileEntry.SetPlaceHolder("Target file name")
	targetFileEntry.SetText("target")
	targetFileEntry.Disable()

	input.SetPlaceHolder("Enter Label")
	input.Validator = labelValidator

	exportBtn.Importance = widget.MediumImportance

	matlabSaveCheck.Checked = false
	flatMatrixCheck.Checked = false

	rowInput.SetPlaceHolder("Rows")
	rowInput.SetText(strconv.Itoa(Options.MatrixRow))
	rowInput.Validator = rowValidator

	colInput.SetPlaceHolder("Columns")
	colInput.SetText(strconv.Itoa(Options.MatrixCol))
	colInput.Validator = colValidator
	counterLabel.SetText("0")
	addBtn.Importance = widget.MediumImportance
	addAndClearPaintBtn.Importance = widget.DangerImportance

	// Main content with padding
	content := container.NewBorder(
		nil,
		container.NewPadded(bottomContainer),
		nil,
		nil,
		nil,
	)

	// Set window content and size
	window.SetContent(content)
	window.SetMaster()
	window.Resize(fyne.NewSize(800, 500))
	window.CenterOnScreen()

	// Configure application lifecycle handlers
	// OnStarted: Initialize matrix display
	mainApp.Lifecycle().SetOnStarted(onStartedApplication)
	// OnStopped: Clean up temporary files
	mainApp.Lifecycle().SetOnStopped(onStoppedApplication)

	// Start the application
	window.Show()
	paintWindow.Show()
	mainApp.Run()
}
