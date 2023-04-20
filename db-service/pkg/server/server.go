package server

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"

	pb "github.com/RestinGreen/protobuf/generated"
	"google.golang.org/grpc"
)

type Server struct {
	pb.UnimplementedDatabaseServer
}

func (s *Server) GetAllDex(ctx context.Context, in *pb.GetAllDexRequest) (*pb.GetAllDexResponse, error) {
	fmt.Println("getall invoked")
	return &pb.GetAllDexResponse{
		
	}, nil
}

func NewServer() {

	port := flag.Int("port", 50051, "The server port")
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterDatabaseServer(s, &Server{})
	log.Printf("server listening at %v", lis.Addr())
	
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
