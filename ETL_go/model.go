package main

import (
	"fmt"
	"log"
	"path/filepath"
	"strconv"
	"strings"

	"etl_go/extract"
	"etl_go/load"
	"etl_go/transform"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	steps       []step
	outputLines []string
	input       string
	width       int
	height      int
	dataset     *extract.DataSet
	focused     string
	scroll      scrollModel
}

type step struct {
	name   string
	status bool
}

func initialModel(inputFile string) model {
	steps := []step{
		{name: "Extract CSV/XLSX", status: false},
		{name: "Drop Columns", status: false},
		{name: "Clean Addresses", status: false},
		{name: "Clean Names", status: false},
		{name: "Normalize Phones", status: false},
		{name: "Populate Geo", status: false},
		{name: "Validate States", status: false},
		{name: "Final Validation", status: false},
		{name: "Write CSV", status: false},
		{name: "Write Report", status: false},
	}

	// Load the initial dataset
	var ds *extract.DataSet
	var err error
	ext := filepath.Ext(inputFile)

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

	outputLines := []string{
		fmt.Sprintf("Loaded %d rows from %s", len(ds.Rows), ds.Source),
	}

	m := model{
		steps:       steps,
		outputLines: outputLines,
		input:       "",
		dataset:     ds,
		focused:     "input",
		scroll:      newScrollModel(),
	}

	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

// In model.go, update the Update method to ensure we return the modified model:
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.focused == "input" {
				m.focused = "output"
			} else {
				m.focused = "input"
			}
			return m, nil
		case "enter":
			if m.focused == "input" {
				cmd := strings.TrimSpace(m.input)
				if cmd != "" {
					return m.processCommand(cmd)
				}
			}
			return m, nil
		case "backspace":
			if m.focused == "input" && len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
			return m, nil
		// Handle scrolling when output is focused
		case "left", "h":
			if m.focused == "output" {
				m.scroll.ScrollLeft()
				return m, nil
			} else {
				m.input += msg.String()
				return m, nil
			}
		case "right", "l":
			if m.focused == "output" {
				m.scroll.ScrollRight()
				return m, nil
			} else {
				m.input += msg.String()
				return m, nil
			}
		case "up", "k":
			if m.focused == "output" {
				m.scroll.ScrollUp()
				return m, nil
			} else {
				m.input += msg.String()
				return m, nil
			}
		case "down", "j":
			if m.focused == "output" {
				m.scroll.ScrollDown()
				return m, nil
			} else {
				m.input += msg.String()
				return m, nil
			}
		default:
			if m.focused == "input" {
				m.input += msg.String()
			}
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Set scroll area size (output window dimensions)
		outputWidth := m.width - 35
		if outputWidth < 20 {
			outputWidth = 20
		}
		outputHeight := 20
		m.scroll.SetSize(outputWidth, outputHeight)
		return m, nil
	}

	return m, nil
}

func (m model) processCommand(cmd string) (tea.Model, tea.Cmd) {
	m.outputLines = []string{}
	m.outputLines = append(m.outputLines, fmt.Sprintf("> %s", cmd))

	args := strings.Fields(cmd)
	if len(args) == 0 {
		return m, nil
	}

	switch strings.ToLower(args[0]) {
	case "show":
		previewLines := m.dataset.FirstNLines(5, m.width-35)
		m.outputLines = append(m.outputLines, previewLines...)

	case "drop":
		if len(args) < 2 {
			m.outputLines = append(m.outputLines, "Usage: drop <colIndex1> <colIndex2> ...")
		} else {
			indexes := []int{}
			for _, s := range args[1:] {
				i, err := strconv.Atoi(s)
				if err == nil {
					indexes = append(indexes, i)
				}
			}
			m.dataset = transform.DropColumns(m.dataset, indexes)
			m.steps[1].status = true
			m.outputLines = append(m.outputLines, fmt.Sprintf("Dropped columns: %v", indexes))
		}

	case "clean-address":
		m.dataset = transform.CleanAddresses(m.dataset)
		m.steps[2].status = true
		m.outputLines = append(m.outputLines, "Cleaned address fields.")

	case "clean-names":
		m.dataset = transform.CleanNames(m.dataset)
		m.steps[3].status = true
		m.outputLines = append(m.outputLines, "Cleaned name fields.")

	case "normalize-phones":
		m.dataset = transform.NormalizePhones(m.dataset)
		m.steps[4].status = true
		m.outputLines = append(m.outputLines, "Normalized phone numbers.")

	case "populate-geo":
		m.dataset = transform.PopulateGeo(m.dataset)
		m.steps[5].status = true
		m.outputLines = append(m.outputLines, "Populated missing state/ZIP data.")

	case "validate-states":
		result := transform.ValidateStates(m.dataset)
		m.dataset = result.Cleaned
		m.steps[6].status = true
		m.outputLines = append(m.outputLines, fmt.Sprintf("Removed %d invalid-state rows.", result.DropCount))

	case "final-validate":
		result := load.FinalValidate(m.dataset)
		m.dataset = result.Cleaned
		m.steps[7].status = true
		m.outputLines = append(m.outputLines, fmt.Sprintf("Removed %d invalid rows.", result.DropCount))

	case "write-csv":
		if err := load.WriteCSV(m.dataset, ""); err != nil {
			m.outputLines = append(m.outputLines, fmt.Sprintf("Error writing CSV: %v", err))
		} else {
			m.steps[8].status = true
			m.outputLines = append(m.outputLines, "Output CSV written successfully.")
		}

	case "write-report":
		report := load.ReportSummary{
			TotalProcessed: len(m.dataset.Rows),
			FinalRowCount:  len(m.dataset.Rows),
		}
		reportLines := load.WriteReport(report)
		m.outputLines = append(m.outputLines, reportLines...)
		m.steps[9].status = true

	case "help":
		m.outputLines = append(m.outputLines, "Available commands:")
		m.outputLines = append(m.outputLines, "show, drop, clean-address, normalize-phones, populate-geo,")
		m.outputLines = append(m.outputLines, "validate-states, final-validate, write-csv, write-report, exit")

	case "exit", "quit":
		return m, tea.Quit

	default:
		m.outputLines = append(m.outputLines, fmt.Sprintf("Unknown command: %s", cmd))
	}

	m.input = ""
	return m, nil
}

func (m model) removePreviewLines(lines []string) []string {
	var result []string
	inPreview := false

	for _, line := range lines {
		if strings.HasPrefix(line, "> show") {
			result = append(result, line)
			inPreview = true
			continue
		}

		if inPreview {
			if strings.Contains(line, "Loaded") && strings.Contains(line, "rows from") && len(result) > 0 {
				continue
			}
			if strings.Contains(line, "Showing") && strings.Contains(line, "rows from") {
				inPreview = false
				continue
			}
			if strings.HasPrefix(line, "|") || strings.HasPrefix(line, "+") {
				continue
			}
		}

		result = append(result, line)
	}
	return result
}
