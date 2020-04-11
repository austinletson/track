package cmd

import (
	"fmt"

	"github.com/austinletson/track/common"
	"github.com/austinletson/track/core"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Lists active tasks by default",
	Long: `List is used to quickly view and analyze everything you are tracking. By default it lists the active tasks but there are flags for many options. Below are a few examples of what list can do:

	track list (lists all active tasks)
	track list -a (lists all tasks)
	track list --all --verbose --tags="bugs, features" (lists all the tasks with tags 'bugs' and 'features.' Verbose gives extra information)`,
	Run: func(cmd *cobra.Command, args []string) {

		trimmedTags := common.TrimList(tagsFlagList)
		taskRecord := core.ReadTasksFromTasksFile()

		taskList := core.ListTasks(taskRecord, allFlag, verboseFlag, trimmedTags)
		fmt.Print(taskList)
	},
}

var allFlag bool
var verboseFlag bool
var tagsFlagList []string

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Lists all tasks")
	listCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Displays verbose list of tasks")
	listCmd.Flags().StringSliceVarP(&tagsFlagList, "tags", "g", nil, "Takes a list of tags to list. If no option is given, all tags are shown.")
}
