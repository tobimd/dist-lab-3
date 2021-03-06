package leia

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
)

var (
	f = "leia.log"

	cmdHistory = make([]data.CommandHistory, 0)
)

func Run() {

	conn, client := util.SetupClient(&f, data.Address.BROKER)
	defer conn.Close()

	leia_client := &Client{Client: client}

	// main client loop
	for {

		cmd, planet, city, _ := util.ReadUserInput(&f, "Leia $ ")

		command := pb.CommandParams{
			Command:    cmd,
			PlanetName: &planet,
			CityName:   &city,
		}

		serverResponse := leia_client.RunCommand(&command)

		numRebels := serverResponse.GetNumOfRebels()
		timeVector := serverResponse.TimeVector.GetTime()
		fulcrumAddr := serverResponse.GetFulcrumRedirectAddr()

		event := data.CommandHistory{
			Command:        *cmd,
			Planet:         planet,
			FulcrumAddress: fulcrumAddr,
			TimeVector:     timeVector,
		}

		// save event in command history
		cmdHistory = append(cmdHistory, event)

		log.Print(&f, "Número de rebeldes: %d\n", numRebels)
		log.Log(&f, "History to date: %+v", cmdHistory)
	}
}
