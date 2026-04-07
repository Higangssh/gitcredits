package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func generateGIF(outputPath, theme, dir string, lines []string, cardCount int) error {
	vhsPath, err := exec.LookPath("vhs")
	if err != nil {
		return fmt.Errorf("vhs is required for GIF output. Install: brew install vhs")
	}

	selfPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("cannot find executable path: %w", err)
	}
	selfPath, _ = filepath.Abs(selfPath)

	var duration int
	switch theme {
	case "matrix":
		duration = cardCount*5 + 2
	default:
		duration = len(lines)/8 + 3
		if duration < 10 {
			duration = 10
		}
	}

	ffmpegPath, _ := exec.LookPath("ffmpeg")

	absOutput, err := filepath.Abs(outputPath)
	if err != nil {
		return fmt.Errorf("invalid output path: %w", err)
	}

	var vhsOutput string
	if ffmpegPath != "" {
		vhsOutput = absOutput + ".mp4"
	} else {
		vhsOutput = absOutput
	}

	var tape strings.Builder
	tape.WriteString(fmt.Sprintf("Output \"%s\"\n", vhsOutput))
	tape.WriteString("Set Width 960\n")
	tape.WriteString("Set Height 600\n")
	tape.WriteString("Set Padding 0\n")
	tape.WriteString("Set FontSize 16\n")
	tape.WriteString("Set Theme \"Builtin Dark\"\n")
	tape.WriteString("Set TypingSpeed 0\n")

	cmdParts := []string{selfPath}
	if theme != "default" {
		cmdParts = append(cmdParts, "--theme", theme)
	}
	if dir != "" {
		cmdParts = append(cmdParts, dir)
	}
	cmd := strings.Join(cmdParts, " ")
	tape.WriteString(fmt.Sprintf("Type %q\n", cmd))
	tape.WriteString("Enter\n")
	tape.WriteString(fmt.Sprintf("Sleep %ds\n", duration))

	tmpFile, err := os.CreateTemp("", "gitcredits-*.tape")
	if err != nil {
		return fmt.Errorf("cannot create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.WriteString(tape.String()); err != nil {
		tmpFile.Close()
		return fmt.Errorf("cannot write tape: %w", err)
	}
	tmpFile.Close()

	vhsCmd := exec.Command(vhsPath, tmpPath)
	if dir != "" {
		vhsCmd.Dir = dir
	} else {
		cwd, _ := os.Getwd()
		vhsCmd.Dir = cwd
	}
	vhsCmd.Env = append(os.Environ(), "TERM=xterm-256color")

	vhsOut, err := vhsCmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("vhs failed: %s\n%s", err, string(vhsOut))
	}

	if ffmpegPath != "" {
		defer os.Remove(vhsOutput)
		palettePath := absOutput + ".palette.png"
		defer os.Remove(palettePath)

		p1 := exec.Command(ffmpegPath, "-y", "-i", vhsOutput,
			"-vf", "fps=25,palettegen=max_colors=256:stats_mode=diff",
			palettePath)
		if out, err := p1.CombinedOutput(); err != nil {
			return fmt.Errorf("ffmpeg palette failed: %s\n%s", err, string(out))
		}

		p2 := exec.Command(ffmpegPath, "-y", "-i", vhsOutput, "-i", palettePath,
			"-filter_complex", "fps=25[v];[v][1:v]paletteuse=dither=floyd_steinberg",
			absOutput)
		if out, err := p2.CombinedOutput(); err != nil {
			return fmt.Errorf("ffmpeg gif failed: %s\n%s", err, string(out))
		}
	}

	return nil
}
