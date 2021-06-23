# diary

Text file based diary command.

## Usage

```
Text file based diary command

Usage:
  diary [command]

Available Commands:
  help         Help about any command
  list         Show file list
  open         Open file
  self-update  Self update

Flags:
      --config string   config file (default is $HOME/.config/diary/config.toml)
  -h, --help            help for diary
  -t, --toggle          Help message for toggle
  -v, --version         version for diary
```

## Installation

You can download binary from [release page](https://github.com/longkey1/diary/releases).

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
# $HOME/.config/diary/config.toml

base_directory = "/home/longkey1/Dropbox/Documents/Diary"
date_format = "2006/01/02"
file_name_format = "2006-01-02.md"
file_template = "# 2006-01-02\n"
open_command = "vim {{ .File }}"
list_command = "ranger {{ .BaseDir }}"
find_command = "selected=$(pt \"{{ .Pattern }}\" \"{{ .BaseDir }}\" | fzf --query \"$LBUFFER\" | awk -F : '{print \"-c \" $2 \" \" $1}'); [[ -n ${selected} ]] && echo $selected || true"
save_command = "git commit -m \"Auto commit by diary command\""
```

`file_name` or `file_template` are using [golang's time format](https://golang.org/src/time/format.go).
