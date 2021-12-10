package fulcrum

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

func (s *Server) RunCommand(ctx context.Context, command *pb.CommandParams) (*pb.Empty, error) {
	log.Log(&f, "[server:RunCommand] Called with argument: command=\"%v\"", command.String())

	switch command.Command.Enum() {

	case pb.Command_ADD_CITY.Enum():

	case pb.Command_UPDATE_NAME.Enum():

	case pb.Command_UPDATE_NUMBER.Enum():

	case pb.Command_DELETE_CITY.Enum():

	case pb.Command_CHECK_CONSISTENCY.Enum():
		// When recieving this command, gather all history from
		// logs and call `BroadcastChanges` to sender with defer

		for planet, _ /*vector*/ := range planetVectors {
			/*info :=*/ ReadPlanetLog(planet)

		}

	}

	return &pb.Empty{}, nil
}

// **** CLIENT FUNCTIONS ****
type Client data.GrpcClient

func (c *Client) RunCommand(command *pb.CommandParams) *pb.Empty {
	client := *((*c).Client)
	log.Log(&f, "[client:RunCommand] Called with argument: command=\"%v\"", command.String())

	ctx, cancel := util.GetContext()
	defer cancel()

	res, err := client.RunCommand(ctx, command)
	log.FailOnError(&f, err, "Couldn't call \"RunCommand\" as client")

	return res
}
