/*
 * File share api
 *
 * API version: 2.0.0
 * Contact: enponsba@gmail.com
 */

package raspi

import (
	"fmt"
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

func SeedDb() error {
	var searchDir = os.Getenv("assetsDir")

	var schema = `
	DROP TABLE IF EXISTS file;
	CREATE TABLE file (
		id INTEGER PRIMARY KEY,
		title VARCHAR DEFAULT '',
		filename VARCHAR DEFAULT '',
		route VARCHAR DEFAULT '',
		type VARCHAR DEFAULT '',
		summary VARCHAR DEFAULT '',
		image VARCHAR DEFAULT ''
	);
	`
	// exec the schema or fail; multi-statement Exec behavior varies between
	// database drivers;  pq will exec them all, sqlite3 won't, ymmv
	db.MustExec(schema)

	if _, err := os.Stat(fmt.Sprintf("./%s", searchDir)); !os.IsNotExist(err) {
		sqlStr := "INSERT INTO file (title, route, filename, type) VALUES "
		vals := []interface{}{}

		e := filepath.Walk(searchDir, func(path string, f os.FileInfo, err error) error {
			if !f.IsDir() {

				sqlStr += "(?, ?, ?, ?),"
				title := strings.TrimRight(filepath.Base(path), filepath.Ext(path))
				route := path
				filename := filepath.Base(path)
				type_ := "film"
				vals = append(vals, title, route, filename, type_)
			}
			return err
		})

		if e != nil {
			log.Fatalln(e.Error())
		}

		stmt, _ := db.Prepare(sqlStr[:len(sqlStr)-1])

		stmt.Exec(vals...)
	}
	return nil
}
