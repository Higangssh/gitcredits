package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"sort"
	"strconv"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type repoInfo struct {
	name         string
	description  string
	totalCommits int
	contributors []contributor
	highlights   []string
	stars        int
	license      string
	language     string
}

type contributor struct {
	name    string
	commits int
}

func getRepoInfo() repoInfo {
	info := repoInfo{}

	if dir, err := os.Getwd(); err == nil {
		parts := strings.Split(dir, string(os.PathSeparator))
		info.name = parts[len(parts)-1]
	}

	if desc, err := os.ReadFile(".git/description"); err == nil {
		d := strings.TrimSpace(string(desc))
		if d != "" && d != "Unnamed repository; edit this file 'description' to name the repository." {
			info.description = d
		}
	}

	if info.description == "" {
		if out, err := exec.Command("gh", "repo", "view", "--json", "description", "-q", ".description").Output(); err == nil {
			d := strings.TrimSpace(string(out))
			if d != "" {
				info.description = d
			}
		}
	}

	if out, err := exec.Command("git", "rev-list", "--count", "HEAD").Output(); err == nil {
		if n, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.totalCommits = n
		}
	}

	if out, err := exec.Command("git", "shortlog", "-sn", "--no-merges", "HEAD").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, "\t", 2)
			if len(parts) == 2 {
				n, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
				info.contributors = append(info.contributors, contributor{
					name:    strings.TrimSpace(parts[1]),
					commits: n,
				})
			}
		}
	}

	sort.Slice(info.contributors, func(i, j int) bool {
		return info.contributors[i].commits > info.contributors[j].commits
	})

	if out, err := exec.Command("git", "log", "--oneline", "--no-merges", "-50", "--format=%s").Output(); err == nil {
		lines := strings.Split(strings.TrimSpace(string(out)), "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if strings.HasPrefix(line, "feat:") || strings.HasPrefix(line, "fix:") {
				msg := line
				if strings.HasPrefix(line, "feat: ") {
					msg = line[6:]
				} else if strings.HasPrefix(line, "fix: ") {
					msg = line[5:]
				}
				info.highlights = append(info.highlights, msg)
				if len(info.highlights) >= 8 {
					break
				}
			}
		}
	}

	if out, err := exec.Command("gh", "repo", "view", "--json", "stargazerCount", "-q", ".stargazerCount").Output(); err == nil {
		if n, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.stars = n
		}
	}

	if out, err := exec.Command("gh", "repo", "view", "--json", "licenseInfo", "-q", ".licenseInfo.name").Output(); err == nil {
		l := strings.TrimSpace(string(out))
		if l != "" {
			info.license = l
		}
	}

	if out, err := exec.Command("gh", "repo", "view", "--json", "primaryLanguage", "-q", ".primaryLanguage.name").Output(); err == nil {
		l := strings.TrimSpace(string(out))
		if l != "" {
			info.language = l
		}
	}

	return info
}

// simple block letter generator for title
func bigText(s string) []string {
	letters := map[rune][]string{
		'A': {"  ██  ", " █  █ ", " ████ ", " █  █ ", " █  █ "},
		'B': {" ███  ", " █  █ ", " ███  ", " █  █ ", " ███  "},
		'C': {"  ███ ", " █    ", " █    ", " █    ", "  ███ "},
		'D': {" ███  ", " █  █ ", " █  █ ", " █  █ ", " ███  "},
		'E': {" ████ ", " █    ", " ███  ", " █    ", " ████ "},
		'F': {" ████ ", " █    ", " ███  ", " █    ", " █    "},
		'G': {"  ███ ", " █    ", " █ ██ ", " █  █ ", "  ███ "},
		'H': {" █  █ ", " █  █ ", " ████ ", " █  █ ", " █  █ "},
		'I': {" ███ ", "  █  ", "  █  ", "  █  ", " ███ "},
		'J': {"  ███ ", "    █ ", "    █ ", " █  █ ", "  ██  "},
		'K': {" █  █ ", " █ █  ", " ██   ", " █ █  ", " █  █ "},
		'L': {" █    ", " █    ", " █    ", " █    ", " ████ "},
		'M': {" █   █ ", " ██ ██ ", " █ █ █ ", " █   █ ", " █   █ "},
		'N': {" █   █ ", " ██  █ ", " █ █ █ ", " █  ██ ", " █   █ "},
		'O': {"  ██  ", " █  █ ", " █  █ ", " █  █ ", "  ██  "},
		'P': {" ███  ", " █  █ ", " ███  ", " █    ", " █    "},
		'Q': {"  ██  ", " █  █ ", " █  █ ", " █ █  ", "  █ █ "},
		'R': {" ███  ", " █  █ ", " ███  ", " █ █  ", " █  █ "},
		'S': {"  ███ ", " █    ", "  ██  ", "    █ ", " ███  "},
		'T': {" █████ ", "   █   ", "   █   ", "   █   ", "   █   "},
		'U': {" █  █ ", " █  █ ", " █  █ ", " █  █ ", "  ██  "},
		'V': {" █  █ ", " █  █ ", " █  █ ", "  ██  ", "  ██  "},
		'W': {" █   █ ", " █   █ ", " █ █ █ ", " ██ ██ ", " █   █ "},
		'X': {" █  █ ", " █  █ ", "  ██  ", " █  █ ", " █  █ "},
		'Y': {" █  █ ", " █  █ ", "  ██  ", "  █   ", "  █   "},
		'Z': {" ████ ", "   █  ", "  █   ", " █    ", " ████ "},
		'-': {"      ", "      ", " ──── ", "      ", "      "},
		' ': {"   ", "   ", "   ", "   ", "   "},
		'_': {"      ", "      ", "      ", "      ", " ████ "},
	}

	upper := strings.ToUpper(s)
	rows := make([]string, 5)

	for _, ch := range upper {
		letter, ok := letters[ch]
		if !ok {
			letter = letters[' ']
		}
		for row := 0; row < 5; row++ {
			rows[row] += letter[row]
		}
	}

	return rows
}

