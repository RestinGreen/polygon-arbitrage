package memory

import (
	"sync"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/types"
	"github.com/ethereum/go-ethereum/common"
)

type TokenMemory struct {
	Tokens     map[common.Address]*types.Token
	TokemMutex *sync.Mutex
}
