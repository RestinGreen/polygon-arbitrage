package memory

import (
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/RestinGreen/polygon-arbitrage/pkg/binding"
	"github.com/RestinGreen/polygon-arbitrage/pkg/binding/univ2factory"
	"github.com/RestinGreen/polygon-arbitrage/pkg/chain"
	"github.com/RestinGreen/polygon-arbitrage/pkg/database"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

type DexMonitor struct {
	DexMemory DexMemory
	txChan    chan *types.Transaction
	binding   *binding.Binding
	read      *bind.CallOpts

	db          *database.Database
	fullNode    *ethclient.Client
	archiveNode *ethclient.Client
}

func NewDexMonirot(txChan chan *types.Transaction, fullNode *ethclient.Client, archiveNode *ethclient.Client, db *database.Database) *DexMonitor {

	return &DexMonitor{
		txChan:    txChan,
		DexMemory: *NewDexMemory(),
		binding:   binding.NewBinding(fullNode),
		read:      &bind.CallOpts{Pending: false},

		db: db,
		fullNode:    fullNode,
		archiveNode: archiveNode,
	}

}


func (m *DexMonitor) LoadFromDB() {

	m.db.GetAllDexs()
}

func (m *DexMonitor) Start() {

	for {
		select {

		case tx := <-m.txChan:
			routerAddress := tx.To()
			if routerAddress == nil {
				fmt.Println("to address is nil")
				continue
			}
			txInputData := tx.Data()
			routerContract, exists := m.binding.Routers[*routerAddress]
			if !exists {
				routerContract = m.binding.AddRouter(*routerAddress)
			}
			factoryAddress, err := routerContract.Factory(m.read)
			if err != nil {
				fmt.Println("Failed to read factory address from router ", routerAddress)
			}

			factoryContract, exists := m.binding.Factories[factoryAddress]
			if !exists {
				factoryContract = m.binding.AddFactory(factoryAddress)
			}

			//add dex to memory
			m.DexMemory.AddDex(*routerAddress, factoryAddress)

			// pairCnt, err := factoryContract.AllPairsLength(m.read)
			// if err != nil {
			// 	fmt.Println("Failed to get number of pairs in factory")
			// 	continue
			// }
			// for i := 0; i < pairCnt.Int64(); i++ {

			// }

			//decoding pairs from tx data and adding to memory
			switch hex.EncodeToString(txInputData[0:4]) {
			case "fb3bdb41", "7ff36ab5", "b6f9de95":
				m.decodeInputData(*routerAddress, factoryContract, txInputData, 4)

			case "18cbafe5", "791ac947", "38ed1739", "5c11d795", "4a25d94a", "8803dbee":
				m.decodeInputData(*routerAddress, factoryContract, txInputData, 5)
			}
		}
	}
}

func (m *DexMonitor) decodeInputData(router common.Address, factoryContract *univ2factory.UniV2Factory, data []byte, pathLengthIndex int) {

	length := int(new(big.Int).SetBytes(getithParam(data, pathLengthIndex)).Int64())
	for i := 0; i < length-1; i++ {
		tokenAByte := getithParam(data, pathLengthIndex+1+i)
		tokenBByte := getithParam(data, pathLengthIndex+1+i+1)

		addressA := common.BytesToAddress(tokenAByte)
		addressB := common.BytesToAddress(tokenBByte)

		token0, token1 := chain.SortAddress(addressA, addressB)

		// pairAddress := m.binding.
		pairAddress, err := factoryContract.GetPair(m.read, token0, token1)
		m.binding.AddPair(pairAddress)
		if err != nil {
			fmt.Println("Failed to get pairAddress")
			continue
		}
		reserves, err := m.binding.Pairs[pairAddress].GetReserves(m.read)
		if err != nil {
			fmt.Println("Failed to get reserves.")
			continue
		}

		m.DexMemory.AddPair(router, pairAddress, token0, token1, reserves.Reserve0, reserves.Reserve1, reserves.BlockTimestampLast)

	}
}

// Param order in bytecode:
// 0
func getithParam(bytes []byte, i int) []byte {

	firstIndex := 4 + i*32
	lastIndex := 4 + (i+1)*32
	return bytes[firstIndex:lastIndex]
}

func getithTopicData(bytes []byte, i int) []byte {
	firstIndex := i * 32
	lastIndex := (i + 1) * 32
	return bytes[firstIndex:lastIndex]
}
