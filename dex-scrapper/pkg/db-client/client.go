package dbclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/RestinGreen/protobuf/generated"
	pb "github.com/RestinGreen/protobuf/generated"
	"github.com/Restingreen/polygon-arbitrage/dex-scrapper/pkg/types"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	connection *grpc.ClientConn
	client     pb.DatabaseClient
}

func NewClient(host, port string) *GRPCClient {

	conn, err := grpc.Dial(host+":"+port, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	dbClinet := pb.NewDatabaseClient(conn)

	return &GRPCClient{
		connection: conn,
		client:     dbClinet,
	}
}

func (c *GRPCClient) GetAllDex() []*generated.Dex {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	resp, err := c.client.GetAllDex(ctx, &pb.GetAllDexRequest{})
	if err != nil {
		fmt.Println("Failed GetAllDex grpc.\n", err)
		return nil
	}
	return resp.Dexs

}

func (c *GRPCClient) InsertDex(dex *types.Dex, pairs []*types.Pair) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	pairsMessage := make([]*pb.Pair, 0)

	for _, p := range pairs {
		pairsMessage = append(pairsMessage, &pb.Pair{
			Token0: &pb.Token{
				Address: p.Token0.Address.Hex(),
				Name:    *p.Token0.Name,
				Symbol:  *p.Token0.Symbol,
				Decimal: int64(*p.Token0.Decimal),
			},
			Token1: &pb.Token{
				Address: p.Token1.Address.Hex(),
				Name:    *p.Token1.Name,
				Symbol:  *p.Token1.Symbol,
				Decimal: *p.Token1.Decimal,
			},
			Address:     p.PairAddress.Hex(),
			Reserve0:    p.Reserve0.Bytes(),
			Reserve1:    p.Reserve1.Bytes(),
			LastUpdated: *p.LastUpdated,
		})
	}

	_, err := c.client.InsertDex(ctx, &pb.InsertDexRequest{
		Dex: &pb.Dex{
			FactoryAddress: dex.Factory.Hex(),
			RouterAddress:  dex.Router.Hex(),
			NumPairs:       *dex.NumPairs,
			Pairs:          pairsMessage,
		},
	})
	if err != nil {
		fmt.Println("Failed to call gRPC InsertDex")
		fmt.Println(err)
	}
}

func (c *GRPCClient) UpdatePair(pairAddress string, reserve0, reserve1 []byte, lastUpdated *int64) {

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()


	_, err := c.client.UpdatePair(ctx, &pb.UpdatePairRequest{
		Address: pairAddress,
		Reserve0: reserve0,
		Reserve1: reserve1,
		LastUpdated: *lastUpdated,

	})
	if err != nil {
		fmt.Println("Failed pair update.\n", err)
		return 
	}
}