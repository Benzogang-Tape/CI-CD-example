package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/Benzogang-Tape/CI-CD-example/internal/models"
	"github.com/Benzogang-Tape/CI-CD-example/internal/models/errs"
)

type NoteRepo struct {
	mu    sync.RWMutex
	notes map[uint64]*models.Note
}

func NewNoteRepo() *NoteRepo {
	return &NoteRepo{
		mu:    sync.RWMutex{},
		notes: make(map[uint64]*models.Note),
	}
}

func (r *NoteRepo) CreateNote(ctx context.Context, text string) (models.Note, error) { //nolint:unparam
	newNote := models.NewNote(text)
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notes[newNote.ID] = newNote
	return *newNote, nil
}

func (r *NoteRepo) GetNoteByID(ctx context.Context, id uint64) (models.Note, error) { //nolint:unparam
	r.mu.RLock()
	defer r.mu.RUnlock()
	note, ok := r.notes[id]
	if !ok {
		return models.Note{}, errs.ErrNoNote
	}
	return *note, nil
}

func (r *NoteRepo) GetAllNotes(ctx context.Context, orderBy string) ([]models.Note, error) { //nolint:unparam
	r.mu.RLock()
	defer r.mu.RUnlock()
	notes := make([]models.Note, 0, len(r.notes))
	for _, note := range r.notes {
		notes = append(notes, *note)
	}

	notes = models.SortNotes(notes, orderBy)

	return notes, nil
}

func (r *NoteRepo) DeleteNote(ctx context.Context, id uint64) error { //nolint:unparam
	r.mu.Lock()
	defer r.mu.Unlock()
	delete(r.notes, id)
	return nil
}

func (r *NoteRepo) UpdateNote(ctx context.Context, id uint64, newText string) (models.Note, error) {
	note, err := r.GetNoteByID(ctx, id)
	if err != nil {
		return models.Note{}, fmt.Errorf("UpdateNote: %w", err)
	}

	note.Text = newText
	r.mu.Lock()
	defer r.mu.Unlock()
	r.notes[id] = &note
	return note, nil
}
