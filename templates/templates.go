package templates

import (
	"bytes"
	lua "github.com/Shopify/go-lua"
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

func GetLuaRenderingContext(derivedFrom Context) func(l *lua.State) int {
	return func(l *lua.State) int {
		err := lua.LoadString(l, `
			local strings = require("goluago/strings")
			local indentCount, writeFunc, escapeHTML, mdownFunc, themeName = ...

			local ctx = {}

			function ctx.start_line() 
				for i=1,indentCount,1 do
					ctx.write("\t")
				end
			end

			function ctx.end_line()
				writeFunc("\n")
			end

			function ctx.write(str) 
			    writeFunc(str)
			end

			function ctx.write_line(str) 
			    writeFunc(str)
			end

			function ctx.write_tags(inner) 
				indentCount = indentCount + 1
				for i=1,#inner do
				    inner[i]()
				end
				indentCount = indentCount - 1
			end

			function ctx.write_attributes(attributes) 
				for i=1,#attributes do 
					local attr = attributes[i]
				    ctx.write(" "..attr.Key..'="')
				    if attr.Trusted then
				    	ctx.write(attr.Value)
				    else
				    	ctx.write(escapeHTML(attr.Value))
				    end
				    ctx.write('"')
				end
			end

			function write_void_tag(tagname, attr)
				return function() 
					ctx.start_line()
					ctx.write("<"..tagname)
					ctx.write_attributes(attr)
					ctx.write(">")
					ctx.end_line()
				end
			end

			function write_tag(tagname, attr, ...) 
				local inner = {...}
				return function()
					ctx.start_line()
					ctx.write("<"..tagname)
					ctx.write_attributes(attr)
					ctx.write(">")
					ctx.end_line()
					ctx.write_tags(inner)
					ctx.write_line("</"..tagname..">")
				end
			end

			local AttributeMT = {
			   __index = function(t, k)
			        return function(value)
			        	table.insert(t, {Key = string.lower(k), Value = value, Trusted = false})
			        	return t
				    end
			   end
			}

			function Atr() 
				local atrChain = {}
				atrChain.AddUnsafe = function(key, value) 
				    table.insert(atrChain, {Key = string.lower(key), Value = value, Trusted = true})
				end
				setmetatable(atrChain, AttributeMT)
				return atrChain
			end

			function InlineStr(content)
				return function() 
				    ctx.write(escapeHTML(content))
				end
			end

			function Str(content)
				return function() 
					ctx.start_line()
				    ctx.write(escapeHTML(content))
				    ctx.end_line()
				end
			end

			function StrBr(content)
				return function()
				    local lines = strings.split(strings.trim(content, " \r\n"), "\n")
				    for i=1,#lines do
						ctx.write(escapeHTML(lines[i].." "))
				    end
				end
			end

			function Markdown(content)
			    return function() 
			       ctx.write(mdownFunc(content))
			    end
			end 


			local void_tags = {"br", "hr", "input"}

			local normal_tags = {
				"div", "span", "h1", "h2", "h3", "h4", 
				"title", "a", "table", "td", "form", "label", 
				"button", 
			}

			-- Add the above tags as functions to _ENV. 
			-- This is mostly intended for pages that need to use 
			-- the HTML renderer
			for i=1,#void_tags do
				local t = void_tags[i]
				local fn_name = string.upper(string.sub(t, 1,1))..string.sub(t, 2)
				_ENV[fn_name] = function(attributes, ...) 
					return write_void_tag(t, attributes, ...)
				end
			end 
			for i=1,#normal_tags do
				local t = normal_tags[i]
				local fn_name = string.upper(string.sub(t, 1,1))..string.sub(t, 2)
				_ENV[fn_name] = function(attributes, ...) 
					return write_tag(t, attributes, ...) 
				end
			end
		`)
		if err != nil {
			str, _ := l.ToString(l.Top())
			panic(str)
		}
		l.PushInteger(derivedFrom.indentCount)
		l.PushGoFunction(func(l *lua.State) int {
			str, ok := l.ToString(1)
			if !ok {
				l.PushString("Cannot convert argument for ctx.write to string!")
				l.Error()
			}
			_, err := derivedFrom.w.Write([]byte(str))
			if err != nil {
				l.PushString("Error while writing output: " + err.Error())
				l.Error()
			}
			return 0
		})
		l.PushGoFunction(func(l *lua.State) int {
			str, ok := l.ToString(1)
			if !ok {
				l.PushString("Cannot convert argument for escapeHTML to string!")
				l.Error()
			}
			buf := &bytes.Buffer{}
			template.HTMLEscape(buf, []byte(str))
			l.PushString(buf.String())
			return 1
		})
		l.PushGoFunction(LuaMarkdown)
		l.PushString(derivedFrom.themeName)
		err = l.ProtectedCall(5, 0, 0)
		if err != nil {
			l.PushString(err.Error())
			l.Error()
		}
		return 0
	}
}

func RenderWithTargetAndTheme(w io.Writer, themeName string, template func(Context)) (err error) {
	template(Context{0, w, themeName})
	return nil
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

func BasePage(title string, inner ...func(Context)) func(Context) {
	return Html(Title(Atr, Str(title)), Body(inner...))
}

func Html(inner ...func(Context)) func(Context) {
	return func(ctx Context) {
		ctx.writeLine("<!DOCTYPE html>")
		ctx.writeLine("<html>")
		ctx.writeLine("<meta charset=\"utf-8\">")
		ctx.WriteTags(inner...)
		ctx.writeLine("</html>")
	}
}

func RawStr(content string) func (Context) {
	return func (ctx Context) {
		ctx.write(content)
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

var mdown_flags = 0 |
	blackfriday.HTML_USE_SMARTYPANTS |
	blackfriday.HTML_SMARTYPANTS_FRACTIONS |
	blackfriday.HTML_FOOTNOTE_RETURN_LINKS

var mdown_extensions = blackfriday.EXTENSION_TABLES |
	blackfriday.EXTENSION_FOOTNOTES |
	blackfriday.EXTENSION_FENCED_CODE |
	blackfriday.EXTENSION_AUTOLINK |
	blackfriday.EXTENSION_STRIKETHROUGH |
	blackfriday.EXTENSION_HARD_LINE_BREAK |
	0

func LuaMarkdown(l *lua.State) int {
	content, ok := l.ToString(1)
	if !ok {
		l.PushString("Cannot render markdown from non-string argument")
		l.Error()
	}
	renderer := blackfriday.HtmlRenderer(mdown_flags, "", "")
	output := blackfriday.Markdown([]byte(content), renderer, mdown_extensions)
	l.PushString(string(output))
	return 1
}

func Markdown(content string) func(Context) {
	return func(ctx Context) {
		renderer := blackfriday.HtmlRenderer(mdown_flags, "", "")
		output := blackfriday.Markdown([]byte(content), renderer, mdown_extensions)
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
	return WriteTag("body", Atr, append([]func(Context){baseStyle}, inner...)...)
}

func (ctx Context) indentMultiline(str string) {
	var lines = strings.Split(strings.Trim(str, " \r\n"), "\n")
	for _, line := range lines {
		ctx.writeLine(strings.Trim(line, " \r\n"))
	}
}

var baseStyle = WriteTag("style", Atr, func(ctx Context) {
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

func WriteVoidTag(tagname string, attributes AttributeChain) func(Context) {
	return func(ctx Context) {
		ctx.startLine()
		ctx.write("<" + tagname)
		ctx.writeAttributes(attributes)
		ctx.write(">")
		ctx.endLine()
	}
}

func (ctx Context) WriteTags(inner ...func(ctx Context)) {
	innerCtx := Context{ctx.indentCount + 1, ctx.w, ctx.themeName}
	for _, fn := range inner {
		fn(innerCtx)
	}
}

func WriteTag(tagname string, attributes AttributeChain, inner ...func(ctx Context)) func(ctx Context) {
	return func(ctx Context) {
		ctx.startLine()
		ctx.write("<")
		ctx.write(tagname)
		ctx.writeAttributes(attributes)
		ctx.write(">")
		ctx.endLine()
		ctx.WriteTags(inner...)
		ctx.writeLine("</" + tagname + ">")
	}
}
