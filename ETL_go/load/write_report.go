package load

import (
	"fmt"

	"etl_go/types"
)

// ReportSummary holds all stats for the ETL run summary.
type ReportSummary struct {
	TotalProcessed      int
	TotalRemoved        int
	RemovedNoPhone      int
	RemovedNoName       int
	RemovedInvalidState int
	RemovedDuplicates   int
	GeoStats			types.GeoStats
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
		fmt.Sprintf("    - %d removed for duplicate phone numbers", report.RemovedDuplicates),
		"",
		"  Geographic Data Cleaning:",
		fmt.Sprintf("    - %d ZIP codes cleaned (contained letters)", report.GeoStats.CleanedZipLetters),
		fmt.Sprintf("    - %d ZIP codes cleaned (too short)", report.GeoStats.CleanedZipTooShort),
		fmt.Sprintf("    - %d missing ZIP codes populated", report.GeoStats.PopulatedZip),
		fmt.Sprintf("    - %d missing states populated", report.GeoStats.PopulatedState),
		fmt.Sprintf("    - %d ZIP-State mismatches corrected", report.GeoStats.CorrectedMismatches),
		fmt.Sprintf("    - %d state/ZIP pairs fixed from area codes", report.GeoStats.FixedFromAreaCode),
		"",
		fmt.Sprintf("Total rows in final, ready-to-load file: %d", report.FinalRowCount),
		"",
		"====================================================",
	}

	return lines
}
