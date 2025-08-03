package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

var TempData struct {
	file  *os.File
	dir   string
	saved bool
}

func ConvertToFlattenMatric(matric [][]int8) []int8 {
	result := make([]int8, 0)
	for _, i := range matric {
		for _, j := range i {
			result = append(result, j)
		}
	}
	return result
}

func InitializeTemps() {
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
	TempData.saved = false
}

func AddToFile(inputData [][]int8, outputData string) error {
	inp := csv.NewWriter(TempData.file)
	defer inp.Flush()
	inp.UseCRLF = true
	err := inp.Write([]string{"Input", "Target"})
	if err != nil {
		return err
	}
	if Options.FlatMatrix {
		err := inp.Write([]string{fmt.Sprintf("%d", ConvertToFlattenMatric(inputData)), outputData})
		if err != nil {
			return err
		}
		return nil
	}
	err = inp.Write([]string{fmt.Sprintf("%d", inputData), outputData})
	if err != nil {
		return err
	}
	return nil
}

func SaveFile(DirPath string) error {
	path := filepath.Join(DirPath, "data.csv")
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	tempFile, err := os.OpenFile(TempData.file.Name(), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	_, err = io.Copy(file, tempFile)
	if err != nil {
		return err
	}

	TempData.saved = true

	defer file.Close()
	defer tempFile.Close()
	return nil
}
