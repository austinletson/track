package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

const (
	storageFile = "tasks.tt"
)

func WriteTasksToTaskFile(taskStamp TaskStartStamp) {

	taskName := taskStamp.Task.Name
	var newRecord TaskRecords
	var _, err = os.Stat(storageFile)
	if err == nil {
		previousRecord := ReadTasksFromTasksFile()
		if previousRecord == nil {
			previousRecord = make(map[string][]TaskStartStamp)
		}
		if previousRecord[taskStamp.Task.Name] == nil {
			// Add a new stamp list
			previousRecord[taskName] = []TaskStartStamp{taskStamp}
		} else {
			// Add the stamp to the old list
			oldList := previousRecord[taskName]
			previousRecord[taskName] = append(oldList, taskStamp)
		}
		newRecord.Record = previousRecord
		err = os.Remove(storageFile)
		isError(err)
	} else {
		fmt.Print("Generating tasks file...")
		mapRecord := make(map[string][]TaskStartStamp)
		mapRecord[taskStamp.Task.Name] = []TaskStartStamp{taskStamp}
		newRecord.Record = mapRecord

	}

	file, err := os.Create(storageFile)
	isError(err)
	defer file.Close()

	formattedTask, err := yaml.Marshal(newRecord)
	isError(err)

	//	writer := bufio.NewWriter(file)
	//	_, err = writer.Write(formattedTask)
	//	isError(err)
	//err = writer.Flush()
	//	isError(err)

	err = ioutil.WriteFile(storageFile, formattedTask, 0644)
	isError(err)
}

// Read tasks edn from file and return them
func ReadTasksFromTasksFile() map[string][]TaskStartStamp {
	file, err := os.Open(storageFile)
	isError(err)
	defer file.Close()

	var bytes []byte
	bytes, err = ioutil.ReadAll(file)
	isError(err)

	var stampsList TaskRecords
	err = yaml.Unmarshal(bytes, &stampsList)
	isError(err)
	return stampsList.Record
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
