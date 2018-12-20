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

// Use this to mark upload display notes in a way that won't get accidentally set by a human
const UploadUUID = "8afee0a6-9ec4-46c1-a530-89287f579300"

func InitDB(path string) error {
	var err error
	db, err = sqlx.Open("sqlite3", path)
	if err != nil {
		return err
	}
	return buildSchema(db)
}

func queryNotes(query, searchTerms string) ([]Note, error) {
	notes := &[]noteDB{}
	err := sqlx.Select(db, notes, query, searchTerms)
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
	return retVals, nil
}

func SearchUploadNotes(searchTerms string) ([]Note, error) {
	return queryNotes(`
		Select NoteId as Id, Content, Created, Updated 
		from Notes 
			where Content like '%' || ? || '%'  
			AND Content like '%@upload-8afee0a6-9ec4-46c1-a530-89287f579300%'
		order by Created DESC
		`, searchTerms)

}

func SearchRecentNotes(searchTerms string) ([]Note, error) {
	return queryNotes(`
		Select NoteId as Id, Content, Created, Updated 
		from Notes 
			where Content like '%' || ? || '%'  
			AND Content not like '%@archive%'
		order by Created DESC
		`, searchTerms)
}

// Assumes an open DB connection
func SearchNotes(searchTerms string) ([]Note, error) {
	return queryNotes(`
		Select NoteId as Id, Content, Created, Updated 
		from Notes 
		where 
			Content like '%' || ? || '%'  
			AND Content not like '%@archive%'
		order by Created DESC
		`, searchTerms)
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

func GetNoteBy(id int64) (Note, error) {
	var n = &noteDB{}
	err := db.Get(n, `Select NoteId as Id, Content, Created, Updated, Deleted from Notes where NoteId = ?`, id)
	if err != nil {
		return Note{}, err
	}
	var note = Note{
		Id:      n.Id,
		Content: n.Content,
		Created: time.Unix(n.Created.Int64, 0),
	}

	if n.Updated.Valid {
		note.Updated = time.Unix(n.Created.Int64, 0)
	}

	return note, err
}

func DeleteNote(id int64) error {
	tx := db.MustBegin()
	defer tx.Rollback()
	tx.MustExec(`
		Insert into NoteHistory (NoteID, Content, Created, Deleted)
		Select ?, Content, Created, strftime('%s', 'now')
		from Notes where NoteID = ? 
		`, id, id)
	tx.MustExec(`Delete from Notes where NoteId = ?`, id)
	return tx.Commit()
}

func SaveNote(note Note) (int64, error) {
	// 0 ID means that this note isn't in the database
	// https://www.sqlite.org/autoinc.html
	// >  If the table is initially empty, then a ROWID of 1 is used
	if note.Id != 0 {
		db.MustExec(`
			Insert into NoteHistory 
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
