package test

import (
	"dist/common/log"
	"dist/entities/fulcrum"
	"fmt"
	"os"
	"os/exec"
)

var (
	planets = []string{
		"Abafar",
		"Aleen",
		"Bogano",
		"Concord_Dawn",
		"Devaron",
		"Exegol",
		"Kuat",
		"Ord_Mantell",
		"Savareen",
		"Tatooine",
		"Vardos",
		"Zeffo",
	}

	cities = []string{
		"Adascopolis",
		"Aldera",
		"Dac_City",
		"Dearic",
		"Falleen_Throne",
		"Nuba_City",
		"Montrol_City",
	}

	f = "test.log"
)

type Entity struct {
	name string
	id   int
}

func CreateEntity(entity Entity) {
	cmd := fmt.Sprintf("%s run %s %d &", os.Args[0], entity.name, entity.id)
	exec.Command("/bin/bash", "-c", cmd)
}

func CreateEntities(entities ...Entity) {
	for _, entity := range entities {
		CreateEntity(entity)
	}
}

func Run(testNumber int) {
	switch testNumber {

	case 0:
		// Test if a fulcrum server runs utility functions
		// correctly (SavePlanetData, etc)
		//
		// Note, will fail if previous <planet>.txt logs aren't
		// removed
		log.Print(&f, "planet:%s -> time vector:%v", planets[0], fulcrum.SavePlanetData(planets[0], cities[0], 1))
		log.Print(&f, "planet:%s -> time vector:%v", planets[0], fulcrum.SavePlanetData(planets[0], cities[1], 5))
		log.Print(&f, "planet:%s -> time vector:%v", planets[2], fulcrum.SavePlanetData(planets[1], cities[4], 9))
		log.Print(&f, "planet:%s -> time vector:%v", planets[0], fulcrum.SavePlanetData(planets[0], cities[0], 2))

		// Expected output:
		//     planet:Abafar -> time vector:{[1 0 0]}
		//     planet:Abafar -> time vector:{[2 0 0]}
		//     planet:Bogano -> time vector:{[1 0 0]}
		//     planet:Abafar -> time vector:{[3 0 0]}

	}

}
