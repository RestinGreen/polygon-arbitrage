package memory

import (
	"context"
	"fmt"
	"math/big"
	"sync"

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
	startSync *sync.Mutex
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

		startSync: &sync.Mutex{},
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

			//decoding pairs from tx data and adding to memory
			// switch hex.EncodeToString(txInputData[0:4]) {
			// case "fb3bdb41", "7ff36ab5", "b6f9de95":
			// 	m.decodeInputData(*routerAddress, factoryContract, txInputData, 4)

			// case "18cbafe5", "791ac947", "38ed1739", "5c11d795", "4a25d94a", "8803dbee":
			// 	m.decodeInputData(*routerAddress, factoryContract, txInputData, 5)
			// }
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
		pairData := m.getPairData(&pairAddress)
		pairData.RouterAddress = routerAddress
		m.Memory.AddPairStruct(pairData)

	}
	fmt.Println("Loading", numPairs, "pairs for factory", factoryAddress, "finished")
	//save to db
	m.db.InsertFullDex(m.Memory.DexMemory.Dexs[*routerAddress], m.Memory.PairMemory.Pairs, m.Memory.CreationMutex)
}

func (m *DataMonitor) getPairData(pairAddress *common.Address) *mem.Pair {

	pairContract := m.binding.AddPairContract(pairAddress)
	reserves, err := pairContract.GetReserves(m.read)
	if err != nil {
		fmt.Println("Failed to get reserves.")
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
	}

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

	// m.eventWatchlist = make([]common.Address, 0)

	for _, pair := range m.Memory.PairMemory.Pairs {
		// key :=
		m.eventWatchlist = append(m.eventWatchlist, *pair.PairAddress)
	}

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
			r0 := new(big.Int).SetBytes(log.Data[0:32])
			r1 := new(big.Int).SetBytes(log.Data[32:64])
			bn := uint32(log.BlockNumber)

			key := m.Memory.PairMemory.PairMap[log.Address]
			m.Memory.PairMemory.PairMutex[*key].Lock()
			m.Memory.PairMemory.Pairs[*key].Reserve0 = r0
			m.Memory.PairMemory.Pairs[*key].Reserve1 = r1
			m.Memory.PairMemory.Pairs[*key].LastUpdated = &bn
			m.Memory.PairMemory.PairMutex[*key].Unlock()
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
				fmt.Println("Updating", *pair.LastUpdated, " -> ", newPairData.BlockTimestampLast)
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

// func (m *DexMonitor) decodeInputData(router common.Address, factoryContract *univ2factory.UniV2Factory, data []byte, pathLengthIndex int) {

// 	length := int(new(big.Int).SetBytes(getithParam(data, pathLengthIndex)).Int64())
// 	for i := 0; i < length-1; i++ {
// 		tokenAByte := getithParam(data, pathLengthIndex+1+i)
// 		tokenBByte := getithParam(data, pathLengthIndex+1+i+1)

// 		addressA := common.BytesToAddress(tokenAByte)
// 		addressB := common.BytesToAddress(tokenBByte)

// 		token0, token1 := chain.SortAddress(addressA, addressB)

// 		pairAddress, err := factoryContract.GetPair(m.read, token0, token1)
// 		m.binding.AddPair(pairAddress)
// 		if err != nil {
// 			fmt.Println("Failed to get pairAddress")
// 			continue
// 		}
// 		reserves, err := m.binding.Pairs[pairAddress].GetReserves(m.read)
// 		if err != nil {
// 			fmt.Println("Failed to get reserves.")
// 			continue
// 		}

// 		m.DexMemory.AddPair(router, pairAddress, token0, token1, reserves.Reserve0, reserves.Reserve1, reserves.BlockTimestampLast)

// 	}
// }

// Param order in bytecode:
// 0
// func getithParam(bytes []byte, i int) []byte {

// 	firstIndex := 4 + i*32
// 	lastIndex := 4 + (i+1)*32
// 	return bytes[firstIndex:lastIndex]
// }

// func getithTopicData(bytes []byte, i int) []byte {
// 	firstIndex := i * 32
// 	lastIndex := (i + 1) * 32
// 	return bytes[firstIndex:lastIndex]
// }
