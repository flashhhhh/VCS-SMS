package grpc

import (
	"log"
	"net"
	grpchandler "user_service/internal/grpc_handler"
	"user_service/pb"

	"google.golang.org/grpc"
)

func StartGRPCServer(userHandler *grpchandler.UserHandler, port string) {
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic(err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcServer, userHandler)

	log.Println("gRPC server is running on port: ", port)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}