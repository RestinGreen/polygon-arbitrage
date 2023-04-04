package binding

import (
	"fmt"

	"github.com/RestinGreen/polygon-arbitrage/pkg/binding/univ2factory"
	"github.com/RestinGreen/polygon-arbitrage/pkg/binding/univ2pair"
	"github.com/RestinGreen/polygon-arbitrage/pkg/binding/univ2router"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Binding struct {
	ethClient *ethclient.Client

	Pairs     map[common.Address]*univ2pair.UniV2Pair
	Routers   map[common.Address]*univ2router.UniV2Router
	Factories map[common.Address]*univ2factory.UniV2Factory
}

func NewBinding(ethClient *ethclient.Client) *Binding {

	return &Binding{
		ethClient: ethClient,
		Pairs:     map[common.Address]*univ2pair.UniV2Pair{},
		Routers:   map[common.Address]*univ2router.UniV2Router{},
		Factories: map[common.Address]*univ2factory.UniV2Factory{},
	}
}

func (b *Binding) AddPair(pairAddress common.Address) {

	newPair, err := univ2pair.NewUniV2Pair(pairAddress, b.ethClient)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to add %s pair binding.", pairAddress.Hex()))
	}
	b.Pairs[pairAddress] = newPair
}

func (b *Binding) AddRouter(routerAddress common.Address) *univ2router.UniV2Router {
	newRouter, err := univ2router.NewUniV2Router(routerAddress, b.ethClient)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to add %s router binding.", routerAddress.Hex()))
	}
	b.Routers[routerAddress] = newRouter
	return newRouter
}

func (b *Binding) AddFactory(factoryAddress common.Address) *univ2factory.UniV2Factory {
	newFactory, err := univ2factory.NewUniV2Factory(factoryAddress, b.ethClient)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to add %s router binding.", factoryAddress.Hex()))
	}
	b.Factories[factoryAddress] = newFactory
	return newFactory
}
