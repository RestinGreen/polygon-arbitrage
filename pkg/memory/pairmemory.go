package memory

import (
	"sync"

	"github.com/RestinGreen/polygon-arbitrage/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type PairMemory struct {
	PairMap   map[common.Address]*string
	
	//key is t0 + t1 + router
	Pairs     map[string]*types.Pair
	PairMutex map[string]*sync.Mutex
}
