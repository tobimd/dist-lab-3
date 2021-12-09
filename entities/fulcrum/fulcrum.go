package fulcrum

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"fmt"
)

var (
	id int
	f  = ""

	// Contains a growing map where the keys are the names of
	// each planet registered, and the values are the
	// corresponding time vector for each one
	planetTimeVectors = make(map[string]TimeVector)
)

type TimeVector struct {
	times []int
}

// Open (or create) a registry file corresponding to a specific
// planet and add the corresponding information onto it.
//
// For planet "Tatooine", city "Mos Eisley" and number of rebels
// "5", all info is stored in the following way:
//     -- Tatooine.txt: Tatooine Mos_Eisley 5
func SavePlanetData(planet string, city string, numRebels int) TimeVector {
	// Set the filename and the data to save
	filename := planet + ".txt"
	info := fmt.Sprintf("%s %s %d", planet, city, numRebels)

	// Write data into file
	fileExisted, err := util.WriteLines(filename, info)
	log.FailOnError(&f, err, "Couldn't write to file \"%s\"", filename)

	// Create time vector for new planet, because there was no
	// file for it before
	if !fileExisted {
		planetTimeVectors[planet] = TimeVector{times: make([]int, 3)}
	}

	// Update time vector for that planet
	planetTimeVectors[planet].times[id] += 1

	// Return the corresponding time vector
	return planetTimeVectors[planet]
}

func UpdatePlanetLog(command data.CommandEnum, planet string, city string, value interface{}) {
	filename := fmt.Sprintf("log.%s.txt", planet)
	info := ""

	if command.Enum != data.Command.DELETE_CITY.Enum {
		info = fmt.Sprintf("%s %s %s %v", command.Name, planet, city, value)
	} else {
		info = fmt.Sprintf("%s %s %s", command.Name, planet, city)
	}

	_, err := util.WriteLines(filename, info)
	log.FailOnError(&f, err, "Couldn't write to file \"%s\"", filename)
}

func Run(fulcrumId int) {
	id = fulcrumId
	f = fmt.Sprintf("fulcrum_%d.log", id)

	util.SetupServer(&f, data.Address.FULCRUM[id], &server{})

	conn1, client1 := util.SetupClient(&f, data.Address.FULCRUM[(id+1)%3])
	conn2, client2 := util.SetupClient(&f, data.Address.FULCRUM[(id+2)%3])

	defer conn1.Close()
	defer conn2.Close()

	log.Print(&f, "About to call both \"HelloTest\"s for both fulcrums")
	res1 := HelloTest(client1, &pb.Empty{})
	res2 := HelloTest(client2, &pb.Empty{})

	log.Print(&f, "Successfully recieved both results as res1=\"%s\" and res2=\"%s\"", res1.String(), res2.String())

	forever := make(chan bool)
	<-forever
}
