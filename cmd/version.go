package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version string = "0.1.0"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of Track",
	Long:  `Prints the version of  Track`,

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Track: v" + version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
