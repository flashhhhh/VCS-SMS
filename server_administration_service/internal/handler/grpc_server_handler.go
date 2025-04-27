package handler

import (
	"context"
	"server_administration_service/internal/service"
	"server_administration_service/pb"
	"strconv"
	"time"
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
		addressLink := address.IPv4
		if address.Port >= 0 {
			addressLink += ":" + strconv.Itoa(address.Port)
		}

		addressInfo[i] = &pb.AddressInfo{
			Id:  int64(address.ID),
			Address: addressLink,
		}
	}

	response := &pb.AddressesResponse{
		Addresses: addressInfo,
	}

	return response, nil
}

func (grpcHandler *GRPCServerHandler) GetServerInformation(ctx context.Context, req *pb.GetServerInformationRequest) (*pb.GetServerInformationResponse, error) {
	numOnServers, _ := grpcHandler.serverService.GetNumOnServers()
	numServers, _ := grpcHandler.serverService.GetNumServers()
	numOffServers := numServers - numOnServers

	startTime := req.GetStartTime()
	endTime := req.GetEndTime()

	// Convert int64 timestamps to time.Time
	startTimeObj := time.Unix(startTime, 0)
	endTimeObj := time.Unix(endTime, 0)

	// Call the service method to get the server uptime ratio
	uptimeRatio, err := grpcHandler.serverService.GetServerUptimeRatio(startTimeObj, endTimeObj)
	if err != nil {
		return nil, err
	}

	response := &pb.GetServerInformationResponse{
		NumServers: int64(numServers),
		NumOnServers: int64(numOnServers),
		NumOffServers: int64(numOffServers),
		MeanUptimeRatio: float32(uptimeRatio),
	}
	
	return response, nil
}