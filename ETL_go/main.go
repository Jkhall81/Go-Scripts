package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"etl_go/extract"
	"etl_go/load"
	"etl_go/transform"
	"etl_go/ui"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: etl_go <inputfile.csv | inputfile.xlsx>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	ext := filepath.Ext(inputFile)

	var ds *extract.DataSet
	var err error

	// Track checklist progress
	status := make(map[string]bool)
	for _, step := range ui.Steps {
		status[step] = false
	}

	// --- Extract Step ---
	switch ext {
	case ".csv":
		ds, err = extract.ReadCSV(inputFile)
	case ".xlsx":
		ds, err = extract.ReadXLSX(inputFile)
	default:
		log.Fatalf("Unsupported file type: %s", ext)
	}
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}
	status["Extract CSV/XLSX"] = true
	ui.AddToOutput(fmt.Sprintf("Loaded %d rows from %s", len(ds.Rows), ds.Source))

	// --- Launch the UI ---
	ui.DrawUI(status, func(cmd string) {
		args := strings.Fields(cmd)
		if len(args) == 0 {
			return
		}

		switch strings.ToLower(args[0]) {

		// 🟡 Preview data
		case "show":
			for _, line := range ds.FirstNLines(5) {
				ui.AddToOutput(line)
			}

		// 🟡 Drop columns by index
		case "drop":
			if len(args) < 2 {
				ui.AddToOutput("Usage: drop <colIndex1> <colIndex2> ...")
				return
			}
			indexes := []int{}
			for _, s := range args[1:] {
				i, err := strconv.Atoi(s)
				if err == nil {
					indexes = append(indexes, i)
				}
			}
			ds = transform.DropColumns(ds, indexes)
			status["Drop Columns"] = true
			ui.AddToOutput(fmt.Sprintf("Dropped columns: %v", indexes))

		// 🟢 Clean addresses
		case "clean-address":
			ds = transform.CleanAddresses(ds)
			status["Clean Addresses"] = true
			ui.AddToOutput("Cleaned address fields.")

		// 🟢 Normalize phone numbers
		case "normalize-phones":
			ds = transform.NormalizePhones(ds)
			status["Normalize Phones"] = true
			ui.AddToOutput("Normalized phone numbers.")

		// 🟢 Populate geo
		case "populate-geo":
			ds = transform.PopulateGeo(ds)
			status["Populate Geo"] = true
			ui.AddToOutput("Populated missing state/ZIP data.")

		// 🟢 Validate states
		case "validate-states":
			result := transform.ValidateStates(ds)
			ds = result.Cleaned
			status["Validate States"] = true
			ui.AddToOutput(fmt.Sprintf("Removed %d invalid-state rows.", result.DropCount))

		// 🟢 Final validation
		case "final-validate":
			result := load.FinalValidate(ds)
			ds = result.Cleaned
			status["Final Validation"] = true
			ui.AddToOutput(fmt.Sprintf("Removed %d invalid rows.", result.DropCount))

		// 🟢 Write CSV
		case "write-csv":
			if err := load.WriteCSV(ds, ""); err != nil {
				ui.AddToOutput(fmt.Sprintf("Error writing CSV: %v", err))
				return
			}
			status["Write CSV"] = true
			ui.AddToOutput("Output CSV written successfully.")

		// 🟢 Write report
		case "write-report":
			report := load.ReportSummary{
				TotalProcessed: len(ds.Rows),
				FinalRowCount:  len(ds.Rows),
			}
			load.WriteReport(report)
			status["Write Report"] = true
			ui.AddToOutput("Report generated successfully.")

		// 🧭 Help
		case "help":
			ui.AddToOutput("Available commands:")
			ui.AddToOutput("show, drop, clean-address, normalize-phones, populate-geo,")
			ui.AddToOutput("validate-states, final-validate, write-csv, write-report, exit")

		// 🚪 Exit
		case "exit", "quit":
			ui.AddToOutput("Exiting ETL tool...")
			os.Exit(0)

		default:
			ui.AddToOutput(fmt.Sprintf("Unknown command: %s", cmd))
		}
	})
}
