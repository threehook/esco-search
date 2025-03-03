package csv

import (
	"encoding/csv"
	"fmt"
	"os"
	"slices"
)

func ReadCSV(filePath string, columns ...string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	header := rows[0]
	table := make([][]string, len(rows)-1)

	for i, row := range rows[1:] {
		var content []string
		for j, value := range row {
			if columns != nil &&
				len(columns) > 0 &&
				!slices.Contains(columns, header[j]) {
				continue
			}
			content = append(content, value)
			table[i] = content
		}
	}

	return table, nil
}

//func MergeCSVs(csvRows [][]string, linkColumnIndexFile1 int, file2 string, linkColumnIndexFile2 int) ([][]string, error) {
//	// Load second CSV into a map
//	file2Data, err := readCSVToMap(file2, linkColumnIndexFile2)
//	if err != nil {
//		return nil, fmt.Errorf("Error reading file2: %w")
//	}
//
//	// Open and process the second CSV file
//	joinedCSVs, err := joinCSVs(csvRows, file2Data, linkColumnIndexFile1)
//	if err != nil {
//		return nil, fmt.Errorf("Error joining CSVs: %w")
//	}
//
//	return joinedCSVs, nil
//}

func ReadCSVToMap(filePath string, keyColumn int, multi bool) (map[string][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	rows, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	// Remove the headers
	rows = slices.Delete(rows, 0, 1)

	dataMap := make(map[string][]string)
	for _, row := range rows {
		if keyColumn >= len(row) {
			return nil, fmt.Errorf("keyColumn index out of range for row: %v", row)
		}
		key := row[keyColumn]
		dataMap[key] = row
	}

	return dataMap, nil
}

//func joinCSVs(csvRows [][]string, mapData map[string][]string, linkColumnIndex int) ([][]string, error) {
//	joinedRows := make([][]string, len(csvRows))
//	for i, row := range csvRows {
//		if linkColumnIndex >= len(row) {
//			return nil, fmt.Errorf("linkColumnIndex out of range for row: %v", row)
//		}
//		key := row[linkColumnIndex]
//
//		// Check if the key exists in the first file's map
//		if data, exists := mapData[key]; exists {
//			// Combine data from both rows
//			row[linkColumnIndex] = data[3]
//			joinedRows[i] = append(row)
//		}
//	}
//
//	return joinedRows, nil
//}
