package fulcrum

import (
	"context"
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"errors"
)

// **** SERVER FUNCTIONS ****
type Server struct {
	pb.UnimplementedCommunicationServer
}

func (s *Server) RunCommand(ctx context.Context, command *pb.CommandParams) (*pb.FulcrumResponse, error) {
	log.Log(&f, "[server:RunCommand] Called with argument: command=\"%v\"", command.String())

	var planet string
	var response pb.FulcrumResponse
	switch *command.Command {

	case pb.Command_ADD_CITY:

		log.Log(&f, "Adding city: %v", command)

		planet = command.GetPlanetName()

		SavePlanetData(planet, command.GetCityName(), int(command.GetNewNumOfRebels()), "", StoreMethod.Create)
		UpdatePlanetLog(command.Command, planet, command.GetCityName(), command.GetNewNumOfRebels())

	case pb.Command_UPDATE_NAME:

		log.Log(&f, "Updating city: %v", command)

		planet = command.GetPlanetName()

		if command.NumOfRebels == nil {
			log.Print(&f, "nil nil")
		}

		log.Log(&f, "new city name: %s", command.GetNewCityName())

		SavePlanetData(planet, command.GetCityName(), 0, command.GetNewCityName(), StoreMethod.Update)
		UpdatePlanetLog(command.Command, planet, command.GetCityName(), command.GetNewCityName())

	case pb.Command_UPDATE_NUMBER:

		log.Log(&f, "Updating rebels: %v", command)

		planet = command.GetPlanetName()

		SavePlanetData(planet, command.GetCityName(), int(command.GetNewNumOfRebels()), "", StoreMethod.Update)
		UpdatePlanetLog(command.Command, planet, command.GetCityName(), command.GetNewNumOfRebels())

	case pb.Command_DELETE_CITY:

		log.Log(&f, "Deleting city: %v", command)

		planet = command.GetPlanetName()

		SavePlanetData(planet, command.GetCityName(), int(command.GetNewNumOfRebels()), "", StoreMethod.Delete)
		UpdatePlanetLog(command.Command, planet, command.GetCityName(), command.GetNewNumOfRebels())

	case pb.Command_CHECK_CONSISTENCY:
		// When recieving this command, gather all history from
		// logs and call `BroadcastChanges` to sender with defer
		info := make([]*pb.CommandParams, 0)

		for planet, vector := range planetVectors {
			// Append all updates done to a planet
			info = append(info, ReadPlanetLog(planet)...)

			// For the last Command, set LastTimeVector with these
			info[len(info)-1].LastTimeVector = &pb.TimeVector{
				Time: vector,
			}
		}

		// Now that all history is saved within 'info', call
		// 'BroadcastChanges' to fulcrum 0
		defer clients[data.Address.FULCRUM[0]].BroadcastChanges(&pb.FulcrumHistory{
			History: info,
		})

	case pb.Command_GET_PLANET_TIME:
		planet = command.GetPlanetName()

		log.Log(&f, "Returning planet %v TimeVector: %v", planet, planetVectors[planet])

	case pb.Command_GET_NUMBER_REBELS:

		planet = command.GetPlanetName()

		rebelNumber := ReadCityData(planet, command.GetCityName())

		timeVector := pb.TimeVector{Time: planetVectors[planet]}
		return &pb.FulcrumResponse{TimeVector: &timeVector, NumOfRebels: &rebelNumber}, nil

	default:
		log.Log(&f, "Unknown command: %v", command)

		return &pb.FulcrumResponse{}, errors.New("unknown command")
	}

	timevector := pb.TimeVector{Time: planetVectors[planet]}
	response = pb.FulcrumResponse{TimeVector: &timevector}

	log.Print(&f, "PlanetVectors before sending: %+v", planetVectors)
	log.Print(&f, "Response before sending: %+v", response)

	return &response, nil
}

// Should only be called to server Fulcrum with id 0, because
// it's the only one assign to merge all three action histories.
//
// The lower the id, the bigger priority it has over merge
// conflicts, so that way we can avoid having to recieve both of
// the neighbor's histories and just deal with them separately.
func (s *Server) BroadcastChanges(ctx context.Context, history *pb.FulcrumHistory) (*pb.FulcrumHistory, error) {
	log.Log(&f, "[server:BroadcastChanges] Called")

	// TODO: Check if it makes sense to use this logic inside
	// an anonymous function (and also that it's async)
	go func() {
		if len(neighborHistories.hist1) == 0 {
			neighborHistories.hist1 = history.GetHistory()

		} else {
			neighborHistories.hist2 = history.GetHistory()
			neighborChangesCh <- true
		}
	}()

	// Stop until second fulcrum sends history, making the
	// channel 'neighborChangesCh' recieve a value
	<-neighborChangesCh

	newHistory := MergeHistories()
	neighborHistories.hist1 = []*pb.CommandParams{}
	neighborHistories.hist2 = []*pb.CommandParams{}

	return &pb.FulcrumHistory{History: newHistory}, nil
}

// **** CLIENT FUNCTIONS ****
type Client data.GrpcClient

func (c *Client) RunCommand(command *pb.CommandParams) *pb.FulcrumResponse {
	client := *((*c).Client)
	log.Log(&f, "[client:RunCommand] Called with argument: command=\"%v\"", command.String())

	ctx, cancel := util.GetContext()
	defer cancel()

	res, err := client.RunCommand(ctx, command)
	log.FailOnError(&f, err, "Couldn't call \"RunCommand\" as entity Fulcrum (id=%d)", id)

	return res
}

func (c *Client) BroadcastChanges(history *pb.FulcrumHistory) {
	client := *((*c).Client)
	log.Log(&f, "[client:BroadcastChanges] Called")

	ctx, cancel := util.GetContext()
	defer cancel()

	res, err := client.BroadcastChanges(ctx, history)
	log.FailOnError(&f, err, "Couldn't call \"BroadcastChanges\" as entity Fulcrum (id=%d)", id)

	// The returned value should be considered as the new
	// history for each planet. Because fulcrum with id=0 did
	// the merging strategy on all three histories and returned
	// the new values

	currPlanet := ""

	for _, history := range res.GetHistory() {
		planet := history.GetPlanetName()
		city := history.GetCityName()
		numRebels := int(history.GetNumOfRebels())

		method := StoreMethod.Create
		if planet != currPlanet {
			method = StoreMethod.Rewrite
			currPlanet = planet
		}

		SavePlanetData(planet, city, numRebels, "", method)
	}
}
