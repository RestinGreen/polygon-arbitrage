package main

import (
	"database/sql"
	"fmt"
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
		panic(err)
	}
	fmt.Println("Login credentials are valid.")

	err = db.Ping()
	if err != nil {
		panic(err)
	}
	fmt.Println("connected")

	query := `SELECT factory_address, router_address, num_pairs FROM dexs`

	rows, err := db.Query(query)
	if err != nil {
		fmt.Println("Failed to do 'GetAllDexs' query")
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var factory, router string
		var numPairs int
		if err := rows.Scan(&factory, &router, &numPairs); err != nil {
			fmt.Println("Error in scanning rows in 'GetAlLDexs'.", err)
			panic(err)
		}
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error in reading 'GetAllDexs' rows.", err)
		return nil
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
