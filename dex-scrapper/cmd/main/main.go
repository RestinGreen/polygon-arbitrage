package main

import (
	"context"
	"flag"
	"log"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "github.com/RestinGreen/protobuf/generated"
)

func main() {
	addr := flag.String("addr", "localhost:50051", "the address to connect to")
	flag.Parse()
	// Set up a connection to the server.
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewDatabaseClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	_, err = c.GetAllDex(ctx, &pb.GetAllDexRequest{})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
}
