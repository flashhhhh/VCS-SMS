package grpcclient

import (
	"context"
	"healthcheck_service/pb"
)

type HealthCheckClient interface {
	GetAllAddresses() (*pb.AddressesResponse, error)
}

type healthCheckClient struct {
	client pb.ServerAdministrationServiceClient
}

func NewHealthCheckClient(client pb.ServerAdministrationServiceClient) HealthCheckClient {
	return &healthCheckClient{
		client: client,
	}
}

func (h *healthCheckClient) GetAllAddresses() (*pb.AddressesResponse, error) {
	resp, err := h.client.GetAllAddresses(
		context.Background(),
		&pb.EmptyRequest{},
	)

	if err != nil {
		return nil, err
	}

	return resp, nil
}