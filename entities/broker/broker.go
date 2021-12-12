package broker

import (
	"dist/common/data"
	"dist/common/util"
)

var (
	// log's filename
	f = ""

	// maps fulcrum ids to client objects for comms with fulcrum servers
	fulcrum_client = make(map[int]*Client)
)

func Run() {
	f = "broker.log"

	util.SetupServer(&f, data.Address.BROKER, &Server{})

	for i := 0; i < 3; i++ {
		conn, client := util.SetupClient(&f, data.Address.FULCRUM[i])
		defer conn.Close()

		fulcrum_client[i] = &Client{Client: client}
	}

	forever := make(chan bool)
	<-forever
}
