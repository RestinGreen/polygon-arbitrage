package memory

import (
	"fmt"
	"math/big"
	"sync"

	"github.com/RestinGreen/polygon-arbitrage/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type DexMemory struct {
	//key is router
	DexMap map[common.Address]*types.Dex
}

func NewDexMemory() *DexMemory {

	return &DexMemory{
		DexMap: map[common.Address]*types.Dex{},
	}
}

func (m *DexMemory) AddDexStruct(dex *types.Dex) {
	_, exists := m.DexMap[dex.Router]
	if !exists {
		m.DexMap[dex.Router] = dex
		fmt.Println("Dex added to memory.")
		fmt.Println("\trouter:\t", dex.Router)
		fmt.Println("\tfactory:\t", dex.Factory)
	}

}
func (m *DexMemory) AddDex(router, factory common.Address, numPairs int) {
	_, exists := m.DexMap[router]
	if !exists {
		m.DexMap[router] = &types.Dex{
			Factory:   factory,
			Router:    router,
			Pairs:     map[string]*types.Pair{},
			PairMutex: map[string]*sync.Mutex{},
			NumPairs:  numPairs,
		}
		fmt.Println("Dex added to memory.")
		fmt.Println("\trouter:\t", router)
		fmt.Println("\tfactory:\t", factory)
	}
}

func (m *DexMemory) AddPairStruct(routerAddress common.Address, pair *types.Pair) {

	//sorted token0 + token1 is the key
	key := pair.Token0Address.Hex() + pair.Token1Address.Hex()
	if _, exists := m.DexMap[routerAddress].Pairs[key]; !exists {
		m.DexMap[routerAddress].Pairs[key] = pair
	}
	if _, exists := m.DexMap[routerAddress].PairMutex[key]; !exists {
		m.DexMap[routerAddress].PairMutex[key] = &sync.Mutex{}
	}
}

func (m *DexMemory) AddPair(router common.Address, pairAddress, token0, token1 common.Address, reserve0, reserve1 *big.Int, lastupdated uint32) {

	_, exists := m.DexMap[router].Pairs[token0.Hash().Hex()+token1.Hash().Hex()]
	if !exists {
		m.DexMap[router].Pairs[token0.Hash().Hex()+token1.Hash().Hex()] = &types.Pair{
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
