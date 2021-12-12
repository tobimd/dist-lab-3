package informant

import (
	"context"
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
)

// **** SERVER FUNCTIONS ****
type Server struct {
	pb.UnimplementedCommunicationServer
}

func (s *Server) RunCommand(ctx context.Context, command *pb.CommandParams) (*pb.FulcrumResponse, error) {
	log.Log(&f, "[server:RunCommand] Called with argument: command=\"%v\"", command.String())

	switch command.Command.Enum() {

	case pb.Command_ADD_CITY.Enum():

	case pb.Command_UPDATE_NAME.Enum():

	case pb.Command_UPDATE_NUMBER.Enum():

	case pb.Command_DELETE_CITY.Enum():
	}

	return &pb.FulcrumResponse{}, nil
}

// **** CLIENT FUNCTIONS ****
type Client data.GrpcClient

func (c *Client) RunCommand(command *pb.CommandParams) *pb.FulcrumResponse {
	client := *((*c).Client)
	log.Log(&f, "[client:RunCommand] Called with argument: command=\"%v\"", command.String())

	ctx, cancel := util.GetContext()
	defer cancel()

	res, err := client.RunCommand(ctx, command)
	log.FailOnError(&f, err, "Couldn't call \"RunCommand\" as client")

	return res
}
