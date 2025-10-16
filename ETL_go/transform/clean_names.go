package transform

import (
	"fmt"
	"regexp"
	"strings"

	"etl_go/extract"
)

// CleanNames removes numeric values and special characters from first, middle, and last name fields.
func CleanNames(ds *extract.DataSet) *extract.DataSet {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return ds
	}

	// Name column indices (0-based) - adjust these based on your actual column positions
	const firstNameIdx = 1
	const middleNameIdx = 2
	const lastNameIdx = 3

	// Regex: keep only letters, spaces, hyphens, and apostrophes for names
	re := regexp.MustCompile(`[^A-Za-z\s\-']`)

	stats := struct {
		cleanedNumeric int
		cleanedSpecial int
	}{}

	newRows := make([][]string, len(ds.Rows))
	for i, row := range ds.Rows {
		newRow := make([]string, len(row))
		copy(newRow, row)

		// Process first name
		if firstNameIdx < len(row) && row[firstNameIdx] != "" {
			cleaned := cleanNameField(row[firstNameIdx], re, &stats)
			newRow[firstNameIdx] = cleaned
		}

		// Process middle name
		if middleNameIdx < len(row) && row[middleNameIdx] != "" {
			cleaned := cleanNameField(row[middleNameIdx], re, &stats)
			newRow[middleNameIdx] = cleaned
		}

		// Process last name
		if lastNameIdx < len(row) && row[lastNameIdx] != "" {
			cleaned := cleanNameField(row[lastNameIdx], re, &stats)
			newRow[lastNameIdx] = cleaned
		}

		newRows[i] = newRow
	}

	return &extract.DataSet{
		Headers: ds.Headers,
		Rows:    newRows,
		Source:  ds.Source,
	}
}

// cleanNameField processes a single name field
func cleanNameField(value string, re *regexp.Regexp, stats *struct {
	cleanedNumeric int
	cleanedSpecial int
}) string {
	original := strings.TrimSpace(value)
	if original == "" {
		return ""
	}

	// Check if the value contains ANY numeric characters
	if containsAnyNumbers(original) {
		stats.cleanedNumeric++
		return "" // Remove the entire value if it contains any numbers
	}

	// Remove special characters (keep only letters, spaces, hyphens, apostrophes)
	cleaned := re.ReplaceAllString(original, "")
	cleaned = strings.TrimSpace(cleaned)

	// Track if we removed special characters
	if cleaned != original {
		stats.cleanedSpecial++
	}

	return cleaned
}

// containsAnyNumbers checks if a string contains any numeric characters
func containsAnyNumbers(s string) bool {
	for _, char := range s {
		if char >= '0' && char <= '9' {
			return true
		}
	}
	return false
}
