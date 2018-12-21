package templates

import (
	"html/template"
	"strings"
)

func Br() func(Context) {
	return writeVoidTag("br", Atr)
}

func Div(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("div", attributes, inner...)
}

func Span(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("span", attributes, inner...)
}

func Style(inner string) func(Context) {
	return writeTag("style", Atr, func(ctx Context) {
		ctx.indentMultiline(inner)
	})
}

func H1(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("h2", attributes, inner...)
}

func H2(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("h2", attributes, inner...)
}

func Title(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("title", attributes, inner...)
}

func A(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("a", attributes, inner...)
}
func P(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("p", attributes, inner...)
}

func Table(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("Table", attributes, inner...)
}

func Td(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("td", attributes, inner...)
}

func Tr(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("tr", attributes, inner...)
}

func Form(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("form", attributes, inner...)
}

func Label(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("label", attributes, inner...)
}

func Button(attributes AttributeChain, inner ...func(Context)) func(Context) {
	return writeTag("button", attributes, inner...)
}

func Input(attributes AttributeChain) func(Context) {
	return writeVoidTag("input", attributes)
}

func JS(script string) func(Context) {
	return writeTag("script", Atr, func(ctx Context) {
		ctx.write(script)
	})
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
