package types

import (
	"sync"

	"github.com/ethereum/go-ethereum/common"
)

type Dex struct {
	Factory  common.Address
	Router   common.Address
	Pairs    map[string]*Pair
	NumPairs int

	PairMutex map[string]*sync.Mutex
}
