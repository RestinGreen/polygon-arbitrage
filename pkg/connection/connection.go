package connection

import (
	"fmt"

	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

type Connection struct {
	RpcClient *rpc.Client
	EthClient *ethclient.Client
}

func NewConnection(endpoint string) *Connection {

	var c Connection
	var err error
		c.RpcClient, err = rpc.Dial(endpoint)
		if err != nil {
			fmt.Println("Failed to open rpc client connection. Exiting.")
			panic(err)
		}
		c.EthClient, err = ethclient.Dial(endpoint)
		if err != nil {
			fmt.Println("Failed to open eth client connection. Exiting.")
			panic(err)
		}

	return &c
}
