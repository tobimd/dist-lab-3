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
		cmd, planet, city, _ := util.ReadUserInput(&f, "> ")

		command := pb.CommandParams{
			Command:    cmd,
			PlanetName: &planet,
			CityName:   &city,
		}

		serverResponse := leia_client.RunCommand(&command)

		log.Log(&f, "Broker Response: %v", *serverResponse)

		///////////////////////  Uncomment when broker responds with something... //////////////////////////
		// numRebels := *serverResponse.NumOfRebels
		// timeVector := serverResponse.TimeVector.Time
		// fulcrumAddr := *serverResponse.FulcrumRedirectAddr

		numRebels := 20
		timeVector := data.TimeVector{1, 1, 1}
		fulcrumAddr := data.Address.FULCRUM[2]

		log.Log(&f, "%v, %v, %v", numRebels, timeVector, fulcrumAddr)

		event := data.CommandHistory{
			Command:        *cmd,
			Planet:         planet,
			FulcrumAddress: fulcrumAddr,
			TimeVector:     timeVector,
		}

		// save event in command history
		cmdHistory = append(cmdHistory, event)

		log.Print(&f, "Number of Rebels: %d\n", numRebels)

		log.Log(&f, "History to date: %+v", cmdHistory)
	}
}
