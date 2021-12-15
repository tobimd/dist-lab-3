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

	curr_id, allVectors := GetLatestVector(command.PlanetName)

	log.Log(&f, "[server:RunCommand] Most recent TimeVector for planet %s found in server",
		command.GetPlanetName(),
	)

	log.Log(&f, "%d", curr_id)

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
		candidates := make([]string, 0)

		informantVector := command.LastTimeVector.GetTime()
		log.Log(&f, "allVectors=%v, informantVector=%v", allVectors, informantVector)

		for fulcrumId, vector := range allVectors {

			if vector.GreaterThanOrEqual(informantVector) {
				candidates = append(candidates, data.Address.FULCRUM[fulcrumId])
			}
		}

		if len(candidates) == 0 {
			candidates = data.Address.FULCRUM
		}

		// redirect client to fulcrum server address with most recent writes or current one
		fulcrumAddr := candidates[util.RandInt(0, len(candidates))]
		log.Log(&f, "Length of candidates for redirecting informant: %d", len(candidates))

		peerMd, ok := peer.FromContext(ctx)
		peerAddr := "NO_ADDR"

		if ok {
			peerAddr = peerMd.Addr.String()
		}

		log.Log(&f, "[server:RunCommand] Redirecting informant %v to %v", peerAddr, fulcrumAddr)

		return &pb.FulcrumResponse{FulcrumRedirectAddr: &fulcrumAddr}, nil
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
func GetLatestVector(planet *string) (int, map[int]data.TimeVector) {
	var wg sync.WaitGroup
	var fulcrumNum = 3
	var timeVectors = make(map[int]data.TimeVector)
	var safeTimeVectors = sync.Map{}

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
			safeTimeVectors.Store(i, resp.TimeVector.GetTime())
			log.Log(&f, "response vectors %v", resp.TimeVector.Time)
		}(i)
	}
	wg.Wait()

	// pass vals from concurrent-safe map to regular one for compliance with rest of code,
	// also handle border cases when fulcrum servers have just started running
	for i := 0; i < fulcrumNum; i++ {
		tm, _ := safeTimeVectors.Load(i)

		// assert is correct TimeVector type and then get underlying value (zero val for type on error)
		dataTm, ok := tm.([]uint32)

		if !ok || len(dataTm) == 0 {
			timeVectors[i] = data.TimeVector{0, 0, 0}
		} else {
			timeVectors[i] = data.TimeVector{dataTm[0], dataTm[1], dataTm[2]}
		}
	}

	log.Log(&f, "TimeVectors: %v", timeVectors)

	randIds := rand.Perm(fulcrumNum)
	index := 0

	// get latest time vector between all three
	if timeVectors[randIds[0]].GreaterThanOrEqual(timeVectors[randIds[1]]) {
		if !timeVectors[randIds[0]].GreaterThanOrEqual(timeVectors[randIds[2]]) {
			index = 2
		}
	} else {
		if timeVectors[randIds[1]].GreaterThanOrEqual(timeVectors[randIds[2]]) {
			index = 1
		} else {
			index = 2
		}
	}

	return randIds[index], timeVectors
}
