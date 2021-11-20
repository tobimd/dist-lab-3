package util

import (
	"dist/common/log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(f *string, addrMap *map[string]string) {
	err := godotenv.Load()
	log.FailOnError(f, err, "Couldn't load variables from \".env\" at the root of the project")

	for env := range *addrMap {
		(*addrMap)[env] = os.Getenv(env)
		log.Log(f, "Loaded environment variable \"%s\" = \"%s\"", env, (*addrMap)[env])
	}
}
