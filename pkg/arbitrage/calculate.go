package arbitrage

import (
	"encoding/hex"
	"math/big"

	"github.com/RestinGreen/polygon-arbitrage/pkg/chain"
	"github.com/RestinGreen/polygon-arbitrage/pkg/memory"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Arbitrage struct {
	Memory *memory.Memory
	txCh   chan *types.Transaction
}

func NewArbitrage(txCh chan *types.Transaction, memory *memory.Memory) *Arbitrage {

	return &Arbitrage{
		Memory: memory,
		txCh:   txCh,
	}
}

/*
	precalculate the best path for every token
	precalculate the minimum amount of input from user that opens up an arb oppurtunity

*/

func (a *Arbitrage) StartArbitrage() {

	for {
		select {
		case tx := <-a.txCh:

			txInputData := tx.Data()
			routerAddress := tx.To()
			// decoding pairs from tx data and adding to memory
			switch hex.EncodeToString(txInputData[0:4]) {
			case "fb3bdb41", "7ff36ab5", "b6f9de95":
				a.decodeAndCalculate(txInputData, routerAddress, 4)

			case "18cbafe5", "791ac947", "38ed1739", "5c11d795", "4a25d94a", "8803dbee":
				a.decodeAndCalculate(txInputData, routerAddress, 5)
			}
		}
	}
}

func (a *Arbitrage) decodeAndCalculate(data []byte, router *common.Address, pathLengthIndex int) {

	length := int(new(big.Int).SetBytes(getithParam(data, pathLengthIndex)).Int64())
	for i := 0; i < length-1; i++ {
		tokenAByte := getithParam(data, pathLengthIndex+1+i)
		tokenBByte := getithParam(data, pathLengthIndex+1+i+1)

		fromToken := common.BytesToAddress(tokenAByte)
		toToken := common.BytesToAddress(tokenBByte)

		_, _, inverse := chain.SortAddress(fromToken, toToken)

		if !inverse {
			
			// sourceKey := fromToken.Hex()+toToken.Hex()+router.Hex()
			// source := a.Memory.PairMemory.Pairs[sourceKey]


		}

	}
}

// func (a *Arbitrage) calculate(tok)

// Param order in bytecode:
// 0
func getithParam(bytes []byte, i int) []byte {

	firstIndex := 4 + i*32
	lastIndex := 4 + (i+1)*32
	return bytes[firstIndex:lastIndex]
}
