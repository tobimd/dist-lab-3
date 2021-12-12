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
	id      int
	f       = ""
	clients = make(map[string]*Client)
)

func ExecuteCommand(command *pb.Command, planet string, city string, value interface{}) {
	// Send (command, planet, city, value) to Broker
	// Receive fulcrum server's address
	fn := "<ExecuteCommand>"
	var fulcrumResponse *pb.FulcrumResponse

	str := fmt.Sprintf("%v", value)
	broker := clients[data.Address.BROKER]
	addressRes := broker.RunCommand(&pb.CommandParams{
		Command: command,
	}).FulcrumRedirectAddr

	log.Log(&f, "%s Received fulcrum address from broker: %s", fn, *addressRes)

	// Send (command, planet, city, value) to fulcrum
	// Receive time vector, associated with that planet
	fulcrum := clients[*addressRes]

	log.Log(&f, "%s Sending command \"%s\" to fulcrum", fn, command)
	switch *command {

	case pb.Command_ADD_CITY:
		number, err := strconv.Atoi(str)
		rebelNumber := uint32(number)
		log.FailOnError(&f, err, "%s Failed to convert string to int (number of rebels)", fn)
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:        command,
			PlanetName:     &planet,
			CityName:       &city,
			NewNumOfRebels: &rebelNumber,
		})
	case pb.Command_UPDATE_NAME:
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:     command,
			PlanetName:  &planet,
			CityName:    &city,
			NewCityName: &str,
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
		})
	case pb.Command_DELETE_CITY:
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:    command,
			PlanetName: &planet,
			CityName:   &city,
		})
	default:
		log.Log(&f, "%s Received command unknown to informant", fn)
	}
	log.Log(&f, "%s Received time vector: %s", fn, fulcrumResponse.TimeVector)

}

func ConsoleInteraction() {
	// Interface between user & program
	fmt.Print("Saludos Informante\n")

	for {
		log.Log(&f, "<ConsoleInteraction> Requesting user input")
		command, planet, city, value := util.ReadUserInput(&f, "Ingresa el comando y argumentos que quieres usar:")
		log.Log(&f, "<ConsoleInteraction> Parsed command: %s %s %s %s", command, planet, city, value)
		ExecuteCommand(command, planet, city, value)
	}
}

func Run(informantId int) {
	id = informantId
	f = fmt.Sprintf("informant_%d.log", id)
	util.SetupServer(&f, data.Address.INFORMANT[id], &Server{})

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
