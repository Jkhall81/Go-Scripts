package transform

import (
	"fmt"
	"regexp"

	"etl_go/extract"
)

// CleanEmails replaces numeric-only email values with an empty string.
// If the value in the email column is a number (e.g. "12345"), it will be cleared.
func CleanEmails(ds *extract.DataSet) *extract.DataSet {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return ds
	}

	// Email column index (0-based)
	const emailIdx = 11

	// Regex: matches strings that consist only of digits
	re := regexp.MustCompile(`^[0-9]+$`)

	newRows := make([][]string, len(ds.Rows))
	for i, row := range ds.Rows {
		newRow := make([]string, len(row))
		copy(newRow, row)

		if emailIdx < len(row) && row[emailIdx] != "" {
			if re.MatchString(row[emailIdx]) {
				newRow[emailIdx] = ""
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
