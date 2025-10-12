package transform

import (
	"fmt"
	"regexp"
	"strings"

	"etl_go/extract"
)

// DedupResult holds the results of phone deduplication
type DedupResult struct {
	Cleaned    *extract.DataSet
	Duplicates int
}

// DedupPhones removes duplicate rows based on normalized phone numbers
func DedupPhones(ds *extract.DataSet) *DedupResult {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return &DedupResult{Cleaned: ds, Duplicates: 0}
	}

	// Find phone column index (adjust based on your column position)
	const phoneColIdx = 8 // Assuming phone is at index 8, adjust as needed

	// If we can't find the phone column in headers, try to locate it
	actualPhoneIdx := phoneColIdx
	if len(ds.Headers) > phoneColIdx && !strings.Contains(strings.ToLower(ds.Headers[phoneColIdx]), "phone") {
		// Try to find phone column by name
		for i, header := range ds.Headers {
			if strings.Contains(strings.ToLower(header), "phone") {
				actualPhoneIdx = i
				break
			}
		}
	}

	seen := make(map[string]bool)
	var uniqueRows [][]string
	duplicates := 0

	// Always keep the header
	uniqueRows = append(uniqueRows, ds.Headers)

	for _, row := range ds.Rows {
		if actualPhoneIdx >= len(row) {
			// If phone column doesn't exist in this row, keep it
			uniqueRows = append(uniqueRows, row)
			continue
		}

		phone := normalizePhone(row[actualPhoneIdx])
		if phone == "" {
			// If no phone number, keep the row
			uniqueRows = append(uniqueRows, row)
			continue
		}

		if seen[phone] {
			duplicates++
			continue // skip duplicate
		}

		seen[phone] = true
		// Update the row with normalized phone number
		newRow := make([]string, len(row))
		copy(newRow, row)
		newRow[actualPhoneIdx] = phone
		uniqueRows = append(uniqueRows, newRow)
	}

	return &DedupResult{
		Cleaned: &extract.DataSet{
			Headers: ds.Headers,
			Rows:    uniqueRows[1:], // Skip header row
			Source:  ds.Source,
		},
		Duplicates: duplicates,
	}
}

// normalizePhone cleans and normalizes phone numbers
func normalizePhone(phone string) string {
	re := regexp.MustCompile(`\D`)
	num := re.ReplaceAllString(phone, "")
	if len(num) == 11 && num[0] == '1' {
		num = num[1:]
	}
	return num
}