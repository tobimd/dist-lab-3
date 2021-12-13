package informant

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
)

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
