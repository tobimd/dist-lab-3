package leia

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
)

var (
	f = ""

	cmdHistory = make([]data.CommandHistory, 1)
)

func Run() {
	f = "leia.log"

	conn, client := util.SetupClient(&f, data.Address.BROKER)
	defer conn.Close()

	leia_client := &Client{Client: client}

	// main client loop
	for {
		cmd, planet, city, _ := util.ReadUserInput(&f, "> ")

		command := pb.CommandParams{
			Command:    cmd,
			PlanetName: &planet,
			CityName:   &city,
		}

		serverResponse := leia_client.RunCommand(&command)

		numRebels := serverResponse.NumOfRebels
		timeVector := serverResponse.TimeVector.Time
		fulcrumAddr := serverResponse.FulcrumRedirectAddr

		event := data.CommandHistory{
			Command:        *cmd,
			Planet:         planet,
			FulcrumAddress: *fulcrumAddr,
			TimeVector:     timeVector,
		}

		// save event in command history
		cmdHistory = append(cmdHistory, event)

		log.Print(&f, "Number of Rebels: %d\n", *numRebels)
	}
}
