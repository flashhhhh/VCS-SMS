package grpc

import (
	"healthcheck_service/pb"

	"google.golang.org/grpc"
)

func StartGRPCClient() (pb.ServerAdministrationServiceClient, error) {
	// Create a connection to the server.
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())

	// Create a new client
	client := pb.NewServerAdministrationServiceClient(conn)

	return client, err
}