package core

import (
	"errors"
	"time"
)

const (
	START = 1
	END   = 0
)

var NIL_TIME time.Time = time.Time{}
var NIL_RECORDS TaskRecords

type Task struct {
	Name          string
	Priority      int
	Tags          []string
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

func MakeTask(name string, priority int, tags []string) (task Task) {

	task = Task{
		Name:          name,
		Priority:      priority,
		Tags:          tags,
		TaskIntervals: nil,
	}
	return task
}

// Function to start a task and write it to file
func ClockIn(taskValue Task, timeStampValue time.Time) error {
	taskRecords := ReadTasksFromTasksFile()

	if taskRecords.Record == nil {
		taskRecords.Record = make(map[string]*Task)
	}
	if task, ok := taskRecords.Record[taskValue.Name]; ok {
		lastInterval := task.TaskIntervals[len(task.TaskIntervals)-1]
		if lastInterval.EndTime == NIL_TIME {
			return errors.New("Task already active")
		}
		if len(task.TaskIntervals) > 1 {
			lastTimeStamp := task.TaskIntervals[len(task.TaskIntervals)-1].EndTime
			if timeStampValue.Before(lastTimeStamp) {
				return errors.New("Cannot start task before the end of the last time interval")
			}
		}
		newInterval := TaskInterval{
			StartTime: timeStampValue,
			EndTime:   NIL_TIME,
		}
		task.TaskIntervals = append(task.TaskIntervals, newInterval)

	} else { // If the task doesn't exist in the record
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
func ClockOut(taskName string, timeStampValue time.Time) error {
	taskRecords := ReadTasksFromTasksFile()
	task := taskRecords.Record[taskName]

	if task.TaskIntervals[len(task.TaskIntervals)-1].EndTime != NIL_TIME {
		return errors.New("Task is not active")
	}

	if timeStampValue.Before(task.TaskIntervals[len(task.TaskIntervals)-1].StartTime) {
		return errors.New("Cannot end task before it started")
	}

	task.TaskIntervals[len(task.TaskIntervals)-1].EndTime = timeStampValue
	WriteTasksToTaskFile(taskRecords)
	return nil
}

// Helper methods

// Gets all tasks that are most recent and have tag START
func GetActiveTasks(taskRecords TaskRecords) (activeTasks []Task) {
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
