package transform

import (
	"fmt"
	"regexp"

	"etl_go/extract"
)

// CleanAddresses removes unwanted special characters from the address1 column.
// It keeps alphanumeric characters, spaces, and these symbols: # / - .
// Commas are removed to prevent CSV misalignment.
func CleanAddresses(ds *extract.DataSet) *extract.DataSet {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return ds
	}

	// Address1 column index (0-based)
	const address1Idx = 4

	// Regex: remove everything except A–Z, 0–9, spaces, # / - .
	re := regexp.MustCompile(`[^A-Za-z0-9\s#\/\-\.]`)

	newRows := make([][]string, len(ds.Rows))
	for i, row := range ds.Rows {
		newRow := make([]string, len(row))
		copy(newRow, row)

		if address1Idx < len(row) && row[address1Idx] != "" {
			cleaned := re.ReplaceAllString(row[address1Idx], "")
			newRow[address1Idx] = cleaned
		}

		newRows[i] = newRow
	}

	return &extract.DataSet{
		Headers: ds.Headers,
		Rows:    newRows,
		Source:  ds.Source,
	}
}
