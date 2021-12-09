package fulcrum

import (
	"context"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
)

// **** SERVER FUNCTIONS ****
type server struct {
	pb.UnimplementedCommunicationServer
}

func (s *server) HelloTest(ctx context.Context, empty *pb.Empty) (*pb.Empty, error) {
	log.Print(&f, "#server:HelloTest# recieved a hello from another fulcrum server")

	return &pb.Empty{}, nil
}

// **** CLIENT FUNCTIONS ****
type client struct {
}

func HelloTest(client *pb.CommunicationClient, in *pb.Empty) *pb.Empty {
	ctx, cancel := util.GetContext()
	defer cancel()

	log.Log(&f, "#client:HelloTest# calling HelloTest to server")
	res, err := (*client).HelloTest(ctx, in)
	log.FailOnError(&f, err, "Recieved error from server when calling \"HelloTest\"")

	log.Log(&f, "#client:HelloTest# server responded with \"%s\"", res.String())
	log.Log(&f, "#client:HelloTest# returning that response")

	return res
}
