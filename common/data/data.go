package data

// Contains structs, interfaces, reciever functions, constants,
// etc.

import (
	"dist/common/log"
	"dist/pb"
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

	Address = struct {
		FULCRUM, INFORMANT []string
		BROKER, LEIA       string
	}{
		FULCRUM:   []string{"0.0.0.0:10000", "0.0.0.0:10001", "0.0.0.0:10002"},
		INFORMANT: []string{"0.0.0.0:10010", "0.0.0.0:10011"},
		BROKER:    "0.0.0.0:10020",
		LEIA:      "0.0.0.0:10030",
	}

	Channel = struct {
		STOP uint8
	}{
		STOP: 0,
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

type GrpcData struct {
}
