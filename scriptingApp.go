package main

import (
	"net/http"
)

// List the lua scripts
func ListLuaScripts(w http.ResponseWriter, r *http.Request) {
	var state, err = LoadState()
	die(err)
	scripts, err := ListScripts()
	die(err)

	die(ScriptListView(w, scripts, state))
}

func RunNewLuaScript(w http.ResponseWriter, r *http.Request) {
}

func EditLuaScript(w http.ResponseWriter, r *http.Request) {
}

func CreateLuaScript(w http.ResponseWriter, r *http.Request) {
}

func SaveLuaScript(w http.ResponseWriter, r *http.Request) {
}

func RunLuaScript(w http.ResponseWriter, r *http.Request) {
}

func RunLuaPostScript(w http.ResponseWriter, r *http.Request) {
}

// Execute code
func RunLuaTestForm(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	code := r.PostFormValue("code")
	state, err := LoadState()
	die(err)
	die(LuaExecutionResultsView(w, state, doLuaScript(code), code))
}

// Show the new lua form
func NewLuaScriptForm(w http.ResponseWriter, r *http.Request) {
	var state, err = LoadState()
	die(err)
	die(LuaNewScriptView(w, state))
}
