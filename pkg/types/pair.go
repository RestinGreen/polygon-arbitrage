package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type SimplePair struct {
	Token0Address common.Address
	Token1Address common.Address
}

type Pair struct {
	PairAddress   common.Address
	Token0Address common.Address
	Token1Address common.Address
	Reserve0      *big.Int
	Reserve1      *big.Int
	LastUpdated   uint32
}
