package main

import "fmt"

const (
	timeLayout = "15:04:05"
)

func generateGanttChart(records TaskRecords) (chart string) {

	for _, task := range records.Record {
		namePrefix := fmt.Sprintf("[%v]", task.Name)

		startTime := task.TaskIntervals[0].StartTime.Format(timeLayout)
		endTime := task.TaskIntervals[0].EndTime.Format(timeLayout)

		timeString := fmt.Sprintf("%v -- %v", startTime, endTime)

		lineString := fmt.Sprintf("%v %v\n", namePrefix, timeString)
		chart = chart + lineString

	}
	return chart

}
