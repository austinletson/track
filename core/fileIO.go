package core

import (
	"fmt"
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"

	"github.com/austinletson/track/common"
)

const (
	configDirPath = "~/.config/track/"
)

var storageFile string = "tasks.tt"
var notesFile string = "notes"

func WriteTasksToTaskFile(taskRecords TaskRecords) {

	if _, err := os.Stat(storageFile); os.IsExist(err) {
		os.Remove(storageFile)
	}

	file, err := os.Create(storageFile)
	common.IsError(err)
	defer file.Close()

	formattedTask, err := yaml.Marshal(taskRecords)
	common.IsError(err)

	err = ioutil.WriteFile(storageFile, formattedTask, 0644)
	common.IsError(err)
}

// Read tasks edn from file and return them
func ReadTasksFromTasksFile() (taskRecords TaskRecords) {
	if _, err := os.Stat(storageFile); os.IsNotExist(err) {
		return taskRecords
	}
	file, err := os.Open(storageFile)
	common.IsError(err)
	defer file.Close()

	var bytes []byte
	bytes, err = ioutil.ReadAll(file)
	common.IsError(err)

	err = yaml.Unmarshal(bytes, &taskRecords)
	common.IsError(err)
	return taskRecords
}

func WriteNotesToNotesFile(notesRecord NotesRecord) error {
	if _, err := os.Stat(notesFile); os.IsExist(err) {
		os.Remove(storageFile)
	}

	file, err := os.Create(notesFile)
	common.IsError(err)
	defer file.Close()

	formattedNote, err := yaml.Marshal(notesRecord)
	common.IsError(err)

	fmt.Println(notesRecord)
	fmt.Println(string(formattedNote))
	var notesRecordCheck NotesRecord
	err = yaml.Unmarshal(formattedNote, &notesRecordCheck)
	common.IsError(err)

	err = ioutil.WriteFile(notesFile, formattedNote, 0644)
	common.IsError(err)
	return nil
}

// Read tasks edn from file and return them
func ReadNotesFromNotesFile() (notesRecord NotesRecord) {
	if _, err := os.Stat(notesFile); os.IsNotExist(err) {
		return notesRecord
	}
	file, err := os.Open(notesFile)
	common.IsError(err)
	defer file.Close()

	var bytes []byte
	bytes, err = ioutil.ReadAll(file)
	common.IsError(err)

	err = yaml.Unmarshal(bytes, &notesRecord)
	common.IsError(err)
	return notesRecord
}
