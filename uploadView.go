package main

import (
	"io"
	. "yumaikas/gills/templates"
)

var uploadViewLayout = `
	input[type="file"] {
		display: none;
	}

	#submit-upload {
		display: none;
	}

	#q, #search-button, .custom-file-upload {
		font-size: 24px;
	}

	#q {
		width: 60%;
	}
	#search-button {
		width: 35%;
	}
	.custom-file-upload {
		display: block;
	    width: 95%;
		padding: 6px 12px;
		font-size: bigger;
		border: 2px solid grey;
		margin-bottom: 20px;
		text-align: center;
	}
`

var uploadScript = `
    function byId(sel) { return document.getElementById(sel) ; };
    byId("take-photo-text").innerHTML = "Take Photo";
	byId("upload").onchange = function() {
		byId("submit-upload").click();
	};
`

func UploadView(w io.Writer, state AppState, searchTerms string, recentUploadNotes []Note) error {
	var template = BasePage(state.AppName(),
		Style(uploadViewLayout),
		H2(Atr, Str("Upload images")),
		Form(Atr.Action("/admin/upload").Method("GET").EncType("multipart/form-data"),
			Label(Atr.For("upload").Class("custom-file-upload"),
				Span(Atr.Id("take-photo-text"), Str("Loading...")),
				Input(Atr.Id("upload").Type("file").Name("upload").Accept("image/*").Multiple()),
			),
			Button(Atr.Id("submit-upload").FormAction("/admin/upload").FormMethod("POST"), Str("Upload Image")),
			Div(Atr,
				Label(Atr.For("q"),
					Input(Atr.Type("text").Id("q").Name("q").Value(searchTerms)),
				),
				Input(Atr.Type("Submit").Id("search-button").Value("Search Notes")),
			),
			RecentNotes(recentUploadNotes, 5),
		),
		JS(uploadScript),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}
