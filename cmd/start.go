package cmd

import (
	"fmt"
	"time"

	"github.com/austinletson/track/common"
	"github.com/austinletson/track/core"
	"github.com/jinzhu/now"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Starts a task(s)",
	Long: `Start is used to start a new or existing inactive task(s). By default the task(s) starts at the current time, however the user can provide a time with the --time flag. Here are a few examples of what start can do:

	task start bugs (starts a task bugs at the current time)
	task start mobile_research -t 15:04 (starts a task mobile_research at 15:04)
	task start bugs mobile_research (starts tasks bugs and mobile_research at the current time)
	task start retro -g "team, meeting" (starts a task with tags team and meeting)`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Usage()
			return
		}

		trimmedFlags := common.TrimList(tagsFlagStart)

		var startTime time.Time
		if timeFlagStart != "" {
			startTimeParsed, err := now.Parse(timeFlagStart)
			startTime = startTimeParsed
			if err != nil {
				fmt.Print(err)
			}
		} else {
			startTime = time.Now()
		}

		taskRecords := core.ReadTasksFromTasksFile()
		for _, taskName := range args {
			taskRecords = core.UpdateTask(taskRecords, taskName, priorityFlag, trimmedFlags)

			var startErr error
			taskRecords, startErr = core.ClockIn(taskRecords, taskName, startTime)
			if startErr != nil {
				fmt.Println(startErr)
			}
		}
		core.WriteTasksToTaskFile(taskRecords)
	},
}

var priorityFlag int
var timeFlagStart string
var tagsFlagStart []string

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVarP(&priorityFlag, "priority", "p", 0, "Assign a priority to the given task")
	startCmd.Flags().StringVarP(&timeFlagStart, "time", "t", "", "Change the start time of the given task. Format 15:04.")
	startCmd.Flags().StringSliceVarP(&tagsFlagStart, "tags", "g", nil, "Takes a list of the tags to attach to the given task")
}
