package main

import (
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

var (
	countValue       = binding.NewString()
	statusLabel      = canvas.NewText("start", color.Black)
	counterLabelText = widget.NewLabel("count: ")
	counterLabel     = widget.NewLabelWithData(countValue)
	refreshBtn       = widget.NewButtonWithIcon("Clear Paint", theme.DeleteIcon(), func() {
		Application.paintObject.Clear()
	})
	savePath        = widget.NewEntry()
	dataFileEntry   = widget.NewEntry()
	targetFileEntry = widget.NewEntry()
	openPaint       = widget.NewButtonWithIcon("OpenPaint", theme.WindowMaximizeIcon(), openPaintWindowOperation)
	changePath      = widget.NewButtonWithIcon("Browse", theme.FolderIcon(), browseOperation)
	saveBtn         = widget.NewButtonWithIcon("Save file", theme.DocumentSaveIcon(), exportFileOperation)
	input           = widget.NewEntry()
	exportBtn       = widget.NewButtonWithIcon("Export PNG", theme.FileImageIcon(), expertPNGOperation)
	flatMatrixCheck = widget.NewCheck("Flat Matrix", func(b bool) {
		Options.FlatMatrix = b
	})
	matlabSaveCheck           = widget.NewCheck("Matlab Save Format", matlabSaveCheckBoxFunction)
	dotMFileWithVariableCheck = widget.NewCheck(".m file save", DotMFileWithVariableCheck)
	oneHotEncodingSaveCheck   = widget.NewCheck("One Hot Encoding Save", oneHotEncodingCheckBoxFunction)
	colInput                  = widget.NewEntry()
	rowInput                  = widget.NewEntry()
	addBtn                    = widget.NewButtonWithIcon("Add", theme.ContentAddIcon(), addButtonFunction)
	addAndClearPaintBtn       = widget.NewButtonWithIcon("Add & Clear Paint", theme.ContentCutIcon(), func() {
		addButtonFunction()
		Application.paintObject.Clear()
	})
	saveOptionsBtn = widget.NewButtonWithIcon("Save Settings", theme.SettingsIcon(), func() {
		applyProjectSetting(true)
	})
	resetProjectBtn = widget.NewButtonWithIcon("Reset Project", theme.ContentClearIcon(), resetProjectSetting)
	toolbar         = widget.NewToolbar(
		widget.NewToolbarAction(theme.DocumentSaveIcon(), saveProjectFileFunction),
		widget.NewToolbarAction(theme.ContentUndoIcon(), loadProjectFileFunction),
		widget.NewToolbarAction(theme.InfoIcon(), aboutBtn))
)

// Layout containers
var (
	settingsContainer = container.NewVBox(
		openPaint,
		widget.NewLabel("Matrix Settings:"),
		container.NewGridWithColumns(2, rowInput, colInput),
		container.NewGridWithColumns(4, flatMatrixCheck, matlabSaveCheck, dotMFileWithVariableCheck, oneHotEncodingSaveCheck),
		container.NewGridWithColumns(2, resetProjectBtn, saveOptionsBtn),
	)

	pathContainer = container.NewVBox(
		container.NewGridWithColumns(2, savePath, changePath),
		dataFileEntry,
		targetFileEntry,
	)
	actionContainer = container.NewVBox(
		widget.NewLabel("Actions:"),
		container.NewGridWithColumns(2, refreshBtn, exportBtn),
		pathContainer, saveBtn,
	)
	statusContainer = container.NewHBox(
		container.NewPadded(container.NewGridWithColumns(2, counterLabelText, counterLabel)),
		layout.NewSpacer(),
		container.NewPadded(statusLabel),
	)

	labelContainer = container.NewVBox(
		widget.NewLabel("Label:"),
		container.NewBorder(nil, statusContainer, addBtn, addAndClearPaintBtn, input),
	)

	bottomContainer = container.NewVBox(
		container.NewPadded(toolbar),
		container.NewPadded(settingsContainer),
		container.NewPadded(actionContainer),
		container.NewPadded(labelContainer),
	)
)
