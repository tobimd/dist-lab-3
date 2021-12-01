package util

import (
	"bufio"
	"dist/common/log"
	"errors"
	"os"

	"github.com/joho/godotenv"
)

// [Internal] Return true if file exists, false otherwise
func fileExists(filename string) bool {
	_, err := os.Stat(filename)
	return !errors.Is(err, os.ErrNotExist)
}

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
	fileExisted := fileExists(filename)

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
	fileExisted := fileExists(filename)

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
