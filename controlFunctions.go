package main

import (
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/dialog"
	"image/color"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

var SavedProject struct {
	Options struct {
		FlatMatrix         bool
		MatlabSaveFormat   bool
		MatrixCol          int
		MatrixRow          int
		SettingsSaved      bool
		OneHotEncodingSave bool
	}
	TempData struct {
		Saved      bool
		buffer     bytes.Buffer
		TempMatrix [][]int8
		TempTarget []string
	}
	OneHotDictionary struct {
		Dictionary map[string]interface{}
		Values     []string
	}
	CounterValue   string
	DataFilePath   string
	TargetFilePath string
	Buffer         []byte
}

func addLabelAnimation(obj *canvas.Text) {
	green := color.NRGBA{G: 0xff, A: 0xff}
	canvas.NewColorRGBAAnimation(green, color.Black, time.Second*2, func(c color.Color) {
		obj.Color = c
		obj.Refresh()
	}).Start()
}

func exportFileOperation() {
	if !Options.SettingsSaved {
		dialog.ShowError(fmt.Errorf("please first save settings"), Application.mainWindow)
		return
	}
	if counterLabel.Text == "0" {
		dialog.ShowError(fmt.Errorf("Please first add at least 1 label"), Application.mainWindow)
		return
	}

	// Validate save path
	path := savePath.Text
	if path == "" {
		dialog.ShowError(errors.New("path is empty"), Application.mainWindow)
		return
	}

	// Validate data filename
	dataFileName := dataFileEntry.Text
	if dataFileName == "" {
		dialog.ShowError(errors.New("data file name is empty"), Application.mainWindow)
		return
	}

	// Handle MATLAB format saving
	if Options.MatlabSaveFormat {
		targetFileName := targetFileEntry.Text
		if targetFileName == "" {
			dialog.ShowError(errors.New("target file name is empty"), Application.mainWindow)
			return
		}
		if err := SaveFileForMatlab(path, dataFileName, targetFileName); err != nil {
			statusLabel.Text = "Not Saved!"
			return
		}
		statusLabel.Text = "Saved!"
		return
	}

	// Handle CSV format saving
	filePath := filepath.Join(path, dataFileName+".csv")
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// file doesn't exist, save directly
		if err = SaveFile(path, dataFileName+".csv"); err != nil {
			statusLabel.Text = "Not Saved!"
			return
		}
		statusLabel.Text = "Saved!"
	} else {
		// file exists, ask for confirmation
		dialog.ShowConfirm("Warning", "file exists. Do you want to replace it?", func(b bool) {
			if b {
				if err = SaveFile(path, "data.csv"); err != nil {
					statusLabel.Text = "Not Saved!"
					return
				}
				statusLabel.Text = "Saved!"
			} else {
				statusLabel.Text = "Not Saved!"
			}
		}, Application.mainWindow)
	}
}
func browseOperation() {
	dialog.ShowFolderOpen(func(uc fyne.ListableURI, err error) {
		if err != nil {
			dialog.ShowError(err, Application.mainWindow)
			return
		}
		if uc == nil {
			return
		}
		DirPath := uc.Path()
		savePath.SetText(DirPath)
	}, Application.mainWindow)
}
func openPaintWindowOperation() {
	if Application.paintWindow != nil {
		Application.paintWindow.Close()
	}
	Application.paintWindow = NewPaintWindow(mainApp, Application.paintObject)
	Application.paintWindow.Show()
}
func matlabSaveCheckBoxFunction(b bool) {
	Options.MatlabSaveFormat = b
	if b {
		targetFileEntry.Enable()
		flatMatrixCheck.Disable()
		flatMatrixCheck.SetChecked(false)
		Options.FlatMatrix = false
		Application.mainWindow.Canvas().Refresh(Application.mainWindow.Content())
	} else {
		targetFileEntry.Disable()
		flatMatrixCheck.Enable()
		Application.mainWindow.Canvas().Refresh(Application.mainWindow.Content())
	}
}
func expertPNGOperation() {
	filename := "draw.png"
	if input.Text != "" {
		filename = input.Text + ".png"
	}
	err := Application.paintObject.ExportToPNG(Application.paintWindow, filename)
	if err != nil {
		fmt.Printf("Export error: %s", err)
	}
}
func oneHotEncodingCheckBoxFunction(b bool) {
	Options.OneHotEncodingSave = b
	if b {
		flatMatrixCheck.Disable()
		flatMatrixCheck.SetChecked(false)
		matlabSaveCheck.SetChecked(true)
	} else {
		flatMatrixCheck.Enable()
	}

}
func addButtonFunction() {
	if !Options.SettingsSaved {
		dialog.ShowError(fmt.Errorf("please first save settings"), Application.mainWindow)
		return
	}
	if input.Text != "" {
		if Options.MatlabSaveFormat {
			AddToFileForMatlab(Application.paintObject.GetMatrix(Application.paintWindow), input.Text)
			count, err := strconv.Atoi(counterLabel.Text)
			if err != nil {
				return
			}
			counterLabel.SetText(strconv.Itoa(count + 1))
			addLabelAnimation(statusLabel)
			statusLabel.Text = "Added!"
			return
		}
		err := AddToFile(Application.paintObject.GetMatrix(Application.paintWindow), input.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("error to add matrix"), Application.mainWindow)
			return
		}
		count, err := strconv.Atoi(counterLabel.Text)
		if err != nil {
			return
		}
		counterLabel.SetText(strconv.Itoa(count + 1))
		addLabelAnimation(statusLabel)
		statusLabel.Text = "Added!"
		return
	}
	dialog.ShowError(fmt.Errorf("please enter valid label"), Application.mainWindow)
}
func applyProjectSetting() {
	rowInput.Disable()
	colInput.Disable()
	flatMatrixCheck.Disable()
	matlabSaveCheck.Disable()
	oneHotEncodingSaveCheck.Disable()
	Options.SettingsSaved = true
	InitializeTemps()
}
func resetProjectSetting() {
	dialog.ShowConfirm("Warning", "Are you sure you want to do that?\nthis is delete your added matrix if you dont saves it. ",
		func(choice bool) {
			rowInput.Enable()
			colInput.Enable()
			flatMatrixCheck.Enable()
			matlabSaveCheck.Enable()
			oneHotEncodingSaveCheck.Enable()
			counterLabel.SetText("0")
			Options.SettingsSaved = false
			TempData.TempTarget = nil
			TempData.TempMatrix = nil
			OneHotDictionary.Dictionary = nil
			OneHotDictionary.Values = nil
			if &TempData.buffer != nil {
				TempData.buffer = bytes.Buffer{}
				TempData.buffer.Reset()
			}

		}, Application.mainWindow,
	)
}
func labelValidator(s string) error {
	if len(s) > 20 {
		return fmt.Errorf("label too long")
	}
	return nil
}
func rowValidator(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil || val <= 0 {
		return fmt.Errorf("enter number ")
	}
	Options.MatrixRow = val + 1
	return nil
}
func colValidator(s string) error {
	val, err := strconv.Atoi(s)
	if err != nil || val <= 0 {
		return fmt.Errorf("enter number")
	}
	Options.MatrixCol = val + 1
	return nil
}
func onStartedApplication() {
	// Temporarily disable stdout to prevent matrix printing
	oldStdOut := os.Stdout
	os.Stdout = nil
	Application.paintObject.PrintMatrix(Application.paintWindow, Options.FlatMatrix)
	os.Stdout = oldStdOut
}

