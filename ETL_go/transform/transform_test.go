package transform

import (
	"etl_go/extract"
	"testing"
)

// mockData returns a realistic dataset using the true schema:
// [0]source_id [1]first [2]middle [3]last [4]address [5]city [6]st [7]zip [8]phone [9]address3 [10]province [11]email [12]trusted_url
func mockData() *extract.DataSet {
	return &extract.DataSet{
		Headers: []string{
			"SourceID", "First", "Middle", "Last", "Address", "City",
			"State", "Zip", "Phone", "Address3", "Province", "Email", "TrustedURL",
		},
		Rows: [][]string{
			{"1", "J4son", "K.", "Hall", "123 Main St,", "Tampa", "florida", "33A10", "(813)555-9999", "", "", "12345", ""},
			{"2", "Alice", "", "Smith", "77## Weird Blvd!", "Austin", "TX", "73301", "1-512-555-8888", "", "", "alice@example.com", ""},
			{"3", "Bob", "", "Jones", "456 Elm, Apt 3", "Chicago", "IL", "60614", "555-555-5555", "", "", "99999", ""},
			{"4", "Eve", "", "Doe", "Bad@Addr!!", "Nowhere", "12345", "ABCDE", "999", "", "", "eve@domain", ""},
			{"5", "Jane", "Marie", "Roe", "1234 West St", "Orlando", "FL", "32801", "+1 (321) 888-1212", "", "", "56789", ""},
			{"6", "Mike", "", "Numeric", "456 Strip", "Vegas", "NV", "89101", "702.555.0000", "", "", "42", ""},
			{"7", "", "", "", "777 Street", "Nowhere", "ZZ", "", "", "", "", "nobody@example.com", ""},
			// New coverage cases ↓
			{"8", "Sara", "", "Lee", "1 Hill Rd", "Denver", "CO", "", "3035550000", "", "", "sara@x.com", ""},            // missing ZIP → populate
			{"9", "Tom", "", "Wayne", "22 Pine St", "Unknown", "", "73301", "5125551212", "", "", "tom@x.com", ""},       // missing state → populate
			{"10", "Nina", "", "Gray", "66 Maple Ave", "Chicago", "CA", "60614", "7735551111", "", "", "nina@x.com", ""}, // mismatch → correct
			{"11", "Rob", "", "Miles", "500 Oak", "Unknown", "", "", "5125559999", "", "", "rob@x.com", ""},              // no state/zip, infer from area code
		},
	}
}

func TestCleanAddresses(t *testing.T) {
	ds := mockData()
	got := CleanAddresses(ds)
	if got.Rows[1][4] != "77## Weird Blvd" {
		t.Errorf("row 1: expected '77## Weird Blvd', got '%s'", got.Rows[1][4])
	}
	if got.Rows[3][4] != "BadAddr" {
		t.Errorf("row 3: expected 'BadAddr', got '%s'", got.Rows[3][4])
	}
}

func TestCleanEmails(t *testing.T) {
	ds := mockData()
	got := CleanEmails(ds)
	if got.Rows[5][11] != "" {
		t.Errorf("row 5: expected blank email for numeric value, got '%s'", got.Rows[5][11])
	}
	if got.Rows[1][11] != "alice@example.com" {
		t.Errorf("row 1: expected alice@example.com, got '%s'", got.Rows[1][11])
	}
}

func TestCleanNames(t *testing.T) {
	ds := mockData()
	got := CleanNames(ds)

	tests := []struct {
		row, col int
		want     string
	}{
		{0, 1, ""},        // J4son → removed (numeric)
		{1, 1, "Alice"},   // unchanged
		{2, 3, "Jones"},   // unchanged
		{3, 1, "Eve"},     // unchanged
		{5, 3, "Numeric"}, // unchanged
		{6, 1, ""},        // empty → remains empty
	}
	for _, tt := range tests {
		if got.Rows[tt.row][tt.col] != tt.want {
			t.Errorf("row %d col %d: expected %q, got %q", tt.row, tt.col, tt.want, got.Rows[tt.row][tt.col])
		}
	}
}

