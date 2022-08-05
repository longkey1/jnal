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

var pathCheck bool
var pathDay string
var pathType string

const (
	FileType = "file"
	BaseType = "base"
)

// pathCmd represents the path command
var pathCmd = &cobra.Command{
	Use:   "path",
	Short: "Show path",
	Run: func(cmd *cobra.Command, args []string) {
		j := jnal.NewJnal(config)

		targetDay, err := time.Parse("2006-01-02", pathDay)
		if err != nil {
			log.Fatalf("target day format error %s, %v", pathDay, err)
		}

		targetPath := j.GetDayFilePath(targetDay)
		if pathType == BaseType {
			targetPath = j.GetBaseDirPath()
		}

		if pathCheck {
			if _, err := os.Stat(targetPath); os.IsNotExist(err) {
				log.Fatalf("Not found %s, %v", targetPath, err)
			}
		}

		fmt.Println(targetPath)
	},
}

func init() {
	rootCmd.AddCommand(pathCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// openCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// openCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	pathCmd.Flags().BoolVarP(&pathCheck, "check", "c", false, "directory or file exist check")
	pathCmd.Flags().StringVarP(&pathDay, "day", "d", time.Now().Format("2006-01-02"), "target day (ISO 8601)")
	pathCmd.Flags().StringVarP(&pathType, "type", "t", FileType, fmt.Sprintf("type fmt.F[%s, %s]", FileType, BaseType))
}
