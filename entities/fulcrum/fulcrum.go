package fulcrum

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"fmt"
	"strconv"
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

	// When using 'SavePlanetData', determine wether to append
	// to the planet's history file, or rewrite it completely
	// which is used when merging.
	StoreMethod = struct{ Create, Update, Delete, Rewrite Method }{Create: 0, Update: 1, Delete: 2, Rewrite: 3}

	// When leader fulcrum (id=0) tells neighbors to send their
	// history to it, the server function BroadcastChanges will
	// be called on leader and the response should be the merge
	// from all three histories. To solve that issue, this
	// channel will recieve boolean true when both histories
	// have been stored
	neighborChangesCh = make(chan bool)
	neighborHistories = struct{ hist1, hist2 []*pb.CommandParams }{}
)

type Method uint8

// Open (or create) a registry file corresponding to a specific
// planet and add the corresponding information onto it.
//
// For planet "Tatooine", city "Mos Eisley" and number of rebels
// "5", all info is stored in the following way:
//     -- Tatooine.txt: Tatooine Mos_Eisley 5
func SavePlanetData(planet string, city string, numRebels int, newCityName string, storeMethod Method) data.TimeVector {
	log.Log(&f, "<SavePlanetData(planet:\"%s\", city:\"%s\", numRebels:%d, newCityName:%s , storeMethod:\"%v\")>", planet, city, numRebels, newCityName, storeMethod)

	// Set the filename and the data to save
	filename := planet + ".txt"
	var info string

	if newCityName == "" {
		info = fmt.Sprintf("%s %s %d", planet, city, numRebels)
	} else {
		info = fmt.Sprintf("%s %s", planet, newCityName)
	}

	if storeMethod == StoreMethod.Rewrite {
		log.Log(&f, "<SavePlanetData> Because storeMethod is Rewrite, delete file and set storeMethod to Create")
		util.DeleteFile(filename)
		storeMethod = StoreMethod.Create
	}

	var err error

	switch storeMethod {
	case StoreMethod.Create:
		// Write data into file
		err = util.WriteLines(filename, false, info)
		log.Log(&f, "<SavePlanetData> storeMethod is Create, written to file")

	case StoreMethod.Update:
		replacedLine := ""
		err = util.ReplaceLines(filename, func(line string) string {
			values := strings.Split(line, " ")
			log.Log(&f, "values: %v", values)
			if values[1] == city {
				replacedLine = line
				if newCityName != "" {
					return strings.Join([]string{info, values[2]}, " ")
				}

				return info
			}

			return line
		})

		log.Log(&f, "<SavePlanetData> storeMethod is Update, replacedLine=\"%s\" with \"%s\"", replacedLine, info)

	case StoreMethod.Delete:
		deletedLine := ""
		err = util.DeleteLines(filename, func(line string) bool {
			values := strings.Split(line, " ")
			if values[1] == city {
				deletedLine = line
				return false
			}

			return true
		})
		log.Log(&f, "<SavePlanetData> storeMethod is Delete, deletedLine=\"%s\"", deletedLine)
	}
	log.FailOnError(&f, err, "Couldn't write to file \"%s\"", filename)

	// Create time vector for new planet, because there was no
	// file for it before
	if _, ok := planetVectors[planet]; !ok {
		planetVectors[planet] = data.TimeVector{0, 0, 0}
	}

	// Update time vector for that planet
	planetVectors[planet][id] += 1

	log.Log(&f, "<SavePlanetData> time vector for planet is now [ %d, %d, %d ]", planetVectors[planet][0], planetVectors[planet][1], planetVectors[planet][2])
	// Return the corresponding time vector
	return planetVectors[planet]
}

func UpdatePlanetLog(command *pb.Command, planet string, city string, value interface{}) {
	log.Log(&f, "<UpdatePlanetLog(command: %v, planet:\"%s\", city:\"%s\", value: %v)>", *command, planet, city, value)

	filename := fmt.Sprintf("log.%s.txt", planet)
	info := ""

	if *command.Enum() != pb.Command_DELETE_CITY {
		info = fmt.Sprintf("%s %s %s %v", command.ToString(), planet, city, value)
	} else {
		info = fmt.Sprintf("%s %s %s", command.ToString(), planet, city)
	}

	log.Log(&f, "<UpdatePlanetLog> new line: \"%s\"", info)

	err := util.WriteLines(filename, false, info)
	log.FailOnError(&f, err, "Couldn't write to file \"%s\"", filename)
}

// Opens planetary registry for a specific planet
// and retrieves rebel data for a city.
func ReadCityData(planet string, city string) uint32 {
	filename := fmt.Sprintf("%s.txt", planet)
	var rebelNumber uint32

	util.ReadLines(filename, func(line string) bool {
		values := strings.Split(line, " ")
		if values[1] == city {
			number, _ := strconv.Atoi(values[2])
			rebelNumber = uint32(number)
			return true
		}
		return false
	})

	return rebelNumber
}

