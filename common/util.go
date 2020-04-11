package common

import (
	"fmt"
	"strings"
)

func IsError(err error) bool {
	if err != nil {
		fmt.Println(err.Error())
	}
	return (err != nil)
}

func ContainsString(str string, strs []string) bool {
	for _, s := range strs {
		if s == str {
			return true
		}
	}
	return false
}

func TrimList(stringList []string) (trimmedStrings []string) {
	for _, str := range stringList {
		trimmedStrings = append(trimmedStrings, strings.Trim(str, " "))
	}
	return trimmedStrings
}
