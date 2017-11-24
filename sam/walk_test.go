package sam

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"testing"
	"time"
)

type fileInfoMock bool

func (m fileInfoMock) Name() string       { return "" }
func (m fileInfoMock) Size() int64        { return 0 }
func (m fileInfoMock) Mode() os.FileMode  { return 0 }
func (m fileInfoMock) ModTime() time.Time { return time.Now() }
func (m fileInfoMock) IsDir() bool        { return bool(m) }
func (m fileInfoMock) Sys() interface{}   { return nil }

var aDir = fileInfoMock(true)
var aFile = fileInfoMock(false)

func TestWalkFuncError(t *testing.T) {
	err := errors.New("my error")
	_, actual := walkfunc("", aDir, err)
	if actual != err {
		t.Fatalf("expected original error, got %s", actual)
	}
}

func TestWalkFuncDir(t *testing.T) {
	_, actual := walkfunc("", aDir, nil)
	if actual != filepath.SkipDir {
		t.Fatalf("expected SkipDir, got %s", actual)
	}
}

func TestWalkFuncWrongFilename(t *testing.T) {
	cases := [][2]string{
		{"", "empty"},
		{"a/00-sql", "ext"},
		{"a/00-.sql", "desc"},
		{"a/-a.sql", "id"},
		{"a/00a.sql", "dash"},
	}

	for _, c := range cases {
		t.Run(c[1], func(t *testing.T) {
			_, actual := walkfunc(c[0], aFile, nil)
			if actual == nil {
				t.Fatal("expected an error, but nothing happend")
			}

			if _, ok := actual.(errMismatch); !ok {
				t.Fatalf("expected mismatch error, got %s", actual)
			}
		})
	}
}

func TestWalkFuncID(t *testing.T) {
	cases := []struct {
		fn string
		id int
	}{
		{fn: "a/00-a.sql", id: 0},
		{fn: "a/01-a.sql", id: 1},
		{fn: "a/1-a.sql", id: 1},
		{fn: "a/09-a.sql", id: 9},
		{fn: "a/10-a.sql", id: 10},
	}

	for _, c := range cases {
		t.Run(c.fn, func(t *testing.T) {
			id, err := walkfunc(c.fn, aFile, nil)
			if err != nil {
				t.Fatalf("unexpected error: %s", err)
			}

			if id != c.id {
				t.Fatalf("expected %d, got %d", c.id, id)
			}
		})
	}
}

func TestSQLFiles(t *testing.T) {
	exp := []sqlFile{
		{id: 1, fn: "testdata/1-a.sql"},
		{id: 2, fn: "testdata/02-b.sql"},
		{id: 3, fn: "testdata/3-c.sql"},
	}

	cases := []struct {
		expect []sqlFile
		cur    int
	}{
		{expect: exp, cur: 0},
		{expect: exp[1:], cur: 1},
		{expect: exp[2:], cur: 2},
		{expect: []sqlFile{}, cur: 3},
		{expect: []sqlFile{}, cur: 4},
	}

	for _, c := range cases {
		t.Run(fmt.Sprintf("cur#%d", c.cur), func(t *testing.T) {
			actual := sqlFiles(filepath.Walk, "testdata/", c.cur)
			if reflect.DeepEqual(c.expect, actual) {
				return
			}

			t.Fatalf("expected [%#v], got [%#v]", c.expect, actual)
		})
	}
}
