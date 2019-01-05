/*
 * File share api
 *
 * API version: 2.0.0
 * Contact: enponsba@gmail.com
 */

package raspi

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var db *sqlx.DB

func NewDbConnection() *sqlx.DB {
	var err error

	db, err = sqlx.Connect("sqlite3", "__data.db")

	if err != nil {
		log.Fatalln(err)
	}

	log.Printf("Database connected")

	return db
}

func SeedDb(searchDir string) {
	var schema = `
	DROP TABLE IF EXISTS file;
	CREATE TABLE file (
		Id INTEGER PRIMARY KEY,
		Alias VARCHAR DEFAULT '',
		Name VARCHAR DEFAULT '',
		Route VARCHAR DEFAULT '',
		Type VARCHAR DEFAULT ''
	);
	`
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	sqlStr := "INSERT INTO file (Alias, Route, Name, Type) VALUES "
	vals := []interface{}{}

	e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {

			sqlStr += "(?, ?, ?, ?),"
			alias := strings.TrimRight(filepath.Base(path), filepath.Ext(path))
			route := path
			name := filepath.Base(path)
			type_ := "film"
			vals = append(vals, alias, route, name, type_)
		}
		return err
	})

	if e != nil {
		log.Fatalln(e.Error())
	}

	stmt, _ := db.Prepare(sqlStr[:len(sqlStr)-1])

	stmt.Exec(vals...)
}
