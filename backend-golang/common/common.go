package common

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var (
	Config config
	DB     *sql.DB
)

type config struct {
	BaseURL   string
	Env       string
	JwtSecret []byte
	db        string
}

func init() {
	Config = config{
		BaseURL:   "",                                                                        //os.Getenv("BASEURL"),
		Env:       "dev",                                                                     // os.Getenv("ENV"),
		JwtSecret: []byte(`D3YRZ|X!!R74+6n%B@(VAv&%1|(?QfsLatVxfEpTsa_?(2\*X+*u&}rfY}oh;3Z`), //os.Getenv("JWT_SECRET")
		db:        "root:testpass@tcp(db:3306)/challenge?parseTime=true",
	}

	DB = getDb()
}

func getDb() *sql.DB {
	log.Println("connecting to DB......")
	db, err := sql.Open("mysql", Config.db)
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
