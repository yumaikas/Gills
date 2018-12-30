package main

import (
	"io"
	. "yumaikas/gills/templates"
)

func ScriptListView(w io.Writer, scripts []Script, state AppState) error {
	if len(scripts) == 0 {
		var template = BasePage(state.AppName(),
			Div(Atr, Str("Looks like you don't have any scripts yet."),
				Str("Why don't you "), A(Atr.Href("/admin/scriptes/new"), Str("make one?")),
			),
		)
		return RenderWithTargetAndTheme(w, "AQUA", template)
	}

	var template = BasePage(state.AppName())

	return RenderWithTargetAndTheme(w, "AQUA", template)

}

func LuaNewScriptView(w io.Writer, state AppState) error {
	var template = BasePage(state.AppName(),
		Form(Atr.Action("/admin/scripting/test").Method("POST"),
			Label(Atr.For("script-name"),
				Input(Atr.Type("text").Name("script-name").Id("script-name")),
			),
			Input(Atr.Type("submit").Value("Save Script").FormMethod("POST").FormAction("/admin/scripts/new/run")),
			TextArea(Atr.Name("code").Cols("80").Rows("40"), ""),
			Input(Atr.Type("submit").Value("Execute Code!")),
		),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}

func LuaExecutionResultsView(w io.Writer, state AppState, doScript func(Context), code string) error {
	var template = BasePage(state.AppName(),
		Form(Atr.Action("/admin/scripting/test").Method("POST"),
			TextArea(Atr.Name("code").Cols("80").Rows("40"), code),
			Input(Atr.Type("submit").Value("Execute Code!")),
		),
		Div(Atr.Id("code-results"),
			doScript,
		),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}
