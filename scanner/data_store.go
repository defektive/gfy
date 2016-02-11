package scanner

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"log"
	"os"
	"path"
	"time"
)

const dbName = "gfy.sqlite"

type DataStore struct {
	db *sql.DB
}

func (ds *DataStore) Add(path, date, hash string) {
	tx, err := ds.db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	stmt, err := tx.Prepare("insert into photos(path, date, hash) values(?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	// Mon Jan 2 15:04:05 -0700 MST 2006
	const shortForm = "2006:01:02 15:04:05 -0700 MST"
	t, _ := time.Parse(shortForm, date+" -0700 MST")

	_, err = stmt.Exec(path, t.Format("2006-01-03 00:04:05"), hash)
	if err != nil {
		log.Fatal(err)
	}
	tx.Commit()
}

func (ds *DataStore) Close() {
	ds.db.Close()
}

func OpenDb(dir string) *DataStore {
	dbPath := path.Join(dir, dbName)
	os.Remove(dbPath)

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatal(err)
	}
	sqlStmt := `
  CREATE TABLE photos (
  	id	INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
  	path	TEXT NOT NULL,
  	date	datetime NOT NULL,
  	hash	TEXT NOT NULL UNIQUE
  );
  CREATE  INDEX photo_date ON photos (date ASC);
	`
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		panic(err)
	}

	ds := &DataStore{
		db: db,
	}
	return ds
	//
	//
	// rows, err := db.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	rows.Scan(&id, &name)
	// 	fmt.Println(id, name)
	// }
	//
	// stmt, err = db.Prepare("select name from foo where id = ?")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer stmt.Close()
	// var name string
	// err = stmt.QueryRow("3").Scan(&name)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(name)
	//
	// _, err = db.Exec("delete from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// _, err = db.Exec("insert into foo(id, name) values(1, 'foo'), (2, 'bar'), (3, 'baz')")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// rows, err = db.Query("select id, name from foo")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer rows.Close()
	// for rows.Next() {
	// 	var id int
	// 	var name string
	// 	rows.Scan(&id, &name)
	// 	fmt.Println(id, name)
	// }
}
