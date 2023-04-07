package memory

import (
	"sync"

	"github.com/RestinGreen/polygon-arbitrage/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type PairMemory struct {
	PairMap   map[common.Address]*types.SimplePair
	Pairs     map[string]*types.Pair
	PairMutex map[string]*sync.Mutex
}
