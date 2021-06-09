package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/blang/semver"
	"github.com/mitchellh/go-homedir"
	"github.com/rhysd/go-github-selfupdate/selfupdate"
	"github.com/urfave/cli"
)

const (
	Version string = "0.6.0"
	ExitCodeOK int = 0
	ExitCodeError int = 1
	DefaultConfigFileName string = "config.toml"
)

type Config struct {
	BaseDirectory  string `toml:"base_directory"`
	DateFormat     string `toml:"date_format"`
	FileNameFormat string `toml:"file_name_format"`
	FileTemplate   string `toml:"file_template"`
	OpenCommand    string `toml:"open_command"`
	ListCommand    string `toml:"list_command"`
	FindCommand    string `toml:"find_command"`
	SaveCommand    string `toml:"save_command"`
}

type CLI struct {
	inStream  io.Reader
	outStream io.Writer
	errStream io.Writer
}

func (c *CLI) Run(args []string) int {
	var configPath string

	app := cli.NewApp()
	app.Name = "diary"
	app.Version = Version
	app.Usage = "Text file based diary command"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "config, c",
			Usage:       "Load configuration from `FILE`",
			Destination: &configPath,
			Value:       defaultConfigPath(),
		},
	}
	app.Commands = []cli.Command{
		{
			Name:    "open",
			Aliases: []string{"o"},
			Usage:   "open file",
			Action: func(ctx *cli.Context) error {
				cnf, err := loadConfig(configPath)
				if err != nil {
					return err
				}

				before := 0
				if ctx.Bool("yesterday") {
					before = -1
				}
				targetDay := time.Now().AddDate(0, 0, before)
				dayFile := buildTargetDayFileName(cnf.BaseDirectory, cnf.FileNameFormat, targetDay)
				dayDir := filepath.Dir(dayFile)
				dayDate := targetDay.Format(cnf.DateFormat)
				if _, err = os.Stat(dayDir); os.IsNotExist(err) {
					if err = os.MkdirAll(dayDir, 0755); err != nil {
						return err
					}
				}
				if _, err = os.Stat(dayFile); os.IsNotExist(err) {
					file, err := os.OpenFile(dayFile, os.O_WRONLY|os.O_CREATE, 0644)
					if err != nil {
						return err
					}
					_, err = fmt.Fprintln(file, buildTargetDayFileContent(cnf.FileTemplate, dayDate))
					if err != nil {
						return err
					}
					err = file.Close()
					if err != nil {
						return err
					}
				}

				cmd, err := c.buildCommand(cnf.OpenCommand, cnf.BaseDirectory, dayDate, dayFile, "")
				if err != nil {
					return err
				}

				err = cmd.Run()
				if err != nil {
					return err
				}

				return nil
			},
			Flags: []cli.Flag{
				cli.BoolFlag{
					Name:  "yesterday, y",
				},
			},
		}, {
			Name:    "list",
			Aliases: []string{"l"},
			Usage:   "list files",
			Action: func(ctx *cli.Context) error {
				cnf, err := loadConfig(configPath)
				if err != nil {
					return err
				}

				cmd, err := c.buildCommand(cnf.ListCommand, cnf.BaseDirectory, "", "", "")
				if err != nil {
					return err
				}

				err = cmd.Run()
				if err != nil {
					return err
				}

				return nil
			},
		}, {
			Name:    "find",
			Aliases: []string{"f"},
			Usage:   "find files",
			Action: func(ctx *cli.Context) error {
				cnf, err := loadConfig(configPath)
				if err != nil {
					return err
				}

				cmd, err := c.buildCommand(cnf.FindCommand, cnf.BaseDirectory, "", "", ctx.Args().First())
				if err != nil {
					return err
				}

				err = cmd.Run()
				if err != nil {
					return err
				}

				return nil
			},
		}, {
			Name:    "save",
			Aliases: []string{"s"},
			Usage:   "save files",
			Action: func(ctx *cli.Context) error {
				cnf, err := loadConfig(configPath)
				if err != nil {
					return err
				}

				cmd, err := c.buildCommand(cnf.SaveCommand, cnf.BaseDirectory, "", "", "")
				if err != nil {
					return err
				}

				err = cmd.Run()
				if err != nil {
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
				latest, err := selfupdate.UpdateSelf(v, "longkey1/diary")
				if err != nil {
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
		_, _ = fmt.Fprintln(c.errStream, err)
		return ExitCodeError
	}
	return ExitCodeOK
}

func (c *CLI) buildCommand(tpl string, dir string, date string, file string, pattern string) (*exec.Cmd, error) {
	t := template.Must(template.New("").Parse(tpl))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, map[string]interface{}{
		"BaseDir": dir,
		"Date": date,
		"File": file,
		"Pattern": pattern,
	})
	if err != nil {
		return nil, err
	}

	cmd := exec.Command("sh", "-c", buf.String())
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

func buildTargetDayFileName(dir string, file string, day time.Time) string {
	return fmt.Sprintf("%s/%s", dir, day.Format(file))
}

func buildTargetDayFileContent(tpl string, date string) string {
	t := template.Must(template.New("").Parse(tpl))
	buf := new(bytes.Buffer)
	err := t.Execute(buf, map[string]interface{}{
		"Date": date,
	})
	if err != nil {
		panic(err)
	}

	return buf.String()
}
