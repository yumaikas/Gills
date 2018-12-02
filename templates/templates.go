package templates

import (
	"html/template"
	"io"
	"strings"
)

// This is a content string,
type Content struct {
	Raw string
	// If false, escape the string, if true, emit it raw
	Trusted bool
}

type Attribute struct {
	Key   string
	Value string
	// If false, escape the string, if true, emit it raw
	Trusted bool
}

type Context struct {
	indentCount int
	w           io.Writer
	themeName   string
}

func WithTargetAndTheme(w io.Writer, themeName string) Context {
	return Context{0, w, themeName}
}

type Kv map[string]string

// Half of me is tempted to create a tag-writing Context....

func (ctx Context) startLine() {
	for i := 0; i < ctx.indentCount; i++ {
		ctx.w.Write([]byte("\t"))
	}
}
func (ctx Context) endLine() {
	ctx.w.Write([]byte("\n"))
}
func (ctx Context) write(content string) {
	ctx.w.Write([]byte(content))
}

func (ctx Context) writeLine(content string) {
	ctx.startLine()
	ctx.write(content)
	ctx.endLine()
}

func nothing(ctx Context) {}

var NoAttr = make(Kv)

func BasePage(inner func(Context)) func(Context) {
	return Html(Body(inner))
}

func Html(inner func(Context)) func(Context) {
	return func(ctx Context) {
		ctx.writeLine("<!DOCTYPE html>")
		ctx.writeLine("<html>")
		inner(Context{ctx.indentCount + 1, ctx.w, ctx.themeName})
		ctx.writeLine("</html>")
	}
}

func Str(content string) func(Context) {
	return func(ctx Context) {
		template.HTMLEscape(ctx.w, []byte(content))
	}
}

func Body(inner func(ctx Context)) func(ctx Context) {
	return writeTag("body", NoAttr, Tags(writeStyle, inner))
}

func Tags(inner ...func(ctx Context)) func(ctx Context) {
	return func(ctx Context) {
		innerCtx := Context{ctx.indentCount + 1, ctx.w, ctx.themeName}
		for _, fn := range inner {
			fn(innerCtx)
		}
	}
}

func (ctx Context) indentMultiline(str string) func(Context) {
	return func(ctx Context) {
		var lines = strings.Split(str, "\n")
		for _, line := range lines {
			strings.Trim(line, " \t\r\n")
			ctx.writeLine(line)
		}
	}
}

func writeStyle(ctx Context) {
	if ctx.themeName == "AQUA" {
		innerCtx := Context{ctx.indentCount + 2, ctx.w, ctx.themeName}
		styleMultiline := innerCtx.indentMultiline(`
			body {
			    max-width: 800px;
			    width: 80%;
			}
			body,input,textarea {
			       font-family: Iosevka, monospace;
			       background: #191e2a;
			       color: #21EF9F;
			}
			a { color: aqua; }
			a:visited { color: darkcyan; }`)
		writeTag("style", NoAttr, styleMultiline)(ctx)
	}
}

func writeVoidTag(tagname string, attributes Kv) func(Context) {
	return func(ctx Context) {
		ctx.startLine()
		ctx.write("<" + tagname)
		for key, val := range attributes {
			ctx.write(" " + key)
			ctx.write("=\"")
			ctx.write(val)
			ctx.write("\"")
		}
		ctx.write(">")
		ctx.endLine()
	}
}

func writeTag(tagname string, attributes Kv, inner func(ctx Context)) func(ctx Context) {
	return func(ctx Context) {
		ctx.startLine()
		ctx.write("<")
		ctx.write(tagname)
		for key, val := range attributes {
			ctx.write(" " + key)
			ctx.write("=\"")
			ctx.write(val)
			ctx.write("\"")
		}
		ctx.write(">")
		ctx.endLine()

		inner(Context{ctx.indentCount + 1, ctx.w, ctx.themeName})

		ctx.writeLine("</" + tagname + ">")
	}
}
