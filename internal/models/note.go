package models

import (
	"cmp"
	"slices"
	"strings"
	"sync/atomic"
	"time"
)

type Note struct {
	ID        uint64    `json:"id"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type NotePayload struct {
	Text string `json:"text"`
}

var counter uint64

func NewNote(text string) *Note {
	return &Note{
		ID:        atomic.AddUint64(&counter, 1),
		Text:      text,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

func SortNotes(notes []Note, field string) []Note {
	sortByField := func(sortFunc func(a, b Note) int) {
		slices.SortFunc(notes, sortFunc)
	}

	switch field {
	case "id":
		sortByField(compareNotesByID)
	case "text":
		sortByField(compareNotesByText)
	case "created_at":
		sortByField(compareNotesByCreationTime)
	case "updated_at":
		sortByField(compareNotesByUpdateTime)
	}

	return notes
}

func compareNotesByID(a, b Note) int {
	return cmp.Compare(a.ID, b.ID)
}

func compareNotesByText(a, b Note) int {
	return strings.Compare(strings.ToLower(a.Text), strings.ToLower(b.Text))
}

func compareNotesByCreationTime(a, b Note) int {
	return a.CreatedAt.Compare(b.CreatedAt)
}

func compareNotesByUpdateTime(a, b Note) int {
	return a.UpdatedAt.Compare(b.UpdatedAt)
}
