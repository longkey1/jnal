# jnal

A simple CLI tool for daily journaling in Markdown.

## Usage

```
A simple CLI tool for daily journaling in Markdown

Usage:
  jnal [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize jnal configuration
  new         Create a journal entry
  path        Show file or directory path
  serve       Start a local preview server
  version     Show version information

Flags:
      --config string   config file (default is $HOME/.config/jnal/config.toml)
  -h, --help            help for jnal

Use "jnal [command] --help" for more information about a command.
```

## Installation

You can download binary from [release page](https://github.com/longkey1/jnal/releases).

## Quick Start

```bash
# Initialize configuration
jnal init

# Edit the config file
vim ~/.config/jnal/config.toml

# Create today's journal entry
jnal new

# Create entry for a specific date
jnal new --date 2024-01-15

# Preview journal entries in browser
jnal serve
```

## Configuration

### Path Format

`path_format` defines the file path structure using [Go's time format](https://golang.org/src/time/format.go):

- `2006-01-02.md` → `2024-01-15.md`
- `2006/2006-01-02.md` → `2024/2024-01-15.md`
- `2006/01/2006-01-02.md` → `2024/01/2024-01-15.md`

### Template Placeholders

**file_template:**
- `{{ .Date }}` - Formatted date (using `date_format`)
- `{{ .Env.<NAME> }}` - Environment variable (e.g., `{{ .Env.HOME }}`)

### CSS Customization

`css` can be a URL (downloaded at startup) or inline CSS:

```toml
[serve]
# URL (classless CSS frameworks work great)
css = "https://cdn.jsdelivr.net/npm/water.css@2/out/water.css"

# Or inline CSS
css = """
body { max-width: 800px; margin: 0 auto; }
"""
```

### Sample Configuration

```toml
# $HOME/.config/jnal/config.toml

base_directory = "/home/user/journal"
date_format = "2006-01-02"
path_format = "2006/2006-01-02.md"
file_template = "# {{ .Date }}\n"

[serve]
port = 8080
sort = "desc"
css = "https://cdn.jsdelivr.net/npm/water.css@2/out/water.css"
```

## Commands

### new

Create a journal entry for the specified date:

```bash
jnal new                       # Today's entry
jnal new --date 2024-01-15     # Specific date
```

### path

Show file or directory path:

```bash
jnal path                      # Today's file path
jnal path --date 2024-01-15    # Specific date's file path
jnal path --type base          # Base directory path
jnal path --check              # Check if path exists
```

### serve

Start a local preview server with hot reload:

```bash
jnal serve                     # Default port 8080
jnal serve --port 3000         # Custom port
jnal serve --sort asc          # Oldest first
```

### build

Generate static HTML files:

```bash
jnal build                     # Output to public/
jnal build --output dist       # Custom output directory
```

### init

Initialize configuration file:

```bash
jnal init           # Create default config
jnal init --force   # Overwrite existing config
```
