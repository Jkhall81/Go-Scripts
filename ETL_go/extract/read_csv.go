package extract

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// DataSet represents the standardized structure returned from the extract phase.
// We'll reuse this for both CSV and XLSX extractors.
type DataSet struct {
	Headers []string
	Rows    [][]string
	Source  string // file name or path
}

// ReadCSV opens a CSV file, reads all rows, and returns a DataSet.
// It trims whitespace, ignores blank lines, and safely handles quoted fields.
func ReadCSV(path string) (*DataSet, error) {
	// --- Open file ---
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open CSV file: %w", err)
	}
	defer f.Close()

	// --- Initialize reader ---
	reader := csv.NewReader(f)
	reader.FieldsPerRecord = -1 // allow variable-length rows
	reader.TrimLeadingSpace = true
	reader.LazyQuotes = true

	// --- Read all rows ---
	rawRows, err := reader.ReadAll()
	if err != nil {
		// If it’s an unexpected EOF, try reading line by line instead
		if err == io.EOF {
			return nil, fmt.Errorf("unexpected end of file while reading CSV: %w", err)
		}
		return nil, fmt.Errorf("error reading CSV: %w", err)
	}

	if len(rawRows) == 0 {
		return nil, fmt.Errorf("CSV file is empty: %s", path)
	}

	// --- Extract headers ---
	headers := cleanRow(rawRows[0])

	// --- Process data rows ---
	var rows [][]string
	for _, r := range rawRows[1:] {
		// Skip completely empty lines
		if len(strings.TrimSpace(strings.Join(r, ""))) == 0 {
			continue
		}
		rows = append(rows, cleanRow(r))
	}

	data := &DataSet{
		Headers: headers,
		Rows:    rows,
		Source:  filepath.Base(path),
	}

	return data, nil
}

// cleanRow trims whitespace and normalizes cell contents in a row.
func cleanRow(row []string) []string {
	cleaned := make([]string, len(row))
	for i, val := range row {
		val = strings.TrimSpace(val)
		// Normalize line breaks and stray carriage returns
		val = strings.ReplaceAll(val, "\r", "")
		cleaned[i] = val
	}
	return cleaned
}
