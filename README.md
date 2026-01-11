# jnal

Text file based journal command.

## Usage

```
Text file based journal command

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
jnal new -d 2024-01-15

# Preview journal entries in browser
jnal serve
```

## Configuration

### File Naming

Journal files are named in `yyyy-mm-dd.md` format:
- `2024-01-15.md`

### Template Placeholders

**file_template:**
- `{{ .Date }}` - Formatted date
- `{{ .Env.<NAME> }}` - Environment variable (e.g., `{{ .Env.HOME }}`)

### Sample Configuration

```toml
# $HOME/.config/jnal/config.toml

base_directory = "/home/user/journal"
date_format = "2006-01-02"
file_template = "# {{ .Date }}\n"

[serve]
port = 8080
sort = "desc"      # desc, asc
```

`file_template` uses [Go's time format](https://golang.org/src/time/format.go).

## Commands

### new

Create a journal entry for the specified date:

```bash
jnal new              # Today's entry
jnal new -d 2024-01-15  # Specific date
```

### path

Show file or directory path:

```bash
jnal path                    # Today's file path
jnal path -d 2024-01-15      # Specific date's file path
jnal path -t base            # Base directory path
jnal path -c                 # Check if path exists
```

### serve

Start a local preview server with hot reload:

```bash
jnal serve                   # Default port 8080
jnal serve -p 3000           # Custom port
jnal serve -s asc            # Oldest first
```

### init

Initialize configuration file:

```bash
jnal init           # Create default config
jnal init --force   # Overwrite existing config
```
