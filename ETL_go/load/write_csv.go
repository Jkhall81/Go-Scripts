package load

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"etl_go/extract"
)

// WriteCSV writes the cleaned dataset to a .csv file.
// If outFile is blank, it auto-generates a name based on the source file.
func WriteCSV(ds *extract.DataSet, outFile string) error {
	if ds == nil || len(ds.Rows) == 0 {
		return fmt.Errorf("no data to write")
	}

	// If no file name provided, build one based on source
	if outFile == "" {
		base := "output"
		if ds.Source != "" {
			base = strings.TrimSuffix(filepath.Base(ds.Source), filepath.Ext(ds.Source))
		}
		outFile = fmt.Sprintf("%s_cleaned.csv", base)
	}

	// Create output file
	f, err := os.Create(outFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer f.Close()

	w := csv.NewWriter(f)
	defer w.Flush()

	// Write header row first
	if err := w.Write(ds.Headers); err != nil {
		return fmt.Errorf("failed to write headers: %v", err)
	}

	// Write all rows
	if err := w.WriteAll(ds.Rows); err != nil {
		return fmt.Errorf("failed to write rows: %v", err)
	}

	fmt.Printf("✅ %d rows written to %s\n", len(ds.Rows), outFile)
	return nil
}
