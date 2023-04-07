package types

import "github.com/ethereum/go-ethereum/common"

type Dex struct {
	Factory  *common.Address
	Router   *common.Address
	NumPairs *int
}
