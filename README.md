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
  list        Show file list
  open        Open file
  path        Show path
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

### Sample

```toml
# $HOME/.config/jnal/config.toml

base_directory = "/home/longkey1/Dropbox/Documents/Journal"
date_format = "2006/01/02"
file_name_format = "2006-01-02.md"
file_template = "# 2006-01-02\n"
open_command = "vim {{ .File }}"
list_command = "ranger {{ .BaseDir }}"
```

`file_name` or `file_template` are using [golang's time format](https://golang.org/src/time/format.go).
