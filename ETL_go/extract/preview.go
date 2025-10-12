package extract

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func (ds *DataSet) FirstNLines(n int, availableWidth int) []string {
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

	// --- Calculate column widths ---
	colWidths := calculateAdaptiveColumnWidths(ds.Headers, ds.Rows[:limit], availableWidth)
	lines := []string{}

	// --- 🆕 Zero-indexed column numbers (above header, blue) ---
	indexLine := "|"
	for j := range ds.Headers {
		width := colWidths[j]
		indexLine += fmt.Sprintf(" %-*d |", width, j)
	}
	blue := lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
	lines = append(lines, blue.Render(indexLine))

	// --- Header row ---
	header := "|"
	for j, h := range ds.Headers {
		width := colWidths[j]
		header += fmt.Sprintf(" %-*s |", width, h)
	}
	lines = append(lines, header)

	// --- Separator line ---
	sep := "+"
	for _, w := range colWidths {
		sep += strings.Repeat("-", w+2) + "+"
	}
	lines = append(lines, sep)

	// --- Data rows ---
	for i := 0; i < limit; i++ {
		row := ds.Rows[i]
		line := "|"
		for j := 0; j < len(ds.Headers); j++ {
			val := ""
			if j < len(row) {
				val = strings.ReplaceAll(row[j], "\n", " ")
			}
			width := colWidths[j]
			if len(val) > width {
				if width > 3 {
					val = val[:width-3] + "..."
				} else if width > 0 {
					val = val[:width]
				} else {
					val = ""
				}
			}
			line += fmt.Sprintf(" %-*s |", width, val)
		}
		lines = append(lines, line)
	}

	lines = append(lines,
		fmt.Sprintf("Showing %d of %d rows from %s", limit, len(ds.Rows), ds.Source),
		fmt.Sprintf("Table width: %d characters (available: %d)",
			calculateTotalTableWidth(colWidths), availableWidth),
	)

	return lines
}

// calculateAdaptiveColumnWidths optimizes column widths for the available space
func calculateAdaptiveColumnWidths(headers []string, rows [][]string, availableWidth int) []int {
	colWidths := make([]int, len(headers))

	// First pass: calculate ideal widths from data
	for j, h := range headers {
		colWidths[j] = len(h)
	}

	for _, row := range rows {
		for j, val := range row {
			if j < len(colWidths) && len(val) > colWidths[j] {
				colWidths[j] = len(val)
			}
		}
	}

	// Calculate total ideal width
	totalIdealWidth := calculateTotalTableWidth(colWidths)

	// If table fits in available width, use ideal widths
	if totalIdealWidth <= availableWidth {
		return colWidths
	}

	// If too wide, scale down proportionally
	scaleFactor := float64(availableWidth) / float64(totalIdealWidth)

	// Apply scaling with minimum widths
	for j := range colWidths {
		scaledWidth := int(float64(colWidths[j]) * scaleFactor)
		// Ensure minimum readable width
		if scaledWidth < 8 {
			scaledWidth = 8
		}
		// Don't scale below header width
		if scaledWidth < len(headers[j]) {
			scaledWidth = len(headers[j])
		}
		colWidths[j] = scaledWidth
	}

	// Final adjustment to ensure we don't exceed max width
	finalWidth := calculateTotalTableWidth(colWidths)
	if finalWidth > availableWidth {
		// Trim the widest column until it fits
		for finalWidth > availableWidth {
			maxWidthIndex := -1
			maxWidth := 0
			for j, width := range colWidths {
				if width > maxWidth && width > len(headers[j]) {
					maxWidth = width
					maxWidthIndex = j
				}
			}
			if maxWidthIndex == -1 {
				break // Can't reduce further without cutting headers
			}
			colWidths[maxWidthIndex]--
			finalWidth = calculateTotalTableWidth(colWidths)
		}
	}

	for j := range colWidths {
		if colWidths[j] < 1 {
			colWidths[j] = 1
		}
	}

	return colWidths
}

// calculateTotalTableWidth calculates the total character width of the table
func calculateTotalTableWidth(colWidths []int) int {
	if len(colWidths) == 0 {
		return 0
	}

	total := 0
	for _, width := range colWidths {
		total += width + 3 // +3 for " | " between columns
	}
	total += 1 // for the starting "|"

	return total
}
