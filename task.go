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

type Task struct {
	Name     string
	Priority int
}

type TaskInterval struct {
	Task      Task
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
	Record map[string][]TaskStartStamp
}

var activeTasks = make([]Task, 0)

// Function to start a task and write it to file
func clockIn(taskValue Task, timeStampValue time.Time) error {
	for _, task := range getActiveTasks() {
		if task.Name == taskValue.Name {
			return errors.New("Task already active")
		}
	}

	taskStamp := TaskStartStamp{
		Task:       taskValue,
		TimeStamp:  time.Now(),
		StartOrEnd: START,
	}

	fmt.Print(taskStamp)
	WriteTasksToTaskFile(taskStamp)
	return nil
}

// Mabeu consolidate this function
// Function to end a task that has already been started
func clockOut(taskName string, timeStampValue time.Time) error {
	for _, task := range getActiveTasks() {
		if task.Name == taskName {
			taskStamp := TaskStartStamp{
				Task: Task{
					Name:     taskName,
					Priority: 0,
				},
				TimeStamp:  timeStampValue,
				StartOrEnd: END,
			}

			WriteTasksToTaskFile(taskStamp)
			return nil
		}
	}
	return errors.New("Task is not active")
}

// Function to lists active tasks
func listTasks(tasksStamps map[string][]TaskStartStamp) {
	for _, task := range getActiveTasks() {
		fmt.Print(task.Name)
	}
}

// Helper methods

// Gets all tasks that are most recent and have tag START
func getActiveTasks() (activeTasks []Task) {
	taskRecord := ReadTasksFromTasksFile()
	for _, taskStamps := range taskRecord {
		if len(taskStamps) != 0 && taskStamps[len(taskStamps)-1].StartOrEnd == START {
			activeTasks = append(activeTasks, taskStamps[len(taskStamps)-1].Task)
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
