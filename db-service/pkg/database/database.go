package database

import (
	"database/sql"
	"fmt"
	"math/big"

	"github.com/RestinGreen/protobuf/generated"
	"github.com/ethereum/go-ethereum/common"
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
	fmt.Println("Connected to postgres database")

	return &Database{db}
}

func (db *Database) InsertDex(factoryAddress, routerAddress *string, numPairs *int64, pairs []*generated.Pair) {

	tx, err := db.db.Begin()
	if err != nil {
		fmt.Println("Failed to being transaction ", err)
		tx.Rollback()
		return
	}
	// Insert Dex records
	var dexId int
	tx.QueryRow(`
        INSERT INTO dexs (factory_address, router_address, num_pairs)
        VALUES ($1, $2, $3)
		RETURNING id
    `, *factoryAddress, *routerAddress, *numPairs).Scan(&dexId)

	// Insert Pair records
	for _, pair := range pairs {
		r0 := new(big.Int).SetBytes(pair.Reserve0)
		r1 := new(big.Int).SetBytes(pair.Reserve1)


		var token0Id, token1Id int

		tx.QueryRow(`
		INSERT INTO tokens (address, symbol, name, decimal)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (address) DO UPDATE SET id=tokens.id
		RETURNING id
		`, pair.Token0.Address, pair.Token0.Symbol, pair.Token0.Name, pair.Token0.Decimal).Scan(&token0Id)
		tx.QueryRow(`
		INSERT INTO tokens (address, symbol, name, decimal)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (address) DO UPDATE SET id=tokens.id
		RETURNING id
		`, pair.Token1.Address, pair.Token1.Symbol, pair.Token1.Name, pair.Token1.Decimal).Scan(&token1Id)

		fmt.Println(pair.Address)
		fmt.Println(r0.String())
		fmt.Println(r1.String())
		fmt.Println(pair.LastUpdated)
		fmt.Println(dexId)
		fmt.Println(token0Id)
		fmt.Println(token1Id)

		_, err = tx.Exec(`
						INSERT INTO pairs (
							pair_address,
							reserve0,
							reserve1,
							last_updated,
							dex_id,
							token0_id,
							token1_id)
						VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			pair.Address,
			r0.String(),
			r1.String(),
			pair.LastUpdated,
			dexId,
			token0Id,
			token1Id)
		if err != nil {
			fmt.Println("Error inserting into db.", err)
			tx.Rollback()
			return
		}
	}

	err = tx.Commit()
	if err != nil {
		fmt.Println("Failed to commit changed to the database", err)
		tx.Rollback()
		return
	}

}

func (db *Database) GetAllDex() {

	query := `
	SELECT p.pair_address, p.reserve0, p.reserve1, p.last_updated, t0.address, t0.symbol, t0.name, t0.decimal, t1.address, t1.symbol, t1.name, t1.decimal, d.factory_address, d.router_address, d.num_pairs 
	FROM pairs p 
	JOIN tokens t0 ON p.token0_id=t0.id 
	JOIN tokens t1 ON t1.id=p.token1_id 
	JOIN dexs d ON d.id = p.dex_id;
	`
	rows, err := db.db.Query(query)
	if err != nil {
		fmt.Println("Error in querying Pairs for Dex.", err)
	}

	var dex []*generated.Dex

	for rows.Next() {
		var pairAddress, reserve0, reserve1, t0Address, t1Address, t0Symbol, t1Symbol, t0Name, t1Name, factoryAddress, routerAddress string
		var lastUpdated, t0Decimal, t1Decimal, numPairs uint64
		if err := rows.Scan(&pairAddress, &reserve0, &reserve1, &lastUpdated, &t0Address, &t0Symbol, &t0Name, t0Decimal, &t1Address, &t1Symbol, &t1Name, &t1Decimal, &factoryAddress, &routerAddress, &numPairs); err != nil {
			fmt.Println("Error in scanning rows in 'GetPair'.", err)
			return
		}
		pa := common.HexToAddress(pairAddress)
		r0, _ := new(big.Int).SetString(reserve0, 10)
		r1, _ := new(big.Int).SetString(reserve1, 10)
		t0A := common.HexToAddress(t0Address)
		t1A := common.HexToAddress(t1Address)
		fA := common.HexToAddress(factoryAddress)
		rA := common.HexToAddress(routerAddress)
	

		
		dex = append(dex, &generated.Dex{
			
		})
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error in reading 'GetPair' rows.", err)
		return
	}
}