package grpcclient

import (
	"context"
	"mail_service/pb"
)

type ServerAdministrationServiceClient interface {
	GetServerInformation(startTime, endTime int64) (*pb.GetServerInformationResponse, error)
}

type serverAdministrationServiceClient struct {
	client pb.ServerAdministrationServiceClient
}

func NewServerAdministrationServiceClient(client pb.ServerAdministrationServiceClient) ServerAdministrationServiceClient {
	return &serverAdministrationServiceClient{
		client: client,
	}
}

func (s *serverAdministrationServiceClient) GetServerInformation(startTime, endTime int64) (*pb.GetServerInformationResponse, error) {
	resp, err := s.client.GetServerInformation(
		context.Background(),
		&pb.GetServerInformationRequest{
			StartTime: startTime,
			EndTime:   endTime,
		},
	)

	if err != nil {
		panic("Failed to get server information: " + err.Error())
	}

	return resp, nil
}