package main

import (
	"fmt"
	"github.com/go-chi/chi"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func Route(r chi.Router) {
	r.Get("/", Search)
	r.Get("/admin", Home)
	r.Get("/admin/", Home)
	r.Post("/admin/search", DoSaveState)
	r.Post("/admin/create", CreateNote)

	// Have the GET and POST forms so that we can keep the URLs clean when we're just linking
	// vs coming here from the main admin page (which wants to save some state on the way)
	r.Get("/admin/note/{noteID}", ShowNote)
	r.Post("/admin/note/{noteID}/show", SaveStateAndRedirToNote)
	r.Post("/admin/note/{noteID}", DoSaveNote)
	r.Post("/admin/note/{noteID}/delete", DoDeleteNote)
	r.Get("/admin/upload", ShowUploadForm)
	r.Get("/admin/upload/{filename}", ShowUploadedFile)
	r.Post("/admin/upload", ProcessUpload)
	r.Get("/admin/upload/list", ShowUploadNotes)

	// Scripting stuffs
	r.Get("/admin/scripting/", ListLuaScripts)

	// Run/Test new lua scripts
	r.Get("/admin/scripts/new/", NewLuaScriptForm)
	r.Post("/admin/scripts/new/run", RunNewLuaScript)
	r.Post("/admin/scripts/new/", CreateLuaScript)

	// Edit scripts
	r.Get("/admin/scripts/edit/{script-name}", EditLuaScript)
	r.Post("/admin/scripts/edit/{script-name}", SaveLuaScript)

	// Run scripts that have been saved
	r.Get("/admin/pages/s/{script-name}/*", RunLuaScript)
	r.Post("/admin/pages/s/{script-name}/*", RunLuaPostScript)

	r.Get("/admin/wild-test*", TestWildCard)
}

func TestWildCard(w http.ResponseWriter, r *http.Request) {
	param := chi.URLParam(r, "*")
	fmt.Fprint(w, "The URL param was:")
	fmt.Fprint(w, param)
}

func Home(w http.ResponseWriter, r *http.Request) {
	appState, err := LoadState()
	die(err)
	searchTerms := appState.GetOr(recentSearchKey, "")
	notes, err := SearchRecentNotes(searchTerms)
	die(err)
	die(HomeView(w, appState, notes))
}

func ShowUploadForm(w http.ResponseWriter, r *http.Request) {
	appState, err := LoadState()
	searchTerms := r.URL.Query().Get("q")
	die(err)
	notes, err := SearchUploadNotes(searchTerms)
	die(err)
	die(UploadView(w, appState, searchTerms, notes))
}

func ShowUploadNotes(w http.ResponseWriter, r *http.Request) {
	state, err := LoadState()
	die(err)
	searchTerms := r.URL.Query().Get("q")
	notes, err := SearchUploadNotes(searchTerms)
	die(err)
	die(SearchView(w, state.AppName(), searchTerms, notes))
}

func ShowUploadedFile(w http.ResponseWriter, r *http.Request) {
	fileName := chi.URLParam(r, "filename")
	http.ServeFile(w, r, PathForName(fileName))
}

func ProcessUpload(w http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(1000000)
	forms := r.MultipartForm
	files := forms.File["upload"]
	fileNames := make([]string, len(files))
	for i, header := range files {
		file, err := files[i].Open()
		defer file.Close()
		die(err)
		parts := strings.Split(header.Filename, ".")
		ext := strings.ToLower(parts[len(parts)-1])
		name, err := SaveUploadedFile(file, ext)
		die(err)
		fileNames[i] = name
	}
	_, err = SaveNote(NoteForFileNames(fileNames))
	die(err)
	http.Redirect(w, r, "/admin/upload", 301)
}

func Search(w http.ResponseWriter, r *http.Request) {
	state, err := LoadState()
	searchTerms := r.URL.Query().Get("q")
	notes, err := SearchNotes(searchTerms)
	die(err)
	die(SearchView(w, state.AppName(), searchTerms, notes))
}

func DoSaveState(w http.ResponseWriter, r *http.Request) {
	die(saveMainPageState(r))
	http.Redirect(w, r, "/admin/", 301)
}

func SaveStateAndRedirToNote(w http.ResponseWriter, r *http.Request) {
	die(saveMainPageState(r))
	strId := chi.URLParam(r, "noteID")
	http.Redirect(w, r, fmt.Sprint("/admin/note/", strId), 301)
}

func ShowNote(w http.ResponseWriter, r *http.Request) {
	state, loadErr := LoadState()
	die(loadErr)
	strId := chi.URLParam(r, "noteID")
	id, convErr := strconv.ParseInt(strId, 10, 64)
	die(convErr)
	note, repoErr := GetNoteBy(id)
	die(repoErr)
	NoteDetailsView(w, state, note)
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
	r.Form.Set("draftnote", "")
	die(saveMainPageState(r))
	http.Redirect(w, r, "/admin/", 301)
}

func DoSaveNote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	content := r.PostFormValue("note-content")
	strId := chi.URLParam(r, "noteID")
	appState, err := LoadState()
	die(err)
	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		InvalidIdView(w,
			appState.AppName(),
			"Cannot save note with invalid id: ",
			strId)
		return
	}

	var note = Note{
		Id:      id,
		Content: content,
		Created: time.Now(),
	}
	note.Id, err = SaveNote(note)
	die(err)
	http.Redirect(w, r, fmt.Sprint("/admin/note/", note.Id), 301)
}

func DoDeleteNote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	strId := chi.URLParam(r, "noteID")
	appState, err := LoadState()
	die(err)

	id, err := strconv.ParseInt(strId, 10, 64)
	if err != nil {
		w.WriteHeader(400)
		InvalidIdView(w,
			appState.AppName(),
			"Cannot delete note from invalid id:",
			strId)
		return
	}
	die(DeleteNote(id))
	http.Redirect(w, r, "/admin/upload", 301)
}

func stateEntry(key string, cb ChainBag) KV {
	return KV{key, cb.GetOr(key, "")}
}

func saveMainPageState(r *http.Request) error {
	fallback, err := LoadState()
	die(err)
	die(r.ParseForm())
	state := MultiBag(r.Form).BackedBy(fallback)

	toSave := []KV{
		stateEntry(scratchpadKey, state),
		stateEntry(draftnoteKey, state),
		stateEntry(taglineKey, state),
		stateEntry(recentSearchKey, state),
		stateEntry(appNameKey, state),
	}
	return SaveState(toSave)
}

func die(e error) {
	if e != nil {
		panic(e)
	}
}
