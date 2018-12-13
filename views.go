package main

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"
	. "yumaikas/gills/templates"
)

var homeLayout = `
@media (min-aspect-ratio: 2/1) {
	.side-by-side {
		display: grid;
		grid-template-columns: 30px auto 10px auto 10px auto 30px;
		grid-template-areas: ". left . center . right .";
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
	.app-name {
	margin-left: 30px;
	}
}

@media (max-aspect-ratio: 2/1) {
	body {
		font-size: 26px;
		width: 100%;
	}

	input[type=submit] {
		margin-left: 20px;
		padding: 5px 10px;
		font-size: 24px;
	}
	input[type=text] {
		font-size: 24px;
	}
	
	input[type=submit].inline-form {
		margin-left: 5px;
		padding: 5px 10px;
	}
	
	input {
		font-size: 22px;
		margin-top: 20px;
		margin-bottom: 20px;
	}
	textarea {
		font-size: 22px;
		width: 95%;
	}
}
`

func HomeView(w io.Writer, state AppState, recentNotes []Note) error {
	var template = BasePage(state.AppName(),
		Style(homeLayout),
		H2(Atr.Class("app-name"), Str(state.AppName())),
		// TODO: Find a better place to put this, probably make a better header-section
		// A(Atr.Href("http://jessicaabel.com/ja/growing-gills/"), Str("Growing Gills"))),
		Form(Atr.Class("side-by-side").Action("/admin/search"),
			Div(Atr.Class("col-left"),
				Label(Atr.For(scratchpadKey),
					Div(Atr, Str("Scratchpad:")),
					TextArea(Atr.Name(scratchpadKey).Id("scratchpad").Cols("48").Rows("20"), state.ScratchPad()),
				),
			),
			Div(Atr.Class("col-center"),
				Label(Atr.For(draftnoteKey),
					Div(Atr, Str("What is on your mind?")),
					TextArea(Atr.Name(draftnoteKey).Id("draftnote").Cols("48").Rows("15"), state.DraftNote()),
				),
				Label(Atr.For(taglineKey),
					Div(Atr,
						Str("Tagline:"),
						Input(Atr.Name(taglineKey).Type("textbox").Id("tagline").Size("45").Value(state.Tagline())),
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
					Input(Atr.Name(recentSearchKey).Id("recentSearch").Size("45").Value(state.RecentSearchTerms())),
				),
				RecentNotes(recentNotes, 5),
			),
		),
		NoteDeletionForms(recentNotes, 5),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}

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

func NoteDetailsView(w io.Writer, state AppState, note Note) error {
	noteURL := fmt.Sprint("/admin/note/", note.Id)
	var template = BasePage(state.AppName(),
		Form(
			Atr.Action(noteURL).Method("POST").Class("note-container"),
			Str("On "+tfmt(note.Created)+" you said:"),
			TextArea(DimensionsOf(note).Name("note-content"), note.Content),
			Input(Atr.Type("Submit").Value("Save Changes")),
			Input(Atr.Type("Submit").Value("Delete Note").FormAction(noteURL).FormMethod("DELETE")),
		),
	)

	return RenderWithTargetAndTheme(w, "AQUA", template)
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
