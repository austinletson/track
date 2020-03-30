package main

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	storageFile = "tasks.tt"
)

func WriteTasksToTaskFile(taskRecords TaskRecords) {

	err := os.Remove(storageFile)
	isError(err)

	file, err := os.Create(storageFile)
	isError(err)
	defer file.Close()

	formattedTask, err := yaml.Marshal(taskRecords)
	isError(err)

	err = ioutil.WriteFile(storageFile, formattedTask, 0644)
	isError(err)
}

// Read tasks edn from file and return them
func ReadTasksFromTasksFile() TaskRecords {
	file, err := os.Open(storageFile)
	isError(err)
	defer file.Close()

	var bytes []byte
	bytes, err = ioutil.ReadAll(file)
	isError(err)

	var taskRecords TaskRecords
	err = yaml.Unmarshal(bytes, &taskRecords)
	isError(err)
	return taskRecords
}

func setTasksFile(newTaskFile string) {
	tasksFile = newTaskFile
}

const (
	timeLayoutDefault = "Jan 2 15:04:05"
)

// Function to format a task that has been started and ended
// To format a task that hasn't ended yet, pass nil for timeEnd
//func formatTaskLine(taskStamp TaskStamp) string {
//	if taskStamp.startOrEnd == START {
//
//		return fmt.Sprintf("%v started with priority %v at %v\n", taskStamp.task.name, strconv.Itoa(taskStamp.task.priority), taskStamp.timeStamp.Format(timeLayoutDefualt))
//	} else {
//		return fmt.Sprintf("%v ended at %v\n", taskStamp.task.name, taskStamp.timeStamp.Format(timeLayoutDefualt))
//	}
//}
