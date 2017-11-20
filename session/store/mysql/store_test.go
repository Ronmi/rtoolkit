package mysql

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB

func init() {
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		log.Fatal("You must set MYSQL_DSN before running test")
	}

	var err error
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Cannot connect to mysql: %s", err)
	}

	db.Exec("DROP TABLE IF EXISTS session_store")
	qstr := `CREATE TABLE session_store (
  sid CHAR(32) PRIMARY KEY,
  seed CHAR(32),
  data TEXT,
  ttl INT(8),
  expire TIMESTAMP,
  INDEX time_to_live (expire ASC)
) DEFAULT CHARACTER SET utf8 DEFAULT COLLATE utf8_general_ci`
	db.Exec(qstr)
}

func createStore() *mysqlStore {
	return NewStore(db, "session_store", "sid", "seed", "data", "ttl", "expire").(*mysqlStore)
}

func TestAllocate(t *testing.T) {
	s := createStore()
	seed := "ineedaseedbutidontknowwhatwillbe"

	_, err := s.Allocate(seed)
	if err != nil {
		t.Fatalf("error occured allocating storage: %s", err)
	}
}
