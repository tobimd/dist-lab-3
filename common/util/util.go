package util

import (
	"bufio"
	"dist/common/log"
	"os"

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
func ReadLines(filename string, readLineCallback func(string) bool) error {
	file, err := os.OpenFile(filename, log.LstdAppendFlags, 0600)
	defer file.Close()
	if err != nil {
		return err
	}

	shouldStop := false
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		shouldStop = readLineCallback(scanner.Text())

		if shouldStop {
			break
		}
	}

	return nil
}

func WriteLines(filename string, lines ...string) error {
	file, err := os.OpenFile(filename, log.LstdAppendFlags, 0600)
	defer file.Close()
	if err != nil {
		return err
	}

	finalString := ""
	for _, line := range lines {
		finalString += line + "\n"
	}

	file.WriteString(finalString)
	file.Sync()

	return nil
}
