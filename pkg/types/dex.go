package types

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type DexMemory struct {
	Factory  common.Address
	Router   common.Address
	NumPairs int

	SimplePair map[common.Address]*SimplePair
	Pairs      map[string]*Pair
	PairMutex  map[string]*sync.Mutex
}

type PairMemory struct {
	PairMap   map[common.Address]*SimplePair
	Pairs     map[string]*Pair
	PairMutex map[string]*sync.Mutex
}
