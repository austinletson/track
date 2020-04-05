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

	justActiveTasks = 0
	allTasks        = 1

	verbose    = 0
	notVerbose = 1
)

func listTasks(taskRecords TaskRecords, allOrActive int, verboseOrNot int) (taskString string) {
	var tasks []Task
	if allOrActive == justActiveTasks {
		tasks = getActiveTasks(taskRecords)
	} else {
		for _, task := range taskRecords.Record {
			tasks = append(tasks, *task)
		}
	}

	if verboseOrNot == notVerbose {
		for _, task := range tasks {
			taskString = taskString + task.Name + "\n"
		}
	} else {
		taskString = generateBasicReport(tasks)
	}

	return taskString
}

// Generate text report of all tasks
func generateBasicReport(tasks []Task) (report string) {

	// Agregate all time and calculate maxNameLength
	maxNameLength := 0
	for _, task := range tasks {
		if len(task.Name) > maxNameLength {
			maxNameLength = len(task.Name)
		}

	}

	for _, task := range tasks {
		var namePrefix string
		if task.Priority != 0 {
			namePrefix = Bold(fmt.Sprintf("[%v %v]", task.Name, task.Priority))
		} else {
			namePrefix = Bold(fmt.Sprintf("[%v]", task.Name))
		}

		var timeString string
		firstOne := true
		for _, interval := range task.TaskIntervals {
			startTime := interval.StartTime.Format(timeLayout)
			endTime := interval.EndTime.Format(timeLayout)

			if firstOne {
				timeString = generateCharacters(" ", maxNameLength-len(task.Name))
				firstOne = false
			} else {
				// +2 for two brackets and -4 to make it less indented
				timeString = timeString + generateCharacters(" ", maxNameLength+2-4)
			}
			if interval.EndTime == NIL_TIME {
				// TODO make this not specific to the current time format
				timeString += color.New(color.FgGreen).Sprintf("%v ------->\n", startTime)
			} else {
				timeString += fmt.Sprintf("%v -- %v\n", startTime, endTime)
			}

		}

		lineString := fmt.Sprintf("%v %v\n", namePrefix, timeString)
		report = report + lineString

	}
	return report

}

func generateGanttChart(records TaskRecords) (chart string) {

	// Agregate all time and calculate maxNameLength
	maxNameLength := 4 // 4 because the header "TASK" has 4 chars
	var allTimes []time.Time
	anyHavePriority := false
	for _, task := range records.Record {
		for _, interval := range task.TaskIntervals {
			allTimes = append(allTimes, interval.StartTime, interval.EndTime)
		}
		if len(task.Name) > maxNameLength {
			maxNameLength = len(task.Name)
		}
		if task.Priority != 0 {
			anyHavePriority = true
		}

	}

	ganttGraphLength := TerminalWidth() - maxNameLength - 10
	minTime, maxTime, containsActiveTask := timeSpan(allTimes)

	// If there is an active task extend the max time to show activity
	if containsActiveTask {
		oneAndOneFourthTimeInterval := time.Unix(maxTime.Unix()+(maxTime.Unix()-minTime.Unix())/4, 0)
		if time.Now().Unix() < maxTime.Unix() {
			maxTime = time.Now()
		} else {
			maxTime = oneAndOneFourthTimeInterval
		}
	}

	var headShift int
	if anyHavePriority {
		headShift = 6
	} else {
		headShift = 4
	}
	//totalLength := ganttGraphLength + maxNameLength + headShift
	chart = ""

	for _, task := range records.Record {
		var lineString string
		if task.Priority == 0 {
			lineString = Bold(fmt.Sprintf(" [%v] %v", task.Name, generateCharacters(" ", maxNameLength-len(task.Name))))
		} else {
			lineString = Bold(fmt.Sprintf(" [%v %v] %v", task.Name, task.Priority, generateCharacters(" ", maxNameLength-len(task.Name))))
		}
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

			// Edge case where a small interval(*) is close to another interval
			// TODO fix bug that this isn't cyan
			if scaledLocationOfStart == previousEnd {
				lineString = lineString[:len(lineString)-1]
				lineString += "|"
			}

			currentChar := " "
			for i := previousEnd + 1; i <= scaledLocationOfEnd; i++ {
				if i == scaledLocationOfEnd && isActive {
					currentChar = color.New(color.FgCyan).Add(color.Bold).Sprint(">")
				} else if i == scaledLocationOfStart && scaledLocationOfStart == scaledLocationOfEnd {
					currentChar = "*"
				} else if i == scaledLocationOfStart && i == 0 {
					currentChar = color.New(color.FgCyan).Add(color.Bold).Sprint("|")
				} else if i == scaledLocationOfStart || i == scaledLocationOfEnd {
					currentChar = color.New(color.FgCyan).Add(color.Bold).Sprint("|")
				} else if i > scaledLocationOfStart && i < scaledLocationOfEnd {
					if isActive {
						currentChar = color.New(color.FgHiGreen).Add(color.Bold).Sprint("=")
					} else {
						currentChar = color.New(color.FgHiRed).Add(color.Bold).Sprint("=")
					}
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
	//	chart = chart + foot
	shiftedHeader := " " + generateLogo() + generateCharacters(" ", maxNameLength+headShift-13) + generateChartHeader(minTime, maxTime, ganttGraphLength)
	return shiftedHeader + Bold(chart)

}

func generateLogo() (logo string) {
	logo = color.New(color.FgBlack, color.BgHiWhite).Sprint(Bold("===")) + color.New(color.FgHiRed, color.BgHiWhite).Sprint((Bold("track"))) + color.New(color.FgBlack, color.BgHiWhite).Sprint(Bold("==="))
	return logo
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

	headerString = fmt.Sprintf("%v %v %v %v %v %v %v %v %v \n", startTime.Format(timeLayout), firstChars, fourthWayTime, secondAndThirdChars, halfWayTime, secondAndThirdChars, threeFourthsWayTime, secondAndThirdChars, endTime.Format(timeLayout))

	return headerString
}

func generateCharacters(character string, count int) (characters string) {
	for i := 0; i < count; i++ {
		characters += character
	}
	return characters
}

// Return min and max of slice of times
func timeSpan(times []time.Time) (min time.Time, max time.Time, containsNileTime bool) {
	min = times[0]
	max = times[0]
	for _, time := range times {
		if time == NIL_TIME {
			containsNileTime = true
		} else if time.Before(min) {
			min = time
		} else if time.After(max) {
			max = time
		}
	}
	return min, max, containsNileTime
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
