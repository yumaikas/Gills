package templates

import (
	"html/template"
	"strings"
)

func Br() func(Context) {
	return WriteVoidTag("br", Atr)
}
func Hr() func(Context) {
	return WriteVoidTag("hr", Atr)
}

func Div(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("div", attributes, inner...)
}

func Span(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("span", attributes, inner...)
}

func Style(inner string) func(Context) {
	return WriteTag("style", Atr, func(ctx Context) {
		ctx.indentMultiline(inner)
	})
}

func StyleLink(href string) func(Context) {
	return WriteVoidTag("link", Atr.Rel("stylesheet").Type("text/css").Href(href))
}

func H1(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("h2", attributes, inner...)
}

func H2(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("h2", attributes, inner...)
}

func H3(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("h3", attributes, inner...)
}

func Title(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("title", attributes, inner...)
}

func A(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("a", attributes, inner...)
}
func P(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("p", attributes, inner...)
}

func Table(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("Table", attributes, inner...)
}

func Td(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("td", attributes, inner...)
}

func Tr(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("tr", attributes, inner...)
}

func Form(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("form", attributes, inner...)
}

func Label(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("label", attributes, inner...)
}

func Button(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return WriteTag("button", attributes, inner...)
}

func Input(attributes AttributeChain) func(Context) {
	return WriteVoidTag("input", attributes)
}

func JS(script string) func(Context) {
	return WriteTag("script", Atr, func(ctx Context) {
		ctx.write(script)
	})
}

func JSLink(src string) func(Context) {
	return WriteTag("script", Atr.Src(src), func(ctx Context) {})
}

func TextArea(attributes AttributeChain, inner string) func(Context) {
	return func(ctx Context) {
		ctx.startLine()
		ctx.write("<textarea")
		ctx.writeAttributes(attributes)
		ctx.write(">")
		template.HTMLEscape(ctx.w, []byte(strings.Trim(inner, " \t\r\n")))
		ctx.write("</textarea>")
		ctx.endLine()
	}
}
