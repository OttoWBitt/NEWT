package db

import (
	"database/sql"
	"log"
)

var DB *sql.DB

//dataBase opens the data base connection to master
func OpenDataBase(dsn string) {

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}

}
