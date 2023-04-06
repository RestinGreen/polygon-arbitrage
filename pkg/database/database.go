package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	_ "github.com/lib/pq"

	"github.com/RestinGreen/polygon-arbitrage/pkg/general"
	"github.com/RestinGreen/polygon-arbitrage/pkg/types"
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
	// defer db.Close()

	return &Database{db}
}

func (db *Database) GetAllDexs() []*types.DexMemory {

	query := `SELECT factory_address, router_address, num_pairs
		FROM dexs`

	rows, err := db.db.Query(query)
	if err != nil {
		fmt.Println("Failed to do 'GetAllDexs' query", err)
		return nil
	}
	defer rows.Close()

	dexs := make([]*types.DexMemory, 0)
	for rows.Next() {
		var factory, router string
		var numPairs int
		if err := rows.Scan(&factory, &router, &numPairs); err != nil {
			fmt.Println("Error in scanning rows in 'GetAlLDexs'.", err)
			return nil
		}
		dexs = append(dexs, &types.DexMemory{
			Factory:    common.HexToAddress(factory),
			Router:     common.HexToAddress(router),
			NumPairs:   numPairs,
			SimplePair: map[common.Address]*types.SimplePair{},
			Pairs:      map[string]*types.Pair{},
			PairMutex:  map[string]*sync.Mutex{},
		})
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error in reading 'GetAllDexs' rows.", err)
		return nil
	}
	return dexs
}

func (db *Database) GetPairsForDex(factoryAddress string) []*types.Pair {
	pairs := make([]*types.Pair, 0)

	query := `SELECT p.pair_address, p.token0_address, p.token1_address, p.reserve0, p.reserve1, p.last_updated
		from dexs d 
		join pairs p
		on d.id = p.dex_id
		where d.factory_address = $1`

	rows, err := db.db.Query(query, factoryAddress)
	if err != nil {
		fmt.Println("Error in querying Pairs for Dex.", err)
		return nil
	}
	for rows.Next() {
		var pairAddress, token0Address, token1Address, reserve0, reserve1 string
		var lastUpdated uint32
		if err := rows.Scan(&pairAddress, &token0Address, &token1Address, &reserve0, &reserve1, &lastUpdated); err != nil {
			fmt.Println("Error in scanning rows in 'GetPair'.", err)
			return nil
		}
		r0, _ := new(big.Int).SetString(reserve0, 10)
		r1, _ := new(big.Int).SetString(reserve1, 10)
		pairs = append(pairs, &types.Pair{
			PairAddress:   common.HexToAddress(pairAddress),
			Token0Address: common.HexToAddress(token0Address),
			Token1Address: common.HexToAddress(token1Address),
			Reserve0:      r0,
			Reserve1:      r1,
			LastUpdated:   lastUpdated,
		})
	}
	if err := rows.Err(); err != nil {
		fmt.Println("Error in reading 'GetPair' rows.", err)
		return nil
	}
	return pairs

}

func (db *Database) GetPairsForFactory(factoryAddress string) []common.Address {
	return nil
}

func (db *Database) InsertPair(pairAddress, token0Address, token1Address string, reserve0, reserve1 *big.Int, lastUpdated uint32) error {
	query := `
			INSERT INTO mytable(pair_address, token0_address, token1_address, reserve0, reserve1, last_updated) 
			VALUES ($1, $2, $3, $4, $5, $6)
			`
	_, err := db.db.Exec(query, pairAddress, token0Address, token1Address, reserve0, reserve1, lastUpdated)
	if err != nil {
		return fmt.Errorf("failed to insert data: %v", err)
	}
	return nil
}

func (db *Database) UpdatePair(pairAddress string, reserve0 *big.Int, reserve1 *big.Int, lastUpdated uint32) {

	query := `
				UPDATE pairs 
				SET reserve0=$1, reserve1=$2, last_updated=$3
				where pair_address=$4
				`
	_, err := db.db.Exec(query, reserve0.String(), reserve1.String(), lastUpdated, pairAddress)
	if err != nil {
		fmt.Println("Error updating pair.", err)
		return
	}

}

func (db *Database) InsertFullDex(dex *types.DexMemory) {
	tx, err := db.db.Begin()
	if err != nil {
		fmt.Println("Failed to being transaction ", err)
		tx.Rollback()
		return
	}
	// Insert or update the Dex record
	var dexId int
	tx.QueryRow(`
        INSERT INTO dexs (factory_address, router_address, num_pairs)
        VALUES ($1, $2, $3)
		RETURNING id
    `, dex.Factory.Hex(), dex.Router.Hex(), dex.NumPairs).Scan(&dexId)

	// Insert or update the Pair records
	for _, pair := range dex.Pairs {
		_, err = tx.Exec(`
			INSERT INTO pairs (pair_address, token0_address, token1_address, reserve0, reserve1, last_updated, dex_id)
			VALUES ($1, $2, $3, $4, $5, $6, $7)`,
			pair.PairAddress.Hex(),
			pair.Token0Address.Hex(),
			pair.Token1Address.Hex(),
			pair.Reserve0.String(),
			pair.Reserve1.String(),
			pair.LastUpdated,
			dexId)
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
