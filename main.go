package main

import (
	"fmt"
	"strings"
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
	dir := ""
	for i, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-v" {
			fmt.Printf("gitcredits %s (%s)\n", version, commit)
			os.Exit(0)
		}
		if arg == "--help" || arg == "-h" {
			fmt.Println("gitcredits - Turn your Git repo into movie-style rolling credits")
			fmt.Println()
			fmt.Printf("Usage: gitcredits [options] <directory>\n\n")
			fmt.Println("Arguments:")
			fmt.Println("  <directory>      Directory to look for a git repository (if not supplied, defaults to current directory)")
			fmt.Println()
			fmt.Println("Options:")
			fmt.Println("  --theme <name>   Theme: default, matrix, spiderman")
			fmt.Println("  --output <file>  Export credits as GIF")
			fmt.Println("  --version, -v    Show version")
			fmt.Println("  --help, -h       Show this help")
			os.Exit(0)
		}
		if arg == "--theme" && i+1 < len(os.Args[1:]) {
			theme = os.Args[i+2]
		}
		if arg == "--output" && i+1 < len(os.Args[1:]) {
			output = os.Args[i+2]
		}
	}

	if len(os.Args) > 1 {
		if !strings.HasPrefix(os.Args[len(os.Args)-2], "--") {
			dir = os.Args[len(os.Args)-1]
		}
	}

	info := getRepoInfo(dir)

	width := 80
	height := 24
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w
		height = h
	}

	// GIF output mode
	if output != "" {
		credits := buildCredits(info, 80)
		var cards []matrixCard
		switch theme {
		default:
			cards = buildMatrixCards(info, 80, 24)
		}
		if err := generateGIF(output, theme, dir, credits, len(cards)); err != nil {
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
	case "spiderman":
		cards := buildSpidermanCards(info, width, height)
		wf := newWebField(width, height*len(cards))
		m = model{
			height:   height,
			width:    width,
			theme:    theme,
			cards:    cards,
			cardIdx:  0,
			mState:   mvsRain,
			webField: wf,
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
