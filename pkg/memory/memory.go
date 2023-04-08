package memory

import (
	"fmt"
	"sync"

	"github.com/RestinGreen/polygon-arbitrage/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type Memory struct {
	//key is router
	DexMemory  *DexMemory
	PairMemory *PairMemory

	// from -> to = pair route list
	Routes map[common.Address]map[common.Address][]*types.Route
	hop0Mu *sync.Mutex

	CreationMutex *sync.Mutex
}

func NewMemory() *Memory {

	return &Memory{
		DexMemory: &DexMemory{
			Dexs:     map[common.Address]*types.Dex{},
			DexMutex: map[common.Address]*sync.Mutex{},
		},
		PairMemory: &PairMemory{
			PairMap:   map[common.Address]*string{},
			Pairs:     map[string]*types.Pair{},
			PairMutex: map[string]*sync.Mutex{},
		},
		CreationMutex: &sync.Mutex{},
	}
}

func (m *Memory) AddDexStruct(dex *types.Dex) {
	if _, exists := m.DexMemory.DexMutex[*dex.Router]; !exists {
		m.DexMemory.DexMutex[*dex.Router] = &sync.Mutex{}
	}
	if _, exists := m.DexMemory.Dexs[*dex.Router]; !exists {

		m.DexMemory.Dexs[*dex.Router] = dex
		fmt.Println("Dex added to memory.")
		fmt.Println("\tpairs:\t\t", *dex.NumPairs)
		fmt.Println("\trouter:\t\t", dex.Router)
		fmt.Println("\tfactory:\t", dex.Factory)
	}

}
func (m *Memory) AddDex(router, factory *common.Address, numPairs *int) {
	if _, exists := m.DexMemory.DexMutex[*router]; !exists {
		m.DexMemory.DexMutex[*router] = &sync.Mutex{}
	}
	if _, exists := m.DexMemory.Dexs[*router]; !exists {
		m.DexMemory.Dexs[*router] = &types.Dex{
			Factory:  factory,
			Router:   router,
			NumPairs: numPairs,
		}
		fmt.Println("Dex added to memory.")
		fmt.Println("\tpairs:\t\t", *numPairs)
		fmt.Println("\trouter:\t\t", router)
		fmt.Println("\tfactory:\t", factory)
	}
}

func (m *Memory) AddPairStruct(pair *types.Pair) {

	m.CreationMutex.Lock()
	defer m.CreationMutex.Unlock()
	//sorted, token0 + token1 + router is the key
	key := pair.Token0Address.Hex() + pair.Token1Address.Hex() + pair.RouterAddress.Hex()

	if _, exists := m.PairMemory.PairMutex[key]; !exists {
		m.PairMemory.PairMutex[key] = &sync.Mutex{}
	}
	if _, exists := m.PairMemory.Pairs[key]; !exists {
		m.PairMemory.Pairs[key] = pair
	}

	if _, exists := m.PairMemory.PairMap[*pair.PairAddress]; !exists {
		m.PairMemory.PairMap[*pair.PairAddress] = &key
	}
}
