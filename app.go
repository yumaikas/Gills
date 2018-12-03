package main

import (
	"github.com/go-chi/chi"
	"net/http"
)

func Route(r chi.Router) {
	r.Get("/", Search)
	r.Get("/admin/", Home)
	r.Post("/admin/create", CreateNote)
	r.Post("/admin/note/save/{noteID}", DoSaveNote)
}

// TODO: This is going to provide HTTP handlers that are routed by main.go
// And then it's going to load state from the database and feed it into
// the various routers.
func Home(w http.ResponseWriter, r *http.Request) {
	// TODO Merge query search parameters into this
	appState, err := LoadState()
	die(err)
	searchTerms := appState.getOr("recentSearchTerms", "")
	notes, err := SearchNotes(searchTerms)
	die(err)
	die(HomeView(w, appState, notes))
}

func Search(w http.ResponseWriter, r *http.Request) {
	// TODO:
	panic("NOT DONE!")
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	panic("NOT DONE!")
}

func DoSaveNote(w http.ResponseWriter, r *http.Request) {
	// appState, err := LoadState()
	// die(err)
	panic("NOT DONE!")
}

func die(e error) {
	if e != nil {
		panic(e)
	}
}
