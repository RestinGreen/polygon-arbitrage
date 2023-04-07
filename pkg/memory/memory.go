package memory

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/RestinGreen/polygon-arbitrage/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type Memory struct {
	//key is router
	DexMemory  *DexMemory
	PairMemory *PairMemory
}

func NewMemory() *Memory {

	return &Memory{
		DexMemory: &DexMemory{
			Dexs:     map[common.Address]*types.Dex{},
			DexMutex: map[common.Address]*sync.Mutex{},
		},
		PairMemory: &PairMemory{
			PairMap:   map[common.Address]*types.SimplePair{},
			Pairs:     map[string]*types.Pair{},
			PairMutex: map[string]*sync.Mutex{},
		},
	}
}

func (m *Memory) AddDexStruct(dex *types.Dex) {
	if _, exists := m.DexMemory.DexMutex[*dex.Router]; !exists {
		m.DexMemory.DexMutex[*dex.Router] = &sync.Mutex{}
	}
	if _, exists := m.DexMemory.Dexs[*dex.Router]; !exists {

		m.DexMemory.Dexs[*dex.Router] = dex
		fmt.Println("Dex added to memory.")
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
		fmt.Println("\trouter:\t\t", router)
		fmt.Println("\tfactory:\t", factory)
	}
}

func (m *Memory) AddPairStruct(pair *types.Pair) {

	//sorted, token0 + token1 is the key
	key := pair.Token0Address.Hex() + pair.Token1Address.Hex()
	if _, exists := m.PairMemory.Pairs[key]; !exists {
		m.PairMemory.Pairs[key] = pair
	}
	if _, exists := m.PairMemory.PairMutex[key]; !exists {
		m.PairMemory.PairMutex[key] = &sync.Mutex{}
	}
	if _, exists := m.PairMemory.PairMap[*pair.PairAddress]; !exists {
		m.PairMemory.PairMap[*pair.PairAddress] = &types.SimplePair{
			Token0Address: pair.Token0Address,
			Token1Address: pair.Token1Address,
		}
	}
}

func (m *Memory) AddPair(router, pairAddress, token0, token1 *common.Address, reserve0, reserve1 *big.Int, lastupdated *uint32) {

	_, exists := m.PairMemory.Pairs[token0.Hash().Hex()+token1.Hash().Hex()]
	if !exists {
		m.PairMemory.Pairs[token0.Hash().Hex()+token1.Hash().Hex()] = &types.Pair{
			PairAddress:   pairAddress,
			Token0Address: token0,
			Token1Address: token1,
			Reserve0:      reserve0,
			Reserve1:      reserve1,
			LastUpdated:   lastupdated,
		}
		fmt.Println("Inserted to", router, "the pair", pairAddress)
	}

}
