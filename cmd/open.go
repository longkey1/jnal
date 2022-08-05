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
	"github.com/longkey1/jnal/jnal"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var openDay string
var openWithNoCreate bool

var openCmd = &cobra.Command{
	Use:   "open",
	Short: "Open file",
	Run: func(cmd *cobra.Command, args []string) {
		j := jnal.NewJnal(config)

		targetDay, err := time.Parse("2006-01-02", openDay)
		if err != nil {
			log.Fatalf("target day format error %s, %v", pathDay, err)
		}

		dayFile := j.GetDayFilePath(targetDay)
		if openWithNoCreate == false {
			dayFile = j.CreateDayFile(targetDay)
		}
		if _, err := os.Stat(dayFile); os.IsNotExist(err) {
			log.Fatalf("Not found %s file, %v", dayFile, err)
		}

		c, err := j.BuildOpenCommand(targetDay)
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
	openCmd.Flags().StringVarP(&openDay, "day", "d", time.Now().Format("2006-01-02"), "target day (ISO 8601)")
	openCmd.Flags().BoolVarP(&openWithNoCreate, "with-no-create", "", false, "not with creating day file")
}
