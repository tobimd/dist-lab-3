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
func ReadLines(filename string, readLineCallback func(string) bool) error {
	file, err := os.OpenFile(filename, log.LstdAppendFlags, 0644)

	if err != nil {
		return err
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

	return nil
}

// Create file (or append) each line in arguments where '\n' is
// added at the end of each one.
// Returns true if file was created before attempting to open.
func WriteLines(filename string, overwrite bool, lines ...string) error {
	flags := log.LstdAppendFlags
	if overwrite {
		flags = log.LstdWriteFlags
	}
	file, err := os.OpenFile(filename, flags, 0644)

	if err != nil {
		return err
	}

	defer file.Close()

	finalString := ""
	for _, line := range lines {
		finalString += line + "\n"
	}

	file.WriteString(finalString)
	file.Sync()

	return nil
}

// Read through each line of `filename`, and call the function
// `replaceCallback` using the current line as a parameter, which
// should return the line replacement and a boolean value that
// tells this function to stop reading lines and just save.
//
// Returns true if file was created before attempting to read,
// false otherwise.
func ReplaceLines(filename string, replaceCallback func(string) string) error {
	result := make([]string, 1)

	ReadLines(filename, func(line string) bool {
		repl := replaceCallback(line)
		result = append(result, repl)
		return false
	})

	for i, l := range result {
		if l == "" {
			result = append(result[:i], result[i+1:]...)
		}
	}

	err := WriteLines(filename, true, result...)

	return err
}

// Read through each line of `filename`, and call the function
// `deleteCallback` using the current line as a parameter, which
// should return the a boolean value that tells this function to
// delete de current line and another to stop reading lines and
// just save.
//
// Returns true if file was created before attempting to read,
// false otherwise.
func DeleteLines(filename string, deleteCallback func(string) bool) error {
	file, err := os.OpenFile(filename, log.LstdAppendFlags, 0644)

	if err != nil {
		return err
	}

	result := make([]string, 1)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var include bool
		text := scanner.Text()
		log.Log(nil, "<DeleteLines> current line is \"%s\"", text)
		include = deleteCallback(text)
		log.Log(nil, "<DeleteLines> line included ? %v", include)

		if include {
			result = append(result, text)
		}
	}
	file.Close()

	err = WriteLines(filename, true, result...)

	return err
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
func ReadUserInput(f *string, msg string, a ...interface{}) (*pb.Command, string, string, string) {
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

	if *cmd == pb.Command_UPDATE_NAME || *cmd == pb.Command_UPDATE_NUMBER || *cmd == pb.Command_ADD_CITY {
		// If command is UpdateName, UpdateNumber or AddCity, then
		// return the string of the 4th argument
		return cmd, planet, city, strs[3]

	} else {
		// Otherwise, the command does not use the 4th argument
		return cmd, planet, city, ""
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
