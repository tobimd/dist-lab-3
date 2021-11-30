package fulcrum

import (
	"dist/common/data"
	"dist/common/log"
	"dist/common/util"
	"fmt"
)

var (
	id int
	f  = ""
)

// Open (or create) a registry file corresponding to a specific
// planet and add the corresponding information onto it.
//
// For planet "Tatooine", city "Mos Eisley" and number of rebels
// "5", all info is stored in the following way:
//
//     -- Tatooine.txt
//     Tatooine Mos_Eisley 5
//
func SavePlanetData(planet string, city string, numRebels int) {
	filename := planet + ".txt"
	data := fmt.Sprintf("%s %s %d", planet, city, numRebels)

	err := util.WriteLines(filename, data)
	log.FailOnError(&f, err, "Couldn't write to file \"%s\"", filename)
}

func LoadPlanetData(planet string) {

}

func UpdatePlanetLog(command data.CommandEnum, planet string, city string, value interface{}) {
	filename := fmt.Sprintf("log.%s.txt", planet)
	data := ""

	if value != nil {
		data = fmt.Sprintf("%s %s %s %v", command.Name, planet, city, value)
	} else {
		data = fmt.Sprintf("%s %s %s", command.Name, planet, city)
	}

	err := util.WriteLines(filename, data)
	log.FailOnError(&f, err, "Couldn't write to file \"%s\"", filename)
}

func Run(fulcrumId int) {
	id = fulcrumId
	f = fmt.Sprintf("fulcrum_%d", id)

}
