package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/Benzogang-Tape/CI-CD-example/internal/repository"
	"github.com/Benzogang-Tape/CI-CD-example/internal/transport/handlers"
)

func main() {
	notesRepo := repository.NewNoteRepo()
	h := handlers.NewHandler(notesRepo)

	r := mux.NewRouter()
	r.HandleFunc("/note/{ID}", h.Get).Methods(http.MethodGet)
	r.HandleFunc("/note", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/note/{ID}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/note/{ID}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/note", h.GetAll).Methods(http.MethodGet)

	addr := ":8080"
	log.Printf("Starting server on %s\n", addr)
	log.Panic(http.ListenAndServe(addr, r))
}
