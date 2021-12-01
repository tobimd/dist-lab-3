package data

// Contains structs, interfaces, reciever functions, constants,
// etc.

import (
	"dist/common/log"
	pb "dist/proto"
)

var (
	Command = struct {
		ADD_CITY, DELETE_CITY, UPDATE_NAME, UPDATE_NUMBER, GET_NUMBER_REBELS, CHECK_CONSISTENCY CommandEnum
	}{
		ADD_CITY:          CommandEnum{"AddCity", 0},
		DELETE_CITY:       CommandEnum{"DeleteCity", 1},
		UPDATE_NAME:       CommandEnum{"UpdateName", 2},
		UPDATE_NUMBER:     CommandEnum{"UpdateNumber", 3},
		GET_NUMBER_REBELS: CommandEnum{"GetNumRebels", 4},
		CHECK_CONSISTENCY: CommandEnum{"", 5},
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

func ToCommandEnum(pbCommand *pb.Command) CommandEnum {
	switch *pbCommand {
	case pb.Command_ADD_CITY:
		return Command.ADD_CITY

	case pb.Command_DELETE_CITY:
		return Command.DELETE_CITY

	case pb.Command_UPDATE_NAME:
		return Command.UPDATE_NAME

	case pb.Command_UPDATE_NUMBER:
		return Command.UPDATE_NUMBER

	case pb.Command_GET_NUMBER_REBELS:
		return Command.GET_NUMBER_REBELS

	case pb.Command_CHECK_CONSISTENCY:
		return Command.CHECK_CONSISTENCY
	}

	log.Fatal(nil, "Couldn't recognize command \"%s\"", pbCommand.String())
	return CommandEnum{}
}

type CommandEnum struct {
	Name string
	Enum uint8
}

type UserInput struct {
	Command CommandEnum
	planet  string
	city    string
	value   interface{}
}
