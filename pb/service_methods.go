package pb

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
	switch command {
	case "AddCity":
		return Command_ADD_CITY.Enum()

	case "UpdateName":
		return Command_UPDATE_NAME.Enum()

	case "UpdateNumber":
		return Command_UPDATE_NUMBER.Enum()

	case "DeleteCity":
		return Command_DELETE_CITY.Enum()

	default:
		return Command_GET_NUMBER_REBELS.Enum()

	}
}
