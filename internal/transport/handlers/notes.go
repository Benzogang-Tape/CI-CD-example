package handlers

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"github.com/Benzogang-Tape/CI-CD-example/internal/models"
)

//go:generate mockgen -source=notes.go -destination=../../repository/mocks/notes_repo_mock.go -package=mocks NotesAPI
type NotesAPI interface {
	CreateNote(ctx context.Context, text string) (models.Note, error)
	GetNoteByID(ctx context.Context, id uint64) (models.Note, error)
	GetAllNotes(ctx context.Context, orderBy string) ([]models.Note, error)
	DeleteNote(ctx context.Context, id uint64) error
	UpdateNote(ctx context.Context, id uint64, newText string) (models.Note, error)
}

type Handler struct {
	repo NotesAPI
}

func NewHandler(n NotesAPI) *Handler {
	return &Handler{n}
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	notePayload := models.NotePayload{}
	if err = json.Unmarshal(body, &notePayload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newNote, err := h.repo.CreateNote(r.Context(), notePayload.Text)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(newNote)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusCreated)
	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	note, err := h.repo.GetNoteByID(r.Context(), uint64(id))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := h.repo.DeleteNote(r.Context(), uint64(id)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["ID"])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	defer r.Body.Close()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	notePayload := models.NotePayload{}
	if err = json.Unmarshal(body, &notePayload); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	note, err := h.repo.UpdateNote(r.Context(), uint64(id), notePayload.Text)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	resp, err := json.Marshal(note)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}

	if _, err := w.Write(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (h *Handler) GetAll(w http.ResponseWriter, r *http.Request) {
	orderField := mux.Vars(r)["order_by"]

	notes, err := h.repo.GetAllNotes(r.Context(), orderField)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err := json.NewEncoder(w).Encode(notes); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
