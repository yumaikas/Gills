package main

import (
	. "./templates"
	"fmt"
	"io"
	"time"
)

var homeLayout = `
.side-by-side {
	display: grid;
	grid-template-columns: 30px auto 10px 1fr 30px;
	grid-template-areas: ". left . right .";
}

.side-left {
	grid-area: left;
}

.side-right {
	grid-area: right;
}
`

func HomeView(w io.Writer, appState Bag, recentNotes []Note) error {
	appName := appState.getOr("appName", "Gills")
	scratchpadContent := appState.getOr("scratchpad", "")
	draftNote := appState.getOr("draftnote", "")
	tagline := appState.getOr("tagline", "")

	var template = BasePage(appName,
		Style(homeLayout),
		H2(Atr, Str(appName)),
		Div(Atr.Id("subheader-callout"),
			Str("The tagged journal for keeping up with things. Name inspired by "),
			A(Atr.Href("http://jessicaabel.com/ja/growing-gills/"), Str("Growing Gills"))),
		Form(Atr.Class("side-by-side"),
			Div(Atr.Class("side-left"),
				Label(Atr.For("scratchpad"),
					Div(Atr, Str("Scratchpad:")),
					TextArea(Atr.Name("scratchpad").Id("scratchpad"), scratchpadContent),
				),
			),
			Div(Atr.Class("side-right"),
				Input(Atr.Type("submit").Value("Save Note")),
				Label(Atr.For("draftnote"),
					Div(Atr, Str("What is on your mind?")),
					TextArea(Atr.Name("draftnote").Id("draftnote"), draftNote),
				),
				Label(Atr.For("tagline"),
					Div(Atr,
						Str("Tagline:"),
						Input(Atr.Name("tagline").Type("textbox").Id("tagline").Value(tagline)),
					),
				),
				RecentNotes(recentNotes, 5),
			),
		),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
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

func min(x, y int) int {
	if x > y {
		return y
	}
	return x
}

var NoteDetails = Div(Atr.Class("note-card"))

func RecentNotes(notes []Note, count int) func(Context) {
	return func(ctx Context) {
		numNotes := min(len(notes), count)
		for i := 0; i < numNotes; i++ {
			renderNote(notes[i], ctx)
		}
	}
}

func renderNote(n Note, ctx Context) {
	Div(Atr.Add("data-note-id", fmt.Sprint(n.Id)).Class("note-card"),
		Div(Atr,
			Str("Created on: "+tfmt(n.Created)),
			A(Atr.Href("#"), Str("Edit")),
			A(Atr.Href("#"), Str("Mark ForRemoval")),
		),
		Str(n.Content),
	)
}
