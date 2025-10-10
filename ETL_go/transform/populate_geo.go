package transform

import (
	"fmt"
	"strconv"
	"strings"

	"etl_go/extract"
)

// --- STATE → ZIP ---
var stateZip = map[string]string{
	"AL": "35007", "AK": "99501", "AZ": "85304", "AR": "71602", "CA": "90005",
	"CO": "80001", "CT": "06001", "DE": "19701", "DC": "20012", "FL": "32003",
	"GA": "30002", "HI": "96701", "ID": "83203", "IL": "61081", "IN": "46011",
	"IA": "50005", "KS": "66008", "KY": "40007", "LA": "70001", "ME": "04750",
	"MD": "20601", "MA": "05544", "MI": "48706", "MN": "54403", "MS": "38601",
	"MO": "64722", "MT": "59001", "NE": "68001", "NV": "88905", "NH": "03031",
	"NJ": "07753", "NC": "28376", "NM": "87001", "NY": "10028", "ND": "58001",
	"OH": "45434", "OK": "73002", "OR": "97009", "PA": "15001", "RI": "02823",
	"SC": "29001", "SD": "57002", "TN": "37011", "TX": "73344", "TX2": "79901",
	"UT": "84002", "VT": "05009", "VA": "20101", "WA": "98001", "WV": "24712",
	"WI": "54990", "WY": "82002", "PR": "00999", "VI": "00851",
}

// --- ZIP → STATE RANGE ---
var zipCodeRanges = map[string][][2]int{
	"AL": {{35000, 36999}}, "AK": {{99500, 99999}}, "AZ": {{85000, 86999}},
	"AR": {{71600, 72999}}, "CA": {{90000, 96699}}, "CO": {{80000, 81999}},
	"CT": {{6000, 6999}}, "DE": {{19700, 19999}}, "FL": {{32000, 34999}},
	"GA": {{30000, 31999}, {39800, 39999}}, "HI": {{96700, 96999}},
	"ID": {{83200, 83999}}, "IL": {{60000, 62999}}, "IN": {{46000, 47999}},
	"IA": {{50000, 52999}}, "KS": {{66000, 67999}}, "KY": {{40000, 42999}},
	"LA": {{70000, 71599}}, "ME": {{3900, 4999}}, "MD": {{20600, 21999}},
	"MA": {{1000, 2799}, {5501, 5544}}, "MI": {{48000, 49999}},
	"MN": {{55000, 56899}}, "MS": {{38600, 39999}}, "MO": {{63000, 65999}},
	"MT": {{59000, 59999}}, "NC": {{27000, 28999}}, "ND": {{58000, 58999}},
	"NE": {{68000, 69999}}, "NV": {{88900, 89999}}, "NH": {{3000, 3899}},
	"NJ": {{7000, 8999}}, "NM": {{87000, 88499}},
	"NY": {{10000, 14999}, {6390, 6390}, {501, 501}, {544, 544}},
	"OH": {{43000, 45999}}, "OK": {{73000, 74999}}, "OR": {{97000, 97999}},
	"PA": {{15000, 19699}}, "RI": {{2800, 2999}}, "SC": {{29000, 29999}},
	"SD": {{57000, 57999}}, "TN": {{37000, 38599}},
	"TX": {{75000, 79999}, {73301, 73399}, {88500, 88599}},
	"UT": {{84000, 84999}}, "VT": {{5000, 5999}}, "VA": {{20100, 24699}},
	"DC": {{20000, 20099}, {20200, 20599}, {56900, 56999}},
	"WA": {{98000, 99499}}, "WV": {{24700, 26999}}, "WI": {{53000, 54999}},
	"WY": {{82000, 83199}, {83414, 83414}}, "PR": {{600, 999}}, "VI": {{801, 851}},
}

// --- AREA CODE → STATE ---
var stateAreaCodes = map[string][]string{
	"AL": {"205", "251", "256", "334"},
	"AK": {"907"},
	"AZ": {"480", "520", "602", "623", "928"},
	"AR": {"479", "501", "870"},
	"CA": {"209", "213", "310", "323", "408", "415", "424", "442", "510", "559", "562", "619", "626", "650", "657", "661", "707", "714", "760", "805", "818", "831", "858", "909", "916", "925", "949", "951"},
	"CO": {"303", "719", "720", "970"},
	"CT": {"203", "475", "860", "959"},
	"DC": {"202"},
	"DE": {"302"},
	"FL": {"305", "321", "352", "386", "407", "561", "727", "754", "772", "786", "813", "850", "863", "904", "941", "954"},
	"GA": {"229", "404", "470", "478", "678", "706", "770", "912"},
	"HI": {"808"},
	"IL": {"217", "224", "309", "312", "331", "618", "630", "708", "773", "779", "815", "847"},
	"NY": {"315", "332", "347", "516", "518", "585", "607", "631", "646", "716", "718", "845", "914", "917", "929"},
	"TX": {"210", "214", "254", "281", "325", "346", "361", "409", "430", "432", "469", "512", "682", "713", "726", "737", "806", "817", "830", "832", "903", "936", "940", "956", "972", "979"},
}

// --- HELPERS ---
func normalizeState(state string) string {
	return strings.ToUpper(strings.TrimSpace(state))
}

func truncateZip(zip string) string {
	if len(zip) > 5 {
		return zip[:5]
	}
	return zip
}

func populateZip(row []string) {
	state := row[6]
	if state != "" && row[7] == "" {
		if zip, ok := stateZip[state]; ok {
			row[7] = zip
		}
	}
}

func populateStateFromZip(row []string) {
	zipStr := row[7]
	if zipStr == "" {
		return
	}
	zipInt, err := strconv.Atoi(zipStr)
	if err != nil {
		return
	}

	for state, ranges := range zipCodeRanges {
		for _, r := range ranges {
			if zipInt >= r[0] && zipInt <= r[1] {
				row[6] = state
				return
			}
		}
	}
}

func populateStateZipFromAreaCode(row []string) {
	if len(row[8]) < 3 {
		return
	}
	ac := row[8][:3]
	for state, codes := range stateAreaCodes {
		for _, c := range codes {
			if c == ac {
				row[6] = state
				row[7] = stateZip[state]
				return
			}
		}
	}
}

// --- MAIN TRANSFORM FUNCTION ---
func PopulateGeo(ds *extract.DataSet) *extract.DataSet {
	if ds == nil {
		fmt.Println("No dataset loaded.")
		return ds
	}

	newRows := make([][]string, len(ds.Rows))
	for i, row := range ds.Rows {
		newRow := make([]string, len(row))
		copy(newRow, row)

		newRow[6] = normalizeState(newRow[6])
		newRow[7] = truncateZip(newRow[7])

		if newRow[7] == "" && newRow[6] != "" {
			populateZip(newRow)
		}
		if newRow[6] == "" && newRow[7] != "" {
			populateStateFromZip(newRow)
		}
		if (newRow[6] == "" || len(newRow[6]) != 2) && newRow[7] == "" {
			populateStateZipFromAreaCode(newRow)
		}

		newRows[i] = newRow
	}

	fmt.Println("Geographic fields populated successfully.")
	return &extract.DataSet{
		Headers: ds.Headers,
		Rows:    newRows,
		Source:  ds.Source,
	}
}
