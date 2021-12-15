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
		city := command.GetCityName()

		if CityDoesExistOn(planet, city) {
			timeVector := pb.TimeVector{Time: []uint32{0, 0, 0}}
			return &pb.FulcrumResponse{TimeVector: &timeVector}, nil
		}

		SavePlanetData(planet, command.GetCityName(), int(command.GetNewNumOfRebels()), "", StoreMethod.Append)
		UpdatePlanetLog(command.Command, planet, command.GetCityName(), command.GetNewNumOfRebels())

	case pb.Command_UPDATE_NAME:

		log.Log(&f, "Updating city: %v", command)

		planet = command.GetPlanetName()
		city := command.GetCityName()

		if !CityDoesExistOn(planet, city) {
			timeVector := pb.TimeVector{Time: []uint32{0, 0, 0}}
			return &pb.FulcrumResponse{TimeVector: &timeVector}, nil
		}

		log.Log(&f, "new city name: %s", command.GetNewCityName())

		SavePlanetData(planet, command.GetCityName(), 0, command.GetNewCityName(), StoreMethod.Update)
		UpdatePlanetLog(command.Command, planet, command.GetCityName(), command.GetNewCityName())

	case pb.Command_UPDATE_NUMBER:

		log.Log(&f, "Updating rebels: %v", command)

		planet = command.GetPlanetName()
		city := command.GetCityName()

		if !CityDoesExistOn(planet, city) {
			timeVector := pb.TimeVector{Time: []uint32{0, 0, 0}}
			return &pb.FulcrumResponse{TimeVector: &timeVector}, nil
		}

		SavePlanetData(planet, command.GetCityName(), int(command.GetNewNumOfRebels()), "", StoreMethod.Update)
		UpdatePlanetLog(command.Command, planet, command.GetCityName(), command.GetNewNumOfRebels())

	case pb.Command_DELETE_CITY:

		log.Log(&f, "Deleting city: %v", command)

		planet = command.GetPlanetName()
		city := command.GetCityName()

		if !CityDoesExistOn(planet, city) {
			timeVector := pb.TimeVector{Time: []uint32{0, 0, 0}}
			return &pb.FulcrumResponse{TimeVector: &timeVector}, nil
		}
		SavePlanetData(planet, command.GetCityName(), int(command.GetNewNumOfRebels()), "", StoreMethod.Delete)
		UpdatePlanetLog(command.Command, planet, command.GetCityName(), command.GetNewNumOfRebels())

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

	return &response, nil
}

// Should only be called to server Fulcrum with id 0, because
// it's the only one assign to merge all three action histories.
//
// The lower the id, the bigger priority it has over merge
// conflicts, so that way we can avoid having to recieve both of
// the neighbor's histories and just deal with them separately.
func (s *Server) BroadcastChanges(ctx context.Context, fulcrumHistory *pb.FulcrumHistory) (*pb.FulcrumHistory, error) {
	log.Log(&f, "[server:BroadcastChanges] Called")

	// If length is 0, then it means fulcrum should respond with
	// it's own history
	if len(fulcrumHistory.GetHistory()) == 0 {
		return &pb.FulcrumHistory{History: GetHistory()}, nil

		// Otherwise, update history
	} else {
		SetHistory(fulcrumHistory.GetHistory())

		return &pb.FulcrumHistory{History: []*pb.CommandParams{}}, nil
	}
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

func (c *Client) BroadcastChanges(history []*pb.CommandParams) []*pb.CommandParams {
	client := *((*c).Client)
	log.Log(&f, "[client:BroadcastChanges] Called")

	ctx, cancel := util.GetContext()
	defer cancel()

	res, err := client.BroadcastChanges(ctx, &pb.FulcrumHistory{History: history})
	log.FailOnError(&f, err, "Couldn't call \"BroadcastChanges\" as entity Fulcrum (id=%d)", id)

	return res.GetHistory()
}
