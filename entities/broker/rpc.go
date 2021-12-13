package broker

import (
	"context"
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"math/rand"
	"sync"

	"google.golang.org/grpc/peer"
)

type Server struct {
	pb.UnimplementedCommunicationServer
}

type Client data.GrpcClient

// **** SERVER FUNCTIONS ****
// recieve command and take action according to type of node
func (s *Server) RunCommand(ctx context.Context, command *pb.CommandParams) (*pb.FulcrumResponse, error) {
	log.Log(&f, "[server:RunCommand] Called with argument: command=\"%v\"", command.String())

	curr_id := GetLatestVector(command.PlanetName)

	log.Log(&f, "[server:RunCommand] Most recent TimeVector for planet %s found in server %d",
		*command.PlanetName,
		curr_id,
	)

	switch *command.Command {

	case pb.Command_GET_NUMBER_REBELS:
		// proxy message between leia and fulcrum servers, getting either the most recent version or current one
		log.Log(&f, "[server:RunCommand] Proxying message for Leia")

		client := fulcrum_client[curr_id]
		cmdResponse := client.RunCommand(command)
		cmdResponse.FulcrumRedirectAddr = &data.Address.FULCRUM[curr_id]

		log.Log(&f, "[server:RunCommand] Proxying response for Leia")

		return cmdResponse, nil

	default:
		// redirect client to fulcrum server address with most recent writes or current one
		randAddr := data.Address.FULCRUM[curr_id]

		peerMd, ok := peer.FromContext(ctx)
		peerAddr := "NO_ADDR"

		if ok {
			peerAddr = peerMd.Addr.String()
		}

		log.Log(&f, "[server:RunCommand] Redirecting informant %v to %v", peerAddr, randAddr)

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

// get all current timevectors from all fulcrum servers
func GetLatestVector(planet *string) int {
	var wg sync.WaitGroup
	var fulcrumNum = 3
	var timevectors = make(map[int]data.TimeVector)

	// get all timevectors from all three fulcrums for a given planet
	for i := 0; i < fulcrumNum; i++ {
		wg.Add(i)

		c := fulcrum_client[i]

		go func(i int) {
			defer wg.Done()

			cmd := pb.CommandParams{
				Command:    pb.Command_GET_PLANET_TIME.Enum(),
				PlanetName: planet,
			}

			resp := c.RunCommand(&cmd)
			timevectors[i] = resp.TimeVector.Time
		}(i)
	}

	wg.Wait()

	// handle border cases when fulcrum servers have just started running
	for i := 0; i < fulcrumNum; i++ {
		if len(timevectors[i]) == 0 {
			timevectors[i] = data.TimeVector{0, 0, 0}
		}
	}

	randIds := rand.Perm(fulcrumNum)

	// get latest time vector between all three
	if timevectors[randIds[0]].GreaterThanOrEqual(timevectors[randIds[1]]) {
		if timevectors[randIds[0]].GreaterThanOrEqual(timevectors[randIds[2]]) {
			return randIds[0]
		}
		return randIds[2]
	} else {
		if timevectors[randIds[1]].GreaterThanOrEqual(timevectors[randIds[2]]) {
			return randIds[1]
		}
		return randIds[2]
	}
}
