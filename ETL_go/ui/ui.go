package ui

import (
	"fmt"
	"log"
	"strings"

	"github.com/gdamore/tcell/v2"
)

// Step names (right-side checklist)
var Steps = []string{
	"Extract CSV/XLSX",
	"Drop Columns",
	"Clean Addresses",
	"Normalize Phones",
	"Populate Geo",
	"Validate States",
	"Final Validation",
	"Write CSV",
	"Write Report",
}

var outputLines []string
var scrollXOffset int = 0

// ----------------------
// Main DrawUI entrypoint
// ----------------------
func DrawUI(status map[string]bool, onCommand func(cmd string)) {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Error creating screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Error initializing screen: %v", err)
	}
	defer screen.Fini()

	input := []rune{} // command input buffer
	drawAll(screen, status, string(input))

	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {

		case *tcell.EventResize:
			screen.Sync()
			drawAll(screen, status, string(input))

		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyEscape:
				return

			// --- Horizontal scrolling ---
			case tcell.KeyLeft:
				scrollBy(-1)
				drawAll(screen, status, string(input))

			case tcell.KeyRight:
				scrollBy(1)
				drawAll(screen, status, string(input))

			case tcell.KeyEnter:
				cmd := strings.TrimSpace(string(input))
				if cmd != "" {
					onCommand(cmd)
					input = []rune{}
				}
				drawAll(screen, status, string(input))

			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(input) > 0 {
					input = input[:len(input)-1]
				}
				drawAll(screen, status, string(input))

			default:
				r := ev.Rune()
				if r != 0 {
					input = append(input, r)
				}
				drawAll(screen, status, string(input))
			}
		}
	}
}

// ----------------------
// SCROLL MANAGEMENT
// ----------------------

func scrollBy(direction int) {
	const scrollStep = 20 // roughly one column width
	if direction < 0 && scrollXOffset > 0 {
		scrollXOffset -= scrollStep
		if scrollXOffset < 0 {
			scrollXOffset = 0
		}
	} else if direction > 0 {
		scrollXOffset += scrollStep
		maxLen := 0
		for _, l := range outputLines {
			if len(l) > maxLen {
				maxLen = len(l)
			}
		}
		if scrollXOffset > maxLen-scrollStep {
			scrollXOffset = maxLen - scrollStep
		}
		if scrollXOffset < 0 {
			scrollXOffset = 0
		}
	}
}

// ----------------------
// DRAW FUNCTIONS
// ----------------------

func drawAll(s tcell.Screen, status map[string]bool, input string) {
	s.Clear()
	drawChecklist(s, status)
	drawOutputWindow(s)
	drawLegend(s)
	drawInputBar(s, input)
	s.Show()
}

// --- Output window (bottom-middle)
func drawOutputWindow(s tcell.Screen) {
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite)
	w, h := s.Size()

	startX := 2
	width := w - 4
	height := 18

	// Move the output window higher (more centered between checklist and legend)
	legendHeight := 10
	inputHeight := 3
	offsetFromBottom := legendHeight + inputHeight + 4 // push it up more
	startY := h - height - offsetFromBottom
	if startY < 1 {
		startY = 1
	}

	title := "[ Output Window — Use ← → to scroll horizontally ]"
	for i, ch := range title {
		s.SetContent(startX+i, startY, ch, nil, style.Bold(true))
	}

	// Borders
	for x := startX; x <= startX+width; x++ {
		s.SetContent(x, startY+1, '─', nil, style)
		s.SetContent(x, startY+height, '─', nil, style)
	}
	for y := startY + 1; y <= startY+height; y++ {
		s.SetContent(startX, y, '│', nil, style)
		s.SetContent(startX+width, y, '│', nil, style)
	}

	// --- compute max content width in runes ---
	maxContentWidth := 0
	for _, line := range outputLines {
		runes := []rune(line)
		if len(runes) > maxContentWidth {
			maxContentWidth = len(runes)
		}
	}

	// Clamp scroll offset
	if scrollXOffset < 0 {
		scrollXOffset = 0
	}
	if scrollXOffset > maxContentWidth-width {
		scrollXOffset = maxContentWidth - width
		if scrollXOffset < 0 {
			scrollXOffset = 0
		}
	}

	// --- Draw visible content ---
	y := startY + 2
	maxVisibleLines := height - 3
	endY := y + maxVisibleLines

	for _, line := range outputLines {
		if y >= endY {
			break
		}

		runes := []rune(line)
		if scrollXOffset < len(runes) {
			runes = runes[scrollXOffset:]
		} else {
			runes = []rune{}
		}

		for x, ch := range runes {
			if x >= width-2 {
				break
			}
			s.SetContent(startX+1+x, y, ch, nil, style)
		}
		y++
	}
}

// --- Legend (bottom-left)
func drawLegend(s tcell.Screen) {
	lines := []string{
		"Legend:",
		"  show ............ preview first 5 rows",
		"  drop <indexes> .. remove columns",
		"  clean-address ... sanitize address fields",
		"  normalize-phones. format phone numbers",
		"  populate-geo .... fill missing geo fields",
		"  validate-states . drop non-US states",
		"  final-validate .. drop rows missing name/phone",
		"  write-csv ....... export cleaned CSV",
		"  write-report .... summary report",
		"  exit ............ quit",
	}

	_, h := s.Size()
	y := h - len(lines) - 3 // room for input bar

	for i, line := range lines {
		for x, ch := range line {
			s.SetContent(x+2, y+i, ch, nil, tcell.StyleDefault)
		}
	}
}

// --- Checklist (right side)
func drawChecklist(s tcell.Screen, status map[string]bool) {
	w, _ := s.Size()
	startX := w - 30

	title := "ETL Pipeline Progress:"
	for i, ch := range title {
		s.SetContent(startX+i, 1, ch, nil, tcell.StyleDefault.Bold(true))
	}

	for i, step := range Steps {
		color := tcell.ColorWhite
		if status[step] {
			color = tcell.ColorGreen
		}
		style := tcell.StyleDefault.Foreground(color)
		for j, ch := range step {
			s.SetContent(startX+j, i+3, ch, nil, style)
		}
	}
}

// --- Input bar (bottom)
func drawInputBar(s tcell.Screen, input string) {
	_, h := s.Size()
	prompt := fmt.Sprintf("> %s", input)
	style := tcell.StyleDefault.Foreground(tcell.ColorWhite).Background(tcell.ColorBlack)

	for x, ch := range prompt {
		s.SetContent(x+2, h-2, ch, nil, style)
	}
}

// ----------------------
// Output buffer utilities
// ----------------------

// AddToOutput appends a line to the on-screen output window
func AddToOutput(line string) {
	const maxLines = 30
	line = strings.TrimRight(line, "\n")
	outputLines = append(outputLines, line)
	if len(outputLines) > maxLines {
		outputLines = outputLines[len(outputLines)-maxLines:]
	}
}
