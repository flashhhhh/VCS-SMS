package main

import (
	"healthcheck_service/infrastructure/grpc"
	"healthcheck_service/infrastructure/healthcheck"
	grpcclient "healthcheck_service/internal/grpc_client"
)

func main() {
	grpcClient, err := grpc.StartGRPCClient()
	if err != nil {
		panic(err)
	}

	client := grpcclient.NewHealthCheckClient(grpcClient)
	addressesResponse, err := client.GetAllAddresses()
	if err != nil {
		panic(err)
	}
	for _, address := range addressesResponse.Addresses {
		serverID := address.ServerId
		serverAddress := address.Address

		// Check if the server is On or Off by pinging the address
		status := healthcheck.IsHostUp("localhost:80")
		if status == false {
			// Server is Off
			println("Server ID:", serverID, "at address", serverAddress, "is Off")
		} else {
			// Server is On
			println("Server ID:", serverID, "at address", serverAddress, "is On")
		}
	}
}