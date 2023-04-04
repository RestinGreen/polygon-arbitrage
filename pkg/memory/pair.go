package memory

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Pair struct {
	PairAddress   common.Address
	Token0Address common.Address
	Token1Address common.Address
	Reserve0      *big.Int
	Reserve1      *big.Int
	LastUpdated   uint32
}
