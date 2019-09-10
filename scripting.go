package main

import (
	lua "github.com/Shopify/go-lua"
	"github.com/Shopify/goluago"
	"github.com/Shopify/goluago/util"
	"yumaikas/gills/templates"
	// "github.com/go-chi/chi"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

func doLuaScript(code string, r *http.Request) func(ctx templates.Context) {
	l := lua.NewState()
	lua.OpenLibraries(l)
	goluago.Open(l)
	die(CleanOSFunctions(l))
	RegisterDbFunctions(l)
	RegisterRequestArgFunctions(l, r)

	return func(ctx templates.Context) {
		buf := &bytes.Buffer{}
		// Will panic if things fail to load
		templates.GetLuaRenderingContext(ctx)(l)
		l.Register("echo", LuaDoPrint(buf))
		l.Register("echoln", LuaDoPrintln(buf))
		l.Register("clear_buffer", LuaClear(buf))
		l.Register("flush_markdown", LuaFlush(buf, ctx, templates.Markdown))
		l.Register("flush_plain", LuaFlush(buf, ctx, templates.StrBr))
		l.Register("flush_raw", LuaFlush(buf, ctx, templates.RawStr))
		l.Register("script", LuaEmitJsScript(ctx))
		l.Register("style", LuaEmitCSS(ctx))
		l.Register("require_note", LuaRequireNoteScript)
		// Emit a div with an ID as a JS mount-point for things like Vue.JS
		l.Register("app_div", LuaEmitDiv(ctx))
		// l.Register("write_tag", LuaWriteTag(ctx))
		// l.Register("write_void_tag", LuaWriteTag(ctx))
		// TODO: Build in
		err := lua.LoadString(l, code)
		msg, _ := l.ToString(l.Top())

		// fmt.Println(l)
		if err != nil {
			templates.StrBr("Partial output: " + msg)(ctx)
			templates.StrBr(buf.String())(ctx)
			templates.StrBr("Error:")(ctx)
			templates.StrBr(err.Error())(ctx)
			return
		}
		err = l.ProtectedCall(0, lua.MultipleReturns, 0)
		if err != nil {
			templates.StrBr("Partial output: " + msg)(ctx)
			templates.StrBr(buf.String())(ctx)
			templates.StrBr("Error:")(ctx)
			templates.StrBr(err.Error())(ctx)
			return
		}
		// String up whatever leftover stuff wasn't flushed
		templates.StrBr(buf.String())(ctx)
	}
}

func LuaRequireNoteScript(l *lua.State) int {
	if l.Top() < 1 {
		l.PushString("require_note called without script name!")
		l.Error()
	}
	scriptName, ok := l.ToString(1)
	if !ok {
		l.PushString("require_note called with non-string script name!")
		l.Error()
	}
	note, err := GetScriptByName(scriptName)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
	}
	err = lua.LoadString(l, note.Content)
	if err != nil {
		msg, _ := l.ToString(l.Top())
		l.PushString(msg)
		l.Error()
	}
	l.Call(0, lua.MultipleReturns)
	return 0

}

func CleanOSFunctions(l *lua.State) error {
	return lua.DoString(l, `
		io = nil
		os.execute = nil 
		os.getenv = nil 
		os.remove = nil 
		os.rename = nil
		os.tmpname = nil
		os.exit = nil
		`)
}

func RegisterDbFunctions(l *lua.State) {
	l.Register("search_notes", LuaDoSearch)
	l.Register("note_for_id", LuaNoteForId)
}

func RegisterRequestArgFunctions(l *lua.State, r *http.Request) {
	l.Register("url_query", LuaUrlQuery(r))
	l.Register("form_value", LuaFormValue(r))

}
func LuaFormValue(r *http.Request) func(*lua.State) int {
	return func(l *lua.State) int {
		argName, ok := l.ToString(1)
		if !ok {
			l.PushString("Expected a string for the query name argument!")
			l.Error()
		}
		l.PushString(r.FormValue(argName))
		return 1
	}
}

func LuaUrlQuery(r *http.Request) func(*lua.State) int {
	return func(l *lua.State) int {
		argName, ok := l.ToString(1)
		if !ok {
			l.PushString("Expected a string for the query name argument!")
			l.Error()
		}
		l.PushString(r.URL.Query().Get(argName))
		return 1
	}
}

func LuaClear(buf *bytes.Buffer) func(*lua.State) int {
	return func(l *lua.State) int {
		buf.Reset()
		return 0
	}
}

func LuaFlush(buf *bytes.Buffer, ctx templates.Context, flushUsing func(string) func(templates.Context)) func(*lua.State) int {
	return func(l *lua.State) int {
		flushUsing(buf.String())(ctx)
		buf.Reset()
		return 0
	}
}

