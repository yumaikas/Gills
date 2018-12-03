package main

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type noteDB struct {
	Id      int64  `db:"Id"`
	Content string `db:"Content"`
	Created int64  `db:"Created"`
	Updated int64  `db:"Updated"`
	Deleted int64  `db:"Deleted"`
}
type noteHistoryDB struct {
	Id      int64  `db:"Id"`
	NoteId  int64  `db:"NoteId"`
	Content string `db:"Content"`
	Created int64  `db:"Created"`
	Updated int64  `db:"Updated"`
	Deleted int64  `db:"Deleted"`
}

type Bag map[string]string

func (b Bag) getOr(key, fallback string) string {
	if value, ok := b[key]; ok {
		return value
	}
	return fallback
}

var db *sqlx.DB

func InitDB(path string) error {
	var err error
	db, err = sqlx.Open("sqlite3", path)
	if err != nil {
		return err
	}
	return buildSchema(db)
}

// Assumes an open DB connection
func SearchNotes(searchTerms string) ([]Note, error) {
	notes := &[]noteDB{}
	err := sqlx.Select(db, notes, `
		Select NoteId as Id, Content, Created, Updated 
		from Notes where Content like '%' || ? || '%'  
		`, searchTerms)
	if err != nil {
		return nil, err
	}
	retVals := make([]Note, len(*notes))
	for idx, note := range *notes {
		retVals[idx] = Note{
			Id:      note.Id,
			Content: note.Content,
			Created: time.Unix(note.Created, 0),
			Updated: time.Unix(note.Updated, 0),
		}
	}
	return retVals, nil
}

type KV struct {
	Key     string `db:"Key"`
	Content string `db:"Content"`
}

func LoadState() (Bag, error) {
	results := []KV{}
	err := db.Select(&results, `Select Key, Content from StateKV;`)
	if err != nil {
		return nil, err
	}
	var retVal = make(map[string]string)
	for _, kv := range results {
		retVal[kv.Key] = kv.Content
	}
	return retVal, nil
}

func SaveState(state []KV) error {
	tx := db.MustBegin()
	defer tx.Rollback()
	tx.MustExec("Delete from StateKV")
	for _, kv := range state {
		tx.MustExec(`
			Insert into StateKV 
				(Key, Content, Created) 
			values (?, ?, strftime('%s', 'now')`, kv.Key, kv.Content)
	}
	return tx.Commit()
}

func DeleteNote(id int64) error {
	_, err := db.Exec(`
		Insert into NotesHistory (NoteID, Content, Created, Deleted)
		Select ?, Content, Created, strftime('%s', 'now')
		from Notes where NoteID = ? 
		`, id, id)
	return err
}

func SaveNote(note Note) (int64, error) {
	// 0 ID means that this note isn't in the database
	// https://www.sqlite.org/autoinc.html
	// >  If the table is initially empty, then a ROWID of 1 is used
	if note.Id != 0 {
		db.MustExec(`
			Insert into NotesHistory 
				(NoteId, Content, Created, Updated)
				Select ?, Content, Created, strftime('%s', 'now')
				from Notes where NoteID = ?

				Update Notes
				Set 
					Content = ?,
					Updated = strftime('%s', 'now')
				where NoteID = ? 
			`, note.Id, note.Id, note.Content, note.Id)
		return note.Id, nil
	} else {
		var retVal int64
		err := db.Get(&retVal,
			`Insert into Notes (Content, Created) values (?, strftime('%s', 'now');
			Select last_insert_rowID();`, note.Content)
		return retVal, err
	}
}

func buildSchema(database *sqlx.DB) error {
	schema := `
	Create Table If Not Exists Notes (
		NoteId INTEGER PRIMARY KEY,
		Content text,
		Created int, -- unix timestamp
		Updated int, -- unix timestamp
		Deleted int -- unix timestamp
	);

	Create Table If Not Exists NoteHistory (
		Id INTEGER PRIMARY KEY,
		NoteId integer,
		Content text,
		Created int, -- unix timestamp
		Updated int, -- unix timestamp
		Deleted int -- unix timestamp
	);

	Create Table If Not Exists StateKV (
		KvId INTEGER PRIMARY KEY,
		Key text,
		Content text,
		Created int, -- unix timestamp
		Updated int, -- unix timestamp
		Deleted int -- unix timestamp
	);
	Create Index If Not Exists KVidx ON StateKV(Key, Content);`

	_, err := database.Exec(schema)
	return err
}
