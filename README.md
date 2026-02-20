<h1 align="center">typtea</h1>

<p align="center">Minimal terminal-based typing speed tester with support for dozens of programming languages</p>

<div align="center">
<img src="assets/example.gif" />
</div>

<br>

<p align="center">
  <i>Featured on</i>
</p>
<p align="center">
  <a href="https://terminaltrove.com/typtea/">
    <img width="25%" alt="Terminal Trove logo" src="https://github.com/user-attachments/assets/f85b1470-a574-4bc2-b002-4e6dddb9e277" />
  </a>
</p>

---

## Features

- **Terminal-based typing** with WPM and accuracy tracking
- **Multi-language support** including English and 30+ programming languages
- **Infinite word generation** with smooth 3-line scrolling display
- **Minimalist TUI** built with Bubble Tea and Lipgloss
- **Embedded language data** for easy distribution
- **Accurate metrics** following standard typing test calculations

### Supported Languages

| | | | |
|-----------|-----------|-----------|-----------|
| Bash      | C         | C++       | C#        |
| Crystal   | CSS       | Emacs     | English 1k|
| Erlang    | Go        | Haskell   | HTML      |
| Java      | JavaScript| JSON      | Julia     |
| Lisp      | Lua       | OCaml     | Perl      |
| PHP       | PowerShell| Python    | R         |
| Ruby      | Rust      | SCSS      | SQL       |
| Swift     | TeX       | TypeScript| Vala      |
| Vimscript | Wolfram   | YAML      | Zig       |
| | | | |

---

## Installation

[![Packaging status](https://repology.org/badge/vertical-allrepos/typtea.svg)](https://repology.org/project/typtea/versions)

### AUR

```bash
yay -S typtea
```

### FreeBSD

```bash
pkg install typtea
```

### Homebrew

```bash
brew install typtea
```

### Nix

```bash
nix profile install nixpkgs#typtea
```

### Void Linux

```bash
sudo xbps-install -S typtea
```

### Via `go install`

```bash
go install github.com/ashish0kumar/typtea@latest
```

### Build from Source

```bash
git clone --depth=1 https://github.com/ashish0kumar/typtea
cd typtea/
go build
sudo mv typtea /usr/local/bin/
typtea -h
```

---

## Usage

### Basic Commands

```yaml
# Start a 30-second English typing test (default)
typtea start

# Start a 60-second typing test
typtea start --duration 60

# Start a Rust keywords typing test
typtea start --lang rust

# Combine duration and language
typtea start --duration 45 --lang javascript

# List all available languages
typtea start --list-langs

# Get help
typtea --help
typtea start --help
```

### During the Test

- **The test starts** when you begin typing
- **Backspace** to correct mistakes
- **Enter** to restart after completion
- **Esc** to quit the application

---

## Development

### Prerequisites

[Go 1.19+](https://go.dev/doc/install)

### Setup

```bash
git clone https://github.com/ashish0kumar/typtea.git
cd typtea
go mod tidy
go build
./typtea start
```

### Adding New Languages

1. Create a `JSON` file in `internal/game/data/` with the format:

```json
{
  "name": "Language Name",
  "words": ["word1", "word2", "word3", ...]
}
```

2. Rebuild the application to embed the new language data

---

## Community Extensions

Looking for persistent stats or a database-backed version of typtea?

Check out [PR #1: Database for storing stats](https://github.com/ashish0kumar/typtea/pull/1), which adds SQLite-based session tracking and best-score display.  
This feature is not in the main branch to keep typtea minimal, but you can use it if you want to track your typing history.

## Dependencies

- [**Bubble Tea**](https://github.com/charmbracelet/bubbletea) - TUI framework
- [**Lipgloss**](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [**Cobra**](https://github.com/spf13/cobra) - CLI framework

## Contributing

Contributions are always welcome! If you have ideas, bug reports, or want to submit code, please feel free to open an issue or a pull request.

## Contributors

<a href="https://github.com/ashish0kumar/typtea/graphs/contributors">
  <img src="https://contrib.rocks/image?repo=ashish0kumar/typtea" />
</a>

<br>

<p align="center">
	<img src="https://raw.githubusercontent.com/catppuccin/catppuccin/main/assets/footers/gray0_ctp_on_line.svg?sanitize=true" />
</p>

<p align="center">
        <i><code>&copy 2025-present <a href="https://github.com/ashish0kumar">Ashish Kumar</a></code></i>
</p>

<div align="center">
<a href="https://github.com/ashish0kumar/typtea/blob/main/LICENSE"><img src="https://img.shields.io/github/license/ashish0kumar/typtea?style=for-the-badge&color=CBA6F7&logoColor=cdd6f4&labelColor=302D41" alt="LICENSE"></a>&nbsp;&nbsp;
</div>
