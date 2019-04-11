# diary

Text file based diary command.

## USAGE
```
NAME:
   diary - Text file based diary command

USAGE:
   diary [global options] command [command options] [arguments...]

VERSION:
   0.5.x

COMMANDS:
     open, o          open file
     list, l          list file
     search, s        search file
     self-update, su  self update
     help, h          Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load configration from FILE (default: "/home/moto/.config/diary/config.toml")
   --help, -h              show help
   --version, -v           print the version
```

## Installation

You can download binary from [release page](https://github.com/longkey1/diary-bin/releases).

## Configuration

### Placeholder

- `{{ .BaseDirectory }}`
- `{{ .DayFile }}`
- `{{ .PATTERN }}`

### Sample

```toml
# $HOME/.config/diary/config.toml

base_directory = "/home/longkey1/Dropbox/Documents/Diary"
file_name = "2006-01-02.md"
file_template = "# 2006/01/02\n"
open_command = "vim {{ .DayFile }}"
list_command = "selected=$(pt -g .md \"{{ .BaseDirectory }}\" | fzf --query \"$LBUFFER\"); [[ -n ${selected} ]] && env LESS=\"-R -X\" less ${selected} || true"
search_command = "selected=$(pt \"{{ .Pattern }}\" \"{{ .BaseDirectory }}\" | fzf --query \"$LBUFFER\" | awk -F : '{print \"-c \" $2 \" \" $1}'); [[ -n ${selected} ]] && vim $selected || true"

```

`fila_name` is using [golang's time format](https://golang.org/src/time/format.go).
