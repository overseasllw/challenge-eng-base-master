package common

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	DB *sql.DB
)

func init() {
	DB = getDb()
}

func getDb() *sql.DB {
	log.Println("connecting to DB......")
	db, err := sql.Open("mysql", "root:testpass@tcp(db:3306)/challenge?parseTime=true")
	if err != nil {
		log.Panicln(err)
	}
	err = db.Ping()
	if err != nil {
		log.Panicln(err)
	}
	log.Println("connected!")
	return db
}
