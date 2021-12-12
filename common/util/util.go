package util

import (
	"bufio"
	"context"
	"dist/common/log"
	"dist/pb"
	"os"
	"strconv"
	"strings"
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

	if err != nil {
		return fileExisted, err
	}

	defer file.Close()

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

	if err != nil {
		return fileExisted, err
	}

	defer file.Close()

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

// Read user input with a given message passed to user
func ReadUserInput(f *string, msg string, a ...interface{}) (*pb.Command, string, string, interface{}) {
	// Print formatted message to screen and read input
	log.Print(f, msg, a...)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	log.FailOnError(f, err, "Problem at ReadUserInput while reading line")

	// Clean input and split arguments
	line = strings.TrimSuffix(line, "\n")
	strs := strings.Split(line, " ")

	cmd := pb.CommandFromString(strs[0])
	planet := strs[1]
	city := strs[2]

	// If command is UpdateName, then the 4th argument is
	// the city's new name
	if *cmd == pb.Command_UPDATE_NAME {
		return cmd, planet, city, strs[3]

		// If command is NOT DeleteCity, then 4th argument is
		// number of rebels and could be ommited, defaulting to
		// a value 0
	} else if *cmd != pb.Command_DELETE_CITY {
		numRebels := 0
		if len(strs) == 4 {
			numRebels = StringToInt(strs[3])
		}
		return cmd, planet, city, numRebels

		// Otherwise, the command IS DeleteCity and 4th argument
		// doesn't matter
	} else {
		return cmd, planet, city, nil
	}
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