func ReadPlanetLog(planet string) []*pb.CommandParams {
	filename := fmt.Sprintf("log.%s.txt", planet)
	info := make([]*pb.CommandParams, 1)

	util.ReadLines(filename,
		func(line string) bool {
			// <command> <planet> <city> [<new_value>]
			values := strings.Split(line, " ")

			command := pb.CommandFromString(values[0])
			newName := ""
			var newNumberOfRebels uint32 = 0

			if len(values) == 4 {
				if *command == pb.Command_UPDATE_NAME {
					newName = values[3]
				} else {
					newNumberOfRebels = util.StringToUint32(values[3])
				}
			}

			info = append(info, &pb.CommandParams{
				Command:        command,
				PlanetName:     &planet,
				CityName:       &values[2],
				NewCityName:    &newName,
				NewNumOfRebels: &newNumberOfRebels,
			})

			return false
		},
	)

	return info
}

func MergeHistories() []*pb.CommandParams {
	var newHistory []*pb.CommandParams

	type fInfo struct {
		hist   []*pb.CommandParams
		vector data.TimeVector
	}
	type fHist struct{ f0, f1, f2 fInfo }

	allHistories := make(map[string]fHist, len(planetVectors))

	for i := 0; i < 3; i++ {
		switch i {
		case 0:
			for planet, vector := range planetVectors {
				hist := ReadPlanetLog(planet)

				// Get previous info if already exists in
				// allHistories
				if info, ok := allHistories[planet]; ok {
					info.f0.hist = hist
					info.f0.vector = vector

					allHistories[planet] = info

				} else {
					allHistories[planet] = fHist{
						f0: fInfo{
							hist:   hist,
							vector: vector,
						},
						f1: fInfo{hist: []*pb.CommandParams{}, vector: []uint32{0, 0, 0}},
						f2: fInfo{hist: []*pb.CommandParams{}, vector: []uint32{0, 0, 0}},
					}
				}
				log.Log(&f, "<MergeHistories> fulcrum id=0 added history for planet \"%s\", vector %d, %d, %d", planet, vector[0], vector[1], vector[2])
			}
		case 1:
			for _, cmd := range neighborHistories.hist1 {
				planet := cmd.GetPlanetName()
				vector := cmd.GetLastTimeVector().GetTime()

				if info, ok := allHistories[planet]; ok {
					info.f1.hist = append(info.f1.hist, cmd)
					info.f1.vector = vector

					allHistories[planet] = info
				} else {
					allHistories[planet] = fHist{
						f1: fInfo{
							hist:   []*pb.CommandParams{cmd},
							vector: vector,
						},
						f0: fInfo{hist: []*pb.CommandParams{}, vector: []uint32{0, 0, 0}},
						f2: fInfo{hist: []*pb.CommandParams{}, vector: []uint32{0, 0, 0}},
					}
				}

				log.Log(&f, "<MergeHistories> fulcrum id=1 added history for planet \"%s\", vector %d, %d, %d", planet, vector[0], vector[1], vector[2])

			}
		case 2:
			for _, cmd := range neighborHistories.hist2 {
				planet := cmd.GetPlanetName()
				vector := cmd.GetLastTimeVector().GetTime()

				if info, ok := allHistories[planet]; ok {
					info.f2.hist = append(info.f2.hist, cmd)
					info.f2.vector = vector

					allHistories[planet] = info
				} else {
					allHistories[planet] = fHist{
						f2: fInfo{
							hist:   []*pb.CommandParams{cmd},
							vector: vector,
						},
						f0: fInfo{hist: []*pb.CommandParams{}, vector: []uint32{0, 0, 0}},
						f1: fInfo{hist: []*pb.CommandParams{}, vector: []uint32{0, 0, 0}},
					}
				}
				log.Log(&f, "<MergeHistories> fulcrum id=2 added history for planet \"%s\", vector %d, %d, %d", planet, vector[0], vector[1], vector[2])

			}
		}
	}

	log.Log(&f, "len(allHistories)=%d", len(allHistories))
	for planet, histories := range allHistories {
		vector0 := histories.f0.vector
		vector1 := histories.f1.vector
		vector2 := histories.f2.vector

		// Once all three histories are in memory for a given
		// planet, we have to merge them into one and add them
		// to newHistory
		if vector0.GreaterThanOrEqual(vector1) && vector0.GreaterThanOrEqual(vector2) {
			// Most updated is vector0
			newHistory = append(newHistory, histories.f0.hist...)
			log.Log(&f, "<MergeHistories> vector 0 is greater, so choosing that as history for planet \"%s\"", planet)

		} else if vector1.GreaterThanOrEqual(vector0) && vector1.GreaterThanOrEqual(vector2) {
			// Most updated is vector1
			newHistory = append(newHistory, histories.f1.hist...)
			log.Log(&f, "<MergeHistories> vector 1 is greater, so choosing that as history for planet \"%s\"", planet)

		} else {
			// Most updated is vector2
			newHistory = append(newHistory, histories.f2.hist...)
			log.Log(&f, "<MergeHistories> vector 2 is greater, so choosing that as history for planet \"%s\"", planet)

		}
	}
	log.Log(&f, "len(newHistory)=%d", len(newHistory))

	return newHistory
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
		//go SyncWithEventualConsistency()
	}

	forever := make(chan bool)
	<-forever
}
