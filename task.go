package main

import (
	"errors"
	"fmt"
	"time"
)

const (
	START = 1
	END   = 0
)

var NIL_TIME time.Time = time.Time{}

type Task struct {
	Name          string
	Priority      int
	TaskIntervals []TaskInterval
}

type TaskInterval struct {
	StartTime time.Time
	EndTime   time.Time
}

type TaskStartStamp struct {
	Task       Task
	TimeStamp  time.Time
	StartOrEnd int
}

type TaskEndStamp struct {
	TimeStamp time.Time
}

type TaskRecords struct {
	Record map[string]*Task
}

var activeTasks = make([]Task, 0)

func makeTask(name string, priority int) (task Task) {

	task = Task{
		Name:          name,
		Priority:      priority,
		TaskIntervals: nil,
	}
	return task
}

// Function to start a task and write it to file
func clockIn(taskValue Task, timeStampValue time.Time) error {
	taskRecords := ReadTasksFromTasksFile()

	if taskRecords.Record == nil {
		taskRecords.Record = make(map[string]*Task)
	}
	taskExistsInRecord := false
	if task, ok := taskRecords.Record[taskValue.Name]; ok {
		lastInterval := task.TaskIntervals[len(task.TaskIntervals)-1]
		if lastInterval.EndTime == NIL_TIME {
			return errors.New("Task already active")
		}
		newInterval := TaskInterval{
			StartTime: time.Now(),
			EndTime:   NIL_TIME,
		}
		task.TaskIntervals = append(task.TaskIntervals, newInterval)
		fmt.Print(task.TaskIntervals)
		taskExistsInRecord = true

	}

	if !taskExistsInRecord {
		taskInterval := TaskInterval{
			StartTime: timeStampValue,
			EndTime:   NIL_TIME,
		}
		taskValue.TaskIntervals = append(taskValue.TaskIntervals, taskInterval)
		taskRecords.Record[taskValue.Name] = &taskValue
	}

	WriteTasksToTaskFile(taskRecords)
	return nil
}

// Mabeu consolidate this function
// Function to end a task that has already been started
// TODO figure out weird behavior with task lastInterval
func clockOut(taskName string, timeStampValue time.Time) error {
	taskRecords := ReadTasksFromTasksFile()
	task := taskRecords.Record[taskName]
	if task.TaskIntervals[len(task.TaskIntervals)-1].EndTime == NIL_TIME {

		task.TaskIntervals[len(task.TaskIntervals)-1].EndTime = time.Now()
		WriteTasksToTaskFile(taskRecords)
		return nil
	}
	return errors.New("Task is not active")
}

// Function to lists active tasks
// Maybe consolidate with get active tasks
func listTasks(taskRecords TaskRecords) {
	for _, task := range getActiveTasks(taskRecords) {
		fmt.Print(task.Name)
	}
}

// Helper methods

// Gets all tasks that are most recent and have tag START
func getActiveTasks(taskRecords TaskRecords) (activeTasks []Task) {
	for _, task := range taskRecords.Record {
		if task.TaskIntervals[len(task.TaskIntervals)-1].EndTime == NIL_TIME {
			activeTasks = append(activeTasks, *task)
		}

	}
	return activeTasks
}

// Checks if a slice of tasks contains another task based on name
func containsTask(s []Task, e Task) bool {
	for _, a := range s {
		if a.Name == e.Name {
			return true
		}
	}
	return false
}
