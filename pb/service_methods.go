package pb

import "strings"

func (cmd *Command) ToString() string {
	switch *cmd {
	case Command_ADD_CITY:
		return "AddCity"

	case Command_UPDATE_NAME:
		return "UpdateName"

	case Command_UPDATE_NUMBER:
		return "UpdateNumber"

	case Command_DELETE_CITY:
		return "DeleteCity"

	case Command_GET_NUMBER_REBELS:
		return "GetNumberRebels"

	default:
		return ""
	}
}

func CommandFromString(command string) *Command {

	command = strings.ToLower(command)

	switch command {
	case "addcity":
		return Command_ADD_CITY.Enum()

	case "updatename":
		return Command_UPDATE_NAME.Enum()

	case "updatenumber":
		return Command_UPDATE_NUMBER.Enum()

	case "deletecity":
		return Command_DELETE_CITY.Enum()

	case "getnumberrebels":
		return Command_GET_NUMBER_REBELS.Enum()
	default:
		return nil

	}
}
