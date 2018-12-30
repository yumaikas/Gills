package main

import (
	"strings"
	"time"
)

type Note struct {
	Id      int64
	Content string
	Created time.Time
	Updated time.Time
	Deleted time.Time
}

type Script struct {
	Id      int64
	Name    string
	Content string
	Created time.Time
	Updated time.Time
	Deleted time.Time
}

func (s Script) IsPage() bool {
	return strings.HasSuffix(s.Name, ".page")
}

func (s Script) IsLibrary() bool {
	return strings.HasSuffix(s.Name, ".lib")
}
