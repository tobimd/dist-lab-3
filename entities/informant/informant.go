package informant

import (
	"dist/pb"
	"fmt"
)

var (
	id int
	f  = ""
)

func ExecuteCommand(command *pb.Command, planet string, city string, value interface{}) {
	// Send (command, planet, city, value) to Broker
	// Recieve fulcrum server's address

	// Send (command, planet, city, value) to fulcrum
	// Recieve time vector, associated with that planet
}

func Run(informantId int) {
	id = informantId
	f = fmt.Sprintf("informant_%d.log", id)
}
