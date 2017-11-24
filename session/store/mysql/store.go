package mysql

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/Ronmi/rtoolkit/session/store"
)

type mysqlStore struct {
	ttl     int
	conn    *sql.DB
	stmtGet *sql.Stmt
	stmtPut *sql.Stmt
	stmtTTL *sql.Stmt
	stmtCLR *sql.Stmt
	stmtNEW *sql.Stmt
	stmtGC  *sql.Stmt
	lastgc  int64 // unix timestamp, in seconds
}

// NewStore creates a MySQL store. You have to fill table and columns.
//
// There are few restrictions:
//
//   - expire column MUST be TIMESTAMP type.
//   - sessID column MUST be PRIMARY KEY or UNIQUE KEY.
//   - data column MUST be TEXT type.
//   - ttl column MUST be INT types.
//   - sessID and seed column MUST be CHAR(32) or VARCHAR(32) type.
//
// Here's an example:
//
//    CREATE session_storage (
//      sid CHAR(32) PRIMARY KEY,
//      seed CHAR(32),
//      data TEXT,
//      ttl INT(8),
//      expire TIMESTAMP,
//      INDEX time_to_live (expire ASC)
//    ) DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci;
func NewStore(db *sql.DB, table, sessID, seed, data, ttl, expire string) store.Store {
	ret := &mysqlStore{
		conn:   db,
		lastgc: time.Now().Unix(),
	}
	p := func(qstr string) *sql.Stmt {
		ret, err := db.Prepare(qstr)
		if err != nil {
			panic(err)
		}

		return ret
	}

	qstr := fmt.Sprintf(
		"SELECT `%s`,`%s` FROM `%s` WHERE `%s`=?",
		seed,
		data,
		table,
		sessID,
	)
	ret.stmtGet = p(qstr)

	qstr = fmt.Sprintf(
		"INSERT INTO `%s` (`%s`,`%s`,`%s`,`%s`,`%s`) VALUES (?,?,?,?,NOW()) ON DUPLICATE KEY UPDATE `%s`=VALUES(`%s`),`%s`=VALUES(`%s`),`%s`=VALUES(`%s`),`%s`=NOW()",
		table,
		sessID,
		seed,
		data,
		ttl,
		expire,
		seed, seed,
		data, data,
		ttl, ttl,
		expire,
	)
	ret.stmtPut = p(qstr)

	qstr = fmt.Sprintf(
		"UPDATE `%s` SET `%s`=?, `%s`=DATE_ADD(NOW(), INTERVAL `%s` SECOND) WHERE `%s`=? LIMIT 1",
		table,
		ttl,
		expire,
		ttl,
		sessID,
	)
	ret.stmtTTL = p(qstr)

	qstr = fmt.Sprintf(
		"DELETE FROM `%s` WHERE `%s`=?",
		table,
		sessID,
	)
	ret.stmtCLR = p(qstr)

	qstr = fmt.Sprintf(
		"INSERT INTO `%s` (`%s`,`%s`,`%s`,`%s`,`%s`) VALUES (?,?,'',?,DATE_ADD(NOW(), INTERVAL `%s` SECOND))",
		table,
		sessID,
		seed,
		data,
		ttl,
		expire,
		ttl,
	)
	ret.stmtNEW = p(qstr)

	qstr = fmt.Sprintf(
		"DELETE FROM `%s` WHERE NOW() > `%s`",
		table,
		expire,
	)
	ret.stmtGC = p(qstr)
	return ret
}

// SetTTL decides how long before data to be considered invalid (in seconds)
func (s *mysqlStore) SetTTL(ttl int) {
	s.ttl = ttl
}

func (s *mysqlStore) tryGC() {
	t := time.Now().Unix()
	if t < s.lastgc+int64(s.ttl) {
		return
	}

	s.lastgc = t
	s.stmtGC.Exec()
}

func (s *mysqlStore) tryInsert(sid, seed string) (bool, error) {
	s.tryGC()
	res, err := s.stmtNEW.Exec(sid, seed, s.ttl)
	if err != nil {
		return false, err
	}

	cnt, err := res.RowsAffected()
	if err != nil {
		return false, err
	}
	if cnt != 1 {
		return false, nil
	}

	return true, nil
}

// Allocate creates a new session id, returns error if store is full
//
// Size of session id depends on store.
//
// Implementation MUST allocate space for the session id before returning it,
// and ttl value MUST follow what was set by SetTTL().
//
// seed is used for session validating, see session.Manager.Start() for detail
func (s *mysqlStore) Allocate(seed string) (string, error) {
	var err error
	sid := store.GenerateRandomKey(32, func(id string) bool {
		var ret bool
		ret, err = s.tryInsert(id, seed)
		if err != nil {
			return true
		}

		return ret
	})

	return sid, err
}

// Get returns session data (string), returns error if not found or something goes wrong
//
// It MUST refresh ttl value.
func (s *mysqlStore) Get(sessID string) (seed, data string, err error) {
	res := s.stmtGet.QueryRow(sessID)
	err = res.Scan(&seed, &data)
	return
}

// Set saves session data, returns error if not found or something goes wrong
//
// It MUST refresh ttl value.
func (s *mysqlStore) Set(sessID string, seed, data string) error {
	_, err := s.stmtPut.Exec(sessID, seed, data, s.ttl)
	return err
}

// Release clears a session, never fail
func (s *mysqlStore) Release(sessID string) {
	s.stmtCLR.Exec(sessID)
}
