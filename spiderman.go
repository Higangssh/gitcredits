package main

import (
	"fmt"
	"math/rand"
	"strings"
)

// Spider-Man hero titles based on rank
func spiderTitle(rank int, commits int) string {
	if rank == 0 {
		if commits >= 100 {
			return "THE SPIDER-MAN"
		}
		return "THE WEB-SLINGER"
	}
	titles := []string{
		"SPIDER-GWEN",
		"MILES MORALES",
		"SPIDER-NOIR",
		"SP//DR",
		"SPIDER-HAM",
		"SPIDER-PUNK",
		"SPIDER-WOMAN",
		"SCARLET SPIDER",
	}
	idx := (rank - 1) % len(titles)
	return titles[idx]
}

// Glitch characters
var glitchChars = []rune("█▓▒░▀▄▌▐╔╗╚╝═║╬╣╠╩╦┃━┏┓┗┛")

// glitchLine applies random glitch distortion to a string
func glitchLine(s string, intensity float64) string {
	runes := []rune(s)
	result := make([]rune, len(runes))
	for i, r := range runes {
		if r == ' ' {
			result[i] = r
			continue
		}
		if rand.Float64() < intensity {
			result[i] = glitchChars[rand.Intn(len(glitchChars))]
		} else {
			result[i] = r
		}
	}
	return string(result)
}

// rgbShift shifts text left/right to simulate chromatic aberration
func rgbShift(s string, offset int) (string, string) {
	runes := []rune(s)
	red := make([]rune, len(runes))
	blue := make([]rune, len(runes))

	for i := range runes {
		red[i] = ' '
		blue[i] = ' '
	}

	for i, r := range runes {
		if r == ' ' {
			continue
		}
		// Red channel: shift left
		ri := i - offset
		if ri >= 0 && ri < len(red) {
			red[ri] = r
		}
		// Blue channel: shift right
		bi := i + offset
		if bi >= 0 && bi < len(blue) {
			blue[bi] = r
		}
	}

	return string(red), string(blue)
}

// Web pattern for background
var webChars = []rune("·.·.··")

type webField struct {
	webs []struct {
		x, y int
		ch   rune
	}
}

func newWebField(width, totalHeight int) webField {
	wf := webField{}
	// Very sparse — just faint dots, not dense patterns
	density := (width * totalHeight) / 200
	for i := 0; i < density; i++ {
		ch := webChars[rand.Intn(len(webChars))]
		wf.webs = append(wf.webs, struct {
			x, y int
			ch   rune
		}{
			x:  rand.Intn(width),
			y:  rand.Intn(totalHeight),
			ch: ch,
		})
	}
	return wf
}

func buildSpidermanCards(info repoInfo, width, height int) []matrixCard {
	var cards []matrixCard

	center := func(s string) string {
		return centerText(s, width)
	}

	makeCard := func(content []string) matrixCard {
		lines := make([]string, height)
		startY := (height - len(content)) / 2
		if startY < 0 {
			startY = 0
		}
		for i, line := range content {
			if startY+i < height {
				lines[startY+i] = line
			}
		}
		return matrixCard{lines: lines}
	}

	// Card 0: Title
	var titleContent []string
	titleContent = append(titleContent, center("━━━━━━━━━━━━━━━━━━━━"))
	titleContent = append(titleContent, "")
	titleRows := bigText(info.name)
	for _, row := range titleRows {
		titleContent = append(titleContent, center(row))
	}
	titleContent = append(titleContent, "")
	if info.description != "" {
		titleContent = append(titleContent, center("\""+info.description+"\""))
		titleContent = append(titleContent, "")
	}
	if info.language != "" || info.stars > 0 {
		var meta []string
		if info.language != "" {
			meta = append(meta, info.language)
		}
		if info.stars > 0 {
			meta = append(meta, fmt.Sprintf("★ %d stars", info.stars))
		}
		titleContent = append(titleContent, center("· "+strings.Join(meta, " · ")+" ·"))
	}
	titleContent = append(titleContent, "")
	titleContent = append(titleContent, center("━━━━━━━━━━━━━━━━━━━━"))
	cards = append(cards, makeCard(titleContent))

	// Contributor cards
	for i, c := range info.contributors {
		title := spiderTitle(i, c.commits)
		var content []string
		content = append(content, center("━━━━━━━━━━━━━━━━━━━━"))
		content = append(content, "")
		content = append(content, center(title))
		content = append(content, "")
		content = append(content, center(strings.ToUpper(c.name)))
		content = append(content, "")
		content = append(content, center(fmt.Sprintf("%d webs spun", c.commits)))
		content = append(content, "")
		content = append(content, center("━━━━━━━━━━━━━━━━━━━━"))
		cards = append(cards, makeCard(content))
	}

	// Notable commits card
	if len(info.highlights) > 0 {
		var hlContent []string
		hlContent = append(hlContent, center("━━━━━━━━━━━━━━━━━━━━"))
		hlContent = append(hlContent, "")
		hlContent = append(hlContent, center("N O T A B L E   C O M M I T S"))
		hlContent = append(hlContent, "")
		for _, h := range info.highlights {
			hlContent = append(hlContent, center("· "+h))
			hlContent = append(hlContent, "")
		}
		hlContent = append(hlContent, center("━━━━━━━━━━━━━━━━━━━━"))
		cards = append(cards, makeCard(hlContent))
	}

	// Final card
	totalCommits := 0
	for _, c := range info.contributors {
		totalCommits += c.commits
	}
	// Stats card
	var statsContent []string
	statsContent = append(statsContent, center("━━━━━━━━━━━━━━━━━━━━"))
	statsContent = append(statsContent, "")
	statsContent = append(statsContent, center(fmt.Sprintf("%d  C O M M I T S", totalCommits)))
	statsContent = append(statsContent, "")
	statsContent = append(statsContent, center(fmt.Sprintf("%d  C O N T R I B U T O R S", len(info.contributors))))
	if info.stars > 0 {
		statsContent = append(statsContent, "")
		statsContent = append(statsContent, center(fmt.Sprintf("★  %d  S T A R S  ★", info.stars)))
	}
	if info.language != "" {
		statsContent = append(statsContent, "")
		statsContent = append(statsContent, center("Written in "+info.language))
	}
	if info.license != "" {
		statsContent = append(statsContent, "")
		statsContent = append(statsContent, center("Licensed under "+info.license))
	}
	statsContent = append(statsContent, "")
	statsContent = append(statsContent, center("━━━━━━━━━━━━━━━━━━━━"))
	cards = append(cards, makeCard(statsContent))

	// Final card
	var finalContent []string
	finalContent = append(finalContent, "")
	finalContent = append(finalContent, "")
	finalContent = append(finalContent, center("With great power comes"))
	finalContent = append(finalContent, center("great responsibility"))
	finalContent = append(finalContent, "")
	finalContent = append(finalContent, "")
	cards = append(cards, makeCard(finalContent))

	return cards
}
