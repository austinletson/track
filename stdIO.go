package main

import (
	"fmt"
	"log"
	"os"
	"time"

	flag "github.com/spf13/pflag"
	cli "github.com/urfave/cli/v2"
)

type Command struct {
	flags    []string
	function func([]*flag.Flag)
}

// Task commands
// TODO make sure that a task is provided
var startFlag *string = flag.StringP("start", "s", "", "Start a task at the current time")
var endFlag *string = flag.StringP("end", "e", "", "End a task")
var priorityFlag *int = flag.IntP("priority", "p", 0, "Assign a priority to an existing task")

// TODO impliment
var timeFlag *string = flag.StringP("time", "t", "", "Change the time of a given command")

// Display commands
var listFlag *bool = flag.BoolP("list", "l", false, "help message for flagname")
var allFlag *bool = flag.BoolP("all", "a", false, "Combined with the list flag it lists all")
var verboseFlag *bool = flag.BoolP("verbose", "v", false, "Combined with list flag it displays more detailed output")
var chartFlag *bool = flag.BoolP("chart", "c", false, "Display a Gantt chart of your tasks")

func readArgs() {

	flag.Parse()
	if noneAreSetBut([]string{"start"}) {
		err := clockIn(makeTask(*startFlag, 0), time.Now())
		if err != nil {
			fmt.Println(err)
		}
	} else if noneAreSetBut([]string{"start", "priority"}) {
		err := clockIn(makeTask(*startFlag, *priorityFlag), time.Now())
		if err != nil {
			fmt.Println(err)
		}
	} else if noneAreSetBut([]string{"end"}) {
		err := clockOut(*endFlag, time.Now())
		if err != nil {
			fmt.Println(err)
		}
	} else if noneAreSetBut([]string{"list"}) {
		taskRecord := ReadTasksFromTasksFile()
		taskList := listTasks(taskRecord, justActiveTasks, notVerbose)
		fmt.Print(taskList)
	} else if noneAreSetBut([]string{"list", "all"}) {
		taskRecord := ReadTasksFromTasksFile()
		taskList := listTasks(taskRecord, allTasks, notVerbose)
		fmt.Print(taskList)
	} else if noneAreSetBut([]string{"list", "verbose"}) {
		taskRecord := ReadTasksFromTasksFile()
		taskList := listTasks(taskRecord, justActiveTasks, verbose)
		fmt.Print(taskList)
	} else if noneAreSetBut([]string{"list", "all", "verbose"}) {
		taskRecord := ReadTasksFromTasksFile()
		taskList := listTasks(taskRecord, allTasks, verbose)
		fmt.Print(taskList)
	} else if noneAreSetBut([]string{"chart"}) {
		taskRecord := ReadTasksFromTasksFile()
		fmt.Print(generateGanttChart(taskRecord))
	} else {
		flag.Usage()
	}

}

func noneAreSetBut(setFlags []string) bool {
	noneAreSet := true
	flagsExist := false
	flag.Visit(func(f *flag.Flag) {
		flagsExist = true
		if !containsString(f.Name, setFlags) {
			noneAreSet = false
		}
	})
	return noneAreSet && flagsExist
}
func getErrorText() (errorText string) {
	return fmt.Sprintf(`track: invalid options or input\n 
track: try track --help for more information`)
}

func getHelpText() (helpText string) {
	helpText = "Track is a fast tool to track anything"
	return helpText
}

func readArgsTest() {
	app := &cli.App{
		Name:  "boom",
		Usage: "make an explosive entrance",
		Action: func(c *cli.Context) error {
			fmt.Printf("Hello %q", c.Args().Get(0))
			return nil
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
