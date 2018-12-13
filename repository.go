package main

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type noteDB struct {
	Id      int64         `db:"Id"`
	Content string        `db:"Content"`
	Created sql.NullInt64 `db:"Created"`
	Updated sql.NullInt64 `db:"Updated"`
	Deleted sql.NullInt64 `db:"Deleted"`
}
type noteHistoryDB struct {
	Id      int64  `db:"Id"`
	NoteId  int64  `db:"NoteId"`
	Content string `db:"Content"`
	Created int64  `db:"Created"`
	Updated int64  `db:"Updated"`
	Deleted int64  `db:"Deleted"`
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
func SearchRecentNotes(searchTerms string) ([]Note, error) {
	notes := &[]noteDB{}
	err := sqlx.Select(db, notes, `
		Select NoteId as Id, Content, Created, Updated 
		from Notes where Content like '%' || ? || '%'  
		order by Created DESC
		`, searchTerms)
	if err != nil {
		return nil, err
	}
	retVals := make([]Note, len(*notes))
	for idx, note := range *notes {
		retVals[idx] = Note{
			Id:      note.Id,
			Content: note.Content,
			// I don't care if these are null for now.
			Created: time.Unix(note.Created.Int64, 0),
			Updated: time.Unix(note.Updated.Int64, 0),
		}
	}
	fmt.Println("Search returned", len(retVals), "notes")
	return retVals, nil
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
			// I don't care if these are null for now.
			Created: time.Unix(note.Created.Int64, 0),
			Updated: time.Unix(note.Updated.Int64, 0),
		}
	}
	fmt.Println("Search returned", len(retVals), "notes")
	return retVals, nil
}

type KV struct {
	Key     string `db:"Key"`
	Content string `db:"Content"`
}

func LoadState() (AppState, error) {
	results := []KV{}
	err := db.Select(&results, `Select Key, Content from StateKV;`)
	if err != nil {
		return AppState{}, err
	}
	var retVal = make(map[string]string)
	for _, kv := range results {
		retVal[kv.Key] = kv.Content
	}
	return AppState{SingleBag(retVal)}, nil
}

func SaveState(state []KV) error {
	tx := db.MustBegin()
	defer tx.Rollback()
	tx.MustExec("Delete from StateKV;")
	for _, kv := range state {
		tx.MustExec(`
			Insert into StateKV (Key, Content, Created)
			values (?, ?, strftime('%s', 'now'))`, kv.Key, kv.Content)
	}
	return tx.Commit()
}

func DeleteNote(id int64) error {
	tx := db.MustBegin()
	defer tx.Rollback()
	tx.MustExec(`
		Insert into NotesHistory (NoteID, Content, Created, Deleted)
		Select ?, Content, Created, strftime('%s', 'now')
		from Notes where NoteID = ? 
		`, id, id)
	tx.MustExec(`Delete from Notes where NotesId = ?`, id)
	return tx.Commit()
}

func SaveNote(note Note) (int64, error) {
	// 0 ID means that this note isn't in the database
	// https://www.sqlite.org/autoinc.html
	// >  If the table is initially empty, then a ROWID of 1 is used
	if note.Id != 0 {
		fmt.Print("QI")
		db.MustExec(`
			Insert into NotesHistory 
				(NoteId, Content, Created, Updated)
				Select ?, Content, Created, strftime('%s', 'now')
				from Notes where NoteID = ?;

				Update Notes
				Set 
					Content = ?,
					Updated = strftime('%s', 'now')
				where NoteID = ? 
			`, note.Id, note.Id, note.Content, note.Id)
		return note.Id, nil
	} else {
		fmt.Println("HEX")
		var retVal int64
		tx := db.MustBegin()
		defer tx.Rollback()
		tx.MustExec("Insert into Notes (Content, Created) values (?, strftime('%s', 'now'));", note.Content)
		err := tx.Get(&retVal, "Select last_insert_rowid()")
		if err != nil {
			return 0, err
		}
		return retVal, tx.Commit()
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
