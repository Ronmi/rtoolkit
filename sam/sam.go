// Package sam provides semi-auto database migrating mechanism.
package sam

import (
	"bytes"
	"database/sql"
	"errors"
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
// SAM needs some machanism to "remember" current state, which is defined at
// type Storage. We also provides few simple storages for you to use: take a
// look at package storage.
//
// You are encouraged to write your own storage for fitting your needs.
type SAM struct {
	// For supporting virtual filesystems like go-bindata. Leave nil
	// to use filepath.Walk
	Walker Walker

	// For supporting virtual filesystems. Leave nil to use ioutil.ReadFile
	Reader func(path string) ([]byte, error)

	// Defines how to persists db versions. Leave nil to use storage.Default
	Storage Storage
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

func split(data []byte) [][]byte {
	arr := bytes.Split(data, []byte("\n"))
	qstrs := make([][]byte, 0, len(arr))

	begin := 0
	cur := 0
	for _, line := range arr {
		cur += len(line) + 1

		// test for ;[ \t\r\n]*\n
		l := bytes.TrimRight(line, " \t\r\n")
		if len(l) > 0 && l[len(l)-1] == ';' {
			qstrs = append(
				qstrs,
				bytes.TrimRight(data[begin:cur], " \t\r\n\x00"),
			)
			begin = cur
		}
	}

	if begin != cur {
		qstr := data[begin:cur]
		if q := bytes.Trim(qstr, " \t\r\n\x00"); len(q) > 0 {
			qstrs = append(qstrs, q)
		}
	}

	return qstrs
}

func wrap(fn string, err error) error {
	return errors.New("while executing " + fn + ": " + err.Error())
}

// Execute executes matched sql files in root path
//
// You might have to handle state records on yourself. But luckilly we have pre-made
// helpers, see type Storage and package storage for example:
//
//     id, err := storage.Default.Load("myapp")
//     if err != nil {
//         // handle error
//     }
//
//     if id, err = sam.Execute(db, "path", "myapp", id); err != nil {
//         // handle error
//     }
//
//     if err = sotrage.Default.Save("myapp", id); err != nil {
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
		// no need to check here since sqlFile() has filtered them out
		data, err := s.reader()(file.fn)
		if err != nil {
			return state, wrap(file.fn, err)
		}

		qstrs := split(data)
		for _, q := range qstrs {
			if _, err := tx.Exec(string(q)); err != nil {
				tx.Rollback()
				return state, wrap(file.fn, err)
			}
		}

		state = file.id + 1
	}
	tx.Commit()

	return state, nil
}

// MustExec is Load()+Execute()+Save() for lazy people, panics if any error
func (s *SAM) MustExec(db *sql.DB, root string, app string, st Storage) {
	id, err := st.Load(app)
	if err != nil {
		panic(err)
	}
	if id, err = s.Execute(db, root, id); err != nil {
		panic(err)
	}
	if err = st.Save(app, id); err != nil {
		panic(err)
	}
}
