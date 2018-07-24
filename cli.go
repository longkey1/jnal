package main

import (
	"fmt"
	"io"

	"github.com/BurntSushi/toml"
	"github.com/mitchellh/go-homedir"
	"github.com/urfave/cli"
)

const (
	// Version
	Version string = "0.0.1"
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
	BaseDir string `toml:"base_dir"`
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
				fmt.Println("base directory: ", cnf.BaseDir)
				if err != nil {
					return err
				}
				fmt.Println("open file: ", c.Args().First())
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
