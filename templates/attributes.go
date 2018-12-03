package templates

var Atr = AttributeChain(make([]Attribute, 0))

type AttributeChain []Attribute

func (attrs AttributeChain) Add(key, value string) AttributeChain {
	return append(attrs, Attribute{Key: key, Value: value, Trusted: false})
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

func (attrs AttributeChain) Href(href string) AttributeChain {
	return attrs.Add("href", href)
}

func (attrs AttributeChain) UnsafeHref(href string) AttributeChain {
	return attrs.AddUnsafe("href", href)
}
