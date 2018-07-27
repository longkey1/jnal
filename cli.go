package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/blang/semver"
	"github.com/BurntSushi/toml"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/mitchellh/go-homedir"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
)

const (
	// Version
	Version string = "0.2.2"
	// ExitCodeOK ...
	ExitCodeOK int = 0
	// ExitCodeError ..
	ExitCodeError int = 1
	// DefaultConfigFileName...
	DefaultConfigFileName string = "config.toml"
)

// Config ...
type Config struct {
	BaseDirectory string `toml:"base_directory"`
	FileName      string `toml:"file_name"`
	OpenCommand   string `toml:"open_command"`
	SearchCommand string `toml:"search_command"`
}

// CLI ...
type CLI struct {
	inStream io.Reader
	outStream io.Writer
	errStream io.Writer
}

// Run ...
func (c *CLI) Run(args []string) int {
	var configPath string

	app := cli.NewApp()
	app.Name = "diary"
	app.Version = Version
	app.Usage = "text file generator for diary"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "config, c",
			Usage: "Load configration from `FILE`",
			Destination: &configPath,
			Value: defaultConfigPath(),
		},
	}
	app.Commands = []cli.Command{
		{
			Name: "open",
			Aliases: []string{"o"},
			Usage: "open file",
			Action: func(ctx *cli.Context) error {
				cnf, err := loadConfig(configPath)
				if err != nil {
					return err
				}

				todayFile := buildTodayFile(cnf.BaseDirectory, cnf.FileName)
				if fileExists(todayFile) == false {
					file, err := os.OpenFile(todayFile, os.O_WRONLY|os.O_CREATE, 0644)
					if err != nil {
						return err
					}
					fmt.Fprintln(file, time.Now().Format("# 2006/01/02"))
					file.Close()
				}

				cmd, err := c.buildCommand(cnf.OpenCommand, cnf.BaseDirectory, cnf.FileName, ""); if err != nil {
					return err
				}

				err = cmd.Run(); if err != nil {
					return err
				}

				return nil
			},
		}, {
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search file",
			Action: func(ctx *cli.Context) error {
				cnf, err := loadConfig(configPath)
				if err != nil {
					return err
				}

				cmd, err := c.buildCommand(cnf.SearchCommand, cnf.BaseDirectory, cnf.FileName, ctx.Args().First()); if err != nil {
					return err
				}

				err = cmd.Run(); if err != nil {
					return err
				}

				return nil
			},
		}, {
			Name:    "self-update",
			Aliases: []string{"su"},
			Usage:   "self update",
			Action: func(ctx *cli.Context) error {
			  v := semver.MustParse(Version)
			  latest, err := selfupdate.UpdateSelf(v, "longkey1/diary"); if err != nil {
			  	return err
				}

			  if latest.Version.Equals(v) {
			  	// latest version is the same as current version. It means current binary is up to date.
			  	log.Println("current binary is the latest version", Version)
			  } else {
			  	log.Println("successfully updated to version", latest.Version)
			  	log.Println("release note:\n", latest.ReleaseNotes)
				}

				return nil
			},
	  },
	}
	app.Writer = c.outStream
	app.ErrWriter = c.errStream

	err := app.Run(args)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func (c *CLI) buildCommand(tpl string, dir string, file string, pattern string) (*exec.Cmd, error) {
	t := template.Must(template.New("").Parse(tpl))
	cmdBuf := new(bytes.Buffer)
	err := t.Execute(cmdBuf, map[string]interface{}{
		"TodayFile": buildTodayFile(dir, file),
		"Pattern": pattern,
		"BaseDirectory": dir,
	})
	if err != nil {
		return nil, err
	}

	sw, err := shellwords.Parse(cmdBuf.String()); if err != nil {
		return nil, err
	}

	cmd := &exec.Cmd{}
	switch len(sw) {
	case 0:
		return nil, fmt.Errorf("not defined command string: %s", cmdBuf.String())
	case 1:
		cmd = exec.Command(sw[0])
	default:
		cmd = exec.Command(sw[0], sw[1:]...)
	}
	cmd.Stdin = c.inStream
	cmd.Stdout = c.outStream
	cmd.Stderr = c.errStream

	return cmd, nil
}

func defaultConfigPath() string {
	home, err := homedir.Dir()
	if err != nil {
		panic(err)
	}
	return fmt.Sprintf("%s/.config/diary/%s", home, DefaultConfigFileName)
}

func loadConfig(path string) (*Config, error) {
	c := &Config{}
	if _, err := toml.DecodeFile(path, c); err != nil {
		return nil, err
	}
	return c, nil
}

func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil
}

func buildTodayFile(dir string, file string) string {
	return time.Now().Format(fmt.Sprintf("%s/%s", dir, file))
}
