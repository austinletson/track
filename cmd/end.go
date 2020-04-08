/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

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
	"time"

	"github.com/austinletson/track/core"
	"github.com/jinzhu/now"
	"github.com/spf13/cobra"
)

// endCmd represents the end command
var endCmd = &cobra.Command{
	Use:   "end",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) != 1 {
			fmt.Print(cmd.Usage())
			return
		}
		var endTime time.Time
		if timeFlagEnd != "" {
			endTimeParsed, err := now.Parse(timeFlagEnd)
			endTime = endTimeParsed

			if err != nil {
				fmt.Print(err)
				fmt.Print(endTime)
			}
		} else {
			endTime = time.Now()
		}

		err := core.ClockOut(args[0], endTime)
		if err != nil {
			fmt.Print(err)
		}
	},
}

var timeFlagEnd string

func init() {
	rootCmd.AddCommand(endCmd)

	endCmd.Flags().StringVarP(&timeFlagEnd, "time", "t", "", "Change the time of a given command")
}
