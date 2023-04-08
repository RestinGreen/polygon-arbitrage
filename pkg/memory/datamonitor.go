package memory

import (
	"context"
	"fmt"
	"math/big"
	"sync"
	"time"

	"github.com/RestinGreen/polygon-arbitrage/pkg/binding"
	"github.com/RestinGreen/polygon-arbitrage/pkg/database"
	mem "github.com/RestinGreen/polygon-arbitrage/pkg/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DataMonitor struct {
	Memory  *Memory
	txChan  chan *types.Transaction
	binding *binding.Binding
	read    *bind.CallOpts

	//router is the key
	startSync      *sync.Mutex
	monitorTracker map[common.Address]bool
	eventWatchlist []common.Address

	db          *database.Database
	fullNode    *ethclient.Client
	archiveNode *ethclient.Client
}

func NewDexMonirot(txChan chan *types.Transaction, fullNode *ethclient.Client, archiveNode *ethclient.Client, db *database.Database) *DataMonitor {

	return &DataMonitor{
		txChan:  txChan,
		Memory:  NewMemory(),
		binding: binding.NewBinding(fullNode),
		read:    &bind.CallOpts{Pending: false},

		startSync:      &sync.Mutex{},
		monitorTracker: map[common.Address]bool{},

		db:          db,
		fullNode:    fullNode,
		archiveNode: archiveNode,
	}

}

func (m *DataMonitor) Start() {

	m.loadFromDB()

	if len(m.Memory.DexMemory.Dexs) > 0 {
		go m.listenPairSyncEvents()
		go m.refreshDexData()

		m.getAllRoutes()
	}

	fmt.Println("Mempool monitoring started.")
	for {
		select {

		case tx := <-m.txChan:
			routerAddress := tx.To()

			m.startSync.Lock()
			monitor, exists := m.monitorTracker[*routerAddress]
			m.startSync.Unlock()
			if exists && monitor {
				continue
			}
			if !exists || (exists && !monitor) {
				m.monitorTracker[*routerAddress] = true

			} else {
				continue
			}

			if routerAddress == nil {
				fmt.Println("to address is nil")
				continue
			}
			if _, exists := m.Memory.DexMemory.Dexs[*routerAddress]; !exists {
				go m.loadNewDexDataFromChain(routerAddress)
			} else {
				//TODO get only new pair addresses
			}
			m.monitorTracker[*routerAddress] = false

		}
	}
}
func (m *DataMonitor) loadNewDexDataFromChain(routerAddress *common.Address) {
	ONE := new(big.Int).SetInt64(1)

	routerContract, exists := m.binding.Routers[*routerAddress]
	if !exists {
		routerContract = m.binding.AddRouterContract(routerAddress)
	}
	factoryAddress, err := routerContract.Factory(m.read)
	if err != nil {
		fmt.Println("Failed to read factory address from router ", routerAddress)
		return
	}

	factoryContract, exists := m.binding.Factories[factoryAddress]
	if !exists {
		factoryContract = m.binding.AddFactoryContract(&factoryAddress)
	}

	numPairsB, err := factoryContract.AllPairsLength(m.read)
	if err != nil {
		fmt.Println("Failed to read number of all pairs from factory ", factoryAddress)
		return
	}
	numPairs := int(numPairsB.Uint64())
	//add dex to memory
	m.Memory.AddDex(routerAddress, &factoryAddress, &numPairs)

	fmt.Println("Loading", numPairs, "pairs from factory", factoryAddress)

	cntB := new(big.Int).SetInt64(0)
	for i := 0; i < numPairs; i++ {

		pairAddress, err := factoryContract.AllPairs(m.read, cntB)
		cntB.Add(cntB, ONE)
		if err != nil {
			fmt.Println("Failed to get pair with index", i, "from factory", factoryAddress)
			continue
		}
		pairData, ok := m.getPairData(&pairAddress)
		if !ok {
			continue
		}
		pairData.RouterAddress = routerAddress
		m.Memory.AddPairStruct(pairData)

	}
	fmt.Println("Loading", numPairs, "pairs for factory", factoryAddress, "finished")
	//save to db
	m.db.InsertFullDex(m.Memory.DexMemory.Dexs[*routerAddress], m.Memory.PairMemory.Pairs, m.Memory.CreationMutex)
}

