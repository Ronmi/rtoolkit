package sam

import (
	"database/sql"
	"testing"

	_ "github.com/mxk/go-sqlite/sqlite3"
)

func db(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		t.Fatalf("unexpected error to create sqlite database: %s", err)
	}

	return db
}

func testResult(t *testing.T, db *sql.DB) {
	qstr := `SELECT * FROM c`
	_, err := db.Exec(qstr)
	if err != nil {
		t.Fatalf("table c does not exists?: %s", err)
	}
}

func TestExecuteAll(t *testing.T) {
	db := db(t)
	s := &SAM{}
	id, err := s.Execute(db, "testdata/sqlite/", 0)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if id != 4 {
		t.Fatalf("expected version 4, got %d", id)
	}

	testResult(t, db)
}

func TestExecuteNothing(t *testing.T) {
	db := db(t)
	s := &SAM{}
	s.Execute(db, "testdata/sqlite/", 0)
	testResult(t, db)
	id, err := s.Execute(db, "testdata/sqlite/", 4)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if id != 4 {
		t.Fatalf("expected version 4, got %d", id)
	}
	testResult(t, db)
}

func TestExecuteFromMiddle(t *testing.T) {
	db := db(t)
	db.Exec(`CREATE TABLE a (x INT, y INT)`)
	s := &SAM{}
	id, err := s.Execute(db, "testdata/sqlite/", 2)
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if id != 4 {
		t.Fatalf("expected version 4, got %d", id)
	}

	testResult(t, db)
}
