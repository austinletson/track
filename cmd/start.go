package cmd

import (
	"fmt"
	"time"

	"github.com/austinletson/track/core"
	"github.com/jinzhu/now"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a task",
	Long:  "Starts a task at the current time by default",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Print(cmd.Usage())
			return
		}
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
		taskName := args[0]
		tags := args[1:]
		err := core.ClockIn(core.MakeTask(taskName, priorityFlag, tags), startTime)
		if err != nil {
			fmt.Println(err)
		}
	},
}

var priorityFlag int
var timeFlagStart string

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().IntVarP(&priorityFlag, "priority", "p", 0, "Assign a priority to the given task")
	startCmd.Flags().StringVarP(&timeFlagStart, "time", "t", "", "Change the start time of the given task")
}
