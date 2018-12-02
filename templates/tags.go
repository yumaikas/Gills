package templates

func Br() func(Context) {
	return writeVoidTag("br", NoAttr)
}

func Div(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("div", attributes, inner)
}

func Div_(inner func(Context)) func(Context) {
	return writeTag("div", NoAttr, inner)
}

func H1(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("h2", attributes, inner)
}

func H2(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("h2", attributes, inner)
}

func A(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("a", attributes, inner)
}
func P(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("p", attributes, inner)
}

func Table(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("Table", attributes, inner)
}

func Td(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("td", attributes, inner)
}

func Td_(inner func(Context)) func(Context) {
	return writeTag("td", NoAttr, inner)
}

func Tr(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("tr", attributes, inner)
}

func Tr_(inner func(Context)) func(Context) {
	return writeTag("tr", NoAttr, inner)
}

func Form(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("form", attributes, inner)
}

func Label(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("label", attributes, inner)
}

func Button(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("button", attributes, inner)
}

func Input(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("input", attributes, inner)
}

func TextArea(attributes Kv, inner func(Context)) func(Context) {
	return writeTag("textarea", attributes, inner)
}
