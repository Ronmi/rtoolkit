package sam

import (
	"bytes"
	"database/sql"
	"fmt"
)

// Load tries to load state from DB using settings in SAM
//
// Load never fails: It returns -1 if anything goes wrong.
//
// You are encouraged to use your own implementation!
func (s *SAM) Load(db *sql.DB, name string) int {
	table, app, state := s.schema()

	qstr := fmt.Sprintf(`SELECT %s FROM %s WHERE %s = ?`, state, table, app)
	row := db.QueryRow(qstr, name)

	cur := -1
	if err := row.Scan(&cur); err != nil {
		cur = -1
	}
	return cur
}

// Save tries to save state to DB using settings in SAM
//
// To be compitable with most SQL servers, it tries insert into table and
// try update if failed.
//
// You are encouraged to use your own implementation!
func (s *SAM) Save(db *sql.DB, name string, cur int) error {
	table, app, state := s.schema()

	ins := fmt.Sprintf(`INSERT INTO %s (%s,%s) VALUES (?,?)`, table, state, app)
	upd := fmt.Sprintf(`UPDATE %s SET %s=? WHERE %s=?`, table, state, app)

	if _, err := db.Exec(ins, cur, name); err == nil {
		return nil
	}

	if _, err := db.Exec(upd, cur, name); err != nil {
		return err
	}

	return nil
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

func (s *SAM) run(db *sql.DB, fn string) (cur int, err error) {
	return
}
