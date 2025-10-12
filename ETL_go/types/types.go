package types

// GeoStats holds geographic data cleaning statistics
type GeoStats struct {
	CleanedZipLetters    int
	CleanedZipTooShort   int
	PopulatedZip         int
	PopulatedState       int
	CorrectedMismatches  int
	FixedFromAreaCode    int
}