package core

import (
	"errors"
	"sort"
	"time"

	"github.com/austinletson/track/common"
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

func UpdateTask(taskRecords TaskRecords, name string, priority int, tags []string) TaskRecords {
	if taskRecords.Record == nil {
		taskRecords.Record = make(map[string]*Task)
	}
	if taskRecords.Record[name] == nil {
		newTask := MakeTask(name, priority, tags)
		taskRecords.Record[name] = &newTask
		return taskRecords
	}
	existingTags := taskRecords.Record[name].Tags
	newTags := []string{}
	for _, tag := range tags {
		if !common.ContainsString(tag, existingTags) {
			newTags = append(newTags, tag)
		}
	}
	taskRecords.Record[name].Tags = append(existingTags, newTags...)
	taskRecords.Record[name].Priority = priority
	return taskRecords
}

// Function to start a task and write it to file
func ClockIn(taskRecords TaskRecords, taskName string, timeStampValue time.Time) (newTaskRecords TaskRecords, startError error) {
	taskValue := taskRecords.Record[taskName]

	if taskRecords.Record[taskValue.Name].TaskIntervals != nil {
		task := taskRecords.Record[taskValue.Name]
		lastInterval := task.TaskIntervals[len(task.TaskIntervals)-1]
		if lastInterval.EndTime == NIL_TIME {
			return taskRecords, errors.New("Task " + task.Name + " already active")
		}
		if len(task.TaskIntervals) > 1 {
			lastTimeStamp := task.TaskIntervals[len(task.TaskIntervals)-1].EndTime
			if timeStampValue.Before(lastTimeStamp) {
				return taskRecords, errors.New("Cannot start task " + task.Name + " before the end of the last time interval")
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
		taskRecords.Record[taskValue.Name] = taskValue
	}
	return taskRecords, nil
}

// Function to end a task that has already been started
func ClockOut(taskRecords TaskRecords, taskName string, timeStampValue time.Time) (records TaskRecords, endError error) {
	task := taskRecords.Record[taskName]
	if task == nil {
		return taskRecords, errors.New("Task " + taskName + " does not exist")
	} else if task.TaskIntervals[len(task.TaskIntervals)-1].EndTime != NIL_TIME {
		return taskRecords, errors.New("Task " + taskName + " is not active")
	} else if timeStampValue.Before(task.TaskIntervals[len(task.TaskIntervals)-1].StartTime) {
		return taskRecords, errors.New("Cannot end task " + taskName + " before it started")
	}

	task.TaskIntervals[len(task.TaskIntervals)-1].EndTime = timeStampValue
	taskRecords.Record[taskName] = task
	return taskRecords, endError
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

// Sorts tasks first by active/inactive and then by last activity
func sortTasksByTime(tasks []Task) []Task {
	activeTasks := []Task{}
	inactiveTasks := []Task{}
	for _, task := range tasks {
		if task.TaskIntervals[len(task.TaskIntervals)-1].EndTime == NIL_TIME {
			activeTasks = append(activeTasks, task)
		} else {
			inactiveTasks = append(inactiveTasks, task)
		}
	}

	sort.Slice(activeTasks, func(i, j int) bool {
		timeI := activeTasks[i].TaskIntervals[len(activeTasks[i].TaskIntervals)-1]
		timeJ := activeTasks[j].TaskIntervals[len(activeTasks[j].TaskIntervals)-1]

		return timeI.StartTime.After(timeJ.StartTime)
	})
	sort.Slice(inactiveTasks, func(i, j int) bool {
		timeI := inactiveTasks[i].TaskIntervals[len(inactiveTasks[i].TaskIntervals)-1]
		timeJ := inactiveTasks[j].TaskIntervals[len(inactiveTasks[j].TaskIntervals)-1]

		return timeI.EndTime.After(timeJ.EndTime)
	})

	return append(activeTasks, inactiveTasks...)
}
