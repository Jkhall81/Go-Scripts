package transform

import (
	"fmt"
	"regexp"
	"strings"

	"etl_go/extract"
)

// NormalizePhones cleans and normalizes phone numbers to a 10-digit numeric format.
// It removes all non-digits and trims a leading '1' if the number has 11 digits.
// Invalid or empty numbers are left as blank strings.
func NormalizePhones(ds *extract.DataSet) *extract.DataSet {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return ds
	}

	// Phone number column index (based on 13-column schema)
	const phoneIdx = 8

	reDigits := regexp.MustCompile(`\D`) // matches all non-digit characters

	newRows := make([][]string, len(ds.Rows))
	for i, row := range ds.Rows {
		newRow := make([]string, len(row))
		copy(newRow, row)

		if phoneIdx < len(row) && row[phoneIdx] != "" {
			num := reDigits.ReplaceAllString(row[phoneIdx], "") // keep only digits

			// Remove leading "1" if 11 digits long (e.g. +1 country code)
			if len(num) == 11 && strings.HasPrefix(num, "1") {
				num = num[1:]
			}

			// Keep only valid 10-digit numbers
			if len(num) == 10 {
				newRow[phoneIdx] = num
			} else {
				newRow[phoneIdx] = ""
			}
		}

		newRows[i] = newRow
	}

	fmt.Println("Phone numbers normalized successfully.")
	return &extract.DataSet{
		Headers: ds.Headers,
		Rows:    newRows,
		Source:  ds.Source,
	}
}
