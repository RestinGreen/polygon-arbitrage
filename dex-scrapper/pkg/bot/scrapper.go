package bot

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain/binding"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/memory"
	t "github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/types"
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
	Memory  *memory.Memory

	// keeps track of already loaded dex's
	// true - reading data from chain
	// false - dex is already loaded
	// nil - need to load dex
	// key is router address
	dexTracker   map[common.Address]bool
	dexTrackerMu *sync.Mutex

	txChan       chan *types.Transaction
	univ2Methods *blockchain.UniV2Selector
}

func NewScrapper(conn *blockchain.Connection) *Scrapper {

	return &Scrapper{
		conn:    conn,
		binding: binding.NewBinding(conn.EthClient),
		read:    &bind.CallOpts{Pending: false},
		Memory:  memory.NewMemory(),

		dexTracker:   make(map[common.Address]bool),
		dexTrackerMu: &sync.Mutex{},

		txChan:       make(chan *types.Transaction),
		univ2Methods: blockchain.NewUniV2Selector(),
	}
}

func (s *Scrapper) StartScrapper() {

	fmt.Println("Scrapper bot booted.")
	subscriber := gethclient.New(s.conn.RpcClient)

	_, err := subscriber.SubscribeFullPendingTransactions(context.Background(), s.txChan)
	if err != nil {
		log.Error("Failed to subscribe to full pending transaction.")
		panic(err)
	}

	for tx := range s.txChan {

		if len(tx.Data()) < 4 || !s.univ2Methods.IsUniV2(tx.Data()) {
			continue
		}
		routerAddress := tx.To()
		if routerAddress == nil {
			log.Error(tx.Hash().Hex(), "to address is nil")
			continue
		}
		s.dexTrackerMu.Lock()
		isLoading, exists := s.dexTracker[*routerAddress]
		s.dexTrackerMu.Unlock()
		if exists || isLoading {
			continue
		}
		if !exists {
			s.dexTracker[*routerAddress] = true
			go s.loadNewDexDataFromChain(routerAddress)
		}

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
		fmt.Println(err)
		return
	}
	factoryContract, exists := s.binding.Factories[factoryAddress]
	if !exists {
		factoryContract = s.binding.AddFactoryContract(&factoryAddress)
	}

	numPairsBn, err := factoryContract.AllPairsLength(s.read)
	if err != nil {
		fmt.Printf("Failed to read number of pairs from factory %s\n", factoryAddress)
		fmt.Println(err)
		return
	}

	numPairs := numPairsBn.Int64()
	fmt.Println("Loading", numPairs, "pairs from factory", factoryAddress)
	s.Memory.AddDex(routerAddress, &factoryAddress, &numPairs)
	monthAgo1 := time.Now().AddDate(0, -1, 0)
	for i := int64(0); i < numPairs; i++ {
		pairAddress, err := factoryContract.AllPairs(s.read, big.NewInt(i))
		if err != nil {
			fmt.Printf("Failed to get pair with index %d from  factory %s\n", i, factoryAddress)
			fmt.Println(err)
		}
		pairData, ok := s.getPairData(&pairAddress, monthAgo1)
		if !ok {
			continue
		}
		pairData.RouterAddress = routerAddress
		s.Memory.AddPairStruct(pairData)

	}
	fmt.Println("Loading", numPairs, "pairs for factory", factoryAddress, "finished")
	s.dexTracker[*routerAddress] = false
}

func (s *Scrapper) getPairData(pairAddress *common.Address, monthAgo1 time.Time) (*t.Pair, bool) {
	pairContract := s.binding.AddPairContract(pairAddress)
	reserves, err := pairContract.GetReserves(s.read)
	if err != nil {
		fmt.Println("Failed to get reserves.")
		fmt.Println(err)
	}
	
	if reserves.Reserve0.Uint64() == 0 || reserves.Reserve1.Uint64() == 0 || time.Unix(int64(reserves.BlockTimestampLast), 0).Before(monthAgo1) {
		return nil, false
	}
	token0, err := pairContract.Token0(s.read)
	if err != nil {
		fmt.Println("Failed to read token0 from pair", pairAddress.Hex())
		fmt.Println(err)
	}
	token1, err := pairContract.Token1(s.read)
	if err != nil {
		fmt.Println("Failed to read token1 from pair", pairAddress.Hex())
		fmt.Println(err)
	}
	return &t.Pair{
		PairAddress:   pairAddress,
		Token0Address: &token0,
		Token1Address: &token1,
		Reserve0:      reserves.Reserve0,
		Reserve1:      reserves.Reserve1,
		LastUpdated:   &reserves.BlockTimestampLast,
	}, true
}

func (s *Scrapper) LoadFromDb() {

	fmt.Println("Loading from database started")

	fmt.Println("Loading from database finished")
}
