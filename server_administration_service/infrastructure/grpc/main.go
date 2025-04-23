package grpc

import (
	"net"
	"server_administration_service/internal/handler"
	"server_administration_service/pb"

	"github.com/flashhhhh/pkg/logging"
	"google.golang.org/grpc"
)

func StartGRPCServer(serverHandler *handler.GRPCServerHandler, port string) {
	lis, err := net.Listen("tcp", ":" + port)
	if err != nil {
		panic(err)
	}

	// Create a new gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterServerAdministrationServiceServer(grpcServer, serverHandler)

	logging.LogMessage("server_administration_service", "gRPC server is running on port: "+port, "INFO")
	if err := grpcServer.Serve(lis); err != nil {
		logging.LogMessage("server_administration_service", "Failed to serve: "+err.Error(), "ERROR")
	}
}