func buildCredits(info repoInfo, width int) []string {
	var lines []string

	center := func(s string) string {
		runeLen := len([]rune(s))
		if runeLen >= width {
			return s
		}
		pad := (width - runeLen) / 2
		return strings.Repeat(" ", pad) + s
	}

	blank := func(n int) {
		for i := 0; i < n; i++ {
			lines = append(lines, "")
		}
	}

	// opening - full screen blank
	blank(20)

	// big ASCII title
	titleRows := bigText(info.name)
	for _, row := range titleRows {
		lines = append(lines, center(row))
	}

	blank(2)

	// description
	if info.description != "" {
		lines = append(lines, center("\""+info.description+"\""))
	}

	blank(6)

	// a film by...
	if len(info.contributors) > 0 {
		lines = append(lines, center("A   P R O J E C T   B Y"))
		blank(2)
		lines = append(lines, center(strings.ToUpper(info.contributors[0].name)))
		blank(1)
		lines = append(lines, center(fmt.Sprintf("— %d commits —", info.contributors[0].commits)))
	}

	blank(6)

	// starring (contributors)
	if len(info.contributors) > 1 {
		lines = append(lines, center("S T A R R I N G"))
		blank(2)
		for _, c := range info.contributors[1:] {
			lines = append(lines, center(strings.ToUpper(c.name)))
			lines = append(lines, center(fmt.Sprintf("%d commits", c.commits)))
			blank(1)
		}
	}

	blank(5)

	// highlights as "scenes"
	if len(info.highlights) > 0 {
		lines = append(lines, center("N O T A B L E   S C E N E S"))
		blank(2)
		for _, h := range info.highlights {
			lines = append(lines, center("· "+h+" ·"))
			blank(1)
		}
	}

	blank(5)

	// production stats
	lines = append(lines, center("━━━━━━━━━━━━━━━━━━━━"))
	blank(2)
	lines = append(lines, center(fmt.Sprintf("%d  C O M M I T S", info.totalCommits)))
	blank(1)
	lines = append(lines, center(fmt.Sprintf("%d  C O N T R I B U T O R S", len(info.contributors))))
	if info.stars > 0 {
		blank(1)
		lines = append(lines, center(fmt.Sprintf("★  %d  S T A R G A Z E R S  ★", info.stars)))
	}
	blank(2)
	if info.language != "" {
		lines = append(lines, center("Written in "+info.language))
	}
	if info.license != "" {
		lines = append(lines, center("Licensed under "+info.license))
	}
	blank(2)
	lines = append(lines, center("━━━━━━━━━━━━━━━━━━━━"))

	blank(6)

	// the end
	endRows := bigText("THE END")
	for _, row := range endRows {
		lines = append(lines, center(row))
	}

	blank(20)

	return lines
}

// star field background
type starField struct {
	stars []struct {
		x, y int
		ch   rune
	}
}

func newStarField(width, totalHeight int) starField {
	sf := starField{}
	density := (width * totalHeight) / 40  // denser star field
	for i := 0; i < density; i++ {
		ch := '·'
		bright := rand.Intn(10)
		if bright == 0 {
			ch = '✦'  // bright star
		} else if bright <= 2 {
			ch = '✧'  // medium star
		} else if bright <= 4 {
			ch = '⋆'  // small star
		} else if bright <= 6 {
			ch = '·'  // dot
		} else {
			ch = '.'  // faint
		}
		sf.stars = append(sf.stars, struct {
			x, y int
			ch   rune
		}{
			x:  rand.Intn(width),
			y:  rand.Intn(totalHeight),
			ch: ch,
		})
	}
	return sf
}

// Bubble Tea model

type tickMsg struct{}

type model struct {
	lines     []string
	offset    int
	height    int
	width     int
	done      bool
	starField starField
}

