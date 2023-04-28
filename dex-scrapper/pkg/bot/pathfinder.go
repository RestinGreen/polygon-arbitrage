package bot

import "github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/memory"

type Pathfinder struct {

	memory *memory.Memory

}

func NewPathFinder(memory *memory.Memory) *Pathfinder {
	
	return &Pathfinder{
		memory: memory,
	}
}

