package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os/exec"

	"github.com/RestinGreen/polygon-arbitrage/pkg/connection"
	"github.com/RestinGreen/polygon-arbitrage/pkg/database"
	"github.com/RestinGreen/polygon-arbitrage/pkg/general"
	"github.com/RestinGreen/polygon-arbitrage/pkg/memory"
	"github.com/RestinGreen/polygon-arbitrage/pkg/peek"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
)

func init() {
	fmt.Println("Bot booted")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			exec.Command("stty", "-F", "/dev/tty", "echo").Run()
		}
	}()

	general := general.NewGeneral()
	fullNode := connection.NewConnection(general.IpcEndpoint)
	archiveNode := connection.NewConnection(general.InfuraEndpoint)

	db := database.NewDB(general)

	subscriber := gethclient.New(fullNode.RpcClient)

	newFullTx := make(chan *types.Transaction)
	dexMonitorCh := make(chan *types.Transaction)

	_, err := subscriber.SubscribePendingFullTransactions(context.Background(), newFullTx)
	if err != nil {
		fmt.Println(err)
	}

	newBlocks := make(chan *types.Header)
	fullNode.EthClient.SubscribeNewHead(context.Background(), newBlocks)

	dexMonitor := memory.NewDexMonirot(dexMonitorCh, fullNode.EthClient, archiveNode.EthClient, db)
	go dexMonitor.Start()
	peek := peek.NewPeek()
	peek.SetDexMemory(&dexMonitor.DexMemory)
	peek.StartPeek()

	for tx := range newFullTx {
		if len(tx.Data()) < 4 {
			continue
		}
		_, exists := general.UniV2Methods[hex.EncodeToString(tx.Data()[0:4])]
		if !exists {
			// fmt.Println(err)
			continue
		}
		dexMonitorCh <- tx
		// fmt.Println(tx.To(), method, tx.Hash().Hex())
	}

}
