package abihandler

import (
	"fmt"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

type AbiHandler  struct {
	UniV2RouterAbi abi.ABI
}

func NewAbiHandler() *AbiHandler {
	var a AbiHandler 

	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("Failed to get current working directory.")
		panic(err)
	}

	fileBytes, err := os.ReadFile(wd + "/pkg/abihandler/IUniswapV2Router02.json")
	if err != nil {
		fmt.Println("Failed to read Uniswap Router json.")
		panic(err)
	}
	uniV2ABIString := string(fileBytes)
	a.UniV2RouterAbi, err = abi.JSON(strings.NewReader(uniV2ABIString))


	return &a
}