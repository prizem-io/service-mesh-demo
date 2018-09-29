package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/prizem-io/service-mesh-demo/proto/message"
)

type messageServer struct{}

func (s *messageServer) Hello(_ context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	name := req.Name
	log.Printf("[INFO] Request: %s", name)
	return &pb.HelloResponse{
		Message: fmt.Sprintf("hello, %s (gRPC)", name),
	}, nil
}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterMessageServer(grpcServer, &messageServer{})

	log.Println("[INFO] Started gRPC server on :4002")
	listener, err := net.Listen("tcp", ":4002")
	if err != nil {
		log.Fatalf("[ERROR] Failed to listen: %s", err)
	}

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("[ERROR] Failed to start server")
	}
}
