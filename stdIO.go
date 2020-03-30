package main

import (
	flag "github.com/spf13/pflag"
)

var configDirPath = "~/.config/track/"
var tasksFile = "tasks.tt"

var listFlag *bool = flag.BoolP("list", "l", false, "help message for flagname")
var startFlag *string = flag.StringP("start", "s", "asleep", "Start a task at the current	 time")
var endFlag *string = flag.StringP("end", "e", "asleep", "End a task")
var priorityFlag *int = flag.IntP("priority", "p", 0, "Assign a priority to an existing task")

func readArgs() (startFlagValue *string, listFlagValue *bool, endFlagValue *string, priorityFlagValue *int) {
	flag.Parse()
	startFlagValue = nil
	if isFlagPassed("start") {
		startFlagValue = startFlag
	}
	listFlagValue = nil
	if isFlagPassed("list") {
		listFlagValue = listFlag
	}
	endFlagValue = nil
	if isFlagPassed("end") {
		endFlagValue = endFlag
	}
	priorityFlag = nil
	if isFlagPassed("priority") {
		priorityFlagValue = priorityFlag
	}
	return startFlagValue, listFlagValue, endFlagValue, priorityFlagValue
}

func isFlagPassed(name string) bool {
	found := false
	flag.Visit(func(f *flag.Flag) {
		if f.Name == name {
			found = true
		}
	})
	return found
}
