package memory

import (
	"sync"

	"github.com/RestinGreen/polygon-arbitrage/tradebot/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type PairMemory struct {
	//key is pait address -> t0+t1+router
	PairMap   map[common.Address]*string
	
	//key is t0+t1+router -> pair data
	Pairs     map[string]*types.Pair
	PairMutex map[string]*sync.Mutex
}
