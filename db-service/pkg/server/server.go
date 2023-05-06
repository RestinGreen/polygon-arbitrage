package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	"github.com/RestinGreen/polygon-arbitrage/db-service/pkg/database"
	pb "github.com/RestinGreen/protobuf/generated"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedDatabaseServer
	DB *database.Database
}

func NewServer(db *database.Database) {

	port := flag.Int("port", 50051, "The server port")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()

	pb.RegisterDatabaseServer(s, &Server{DB: db})
	log.Printf("server listening at %v", lis.Addr())

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

func (s *Server) GetAllDex(ctx context.Context, in *pb.GetAllDexRequest) (*pb.GetAllDexResponse, error) {
	fmt.Println("GetAllDex invoked")
	dexs, _ := s.DB.GetAllDex()
	return &pb.GetAllDexResponse{
		Dexs: dexs,
	}, nil
}

func (s *Server) InsertDex(ctx context.Context, in *pb.InsertDexRequest) (*pb.InsertDexResponse, error) {
	fmt.Println("Inserting dex into database")

	s.DB.InsertDex(&in.Dex.FactoryAddress, &in.Dex.RouterAddress, &in.Dex.NumPairs, in.Dex.Pairs)

	return &pb.InsertDexResponse{}, nil
}

func (s *Server) UpdatePair(ctx context.Context, in *pb.UpdatePairRequest) (*pb.UpdatePairResponse, error) {

	s.DB.UpdatePair(in.Address, in.Reserve0, in.Reserve1, &in.LastUpdated)

	return &pb.UpdatePairResponse{}, nil
}
