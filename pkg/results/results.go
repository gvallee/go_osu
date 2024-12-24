//
// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package results

import (
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"strings"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/gvallee/go_util/pkg/notation"
)

type DataPoint struct {
	Size  float64
	Value float64
}

type Result struct {
	DataPoints []*DataPoint
}

type Results struct {
	Result []*Result
}

// ExtractDataFromOutput returns the OSU data from a output file. The first
// array returned represents all the data sizes, while the second returned array
// represents the value for the associated size (the two arrays are assumed to
// be ordered).
func ExtractDataFromOutput(benchmarkOutput []string) ([]float64, []float64, error) {
	var x []float64
	var y []float64

	var val1 float64
	var val2 float64

	var err error

	save := false
	stop := false

	for _, line := range benchmarkOutput {
		val1 = -1.0
		val2 = -1.0

		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			if !save && strings.HasPrefix(line, "# OSU") {
				save = true
			}
			continue
		}
		if !save {
			// We skip whatever is at the beginning of the file until we reach the OSU header
			continue
		}
		if strings.Contains(line, "more processes have sent help message") {
			// Open MPI throwing warnings for whatever reason, skipping
			continue
		}

		// We replace all double spaces with a single space to make it easier to identify the real data
		tokens := strings.Split(line, " ")
		for _, t := range tokens {
			if t == " " || t == "" {
				continue
			}

			if val1 == -1.0 {
				val1, err = strconv.ParseFloat(t, 64)
				if err != nil {
					if len(x) == 0 {
						return nil, nil, fmt.Errorf("unable to convert %s (from %s): %w", t, line, err)
					} else {
						log.Printf("stop parsing, unable to convert %s (from %s): %s", t, line, err)
						stop = true
						break
					}
				}
				x = append(x, val1)
			} else {
				val2, err = strconv.ParseFloat(t, 64)
				if err != nil {
					if len(y) == 0 {
						return nil, nil, fmt.Errorf("unable to convert %s (from %s): %w", t, line, err)
					} else {
						log.Printf("stop parsing, unable to convert %s (from %s): %s", t, line, err)
						stop = true
						break
					}
				}
				y = append(y, val2)
				break // todo: we need to extend this for sub-benchmarks returning more than one value (see what is done in OpenHPCA)
			}
		}

		if stop {
			break
		}
	}

	return x, y, nil
}

func ParseOutputFile(path string) (*Result, error) {
	log.Printf("Parsing result file %s", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	str := string(content)
	dataSize, data, err := ExtractDataFromOutput(strings.Split(str, "\n"))
	if err != nil {
		return nil, fmt.Errorf("unable to parse %s: %w", path, err)
	}
	if len(dataSize) != len(data) {
		return nil, fmt.Errorf("unsupported data format (%s), skipping: %d different sizes with %d values", path, len(dataSize), len(data))
	}

	newResult := new(Result)
	for i := 0; i < len(dataSize); i++ {
		newDataPoint := new(DataPoint)
		newDataPoint.Size = dataSize[i]
		newDataPoint.Value = data[i]
		newResult.DataPoints = append(newResult.DataPoints, newDataPoint)
	}

	return newResult, nil
}

func GetResultsFromFiles(listFiles []string) (*Results, error) {
	res := new(Results)
	for _, file := range listFiles {
		newResult, err := ParseOutputFile(file)
		if err != nil {
			return nil, err
		}
		res.Result = append(res.Result, newResult)
	}

	return res, nil
}

func addValuesToExcel(excelFile *excelize.File, sheetID string, lineStart int, col int, datapoints []*DataPoint) error {
	colID := notation.IntToAA(col)
	lineID := lineStart
	for _, d := range datapoints {
		// Find the correct line where to put the data
		for {
			dataSizeStr := excelFile.GetCellValue(sheetID, fmt.Sprintf("A%d", lineID))
			dataSize, err := strconv.ParseFloat(dataSizeStr, 64)
			if err != nil {
				return fmt.Errorf("unable to parse %s: %w", dataSizeStr, err)
			}
			if dataSize == d.Size {
				break
			}
			lineID++
		}
		excelFile.SetCellValue(sheetID, fmt.Sprintf("%s%d", colID, lineID), d.Value)
		lineID++
	}
	return nil
}

func Excelize(excelFilePath string, results *Results) error {
	excelFile := excelize.NewFile()

	// Add the message sizes into the first column
	lineID := 1 // 1-indexed to match Excel semantics
	for _, dp := range results.Result[0].DataPoints {
		excelFile.SetCellValue("Sheet1", fmt.Sprintf("A%d", lineID), dp.Size)
		lineID++
	}

	// Add the values
	col := 1   // 0-indexed so it can be used with IntToAA
	lineID = 1 // 1-indexed to match Excel semantics
	for _, d := range results.Result {
		err := addValuesToExcel(excelFile, "Sheet1", lineID, col, d.DataPoints)
		if err != nil {
			return fmt.Errorf("addValuesToExcel() failed: %w", err)
		}
		col++
	}

	err := excelFile.SaveAs(excelFilePath)
	if err != nil {
		return err
	}

	return nil
}

// Excelize create a MSExcel spreadsheet with all the data passed in, as well as labels associated to the data.
// The order of the labels is assumed to be the same than the order of the data. The data is included on the sheet
// of index sheetStart (1-indexed).
func ExcelizeWithLabels(sheetStart int, results *Results, labels []string) (*excelize.File, error) {
	if sheetStart <= 0 {
		return nil, fmt.Errorf("invalid sheet start index (must be > 0): %d", sheetStart)
	}
	if results == nil {
		return nil, fmt.Errorf("undefined results")
	}
	if len(results.Result) == 0 {
		return nil, fmt.Errorf("empty result dataset")
	}
	excelFile := excelize.NewFile()
	sheetID := fmt.Sprintf("Sheet%d", sheetStart)

	if sheetID != "Sheet1" {
		excelFile.NewSheet(sheetID)
	}

	// Add the labels
	lineID := 1 // 1-indexed to match Excel semantics
	col := 1    // 0-indexed so it can be used with IntToAA
	for _, label := range labels {
		excelFile.SetCellValue(sheetID, fmt.Sprintf("%s%d", notation.IntToAA(col), lineID), label)
		col++
	}

	// Add the message sizes into the first column
	lineID = 2 // 1-indexed
	for _, dp := range results.Result[0].DataPoints {
		excelFile.SetCellValue(sheetID, fmt.Sprintf("A%d", lineID), dp.Size)
		lineID++
	}

	// Add the values
	col = 1    // 0-indexed so it can be used with IntToAA
	lineID = 2 // 1-indexed
	for _, d := range results.Result {
		err := addValuesToExcel(excelFile, sheetID, lineID, col, d.DataPoints)
		if err != nil {
			return nil, err
		}
		col++
	}

	return excelFile, nil
}
