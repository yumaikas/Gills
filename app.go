package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"time"
)

func Route(r chi.Router) {
	r.Get("/", Search)
	r.Get("/admin", Home)
	r.Get("/admin/", Home)
	r.Post("/admin/search", DoSaveState)
	r.Post("/admin/create", CreateNote)
	r.Post("/admin/note/save/{noteID}", DoSaveNote)
	r.Delete("/admin/note/{noteID}", DoDeleteNote)
}

// TODO: This is going to provide HTTP handlers that are routed by main.go
// And then it's going to load state from the database and feed it into
// the various routers.
func Home(w http.ResponseWriter, r *http.Request) {
	appState, err := LoadState()
	die(err)
	searchTerms := appState.getOr(recentSearchKey, "")
	notes, err := SearchRecentNotes(searchTerms)
	die(err)
	die(HomeView(w, appState, notes))
}

func Search(w http.ResponseWriter, r *http.Request) {
	appState, err := LoadState()
	die(err)
	searchTerms := appState.getOr(recentSearchKey, "")
	notes, err := SearchNotes(searchTerms)
	die(SearchView(w, appState, notes))
}

func DoSaveState(w http.ResponseWriter, r *http.Request) {
	die(saveMainPageState(r))
	http.Redirect(w, r, "/admin/", 301)
}

func CreateNote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	content := r.PostFormValue(draftnoteKey)
	tagline := r.PostFormValue(taglineKey)
	var note = Note{
		Id:      0, // used to signal that this note does *not* have a correspoding database row
		Content: content + "\n" + tagline,
		Created: time.Now(),
	}
	var err error
	note.Id, err = SaveNote(note)
	die(err)
	// Remove the draftnote so that it gets cleared out on savestate
	//
	r.Form.Del("draftnote")
	die(saveMainPageState(r))
	http.Redirect(w, r, "/admin/", 301)
}

func DoSaveNote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	content := r.PostFormValue("draftnote")
	strId := chi.URLParam(r, "noteID")
	appState, err := LoadState()
	die(err)
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		InvalidSaveIdView(w, appState, strId)
		return
	}

	var note = Note{
		Id:      id, // used to signal that this note does *not* have a correspoding database row
		Content: content,
		Created: time.Now(),
	}
	note.Id, err = SaveNote(note)
	die(err)
	die(saveMainPageState(r))
	http.Redirect(w, r, "/admin/", 301)
}

func DoDeleteNote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	strId := chi.URLParam(r, "noteID")
	appState, err := LoadState()
	die(err)

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		InvalidDeleteIdView(w, appState, strId)
		return
	}
	die(DeleteNote(id))
}

func saveMainPageState(r *http.Request) error {
	state := []KV{
		KV{scratchpadKey, r.FormValue(scratchpadKey)},
		KV{draftnoteKey, r.FormValue(draftnoteKey)},
		KV{taglineKey, r.FormValue(taglineKey)},
		KV{recentSearchKey, r.FormValue(recentSearchKey)},
	}
	return SaveState(state)
}

func die(e error) {
	if e != nil {
		panic(e)
	}
}