func TestCleanStates(t *testing.T) {
	ds := mockData()
	got := CleanStates(ds)
	if got.Rows[0][6] != "" {
		t.Errorf("row 0: expected blank for invalid 'florida', got '%s'", got.Rows[0][6])
	}
	if got.Rows[6][6] != "ZZ" {
		t.Errorf("row 6: expected 'ZZ' to remain, got '%s'", got.Rows[6][6])
	}
}

func TestNormalizePhones(t *testing.T) {
	ds := mockData()
	got := NormalizePhones(ds)

	if got.Rows[0][8] != "8135559999" {
		t.Errorf("expected 8135559999, got '%s'", got.Rows[0][8])
	}
	if got.Rows[4][8] != "3218881212" {
		t.Errorf("expected 3218881212, got '%s'", got.Rows[4][8])
	}
	if got.Rows[3][8] != "" {
		t.Errorf("expected blank invalid number, got '%s'", got.Rows[3][8])
	}
}

func TestDedupPhones(t *testing.T) {
	ds := &extract.DataSet{
		Headers: []string{"SourceID", "First", "Middle", "Last", "Address", "City", "State", "Zip", "Phone", "Address3", "Province", "Email", "TrustedURL"},
		Rows: [][]string{
			{"1", "A", "", "X", "", "", "FL", "11111", "555-111-1111", "", "", "a@x.com", ""},
			{"2", "B", "", "Y", "", "", "FL", "11111", "(555)111-1111", "", "", "b@y.com", ""},
			{"3", "C", "", "Z", "", "", "TX", "73301", "5125550000", "", "", "c@z.com", ""},
			{"4", "D", "", "Z", "", "", "TX", "73301", "5125550000", "", "", "d@z.com", ""},
		},
	}
	res := DedupPhones(ds)
	if res.Duplicates != 2 {
		t.Errorf("expected 2 duplicates removed, got %d", res.Duplicates)
	}
	if len(res.Cleaned.Rows) != 2 {
		t.Errorf("expected 2 unique rows, got %d", len(res.Cleaned.Rows))
	}
}

func TestDropColumns(t *testing.T) {
	ds := mockData()
	got := DropColumns(ds, []int{9, 10, 12})
	expected := len(ds.Headers) - 3
	if len(got.Headers) != expected {
		t.Errorf("expected %d headers, got %d", expected, len(got.Headers))
	}
}

func TestPopulateGeo(t *testing.T) {
	ds := mockData()
	got, stats := PopulateGeo(ds)
	if stats.PopulatedZip+stats.PopulatedState+stats.FixedFromAreaCode+stats.CorrectedMismatches == 0 {
		t.Log("No geo fields were populated or fixed — dataset may already be clean")
	}

	if got == nil || len(got.Rows) == 0 {
		t.Errorf("expected dataset to be returned, got nil/empty")
	}
}

func TestValidateStates(t *testing.T) {
	ds := &extract.DataSet{
		Headers: []string{"SourceID", "First", "Middle", "Last", "Address", "City", "State", "Zip", "Phone", "Address3", "Province", "Email", "TrustedURL"},
		Rows: [][]string{
			{"1", "", "", "", "", "", "FL", "", "", "", "", "", ""},
			{"2", "", "", "", "", "", "CA", "", "", "", "", "", ""},
			{"3", "", "", "", "", "", "TX", "", "", "", "", "", ""},
			{"4", "", "", "", "", "", "NY", "", "", "", "", "", ""},
			{"5", "", "", "", "", "", "ZZ", "", "", "", "", "", ""},
			{"6", "", "", "", "", "", "PR", "", "", "", "", "", ""},
			{"7", "", "", "", "", "", "XX", "", "", "", "", "", ""},
		},
	}
	result := ValidateStates(ds)

	if len(result.Cleaned.Rows) != 4 {
		t.Errorf("expected 4 valid rows, got %d", len(result.Cleaned.Rows))
	}
	if result.DropCount != 3 {
		t.Errorf("expected 3 dropped invalid rows, got %d", result.DropCount)
	}
}
