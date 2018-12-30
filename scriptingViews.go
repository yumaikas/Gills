package main

import (
	"fmt"
	"io"
	. "yumaikas/gills/templates"
)

func ScriptListView(w io.Writer, scripts []Script, state AppState) error {
	if len(scripts) == 0 {
		var template = BasePage(state.AppName(),
			Div(Atr, Str("Looks like you don't have any scripts yet."),
				Str("Why don't you "), A(Atr.Href("/admin/scripts/new/"), Str("make one?")),
			),
		)
		return RenderWithTargetAndTheme(w, "AQUA", template)
	}

	var template = BasePage(state.AppName(),
		Style(`div.script-info { padding: 30px; }`),
		Div(Atr,
			H2(Atr, A(Atr.Href("/admin/scripts/new"), Str("Create new script"))),
		),
		Div(Atr,
			H2(Atr, Str("Pages")),
			PageScripts(scripts),
			Hr(),
		),
		Div(Atr,
			H2(Atr, Str("Script library")),
			LibScripts(scripts),
			Hr(),
		),
		Div(Atr,
			H2(Atr, Str("Other scripts")),
			OtherScripts(scripts),
			Hr(),
		),
	)

	return RenderWithTargetAndTheme(w, "AQUA", template)

}
func renderScriptLinks(s Script, ctx Context) {
	Div(Atr.Class("script-info"),
		Div(Atr, H3(Atr, Str(s.Name))),
		Span(Atr, PageLink(s), Str(" "), EditLink(s)),
	)(ctx)
}

func PageLink(s Script) func(Context) {
	if s.IsPage() {
		return A(Atr.Href("/admin/pages/s/"+s.Name), Str("View Page"))
	}
	return func(ctx Context) {}
}

func EditLink(s Script) func(Context) {
	return A(Atr.Href("/admin/scripts/edit/"+s.Name), Str("Edit Script"))
}

func LibScripts(scripts []Script) func(Context) {
	return func(ctx Context) {
		for _, s := range scripts {
			if s.IsLibrary() {
				renderScriptLinks(s, ctx)
			}
		}
	}
}

func PageScripts(scripts []Script) func(Context) {
	return func(ctx Context) {
		for _, s := range scripts {
			if s.IsPage() {
				renderScriptLinks(s, ctx)
			}
		}
	}
}

func OtherScripts(scripts []Script) func(Context) {
	return func(ctx Context) {
		for _, s := range scripts {
			if !s.IsPage() && !s.IsLibrary() {
				renderScriptLinks(s, ctx)
			}
		}
	}

}

func LuaScriptEditView(w io.Writer, state AppState, code, name string) error {
	var template = BasePage(state.AppName(),
		Form(ExecuteActionForScriptName(name),
			Div(Atr,
				Label(Atr.For("script-name"),
					Str("Script Name"),
					Input(Atr.Type("text").Name("script-name").Id("script-name").Value(name)),
				),
				SubmitForScriptName(name),
				Input(Atr.Type("submit").Value("Execute Code!")),
				Input(Atr.Id("old-script-name").Name("old-script-name").Type("hidden").Value(name)),
			),
			TextArea(Atr.Name("code").Cols("80").Rows("40"), code),
		),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}

func ExecuteActionForScriptName(name string) AttributeChain {
	if len(name) > 0 {
		return Atr.Action("/admin/scripts/save-and-run/" + name).Method("POST")
	}
	return Atr.Action("/admin/scripts/new/run").Method("POST")
}

func SubmitForScriptName(name string) func(Context) {
	if len(name) > 0 {
		return Input(Atr.Type("submit").
			Value("Save Script").
			FormMethod("POST").
			FormAction(fmt.Sprint("/admin/scripts/edit/", name)))
	}
	return Input(Atr.Type("submit").
		Value("Create Script").
		FormMethod("POST").
		FormAction("/admin/scripts/new/"))
}

func LuaExecutionOnlyView(w io.Writer, state AppState, doScript func(Context)) error {
	return RenderWithTargetAndTheme(w, "AQUA", BasePage(state.AppName(), doScript))
}

func LuaExecutionResultsView(w io.Writer, state AppState, doScript func(Context), code, name string) error {
	var template = BasePage(state.AppName(),
		Form(Atr.Action("/admin/scripting/test").Method("POST"),
			Label(Atr.For("script-name"),
				Str("Script Name"),
				Input(Atr.Type("text").Name("script-name").Id("script-name").Value(name)),
			),
			Input(Atr.Type("submit").Value("Execute Code!")),
			TextArea(Atr.Name("code").Cols("80").Rows("40"), code),
		),
		Div(Atr.Id("code-results"),
			doScript,
		),
	)
	return RenderWithTargetAndTheme(w, "AQUA", template)
}