func (m model) Init() tea.Cmd {
	return tea.Tick(120*time.Millisecond, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.done = true
			return m, tea.Quit
		case "up":
			m.offset -= 3
			if m.offset < 0 {
				m.offset = 0
			}
			return m, nil
		case "down":
			m.offset += 3
			return m, nil
		}

	case tea.WindowSizeMsg:
		m.height = msg.Height
		m.width = msg.Width
		return m, nil

	case tickMsg:
		m.offset++
		if m.offset > len(m.lines) {
			m.done = true
			return m, tea.Quit
		}
		return m, tea.Tick(120*time.Millisecond, func(_ time.Time) tea.Msg {
			return tickMsg{}
		})
	}

	return m, nil
}

func (m model) View() string {
	if m.done {
		return ""
	}

	// color palette — cinematic night sky
	title := lipgloss.NewStyle().Foreground(lipgloss.Color("#00BFFF")).Bold(true)  // deep sky blue, bold
	silver := lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("#666666"))
	dimmer := lipgloss.NewStyle().Foreground(lipgloss.Color("#444444"))
	bright := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("#87CEEB"))  // sky blue
	scene := lipgloss.NewStyle().Foreground(lipgloss.Color("#B0C4DE"))   // light steel blue
	starBright := lipgloss.NewStyle().Foreground(lipgloss.Color("#8899AA"))
	gold := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFD700")).Bold(true)
	contributor := lipgloss.NewStyle().Foreground(lipgloss.Color("#FFFFFF")).Bold(true)

	var screenLines []string
	start := m.offset
	end := m.offset + m.height

	if start < 0 {
		start = 0
	}

	for i := start; i < end && i < len(m.lines); i++ {
		line := m.lines[i]
		screenIdx := i - start
		trimmed := strings.TrimSpace(line)

		// fade zones
		fadeTop := 4
		fadeBottom := 4
		distFromTop := screenIdx
		distFromBottom := m.height - 1 - screenIdx

		isFaded := distFromTop < fadeTop || distFromBottom < fadeBottom
		isVeryFaded := distFromTop < 2 || distFromBottom < 2

		var styled string
		if trimmed == "" {
			// render star background on blank lines
			starLine := make([]rune, m.width)
			for j := range starLine {
				starLine[j] = ' '
			}
			for _, s := range m.starField.stars {
				if s.y == i && s.x < m.width {
					starLine[s.x] = s.ch
				}
			}
			sl := string(starLine)
			if strings.TrimSpace(sl) != "" {
				styled = starBright.Render(sl)
			} else {
				styled = ""
			}
		} else if strings.Contains(trimmed, "██") {
			// ASCII art title / THE END
			if isVeryFaded {
				styled = dimmer.Render(line)
			} else if isFaded {
				styled = dim.Render(line)
			} else {
				styled = title.Render(line)
			}
		} else if strings.Contains(trimmed, "━") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = accent.Render(line)
			}
		} else if strings.Contains(trimmed, "★") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = gold.Render(line)
			}
		} else if strings.HasPrefix(trimmed, "A   P R O") || strings.HasPrefix(trimmed, "S T A R") ||
			strings.HasPrefix(trimmed, "N O T A B") {
			if isVeryFaded {
				styled = dimmer.Render(line)
			} else if isFaded {
				styled = dim.Render(line)
			} else {
				styled = bright.Render(line)
			}
		} else if strings.Contains(trimmed, "C O M M") || strings.Contains(trimmed, "C O N T R") ||
			strings.Contains(trimmed, "S T A R G") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = silver.Render(line)
			}
		} else if strings.Contains(trimmed, "\"") {
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = accent.Render(line)
			}
		} else if strings.Contains(trimmed, "· ") && strings.HasSuffix(trimmed, " ·") {
			// notable scenes items
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = scene.Render(line)
			}
		} else if strings.Contains(trimmed, "—") && strings.Contains(trimmed, "commits") {
			// commit count under name
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = accent.Render(line)
			}
		} else if strings.Contains(trimmed, "commits") && !strings.Contains(trimmed, "C O M") {
			// contributor commit count
			if isFaded {
				styled = dim.Render(line)
			} else {
				styled = accent.Render(line)
			}
		} else if trimmed == strings.ToUpper(trimmed) && len(trimmed) > 2 && !strings.Contains(trimmed, " O ") {
			// uppercase names (contributors)
			if isVeryFaded {
				styled = dimmer.Render(line)
			} else if isFaded {
				styled = dim.Render(line)
			} else {
				styled = contributor.Render(line)
			}
		} else {
			if isVeryFaded {
				styled = dimmer.Render(line)
			} else if isFaded {
				styled = dim.Render(line)
			} else {
				styled = silver.Render(line)
			}
		}

		screenLines = append(screenLines, styled)
	}

	// pad remaining lines
	for len(screenLines) < m.height {
		screenLines = append(screenLines, "")
	}

	return strings.Join(screenLines, "\n")
}

func main() {
	info := getRepoInfo()

	width := 80
	height := 24
	if w, h, err := term.GetSize(int(os.Stdout.Fd())); err == nil {
		width = w
		height = h
	}

	credits := buildCredits(info, width)
	sf := newStarField(width, len(credits))

	m := model{
		lines:     credits,
		offset:    0,
		height:    height,
		width:     width,
		starField: sf,
	}

	p := tea.NewProgram(m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
