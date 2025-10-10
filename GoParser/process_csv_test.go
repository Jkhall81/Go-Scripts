package main

import (
	"testing"
)

// --- 1. Phone number cleanup ---

func TestCleanPhone_RemovesSymbols(t *testing.T) {
	got := cleanPhone("(555) 123-4567")
	want := "5551234567"
	if got != want {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestCleanPhone_BlankPhoneStaysBlank(t *testing.T) {
	got := cleanPhone("")
	if got != "" {
		t.Errorf("Expected empty string, got %s", got)
	}
}

// --- 2. State normalization ---

func TestNormalizeState_LowercaseToUppercase(t *testing.T) {
	got := normalizeState("ca")
	want := "CA"
	if got != want {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestNormalizeState_TrimsWhitespace(t *testing.T) {
	got := normalizeState("  ny ")
	want := "NY"
	if got != want {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

// --- 3. ZIP code truncation ---

func TestTruncateZip_LongZipTruncated(t *testing.T) {
	got := truncateZip("123456789")
	want := "12345"
	if got != want {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

func TestTruncateZip_ShortZipUnchanged(t *testing.T) {
	got := truncateZip("90210")
	want := "90210"
	if got != want {
		t.Errorf("Expected %s, got %s", want, got)
	}
}

// --- 4. Missing ZIP → populate from state ---

func TestPopulateZip_UsesStateMapping(t *testing.T) {
	row := make([]string, 13)
	row[6] = "CA"
	row[7] = ""
	populateZip(row)
	if row[7] != "90005" {
		t.Errorf("Expected CA ZIP 90005, got %s", row[7])
	}
}

func TestPopulateZip_DoesNothingIfZipAlreadyPresent(t *testing.T) {
	row := make([]string, 13)
	row[6] = "CA"
	row[7] = "94105"
	populateZip(row)
	if row[7] != "94105" {
		t.Errorf("Expected unchanged ZIP 94105, got %s", row[7])
	}
}

// --- 5. ZIP → state inference ---

func TestPopulateStateFromZip_FindsMatchingState(t *testing.T) {
	row := make([]string, 13)
	row[7] = "90210"
	populateStateFromZip(row)
	if row[6] != "CA" {
		t.Errorf("Expected CA for ZIP 90210, got %s", row[6])
	}
}

func TestPopulateStateFromZip_UnknownZipDoesNothing(t *testing.T) {
	row := make([]string, 13)
	row[7] = "999999"
	populateStateFromZip(row)
	if row[6] != "" {
		t.Errorf("Expected blank state for unknown ZIP, got %s", row[6])
	}
}

// --- 6. Area code → state/zip inference ---

func TestPopulateStateZipFromAreaCode_FillsCorrectly(t *testing.T) {
	row := make([]string, 13)
	row[8] = "5125559999" // Texas
	populateStateZipFromAreaCode(row)
	if row[6] != "TX" || row[7] != "73344" {
		t.Errorf("Expected TX/73344, got %s/%s", row[6], row[7])
	}
}

func TestPopulateStateZipFromAreaCode_UnknownAreaCode(t *testing.T) {
	row := make([]string, 13)
	row[8] = "0009998888"
	populateStateZipFromAreaCode(row)
	if row[6] != "" || row[7] != "" {
		t.Errorf("Expected no inference, got %s/%s", row[6], row[7])
	}
}
