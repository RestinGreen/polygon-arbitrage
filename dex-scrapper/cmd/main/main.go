package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain"
	dbclient "github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/db-client"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/log"
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

	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	grpcClient := dbclient.NewClient(host, port)
	grpcClient.GetAllDex()

	fullNode := blockchain.NewConnection("/ipc/bor.ipc")

	newFullTx := make(chan *types.Transaction)

	subscriber := gethclient.New(fullNode.RpcClient)
	_, err := subscriber.SubscribeFullPendingTransactions(context.Background(), newFullTx)
	if err != nil {
		log.Error("Failed to subscribe to full pending transaction.")
		panic(err)
	}

	for tx := range newFullTx {
		fmt.Println(tx.Gas())
	}
	

}
