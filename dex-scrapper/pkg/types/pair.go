package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

type Pair struct {
	PairAddress   *common.Address
	Token0        *Token
	Token1        *Token
	RouterAddress *common.Address
	Reserve0      *big.Int
	Reserve1      *big.Int
	LastUpdated   *int64
}
