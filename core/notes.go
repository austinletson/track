package core

import (
	"time"
)

type NotesRecord struct {
	Notes []Note
}

type Note struct {
	Tasks    []string
	Time     time.Time
	NoteText string
}

func TakeNote(note string, tasks []Task) {
	taskNames := []string{}
	for _, task := range tasks {
		taskNames = append(taskNames, task.Name)
	}
	newNote := Note{
		Tasks:    taskNames,
		Time:     time.Now(),
		NoteText: note,
	}

	notesRecord := ReadNotesFromNotesFile()

	if notesRecord.Notes == nil {
		notesRecord.Notes = []Note{}
	}

	notesRecord.Notes = append(notesRecord.Notes, newNote)
	WriteNotesToNotesFile(notesRecord)
}
