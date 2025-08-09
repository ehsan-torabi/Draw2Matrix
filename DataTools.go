package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
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
	tempMatric [][]int8
}

func transposeMatric(matrix [][]int8) [][]int8 {
	if len(matrix) == 0 {
		return [][]int8{}
	}

	rows := len(matrix)
	cols := len(matrix[0])
	transposed := make([][]int8, cols)
	for i := range transposed {
		transposed[i] = make([]int8, rows)
	}

	for i := 0; i < rows; i++ {
		for j := 0; j < cols; j++ {
			transposed[j][i] = matrix[i][j]
		}
	}

	return transposed
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
		result := fmt.Sprintf("%v", rowFlatten)
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

func addSemicolon(file *os.File) string {
	result := ""
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		result += line + ";\n"
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	before, found := strings.CutSuffix(result, ";\n")
	if !found {
		return before
	}
	return result
}

func processForMatlabString(matrix [][]int8) string {
	var result strings.Builder
	result.Grow(len(matrix))
	result.WriteString("[ ")
	for i, row := range matrix {
		for _, element := range row {
			result.WriteString(fmt.Sprintf("%d ", element))
		}
		if i < len(matrix)-1 {
			result.WriteString(";\n")
		}
	}
	result.WriteString("]")
	return result.String()
}

func AddToFileForMatlab(inputData [][]int8, outputData string) error {
	outputFile := TempData.targetFile
	var tempByte = make([]byte, 1)
	n, err := outputFile.ReadAt(tempByte, 0)
	if err != nil && err != io.EOF {
		return err
	}
	if n == 0 || err == io.EOF {
		_, err = outputFile.WriteString("[ ")
		if err != nil {
			return err
		}
	}
	TempData.tempMatric = append(TempData.tempMatric, ToFlattenMatric(inputData))
	fmt.Println(TempData.tempMatric)
	target := outputData + " "
	if err != nil {
		return err
	}
	_, err = outputFile.WriteString(target)
	if err != nil {
		return err
	}
	return nil

}

func SaveFileForMatlab(DirPath, dataFileName, targetFileName string) error {
	dataPath := filepath.Join(DirPath, dataFileName+".txt")
	targetPath := filepath.Join(DirPath, targetFileName+".txt")
	dataFile, err := os.OpenFile(dataPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	defer dataFile.Close()
	if err != nil {
		return err
	}
	targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	defer targetFile.Close()
	if err != nil {
		return err
	}
	tempTargetFile, err := os.OpenFile(TempData.targetFile.Name(), os.O_APPEND|os.O_RDONLY, 0600)
	defer tempTargetFile.Close()
	if err != nil {
		return err
	}

	finalData := processForMatlabString(transposeMatric(TempData.tempMatric))
	_, err = dataFile.WriteString(finalData)
	if err != nil {
		return err
	}
	_, err = targetFile.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}
	scanner := bufio.NewScanner(tempTargetFile)
	var line string
	for scanner.Scan() {
		line += scanner.Text()
	}
	_, err = targetFile.WriteString(line + "]")
	if err != nil {
		return err
	}
	return nil

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
	if !Options.MatlabSaveFormat {
		_, err = file.WriteString("Input,Target\n")
		if err != nil {
			return err
		}
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
