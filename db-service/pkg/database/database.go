package database

import (
	"database/sql"
	"fmt"
	"math/big"
	"time"

	"github.com/RestinGreen/protobuf/generated"
	_ "github.com/lib/pq"
)

type Database struct {
	db *sql.DB
}

func NewDB(host, port, user, psw, dbName string) *Database {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, psw, dbName)
	fmt.Println("commection string: ", connStr)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	fmt.Println("Login credentials are valid.")

	timeout := 30
	i := 0
	for i = 1; i <= timeout; {
		err = db.Ping()
		if err != nil {
			fmt.Println("waiting for pgsql server to start...")
			fmt.Println(err)
			i += 2
			time.Sleep(time.Second * 2)
		} else {
			break
		}

	}
	if i >= timeout {
		panic("pgsql server connection timeout")
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

func (db *Database) GetAllDex() ([]*generated.Dex, bool) {

	var dexs []*generated.Dex
	dexQuery := `
	SELECT id, factory_address, router_address, num_pairs FROM dexs;
	`
	dexRows, err := db.db.Query(dexQuery)
	if err != nil {
		fmt.Println("Failed to query dex rows in GetAllDex.\n", err)
		return nil, false
	}

	for dexRows.Next() {
		var dexId, numPairs int64
		var factoryAddress, routerAddress string
		if err := dexRows.Scan(&dexId, &factoryAddress, &routerAddress, &numPairs); err != nil {
			fmt.Println("Error scanning dex row in GetAlLDex\n", err)
			return nil, false
		}

		dex := &generated.Dex{
			FactoryAddress: factoryAddress,
			RouterAddress:  routerAddress,
			NumPairs:       numPairs,
			Pairs:          make([]*generated.Pair, 0),
		}

		pairQuery := `
		SELECT pair_address, reserve0, reserve1, last_updated, token0_id, token1_id 
		FROM pairs
		WHERE dex_id = $1;
		`
		pairRows, err := db.db.Query(pairQuery, dexId)
		if err != nil {
			fmt.Println("Failed to query pair rows in GetAllDex.\n", err)
		}

		for pairRows.Next() {

			pair := &generated.Pair{}

			var pairAddress, reserve0, reserve1 string
			var lastUpdated, token0ID, token1ID int64
			if err := pairRows.Scan(&pairAddress, &reserve0, &reserve1, &lastUpdated, &token0ID, &token1ID); err != nil {
				fmt.Println("Error scanning pair row in GetAlLDex\n", err)
				return nil, false
			}

			r0, _ := new(big.Int).SetString(reserve0, 10)
			r1, _ := new(big.Int).SetString(reserve1, 10)

			pair.Address = pairAddress
			pair.Reserve0 = r0.Bytes()
			pair.Reserve1 = r1.Bytes()
			pair.LastUpdated = lastUpdated

			token0Query := `
			SELECT address, symbol, name, decimal
			FROM tokens
			WHERE id=$1;
			`
			token0Rows, err := db.db.Query(token0Query, token0ID)
			if err != nil {
				fmt.Println("Failed to query token0 rows in GetAllDex.\n", err)
			}
			for token0Rows.Next() {
				var address, symbol, name string
				var decimal int64
				if err := token0Rows.Scan(&address, &symbol, &name, &decimal); err != nil {
					fmt.Println("Error scanning token0 row in GetAllDex.\n", err)
					return nil, false
				}
				pair.Token0 = &generated.Token{
					Address: address,
					Symbol:  symbol,
					Name:    name,
					Decimal: decimal,
				}

			}
			token1Query := `
			SELECT address, symbol, name, decimal
			FROM tokens
			WHERE id=$1
			`
			token1Rows, err := db.db.Query(token1Query, token1ID)
			if err != nil {
				fmt.Println("Failed to quert token1 rows in GetAllDex.\n", err)
			}
			for token1Rows.Next() {
				var address, symbol, name string
				var decimal int64
				if err := token1Rows.Scan(&address, &symbol, &name, &decimal); err != nil {
					fmt.Println("Error scanning token1 row in GetAllDex.\n", err)
					return nil, false
				}
				pair.Token1 = &generated.Token{
					Address: address,
					Symbol:  symbol,
					Name:    name,
					Decimal: decimal,
				}
			}

			dex.Pairs = append(dex.Pairs, pair)

		}
		dexs = append(dexs, dex)

	}
	fmt.Println("Finished loading all dex data from database.")
	return dexs, true
}

func (db *Database) UpdatePair(pairAddress string, reserve0, reserve1 []byte, lastUpdated *int64) {

	r0 := new(big.Int).SetBytes(reserve0)
	r1 := new(big.Int).SetBytes(reserve1)

	query := `
		UPDATE pairs 
		SET reserve0=$1, reserve1=$2, last_updated=$3
		where pair_address=$4
	`
	_, err := db.db.Exec(query, r0.String(), r1.String(), *lastUpdated, pairAddress)
	if err != nil {
		fmt.Println("Error updating pair.", err)
		return
	}
	fmt.Println("Pair ", pairAddress, "updated.")

}

func (db *Database) RemovePair(pairAddress string) {
	removeQuery := `
	DELETE FROM pairs
	WHERE pair_address like $1
	`
	_, err := db.db.Exec(removeQuery, pairAddress)
	if err != nil {
		fmt.Println("Failed to remove pair ", pairAddress, "from the database.")
		fmt.Println(err)
	}
}

func (db *Database) RemoveToken(tokenAddress string) {
	removeQuery := `
	DELETE FROM tokens
	WHERE address like $1
	`
	_, err := db.db.Exec(removeQuery, tokenAddress)
	if err != nil {
		fmt.Println("Failed to remove token ", tokenAddress, "from the database.")
		fmt.Println(err)
	}
}

