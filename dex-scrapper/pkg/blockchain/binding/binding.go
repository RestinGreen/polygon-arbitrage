package binding

import (
	"fmt"
	"sync"

	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain/binding/erc20"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain/binding/univ2factory"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain/binding/univ2pair"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/blockchain/binding/univ2router"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Binding struct {
	ethClient *ethclient.Client

	pairsMutex     *sync.Mutex
	routersMutex   *sync.Mutex
	factoriesMutex *sync.Mutex
	tokensMutex    *sync.Mutex

	pairs     map[common.Address]*univ2pair.UniV2Pair
	Routers   map[common.Address]*univ2router.UniV2Router
	Factories map[common.Address]*univ2factory.UniV2Factory
	Tokens    map[common.Address]*erc20.ERC20
}

func NewBinding(ethClient *ethclient.Client) *Binding {

	return &Binding{
		ethClient: ethClient,

		pairsMutex:     &sync.Mutex{},
		routersMutex:   &sync.Mutex{},
		factoriesMutex: &sync.Mutex{},
		tokensMutex:    &sync.Mutex{},

		pairs:     make(map[common.Address]*univ2pair.UniV2Pair),
		Routers:   make(map[common.Address]*univ2router.UniV2Router),
		Factories: make(map[common.Address]*univ2factory.UniV2Factory),
		Tokens:    make(map[common.Address]*erc20.ERC20),
	}
}

func (b *Binding) GetPairContract(pairAddress common.Address) *univ2pair.UniV2Pair {
	b.pairsMutex.Lock()
	defer b.pairsMutex.Unlock()
	return b.pairs[pairAddress]
}

func (b *Binding) AddPairContract(pairAddress *common.Address) *univ2pair.UniV2Pair {

	newPairContract, err := univ2pair.NewUniV2Pair(*pairAddress, b.ethClient)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to add %s pair binding.", pairAddress.Hex()))
	}
	b.pairsMutex.Lock()
	b.pairs[*pairAddress] = newPairContract
	b.pairsMutex.Unlock()

	return newPairContract
}

func (b *Binding) AddRouterContract(routerAddress *common.Address) *univ2router.UniV2Router {

	newRouterContract, err := univ2router.NewUniV2Router(*routerAddress, b.ethClient)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to add %s router binding.", routerAddress.Hex()))
	}
	b.routersMutex.Lock()
	b.Routers[*routerAddress] = newRouterContract
	b.routersMutex.Unlock()

	return newRouterContract
}

func (b *Binding) AddFactoryContract(factoryAddress *common.Address) *univ2factory.UniV2Factory {

	newFactoryContract, err := univ2factory.NewUniV2Factory(*factoryAddress, b.ethClient)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to add %s router binding.", factoryAddress.Hex()))
	}
	b.factoriesMutex.Lock()
	b.Factories[*factoryAddress] = newFactoryContract
	b.factoriesMutex.Unlock()

	return newFactoryContract
}

func (b *Binding) AddTokenContract(tokenAddress *common.Address) *erc20.ERC20 {

	b.tokensMutex.Lock()
	_, exists := b.Tokens[*tokenAddress]
	b.tokensMutex.Unlock()
	if exists {
		return b.Tokens[*tokenAddress]
	}
	newTokenContract, err := erc20.NewERC20(*tokenAddress, b.ethClient)
	if err != nil {
		fmt.Println(fmt.Sprintf("Failed to add %s router binding.", tokenAddress.Hex()))
	}
	b.tokensMutex.Lock()
	b.Tokens[*tokenAddress] = newTokenContract
	b.tokensMutex.Unlock()

	return newTokenContract
}
