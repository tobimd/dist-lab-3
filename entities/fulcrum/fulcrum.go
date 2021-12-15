package fulcrum

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"dist/pb"
	"fmt"
	"os"
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

	// A list of all filenames for log.<planet>.txt used when
	// merging so that deleting them is easier
	planetLogList = make([]string, 0)

	// All grpc clients stored with each entity's address used
	// as keys for the map object. `data.Client` is a pointer
	// to pb.CommuncationClient which should be what gets
	// returned by `util.SetupClient(...)`
	clients = make(map[string]*Client)

	// When using 'SavePlanetData', determine wether to append
	// to the planet's history file, or rewrite it completely
	// which is used when merging.
	StoreMethod = struct{ Append, Update, Delete, Rewrite Method }{Append: 0, Update: 1, Delete: 2, Rewrite: 3}

	// When leader fulcrum (id=0) tells neighbors to send their
	// history to it, the server function BroadcastChanges will
	// be called on leader and the response should be the merge
	// from all three histories. To solve that issue, this
	// channel will recieve boolean true when both histories
	// have been stored
	neighborChangesCh = make(chan bool)
	neighborHistories = make(map[string][]*pb.CommandParams, 2)
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

	var err error

	switch storeMethod {
	case StoreMethod.Append, StoreMethod.Rewrite:
		overwrite := false
		if storeMethod == StoreMethod.Rewrite {
			overwrite = true
		}

		// Write data into file
		err = util.WriteLines(filename, overwrite, info)
		log.Log(&f, "<SavePlanetData> storeMethod is Create (%v) or Rewrite (%v), written to file", overwrite, !overwrite)

	case StoreMethod.Update:
		replacedLine := ""
		err = util.ReplaceLines(filename, func(line string) string {
			values := strings.Split(line, " ")
			log.Log(&f, "values: %v", values)
			if len(values) == 0 {
				return line
			}
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
			if len(values) == 0 {
				return true
			}
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

	// Add filename to planetLogList if it wasn't added before
	exists := false
	for _, planetLogName := range planetLogList {
		if planetLogName == planet {
			exists = true
		}
	}

	if !exists {
		planetLogList = append(planetLogList, planet)
	}
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

// Checks if city exist on planet registry
func CityDoesExistOn(planet string, city string) (exist bool) {
	filename := fmt.Sprintf("%s.txt", planet)
	exist = false
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return
	}
	util.ReadLines(filename, func(line string) bool {
		values := strings.Split(line, " ")
		if values[1] == city {
			exist = true
			return true
		}
		return false
	})
	return
}

func ReadPlanetLog(planet string) []*pb.CommandParams {
	filename := fmt.Sprintf("log.%s.txt", planet)
	info := make([]*pb.CommandParams, 0)

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
				LastTimeVector: nil,
			})

			return false
		},
	)

	log.Log(&f, "Info recieved: %+v", info)
	return info
}

