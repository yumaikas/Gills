package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	. "yumaikas/gills/templates"
)

func SearchView(w io.Writer, state AppState, searchedNotes []Note) error {
	var template = BasePage(state.AppName(),
		H2(Atr, Str(state.AppName()+" Search")),
		Div(Atr,
			Str("Standalone search isn't done yet, check out "),
			A(Atr.Href("/admin/"), Str("the main page")),
			Str(" in the mean time"),
		),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}

func InvalidIdView(w io.Writer, appName, message, invalidID string) error {
	var template = BasePage(appName,
		H2(Atr, Str("400: You sent me a goofy request")),
		Str(fmt.Sprint(message, invalidID)))

	return RenderWithTargetAndTheme(w, "AQUA", template)
}

func renderNoteDeleteForm(note Note, ctx Context) {
	var template = Form(
		Atr.Id(fmt.Sprint("note-delete-", note.Id)).Method("DELETE").Action(fmt.Sprint("/admin/note/", note.Id)))
	template(ctx)
}
func NoteDeletionForms(notes []Note, count int) func(Context) {
	return func(ctx Context) {
		numNotes := min(len(notes), count)
		for i := 0; i < numNotes; i++ {
			renderNoteDeleteForm(notes[i], ctx)
		}
	}
}

func RecentNotes(notes []Note, count int) func(Context) {
	return func(ctx Context) {
		numNotes := min(len(notes), count)
		for i := 0; i < numNotes; i++ {
			renderNote(notes[i], ctx)
		}
	}
}

func DimensionsOf(n Note) AttributeChain {
	lines := strings.Split(n.Content, "\n")
	numLines := strconv.Itoa(len(lines))
	maxLineLen := 0
	for _, l := range lines {
		maxLineLen = max(maxLineLen, len(strings.Trim(l, "\t \r\n")))
	}
	return Atr.Cols(strconv.Itoa(maxLineLen)).Rows(numLines)
}

func renderNote(n Note, ctx Context) {
	noteShowURL := fmt.Sprint("/admin/note/", n.Id, "/show")
	noteActionURL := fmt.Sprint("/admin/note/", n.Id)
	Div(Atr.Add("data-note-id", fmt.Sprint(n.Id)).Class("note-card"),
		Div(Atr,
			Str(tfmt(n.Created)),
			// TODO: these buttons need to be hooked up better, or I need to just break down and
			// write a JS frontend for this.
			Input(Atr.Type("submit").Value("Edit").FormAction(noteShowURL).FormMethod("POST")),
			Input(Atr.Type("submit").Value("Delete").FormAction(noteActionURL).FormMethod("DELETE")),
		),
		StrBr(n.Content),
	)(ctx)
}

func tfmt(t time.Time) string {
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d",
		t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute())
}

func max(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}
