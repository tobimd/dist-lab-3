package fulcrum

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"fmt"
	"strings"
	"time"
)

var (
	id int
	f  = ""

	// Contains a growing map where the keys are the names of
	// each planet registered, and the values are the
	// corresponding time vector for each one
	planetVectors = make(map[string]data.TimeVector)

	// All grpc clients stored with each entity's address used
	// as keys for the map object. `data.Client` is a pointer
	// to pb.CommuncationClient which should be what gets
	// returned by `util.SetupClient(...)`
	clients = make(map[string]*Client)
)

// Open (or create) a registry file corresponding to a specific
// planet and add the corresponding information onto it.
//
// For planet "Tatooine", city "Mos Eisley" and number of rebels
// "5", all info is stored in the following way:
//     -- Tatooine.txt: Tatooine Mos_Eisley 5
func SavePlanetData(planet string, city string, numRebels int) data.TimeVector {
	// Set the filename and the data to save
	filename := planet + ".txt"
	info := fmt.Sprintf("%s %s %d", planet, city, numRebels)

	// Write data into file
	fileExisted, err := util.WriteLines(filename, info)
	log.FailOnError(&f, err, "Couldn't write to file \"%s\"", filename)

	// Create time vector for new planet, because there was no
	// file for it before
	if !fileExisted {
		planetVectors[planet] = data.TimeVector{Time: [3]int{0, 0, 0}}
	}

	// Update time vector for that planet
	vector := planetVectors[planet]
	vector.Time[id] += 1
	planetVectors[planet] = vector

	// Return the corresponding time vector
	return planetVectors[planet]
}

func UpdatePlanetLog(command *pb.Command, planet string, city string, value interface{}) {
	filename := fmt.Sprintf("log.%s.txt", planet)
	info := ""

	if *command.Enum() != pb.Command_DELETE_CITY {
		info = fmt.Sprintf("%s %s %s %v", command.ToString(), planet, city, value)
	} else {
		info = fmt.Sprintf("%s %s %s", command.ToString(), planet, city)
	}

	_, err := util.WriteLines(filename, info)
	log.FailOnError(&f, err, "Couldn't write to file \"%s\"", filename)
}

func ReadPlanetLog(planet string) []pb.CommandParams {
	filename := fmt.Sprintf("log.%s.txt", planet)
	info := make([]pb.CommandParams, 1)

	util.ReadLines(filename,
		func(line string) bool {
			// <command> <planet> <city> [<new_value>]
			values := strings.Split(line, " ")
			info = append(info, pb.CommandParams{
				Command: pb.CommandFromString(values[0]),
			})

			return false
		},
	)

	return info
}

// Asynchronous function called only by one of the available
// fulcrums. The eventual consistency will be carried out by
// this entity, asking both neighbours to send their history and
// then merge all three states. The final version will be sent
// as a response.
func SyncWithEventualConsistency() {
	for {
		time.Sleep(time.Minute * 2)

		// Send 'RunCommand' rpc call with 'CHECK_CONSISTENCY'
		for i := 1; i < 3; i++ {
			neighbor := data.Address.FULCRUM[(id+1)%3]

			clients[neighbor].RunCommand(&pb.CommandParams{
				Command: pb.Command_CHECK_CONSISTENCY.Enum(),
			})
		}

		// Await for them to  send 'BroadcastChanges' rpc call
		// and use that to respond with the new merged history
	}
}

func Run(fulcrumId int) {
	id = fulcrumId
	f = fmt.Sprintf("fulcrum_%d.log", id)

	util.SetupServer(&f, data.Address.FULCRUM[id], &Server{})

	for i := 1; i < 3; i++ {
		neighbor := data.Address.FULCRUM[(id+i)%3]
		conn, client := util.SetupClient(&f, neighbor)
		defer conn.Close()

		clients[neighbor] = &Client{Client: client}
	}

	if id == 0 {
		go SyncWithEventualConsistency()
	}

	forever := make(chan bool)
	<-forever
}
