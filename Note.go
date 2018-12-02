package main

import "time"

type Note struct {
	Id      int64
	Content string
	Created time.Time
	Updated time.Time
	Deleted time.Time
}
