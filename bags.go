package main

type SingleBag map[string]string
type MultiBag map[string][]string

type ChainBag []Bag

type Bag interface {
	GetOr(key, value string) string
	get(key string) (string, bool)
}

func (mb MultiBag) GetOr(key, fallback string) string {
	if value, found := mb[key]; found && len(value) > 0 {
		return value[0]
	}
	return fallback
}

func (b SingleBag) GetOr(key, fallback string) string {
	if value, ok := b[key]; ok {
		return value
	}
	return fallback
}

func (cb ChainBag) GetOr(key, fallback string) string {
	for _, innerBag := range cb {
		val, found := innerBag.get(key)
		if found {
			return val
		}
	}
	return fallback
}

func (mb MultiBag) get(key string) (string, bool) {
	if value, found := mb[key]; found && len(value) > 0 {
		return value[0], true
	}
	return "", false
}

func (b SingleBag) get(key string) (string, bool) {
	if value, ok := b[key]; ok {
		return value, ok
	}
	return "", false
}

func (cb ChainBag) get(key string) (string, bool) {
	for _, innerBag := range cb {
		val, found := innerBag.get(key)
		if found {
			return val, true
		}
	}
	return "", false
}

func (cb ChainBag) BackedBy(inner Bag) ChainBag {
	return append(cb, inner)
}

func (b MultiBag) BackedBy(inner Bag) ChainBag {
	return []Bag{b, inner}
}

func (b SingleBag) BackedBy(inner Bag) ChainBag {
	return []Bag{b, inner}
}
