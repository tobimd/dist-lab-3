package test

import (
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
		// Check if fulcrum servers save info correctly
		CreateEntities(Entity{"fulcrum", 0}, Entity{"fulcrum", 1}, Entity{"fulcrum", 2})

	}

}
