package load

import (
	"fmt"
)

// ReportSummary holds all stats for the ETL run summary.
type ReportSummary struct {
	TotalProcessed         int
	TotalRemoved           int
	RemovedNoPhone         int
	RemovedNoName          int
	RemovedInvalidState    int
	FinalRowCount          int
}

// WriteReport prints the summary of ETL operations to the console.
func WriteReport(report ReportSummary) {
	fmt.Println("\n================ ETL SUMMARY REPORT ================")

	fmt.Printf("Total rows processed: %d\n", report.TotalProcessed)
	fmt.Printf("Total rows removed:   %d\n", report.TotalRemoved)

	fmt.Println("  Breakdown:")
	fmt.Printf("    - %d removed for missing phone number\n", report.RemovedNoPhone)
	fmt.Printf("    - %d removed for missing first AND last name\n", report.RemovedNoName)
	fmt.Printf("    - %d removed for invalid state\n", report.RemovedInvalidState)

	fmt.Printf("\nTotal rows in final, ready-to-load file: %d\n", report.FinalRowCount)

	fmt.Println("\n====================================================")
}
