package load

import (
	"fmt"
)

// ReportSummary holds all stats for the ETL run summary.
type ReportSummary struct {
	TotalProcessed      int
	TotalRemoved        int
	RemovedNoPhone      int
	RemovedNoName       int
	RemovedInvalidState int
	FinalRowCount       int
}

// WriteReport returns the summary of ETL operations as formatted strings for the output window.
func WriteReport(report ReportSummary) []string {
	lines := []string{
		"================ ETL SUMMARY REPORT ================",
		"",
		fmt.Sprintf("Total rows processed: %d", report.TotalProcessed),
		fmt.Sprintf("Total rows removed:   %d", report.TotalRemoved),
		"",
		"  Breakdown:",
		fmt.Sprintf("    - %d removed for missing phone number", report.RemovedNoPhone),
		fmt.Sprintf("    - %d removed for missing first AND last name", report.RemovedNoName),
		fmt.Sprintf("    - %d removed for invalid state", report.RemovedInvalidState),
		"",
		fmt.Sprintf("Total rows in final, ready-to-load file: %d", report.FinalRowCount),
		"",
		"====================================================",
	}

	return lines
}
