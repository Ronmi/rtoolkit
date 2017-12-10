package sam

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
)

// SAM stands for Semi-Auto Migration
//
// You can focus on changes of you database schemas and SAM executes the SQL
// files for you.
//
// SQL files MUST be placed in same folder, and named in following pattern:
//
//     [0-9]+-.*\.sql
//
// SAM needs some machanism to "remember" current state. Here's virtual code
// demonstrates how SAM load current state (in SAM.Load() and SAM.Save()):
//
//     SELECT state_column_name FROM table_name WHERE app_name_column="app_name"
//
// It uses following schema if you are not providing one:
//
//     CREATE TABLE sam_state (
//         app VARCHAR PRIMARY KEY,
//         state int
//     )
//
// Although every SQL file execution is protected in transaction, SAM is
// neither thread-safe nor reentrant: You SHOULD NEVER run same folder/app
// repeatly (which is also nonsence).
//
// SAM splits SQL statements in same file by ";[ \t\r\n]*\n", for supporting SQL
// drivers that cannot execute multiple statements in one DB.Exec() call.
type SAM struct {
	// For supporting virtual filesystems like go-bindata. Leave nil
	// to use filepath.Walk
	Walker Walker

	// For supporting virtual filesystems. Leave nil to use ioutil.ReadFile
	Reader func(path string) ([]byte, error)

	// Database schema settings, need to be quoted, like "`order`" if using
	// MySQL
	//
	// Ignored if you are not using SAM.Load() or SAM.Save()
	Table    string // default "sam_state"
	ColApp   string // default "app", should be long enough for app names
	ColState string // default "state", should be integer type and big enough for you apps
}

func (s *SAM) walker() Walker {
	if s.Walker == nil {
		return filepath.Walk
	}

	return s.Walker
}

func (s *SAM) reader() func(string) ([]byte, error) {
	if s.Reader == nil {
		return ioutil.ReadFile
	}

	return s.Reader
}

func strDEF(in, def string) string {
	if in == "" {
		return def
	}

	return in
}

func (s *SAM) schema() (table, app, state string) {
	table = strDEF(s.Table, "sam_state")
	app = strDEF(s.ColApp, "app")
	state = strDEF(s.ColState, "state")

	return
}

// Execute executes matched sql files in root path
//
// You might have to handle state records on yourself. But luckilly we have pre-made
// helpers: SAM.Load() and SAM.Save()
//
//     id, err := sam.Execute(db, "path", "myapp", sam.Load(db, "myapp"))
//     if err != nil {
//         // handle error
//     }
//     if err = sam.Save(db, "myapp", id); err != nil {
//         // handle error
//     }
//
// The whole execution is wrapped in transaction, ensuring all-or-nothing behavier.
func (s *SAM) Execute(db *sql.DB, root string, state int) (int, error) {
	files := sqlFiles(s.walker(), root, state)
	if len(files) == 0 {
		return state, nil
	}

	tx, err := db.Begin()
	if err != nil {
		return state, err
	}

	for _, file := range files {
		data, err := s.reader()(file.fn)
		if err != nil {
			return state, err
		}

		qstrs := split(data)
		for _, q := range qstrs {
			if _, err := tx.Exec(string(q)); err != nil {
				tx.Rollback()
				return state, err
			}
		}

		state = file.id
	}
	tx.Commit()

	return state, nil
}
