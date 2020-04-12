package main

import (
	"github.com/austinletson/track/cmd"
	"github.com/austinletson/track/core"
)

func main() {
	pipedInput := getStdin()
	if pipedInput != "" {
		core.TakeNote(pipedInput, core.GetActiveTasks(core.ReadTasksFromTasksFile()))
	}
	cmd.Execute()

}
