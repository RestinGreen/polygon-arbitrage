package bot

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain/binding"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain/binding/erc20"
	dbclient "github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/db-client"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/memory"
	t "github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient/gethclient"
	"github.com/ethereum/go-ethereum/log"
)

type Scrapper struct {
	conn          *blockchain.Connection
	binding       *binding.Binding
	read          *bind.CallOpts
	Memory        *memory.Memory
	syncWatchList []common.Address

	grpc *dbclient.GRPCClient

	// keeps track of already loaded dex's
	// true - currently reading data from chain
	// false - dex is already loaded
	// nil - need to load dex
	// key is router address
	dexTracker   map[common.Address]bool
	dexTrackerMu *sync.Mutex

	txChan       chan *types.Transaction
	univ2Methods *blockchain.UniV2Selector
}

func NewScrapper(conn *blockchain.Connection, grpc *dbclient.GRPCClient) *Scrapper {

	return &Scrapper{
		conn:    conn,
		binding: binding.NewBinding(conn.EthClient),
		read:    &bind.CallOpts{Pending: false},
		Memory:  memory.NewMemory(),

		grpc: grpc,

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
			fmt.Println("new key: ", *routerAddress)
			s.dexTrackerMu.Lock()
			s.dexTracker[*routerAddress] = true
			s.dexTrackerMu.Unlock()
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

	tmpPairs := make([]*t.Pair, 0)
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
		tmpPairs = append(tmpPairs, pairData)

	}
	fmt.Println("Loading", numPairs, "pairs for factory", factoryAddress, "finished")
	s.dexTracker[*routerAddress] = false
	s.grpc.InsertDex(s.Memory.DexMemory.Dexs[*routerAddress], tmpPairs)
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
	newToken0 := s.addTokenToMemory(&token0)
	if newToken0 == nil {
		return nil, false
	}

	token1, err := pairContract.Token1(s.read)
	if err != nil {
		fmt.Println("Failed to read token1 from pair", pairAddress.Hex())
		fmt.Println(err)
	}
	newToken1 := s.addTokenToMemory(&token1)
	if newToken1 == nil {
		return nil, false
	}

	lastTimestamp := int64(reserves.BlockTimestampLast)

	return &t.Pair{
		PairAddress: pairAddress,
		Token0:      newToken0,
		Token1:      newToken1,
		Reserve0:    reserves.Reserve0,
		Reserve1:    reserves.Reserve1,
		LastUpdated: &lastTimestamp,
	}, true
}

func (s *Scrapper) addTokenToMemory(tokenAddress *common.Address) *t.Token {
	tokenContract, err := erc20.NewERC20(*tokenAddress, s.conn.EthClient)
	if err != nil {
		fmt.Println("Failed to create ERC220 contract binding.")
		fmt.Println(err)
		return nil
	}
	name, err := tokenContract.Name(s.read)
	if err != nil {
		fmt.Println("Failed to get token name.")
		fmt.Println(err)
		return nil
	}
	symbol, err := tokenContract.Symbol(s.read)
	if err != nil {
		fmt.Println("Failed to get token symbol.")
		fmt.Println(err)
		return nil
	}
	decimals, err := tokenContract.Decimals(s.read)
	if err != nil {
		fmt.Println("Failed to get token decimals.")
		fmt.Println(err)
		return nil
	}
	dec := int64(decimals)
	newToken := &t.Token{
		Address: tokenAddress,
		Symbol:  &symbol,
		Name:    &name,
		Decimal: &dec,
	}
	s.Memory.AddTokenStruct(newToken)
	return newToken
}

func (s *Scrapper) LoadFromDb() {

	fmt.Println("Loading from database started")

	dexs := s.grpc.GetAllDex()

	for _, dex := range dexs {
		dexAddress := common.HexToAddress(dex.RouterAddress)
		factoryAddress := common.HexToAddress(dex.FactoryAddress)
		routerAddress := common.HexToAddress(dex.RouterAddress)
		s.dexTracker[routerAddress] = false
		s.Memory.DexMemory.Dexs[dexAddress] = &t.Dex{
			Factory:  &factoryAddress,
			Router:   &routerAddress,
			NumPairs: &dex.NumPairs,
		}
		for _, pair := range dex.Pairs {
			key := dex.RouterAddress + pair.Token0.Address + pair.Token1.Address
			pairAddress := common.HexToAddress(pair.Address)
			token0Address := common.HexToAddress(pair.Token0.Address)
			token1Address := common.HexToAddress(pair.Token1.Address)
			token0 := &t.Token{
				Address: &token0Address,
				Name:    &pair.Token0.Name,
				Symbol:  &pair.Token0.Symbol,
				Decimal: &pair.Token0.Decimal,
			}
			token1 := &t.Token{
				Address: &token1Address,
				Name:    &pair.Token1.Name,
				Symbol:  &pair.Token1.Symbol,
				Decimal: &pair.Token1.Decimal,
			}
			s.Memory.PairMemory.Pairs[key] = &t.Pair{
				PairAddress:   &pairAddress,
				RouterAddress: &routerAddress,
				Reserve0:      new(big.Int).SetBytes(pair.Reserve0),
				Reserve1:      new(big.Int).SetBytes(pair.Reserve1),
				LastUpdated:   &pair.LastUpdated,
				Token0:        token0,
				Token1:        token1,
			}
			s.Memory.PairMemory.PairMap[pairAddress] = &key
			s.Memory.TokenMemory.Tokens[token0Address] = token0
			s.Memory.TokenMemory.Tokens[token1Address] = token1
		}
	}

	fmt.Println("Loading from database finished")
}

func (s *Scrapper) ListenPairSyncEvents() {

	for _, pair := range s.Memory.PairMemory.Pairs {
		s.syncWatchList = append(s.syncWatchList, *pair.PairAddress)
	}
	fmt.Println("There are", len(s.syncWatchList), "pairs to watch for reserve updates.")

	query := ethereum.FilterQuery{
		Addresses: s.syncWatchList,
		Topics:    [][]common.Hash{{common.HexToHash("0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1")}},
	}

	logsCh := make(chan types.Log)

	subscribtion, err := s.conn.EthClient.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		fmt.Println("Error subscribing to events.")
		panic(err)
	}

	for {
		select {
		case err := <-subscribtion.Err():
			panic(err)
		case log := <-logsCh:
			b, err := s.conn.EthClient.BlockByHash(context.Background(), log.BlockHash)
			if err != nil {
				fmt.Println("Failed to get block by hash")
			}
			r0 := new(big.Int).SetBytes(log.Data[0:32])
			r1 := new(big.Int).SetBytes(log.Data[32:64])
			timestamp := int64(b.Time())

			key := s.Memory.PairMemory.PairMap[log.Address]
			s.Memory.PairMemory.Pairs[*key].Reserve0 = r0
			s.Memory.PairMemory.Pairs[*key].Reserve1 = r1
			s.Memory.PairMemory.Pairs[*key].LastUpdated = &timestamp
			
			go s.grpc.UpdatePair(log.Address.Hex(), r0.Bytes(), r1.Bytes(), &timestamp)
			fmt.Println("Pair", log.Address, "updated in block", log.BlockNumber)
		}
	}
}

func (s *Scrapper) UpdateExistingDexData() {

}