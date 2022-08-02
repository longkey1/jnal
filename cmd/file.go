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
	"fmt"
	"github.com/longkey1/jnal/jnal"
	"github.com/spf13/cobra"
	"log"
	"os"
	"time"
)

var fileYesterday bool
var fileCreate bool

// openCmd represents the open command
var fileCmd = &cobra.Command{
	Use:   "file",
	Short: "Show file path",
	Run: func(cmd *cobra.Command, args []string) {
		j := jnal.NewJnal(config)

		before := 0
		if fileYesterday {
			before = -1
		}
		targetDay := time.Now().AddDate(0, 0, before)

		dayFile := j.GetFileName(targetDay)
		if fileCreate {
			j.CreateFile(targetDay)
		}
		if _, err := os.Stat(dayFile); os.IsNotExist(err) {
			log.Fatalf("Not found %s file, %v", dayFile, err)
		}

		fmt.Println(dayFile)
	},
}

func init() {
	rootCmd.AddCommand(fileCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	fileCmd.Flags().BoolVarP(&fileYesterday, "yesterday", "y", false, "yesterday")
	fileCmd.Flags().BoolVarP(&fileCreate, "create", "c", false, "create day file")
}
