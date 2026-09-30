package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT '',
    title TEXT NOT NULL DEFAULT '',
    comment TEXT NOT NULL DEFAULT '',
    repeat VARCHAR(128) NOT NULL DEFAULT '' 
);

CREATE INDEX scheduler_date ON scheduler (date);
`

var db *sql.DB

func Init(dbFile string) error {
	var install bool
	if _, err := os.Stat(dbFile); err != nil {
		install = true
	}
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}
	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
