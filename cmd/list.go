package cmd

import (
	"fmt"

	"github.com/austinletson/track/core"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		taskRecord := core.ReadTasksFromTasksFile()
		taskList := core.ListTasks(taskRecord, allFlag, verboseFlag)

		fmt.Print(taskList)
	},
}

var allFlag bool
var verboseFlag bool
var tagFlag []string

func init() {
	rootCmd.AddCommand(listCmd)

	listCmd.Flags().BoolVarP(&allFlag, "all", "a", false, "Lists all tasks")
	listCmd.Flags().BoolVarP(&verboseFlag, "verbose", "v", false, "Displays verbose list of tasks")
	listCmd.Flags().StringArrayVarP(&tagFlag, "tag", "t", nil, "Takes a list of the categories to list. If no option is given, all categories are shown.")
}
