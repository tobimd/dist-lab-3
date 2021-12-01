package util

import (
	"bufio"
	"context"
	"dist/common/data"
	"dist/common/log"
	"os"
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

func ReadUserInput(msg string, a ...interface{}) (data.CommandEnum, string, string, interface{}, error) {
	// finalMsg := fmt.Sprintf(msg, a...)

	// ...

	return data.CommandEnum{}, "planet", "city", 0, nil
}

func GetContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), time.Second)
}
