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
func ClockIn(taskValues []Task, timeStampValue time.Time) (startErrors []error) {
	taskRecords := ReadTasksFromTasksFile()
	if taskRecords.Record == nil {
		taskRecords.Record = make(map[string]*Task)
	}
	for _, taskValue := range taskValues {
		if task, ok := taskRecords.Record[taskValue.Name]; ok {
			lastInterval := task.TaskIntervals[len(task.TaskIntervals)-1]
			if lastInterval.EndTime == NIL_TIME {
				startError := errors.New("Task " + task.Name + " already active")
				startErrors = append(startErrors, startError)
				continue
			}
			if len(task.TaskIntervals) > 1 {
				lastTimeStamp := task.TaskIntervals[len(task.TaskIntervals)-1].EndTime
				if timeStampValue.Before(lastTimeStamp) {
					startError := errors.New("Cannot start task " + task.Name + " before the end of the last time interval")
					startErrors = append(startErrors, startError)
					continue
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
			// Make a copy because of for loop pointer behavior
			taskValueNew := taskValue
			taskRecords.Record[taskValue.Name] = &taskValueNew
		}
	}
	WriteTasksToTaskFile(taskRecords)
	return startErrors
}

// Mabeu consolidate this function
// Function to end a task that has already been started
// TODO figure out weird behavior with task lastInterval
func ClockOut(taskNames []string, timeStampValue time.Time) (endErrors []error) {
	taskRecords := ReadTasksFromTasksFile()
	for _, taskName := range taskNames {
		task := taskRecords.Record[taskName]

		if task.TaskIntervals[len(task.TaskIntervals)-1].EndTime != NIL_TIME {
			endError := errors.New("Task " + taskName + " is not active")
			endErrors = append(endErrors, endError)
		}

		if timeStampValue.Before(task.TaskIntervals[len(task.TaskIntervals)-1].StartTime) {
			endError := errors.New("Cannot end task " + taskName + " before it started")
			endErrors = append(endErrors, endError)
		}

		task.TaskIntervals[len(task.TaskIntervals)-1].EndTime = timeStampValue
	}
	WriteTasksToTaskFile(taskRecords)
	return endErrors
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
