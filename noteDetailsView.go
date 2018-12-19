package main

import (
	"fmt"
	"io"
	. "yumaikas/gills/templates"
)

var NoteDetailsStyle = `
@media (min-aspect-ratio: 2/1) {
	textarea {
		max-width: 800px;
	}
}
@media (max-aspect-ratio: 2/1) {
	textarea {
		width: 95%;
	}
}
`

func NoteDetailsView(w io.Writer, state AppState, note Note) error {
	noteURL := fmt.Sprint("/admin/note/", note.Id)
	var template = BasePage(state.AppName(),
		Form(Atr.Action(noteURL).Method("POST").Class("note-container"),
			H2(Atr,
				A(Atr.Href("/admin/"), Str(state.AppName()+" Home"))),
			Div(Atr, Str("On "+tfmt(note.Created)+" you said:")),
			TextArea(DimensionsOf(note).Name("note-content"), note.Content),
			Div(Atr,
				Input(Atr.Type("Submit").Value("Save Changes")),
				Button(Atr.Type("Submit").FormAction(noteURL+"/delete").FormMethod("POST"),
					Str("Delete Note"),
				),
			),
		),
		Div(Atr, Str("Preview: ")),
		Markdown(note.Content),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}
