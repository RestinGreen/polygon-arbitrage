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
	DexMemory map[common.Address]*types.DexMemory

	PairMemory *types.PairMemory
}

func NewMemory() *Memory {

	return &Memory{
		DexMemory: map[common.Address]*types.DexMemory{},
		PairMemory: &types.PairMemory{
			PairMap: map[common.Address]*types.SimplePair{},
			Pairs:   map[string]*types.Pair{},
		},
	}
}

func (m *Memory) AddDexStruct(dex *types.DexMemory) {
	_, exists := m.DexMemory[dex.Router]
	if !exists {
		m.DexMemory[dex.Router] = dex
		fmt.Println("Dex added to memory.")
		fmt.Println("\trouter:\t\t", dex.Router)
		fmt.Println("\tfactory:\t", dex.Factory)
	}

}
func (m *Memory) AddDex(router, factory common.Address, numPairs int) {
	_, exists := m.DexMemory[router]
	if !exists {
		m.DexMemory[router] = &types.DexMemory{
			Factory:    factory,
			Router:     router,
			Pairs:      map[string]*types.Pair{},
			SimplePair: map[common.Address]*types.SimplePair{},
			PairMutex:  map[string]*sync.Mutex{},
			NumPairs:   numPairs,
		}
		fmt.Println("Dex added to memory.")
		fmt.Println("\trouter:\t\t", router)
		fmt.Println("\tfactory:\t", factory)
	}
}

func (m *Memory) AddPairStruct(routerAddress common.Address, pair *types.Pair) {

	//sorted, token0 + token1 is the key
	key := pair.Token0Address.Hex() + pair.Token1Address.Hex()
	if _, exists := m.DexMemory[routerAddress].Pairs[key]; !exists {
		m.DexMemory[routerAddress].Pairs[key] = pair
	}
	if _, exists := m.DexMemory[routerAddress].PairMutex[key]; !exists {
		m.DexMemory[routerAddress].PairMutex[key] = &sync.Mutex{}
	}
	if _, exists := m.DexMemory[routerAddress].SimplePair[pair.PairAddress]; !exists {
		m.DexMemory[routerAddress].SimplePair[pair.PairAddress] = &types.SimplePair{
			Token0Address: pair.Token0Address,
			Token1Address: pair.Token1Address,
		}
	}
}

func (m *Memory) AddPair(router common.Address, pairAddress, token0, token1 common.Address, reserve0, reserve1 *big.Int, lastupdated uint32) {

	_, exists := m.DexMemory[router].Pairs[token0.Hash().Hex()+token1.Hash().Hex()]
	if !exists {
		m.DexMemory[router].Pairs[token0.Hash().Hex()+token1.Hash().Hex()] = &types.Pair{
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
