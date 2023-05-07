package bot

import (
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/memory"
)

type Pathfinder struct {
	memory *memory.Memory
}

func NewPathFinder(memory *memory.Memory) *Pathfinder {

	return &Pathfinder{
		memory: memory,
	}
}

func (p *Pathfinder) BestRoutes() {

	// map[tokenIn][tokenOut][router] [] available pairs for tokenOut
	// trade := make(map[common.Address]map[common.Address]map[common.Address][]*t.Pair)

	// for pair, pairData := range p.memory.PairMemory.Pairs {

	// 	trade[*pairData.Token0.Address]

	// }

}
