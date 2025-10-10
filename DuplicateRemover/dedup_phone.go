package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"regexp"
)

// Clean non-digit characters out of phone numbers
func normalizePhone(phone string) string {
	re := regexp.MustCompile(`\D`)
	num := re.ReplaceAllString(phone, "")
	if len(num) == 11 && num[0] == '1' {
		num = num[1:]
	}
	return num
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go run dedup_phone.go <input.csv>")
		os.Exit(1)
	}

	inFile := os.Args[1]
	outFile := "clean.csv"

	// --- Open input ---
	f, err := os.Open(inFile)
	if err != nil {
		log.Fatalf("Error opening input file: %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	reader.LazyQuotes = true
	rows, err := reader.ReadAll()
	if err != nil {
		log.Fatalf("Error reading CSV: %v", err)
	}

	if len(rows) < 2 {
		log.Fatalf("CSV appears empty or missing data rows")
	}

	header := rows[0]
	colIndex := -1
	for i, h := range header {
		if h == "phone1" {
			colIndex = i
			break
		}
	}
	if colIndex == -1 {
		log.Fatalf("No 'phone1' column found in CSV headers")
	}

	// --- Deduplicate ---
	seen := make(map[string]bool)
	cleaned := [][]string{header}
	duplicates := 0

	for _, row := range rows[1:] {
		if colIndex >= len(row) {
			continue
		}
		phone := normalizePhone(row[colIndex])
		if phone == "" {
			cleaned = append(cleaned, row)
			continue
		}
		if seen[phone] {
			duplicates++
			continue // skip duplicate
		}
		seen[phone] = true
		row[colIndex] = phone
		cleaned = append(cleaned, row)
	}

	// --- Write output ---
	out, err := os.Create(outFile)
	if err != nil {
		log.Fatalf("Error creating output file: %v", err)
	}
	defer out.Close()

	writer := csv.NewWriter(out)
	defer writer.Flush()

	err = writer.WriteAll(cleaned)
	if err != nil {
		log.Fatalf("Error writing CSV: %v", err)
	}

	fmt.Printf("✅ %d total rows processed\n", len(rows)-1)
	fmt.Printf("🗑️  %d duplicate rows removed\n", duplicates)
	fmt.Printf("📄 Clean CSV written to %s\n", outFile)
}
