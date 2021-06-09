# diary

Text file based diary command.

## Usage

```
diary [global options] command [command options] [arguments...]

COMMANDS:
   open, o          open file
   list, l          list files
   find, f          find files
   save, s          save files
   self-update, su  self update
   help, h          Shows a list of commands or help for one comma

GLOBAL OPTIONS:
   --config FILE, -c FILE  Load configration from FILE (default: "/home/longkey1/.config/diary/config.toml")
   --help, -h              show help
   --version, -v           print the version
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
