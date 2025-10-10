package main

import (
	"strings"
)

type scrollModel struct {
	content         string
	scrollX         int
	scrollY         int
	width           int
	height          int
	maxContentWidth int
}

func newScrollModel() scrollModel {
	return scrollModel{
		scrollX: 0,
		scrollY: 0,
		width:   80,
		height:  20,
	}
}

// All methods must have pointer receivers to modify the actual struct
func (s *scrollModel) SetContent(content string) {
	s.content = content
	// Calculate max line width for horizontal scrolling
	s.maxContentWidth = 0
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if len(line) > s.maxContentWidth {
			s.maxContentWidth = len(line)
		}
	}
	s.ensureBounds()
}

func (s *scrollModel) SetSize(width, height int) {
	s.width = width
	s.height = height
	s.ensureBounds()
}

func (s *scrollModel) ScrollLeft() {
	s.scrollX -= 5
	if s.scrollX < 0 {
		s.scrollX = 0
	}
}

func (s *scrollModel) ScrollRight() {
	s.scrollX += 5
	s.ensureBounds()
}

func (s *scrollModel) ScrollUp() {
	s.scrollY -= 1
	if s.scrollY < 0 {
		s.scrollY = 0
	}
}

func (s *scrollModel) ScrollDown() {
	s.scrollY += 1
	s.ensureBounds()
}

func (s *scrollModel) ensureBounds() {
	// Horizontal bounds
	if s.maxContentWidth > s.width && s.scrollX > s.maxContentWidth-s.width {
		s.scrollX = s.maxContentWidth - s.width
	} else if s.scrollX < 0 {
		s.scrollX = 0
	}

	// Vertical bounds
	lineCount := strings.Count(s.content, "\n") + 1
	if s.scrollY > lineCount-s.height {
		s.scrollY = lineCount - s.height
	}
	if s.scrollY < 0 {
		s.scrollY = 0
	}
}

// View can stay with value receiver since it doesn't modify state
func (s scrollModel) View() string {
	if s.content == "" {
		return ""
	}

	lines := strings.Split(s.content, "\n")

	// Apply vertical scrolling
	startLine := s.scrollY
	endLine := startLine + s.height
	if endLine > len(lines) {
		endLine = len(lines)
	}
	if startLine < 0 {
		startLine = 0
	}

	var visibleLines []string
	for i := startLine; i < endLine && i < len(lines); i++ {
		line := lines[i]

		// Apply horizontal scrolling
		if s.scrollX < len(line) {
			line = line[s.scrollX:]
		} else {
			line = ""
		}

		// Trim to width
		if len(line) > s.width {
			line = line[:s.width]
		}

		visibleLines = append(visibleLines, line)
	}

	return strings.Join(visibleLines, "\n")
}
