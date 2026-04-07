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

type config struct {
	theme  string
	output string
	dir    string
}

func main() {
	cfg, err := parseArgs(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	if cfg == nil {
		return
	}

	info, err := getRepoInfo(cfg.dir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	width := 80
	height := 24
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w
		height = h
	}

	if cfg.output != "" {
		credits := buildCredits(info, 80)
		var cards []matrixCard
		switch cfg.theme {
		default:
			cards = buildMatrixCards(info, 80, 24)
		}
		if err := generateGIF(cfg.output, cfg.theme, cfg.dir, credits, len(cards)); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("GIF saved: %s\n", cfg.output)
		return
	}

	var m model

	switch cfg.theme {
	case "matrix":
		cards := buildMatrixCards(info, width, height)
		m = model{
			height:  height,
			width:   width,
			theme:   cfg.theme,
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
			theme:    cfg.theme,
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
			theme:     cfg.theme,
		}
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func parseArgs(args []string) (*config, error) {
	cfg := &config{theme: "default"}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch arg {
		case "--version", "-v":
			fmt.Printf("gitcredits %s (%s)\n", version, commit)
			return nil, nil
		case "--help", "-h":
			printHelp()
			return nil, nil
		case "--theme":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("missing value for --theme")
			}
			cfg.theme = args[i]
		case "--output":
			i++
			if i >= len(args) {
				return nil, fmt.Errorf("missing value for --output")
			}
			cfg.output = args[i]
		default:
			if len(arg) > 0 && arg[0] == '-' {
				return nil, fmt.Errorf("unknown flag: %s", arg)
			}
			if cfg.dir != "" {
				return nil, fmt.Errorf("only one target directory can be provided")
			}
			cfg.dir = arg
		}
	}

	return cfg, nil
}

func printHelp() {
	fmt.Println("gitcredits - Turn your Git repo into movie-style rolling credits")
	fmt.Println()
	fmt.Printf("Usage: gitcredits [options] [directory]\n\n")
	fmt.Println("Arguments:")
	fmt.Println("  directory       Target git repository directory (defaults to current directory)")
	fmt.Println()
	fmt.Println("Options:")
	fmt.Println("  --theme <name>   Theme: default, matrix, spiderman")
	fmt.Println("  --output <file>  Export credits as GIF")
	fmt.Println("  --version, -v    Show version")
	fmt.Println("  --help, -h       Show this help")
}
