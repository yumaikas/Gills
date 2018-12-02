package templates

func Br() func(context) {
	return writeVoidTag("br", NoAttr)
}

func Div(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("div", attributes, inner)
}

func H1(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("h2", attributes, inner)
}

func H2(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("h2", attributes, inner)
}

func A(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("a", attributes, inner)
}
func P(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("p", attributes, inner)
}

func Table(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("Table", attributes, inner)
}

func Td(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("td", attributes, inner)
}

func Td_(inner func(context)) func(context) {
	return writeTag("td", NoAttr, inner)
}

func Tr(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("tr", attributes, inner)
}

func Tr_(inner func(context)) func(context) {
	return writeTag("tr", NoAttr, inner)
}

func Form(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("form", attributes, inner)
}

func Label(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("label", attributes, inner)
}

func Button(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("button", attributes, inner)
}

func Input(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("input", attributes, inner)
}

func TextArea(attributes []Attribute, inner func(context)) func(context) {
	return writeTag("textarea", attributes, inner)
}
