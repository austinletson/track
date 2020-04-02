package main

import (
	"fmt"
	color "github.com/fatih/color"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const (
	timeLayout  = "15:04"
	_TIOCGWINSZ = 0x5413
)

func generateBasicDisplay(records TaskRecords) (chart string) {

	for _, task := range records.Record {
		namePrefix := fmt.Sprintf("[%v]", task.Name)

		startTime := task.TaskIntervals[0].StartTime.Format(timeLayout)
		endTime := task.TaskIntervals[0].EndTime.Format(timeLayout)

		var timeString string
		if task.TaskIntervals[0].EndTime == NIL_TIME {
			timeString = fmt.Sprintf("%v ---->", startTime)
		} else {
			timeString = fmt.Sprintf("%v -- %v", startTime, endTime)
		}

		lineString := fmt.Sprintf("%v %v\n", namePrefix, timeString)
		chart = chart + lineString

	}
	return chart

}
func generateGanttChart(records TaskRecords) (chart string) {

	var allTimes []time.Time
	maxNameLength := 0
	ganttGraphLength := TerminalWidth() - 20

	for _, record := range records.Record {
		for _, interval := range record.TaskIntervals {
			allTimes = append(allTimes, interval.StartTime, interval.EndTime)
		}
		if len(record.Name) > maxNameLength {
			maxNameLength = len(record.Name)
		}

	}
	minTime, maxTime := timeSpan(allTimes)

	totalLength := ganttGraphLength + maxNameLength + 3
	chart = generateChartHeader(minTime, maxTime, totalLength)

	for _, task := range records.Record {
		lineString := Bold(fmt.Sprintf(" %v: %v", task.Name, generateCharacters(" ", maxNameLength-len(task.Name))))
		previousEnd := -1
		for _, interval := range task.TaskIntervals {
			startTime := interval.StartTime
			endTime := interval.EndTime

			scaledLocationOfStart := fractionOfTimeRange(startTime, minTime, maxTime, ganttGraphLength)
			var scaledLocationOfEnd int
			isActive := false
			if endTime == NIL_TIME {
				scaledLocationOfEnd = ganttGraphLength
				isActive = true
			} else {
				scaledLocationOfEnd = fractionOfTimeRange(endTime, minTime, maxTime, ganttGraphLength)
			}

			currentChar := " "
			for i := previousEnd + 1; i <= scaledLocationOfEnd; i++ {
				if i == scaledLocationOfEnd && isActive {
					currentChar = color.New(color.FgCyan).Add(color.Bold).Sprint(">")
				} else if i == scaledLocationOfStart && i == 0 {
					currentChar = "|"
				} else if i == scaledLocationOfStart && scaledLocationOfStart == scaledLocationOfEnd {
					currentChar = "*"
				} else if i == scaledLocationOfStart || i == scaledLocationOfEnd {
					currentChar = color.New(color.FgCyan).Add(color.Bold).Sprint("|")
				} else if i > scaledLocationOfStart && i < scaledLocationOfEnd {
					currentChar = color.New(color.FgHiRed).Add(color.Bold).Sprint("=")
				} else {
					currentChar = " "
				}
				lineString += currentChar
			}
			previousEnd = scaledLocationOfEnd
		}

		chart = chart + lineString + "\n"

		//lineString := fmt.Sprintf("%v %v\n", namePrefix, timeString)
		//chart = chart + lineString

	}
	//footer := "#" + generateCharacters("#", totalLength-2) + "#"
	//	chart = chart + footer
	return Bold(chart)

}

func generateChartHeader(startTime time.Time, endTime time.Time, width int) (headerString string) {

	//	topHeader := "#" + generateCharacters("#", width-2) + "#\n"
	differenceTimeUnix := endTime.Unix() - startTime.Unix()

	fourthWayTime := time.Unix(differenceTimeUnix/4+startTime.Unix(), 0).Format(timeLayout)
	halfWayTime := time.Unix(differenceTimeUnix/2+startTime.Unix(), 0).Format(timeLayout)
	threeFourthsWayTime := time.Unix(3*differenceTimeUnix/4+startTime.Unix(), 0).Format(timeLayout)

	widthWithOffset := width - 5*5 - 8

	firstThirdOfCharacters := widthWithOffset/4 + widthWithOffset%4
	secondAndThirdOfCharacters := widthWithOffset / 4

	firstChars := generateCharacters("-", firstThirdOfCharacters)
	secondAndThirdChars := generateCharacters("-", secondAndThirdOfCharacters)

	headerString = fmt.Sprintf(" %v %v %v %v %v %v %v %v %v \n", startTime.Format(timeLayout), firstChars, fourthWayTime, secondAndThirdChars, halfWayTime, secondAndThirdChars, threeFourthsWayTime, secondAndThirdChars, endTime.Format(timeLayout))

	return headerString
}

func generateCharacters(character string, count int) (characters string) {
	for i := 0; i < count; i++ {
		characters += character
	}
	return characters
}

// Return min and max of slice of times
func timeSpan(times []time.Time) (min time.Time, max time.Time) {
	min = times[0]
	max = times[0]
	for _, time := range times {
		if time.Before(min) && time.After(NIL_TIME) {
			min = time
		} else if time.After(max) {
			max = time
		}
	}
	return min, max
}

// Returns the scalled numerator of how far between max and min the middle time lies
func fractionOfTimeRange(middle time.Time, min time.Time, max time.Time, denominator int) (scaledFraction int) {
	totalUnixTimeOfRange := max.Unix() - min.Unix()
	middleUnixTimeOfRange := middle.Unix() - min.Unix()

	decimal := float64(middleUnixTimeOfRange) / float64(totalUnixTimeOfRange)
	scaledFraction = int(float64(denominator) * decimal)
	return scaledFraction
}

func Bold(str string) string {
	return "\033[1m" + str + "\033[0m"
}

func TerminalWidth() int {
	cmd := exec.Command("stty", "size")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	splitOnSpace := strings.Split(string(out), " ")[1]
	width, err := strconv.Atoi(splitOnSpace[:len(string(splitOnSpace))-1])
	return width
}
