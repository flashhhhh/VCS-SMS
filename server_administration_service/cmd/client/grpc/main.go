package main

import (
	"context"
	"server_administration_service/pb"

	"github.com/flashhhhh/pkg/env"
	"google.golang.org/grpc"
)

func main() {
	// Load environment variables from the .env file
	environmentFilePath := "./configs/local.env"
	if err := env.LoadEnv(environmentFilePath); err != nil {
		panic("Failed to load environment variables from " + environmentFilePath + ": " + err.Error())
	} else {
		println("Environment variables loaded successfully from " + environmentFilePath)
	}

	// Connect to the gRPC server
	grpcServerAddress := env.GetEnv("GRPC_SERVER_ADMINISTRATION_HOST", "localhost") + ":" + env.GetEnv("SERVER_GRPC_ADMINISTRATION_PORT", "50051")

	conn, err := grpc.Dial(grpcServerAddress, grpc.WithInsecure())
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	client := pb.NewServerAdministrationServiceClient(conn)

	// Example usage of the client
	// Get all addresses
	addressesResponse, err := client.GetAllAddresses(context.Background(), &pb.EmptyRequest{})
	if err != nil {
		panic(err)
	}
	for _, address := range addressesResponse.Addresses {
		println("Server ID:", address.Id, ", Address:", address.Address)
	}
}