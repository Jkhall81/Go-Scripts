package transform

import (
	"strings"
	"unicode"

	"etl_go/extract"
)

// CleanStates ensures that the state field (row[6]) contains only valid 2-letter alphabetic codes.
// If not, it blanks out the state value but does NOT drop the row.
func CleanStates(ds *extract.DataSet) *extract.DataSet {
	if ds == nil {
		return ds
	}

	newRows := make([][]string, len(ds.Rows))

	for i, row := range ds.Rows {
		newRow := make([]string, len(row))
		copy(newRow, row)

		if len(newRow) > 6 {
			state := strings.ToUpper(strings.TrimSpace(newRow[6]))

			if !isTwoLetterAlpha(state) {
				newRow[6] = ""
			} else {
				newRow[6] = state
			}
		}

		newRows[i] = newRow
	}

	return &extract.DataSet{
		Headers: ds.Headers,
		Rows:    newRows,
		Source:  ds.Source,
	}
}

// Helper: returns true only if the string is exactly two letters A–Z
func isTwoLetterAlpha(s string) bool {
	if len(s) != 2 {
		return false
	}
	for _, ch := range s {
		if !unicode.IsLetter(ch) {
			return false
		}
	}
	return true
}
