package util

import (
	"bufio"
	"context"
	"dist/common/log"
	"dist/pb"
	"math/rand"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

func LoadEnvVariables(addrMap *map[string]string) error {
	err := godotenv.Load()
	if err != nil {
		return err
	}

	for env := range *addrMap {
		(*addrMap)[env] = os.Getenv(env)
	}

	return nil
}

// Read through each line of `filename`, and call the function
// `readLineCallback` using the current line as a parameter. If
// that function returns false, then this will stop reading
// lines.
// Returns true if file was created before attempting to read.
func ReadLines(filename string, readLineCallback func(string) bool) (bool, error) {
	fileExisted := log.FileExists(filename)

	file, err := os.OpenFile(filename, log.LstdAppendFlags, 0600)
	defer file.Close()
	if err != nil {
		return fileExisted, err
	}

	shouldStop := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		shouldStop = readLineCallback(scanner.Text())

		if shouldStop {
			break
		}
	}

	return fileExisted, nil
}

// Create file (or append) each line in arguments where '\n' is
// added at the end of each one.
// Returns true if file was created before attempting to open.
func WriteLines(filename string, lines ...string) (bool, error) {
	fileExisted := log.FileExists(filename)

	file, err := os.OpenFile(filename, log.LstdAppendFlags, 0600)
	defer file.Close()
	if err != nil {
		return fileExisted, err
	}

	finalString := ""
	for _, line := range lines {
		finalString += line + "\n"
	}

	file.WriteString(finalString)
	file.Sync()

	return fileExisted, nil
}

// Try to delete file if exists. If it doesn't, no errors will
// be raised and false will be returned. Otherwise, return true
func DeleteFile(filename string) bool {
	exists := log.FileExists(filename)

	if exists {
		err := os.Remove(filename)
		log.FailOnError(nil, err, "Couldn't remove file \"%s\"", filename)

		return true
	}

	return false
}

func ReadUserInput(msg string, a ...interface{}) (*pb.Command, string, string, interface{}, error) {
	// finalMsg := fmt.Sprintf(msg, a...)

	/* code */

	return pb.Command_ADD_CITY.Enum(), "planet", "city", 0, nil
}

func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}

// Same as 'strconv.Atoi(...)' but without having to handle the
// error yourself all the time
func StringToInt(value string) int {
	num, err := strconv.Atoi(value)
	log.FailOnError(nil, err, "Coudln't convert string \"%s\" to a number", value)

	return num
}

// Used when setting values in a protobuf object
func StringToUint32(value string) uint32 {
	return uint32(StringToInt(value))
}

// Return a random number within [min, max] both inclusive.
func RandInt(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}
