package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type FlatDirection int8

const (
	RowFlat FlatDirection = iota
	ColFlat
)

var TempData struct {
	file       *os.File
	targetFile *os.File
	dir        string
	saved      bool
}

func ToFlattenMatric(matric [][]int8) []int8 {
	result := make([]int8, 0)
	for _, i := range matric {
		for _, j := range i {
			result = append(result, j)
		}
	}
	return result
}

func ToFlattenMatricString(matric [][]int8, direction FlatDirection) string {
	rowFlatten := ToFlattenMatric(matric)
	if direction == RowFlat {
		result := fmt.Sprintf("%d", rowFlatten)
		return result
	} else if direction == ColFlat {
		var tempSlice []string
		for _, element := range rowFlatten {
			tempSlice = append(tempSlice, fmt.Sprintf("%d", element))
		}
		return strings.Join(tempSlice, "\n")
	}
	return ""
}

func InitializeTemps(forMatlab bool) {
	var err error
	temp, err := os.MkdirTemp(".", "temp")
	if err != nil {
		return
	}
	TempData.dir = temp
	TempData.file, err = os.CreateTemp(temp, "file")
	if err != nil {
		panic(err)
	}
	if forMatlab {
		TempData.targetFile, err = os.CreateTemp(temp, "target")
	}
	TempData.saved = false
}

func AddToFile(inputData [][]int8, outputData string) error {
	inp := csv.NewWriter(TempData.file)
	defer inp.Flush()
	inp.UseCRLF = true
	if Options.FlatMatrix {
		err := inp.Write([]string{ToFlattenMatricString(inputData, RowFlat), outputData})
		if err != nil {
			return err
		}
		return nil
	}
	err := inp.Write([]string{fmt.Sprintf("%d", inputData), outputData})
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(DirPath, filename string) error {
	path := filepath.Join(DirPath, filename)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	defer file.Close()
	if err != nil {
		return err
	}
	_, err = file.WriteString("Input,Target\n")
	if err != nil {
		return err
	}
	tempFile, err := os.OpenFile(TempData.file.Name(), os.O_RDONLY, os.ModePerm)
	defer tempFile.Close()
	if err != nil {
		return err
	}
	_, err = io.Copy(file, tempFile)
	if err != nil {
		return err
	}
	TempData.saved = true
	return nil
}
