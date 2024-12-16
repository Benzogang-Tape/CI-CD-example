package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Benzogang-Tape/CI-CD-example/internal/models"
	"github.com/Benzogang-Tape/CI-CD-example/internal/models/errs"
	"github.com/Benzogang-Tape/CI-CD-example/internal/repository/mocks"
)

type fakeBody struct {
	data string
}

func (f *fakeBody) Read(p []byte) (int, error) {
	return 0, errors.New("planned read error")
}

func (f *fakeBody) Close() error {
	return errors.New("planned close error")
}

var (
	note = models.Note{
		ID:   1,
		Text: "Text",
	}
)

const (
	notePayload = `{"text":"Text"}`
)

func TestGet(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNotesAPI(ctrl)
	handler := NewHandler(repo)

	// Success
	r := httptest.NewRequest("GET", fmt.Sprintf("/note/%d", note.ID), nil)
	r = mux.SetURLVars(r, map[string]string{
		"ID": strconv.Itoa(int(note.ID)),
	})
	w := httptest.NewRecorder()
	repo.EXPECT().GetNoteByID(r.Context(), note.ID).Return(note, nil)

	handler.Get(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	expectedData, _ := json.Marshal(note) //nolint:errcheck
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, expectedData, body)

	// Bad note id
	r = httptest.NewRequest("GET", "/note/a", nil)
	w = httptest.NewRecorder()

	handler.Get(w, r)
	resp = w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Note doesn't exist
	var fakeID uint64 = 54
	r = httptest.NewRequest("GET", fmt.Sprintf("/note/%d", fakeID), nil)
	r = mux.SetURLVars(r, map[string]string{
		"ID": strconv.Itoa(int(fakeID)),
	})
	w = httptest.NewRecorder()
	repo.EXPECT().GetNoteByID(r.Context(), fakeID).Return(models.Note{}, errs.ErrNoNote)

	handler.Get(w, r)
	resp = w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}

func TestCreate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mocks.NewMockNotesAPI(ctrl)
	handler := NewHandler(repo)

	// Success
	repo.EXPECT().CreateNote(context.Background(), note.Text).Return(note, nil)
	r := httptest.NewRequest("POST", "/note", strings.NewReader(notePayload))
	w := httptest.NewRecorder()

	handler.Create(w, r)
	resp := w.Result()
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)

	expectedData, _ := json.Marshal(note) //nolint:errcheck
	require.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)
	assert.Equal(t, expectedData, body)

	// Read body error
	r = httptest.NewRequest("POST", "/note", bytes.NewReader(nil))
	r.Body = &fakeBody{data: notePayload}
	w = httptest.NewRecorder()
	handler.Create(w, r)
	resp = w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Unmarshal body error
	r = httptest.NewRequest("POST", "/api/posts", bytes.NewReader(nil))
	w = httptest.NewRecorder()
	handler.Create(w, r)
	resp = w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)

	// Success
	repo.EXPECT().CreateNote(context.Background(), note.Text).Return(models.Note{}, errors.New("internal"))
	r = httptest.NewRequest("POST", "/note", strings.NewReader(notePayload))
	w = httptest.NewRecorder()

	handler.Create(w, r)
	resp = w.Result()
	defer resp.Body.Close()

	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)
}
