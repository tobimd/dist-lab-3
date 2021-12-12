package leia

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"fmt"
)

var (
	f = ""
)

func Run() {
	f = "leia.log"

	conn, client := util.SetupClient(&f, data.Address.BROKER)
	defer conn.Close()

	leia_client := &Client{Client: client}

	// main client loop
	for {
		cmd, planet, city, _, err := util.ReadUserInput(&f, "> ")

		log.FailOnError(&f, err, err.Error())

		command := pb.CommandParams{
			Command:    cmd,
			PlanetName: &planet,
			CityName:   &city,
		}

		serverResponse := leia_client.RunCommand(&command)

		numRebels := serverResponse.NumOfRebels
		timeVector := serverResponse.TimeVector

		fmt.Printf("Number of Rebels: %d\n", *numRebels)
	}

	forever := make(chan bool)
	<-forever
}
