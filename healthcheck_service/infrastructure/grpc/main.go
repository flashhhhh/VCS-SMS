package grpc

import (
	"healthcheck_service/pb"

	"github.com/flashhhhh/pkg/env"
	"google.golang.org/grpc"
)

func StartGRPCClient() (pb.ServerAdministrationServiceClient, error) {
	// Create a connection to the server.
	conn, err := grpc.Dial(env.GetEnv("GRPC_SERVER_ADMINISTRATION_SERVER", "localhost") + ":" + env.GetEnv("GRPC_SERVER_ADMINISTRATION_PORT", "50052"), grpc.WithInsecure())

	// Create a new client
	client := pb.NewServerAdministrationServiceClient(conn)

	return client, err
}