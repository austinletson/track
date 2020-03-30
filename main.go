package main

import (
	"fmt"
	"time"
)

var goalTasks []string

func main() {
	startName, list, endName, priority := readArgs()

	fmt.Print("Gotinto main")
	if startName != nil {
		fmt.Print("got into start")
		var priorityValue int
		if priority != nil {
			priorityValue = *priority
		} else {
			priorityValue = 0
		}
		task := Task{
			Name:     *startName,
			Priority: priorityValue,
		}
		err := clockIn(task, time.Now())
		fmt.Println(err)
		if err != nil {
			fmt.Println(err)
		}
	}
	if list != nil {
		taskRecord := ReadTasksFromTasksFile()
		listTasks(taskRecord)
	}
	if endName != nil {
		err := clockOut(*endName, time.Now())
		if err != nil {
			fmt.Println(err)
		}
	}

	//	writeTaskToTasksFile("Write", 2, time.Now(), time.Now())
}

var showList bool
