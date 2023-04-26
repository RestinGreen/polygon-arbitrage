package dbclient

import (
	"context"
	"log"
	"time"

	pb "github.com/RestinGreen/protobuf/generated"
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

func (c *GRPCClient) GetAllDex() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	_, err := c.client.GetAllDex(ctx, &pb.GetAllDexRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}
