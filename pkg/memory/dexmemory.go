package memory

import (
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Dex struct {
	Factory  common.Address
	Router   common.Address
	Pairs    map[string]Pair
}

type DexMemory struct {
	//key is router
	DexMap map[common.Address]*Dex
}

func NewDexMemory() *DexMemory {

	return &DexMemory{
		DexMap: map[common.Address]*Dex{},
	}
}

func (m *DexMemory) AddDex(router, factory common.Address) {
	_, exists := m.DexMap[router]
	if !exists {
		m.DexMap[router] = &Dex{
			Factory: factory,
			Router:  router,
			Pairs:   map[string]Pair{},
		}
		fmt.Println("Inserted to", router, "the factory", factory)
	}
}

func (m *DexMemory) AddPair(router common.Address, pairAddress, token0, token1 common.Address, reserve0, reserve1 *big.Int, lastupdated uint32) {

	_, exists := m.DexMap[router].Pairs[token0.Hash().Hex()+token1.Hash().Hex()]
	if !exists {
		m.DexMap[router].Pairs[token0.Hash().Hex()+token1.Hash().Hex()] = Pair{
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