func LuaEmitDiv(ctx templates.Context) func(*lua.State) int {
	return func(l *lua.State) int {
		if l.Top() != 2 {
			l.PushString("\"app_div\" requires two arguments, id and loading message")
			l.Error()
		}
		id, ok1 := l.ToString(1)
		loadingText, ok2 := l.ToString(2)
		if !(ok1 && ok2) {
			l.PushString("Either id or loadingText isn't a string!")
			l.Error()
		}
		templates.Div(templates.Atr.Id(id), templates.Str(loadingText))(ctx)
		return 0
	}

}

func LuaEmitCSS(ctx templates.Context) func(*lua.State) int {
	return func(l *lua.State) int {
		if l.Top() != 2 {
			l.PushString("\"style\" requires two arguments")
			l.Error()
		}
		mode, _ := l.ToString(1)
		script, ok := l.ToString(2)
		if mode == "text" && ok {
			templates.Style(script)(ctx)
		} else if mode == "link" && ok {
			templates.StyleLink(script)(ctx)
		} else {
			l.PushString("First argument for \"style\" must be \"text\" or \"link\", and the second argument must be a string")
			l.Error()
		}
		return 0
	}
}

func LuaEmitJsScript(ctx templates.Context) func(*lua.State) int {
	return func(l *lua.State) int {
		if l.Top() != 2 {
			l.PushString("script requires two arguments")
			l.Error()
		}

		mode, _ := l.ToString(1)
		script, ok := l.ToString(2)
		if mode == "text" && ok {
			templates.JS(script)(ctx)
		} else if mode == "link" && ok {
			templates.JSLink(script)(ctx)
		} else {
			l.PushString("First argument for scipt must be \"text\" or \"link\", and the second argument must be a string")
			l.Error()
		}
		return 0
	}
}

func LuaDoPrintln(w io.Writer) func(*lua.State) int {
	return func(l *lua.State) int {
		numArgs := l.Top()
		for i := 1; i <= numArgs; i++ {
			str, ok := l.ToString(i)
			if !ok {
				l.PushString(fmt.Sprint("Cannot convert argument at position", i, "to a lua string for printing!"))
				l.Error()
			}
			w.Write([]byte(str))
		}
		w.Write([]byte("\n"))
		return 0
	}
}
func LuaDoPrint(w io.Writer) func(*lua.State) int {
	return func(l *lua.State) int {
		numArgs := l.Top()
		for i := 1; i <= numArgs; i++ {
			str, ok := l.ToString(i)
			if !ok {
				l.PushString(fmt.Sprint("Cannot convert argument at position", i, "to a lua string for printing!"))
				l.Error()
			}
			w.Write([]byte(str))
		}
		return 0
	}
}

func LuaDoSearch(l *lua.State) int {
	numArgs := l.Top()
	fmt.Println("Number of arguments:", numArgs)
	searchTerms, ok := l.ToString(1)
	if !ok {
		l.PushString("Cannot search on a non-string term!")
		l.Error()
	}
	notes, err := SearchNotes(searchTerms)
	if err != nil {
		l.PushString(err.Error())
		l.Error()
	}
	util.DeepPush(l, MapFromNotes(notes))
	return 1
}

func LuaNoteForId(l *lua.State) int {
	id := lua.CheckInteger(l, 1)
	note, err := GetNoteBy(int64(id))
	if err != nil {
		l.PushString("An error happened fetching a note:"+ err.Error())
		l.Error()
	}
	util.DeepPush(l, MapFromNote(note))
	return 1

}

func MapFromNote(n Note) map[string]interface{} {
	return map[string]interface{}{
		"id":      n.Id,
		"content": n.Content,
		"created": MapFromGoTime(n.Created.Local()),
		"updated": MapFromGoTime(n.Updated.Local()),
	}
}

func MapFromNotes(notes []Note) []map[string]interface{} {
	var mappedNotes = make([]map[string]interface{}, len(notes))
	for idx, n := range notes {
		mappedNotes[idx] = MapFromNote(n)
	}

	return mappedNotes
}

func MapFromGoTime(t time.Time) map[string]interface{} {
	return map[string]interface{}{
		"year":       t.Year(),
		"month":      int(t.Month()),
		"day":        t.Day(),
		"hour":       t.Hour(),
		"minute":     t.Minute(),
		"second":     t.Second(),
		"nanosecond": t.Nanosecond(),
		"weekday":    int(t.Weekday()) + 1,
		"unix":       t.Unix(),
	}
}
