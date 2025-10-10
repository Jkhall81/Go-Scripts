package transform

import (
	"fmt"
	"strconv"
	"strings"

	"etl_go/extract"

	"github.com/fatih/color"
)

// PreviewData displays headers and the first 'limit' rows from a DataSet.
// Each header is color-coded and indexed for reference in drop commands.
func PreviewData(ds *extract.DataSet, limit int) {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return
	}

	// --- Display headers with colored indexes ---
	fmt.Println("\nHeaders:")
	for i, h := range ds.Headers {
		label := color.New(color.FgCyan, color.Bold).SprintFunc()
		fmt.Printf("[%s] %s\n", label(strconv.Itoa(i)), h)
	}

	fmt.Println("\nPreview (first", limit, "rows):")

	// --- Show up to 'limit' rows ---
	for i := 0; i < limit && i < len(ds.Rows); i++ {
		row := ds.Rows[i]
		fmt.Printf("%d. %s\n", i+1, strings.Join(row, " | "))
	}

	fmt.Println()
}

// DropColumns removes the specified column indexes from the DataSet and returns a new copy.
func DropColumns(ds *extract.DataSet, indexes []int) *extract.DataSet {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return ds
	}

	// Make a map for fast lookup
	toDrop := make(map[int]bool)
	for _, idx := range indexes {
		toDrop[idx] = true
	}

	// Drop from headers
	var newHeaders []string
	for i, h := range ds.Headers {
		if !toDrop[i] {
			newHeaders = append(newHeaders, h)
		}
	}

	// Drop from rows
	var newRows [][]string
	for _, row := range ds.Rows {
		var newRow []string
		for i, val := range row {
			if !toDrop[i] {
				newRow = append(newRow, val)
			}
		}
		newRows = append(newRows, newRow)
	}

	return &extract.DataSet{
		Headers: newHeaders,
		Rows:    newRows,
		Source:  ds.Source,
	}
}

// ParseIndexes parses user input like "drop 0 3 6" into []int{0, 3, 6}.
func ParseIndexes(input string) ([]int, error) {
	parts := strings.Fields(input)
	if len(parts) < 2 {
		return nil, fmt.Errorf("no indexes provided")
	}

	var indexes []int
	for _, p := range parts[1:] {
		i, err := strconv.Atoi(p)
		if err != nil {
			return nil, fmt.Errorf("invalid index: %s", p)
		}
		indexes = append(indexes, i)
	}
	return indexes, nil
}
