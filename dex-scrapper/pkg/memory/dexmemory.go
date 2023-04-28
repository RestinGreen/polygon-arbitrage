package memory

import (
	"sync"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type DexMemory struct {
	Dexs     map[common.Address]*types.Dex
	DexMutex *sync.Mutex
}
