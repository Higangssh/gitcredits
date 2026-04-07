package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
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

func getRepoInfo(dir string) (repoInfo, error) {
	info := repoInfo{}
	repoDir := dir
	if repoDir == "" {
		wd, err := os.Getwd()
		if err != nil {
			return info, fmt.Errorf("get current directory: %w", err)
		}
		repoDir = wd
	}

	absRepoDir, err := filepath.Abs(repoDir)
	if err != nil {
		return info, fmt.Errorf("resolve repository path %q: %w", repoDir, err)
	}

	stat, err := os.Stat(absRepoDir)
	if err != nil {
		return info, fmt.Errorf("access repository path %q: %w", absRepoDir, err)
	}
	if !stat.IsDir() {
		return info, fmt.Errorf("repository path %q is not a directory", absRepoDir)
	}

	info.name = filepath.Base(absRepoDir)

	if desc, err := os.ReadFile(filepath.Join(absRepoDir, ".git", "description")); err == nil {
		d := strings.TrimSpace(string(desc))
		if d != "" && d != "Unnamed repository; edit this file 'description' to name the repository." {
			info.description = d
		}
	}

	if info.description == "" {
		if out, err := runCommand(absRepoDir, "gh", "repo", "view", "--json", "description", "-q", ".description"); err == nil {
			d := strings.TrimSpace(string(out))
			if d != "" {
				info.description = d
			}
		}
	}

	if out, err := runCommand(absRepoDir, "git", "rev-list", "--count", "HEAD"); err == nil {
		if n, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.totalCommits = n
		}
	}

	if out, err := runCommand(absRepoDir, "git", "shortlog", "-sn", "--no-merges", "HEAD"); err == nil {
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

	if out, err := runCommand(absRepoDir, "git", "log", "--oneline", "--no-merges", "-50", "--format=%s"); err == nil {
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

	if out, err := runCommand(absRepoDir, "gh", "repo", "view", "--json", "stargazerCount", "-q", ".stargazerCount"); err == nil {
		if n, err := strconv.Atoi(strings.TrimSpace(string(out))); err == nil {
			info.stars = n
		}
	}

	if out, err := runCommand(absRepoDir, "gh", "repo", "view", "--json", "licenseInfo", "-q", ".licenseInfo.name"); err == nil {
		l := strings.TrimSpace(string(out))
		if l != "" {
			info.license = l
		}
	}

	if out, err := runCommand(absRepoDir, "gh", "repo", "view", "--json", "primaryLanguage", "-q", ".primaryLanguage.name"); err == nil {
		l := strings.TrimSpace(string(out))
		if l != "" {
			info.language = l
		}
	}

	return info, nil
}

func runCommand(dir, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	return cmd.Output()
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

func centerText(s string, width int) string {
	runeLen := len([]rune(s))
	if runeLen >= width {
		return s
	}
	pad := (width - runeLen) / 2
	return strings.Repeat(" ", pad) + s
}

func formatCommitCount(n int) string {
	return fmt.Sprintf("%d", n)
}
