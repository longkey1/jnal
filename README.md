# jnal

Text file based journal command.

## Usage

```
Text file based journal command

Usage:
  jnal [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  file        Show file path
  help        Help about any command
  list        Show file list
  open        Open file
  self-update self update binary file

Flags:
      --config string   config file (default is $HOME/.config/jnal/config.toml)
  -h, --help            help for jnal
  -t, --toggle          Help message for toggle
  -v, --version         version for jnal

Use "jnal [command] --help" for more information about a command.
```

## Installation

You can download binary from [release page](https://github.com/longkey1/jnal/releases).

## Configuration

### Placeholder

**file template**

- `{{ .Date }}`

**command**
- `{{ .BaseDir }}`
- `{{ .File }}`
- `{{ .Pattern }}`

### Sample

```toml
# $HOME/.config/jnal/config.toml

base_directory = "/home/longkey1/Dropbox/Documents/Journal"
date_format = "2006/01/02"
file_name_format = "2006-01-02.md"
file_template = "# 2006-01-02\n"
open_command = "vim {{ .File }}"
list_command = "ranger {{ .BaseDir }}"
find_command = "selected=$(pt \"{{ .Pattern }}\" \"{{ .BaseDir }}\" | fzf --query \"$LBUFFER\" | awk -F : '{print \"-c \" $2 \" \" $1}'); [[ -n ${selected} ]] && echo $selected || true"
save_command = "git commit -m \"Auto commit by diary command\""
```

`file_name` or `file_template` are using [golang's time format](https://golang.org/src/time/format.go).
