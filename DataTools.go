// Package main provides functionality for matrix operations and data handling in Draw2Matrix
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

// FlatDirection represents the direction for flattening a matrix
type FlatDirection int8

const (
	// RowFlat indicates row-wise flattening of matrix
	RowFlat FlatDirection = iota
	// ColFlat indicates column-wise flattening of matrix
	ColFlat
)

// TempData stores temporary files and data during program execution
var TempData struct {
	file       *os.File // Temporary file for storing matrix data
	targetFile *os.File // Temporary file for storing target data (used in MATLAB format)
	dir        string   // Directory path for temporary files
	saved      bool     // Flag indicating if data has been saved
	tempMatrix [][]int8 // Temporary storage for matrix data
}

// transposeMatrix converts a matrix to its transpose form
func transposeMatrix(matrix [][]int8) [][]int8 {
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

// ToFlattenMatrix converts a 2D matrix into a 1D slice
func ToFlattenMatrix(matrix [][]int8) []int8 {
	result := make([]int8, 0)
	for _, row := range matrix {
		for _, val := range row {
			result = append(result, val)
		}
	}
	return result
}

// ToFlattenMatrixString converts a 2D matrix to a string representation
// based on the specified flattening direction (row-wise or column-wise)
func ToFlattenMatrixString(matrix [][]int8, direction FlatDirection) string {
	flattenedMatrix := ToFlattenMatrix(matrix)
	if direction == RowFlat {
		return fmt.Sprintf("%v", flattenedMatrix)
	} else if direction == ColFlat {
		var elements []string
		for _, element := range flattenedMatrix {
			elements = append(elements, fmt.Sprintf("%d", element))
		}
		return strings.Join(elements, "\n")
	}
	return ""
}

// InitializeTemps creates temporary files and directories for data storage
// If forMatlab is true, additional files for MATLAB format will be created
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

// AddToFileForMatlab appends matrix data and its corresponding output
// to temporary files in MATLAB format
func AddToFileForMatlab(inputData [][]int8, outputData string) error {
	outputFile := TempData.targetFile
	var tempByte = make([]byte, 1)

	// Check if file is empty
	n, err := outputFile.ReadAt(tempByte, 0)
	if err != nil && err != io.EOF {
		return err
	}

	// Initialize file with opening bracket if empty
	if n == 0 || err == io.EOF {
		if _, err = outputFile.WriteString("[ "); err != nil {
			return err
		}
	}

	// Store flattened matrix data
	TempData.tempMatrix = append(TempData.tempMatrix, ToFlattenMatrix(inputData))
	fmt.Println(TempData.tempMatrix)

	// Write target data with space
	target := outputData + " "
	if _, err = outputFile.WriteString(target); err != nil {
		return err
	}

	return nil
}

// SaveFileForMatlab saves the matrix data and target data to separate files
// in MATLAB compatible format
func SaveFileForMatlab(dirPath, dataFileName, targetFileName string) error {
	// Prepare file paths
	dataPath := filepath.Join(dirPath, dataFileName+".txt")
	targetPath := filepath.Join(dirPath, targetFileName+".txt")

	// Create data file
	dataFile, err := os.OpenFile(dataPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer dataFile.Close()

	// Create target file
	targetFile, err := os.OpenFile(targetPath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer targetFile.Close()

	// Open temporary target file
	tempTargetFile, err := os.OpenFile(TempData.targetFile.Name(), os.O_APPEND|os.O_RDONLY, 0600)
	if err != nil {
		return err
	}
	defer tempTargetFile.Close()

	// Write processed matrix data
	finalData := processForMatlabString(transposeMatrix(TempData.tempMatrix))
	if _, err = dataFile.WriteString(finalData); err != nil {
		return err
	}

	// Write target data with closing bracket
	if _, err = targetFile.Seek(0, io.SeekEnd); err != nil {
		return err
	}

	scanner := bufio.NewScanner(tempTargetFile)
	var line string
	for scanner.Scan() {
		line += scanner.Text()
	}

	if _, err = targetFile.WriteString(line + "]"); err != nil {
		return err
	}

	return nil
}

// AddToFile appends matrix data and its corresponding output to a CSV file
// If FlatMatrix option is enabled, the matrix will be flattened before writing
func AddToFile(inputData [][]int8, outputData string) error {
	csvWriter := csv.NewWriter(TempData.file)
	defer csvWriter.Flush()
	csvWriter.UseCRLF = true

	var dataString string
	if Options.FlatMatrix {
		dataString = ToFlattenMatrixString(inputData, RowFlat)
	} else {
		dataString = fmt.Sprintf("%d", inputData)
	}

	if err := csvWriter.Write([]string{dataString, outputData}); err != nil {
		return err
	}
	return nil
}

// SaveFile saves the accumulated data to a final file
// For non-MATLAB format, it includes a header row
func SaveFile(dirPath, filename string) error {
	// Create the final file
	path := filepath.Join(dirPath, filename)
	file, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write header for non-MATLAB format
	if !Options.MatlabSaveFormat {
		if _, err = file.WriteString("Input,Target\n"); err != nil {
			return err
		}
	}

	// Copy temporary file contents to final file
	tempFile, err := os.OpenFile(TempData.file.Name(), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return err
	}
	defer tempFile.Close()

	if _, err = io.Copy(file, tempFile); err != nil {
		return err
	}

	TempData.saved = true
	return nil
}
