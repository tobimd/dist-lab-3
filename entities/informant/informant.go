package informant

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"fmt"
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

}

func inter() {
	// Temp function
	var input string
	for {
		fmt.Scan(&input)
		log.Log(&f, "received input from user: %s", input)
		command, planet, city, rebels, err := util.ReadUserInput(input)
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

	forever := make(chan bool)
	<-forever
}
