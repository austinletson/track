package main

import (
	"bufio"
	"os"

	flag "github.com/spf13/pflag"
)

type Command struct {
	flags    []string
	function func([]*flag.Flag)
}

func getStdin() (dataString string) {
	if isInputFromPipe() {
		scanner := bufio.NewScanner(bufio.NewReader(os.Stdin))
		for scanner.Scan() {
			dataString += scanner.Text()
		}
	}
	return dataString
}

func isInputFromPipe() bool {
	fileInfo, _ := os.Stdin.Stat()
	return fileInfo.Mode()&os.ModeCharDevice == 0
}
