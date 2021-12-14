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
		FULCRUM:   []string{"dist181:10000", "dist182:10001", "dist183:10002"},
		INFORMANT: []string{"dist181:10010", "dist182:10011"},
		BROKER:    "dist184:10020",
		LEIA:      "dist183:10030",
	}
)

type TimeVector []uint32

func (t TimeVector) GreaterThanOrEqual(other TimeVector) bool {
	return t[0] >= other[0] && t[1] >= other[1] && t[2] >= other[2]
}

func (t TimeVector) ToProto() *pb.TimeVector {
	return &pb.TimeVector{Time: []uint32{t[0], t[1], t[2]}}
}

type GrpcClient struct {
	Client *pb.CommunicationClient
}

type CommandHistory struct {
	Command        pb.Command
	Planet         string
	City           string
	FulcrumAddress string
	TimeVector     TimeVector
}

type Planet = string
type City = string