func MergeHistories() []*pb.CommandParams {
	var newHistory []*pb.CommandParams
	allHistories := [3][]*pb.CommandParams{
		GetHistory(),
		neighborHistories[data.Address.FULCRUM[1]],
		neighborHistories[data.Address.FULCRUM[2]],
	}

	type FH struct {
		vector []data.TimeVector
		hist   [][]*pb.CommandParams
	}

	info := make(map[string]FH, 0)

	for i := 0; i < 3; i++ {
		hist := allHistories[i]

		currHist := make([]*pb.CommandParams, 0)

		for _, cmdParams := range hist {
			planet := cmdParams.GetPlanetName()
			vector := cmdParams.GetLastTimeVector()

			currHist = append(currHist, cmdParams)

			if vector != nil {
				if _, ok := info[planet]; !ok {
					fh := FH{
						hist:   [][]*pb.CommandParams{{}, {}, {}},
						vector: []data.TimeVector{{0, 0, 0}, {0, 0, 0}, {0, 0, 0}},
					}
					info[planet] = fh
				}

				v := vector.GetTime()
				info[planet].vector[i] = data.TimeVector{v[0], v[1], v[2]}
				info[planet].hist[i] = currHist

				currHist = []*pb.CommandParams{}
			}
		}
	}

	log.Log(&f, "amount of planets from all fulcrums is %d", len(info))
	for planet, history := range info {
		for i := 0; i < 3; i++ {
			hOwn := history.hist[i]

			vOwn := history.vector[i]
			vTheirs1 := history.vector[(i+1)%3]
			vTheirs2 := history.vector[(i+2)%3]

			// decides which history to take as the doodoo one
			if vOwn.GreaterThanOrEqual(vTheirs1) && vOwn.GreaterThanOrEqual(vTheirs2) {
				log.Log(&f, "<MergeHistories> vector from fulcrum %d is greater, so we'll choose that as history for planet \"%s\"", i, planet)

				newHistory = append(newHistory, hOwn...)
				break
			}
		}

		// merge all vectors into one
		finalVector := &pb.TimeVector{
			Time: []uint32{
				history.vector[0][0],
				history.vector[1][1],
				history.vector[2][2],
			},
		}

		if len(newHistory) != 0 {
			newHistory[len(newHistory)-1].LastTimeVector = finalVector
		}
	}
	log.Log(&f, "len(newHistory)=%d", len(newHistory))

	return newHistory
}

func SetHistory(newHistory []*pb.CommandParams) {
	currPlanet := ""
	planetVectors = map[string]data.TimeVector{}

	for _, history := range newHistory {
		planet := history.GetPlanetName()
		city := history.GetCityName()
		numRebels := int(history.GetNumOfRebels())

		if history.NumOfRebels != nil {
			numRebels = int(history.GetNewNumOfRebels())
		}

		if history.NewCityName != nil {
			city = history.GetCityName()
		}

		method := StoreMethod.Append
		if planet != currPlanet {
			method = StoreMethod.Rewrite
			currPlanet = planet
		}

		SavePlanetData(planet, city, numRebels, "", method)

		// overwrite last vector for a given planet to the final vector agreement
		if history.LastTimeVector != nil {
			planetVectors[planet] = history.LastTimeVector.GetTime()
		}

		currPlanet = planet
	}

	// Delete all log.<planet>.txt files
	for _, planet := range planetLogList {
		util.DeleteFile("log." + planet + ".txt")
	}
	util.DeleteFile(".txt")

	planetLogList = []string{}
}

func GetHistory() []*pb.CommandParams {
	history := make([]*pb.CommandParams, 0)
	for planet := range planetVectors {
		history = append(history, ReadPlanetLog(planet)...)

		if len(history) == 0 {
			continue
		}

		last := history[len(history)-1]
		log.Log(&f, "history variable: %+v", last)

		planet := last.GetPlanetName()
		city := last.GetCityName()
		num := last.GetNumOfRebels()
		newcity := last.GetNewCityName()
		newnum := last.GetNewNumOfRebels()
		command := last.GetCommand()

		history[len(history)-1] = &pb.CommandParams{
			PlanetName:     &planet,
			CityName:       &city,
			NumOfRebels:    &num,
			NewCityName:    &newcity,
			NewNumOfRebels: &newnum,
			Command:        &command,
			LastTimeVector: planetVectors[planet].ToProto(),
		}
	}
	return history
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
			neighbor := data.Address.FULCRUM[i]
			log.Log(&f, "About to call BoradcastChanges to fulcrum %d returns it's history", i)

			neighborHistories[neighbor] = clients[neighbor].BroadcastChanges([]*pb.CommandParams{})
		}

		newHistory := MergeHistories()

		SetHistory(newHistory)
		for i := 1; i < 3; i++ {
			neighbor := data.Address.FULCRUM[i]
			log.Log(&f, "About to call BoradcastChanges to fulcrum %d with new history", i)

			clients[neighbor].BroadcastChanges(newHistory)
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
		go SyncWithEventualConsistency()
	}

	forever := make(chan bool)
	<-forever
}
