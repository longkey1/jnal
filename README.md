# jnal

A simple CLI tool for daily journaling in Markdown.

## Usage

```
A simple CLI tool for daily journaling in Markdown

Usage:
  jnal [command]

Available Commands:
  build       Build static HTML files
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
  init        Initialize jnal configuration
  new         Create a journal entry
  path        Show file or directory path
  serve       Start a local preview server
  version     Show version information

Flags:
      --config string   config file (default is $JNAL_CONFIG or $HOME/.config/jnal/config.toml)
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

### Environment Variables

- `JNAL_CONFIG` - Path to config file (overrides default `$HOME/.config/jnal/config.toml`)

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
[build]
# URL (classless CSS frameworks work great)
css = "https://cdn.jsdelivr.net/npm/water.css@2/out/water.css"

# Or inline CSS
css = """
body { max-width: 800px; margin: 0 auto; }
"""
```

The year navigation is sticky by default. You can override this behavior:

```css
nav { position: static; }
```

### Heading Shift

When rendering multiple journal entries on a single page, headings are shifted to maintain proper HTML hierarchy:

- Page title: H1
- Year: H2
- Month: H3
- Date: H4
- Entry content H1 → H5, H2 → H6, etc.

By default, `heading_shift = 4`. Set to `0` to disable and output entry content as-is:

```toml
[build]
heading_shift = 0  # Disable heading shift
```

### Auto-linking URLs

URLs in journal entries are automatically converted to clickable links. By default, all links open in a new tab with `target="_blank"` and `rel="noopener noreferrer"` for security.

```toml
[build]
linkify = true           # Auto-convert URLs to links (default: true)
link_target_blank = true # Open links in new tab (default: true)
```

### Sample Configuration

```toml
# $HOME/.config/jnal/config.toml

[common]
base_directory = "/home/user/journal"
date_format = "2006-01-02"
path_format = "2006/2006-01-02.md"

[new]
file_template = "# {{ .Date }}\n"

[build]
title = "My Journal"
sort = "desc"
css = "https://cdn.jsdelivr.net/npm/water.css@2/out/water.css"

[serve]
port = 8080
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
jnal path --base               # Base directory path
jnal path --check              # Check if path exists
```

### serve

Start a local preview server:

```bash
jnal serve                     # Default port 8080
jnal serve --port 3000         # Custom port
jnal serve --sort asc          # Oldest first
jnal serve --live-reload       # Enable browser auto-reload on file changes
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

## Docker

A Docker image is available for running jnal in containers:

```bash
docker run -v /path/to/journal:/app -v /path/to/config.toml:/app/config.toml ghcr.io/longkey1/jnal build
```

The image uses `/app` as the working directory and expects:
- Journal files mounted at `/app` (or configured `base_directory`)
- Config file at `/app/config.toml` (via `JNAL_CONFIG` environment variable)
