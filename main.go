package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"golang.org/x/term"
)

var (
	version = "dev"
	commit  = "none"
)

func main() {
	// parse flags
	theme := "default"
	output := ""
	for i, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("gitcredits %s (%s)\n", version, commit)
			os.Exit(0)
		}
		if arg == "--theme" && i+1 < len(os.Args[1:]) {
			theme = os.Args[i+2]
		}
		if arg == "--output" && i+1 < len(os.Args[1:]) {
			output = os.Args[i+2]
		}
	}

	info := getRepoInfo()

	width := 80
	height := 24
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w
		height = h
	}

	// GIF output mode
	if output != "" {
		credits := buildCredits(info, 80)
		cards := buildMatrixCards(info, 80, 24)
		if err := generateGIF(output, theme, credits, len(cards)); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("GIF saved: %s\n", output)
		return
	}

	var m model

	switch theme {
	case "matrix":
		cards := buildMatrixCards(info, width, height)
		m = model{
			height:  height,
			width:   width,
			theme:   theme,
			cards:   cards,
			cardIdx: 0,
			mState:  mvsRain,
		}
		m.initRain()
	default:
		credits := buildCredits(info, width)
		sf := newStarField(width, len(credits))
		m = model{
			lines:     credits,
			offset:    0,
			height:    height,
			width:     width,
			starField: sf,
			theme:     theme,
		}
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
