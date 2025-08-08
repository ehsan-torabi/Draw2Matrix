//go:generate fyne bundle -o data.go Icon.png
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

var Options struct {
	FlatMatrix       bool
	MatlabSaveFormat bool
	MatrixCol        int
	MatrixRow        int
	SettingsSaved    bool
}

func main() {
	a := app.New()
	w := a.NewWindow("Draw2Matrix")

	// Theme settings
	currentTheme := 0 // 0: light, 1: dark, 2: custom
	themes := []fyne.Theme{
		theme.LightTheme(),
		theme.DarkTheme(),
		&customTheme{},
	}
	a.Settings().SetTheme(themes[currentTheme])
	Options.FlatMatrix = true
	Options.MatlabSaveFormat = true
	// Create components
	paint := NewPaintWidget()
	paint.Resize(fyne.NewSize(20, 20))
	statusLabel := widget.NewLabel("start")
	// Styled buttons
	refreshBtn := widget.NewButtonWithIcon("Clear Paint", theme.DeleteIcon(), func() {
		paint.Clear()
	})
	refreshBtn.Importance = widget.HighImportance

	savePath := widget.NewEntry()
	savePath.SetPlaceHolder("Path For save file")
	savePath.Text, _ = os.Getwd()

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

	saveBtn := widget.NewButtonWithIcon("Save File", theme.DocumentSaveIcon(), func() {
		if !Options.SettingsSaved {
			dialog.ShowError(fmt.Errorf("please first save settings"), w)
			return
		}
		path := savePath.Text
		if path == "" {
			dialog.ShowError(errors.New("path is empty"), w)
			return
		}
		_, err := os.Stat(filepath.Join(path, "data.csv"))
		if os.IsNotExist(err) {
			err = SaveFile(path, "data.csv")
			if err != nil {
				return
			}
			statusLabel.SetText("Saved!")
		} else {
			dialog.ShowConfirm("Warning", "File is exists.are you want replace file?", func(b bool) {
				if b {
					err = SaveFile(path, "data.csv")
					if err != nil {
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
	})
	matlabSaveCheck.Checked = true
	flatMatrixCheck.Checked = true
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
	w.Resize(fyne.NewSize(800, 700))
	w.SetFixedSize(true)
	w.CenterOnScreen()

	// Initial setup
	a.Lifecycle().SetOnStarted(func() {
		oldStdOut := os.Stdout
		os.Stdout = nil
		paint.PrintMatrix(w, Options.FlatMatrix)
		os.Stdout = oldStdOut
	})
	a.Lifecycle().SetOnStopped(func() {
		if TempData.file != nil {
			err := os.Remove(TempData.file.Name())
			if err != nil {
				fmt.Println(err)
			}
			err = os.Remove(TempData.targetFile.Name())
			if err != nil {
				fmt.Println(err)
			}
			err = os.RemoveAll(TempData.dir)
			if err != nil {
				fmt.Println(err)
			}
			TempData.file.Close()
		}
	})
	w.ShowAndRun()
}

// Custom theme definition
type customTheme struct{}

func (t *customTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	switch name {
	case theme.ColorNameBackground:
		return color.NRGBA{R: 0x1E, G: 0x1E, B: 0x2E, A: 0xFF}
	case theme.ColorNameForeground:
		return color.NRGBA{R: 0xDD, G: 0xDD, B: 0xFF, A: 0xFF}
	case theme.ColorNamePrimary:
		return color.NRGBA{R: 0x88, G: 0x66, B: 0xFF, A: 0xFF}
	case theme.ColorNameFocus:
		return color.NRGBA{R: 0xAA, G: 0x88, B: 0xFF, A: 0xFF}
	default:
		return theme.DefaultTheme().Color(name, variant)
	}
}

func (t *customTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *customTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *customTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
