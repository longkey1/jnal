package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	shellwords "github.com/mattn/go-shellwords"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

const (
	// Version
	Version string = "0.1.0"
	// ExitCodeOK ...
	ExitCodeOK int = 0
	// ExitCodeError ..
	ExitCodeError int = 1
	// DefaultConfigFileName...
	DefaultConfigFileName string = "config.toml"
)

// CLI ...
type CLI struct {
	outStream io.Writer
	errStream io.Writer
}

// Config ...
type Config struct {
	BaseDirectory string `toml:"base_directory"`
	FileName      string `toml:"file_name"`
	OpenCommand   string `toml:"open_command"`
	SearchCommand string `toml:"search_command"`
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
			Name:        "config, c",
			Usage:       "Load configration from `FILE`",
			Destination: &configPath,
			Value:       defaultConfigPath(),
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "open",
			Aliases: []string{"o"},
			Usage:   "open file",
			Action: func(c *cli.Context) error {
				cnf, err := loadConfig(configPath)
				if err != nil {
					return err
				}
				filepath := time.Now().Format(fmt.Sprintf("%s/%s", cnf.BaseDirectory, cnf.FileName))
				if fileExists(filepath) == false {
					file, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, 0644)
					if err != nil {
						return err
					}
					fmt.Fprintln(file, time.Now().Format("# 2006/01/02"))
					file.Close()
				}

				cmdFmt := template.Must(template.New("cmd").Parse(cnf.OpenCommand))
				openCmd := new(bytes.Buffer)
				err = cmdFmt.Execute(openCmd, map[string]interface{}{
					"TodayFile": filepath,
				})
				if err != nil {
					return err
				}
				cw, err := shellwords.Parse(openCmd.String())
				if err != nil {
					return err
				}

				cmd := &exec.Cmd{}
				switch len(cw) {
				case 0:
					return errors.New("Not defined open command")
				case 1:
					cmd = exec.Command(cw[0])
				default:
					cmd = exec.Command(cw[0], cw[1:]...)
				}
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					return err
				}
				return nil
			},
		}, {
			Name:    "search",
			Aliases: []string{"s"},
			Usage:   "search file",
			Action: func(c *cli.Context) error {
				cnf, err := loadConfig(configPath)
				if err != nil {
					return err
				}

				cmdFmt := template.Must(template.New("cmd").Parse(cnf.SearchCommand))
				searchCmd := new(bytes.Buffer)
				err = cmdFmt.Execute(searchCmd, map[string]interface{}{
					"Pattern":       c.Args().First(),
					"BaseDirectory": cnf.BaseDirectory,
				})
				if err != nil {
					return err
				}
				cw, err := shellwords.Parse(searchCmd.String())
				if err != nil {
					return err
				}

				cmd := &exec.Cmd{}
				switch len(cw) {
				case 0:
					return errors.New("Not defined open command")
				case 1:
					cmd = exec.Command(cw[0])
				default:
					cmd = exec.Command(cw[0], cw[1:]...)
				}
				cmd.Stdin = os.Stdin
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err = cmd.Run()
				if err != nil {
					return err
				}
				return nil
			},
		},
	}

	err := app.Run(args)
	if err != nil {
		fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}
	return ExitCodeOK
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
