package extract

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/xuri/excelize/v2"
)

// ReadXLSX opens a .xlsx file, validates worksheet count,
// and converts the data to a DataSet compatible with CSV reads.
func ReadXLSX(path string) (*DataSet, error) {
	// --- Open the workbook ---
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open xlsx file: %w", err)
	}
	defer func() { _ = f.Close() }()

	// --- Validate sheet count ---
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no worksheets found in file: %s", path)
	}
	if len(sheets) > 1 {
		return nil, fmt.Errorf("multiple worksheets detected (%d). Please ensure only one worksheet", len(sheets))
	}

	sheet := sheets[0]

	// --- Read all rows from the single sheet ---
	rows, err := f.GetRows(sheet)
	if err != nil {
		return nil, fmt.Errorf("failed to read rows from %s: %w", sheet, err)
	}
	if len(rows) == 0 {
		return nil, fmt.Errorf("worksheet '%s' is empty", sheet)
	}

	// --- Clean and normalize rows ---
	headers := cleanRow(rows[0])
	var dataRows [][]string

	for _, row := range rows[1:] {
		// Skip empty rows (Excel sometimes has trailing blanks)
		if isRowEmpty(row) {
			continue
		}
		dataRows = append(dataRows, cleanRow(row))
	}

	data := &DataSet{
		Headers: headers,
		Rows:    dataRows,
		Source:  fmt.Sprintf("%s (sheet: %s)", filepath.Base(path), sheet),
	}

	return data, nil
}

// isRowEmpty checks if a row contains only empty cells.
func isRowEmpty(row []string) bool {
	for _, cell := range row {
		if strings.TrimSpace(cell) != "" {
			return false
		}
	}
	return true
}