func (m *DataMonitor) getPairData(pairAddress *common.Address) (*mem.Pair, bool) {

	pairContract := m.binding.AddPairContract(pairAddress)
	reserves, err := pairContract.GetReserves(m.read)
	if err != nil {
		fmt.Println("Failed to get reserves.")
	}
	if reserves.Reserve0.Uint64() == 0 || reserves.Reserve1.Uint64() == 0 {
		return nil, false
	}
	token0, err := pairContract.Token0(m.read)
	if err != nil {
		fmt.Println("Failed to read token0 from pair", pairAddress.Hex())
	}
	token1, err := pairContract.Token1(m.read)
	if err != nil {
		fmt.Println("Failed to read token1 from pair", pairAddress.Hex())
	}
	return &mem.Pair{
		PairAddress:   pairAddress,
		Token0Address: &token0,
		Token1Address: &token1,
		Reserve0:      reserves.Reserve0,
		Reserve1:      reserves.Reserve1,
		LastUpdated:   &reserves.BlockTimestampLast,
	}, true

}

func (m *DataMonitor) loadFromDB() {

	fmt.Println("Loading database and creating bindings.")

	var wg sync.WaitGroup
	fmt.Println("Loading dexs.")
	for _, dex := range m.db.GetAllData() {
		m.Memory.AddDexStruct(dex)
		go func(dex *mem.Dex) {
			wg.Add(1)
			defer wg.Done()
			m.binding.AddFactoryContract(dex.Factory)
		}(dex)
		go func(dex *mem.Dex) {
			wg.Add(1)
			defer wg.Done()
			m.binding.AddRouterContract(dex.Router)
		}(dex)
	}
	fmt.Println("Loading dexs finished.", len(m.Memory.DexMemory.Dexs), "dexs loaded.")
	fmt.Println("Loading pairs.")
	for _, pair := range m.db.GetPairs() {
		m.Memory.AddPairStruct(pair)
		go func(pair *mem.Pair) {
			wg.Add(1)
			defer wg.Done()
			m.binding.AddPairContract(pair.PairAddress)
		}(pair)
	}
	fmt.Println("Loading pairs finished.", len(m.Memory.PairMemory.Pairs), "pairs loaded.")
	wg.Wait()
	fmt.Println("Database loaded.")
}

func (m *DataMonitor) listenPairSyncEvents() {

	m.Memory.CreationMutex.Lock()
	for _, pair := range m.Memory.PairMemory.Pairs {
		m.eventWatchlist = append(m.eventWatchlist, *pair.PairAddress)
	}
	m.Memory.CreationMutex.Unlock()

	fmt.Println("There are", len(m.eventWatchlist), "pairs to watch.")
	query := ethereum.FilterQuery{
		Addresses: m.eventWatchlist,
		Topics:    [][]common.Hash{{common.HexToHash("0x1c411e9a96e071241c2f21f7726b17ae89e3cab4c78be50e062b03a9fffbbad1")}},
	}

	logsCh := make(chan types.Log)

	subscribtion, err := m.fullNode.SubscribeFilterLogs(context.Background(), query, logsCh)
	if err != nil {
		fmt.Println("Error subscribing to events.")
		panic(err)
	}

	for {
		select {
		case err := <-subscribtion.Err():
			panic(err)
		case log := <-logsCh:
			b, err := m.fullNode.BlockByHash(context.Background(), log.BlockHash)
			if err != nil {
				fmt.Println("Failed to get block by hash")
			}
			r0 := new(big.Int).SetBytes(log.Data[0:32])
			r1 := new(big.Int).SetBytes(log.Data[32:64])
			timestamp := uint32(b.Time())

			key := m.Memory.PairMemory.PairMap[log.Address]
			m.Memory.PairMemory.PairMutex[*key].Lock()
			m.Memory.PairMemory.Pairs[*key].Reserve0 = r0
			m.Memory.PairMemory.Pairs[*key].Reserve1 = r1
			m.Memory.PairMemory.Pairs[*key].LastUpdated = &timestamp
			m.Memory.PairMemory.PairMutex[*key].Unlock()
			go m.db.UpdatePair(log.Address.Hex(), r0, r1, &timestamp)
			fmt.Println("Pair", log.Address, "updated in block", log.BlockNumber)
		}
	}

}

