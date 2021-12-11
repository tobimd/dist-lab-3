package informant

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"fmt"

	"google.golang.org/grpc"
)

var (
	id      int
	f       = ""
	conns   []*grpc.ClientConn
	clients []*pb.CommunicationClient
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
	conns[0], clients[0] = util.SetupClient(&f, data.Address.BROKER)
	conns[1], clients[1] = util.SetupClient(&f, data.Address.FULCRUM[0])
	conns[2], clients[2] = util.SetupClient(&f, data.Address.FULCRUM[1])
	conns[3], clients[3] = util.SetupClient(&f, data.Address.FULCRUM[2])

}
