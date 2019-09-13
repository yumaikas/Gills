package main

import (
	"io"

	. "gills/templates"
)

var homeLayout = `
@media (min-aspect-ratio: 4/3) {
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

@media (max-aspect-ratio: 4/3) {
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
		H2(Atr.Class("app-name"), A(Atr.Href("/"), Str(state.AppName()))),
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
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}
