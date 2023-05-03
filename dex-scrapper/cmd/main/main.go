package main

import (
	"context"
	"fmt"
	"os"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/bot"
	dbclient "github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/db-client"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/peek"
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
	ipcPath := os.Getenv("IPC_PATH")

	grpcClient := dbclient.NewClient(host, port)

	fullNode := blockchain.NewConnection(ipcPath)

	newFullTx := make(chan *types.Transaction)

	subscriber := gethclient.New(fullNode.RpcClient)
	_, err := subscriber.SubscribeFullPendingTransactions(context.Background(), newFullTx)
	if err != nil {
		log.Error("Failed to subscribe to full pending transaction.")
		panic(err)
	}

	scrapper := bot.NewScrapper(fullNode, grpcClient)
	peek := peek.NewPeek(scrapper.Memory)
	peek.StartPeek()
	scrapper.LoadFromDb()
	scrapper.StartScrapper()

}
