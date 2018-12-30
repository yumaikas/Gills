package templates

var Atr = AttributeChain(make([]Attribute, 0))

type AttributeChain []Attribute

func (attrs AttributeChain) Add(key, value string) AttributeChain {
	return append(attrs, Attribute{Key: key, Value: value, Trusted: false})
}

func (attrs AttributeChain) AddVoid(key string) AttributeChain {
	return append(attrs, Attribute{Key: key, Trusted: false, Void: true})
}
func (attrs AttributeChain) AddUnsafe(key, value string) AttributeChain {
	return append(attrs, Attribute{Key: key, Value: value, Trusted: true})
}

func (attrs AttributeChain) Id(id string) AttributeChain {
	return attrs.Add("id", id)
}

func (attrs AttributeChain) Value(value string) AttributeChain {
	return attrs.Add("value", value)
}

func (attrs AttributeChain) Type(value string) AttributeChain {
	return attrs.Add("type", value)
}

func (attrs AttributeChain) Class(class string) AttributeChain {
	return attrs.Add("class", class)
}

func (attrs AttributeChain) For(_for string) AttributeChain {
	return attrs.Add("for", _for)
}

func (attrs AttributeChain) Name(name string) AttributeChain {
	return attrs.Add("name", name)
}

func (attrs AttributeChain) Cols(name string) AttributeChain {
	return attrs.Add("cols", name)
}

func (attrs AttributeChain) Rows(name string) AttributeChain {
	return attrs.Add("rows", name)
}

func (attrs AttributeChain) Size(name string) AttributeChain {
	return attrs.Add("size", name)
}

func (attrs AttributeChain) FormMethod(name string) AttributeChain {
	return attrs.Add("formmethod", name)
}

func (attrs AttributeChain) FormAction(name string) AttributeChain {
	return attrs.Add("formaction", name)
}

func (attrs AttributeChain) Form(name string) AttributeChain {
	return attrs.Add("form", name)
}

func (attrs AttributeChain) Method(name string) AttributeChain {
	return attrs.Add("method", name)
}

func (attrs AttributeChain) Action(name string) AttributeChain {
	return attrs.Add("action", name)
}

func (attrs AttributeChain) Accept(name string) AttributeChain {
	return attrs.Add("accept", name)
}

func (attrs AttributeChain) Multiple() AttributeChain {
	return attrs.AddVoid("multiple")
}

func (attrs AttributeChain) EncType(enctype string) AttributeChain {
	return attrs.Add("enctype", enctype)
}

func (attrs AttributeChain) Href(href string) AttributeChain {
	return attrs.Add("href", href)
}

func (attrs AttributeChain) Rel(rel string) AttributeChain {
	return attrs.Add("rel", rel)
}

func (attrs AttributeChain) Src(src string) AttributeChain {
	return attrs.Add("src", src)
}

func (attrs AttributeChain) UnsafeHref(href string) AttributeChain {
	return attrs.AddUnsafe("href", href)
}
