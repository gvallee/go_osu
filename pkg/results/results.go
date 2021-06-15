//
// Copyright (c) 2021, NVIDIA CORPORATION. All rights reserved.
//
// See LICENSE.txt for license information
//

package results

import (
	"fmt"
	"io/ioutil"
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

func ExtractDataFromOutput(benchmarkOutput []string) ([]float64, []float64, error) {
	var x []float64
	var y []float64

	var val1 float64
	var val2 float64

	var err error

	save := false

	for _, line := range benchmarkOutput {
		val1 = -1.0
		val2 = -1.0

		if line == "" {
			continue
		}
		if strings.HasPrefix(line, "#") {
			if !save {
				save = true
			}
			continue
		}
		if !save {
			// We skip whatever is at the beginning of the file until we reach the OSU header
			continue
		}

		tokens := strings.Split(line, " ")
		for _, t := range tokens {
			if t == " " || t == "" {
				continue
			}

			if val1 == -1.0 {
				val1, err = strconv.ParseFloat(t, 64)
				if err != nil {
					return nil, nil, err
				}
				x = append(x, val1)
			} else {
				val2, err = strconv.ParseFloat(t, 64)
				if err != nil {
					return nil, nil, err
				}
				y = append(y, val2)
			}
		}
	}

	return x, y, nil
}

func rawDataToResults(sizes []float64, values []float64, res *Results) error {
	newResult := new(Result)
	for i := 0; i < len(sizes); i++ {
		newDataPoint := new(DataPoint)
		newDataPoint.Size = sizes[i]
		newDataPoint.Value = values[i]
		newResult.DataPoints = append(newResult.DataPoints, newDataPoint)
	}
	res.Result = append(res.Result, newResult)
	return nil
}

func GetResultsFromFiles(listFiles []string) (*Results, error) {
	res := new(Results)
	for _, file := range listFiles {
		content, err := ioutil.ReadFile(file)
		if err != nil {
			return nil, err
		}
		str := string(content)
		dataSize, data, err := ExtractDataFromOutput(strings.Split(str, "\n"))
		if err != nil {
			return nil, err
		}
		if len(dataSize) != len(data) {
			return nil, fmt.Errorf("inconsistent data")
		}
		err = rawDataToResults(dataSize, data, res)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

func addValuesToExcel(excelFile *excelize.File, lineStart int, col int, datapoints []*DataPoint) error {
	colID := notation.IntToAA(col)
	lineID := lineStart
	for _, d := range datapoints {
		excelFile.SetCellValue("Sheet1", fmt.Sprintf("%s%d", colID, lineID), d.Value)
		lineID++
	}
	return nil
}

func Excelize(excelFilePath string, results *Results) error {
	excelFile := excelize.NewFile()

	// Add the message sizes into the first column
	lineID := 1
	for _, dp := range results.Result[0].DataPoints {
		excelFile.SetCellValue("Sheet1", fmt.Sprintf("A%d", lineID), dp.Size)
		lineID++
	}

	// Add the values
	col := 1 // 1-indexed
	for _, d := range results.Result {
		err := addValuesToExcel(excelFile, 1, col, d.DataPoints)
		if err != nil {
			return err
		}
		col++
	}

	err := excelFile.SaveAs(excelFilePath)
	if err != nil {
		return err
	}

	return nil
}

func ExcelizeWithLabels(excelFilePath string, results *Results, labels []string) error {
	excelFile := excelize.NewFile()

	// Add the labels
	lineID := 1
	col := 2
	for _, label := range labels {
		excelFile.SetCellValue("Sheet1", fmt.Sprintf("%s%d", notation.IntToAA(col), lineID), label)
		col++
	}

	// Add the message sizes into the first column
	lineID = 2 // 1-indexed
	for _, dp := range results.Result[0].DataPoints {
		excelFile.SetCellValue("Sheet1", fmt.Sprintf("A%d", lineID), dp.Size)
		lineID++
	}

	// Add the values
	col = 2    // 1-indexed
	lineID = 2 // 1-indexed
	for _, d := range results.Result {
		err := addValuesToExcel(excelFile, lineID, col, d.DataPoints)
		if err != nil {
			return err
		}
		col++
	}

	err := excelFile.SaveAs(excelFilePath)
	if err != nil {
		return err
	}

	return nil
}
