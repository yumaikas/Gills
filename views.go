package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	. "yumaikas/gills/templates"
)

func SearchView(w io.Writer, appName, searchTerms string, searchedNotes []Note) error {
	var template = BasePage(appName,
		H2(Atr, Str(appName+" Search")),
		Div(Atr,
			A(Atr.Href("/admin/upload"), Str("Upload Photos")),
		),
		Div(Atr,
			A(Atr.Href("/admin/"), Str("View Other Notes")),
		),
		Form(Atr.Action("/").Method("GET"),
			Label(Atr.For("q"),
				Input(Atr.Type("text").Name("q").Value(searchTerms)),
			),
			Input(Atr.Type("Submit").Value("Search Notes")),
			RecentNotes(searchedNotes, len(searchedNotes)),
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
	Div(Atr.Add("data-note-id", fmt.Sprint(n.Id)).Class("note-card"),
		Div(Atr,
			Str(tfmt(n.Created)),
			Input(Atr.Type("submit").Value("Edit").FormAction(noteShowURL).FormMethod("POST")),
		),
		Markdown(n.Content),
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
