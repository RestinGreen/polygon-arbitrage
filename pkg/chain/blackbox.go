package chain

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
)

func SortAddress(tokenA common.Address, tokenB common.Address) (common.Address, common.Address) {

	if tokenA.Hash().Big().Cmp(tokenB.Hash().Big()) < 0 {
		return tokenA, tokenB
	} else {
		return tokenB, tokenA
	}
}

func Crate2(tokenA common.Address, tokenB common.Address, factory common.Address, inithash string) common.Address {

	tA, tB := SortAddress(tokenA, tokenB)
	ff := []byte{255}
	ff = append(ff, factory.Bytes()...)
	ff = append(ff, crypto.Keccak256(append(tA.Bytes(), tB.Bytes()...))...)
	ff = append(ff, common.Hex2Bytes(inithash)...)

	return common.BytesToAddress(crypto.Keccak256(ff))
}
