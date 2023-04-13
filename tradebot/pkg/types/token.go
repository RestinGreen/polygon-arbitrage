package types

import "github.com/ethereum/go-ethereum/common"

type Token struct {
	Address *common.Address
	Symbol  *string
	Name    *string
	Decimal *int
}
