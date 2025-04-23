package handler

import (
	"context"
	"server_administration_service/internal/service"
	"server_administration_service/pb"
)

type GRPCServerHandler struct {
	serverService service.ServerService
	pb.UnimplementedServerAdministrationServiceServer
}

func NewGrpcServerHandler(serverService service.ServerService) *GRPCServerHandler {
	return &GRPCServerHandler{
		serverService: serverService,
	}
}

func (grpcHandler *GRPCServerHandler) GetAllAddresses(ctx context.Context, req *pb.EmptyRequest) (*pb.AddressesResponse, error) {
	addresses, err := grpcHandler.serverService.GetAllAddresses()
	if err != nil {
		return nil, err
	}

	addressInfo := make([]*pb.AddressInfo, len(addresses))
	for i, address := range addresses {
		addressInfo[i] = &pb.AddressInfo{
			ServerId: address[0],
			Address:   address[1],
		}
	}

	response := &pb.AddressesResponse{
		Addresses: addressInfo,
	}

	return response, nil
}