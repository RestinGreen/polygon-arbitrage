package memory

import (
	"sync"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type Memory struct {
	DexMemory   *DexMemory
	PairMemory  *PairMemory
	TokenMemory *TokenMemory
}

func NewMemory() *Memory {

	return &Memory{
		DexMemory: &DexMemory{
			Dexs:     map[common.Address]*types.Dex{},
			DexMutex: &sync.Mutex{},
		},
		PairMemory: &PairMemory{
			PairMap:   map[common.Address]*string{},
			MapMutex:  &sync.Mutex{},
			Pairs:     map[string]*types.Pair{},
			PairMutex: &sync.Mutex{},
		},
		TokenMemory: &TokenMemory{
			Tokens: make([]*types.Token, 0),
		},
	}
}

func (m *Memory) AddDex(router, factory *common.Address, numPairs *int64) {

	_, exists := m.DexMemory.Dexs[*router]
	if !exists {
		m.DexMemory.DexMutex.Lock()
		m.DexMemory.Dexs[*router] = &types.Dex{
			Factory:  factory,
			Router:   router,
			NumPairs: numPairs,
		}
		m.DexMemory.DexMutex.Unlock()
	}
}

func (m *Memory) AddPairStruct(pair *types.Pair) {

	//sorted, router + token0 + token1 is the key
	key := pair.RouterAddress.Hex() + pair.Token0Address.Hex() + pair.Token1Address.Hex()

	m.PairMemory.PairMutex.Lock()
	if _, exists := m.PairMemory.Pairs[key]; !exists {
		m.PairMemory.Pairs[key] = pair
	}
	m.PairMemory.PairMutex.Unlock()

	m.PairMemory.MapMutex.Lock()
	if _, exists := m.PairMemory.PairMap[*pair.PairAddress]; !exists {
		m.PairMemory.PairMap[*pair.PairAddress] = &key
	}
	m.PairMemory.MapMutex.Unlock()
}
