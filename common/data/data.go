package data

// Contains structs, interfaces, reciever functions, constants,
// etc.

var (
	Command = struct {
		ADD_CITY, DELETE_CITY, UPDATE_NAME, UPDATE_NUMBER, GET_NUMBER_REBELS uint8
	}{
		ADD_CITY:          0,
		DELETE_CITY:       1,
		UPDATE_NAME:       2,
		UPDATE_NUMBER:     3,
		GET_NUMBER_REBELS: 4,
	}
)
