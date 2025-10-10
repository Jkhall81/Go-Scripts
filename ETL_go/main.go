package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: etl_go <inputfile.csv | inputfile.xlsx>")
		os.Exit(1)
	}

	inputFile := os.Args[1]
	p := tea.NewProgram(initialModel(inputFile), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error: %v", err)
		os.Exit(1)
	}
}
