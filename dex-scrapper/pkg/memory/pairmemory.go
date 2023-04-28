package memory

import (
	"sync"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type PairMemory struct {
	//key is pair address -> router+t0+t1
	PairMap map[common.Address]*string
	MapMutex *sync.Mutex

	//key is t0+t1+router -> pair data
	Pairs     map[string]*types.Pair
	PairMutex *sync.Mutex
}
