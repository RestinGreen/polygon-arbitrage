package memory

import (
	"sync"

	"github.com/RestinGreen/polygon-arbitrage/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type DexMemory struct {
	Dexs     map[common.Address]*types.Dex
	DexMutex map[common.Address]*sync.Mutex
}
