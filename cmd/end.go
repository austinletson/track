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
	Short: "Ends an active task(s)",
	Long: `Start is used to end an active task(s). By default the task(s) end at the current time, however the user can provide a time with the --time flag. Here are a few examples of what end can do:

	task end bugs (ends active task bugs)
	task end mobile_research -t 15:04 (ends active task mobile_research at 15:04)
	task end bugs mobile_research (ends tasks bugs and mobile_research at the current time)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
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
		taskRecords := core.ReadTasksFromTasksFile()
		for _, taskName := range args {
			var err error
			taskRecords, err = core.ClockOut(taskRecords, taskName, endTime)
			if err != nil {
				fmt.Println(err)
			}
		}
		core.WriteTasksToTaskFile(taskRecords)
	},
}

var timeFlagEnd string

func init() {
	rootCmd.AddCommand(endCmd)

	endCmd.Flags().StringVarP(&timeFlagEnd, "time", "t", "", "Change the time of a given command. Format 15:04")
}
