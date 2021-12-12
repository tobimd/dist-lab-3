package broker

import (
	"context"
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"math/rand"
)

type Server struct {
	pb.UnimplementedCommunicationServer
}

type Client data.GrpcClient

// **** SERVER FUNCTIONS ****
// recieve command and take action according to type of node
func (s *Server) RunCommand(ctx context.Context, command *pb.CommandParams) (*pb.FulcrumResponse, error) {
	randId := rand.Int() % 3

	switch command.Command.Enum() {

	case pb.Command_GET_NUMBER_REBELS.Enum():
		// proxy message between leia and fulcrum servers
		client := fulcrum_client[randId]
		cmdResponse := client.RunCommand(command)
		return cmdResponse, nil

	default:
		// redirect client to fulcrum server address
		randAddr := data.Address.FULCRUM[randId]
		return &pb.FulcrumResponse{FulcrumRedirectAddr: &randAddr}, nil
	}

}

// **** CLIENT FUNCTIONS ****
// send command to fulcrum servers when Leia queries curr number of rebels
func (c *Client) RunCommand(command *pb.CommandParams) *pb.FulcrumResponse {
	client := *((*c).Client)
	log.Log(&f, "[client:RunCommand] Called with argument: command=\"%v\"", command.String())

	ctx, cancel := util.GetContext()
	defer cancel()

	res, err := client.RunCommand(ctx, command)
	log.FailOnError(&f, err, "Couldn't call \"RunCommand\" as entity Broker")

	return res
}
