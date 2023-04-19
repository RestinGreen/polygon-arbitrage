package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os/exec"

	"github.com/RestinGreen/polygon-arbitrage/tradebot/pkg/arbitrage"
	"github.com/RestinGreen/polygon-arbitrage/tradebot/pkg/connection"
	"github.com/RestinGreen/polygon-arbitrage/tradebot/pkg/database"
	"github.com/RestinGreen/polygon-arbitrage/tradebot/pkg/general"
	"github.com/RestinGreen/polygon-arbitrage/tradebot/pkg/memory"
	"github.com/RestinGreen/polygon-arbitrage/tradebot/pkg/peek"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
)

func init() {
	fmt.Println("Bot booted")
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("asd")
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
	arbitrageCh := make(chan *types.Transaction)

	defer close(dexMonitorCh)
	defer close(arbitrageCh)

	_, err := subscriber.SubscribePendingFullTransactions(context.Background(), newFullTx)
	if err != nil {
		fmt.Println(err)
	}

	newBlocks := make(chan *types.Header)
	fullNode.EthClient.SubscribeNewHead(context.Background(), newBlocks)

	monitor := memory.NewDexMonirot(dexMonitorCh, fullNode.EthClient, archiveNode.EthClient, db)
	arbitrage := arbitrage.NewArbitrage(arbitrageCh, monitor.Memory)
	go arbitrage.StartArbitrage()
	go monitor.Start()
	peek := peek.NewPeek(monitor.Memory, db)
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
		arbitrageCh <- tx
		// fmt.Println(tx.To(), method, tx.Hash().Hex())
	}

}
