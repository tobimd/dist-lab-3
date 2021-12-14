package informant

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"fmt"
	"strconv"
)

var (
	id            int
	f             = ""
	clients       = make(map[string]*Client)
	planetHistory = make(map[data.Planet][]data.CommandHistory)
)

func ExecuteCommand(command *pb.Command, planet string, city string, value interface{}) {
	// Send (command, planet, city, value) to Broker
	// Receive fulcrum server's address
	fn := "<ExecuteCommand>"
	var fulcrumResponse *pb.FulcrumResponse

	timeVector := []uint32{0, 0, 0}
	if len(planetHistory[planet]) == 0 {
		log.Log(&f, "%s Latest time vector for planet %s: %v", fn, planet, planetHistory[planet])

	} else {
		log.Log(&f, "%s Latest time vector for planet %s: %v", fn, planet, planetHistory[planet][len(planetHistory[planet])-1])
		timeVector = planetHistory[planet][len(planetHistory[planet])-1].TimeVector

	}

	str := fmt.Sprintf("%v", value)
	broker := clients[data.Address.BROKER]
	fulcrumAddress := broker.RunCommand(&pb.CommandParams{
		Command:        command,
		LastTimeVector: &pb.TimeVector{Time: timeVector},
		PlanetName:     &planet,
	}).FulcrumRedirectAddr

	log.Log(&f, "%s Received fulcrum address from broker: %s", fn, *fulcrumAddress)

	// Send (command, planet, city, value) to fulcrum
	// Receive time vector, associated with that planet
	fulcrum := clients[*fulcrumAddress]

	log.Log(&f, "%s Sending command \"%s\" to fulcrum", fn, command)

	switch *command {
	// Sends necessary parameters to command executed

	case pb.Command_ADD_CITY:
		number, err := strconv.Atoi(str)
		rebelNumber := uint32(number)
		log.FailOnError(&f, err, "%s Failed to convert string to int (number of rebels)", fn)
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:        command,
			PlanetName:     &planet,
			CityName:       &city,
			NewNumOfRebels: &rebelNumber,
			LastTimeVector: &pb.TimeVector{Time: timeVector},
		})
	case pb.Command_UPDATE_NAME:
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:        command,
			PlanetName:     &planet,
			CityName:       &city,
			NewCityName:    &str,
			LastTimeVector: &pb.TimeVector{Time: timeVector},
		})
	case pb.Command_UPDATE_NUMBER:
		number, err := strconv.Atoi(str)
		rebelNumber := uint32(number)
		log.FailOnError(&f, err, "$s Failed to convert string to int (number of rebels)", fn)
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:        command,
			PlanetName:     &planet,
			CityName:       &city,
			NewNumOfRebels: &rebelNumber,
			LastTimeVector: &pb.TimeVector{Time: timeVector},
		})
	case pb.Command_DELETE_CITY:
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:        command,
			PlanetName:     &planet,
			CityName:       &city,
			LastTimeVector: &pb.TimeVector{Time: timeVector},
		})
	default:
		// Stops if command was not recognized
		log.Log(&f, "%s Received command unknown to informant", fn)
	}
	log.Log(&f, "%s Received time vector: %v", fn, fulcrumResponse.GetTimeVector())

	info := new(data.CommandHistory)
	info.Command = *command
	info.City = city
	info.FulcrumAddress = *fulcrumAddress
	//DELETE when fuclrumResponse contains non null TimeVector
	if fulcrumResponse.TimeVector != nil {
		info.TimeVector = fulcrumResponse.TimeVector.GetTime()

	} else {
		info.TimeVector = []uint32{1, 1, 1}
		log.Log(&f, "%s Using default time vector", fn)

	}
	planetHistory[planet] = append(planetHistory[planet], *info)

}

func ConsoleInteraction() {
	// Interface between user & program
	fmt.Print("Saludos Informante\n")

	for {
		log.Log(&f, "<ConsoleInteraction> Requesting user input")
		command, planet, city, value := util.ReadUserInput(&f, "Ingresa el comando y argumentos que quieres usar:")
		log.Log(&f, "<ConsoleInteraction> Parsed command: %s %s %s %s", command, planet, city, value)
		if command != nil {
			ExecuteCommand(command, planet, city, value)
		}
	}
}

func Run(informantId int) {
	id = informantId
	f = fmt.Sprintf("informant_%d.log", id)

	// Setup all four clients with broker and three fulcrums
	conn, client := util.SetupClient(&f, data.Address.BROKER)
	defer conn.Close()

	clients[data.Address.BROKER] = &Client{Client: client}

	for i := 0; i < 3; i++ {
		fulcrumAddress := data.Address.FULCRUM[i]
		conn, client = util.SetupClient(&f, fulcrumAddress)
		defer conn.Close()

		clients[fulcrumAddress] = &Client{Client: client}
	}
	go ConsoleInteraction()

	forever := make(chan bool)
	<-forever
}
