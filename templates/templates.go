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

type context struct {
	indentCount int
	w           io.Writer
	themeName   string
}

func WithTargetAndTheme(w io.Writer, themeName string) context {
	return context{0, w, themeName}
}

type Kv Attribute

// Half of me is tempted to create a tag-writing context....

func (ctx context) startLine() {
	for i := 0; i < ctx.indentCount; i++ {
		ctx.w.Write([]byte("\t"))
	}
}
func (ctx context) endLine() {
	ctx.w.Write([]byte("\n"))
}
func (ctx context) write(content string) {
	ctx.w.Write([]byte(content))
}

func (ctx context) writeLine(content string) {
	ctx.startLine()
	ctx.write(content)
	ctx.endLine()
}

func nothing(ctx context) {}

var NoAttr = []Attribute{}

func Attrs(pairs ...Attribute) []Attribute {
	return pairs
}

func Html(inner func(context)) func(context) {
	return func(ctx context) {
		ctx.writeLine("<!DOCTYPE html>")
		ctx.writeLine("<html>")
		inner(context{ctx.indentCount + 1, ctx.w, ctx.themeName})
		ctx.writeLine("</html>")
	}
}

func Str(content string) func(context) {
	return func(ctx context) {
		template.HTMLEscape(ctx.w, []byte(content))
	}
}

func Body(inner func(ctx context)) func(ctx context) {
	return writeTag("body", NoAttr, Tags(writeStyle, inner))
}

func Tags(inner ...func(ctx context)) func(ctx context) {
	return func(ctx context) {
		innerCtx := context{ctx.indentCount + 1, ctx.w, ctx.themeName}
		for _, fn := range inner {
			fn(innerCtx)
		}
	}
}

func (ctx context) indentMultiline(str string) {
	var lines = strings.Split(str, "\n")
	for _, line := range lines {
		strings.Trim(line, " \t\r\n")
		ctx.writeLine(line)
	}
}

func writeStyle(ctx context) {
	if ctx.themeName == "AQUA" {
		innerCtx = context{ctx.indentCount + 2, ctx.w, ctx.themeName}
		writeTag("style", NoAttr, Str(innerCtx.indentMultiline(
			`
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
a:visited { color: darkcyan; }`)))(ctx)
	}
}

func writeVoidTag(tagname string, attributes []Attribute) func(context) {
	return func(ctx context) {
		ctx.startLine()
		ctx.write("<" + tagname)
		for _, attr := range attributes {
			ctx.write(" " + attr.Key)
			ctx.write("=\"")
			if attr.Trusted {
				ctx.write(attr.Value)
			} else {
				template.HTMLEscape(ctx.w, []byte(attr.Value))
			}
			ctx.write("\"")
		}
		ctx.write(">")
		ctx.endLine()
	}
}

func writeTag(tagname string, attributes []Attribute, inner func(ctx context)) func(ctx context) {
	return func(ctx context) {
		ctx.startLine()
		ctx.write("<")
		ctx.write(tagname)
		for _, attr := range attributes {
			ctx.write(" " + attr.Key)
			ctx.write("=\"")
			if attr.Trusted {
				ctx.write(attr.Value)
			} else {
				template.HTMLEscape(ctx.w, []byte(attr.Value))
			}
			ctx.write("\"")
		}
		ctx.write(">")
		ctx.endLine()

		inner(context{ctx.indentCount + 1, ctx.w})

		ctx.writeLine("</" + tagname + ">")
	}
}
