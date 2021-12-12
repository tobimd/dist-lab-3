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
	// Recieve fulcrum server's address

	// Send (command, planet, city, value) to fulcrum
	// Recieve time vector, associated with that planet

	var fulcrumResponse *pb.FulcrumResponse

	str := fmt.Sprintf("%v", value)
	broker := clients[data.Address.BROKER]
	addressRes := broker.RunCommand(&pb.CommandParams{
		Command: command,
	})

	fulcrum := clients[*addressRes.FulcrumRedirectAddr]
	switch command {
	case pb.Command_ADD_CITY.Enum():
		number, err := strconv.Atoi(str)
		rebelNumber := uint32(number)
		log.FailOnError(&f, err, "Failed to convert string to int (number of rebels")
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:        command,
			PlanetName:     &planet,
			CityName:       &city,
			NewNumOfRebels: &rebelNumber,
		})
	case pb.Command_UPDATE_NAME.Enum():
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:     command,
			PlanetName:  &planet,
			CityName:    &city,
			NewCityName: &str,
		})
	case pb.Command_UPDATE_NUMBER.Enum():
		number, err := strconv.Atoi(str)
		rebelNumber := uint32(number)
		log.FailOnError(&f, err, "Failed to convert string to int (number of rebels")
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:        command,
			PlanetName:     &planet,
			CityName:       &city,
			NewNumOfRebels: &rebelNumber,
		})
	case pb.Command_DELETE_CITY.Enum():
		fulcrumResponse = fulcrum.RunCommand(&pb.CommandParams{
			Command:    command,
			PlanetName: &planet,
			CityName:   &city,
		})
	default:
		log.Fatal(&f, "Received command unknown to informant")
	}
	log.Log(&f, "Received time vector: %s", fulcrumResponse.TimeVector)

}

func inter() {
	// Temp function
	var input string
	for {
		fmt.Scan(&input)
		log.Log(&f, "received input from user: %s", input)
		command, planet, city, rebels, err := util.ReadUserInput(input)
		log.FailOnError(&f, err, "failed to read user input")
		ExecuteCommand(command, planet, city, rebels)
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
	go inter()

	forever := make(chan bool)
	<-forever
}
