//go:generate fyne bundle -o data.go Icon.png

// Package main implements a drawing application that converts drawings to matrices
// The application supports both standard CSV format and MATLAB compatible format
// for saving the drawn patterns and their corresponding labels.
package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
	"os"
	"path/filepath"
	"strconv"
)

// Options stores the global application settings
var Options struct {
	FlatMatrix       bool  // Whether to flatten the matrix when saving
	MatlabSaveFormat bool  // Whether to save in MATLAB compatible format
	MatrixCol        int   // Number of columns in the output matrix
	MatrixRow        int   // Number of rows in the output matrix
	SettingsSaved    bool  // Whether settings have been saved and locked
}

// main initializes and runs the Draw2Matrix application
func main() {
	// Initialize application and main window
	a := app.New()
	w := a.NewWindow("Draw2Matrix")

	// Setup theme configuration
	currentTheme := 0 // Theme indices: 0=light, 1=dark, 2=custom
	themes := []fyne.Theme{
		theme.LightTheme(),
		theme.DarkTheme(),
		&customTheme{},
	}
	a.Settings().SetTheme(themes[currentTheme])

	// Initialize default application options
	Options.FlatMatrix = true       // Default to flat matrix output
	Options.MatlabSaveFormat = true // Default to MATLAB format
	
	// Initialize UI components
	paint := NewPaintWidget()
	paint.Resize(fyne.NewSize(20, 20))
	statusLabel := widget.NewLabel("start")
	// Styled buttons
	refreshBtn := widget.NewButtonWithIcon("Clear Paint", theme.DeleteIcon(), func() {
		paint.Clear()
	})
	refreshBtn.Importance = widget.HighImportance

	savePath := widget.NewEntry()
	savePath.SetPlaceHolder("Directory path For save file")

	dataFileEntry := widget.NewEntry()
	dataFileEntry.SetPlaceHolder("Data file name")
	dataFileEntry.SetText("data")

	targetFileEntry := widget.NewEntry()
	targetFileEntry.SetPlaceHolder("Target file name")
	targetFileEntry.SetText("target")
	targetFileEntry.Hide()

	changePath := widget.NewButtonWithIcon("Browse", theme.FolderIcon(), func() {
		dialog.ShowFolderOpen(func(uc fyne.ListableURI, err error) {
			if err != nil {
				dialog.ShowError(err, w)
				return
			}
			if uc == nil {
				return
			}
			DirPath := uc.Path()
			savePath.SetText(DirPath)
		}, w)
	})

	// Create save button with file saving functionality
	saveBtn := widget.NewButtonWithIcon("Save File", theme.DocumentSaveIcon(), func() {
		// Validate settings are saved
		if !Options.SettingsSaved {
			dialog.ShowError(fmt.Errorf("please first save settings"), w)
			return
		}
		
		// Validate save path
		path := savePath.Text
		if path == "" {
			dialog.ShowError(errors.New("path is empty"), w)
			return
		}
		
		// Validate data filename
		dataFileName := dataFileEntry.Text
		if dataFileName == "" {
			dialog.ShowError(errors.New("data file name is empty"), w)
			return
		}

		// Handle MATLAB format saving
		if Options.MatlabSaveFormat {
			targetFileName := targetFileEntry.Text
			if targetFileName == "" {
				dialog.ShowError(errors.New("target file name is empty"), w)
				return
			}
			if err := SaveFileForMatlab(path, dataFileName, targetFileName); err != nil {
				statusLabel.SetText("Not Saved!")
				return
			}
			statusLabel.SetText("Saved!")
			return
		}

		// Handle CSV format saving
		filePath := filepath.Join(path, dataFileName+".csv")
		if _, err := os.Stat(filePath); os.IsNotExist(err) {
			// File doesn't exist, save directly
			if err := SaveFile(path, dataFileName+".csv"); err != nil {
				statusLabel.SetText("Not Saved!")
				return
			}
			statusLabel.SetText("Saved!")
		} else {
			// File exists, ask for confirmation
			dialog.ShowConfirm("Warning", "File exists. Do you want to replace it?", func(b bool) {
				if b {
					if err := SaveFile(path, "data.csv"); err != nil {
						statusLabel.SetText("Not Saved!")
						return
					}
					statusLabel.SetText("Saved!")
				} else {
					statusLabel.SetText("Not Saved!")
				}
			}, w)
		}

	})
	input := widget.NewEntry()
	input.SetPlaceHolder("Enter Label")
	input.Validator = func(s string) error {
		if len(s) > 20 {
			return fmt.Errorf("label too long")
		}
		return nil
	}
	exportBtn := widget.NewButtonWithIcon("Export PNG", theme.FileImageIcon(), func() {
		filename := "draw.png"
		if input.Text != "" {
			filename = input.Text + ".png"
		}
		err := paint.ExportToPNG(w, filename)
		if err != nil {
			fmt.Printf("Export error: %s", err)
		}
	})
	exportBtn.Importance = widget.MediumImportance

	flatMatrixCheck := widget.NewCheck("Flat Matrix", func(b bool) {
		Options.FlatMatrix = b
	})
	matlabSaveCheck := widget.NewCheck("Matlab Save Format", func(b bool) {
		Options.MatlabSaveFormat = b
		if b {
			targetFileEntry.Show()
			flatMatrixCheck.Disable()
			flatMatrixCheck.SetChecked(false)
			Options.FlatMatrix = false
		} else {
			targetFileEntry.Hide()
			flatMatrixCheck.Enable()
		}
	})
	matlabSaveCheck.Checked = false
	flatMatrixCheck.Checked = false
	Options.MatrixRow = 20
	Options.MatrixCol = 20
	rowInput := widget.NewEntry()
	rowInput.SetPlaceHolder("Rows")
	rowInput.SetText(strconv.Itoa(Options.MatrixRow))
	rowInput.Validator = func(s string) error {
		val, err := strconv.Atoi(s)
		if err != nil || val <= 0 {
			return fmt.Errorf("enter number ")
		}
		Options.MatrixRow = val + 1
		return nil
	}

	colInput := widget.NewEntry()
	colInput.SetPlaceHolder("Columns")
	colInput.SetText(strconv.Itoa(Options.MatrixCol))
	colInput.Validator = func(s string) error {
		val, err := strconv.Atoi(s)
		if err != nil || val <= 0 {
			return fmt.Errorf("enter number")
		}
		Options.MatrixCol = val + 1
		return nil
	}

	submitBtn := widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), func() {
		if !Options.SettingsSaved {
			dialog.ShowError(fmt.Errorf("please first save settings"), w)
			return
		}
		if input.Text != "" {
			if Options.MatlabSaveFormat {
				err := AddToFileForMatlab(paint.GetMatrix(w), input.Text)
				if err != nil {
					dialog.ShowError(fmt.Errorf("error to add matrix"), w)
					return
				}
				statusLabel.SetText("Added!")
				return
			}
			err := AddToFile(paint.GetMatrix(w), input.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error to add matrix"), w)
				return
			}
			statusLabel.SetText("Added!")
			return
		}
		dialog.ShowError(fmt.Errorf("please enter valid label"), w)
	})
	submitBtn.Importance = widget.MediumImportance

	// Theme switcher
	themeBtn := widget.NewButtonWithIcon("Theme", theme.ColorPaletteIcon(), func() {
		currentTheme = (currentTheme + 1) % len(themes)
		a.Settings().SetTheme(themes[currentTheme])
	})
	saveOptionsBtn := widget.NewButtonWithIcon("Save Settings", theme.SettingsIcon(), func() {
		rowInput.Disable()
		colInput.Disable()
		flatMatrixCheck.Disable()
		matlabSaveCheck.Disable()
		Options.SettingsSaved = true
		InitializeTemps(Options.MatlabSaveFormat)
	})
	resetProjectBtn := widget.NewButtonWithIcon("Reset Project", theme.ContentClearIcon(), func() {
		dialog.ShowConfirm("Warning", "Are you sure you want to do that?\nthis is delete your added matrix if you dont saves it. ",
			func(choice bool) {
				rowInput.Enable()
				colInput.Enable()
				flatMatrixCheck.Enable()
				matlabSaveCheck.Enable()
				Options.SettingsSaved = false
				if TempData.file != nil {
					err := os.Remove(TempData.file.Name())
					if err != nil {
						fmt.Println(err)
					}
					err = os.RemoveAll(TempData.dir)
					if err != nil {
						fmt.Println(err)
					}
					TempData.file.Close()
				}
			}, w,
		)
	})
	// Layout containers
	settingsContainer := container.NewVBox(
		widget.NewLabel("Matrix Settings:"),
		container.NewGridWithColumns(2, rowInput, colInput),
		container.NewGridWithColumns(2, flatMatrixCheck, matlabSaveCheck),
		container.NewGridWithColumns(2, resetProjectBtn, saveOptionsBtn),
	)

	pathContainer := container.NewVBox(
		container.NewGridWithColumns(2, savePath, changePath),
		dataFileEntry,
		targetFileEntry,
	)
	actionContainer := container.NewVBox(
		widget.NewLabel("Actions:"),
		container.NewGridWithColumns(2, refreshBtn, exportBtn),
		pathContainer, saveBtn,
	)

	labelContainer := container.NewVBox(
		widget.NewLabel("Label:"),
		container.NewBorder(nil, nil, nil, submitBtn, input),
	)

	themeContainer := container.NewVBox(
		container.NewBorder(nil, nil, nil, statusLabel, themeBtn),
	)

	bottomContainer := container.NewVBox(
		container.NewPadded(settingsContainer),
		container.NewPadded(actionContainer),
		container.NewPadded(labelContainer),
		container.NewPadded(themeContainer),
	)

	// Main content with padding
	content := container.NewBorder(
		nil,
		container.NewPadded(bottomContainer),
		nil,
		nil,
		container.NewPadded(paint),
	)

	// Set window content and size
	w.SetContent(content)
	w.Resize(fyne.NewSize(800, 800))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	// Configure application lifecycle handlers
	
	// OnStarted: Initialize matrix display
	a.Lifecycle().SetOnStarted(func() {
		// Temporarily disable stdout to prevent matrix printing
		oldStdOut := os.Stdout
		os.Stdout = nil
		paint.PrintMatrix(w, Options.FlatMatrix)
		os.Stdout = oldStdOut
	})
	
	// OnStopped: Clean up temporary files
	a.Lifecycle().SetOnStopped(func() {
		if TempData.file != nil {
			// Remove temporary files and directory
			if err := os.Remove(TempData.file.Name()); err != nil {
				fmt.Println("Error removing temp file:", err)
			}
			if err := os.Remove(TempData.targetFile.Name()); err != nil {
				fmt.Println("Error removing temp target file:", err)
			}
			if err := os.RemoveAll(TempData.dir); err != nil {
				fmt.Println("Error removing temp directory:", err)
			}
			TempData.file.Close()
		}
	})
	
	// Start the application
	w.ShowAndRun()
}

// customTheme implements a custom dark theme with purple accents
type customTheme struct{}

// Color returns the color for the specified theme element
func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x1E, G: 0x1E, B: 0x2E, A: 0xFF} // Dark blue-purple background
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xDD, G: 0xDD, B: 0xFF, A: 0xFF} // Light purple text
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x88, G: 0x66, B: 0xFF, A: 0xFF} // Medium purple accents
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xAA, G: 0x88, B: 0xFF, A: 0xFF} // Light purple focus
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

// Font returns the font resource for the specified text style
func (t *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

// Icon returns the icon resource for the specified icon name
func (t *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

// Size returns the size for the specified theme element
func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
