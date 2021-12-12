package data

// Contains structs, interfaces, reciever functions, constants,
// etc.

import (
	"dist/pb"
)

var (
	Address = struct {
		FULCRUM, INFORMANT []string
		BROKER, LEIA       string
	}{
		FULCRUM:   []string{"0.0.0.0:10000", "0.0.0.0:10001", "0.0.0.0:10002"},
		INFORMANT: []string{"0.0.0.0:10010", "0.0.0.0:10011"},
		BROKER:    "0.0.0.0:10020",
		LEIA:      "0.0.0.0:10030",
	}
)

type TimeVector struct {
	Time []uint32
}

type GrpcClient struct {
	Client *pb.CommunicationClient
}

type CommandHistory struct {
	command        pb.Command
	city           string
	fulcrumAddress string
	timeVector     TimeVector
}

type Planet = string
type City = string
