package sam

import (
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
// SAM needs a table to "remember" current state. Here's virtual code
// demonstrates how SAM load current state:
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
// neither thread-safe nor reentrant: You SHOULD NEVER run same folder
// concurrently.
//
// SAM splits SQL statements in same file by ";\n", for supporting SQL drivers
// that does not support multi-statements execution.
type SAM struct {
	// For supporting virtual filesystems like go-bindata. Leave nil
	// to use filepath.Walk
	Walker Walker

	// For supporting virtual filesystems. Leave nil to use ioutil.ReadFile
	Reader func(path string) ([]byte, error)

	// Database schema settings
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
