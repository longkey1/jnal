/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"bytes"
	"fmt"
	"github.com/longkey1/diary/util"
	"github.com/spf13/cobra"
	"log"
	"os"
	"path/filepath"
	"text/template"
	"time"
)

var Yesterday bool

// openCmd represents the open command
var openCmd = &cobra.Command{
	Use:   "open",
	Short: "open file",
	Run: func(cmd *cobra.Command, args []string) {
		before := 0
		if Yesterday {
			before = -1
		}
		targetDay := time.Now().AddDate(0, 0, before)
		dayFile := buildTargetDayFileName(config.BaseDirectory, config.FileNameFormat, targetDay)
		dayDir := filepath.Dir(dayFile)
		dayDate := targetDay.Format(config.DateFormat)
		if _, err := os.Stat(dayDir); os.IsNotExist(err) {
			if err = os.MkdirAll(dayDir, 0755); err != nil {
				log.Fatalf("Unable to make directory, %v", err)
			}
		}
		if _, err := os.Stat(dayFile); os.IsNotExist(err) {
			file, err := os.OpenFile(dayFile, os.O_WRONLY|os.O_CREATE, 0644)
			if err != nil {
				log.Fatalf("Unable to open file, %v", err)
			}
			_, err = fmt.Fprintln(file, buildTargetDayFileContent(config.FileTemplate, dayDate))
			if err != nil {
				log.Fatalf("Unable to build file content, %v", err)
			}
			err = file.Close()
			if err != nil {
				log.Fatalf("Unable to close file, %v", err)
			}
		}

		c, err := util.BuildCommand(config.OpenCommand, config.BaseDirectory, dayDate, dayFile, "")
		if err != nil {
			log.Fatalf("Unable to build open command, %v", err)
		}

		err = c.Run()
		if err != nil {
			log.Fatalf("Unable to execute open command, %#v", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(openCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	openCmd.Flags().BoolVarP(&Yesterday,"yesterday", "y", false, "yesterday")
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
