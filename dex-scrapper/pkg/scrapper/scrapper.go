package scrapper

import (
	"context"
	"fmt"
	"math/big"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain/binding"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/log"
)

type Scrapper struct {
	conn    *blockchain.Connection
	binding *binding.Binding
	read    *bind.CallOpts

	txChan    chan *types.Transaction
	blockChan chan *types.Header
}

func NewScrapper(conn *blockchain.Connection) *Scrapper {

	return &Scrapper{
		conn:    conn,
		binding: binding.NewBinding(conn.EthClient),
		read:    &bind.CallOpts{Pending: false},

		txChan: make(chan *types.Transaction),
	}
}

func (s *Scrapper) StartScrapper() {

	subscriber := gethclient.New(s.conn.RpcClient)

	_, err := subscriber.SubscribeFullPendingTransactions(context.Background(), s.txChan)
	if err != nil {
		log.Error("Failed to subscribe to full pending transaction.")
		panic(err)
	}

	for tx := range s.txChan {

		routerAddress := tx.To()
		if routerAddress == nil {
			log.Error(tx.Hash().Hex(), "to address is nil")
			continue
		}

		go s.loadNewDexDataFromChain(routerAddress)
	}

}

func (s *Scrapper) loadNewDexDataFromChain(routerAddress *common.Address) {
	routerContract, exists := s.binding.Routers[*routerAddress]
	if !exists {
		routerContract = s.binding.AddRouterContract(routerAddress)
	}
	factoryAddress, err := routerContract.Factory(s.read)
	if err != nil {
		fmt.Println("Failed to read factory address from router", routerAddress)
		return
	}
	factoryContract, exists := s.binding.Factories[factoryAddress]
	if !exists {
		factoryContract = s.binding.AddFactoryContract(&factoryAddress)
	}

	numPairsBn, err := factoryContract.AllPairsLength(s.read)
	if err != nil {
		fmt.Sprintf("Failed to read number of pairs from factory %s", factoryAddress)
		return
	}

	numPairs := numPairsBn.Int64()

	for i := int64(0); i < numPairs; i++ {

		pairAddress, err := factoryContract.AllPairs(s.read, big.NewInt(i))
		if err != nil {
			fmt.Println("Failed to get pair with index %d from  factory %s", i, factoryAddress)
		}
		pairData, ok := s.getPairData(&pairAddress)
	}
}

func (s *Scrapper) getPairData(pairAddress *common.Address) {
	
}
