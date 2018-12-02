package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func SearchNotes(database *sql.DB) error {
	// var results =
	return nil
}

func LoadState() (map[string]string, error) {
	//
	return nil, nil
}

func DeleteNote(id int64) error {
	// TODO
}

func SaveNote(note Note) error {
	// TODO
}

func buildSchema(database *sql.DB) error {
	schema := `
	Create Table If Not Exists Notes (
		Id INTEGER PRIMARY KEY,
		Content text,
		Created int, -- unix timestamp
		Update int, -- unix timestamp
		Deleted int -- unix timestamp
	);

	Create Table If Not Exists NoteHistory (
		Id INTEGER PRIMARY KEY,
		NoteId integer,
		Content text,
		Created int, -- unix timestamp
		Update int, -- unix timestamp
		Deleted int -- unix timestamp
	);

	Create Table If Not Exists StateKV (
		KvId INTEGER PRIMARY KEY,
		Key text,
		Context text,
		Created int, -- unix timestamp
		Update int, -- unix timestamp
		Deleted int -- unix timestamp
	);
	Create Index If Not Exists KVidx StateKV(Key, Content);`

	_, err := database.Exec(schema)
	return err
}
