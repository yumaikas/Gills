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
	grid-template-columns: 30px auto 10px auto 10px auto 30px;
	grid-template-areas: ". left . center . right .";
}

.app-name {
	margin-left: 30px;
}

.col-left {
	grid-area: left;
}

.col-center {
	grid-area: center;
}

.col-right {
	grid-area: right;
}

textarea, .note-card {
	font-size: 9pt;
}
`

const scratchpadKey = "scratchpad"
const draftnoteKey = "draftnote"
const taglineKey = "tagline"
const recentSearchKey = "recentSearchTerms"

func HomeView(w io.Writer, appState Bag, recentNotes []Note) error {
	appName := appState.getOr("appName", "Gills")
	scratchpadContent := appState.getOr(scratchpadKey, "")
	draftNote := appState.getOr(draftnoteKey, "")
	tagline := appState.getOr(taglineKey, "")
	recentSearchTerms := appState.getOr(recentSearchKey, "")

	var template = BasePage(appName,
		Style(homeLayout),
		H2(Atr.Class("app-name"), Str(appName)),
		// TODO: Find a better place to put this, probably make a better header-section
		// A(Atr.Href("http://jessicaabel.com/ja/growing-gills/"), Str("Growing Gills"))),
		Form(Atr.Class("side-by-side").Action("/admin/search"),
			Div(Atr.Class("col-left"),
				Label(Atr.For(scratchpadKey),
					Div(Atr, Str("Scratchpad:")),
					TextArea(Atr.Name(scratchpadKey).Id("scratchpad").Cols("48").Rows("20"), scratchpadContent),
				),
			),
			Div(Atr.Class("col-center"),
				Label(Atr.For(draftnoteKey),
					Div(Atr, Str("What is on your mind?")),
					TextArea(Atr.Name(draftnoteKey).Id("draftnote").Cols("48").Rows("15"), draftNote),
				),
				Label(Atr.For(taglineKey),
					Div(Atr,
						Str("Tagline:"),
						Input(Atr.Name(taglineKey).Type("textbox").Id("tagline").Size("45").Value(tagline)),
						Input(Atr.Type("submit").
							Class("inline-form").
							Value("Save Note").
							FormMethod("POST").FormAction("/admin/create")),
					),
				),
			),
			Div(Atr.Class("col-right"),
				Label(Atr.For(recentSearchKey),
					Input(Atr.Type("submit").FormAction("/admin/search").Class("inline-form").FormMethod("POST").Value("Search Recent")),
					Input(Atr.Name(recentSearchKey).Id("recentSearch").Size("45").Value(recentSearchTerms)),
				),
				RecentNotes(recentNotes, 5),
			),
		),
		NoteDeletionForms(recentNotes, 5),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}

func SearchView(w io.Writer, appState Bag, searchedNotes []Note) error {
	appName := appState.getOr("appName", "Gills")
	var template = BasePage(appName,
		H2(Atr, Str(appName+" Search")),
		Div(Atr,
			Str("Standalone search isn't done yet, check out "),
			A(Atr.Href("/admin/"), Str("the main page")),
			Str(" in the mean time"),
		),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}

func InvalidSaveIdView(w io.Writer, appState Bag, invalidID string) error {
	appName := appState.getOr("appName", "Gills")
	var template = BasePage(appName,
		H2(Atr, Str("400: You sent me a goofy request")),
		Str(fmt.Sprint("Cannot save note based on invalid id ", invalidID)))

	return RenderWithTargetAndTheme(w, "AQUA", template)
}

func InvalidDeleteIdView(w io.Writer, appState Bag, invalidID string) error {
	appName := appState.getOr("appName", "Gills")
	var template = BasePage(appName,
		H2(Atr, Str("400: You sent me a goofy request")),
		Str(fmt.Sprint("Cannot delete note based on invalid id ", invalidID)))

	return RenderWithTargetAndTheme(w, "AQUA", template)
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

func renderNote(n Note, ctx Context) {
	Div(Atr.Add("data-note-id", fmt.Sprint(n.Id)).Class("note-card"),
		Div(Atr,
			Str(tfmt(n.Created)),
			// TODO: these buttons need to be hooked up better, or I need to just break down and
			// write a JS frontend for this.
			Input(Atr.Type("submit").Value("Edit")),
			Input(Atr.Form(fmt.Sprint("note-delete-", n.Id)).Type("submit").Value("Delete")),
		),
		StrBr(n.Content),
	)(ctx)
}
