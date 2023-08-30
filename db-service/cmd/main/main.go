package main

import (
	"fmt"
	"os"

	"github.com/RestinGreen/polygon-arbitrage/db-service/pkg/database"
	"github.com/RestinGreen/polygon-arbitrage/db-service/pkg/server"
	"github.com/joho/godotenv"
)

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
	grpcHost := os.Getenv("GRPC_HOST")
	grpcPort := os.Getenv("GRPC_PORT")

	db := database.NewDB(host, port, user, password, dbName)
	server.NewServer(db, grpcHost, grpcPort)
}
