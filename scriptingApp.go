package main

import (
	"github.com/go-chi/chi"
	"net/http"
	"strings"
)

// List the lua scripts
func ListLuaScripts(w http.ResponseWriter, r *http.Request) {
	var state, err = LoadState()
	die(err)
	scripts, err := ListScripts()
	die(err)
	die(ScriptListView(w, scripts, state))
}

func NewLuaScriptForm(w http.ResponseWriter, r *http.Request) {
	var state, err = LoadState()
	die(err)
	die(LuaScriptEditView(w, state, "", ""))
}

func RunDraftLuaScript(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.PostFormValue("code")
	state, err := LoadState()
	die(err)
	die(LuaExecutionResultsView(w, state, doLuaScript(code, r), code, ""))
}

func CreateLuaScript(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.PostFormValue("code")
	name := r.PostFormValue("script-name")
	state, err := LoadState()
	die(err)
	if len(name) <= 0 {
		InvalidScriptNameView(w, state.AppName(), name, "400: Script name cannot be empty")
		return
	}
	if strings.ContainsAny(name, "{}/?") {
		InvalidScriptNameView(w, state.AppName(), name, "400: script name cannot have any of {, }, / or ?")
		return
	}
	die(CreateScript(name, code))
	http.Redirect(w, r, "/admin/scripts/edit/"+name, 301)
}

func EditLuaScript(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	name := chi.URLParam(r, "script-name")
	state, err := LoadState()
	die(err)
	script, err := GetScriptByName(name)
	die(err)
	die(LuaScriptEditView(w, state, script.Content, script.Name))
}

func SaveLuaScript(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.PostFormValue("code")
	name := r.PostFormValue("script-name")
	oldName := r.PostFormValue("old-script-name")

	if name != oldName {
		die(RenameScript(oldName, name))
	}
	die(SaveScript(name, code))
	http.Redirect(w, r, "/admin/scripts/edit/"+name, 301)
}

func SaveAndRunLuaScript(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.PostFormValue("code")
	name := r.PostFormValue("script-name")
	oldName := r.PostFormValue("old-script-name")

	if name != oldName {
		die(RenameScript(oldName, name))
	}
	die(SaveScript(name, code))
	// Ensure that the script is a page-script
	state, err := LoadState()
	die(err)
	die(err)
	die(LuaExecutionOnlyView(w, state, doLuaScript(code, r)))
}

func RunLuaScript(w http.ResponseWriter, r *http.Request) {
	// Ensure that the script is a page-script
	name := chi.URLParam(r, "script-name")
	script, err := GetScriptByName(name)
	state, err := LoadState()
	die(err)
	die(err)
	die(LuaExecutionOnlyView(w, state, doLuaScript(script.Content, r)))
}

func RunLuaPostScript(w http.ResponseWriter, r *http.Request) {
	// Ensure that the script is a page-script
	die(r.ParseForm())
	name := chi.URLParam(r, "script-name")
	script, err := GetScriptByName(name)
	state, err := LoadState()
	die(err)
	die(err)
	die(LuaExecutionOnlyView(w, state, doLuaScript(script.Content, r)))
}
