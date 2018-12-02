package main

import (
	. "./templates"
	"fmt"
	"io"
	"time"
)

var homeViewTemplate = Html(Body(Div_(Str("Home!"))))

func TempHomeView(w io.Writer) {
	homeViewTemplate(WithTargetAndTheme(w, "AQUA"))
}

func HomeView(w io.Writer, appState map[string]string, recentNotes []Note) error {
	Html(Body(Div_(Str("Home!"))))(WithTargetAndTheme(w, "AQUA"))
	return nil
}

func SearchView(w io.Writer, appState map[string]string, recentNotes []Note) error {
	return nil
}

func NoteDetailsView(w io.Writer, note Note) error {
	return nil
}

func tfmt(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func noteRender(n Note) func(Context) {
	return Div(
		Kv{"data-note-id": fmt.Sprint(n.Id), "class": "note-container"},
		Tags(
			Div_(Tags(
				Str("Created on: "+tfmt(n.Created)),
				A(Kv{"href": "#"}, Str("Edit")),
				A(Kv{"href": "#"}, Str("Mark ForRemoval")),
				Str(n.Content)),
			),
		),
	)
}
