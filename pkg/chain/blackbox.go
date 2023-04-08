package chain

import (
	"github.com/ethereum/go-ethereum/common"
)

func SortAddress(tokenA common.Address, tokenB common.Address) (common.Address, common.Address, bool) {

	if tokenA.Hash().Big().Cmp(tokenB.Hash().Big()) < 0 {
		return tokenA, tokenB, false
	} else {
		return tokenB, tokenA, true
	}
}
