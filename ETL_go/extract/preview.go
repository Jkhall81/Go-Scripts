package extract

import (
	"fmt"
	"strings"
)

func (ds *DataSet) FirstNLines(n int) []string {
	if ds == nil {
		return []string{"Dataset is nil."}
	}
	if len(ds.Rows) == 0 {
		return []string{"No data loaded."}
	}

	limit := n
	if len(ds.Rows) < n {
		limit = len(ds.Rows)
	}

	const maxColWidth = 25
	colWidths := make([]int, len(ds.Headers))
	for j, h := range ds.Headers {
		colWidths[j] = len(h)
		if colWidths[j] > maxColWidth {
			colWidths[j] = maxColWidth
		}
	}
	for _, row := range ds.Rows[:limit] {
		for j, val := range row {
			if j < len(colWidths) {
				if l := len(val); l > colWidths[j] {
					colWidths[j] = min(l, maxColWidth)
				}
			}
		}
	}

	lines := []string{}
	lines = append(lines, fmt.Sprintf("Loaded %d rows from %s", len(ds.Rows), ds.Source))

	// header
	header := "|"
	for j, h := range ds.Headers {
		width := min(colWidths[j], maxColWidth)
		header += fmt.Sprintf(" %-*.*s |", width, width, h)
	}
	lines = append(lines, header)

	// separator
	sep := "+"
	for _, w := range colWidths {
		w = min(w, maxColWidth)
		sep += strings.Repeat("-", w+2) + "+"
	}
	lines = append(lines, sep)

	// data rows
	for i := 0; i < limit; i++ {
		row := ds.Rows[i]
		line := "|"
		for j := 0; j < len(ds.Headers); j++ {
			val := ""
			if j < len(row) {
				val = strings.ReplaceAll(row[j], "\n", " ")
			}
			width := min(colWidths[j], maxColWidth)
			line += fmt.Sprintf(" %-*.*s |", width, width, val)
		}
		lines = append(lines, line)
	}

	lines = append(lines, fmt.Sprintf("Showing %d of %d rows from %s", limit, len(ds.Rows), ds.Source))
	return lines
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
