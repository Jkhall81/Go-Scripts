package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	completedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("10"))
	pendingStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("7"))
	titleStyle     = lipgloss.NewStyle().Bold(true)
	outputStyle    = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(0, 1)
	focusedOutputStyle = outputStyle.Copy().BorderForeground(lipgloss.Color("12"))
	inputStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("12"))
)

func (m model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	m.steps[0].status = true

	checklist := m.renderChecklist()
	legend := m.renderLegend()
	output := m.renderOutput()
	input := m.renderInput()

	// Top section: Legend on left, Checklist on right
	topSection := lipgloss.JoinHorizontal(
		lipgloss.Top,
		legend,
		lipgloss.NewStyle().Width(m.width-lipgloss.Width(legend)-30).Render(""), // Spacer
		checklist,
	)

	mainContent := lipgloss.JoinVertical(
		lipgloss.Left,
		topSection,
		"", // Add some spacing
		output,
		"", // Add some spacing
		input,
	)

	return lipgloss.NewStyle().
		Padding(1).
		Render(mainContent)
}

func (m model) renderChecklist() string {
	var steps []string
	steps = append(steps, titleStyle.Render("ETL Pipeline Progress:"))
	steps = append(steps, "")

	for _, step := range m.steps {
		status := "○"
		style := pendingStyle
		if step.status {
			status = "●"
			style = completedStyle
		}
		steps = append(steps, style.Render(fmt.Sprintf("%s %s", status, step.name)))
	}

	return lipgloss.NewStyle().
		Width(30).
		Padding(0, 1).
		Render(lipgloss.JoinVertical(lipgloss.Left, steps...))
}

func (m model) renderOutput() string {
	// Build all output content
	allContent := ""
	start := len(m.outputLines) - 30
	if start < 0 {
		start = 0
	}
	for _, line := range m.outputLines[start:] {
		allContent += line + "\n"
	}

	// Set content to scroll model and get scrolled view
	m.scroll.SetContent(allContent)
	scrolledContent := m.scroll.View()

	// Build the final content with title
	title := "Output Window"
	if m.focused == "output" {
		title = "▶ Output Window — Use ← → ↑ ↓ to scroll"
	} else {
		title = "Output Window — TAB to focus, then ← → ↑ ↓ to scroll"
	}

	content := titleStyle.Render(title)

	// Calculate the exact width for the border line (accounting for padding and borders)
	outputWidth := m.width - 12                      // Adjusted to prevent wrapping
	borderLine := strings.Repeat("─", outputWidth-4) // -4 to account for padding
	content += "\n" + borderLine + "\n\n"
	content += scrolledContent

	// Apply appropriate style based on focus
	style := outputStyle
	if m.focused == "output" {
		style = focusedOutputStyle
	}

	// Output window uses full width but with precise calculation
	return style.
		Width(outputWidth).
		Height(20).
		Render(content)
}

func (m model) renderLegend() string {
	legend := []string{
		"Legend:",
		"  show ............ preview first 5 rows",
		"  drop <indexes> .. remove columns",
		"  clean-address ... sanitize address fields",
		"  clean-names ..... remove numbers & special chars from names",
		"  clean-email ..... make sure there are no numeric values",
		"  clean-states .... make sure there are no numeric values or invalid strings",
		"  normalize-phones. format phone numbers",
		"  dedup-phones .... remove duplicate phone numbers",
		"  populate-geo .... fill missing geo fields",
		"  validate-states . drop non-US states",
		"  final-validate .. drop rows missing name/phone",
		"  clean-all ....... run entire automated pipeline",
		"  write-csv ....... export cleaned CSV",
		"  write-report .... summary report",
		"  exit ............ quit",
		"",
		"Navigation:",
		"  TAB ............ switch focus",
		"  ← → ............ scroll horizontally (when output focused)",
		"  ↑ ↓ ............ scroll vertically (when output focused)",
	}

	return lipgloss.JoinVertical(lipgloss.Left, legend...)
}

func (m model) renderInput() string {
	prompt := "> " + m.input
	if m.focused == "input" {
		prompt = "▶ " + m.input
	}
	return inputStyle.Render(prompt)
}
