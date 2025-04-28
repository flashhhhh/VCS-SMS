package main

import (
	"context"
	"server_administration_service/pb"

	"github.com/flashhhhh/pkg/env"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Load environment variables from the .env file
	environmentFilePath := "./configs/local.env"
	if err := env.LoadEnv(environmentFilePath); err != nil {
		panic("Failed to load environment variables from " + environmentFilePath + ": " + err.Error())
	} else {
		println("Environment variables loaded successfully from " + environmentFilePath)
	}

	grpcServerAddress := "localhost:50052"

	conn, err := grpc.Dial(grpcServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
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

	// Get server information
	startTime := int64(1672531199) // Example start time
	endTime := int64(1672617599)   // Example end time
	serverInfoResponse, err := client.GetServerInformation(context.Background(), &pb.GetServerInformationRequest{
		StartTime: startTime,
		EndTime:   endTime,
	})

	if err != nil {
		panic(err)
	}
	println("Number of ON servers:", serverInfoResponse.NumOnServers)
	println("Number of OFF servers:", serverInfoResponse.NumOffServers)
}