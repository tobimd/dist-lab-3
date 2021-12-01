package data

// Contains structs, interfaces, reciever functions, constants,
// etc.

var (
	Command = struct {
		ADD_CITY, DELETE_CITY, UPDATE_NAME, UPDATE_NUMBER, GET_NUMBER_REBELS CommandEnum
	}{
		ADD_CITY:          CommandEnum{"AddCity", 0},
		DELETE_CITY:       CommandEnum{"DeleteCity", 1},
		UPDATE_NAME:       CommandEnum{"UpdateName", 2},
		UPDATE_NUMBER:     CommandEnum{"UpdateNumber", 3},
		GET_NUMBER_REBELS: CommandEnum{"GetNumRebels", 4},
	}

	Entity = struct {
		FULCRUM_0, FULCRUM_1, FULCRUM_2, INFORMANT_0, INFORMANT_1, BROKER, LEIA string
	}{
		FULCRUM_0:   "FULCRUM_0",
		FULCRUM_1:   "FULCRUM_1",
		FULCRUM_2:   "FULCRUM_2",
		INFORMANT_0: "INFORMANT_0",
		INFORMANT_1: "INFORMANT_0",
		BROKER:      "BROKER",
		LEIA:        "LEIA",
	}
)

type CommandEnum struct {
	Name string
	Enum uint8
}
