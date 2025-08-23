// Package main provides functionality for matrix operations and data handling in Draw2Matrix
package main

import (
	"bytes"
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

var OneHotDictionary struct {
	Dictionary map[string]interface{}
	Values     []string
}

const (
	// RowFlat indicates row-wise flattening of matrix
	RowFlat FlatDirection = iota
	// ColFlat indicates column-wise flattening of matrix
	ColFlat
)

// TempData stores temporary files and data during program execution
var TempData struct {
	Saved      bool // Flag indicating if data has been Saved
	buffer     bytes.Buffer
	TempMatrix [][]int8 // Temporary storage for matrix data
	TempTarget []string // Temporary storage for matrix label
}

// InitializeTemps creates temporary files and directories for data storage
// If forMatlab is true, additional files for MATLAB format will be created
func InitializeTemps() {
	TempData.Saved = false
	if Options.OneHotEncodingSave {
		OneHotDictionary.Dictionary = map[string]interface{}{}
		OneHotDictionary.Values = []string{}
	}
	if Options.MatlabSaveFormat {
		TempData.TempTarget = make([]string, 0)
		TempData.TempMatrix = make([][]int8, 0)

	}
	TempData.buffer = bytes.Buffer{}
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

func oneHotEncoder(label string) []int8 {
	result := make([]int8, len(OneHotDictionary.Values))
	for i, value := range OneHotDictionary.Values {
		if value == label {
			result[i] = 1
			continue
		}
		result[i] = 0
	}
	return result
}

func oneHotSaveUtil() string {
	tempResult := make([][]int8, 0)
	for _, value := range TempData.TempTarget {
		tempResult = append(tempResult, oneHotEncoder(value))
	}
	result := transposeMatrix(tempResult)
	return processForMatlabString(result)
}

// AddToFileForMatlab appends matrix data and its corresponding output
// to temporary files in MATLAB format
func AddToFileForMatlab(inputData [][]int8, outputData string) {

	// Store flattened matrix data
	TempData.TempMatrix = append(TempData.TempMatrix, ToFlattenMatrix(inputData))
	TempData.TempTarget = append(TempData.TempTarget, outputData)
	if Options.OneHotEncodingSave {
		if _, ok := OneHotDictionary.Dictionary[outputData]; !ok {
			OneHotDictionary.Dictionary[outputData] = true
			OneHotDictionary.Values = append(OneHotDictionary.Values, outputData)
		}
	}
}

// SaveFileForMatlab saves the matrix data and target data to separate files
// in MATLAB compatible format
func SaveFileForMatlab(dirPath, dataFileName, targetFileName string) error {
	extension := ".txt"
	if Options.DotMFileWithVariable {
		extension = ".m"
	}
	// Prepare file paths
	dataPath := filepath.Join(dirPath, dataFileName+extension)
	targetPath := filepath.Join(dirPath, targetFileName+extension)

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

	// Write processed matrix data
	finalData := processForMatlabString(transposeMatrix(TempData.TempMatrix))
	if Options.DotMFileWithVariable {
		finalData = dataFileName + "_variable = " + finalData + ";"
	}
	if _, err = dataFile.WriteString(finalData); err != nil {
		return err
	}

	if Options.OneHotEncodingSave {
		finalTarget := oneHotSaveUtil()
		if Options.DotMFileWithVariable {
			finalTarget = targetFileName + "_variable = " + finalTarget + ";"
		}
		_, err = targetFile.WriteString(finalTarget)
		if err != nil {
			log.Println(err)
			return err
		}
	} else {
		finalTarget := fmt.Sprintf("%v", TempData.TempTarget)
		if Options.DotMFileWithVariable {
			finalTarget = targetFileName + " = " + finalTarget + ";"
		}
		_, err = targetFile.WriteString(finalTarget)
		if err != nil {
			log.Println(err)
			return err
		}
	}

	return nil
}

// AddToFile appends matrix data and its corresponding output to a CSV file
// If FlatMatrix option is enabled, the matrix will be flattened before writing
func AddToFile(inputData [][]int8, outputData string) error {
	csvWriter := csv.NewWriter(&TempData.buffer)
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

	if _, err = io.Copy(file, &TempData.buffer); err != nil {
		return err
	}

	TempData.Saved = true
	return nil
}
