package transform

import (
	"fmt"
	"strings"

	"etl_go/extract"
)

// AllowedStates is the list of valid US state codes + DC.
var AllowedStates = map[string]bool{
	"AL": true, "AK": true, "AZ": true, "AR": true, "CA": true,
	"CO": true, "CT": true, "DE": true, "DC": true, "FL": true,
	"GA": true, "HI": true, "ID": true, "IL": true, "IN": true,
	"IA": true, "KS": true, "KY": true, "LA": true, "ME": true,
	"MD": true, "MA": true, "MI": true, "MN": true, "MS": true,
	"MO": true, "MT": true, "NE": true, "NV": true, "NH": true,
	"NJ": true, "NM": true, "NY": true, "NC": true, "ND": true,
	"OH": true, "OK": true, "OR": true, "PA": true, "RI": true,
	"SC": true, "SD": true, "TN": true, "TX": true, "UT": true,
	"VT": true, "VA": true, "WA": true, "WV": true, "WI": true, "WY": true,
}

// ValidationResult holds the cleaned dataset and any rows that were dropped.
type ValidationResult struct {
	Cleaned   *extract.DataSet
	Dropped   [][]string
	DropCount int
}

// ValidateStates removes rows with invalid state codes (non-50 states + DC).
// It returns a ValidationResult for later reporting.
func ValidateStates(ds *extract.DataSet) *ValidationResult {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return &ValidationResult{Cleaned: ds}
	}

	var (
		validRows [][]string
		dropped   [][]string
	)

	for _, row := range ds.Rows {
		if len(row) <= 6 {
			// Malformed row — drop it
			dropped = append(dropped, row)
			continue
		}

		state := strings.ToUpper(strings.TrimSpace(row[6]))
		if AllowedStates[state] {
			validRows = append(validRows, row)
		} else {
			dropped = append(dropped, row)
		}
	}

	dropCount := len(dropped)

	cleaned := &extract.DataSet{
		Headers: ds.Headers,
		Rows:    validRows,
		Source:  ds.Source,
	}

	return &ValidationResult{
		Cleaned:   cleaned,
		Dropped:   dropped,
		DropCount: dropCount,
	}
}
