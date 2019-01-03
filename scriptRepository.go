package main

import (
	"database/sql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
	"time"
)

type scriptDB struct {
	Id        int64         `db:"Id"`
	Name      string        `db:"Name"`
	Content   string        `db:"Content"`
	Created   sql.NullInt64 `db:"Created"`
	Updated   sql.NullInt64 `db:"Updated"`
	Deleted   sql.NullInt64 `db:"Deleted"`
	IsLibrary string        `db:"IsPage"`
}

func GetScriptByName(name string) (Script, error) {
	var s = &scriptDB{}
	err := db.Get(s, `
		Select 
		    ScriptID as Id,
		    Name,
		    Content,
		    Created,
		    Updated,
		    Deleted
		from Scripts 
		where Name = ?
		`, name)
	if err != nil {
		return Script{}, err
	}

	var script = Script{
		Id:      s.Id,
		Name:    s.Name,
		Content: s.Content,
		Created: time.Unix(s.Created.Int64, 0),
	}
	if s.Updated.Valid {
		script.Updated = time.Unix(s.Updated.Int64, 0)
	}
	return script, err
}

func ListScripts() ([]Script, error) {
	var scripts = &[]scriptDB{}
	err := sqlx.Select(db, scripts,
		`Select 
		    ScriptID as Id,
		    Name,
		    Content,
		    Created,
		    Updated,
		    Deleted
		from Scripts`)
	if err != nil {
		return nil, err
	}

	var retScripts = make([]Script, len(*scripts))
	for i, s := range *scripts {
		retScripts[i] = Script{
			Id:      s.Id,
			Name:    s.Name,
			Content: s.Content,
			Created: time.Unix(s.Created.Int64, 0),
		}
		if s.Updated.Valid {
			retScripts[i].Updated = time.Unix(s.Updated.Int64, 0)
		}
	}

	return retScripts, nil
}

func CreateScript(name, code string) error {
	_, err := db.Exec(`
		INSERT OR FAIL 
		into Scripts(Name, Content, Created) 
		values (?, ?, strftime('%s', 'now'));`, name, code)
	return err
}

func RenameScript(currentName, newName string) error {
	db.MustExec(`
		Insert Into ScriptHistory (ScriptId, Name, Content, Created, Updated)
		Select Id, Name, Content, Created, strftime('%s', 'now') from Scripts where Name = ?;
		Update Scripts Set Name = ?, Updated where Name = ?;
		`, currentName, newName, currentName)
	return nil
}

func SaveScript(name, code string) error {
	_, err := db.Exec(`
		INSERT OR IGNORE
		into Scripts(Name, Content, Updated) 
		values (?, ?, strftime('%s', 'now'));

		Insert Into ScriptHistory (ScriptId, Name, Content, Created, Updated)
		Select ScriptId, Name, Content, Created, strftime('%s', 'now') from Scripts where Name = ?;

		Update Scripts
		Set 
		  Content = ?,
		  Updated = strftime('%s', 'now')
		where Name = ?;
		`, name, code, name, code, name)
	return err
}
