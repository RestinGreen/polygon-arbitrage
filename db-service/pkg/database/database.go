package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDB(host, port, user, psw, dbName string) *Database {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, psw, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Login credentials are valid.")

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("connected")

	return &Database{db}
}
