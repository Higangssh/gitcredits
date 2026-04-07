# gitcredits

**`Git log` doesn't do them justice. Turn your contributors into movie stars.**

<p align="center">
  <img src="assets/demo.gif" alt="gitcredits demo" width="720">
</p>

## Install

### Quick install (recommended)

```bash
curl -sL https://raw.githubusercontent.com/Higangssh/gitcredits/main/install.sh | sh
```

Supports macOS (Intel/Apple Silicon), Linux (x86_64/ARM64), and Windows (x86_64 via Git Bash).

### Pre-built binaries

Download the latest release for your platform from [GitHub Releases](https://github.com/Higangssh/gitcredits/releases).

### Go

```bash
go install github.com/Higangssh/gitcredits@latest
```

### From source

```bash
git clone https://github.com/Higangssh/gitcredits.git
cd gitcredits
go build -o gitcredits .
```

## Usage

```bash
cd your-repo
gitcredits

# or target another repository directly
gitcredits /path/to/your-repo
```

That's it. Run `gitcredits` inside a Git repository, or pass a target repository path directly.

### Themes

**Matrix** — digital rain with text resolve effect:

```bash
gitcredits --theme matrix
```

<p align="center">
  <img src="assets/matrix-demo.gif" alt="gitcredits matrix theme" width="720">
</p>

**Spider-Man** — glitch effects with radial web transitions:

```bash
gitcredits --theme spiderman
```

<p align="center">
  <img src="assets/spiderman-demo.gif" alt="gitcredits spider-man theme" width="720">
</p>

### Export to GIF

Save the credits as a high-quality GIF — perfect for READMEs, presentations, or sharing.

```bash
gitcredits --output credits.gif
gitcredits --output credits.gif --theme matrix
gitcredits --output credits.gif --theme spiderman
gitcredits /path/to/your-repo --output credits.gif
```

Requires [VHS](https://github.com/charmbracelet/vhs) and [ffmpeg](https://ffmpeg.org/):

```bash
brew install vhs ffmpeg
```

VHS records the terminal in real-time, and ffmpeg converts it to an optimized GIF with 2-pass palette generation for maximum quality.

### Controls

| Key | Action |
|-----|--------|
| `↑` / `↓` | Manual scroll |
| `q` / `Esc` | Quit |

## What it shows

- **ASCII art title** from your repo name
- **Project lead** — top contributor by commits
- **Contributors** — everyone who committed
- **Notable scenes** — recent `feat:` and `fix:` commits
- **Stats** — total commits, contributors, GitHub stars, language, license

GitHub metadata (stars, description, license) requires [`gh` CLI](https://cli.github.com/) to be installed and authenticated. Works without it — you'll just get git-only data.

## Requirements

- **git** (required) — commit history, contributors, repo info
- **Go 1.21+** — for `go install`
- [`gh` CLI](https://cli.github.com/) (optional) — enables GitHub stars, license, language, and description
- [VHS](https://github.com/charmbracelet/vhs) + [ffmpeg](https://ffmpeg.org/) (optional) — required for `--output` GIF export

## License

MIT