func (m *DataMonitor) refreshDexData() {
	fmt.Println("Refreshing started.")
	for _, pair := range m.Memory.PairMemory.Pairs {
		newPairData, err := m.binding.GetPairContract(*pair.PairAddress).GetReserves(m.read)
		if err != nil {
			fmt.Println("Failed to get reservers in refreshing.", err)
			continue
		}
		go func(pair *mem.Pair) {
			if *pair.LastUpdated < newPairData.BlockTimestampLast {
				fmt.Println("Updating", *pair.LastUpdated, " -> ", newPairData.BlockTimestampLast, "pair", *pair.PairAddress)
				m.Memory.PairMemory.PairMutex[getPairKey(pair)].Lock()
				pair.Reserve0 = newPairData.Reserve0
				pair.Reserve1 = newPairData.Reserve1
				pair.LastUpdated = &newPairData.BlockTimestampLast
				m.Memory.PairMemory.PairMutex[getPairKey(pair)].Unlock()
				go m.db.UpdatePair(pair.PairAddress.Hex(), pair.Reserve0, pair.Reserve1, pair.LastUpdated)
			}
		}(pair)
	}
	fmt.Println("Refreshing finished.")
}

func getPairKey(pair *mem.Pair) string {
	return pair.Token0Address.Hex() + pair.Token1Address.Hex() + pair.RouterAddress.Hex()
}

func (m *Memory) dbSaverJob() {

}

func (m *DataMonitor) getAllRoutes() {
	fmt.Println("Getting routes")

	m.Memory.Routes = make(map[common.Address]map[common.Address][]*mem.Route)
	i := 0
	m.Memory.hop0Mu = &sync.Mutex{}
	wg := &sync.WaitGroup{}

	tokens := map[common.Address]bool{}
	for _, x := range m.Memory.PairMemory.Pairs {
		tokens[*x.Token0Address] = true
		tokens[*x.Token1Address] = true
	}

	fmt.Println("nr of tokens ", len(tokens))

	a := 1
	for t1 := range tokens {
		fmt.Println(a)
		a++
		for t2 := range tokens {
			if _, exists := m.Memory.Routes[t1]; !exists {
				m.Memory.Routes[t1] = make(map[common.Address][]*mem.Route)
			}
			if _, exists := m.Memory.Routes[t1][t2]; !exists {
				m.Memory.Routes[t1][t2] = make([]*mem.Route, 0)
			}
		}
	}

	for _, pair := range m.Memory.PairMemory.Pairs {
		go func(a *common.Address) {
			wg.Add(1)
			defer wg.Done()
			m.find0HopRoutes(a)
		}(pair.Token0Address)
		i++
		if i == 1000 {
			break
		}
	}
	wg.Wait()
	fmt.Println("0 hop routes finished,")

}

func (m *DataMonitor) find0HopRoutes(start *common.Address) {

	sum := 0
	for _, pair := range m.Memory.PairMemory.Pairs {
		sstart := time.Now()
		route := &mem.Route{
			Path: []*mem.Node{},
		}
		if *start == *pair.Token0Address {
			node := &mem.Node{
				Pair:      pair,
				IsInverse: false,
			}
			route.Path = append(route.Path, node)
			// routes = append(routes, route)
		} else if *start == *pair.Token1Address {
			node := &mem.Node{
				Pair:      pair,
				IsInverse: true,
			}
			route.Path = append(route.Path, node)
			// routes = append(routes, route)
		} else {
			continue
		}
		m.Memory.hop0Mu.Lock()
		if _, exists := m.Memory.Routes[*start]; !exists {
			m.Memory.Routes[*start] = make(map[common.Address][]*mem.Route)
		}

		if _, exists := m.Memory.Routes[*start][*pair.Token0Address]; !exists {
			m.Memory.Routes[*start][*pair.Token0Address] = make([]*mem.Route, 0)
		}
		m.Memory.Routes[*start][*pair.Token0Address] = append(m.Memory.Routes[*start][*pair.Token0Address], route)
		m.Memory.hop0Mu.Unlock()
		since := int(time.Since(sstart))
		sum += since
	}
	fmt.Println("avg ", sum/len(m.Memory.PairMemory.Pairs))
}
