package main

import (
	"errors"
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

func saveOperation() {
	if !Options.SettingsSaved {
		dialog.ShowError(fmt.Errorf("please first save settings"), Application.mainWindow)
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
		if err = SaveFile(path, dataFileName+".csv"); err != nil {
			statusLabel.SetText("Not Saved!")
			return
		}
		statusLabel.SetText("Saved!")
	} else {
		// File exists, ask for confirmation
		dialog.ShowConfirm("Warning", "File exists. Do you want to replace it?", func(b bool) {
			if b {
				if err = SaveFile(path, "data.csv"); err != nil {
					statusLabel.SetText("Not Saved!")
					return
				}
				statusLabel.SetText("Saved!")
			} else {
				statusLabel.SetText("Not Saved!")
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
func expertOperation() {
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
	Options.oneHotEncodingSave = b
	if b {
		flatMatrixCheck.Disable()
		flatMatrixCheck.SetChecked(false)
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
			err := AddToFileForMatlab(Application.paintObject.GetMatrix(Application.paintWindow), input.Text)
			if err != nil {
				dialog.ShowError(fmt.Errorf("error to add matrix"), Application.mainWindow)
				return
			}
			statusLabel.SetText("Added!")
			return
		}
		err := AddToFile(Application.paintObject.GetMatrix(Application.paintWindow), input.Text)
		if err != nil {
			dialog.ShowError(fmt.Errorf("error to add matrix"), Application.mainWindow)
			return
		}
		statusLabel.SetText("Added!")
		return
	}
	dialog.ShowError(fmt.Errorf("please enter valid label"), Application.mainWindow)
}
func saveProjectButtonFunction() {
	rowInput.Disable()
	colInput.Disable()
	flatMatrixCheck.Disable()
	matlabSaveCheck.Disable()
	oneHotEncodingSaveCheck.Disable()
	Options.SettingsSaved = true
	InitializeTemps(Options.MatlabSaveFormat)
}
func resetProjectButtonFunction() {
	dialog.ShowConfirm("Warning", "Are you sure you want to do that?\nthis is delete your added matrix if you dont saves it. ",
		func(choice bool) {
			rowInput.Enable()
			colInput.Enable()
			flatMatrixCheck.Enable()
			matlabSaveCheck.Enable()
			oneHotEncodingSaveCheck.Enable()
			Options.SettingsSaved = false
			if TempData.file != nil {
				err := os.Remove(TempData.file.Name())
				if err != nil {
					fmt.Println(err)
				}
			}
			if TempData.targetFile != nil {
				err := os.Remove(TempData.targetFile.Name())
				if err != nil {
					fmt.Println(err)
				}
			}
			err := os.RemoveAll(TempData.dir)
			if err != nil {
				fmt.Println(err)
			}
			TempData.file.Close()

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
func onStoppedApplication() {
	if TempData.file != nil {
		err := TempData.file.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}
	if TempData.targetFile != nil {
		err := TempData.targetFile.Close()
		if err != nil {
			log.Println(err)
			return
		}
	}
	if err := os.RemoveAll(TempData.dir); err != nil {
		log.Println("Error removing temp directory:", err)
	}
}
