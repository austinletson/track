package main

import (
	"fmt"
	"time"

	flag "github.com/spf13/pflag"
)

type Command struct {
	flags    []string
	function func([]*flag.Flag)
}

var listFlag *bool = flag.BoolP("list", "l", false, "help message for flagname")

// TODO make sure that a task is provided
var startFlag *string = flag.StringP("start", "s", "asleep", "Start a task at the current time")
var endFlag *string = flag.StringP("end", "e", "asleep", "End a task")
var priorityFlag *int = flag.IntP("priority", "p", 0, "Assign a priority to an existing task")
var allFlag *bool = flag.BoolP("all", "a", false, "Combined with the list flag it lists all")
var chartFlag *bool = flag.BoolP("chart", "c", false, "Display a Gantt chart of your tasks")
var timeFlag *string = flag.StringP("time", "t", "", "Change the time of a given command")

func readArgs() {

	flag.Parse()
	if noneAreSetBut([]string{"start"}) {
		err := clockIn(makeTask(*startFlag, 0), time.Now())
		fmt.Println(err)
	} else if noneAreSetBut([]string{"start", "priority"}) {
		err := clockIn(makeTask(*startFlag, *priorityFlag), time.Now())
		fmt.Println(err)
	} else if noneAreSetBut([]string{"end"}) {
		err := clockOut(*endFlag, time.Now())
		fmt.Println(err)
	} else if noneAreSetBut([]string{"list"}) {
		taskRecord := ReadTasksFromTasksFile()
		listTasks(taskRecord)
	} else if noneAreSetBut([]string{"chart"}) {
		taskRecord := ReadTasksFromTasksFile()
		fmt.Print(generateGanttChart(taskRecord))
	} else {

		fmt.Println(getErrorText())
	}

}

func noneAreSetBut(setFlags []string) bool {
	noneAreSet := true
	flag.Visit(func(f *flag.Flag) {
		fmt.Print(f.Name)
		if !containsString(f.Name, setFlags) {
			noneAreSet = false
		}
	})
	return noneAreSet
}
func getErrorText() (errorText string) {
	return fmt.Sprintf(`track: invalid options or input\n 
track: try track --help for more information`)
}

func getHelpText() (helpText string) {
	helpText = "Track is a fast tool to track anything"
	return helpText
}

func listTasks(taskRecords TaskRecords) {
	for _, task := range getActiveTasks(taskRecords) {
		fmt.Println(task.Name)
	}
}
