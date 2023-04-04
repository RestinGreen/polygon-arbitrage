package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"

	_ "github.com/lib/pq"

	"github.com/RestinGreen/polygon-arbitrage/pkg/general"
)

type Database struct {
	db *sql.DB
}

func NewDB(gen *general.General) *Database {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", gen.Host, gen.Port, gen.User, gen.Password, gen.DBName)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	return &Database{db}
}

func (db *Database) GetAllDexs()  {

	query := `SELECT dex.factory_address, dex.router_address, dex.number_of_pairs, pairs.pair_address, pairs.token0_address, pairs.token1_address, pairs.reserve0, pairs.reserve1, pairs.last_updated 
		FROM dex
  		JOIN pairs ON dex.factory_address = pairs.factory_address`
		
	rows, err := db.db.Query(query)
	if err != nil {
		fmt.Println(err)
	}
	defer rows.Close()
	for rows.Next() {
        var factoryAddress, routerAddress, pairAddress, token0Address, token1Address string
        if err := rows.Scan(&factoryAddress, &routerAddress, &pairAddress, &token0Address, &token1Address); err != nil {
            panic(err)
        }
        fmt.Printf("Dex: %s %s, Pair: %s %s %s\n", factoryAddress, routerAddress, pairAddress, token0Address, token1Address)
    }
    if err := rows.Err(); err != nil {
        panic(err)
    }

}

func (db *Database) InsertDex(rou)

func (db *Database) InsertPair(pairAddress, token0Address, token1Address string, reserve0, reserve1 *big.Int, lastUpdated uint32) error {
	query := `INSERT INTO mytable(pair_address, token0_address, token1_address, reserve0, reserve1, last_updated) 
				VALUES ($1, $2, $3, $4, $5, $6)`
	_, err := db.db.Exec(query, pairAddress, token0Address, token1Address, reserve0, reserve1, lastUpdated)
	if err != nil {
		return fmt.Errorf("failed to insert data: %v", err)
	}
	return nil
}
