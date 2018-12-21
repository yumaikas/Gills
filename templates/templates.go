package templates

import (
	"bytes"
	"github.com/russross/blackfriday"
	"html/template"
	"io"
	"strings"
)

// TODO: Come back to this if I find it works better "gopkg.in/russross/blackfriday.v2"
// TODO: integrate "github.com/microcosm-cc/bluemonday" into this code if I start trusting more than one user.

// This is a content string
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
	// If true, then this is an attribute without a value
	Void bool
}

type Context struct {
	indentCount int
	w           io.Writer
	themeName   string
}

func RenderWithTargetAndTheme(w io.Writer, themeName string, template func(Context)) (err error) {
	defer func() {
		val := recover()
		if errInner, ok := val.(error); ok {
			err = errInner
		}
	}()
	err = nil
	template(Context{0, w, themeName})
	return
}

func (ctx Context) startLine() {
	for i := 0; i < ctx.indentCount; i++ {
		_, err := ctx.w.Write([]byte("\t"))
		if err != nil {
			panic(err)
		}
	}
}
func (ctx Context) endLine() {
	_, err := ctx.w.Write([]byte("\n"))
	if err != nil {
		panic(err)
	}
}
func (ctx Context) write(content string) {
	_, err := ctx.w.Write([]byte(content))
	if err != nil {
		panic(err)
	}
}

func (ctx Context) writeLine(content string) {
	ctx.startLine()
	ctx.write(content)
	ctx.endLine()
}

func nothing(ctx Context) {}

func BasePage(title string, inner ...func(Context)) func(Context) {
	return Html(Title(Atr, Str(title)), Body(inner...))
}

func Html(inner ...func(Context)) func(Context) {
	return func(ctx Context) {
		ctx.writeLine("<!DOCTYPE html>")
		ctx.writeLine("<html>")
		ctx.writeLine("<meta charset=\"utf-8\">")
		ctx.writeTags(inner...)
		ctx.writeLine("</html>")
	}
}

func StrBr(content string) func(Context) {
	return func(ctx Context) {
		var buf bytes.Buffer
		template.HTMLEscape(&buf, []byte(content))
		ctx.startLine()
		ctx.write(strings.Replace(buf.String(), "\n", "<br/>", -1))
		ctx.endLine()
	}
}

func Markdown(content string) func(Context) {
	return func(ctx Context) {
		flags := 0 |
			blackfriday.HTML_USE_SMARTYPANTS |
			blackfriday.HTML_SMARTYPANTS_FRACTIONS

		extensions := blackfriday.EXTENSION_TABLES |
			blackfriday.EXTENSION_FENCED_CODE |
			blackfriday.EXTENSION_AUTOLINK |
			blackfriday.EXTENSION_STRIKETHROUGH |
			blackfriday.EXTENSION_HARD_LINE_BREAK |
			0

		renderer := blackfriday.HtmlRenderer(flags, "", "")

		output := blackfriday.Markdown([]byte(content), renderer, extensions)
		ctx.write(string(output))
	}
}

func Str(content string) func(Context) {
	return func(ctx Context) {
		ctx.startLine()
		template.HTMLEscape(ctx.w, []byte(content))
		ctx.endLine()
	}
}

func Body(inner ...func(ctx Context)) func(ctx Context) {
	return writeTag("body", Atr, append([]func(Context){baseStyle}, inner...)...)
}

func (ctx Context) indentMultiline(str string) {
	var lines = strings.Split(strings.Trim(str, " \r\n"), "\n")
	for _, line := range lines {
		ctx.writeLine(strings.Trim(line, " \r\n"))
	}
}

var baseStyle = writeTag("style", Atr, func(ctx Context) {
	if ctx.themeName == "AQUA" {
		ctx.indentMultiline(`
            body {
            	max-width: 1200px;
            	width: 80%;
            }
            body,input,textarea,button {
            	font-family: Iosevka, monospace;
            	background: #191e2a;
            	color: #21EF9F;
            }
            a { color: aqua; }
            a:visited { color: darkcyan; }
            .note-card {
            	border-top: solid 1px #21EF9F; 
	        	margin-top: 5px;
	        	padding-top: 5px;
            }
            img {max-width: 85%;}
            `)
	}
})

func (ctx Context) writeAttributes(attributes AttributeChain) {
	for _, attr := range attributes {
		ctx.write(" " + attr.Key)
		if attr.Void {
			continue
		}
		ctx.write("=\"")
		if attr.Trusted {
			ctx.write(attr.Value)
		} else {
			template.HTMLEscape(ctx.w, []byte(attr.Value))
		}
		ctx.write("\"")
	}
}

func writeVoidTag(tagname string, attributes AttributeChain) func(Context) {
	return func(ctx Context) {
		ctx.startLine()
		ctx.write("<" + tagname)
		ctx.writeAttributes(attributes)
		ctx.write(">")
		ctx.endLine()
	}
}

func (ctx Context) writeTags(inner ...func(ctx Context)) {
	innerCtx := Context{ctx.indentCount + 1, ctx.w, ctx.themeName}
	for _, fn := range inner {
		fn(innerCtx)
	}
}

func writeTag(tagname string, attributes AttributeChain, inner ...func(ctx Context)) func(ctx Context) {
	return func(ctx Context) {
		ctx.startLine()
		ctx.write("<")
		ctx.write(tagname)
		ctx.writeAttributes(attributes)
		ctx.write(">")
		ctx.endLine()
		ctx.writeTags(inner...)
		ctx.writeLine("</" + tagname + ">")
	}
}