func prepareSaveProjectObj() {
	SavedProject.Options = Options
	SavedProject.TempData = TempData
	SavedProject.OneHotDictionary = OneHotDictionary
	SavedProject.CounterValue = counterLabel.Text
	SavedProject.Buffer = TempData.buffer.Bytes()

}

func loadProjectFile(reader io.ReadCloser) error {
	decoder := gob.NewDecoder(reader)
	err := decoder.Decode(&SavedProject)
	if err != nil {
		return err
	}
	Options = SavedProject.Options
	TempData = SavedProject.TempData
	TempData.buffer.Write(SavedProject.Buffer)
	OneHotDictionary = SavedProject.OneHotDictionary
	err = countValue.Set(SavedProject.CounterValue)
	if err != nil {
		log.Println(err)
		return err
	}
	rowInput.Text = strconv.Itoa(Options.MatrixRow - 1)
	colInput.Text = strconv.Itoa(Options.MatrixCol - 1)
	oneHotEncodingSaveCheck.SetChecked(Options.OneHotEncodingSave)
	matlabSaveCheck.SetChecked(Options.MatlabSaveFormat)
	flatMatrixCheck.SetChecked(Options.FlatMatrix)
	Application.mainWindow.Content().Refresh()
	return nil
}

func saveProjectFileFunction() {
	dialog.ShowFileSave(func(writer fyne.URIWriteCloser, err error) {
		if !Options.SettingsSaved {
			dialog.ShowError(fmt.Errorf("Please first save project settings."), Application.mainWindow)
			return
		}
		prepareSaveProjectObj()
		encoder := gob.NewEncoder(writer)
		err = encoder.Encode(SavedProject)
		if err != nil {
			log.Println(err)
			return
		}

		err = writer.Close()
		if err != nil {
			return
		}

	}, Application.mainWindow)

}

func loadProjectFileFunction() {
	dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {

		if Options.SettingsSaved {
			dialog.ShowConfirm("Warning", "Are you sure to load project? Your current session is removed if you dont saved it.", func(b bool) {
				if b {
					loadErr := loadProjectFile(reader)
					if loadErr != nil {
						dialog.ShowError(fmt.Errorf("error loading project file"), Application.mainWindow)
						return
					}
					applyProjectSetting()
				}

			}, Application.mainWindow)
		} else {
			loadErr := loadProjectFile(reader)
			if loadErr != nil {
				dialog.ShowError(fmt.Errorf("error loading project file"), Application.mainWindow)
				return
			}
			applyProjectSetting()
			statusLabel.Text = "Project loaded!"
			addLabelAnimation(statusLabel)
		}

	}, Application.mainWindow)
}
