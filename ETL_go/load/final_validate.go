package load

import (
	"fmt"
	"strings"

	"etl_go/extract"
)

// FinalValidationResult holds cleaned dataset and dropped rows.
type FinalValidationResult struct {
	Cleaned   *extract.DataSet
	Dropped   [][]string
	DropCount int
}

// FinalValidate removes rows missing a phone number or both first and last names.
// Returns a FinalValidationResult so dropped rows can be logged later.
func FinalValidate(ds *extract.DataSet) *FinalValidationResult {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return &FinalValidationResult{Cleaned: ds}
	}

	var (
		validRows [][]string
		dropped   [][]string
	)

	for _, row := range ds.Rows {
		// Defensive check for malformed rows
		if len(row) < 9 {
			dropped = append(dropped, row)
			continue
		}

		first := strings.TrimSpace(row[1])
		last := strings.TrimSpace(row[3])
		phone := strings.TrimSpace(row[8])

		// If missing phone OR both names missing → drop
		if phone == "" || (first == "" && last == "") {
			dropped = append(dropped, row)
		} else {
			validRows = append(validRows, row)
		}
	}

	dropCount := len(dropped)

	if dropCount > 0 {
		fmt.Printf("%d rows removed due to missing required fields (phone/name).\n", dropCount)
	} else {
		fmt.Println("All rows passed final validation.")
	}

	cleaned := &extract.DataSet{
		Headers: ds.Headers,
		Rows:    validRows,
		Source:  ds.Source,
	}

	return &FinalValidationResult{
		Cleaned:   cleaned,
		Dropped:   dropped,
		DropCount: dropCount,
	}
}
