package main

import (
	"dist/common/log"
	"dist/entities/broker"
	"dist/entities/fulcrum"
	"dist/entities/informant"
	"dist/entities/leia"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var (
	f = "main.log"
)

func showHelp() {
	pName := os.Args[0]
	log.Fatal(&f, "\nUsage: %s run <broker | fulcrum | informant | leia> [entity_id]\n  or:  %s test <test_number>", pName, pName)
}

func main() {
	argv := os.Args[:]

	if len(argv) < 3 {
		showHelp()
	}

	cmd := argv[1]
	arg := argv[2]
	id := -1

	if len(argv) > 3 {
		var err error
		id, err = strconv.Atoi(argv[3])
		log.FailOnError(&f, err, "Couldn't convert \"%s\" to a number", argv[3])
	}

	rand.Seed(time.Now().Unix())

	if cmd == "run" {
		switch arg {
		case "broker":
			broker.Run()

		case "fulcrum":
			if id < 0 || id > 2 {
				log.Fatal(&f, "Please make sure that a valid id is given for fulcrum entity (0, 1 or 2)")
			}
			fulcrum.Run(id)

		case "informant":
			if id < 0 || id > 1 {
				log.Fatal(&f, "Please make sure that a valid id is given for informant entity (0 or 1)")
			}
			informant.Run(id)

		case "leia":
			leia.Run()

		default:
			log.Print(&f, "Couldn't recognize entity: \"%s\"", arg)
			showHelp()
		}

	} else if cmd == "test" {
		testNumber, err := strconv.Atoi(arg)
		log.FailOnError(&f, err, "Couldn't convert \"%s\" to a number", arg)

		if testNumber < 0 {
			log.Fatal(&f, "Please make sure that a valid test number is given (0, 1, ...)")
		}

	} else {
		showHelp()
	}
}
