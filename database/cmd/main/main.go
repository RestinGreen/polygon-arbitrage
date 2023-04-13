package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDB(host, port, user, psw, dbName string) *Database {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, psw, dbName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	return &Database{db}
}

func init() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed to load .env file.")
		panic(err)
	}

}

func main() {

	host := os.Getenv("PGSQL_HOST")
	port := os.Getenv("PGSQL_PORT")
	user := os.Getenv("PGSQL_USER")
	password := os.Getenv("PGSQL_PASSWORD")
	dbName := os.Getenv("PGSQL_DBNAME")

	NewDB(host, port, user, password, dbName)
}
