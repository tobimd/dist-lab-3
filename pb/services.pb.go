// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.0
// source: pb/services.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type Command int32

const (
	Command_ADD_CITY          Command = 0 // [Informants] used to add a new city to a planet
	Command_DELETE_CITY       Command = 1 // [Informants] used to delete city from a planet
	Command_UPDATE_NAME       Command = 2 // [Informants] used to change name of a city
	Command_UPDATE_NUMBER     Command = 3 // [Informants] used to change number of rebels in a city
	Command_GET_NUMBER_REBELS Command = 4 // [Leia] return the number of rebels in a given city and planet
	Command_CHECK_CONSISTENCY Command = 5 // [Fulcrum] used internally to tell other fulcrums to send their history logs
	Command_GET_PLANET_TIME   Command = 6 // [Broker] used internally to get last TimeVector for a given planet on a fulcrum server
)

// Enum value maps for Command.
var (
	Command_name = map[int32]string{
		0: "ADD_CITY",
		1: "DELETE_CITY",
		2: "UPDATE_NAME",
		3: "UPDATE_NUMBER",
		4: "GET_NUMBER_REBELS",
		5: "CHECK_CONSISTENCY",
		6: "GET_PLANET_TIME",
	}
	Command_value = map[string]int32{
		"ADD_CITY":          0,
		"DELETE_CITY":       1,
		"UPDATE_NAME":       2,
		"UPDATE_NUMBER":     3,
		"GET_NUMBER_REBELS": 4,
		"CHECK_CONSISTENCY": 5,
		"GET_PLANET_TIME":   6,
	}
)

func (x Command) Enum() *Command {
	p := new(Command)
	*p = x
	return p
}

func (x Command) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Command) Descriptor() protoreflect.EnumDescriptor {
	return file_pb_services_proto_enumTypes[0].Descriptor()
}

func (Command) Type() protoreflect.EnumType {
	return &file_pb_services_proto_enumTypes[0]
}

func (x Command) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Command.Descriptor instead.
func (Command) EnumDescriptor() ([]byte, []int) {
	return file_pb_services_proto_rawDescGZIP(), []int{0}
}

type TimeVector struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Max amount of 'time' objects should be 3 because of each
	// Fulcrum, even if more is allowed due to being a 'repeated'
	// type of object.
	Time []uint32 `protobuf:"varint,1,rep,packed,name=time,proto3" json:"time,omitempty"`
}

func (x *TimeVector) Reset() {
	*x = TimeVector{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_services_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TimeVector) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TimeVector) ProtoMessage() {}

func (x *TimeVector) ProtoReflect() protoreflect.Message {
	mi := &file_pb_services_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TimeVector.ProtoReflect.Descriptor instead.
func (*TimeVector) Descriptor() ([]byte, []int) {
	return file_pb_services_proto_rawDescGZIP(), []int{0}
}

func (x *TimeVector) GetTime() []uint32 {
	if x != nil {
		return x.Time
	}
	return nil
}

type CommandParams struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Command        *Command `protobuf:"varint,1,opt,name=command,proto3,enum=proto.Command,oneof" json:"command,omitempty"`                      // All possible commands that could be used (even if some entities do and some don't)
	PlanetName     *string  `protobuf:"bytes,2,opt,name=planet_name,json=planetName,proto3,oneof" json:"planet_name,omitempty"`                  // The current planet being affected by the command
	CityName       *string  `protobuf:"bytes,3,opt,name=city_name,json=cityName,proto3,oneof" json:"city_name,omitempty"`                        // The current city from that planet being affected, or the name of the city if "AddCity" is used
	NumOfRebels    *uint32  `protobuf:"varint,4,opt,name=num_of_rebels,json=numOfRebels,proto3,oneof" json:"num_of_rebels,omitempty"`            // The current rebels in the city. NOTE: Only used when BroadcastChanges is being returned with the planet's history.
	NewCityName    *string  `protobuf:"bytes,5,opt,name=new_city_name,json=newCityName,proto3,oneof" json:"new_city_name,omitempty"`             // If command is "UpdateName", then this will contain the new city's name
	NewNumOfRebels *uint32  `protobuf:"varint,6,opt,name=new_num_of_rebels,json=newNumOfRebels,proto3,oneof" json:"new_num_of_rebels,omitempty"` // If command is "UpdateNumber" or "AddCity" (with default value 0), then this will be set
	// Used by fulcrums when syncing, append the time vector for
	// merging, and may only be set in the last CommandParams
	// added for a given planet (because time vectors only
	// contain) info about the last update made to the planet.
	// Also used by Informant to comply with read your writes
	LastTimeVector *TimeVector `protobuf:"bytes,7,opt,name=last_time_vector,json=lastTimeVector,proto3,oneof" json:"last_time_vector,omitempty"`
}

func (x *CommandParams) Reset() {
	*x = CommandParams{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_services_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommandParams) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandParams) ProtoMessage() {}

func (x *CommandParams) ProtoReflect() protoreflect.Message {
	mi := &file_pb_services_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandParams.ProtoReflect.Descriptor instead.
func (*CommandParams) Descriptor() ([]byte, []int) {
	return file_pb_services_proto_rawDescGZIP(), []int{1}
}

func (x *CommandParams) GetCommand() Command {
	if x != nil && x.Command != nil {
		return *x.Command
	}
	return Command_ADD_CITY
}

func (x *CommandParams) GetPlanetName() string {
	if x != nil && x.PlanetName != nil {
		return *x.PlanetName
	}
	return ""
}

func (x *CommandParams) GetCityName() string {
	if x != nil && x.CityName != nil {
		return *x.CityName
	}
	return ""
}

func (x *CommandParams) GetNumOfRebels() uint32 {
	if x != nil && x.NumOfRebels != nil {
		return *x.NumOfRebels
	}
	return 0
}

func (x *CommandParams) GetNewCityName() string {
	if x != nil && x.NewCityName != nil {
		return *x.NewCityName
	}
	return ""
}

func (x *CommandParams) GetNewNumOfRebels() uint32 {
	if x != nil && x.NewNumOfRebels != nil {
		return *x.NewNumOfRebels
	}
	return 0
}

func (x *CommandParams) GetLastTimeVector() *TimeVector {
	if x != nil {
		return x.LastTimeVector
	}
	return nil
}

// Corresponds to the history of commands being run over a
// planet. That includes info of arguments such as the new number
// of rebels being assigned to it.
type FulcrumHistory struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	History []*CommandParams `protobuf:"bytes,1,rep,name=history,proto3" json:"history,omitempty"`
}

func (x *FulcrumHistory) Reset() {
	*x = FulcrumHistory{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_services_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FulcrumHistory) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FulcrumHistory) ProtoMessage() {}

func (x *FulcrumHistory) ProtoReflect() protoreflect.Message {
	mi := &file_pb_services_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FulcrumHistory.ProtoReflect.Descriptor instead.
func (*FulcrumHistory) Descriptor() ([]byte, []int) {
	return file_pb_services_proto_rawDescGZIP(), []int{2}
}

func (x *FulcrumHistory) GetHistory() []*CommandParams {
	if x != nil {
		return x.History
	}
	return nil
}

// Simple message that contains the response from a fulcrum server
// (or broker in the first interaction from clients)
type FulcrumResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	FulcrumRedirectAddr *string       `protobuf:"bytes,1,opt,name=fulcrum_redirect_addr,json=fulcrumRedirectAddr,proto3,oneof" json:"fulcrum_redirect_addr,omitempty"` // fulcrum server address for informants
	NumOfRebels         *uint32       `protobuf:"varint,2,opt,name=num_of_rebels,json=numOfRebels,proto3,oneof" json:"num_of_rebels,omitempty"`                        // number of rebels for Leia
	SingleTimeVector    *TimeVector   `protobuf:"bytes,3,opt,name=single_time_vector,json=singleTimeVector,proto3,oneof" json:"single_time_vector,omitempty"`          // used by informants in a read-your-writes operation
	MultTimeVectors     []*TimeVector `protobuf:"bytes,4,rep,name=mult_time_vectors,json=multTimeVectors,proto3" json:"mult_time_vectors,omitempty"`                   // used by broker to get all timevectors for a single planet
}

func (x *FulcrumResponse) Reset() {
	*x = FulcrumResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_services_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *FulcrumResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FulcrumResponse) ProtoMessage() {}

func (x *FulcrumResponse) ProtoReflect() protoreflect.Message {
	mi := &file_pb_services_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FulcrumResponse.ProtoReflect.Descriptor instead.
func (*FulcrumResponse) Descriptor() ([]byte, []int) {
	return file_pb_services_proto_rawDescGZIP(), []int{3}
}

func (x *FulcrumResponse) GetFulcrumRedirectAddr() string {
	if x != nil && x.FulcrumRedirectAddr != nil {
		return *x.FulcrumRedirectAddr
	}
	return ""
}

func (x *FulcrumResponse) GetNumOfRebels() uint32 {
	if x != nil && x.NumOfRebels != nil {
		return *x.NumOfRebels
	}
	return 0
}

func (x *FulcrumResponse) GetSingleTimeVector() *TimeVector {
	if x != nil {
		return x.SingleTimeVector
	}
	return nil
}

func (x *FulcrumResponse) GetMultTimeVectors() []*TimeVector {
	if x != nil {
		return x.MultTimeVectors
	}
	return nil
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pb_services_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_pb_services_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_pb_services_proto_rawDescGZIP(), []int{4}
}

var File_pb_services_proto protoreflect.FileDescriptor

var file_pb_services_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x62, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x73, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x20, 0x0a, 0x0a, 0x54, 0x69,
	0x6d, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0d, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x22, 0xc3, 0x03, 0x0a,
	0x0d, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x12, 0x2d,
	0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32,
	0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x48,
	0x00, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x88, 0x01, 0x01, 0x12, 0x24, 0x0a,
	0x0b, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x48, 0x01, 0x52, 0x0a, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74, 0x4e, 0x61, 0x6d, 0x65,
	0x88, 0x01, 0x01, 0x12, 0x20, 0x0a, 0x09, 0x63, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x48, 0x02, 0x52, 0x08, 0x63, 0x69, 0x74, 0x79, 0x4e, 0x61,
	0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x27, 0x0a, 0x0d, 0x6e, 0x75, 0x6d, 0x5f, 0x6f, 0x66, 0x5f,
	0x72, 0x65, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x03, 0x52, 0x0b,
	0x6e, 0x75, 0x6d, 0x4f, 0x66, 0x52, 0x65, 0x62, 0x65, 0x6c, 0x73, 0x88, 0x01, 0x01, 0x12, 0x27,
	0x0a, 0x0d, 0x6e, 0x65, 0x77, 0x5f, 0x63, 0x69, 0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x18,
	0x05, 0x20, 0x01, 0x28, 0x09, 0x48, 0x04, 0x52, 0x0b, 0x6e, 0x65, 0x77, 0x43, 0x69, 0x74, 0x79,
	0x4e, 0x61, 0x6d, 0x65, 0x88, 0x01, 0x01, 0x12, 0x2e, 0x0a, 0x11, 0x6e, 0x65, 0x77, 0x5f, 0x6e,
	0x75, 0x6d, 0x5f, 0x6f, 0x66, 0x5f, 0x72, 0x65, 0x62, 0x65, 0x6c, 0x73, 0x18, 0x06, 0x20, 0x01,
	0x28, 0x0d, 0x48, 0x05, 0x52, 0x0e, 0x6e, 0x65, 0x77, 0x4e, 0x75, 0x6d, 0x4f, 0x66, 0x52, 0x65,
	0x62, 0x65, 0x6c, 0x73, 0x88, 0x01, 0x01, 0x12, 0x40, 0x0a, 0x10, 0x6c, 0x61, 0x73, 0x74, 0x5f,
	0x74, 0x69, 0x6d, 0x65, 0x5f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x18, 0x07, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x56, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x48, 0x06, 0x52, 0x0e, 0x6c, 0x61, 0x73, 0x74, 0x54, 0x69, 0x6d, 0x65,
	0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x88, 0x01, 0x01, 0x42, 0x0a, 0x0a, 0x08, 0x5f, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x42, 0x0e, 0x0a, 0x0c, 0x5f, 0x70, 0x6c, 0x61, 0x6e, 0x65, 0x74,
	0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x0c, 0x0a, 0x0a, 0x5f, 0x63, 0x69, 0x74, 0x79, 0x5f, 0x6e,
	0x61, 0x6d, 0x65, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x6e, 0x75, 0x6d, 0x5f, 0x6f, 0x66, 0x5f, 0x72,
	0x65, 0x62, 0x65, 0x6c, 0x73, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x6e, 0x65, 0x77, 0x5f, 0x63, 0x69,
	0x74, 0x79, 0x5f, 0x6e, 0x61, 0x6d, 0x65, 0x42, 0x14, 0x0a, 0x12, 0x5f, 0x6e, 0x65, 0x77, 0x5f,
	0x6e, 0x75, 0x6d, 0x5f, 0x6f, 0x66, 0x5f, 0x72, 0x65, 0x62, 0x65, 0x6c, 0x73, 0x42, 0x13, 0x0a,
	0x11, 0x5f, 0x6c, 0x61, 0x73, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x22, 0x40, 0x0a, 0x0e, 0x46, 0x75, 0x6c, 0x63, 0x72, 0x75, 0x6d, 0x48, 0x69, 0x73,
	0x74, 0x6f, 0x72, 0x79, 0x12, 0x2e, 0x0a, 0x07, 0x68, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x18,
	0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x50, 0x61, 0x72, 0x61, 0x6d, 0x73, 0x52, 0x07, 0x68, 0x69, 0x73,
	0x74, 0x6f, 0x72, 0x79, 0x22, 0xbb, 0x02, 0x0a, 0x0f, 0x46, 0x75, 0x6c, 0x63, 0x72, 0x75, 0x6d,
	0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x15, 0x66, 0x75, 0x6c, 0x63,
	0x72, 0x75, 0x6d, 0x5f, 0x72, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x5f, 0x61, 0x64, 0x64,
	0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x48, 0x00, 0x52, 0x13, 0x66, 0x75, 0x6c, 0x63, 0x72,
	0x75, 0x6d, 0x52, 0x65, 0x64, 0x69, 0x72, 0x65, 0x63, 0x74, 0x41, 0x64, 0x64, 0x72, 0x88, 0x01,
	0x01, 0x12, 0x27, 0x0a, 0x0d, 0x6e, 0x75, 0x6d, 0x5f, 0x6f, 0x66, 0x5f, 0x72, 0x65, 0x62, 0x65,
	0x6c, 0x73, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0d, 0x48, 0x01, 0x52, 0x0b, 0x6e, 0x75, 0x6d, 0x4f,
	0x66, 0x52, 0x65, 0x62, 0x65, 0x6c, 0x73, 0x88, 0x01, 0x01, 0x12, 0x44, 0x0a, 0x12, 0x73, 0x69,
	0x6e, 0x67, 0x6c, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x76, 0x65, 0x63, 0x74, 0x6f, 0x72,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54,
	0x69, 0x6d, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x48, 0x02, 0x52, 0x10, 0x73, 0x69, 0x6e,
	0x67, 0x6c, 0x65, 0x54, 0x69, 0x6d, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x88, 0x01, 0x01,
	0x12, 0x3d, 0x0a, 0x11, 0x6d, 0x75, 0x6c, 0x74, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x76, 0x65,
	0x63, 0x74, 0x6f, 0x72, 0x73, 0x18, 0x04, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x11, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x69, 0x6d, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x52, 0x0f,
	0x6d, 0x75, 0x6c, 0x74, 0x54, 0x69, 0x6d, 0x65, 0x56, 0x65, 0x63, 0x74, 0x6f, 0x72, 0x73, 0x42,
	0x18, 0x0a, 0x16, 0x5f, 0x66, 0x75, 0x6c, 0x63, 0x72, 0x75, 0x6d, 0x5f, 0x72, 0x65, 0x64, 0x69,
	0x72, 0x65, 0x63, 0x74, 0x5f, 0x61, 0x64, 0x64, 0x72, 0x42, 0x10, 0x0a, 0x0e, 0x5f, 0x6e, 0x75,
	0x6d, 0x5f, 0x6f, 0x66, 0x5f, 0x72, 0x65, 0x62, 0x65, 0x6c, 0x73, 0x42, 0x15, 0x0a, 0x13, 0x5f,
	0x73, 0x69, 0x6e, 0x67, 0x6c, 0x65, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x5f, 0x76, 0x65, 0x63, 0x74,
	0x6f, 0x72, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x2a, 0x8f, 0x01, 0x0a, 0x07,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x0c, 0x0a, 0x08, 0x41, 0x44, 0x44, 0x5f, 0x43,
	0x49, 0x54, 0x59, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x44, 0x45, 0x4c, 0x45, 0x54, 0x45, 0x5f,
	0x43, 0x49, 0x54, 0x59, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x50, 0x44, 0x41, 0x54, 0x45,
	0x5f, 0x4e, 0x41, 0x4d, 0x45, 0x10, 0x02, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x50, 0x44, 0x41, 0x54,
	0x45, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x10, 0x03, 0x12, 0x15, 0x0a, 0x11, 0x47, 0x45,
	0x54, 0x5f, 0x4e, 0x55, 0x4d, 0x42, 0x45, 0x52, 0x5f, 0x52, 0x45, 0x42, 0x45, 0x4c, 0x53, 0x10,
	0x04, 0x12, 0x15, 0x0a, 0x11, 0x43, 0x48, 0x45, 0x43, 0x4b, 0x5f, 0x43, 0x4f, 0x4e, 0x53, 0x49,
	0x53, 0x54, 0x45, 0x4e, 0x43, 0x59, 0x10, 0x05, 0x12, 0x13, 0x0a, 0x0f, 0x47, 0x45, 0x54, 0x5f,
	0x50, 0x4c, 0x41, 0x4e, 0x45, 0x54, 0x5f, 0x54, 0x49, 0x4d, 0x45, 0x10, 0x06, 0x32, 0x8d, 0x01,
	0x0a, 0x0d, 0x43, 0x6f, 0x6d, 0x6d, 0x75, 0x6e, 0x69, 0x63, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x12,
	0x3a, 0x0a, 0x0a, 0x52, 0x75, 0x6e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x14, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x50, 0x61, 0x72,
	0x61, 0x6d, 0x73, 0x1a, 0x16, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x75, 0x6c, 0x63,
	0x72, 0x75, 0x6d, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x40, 0x0a, 0x10, 0x42,
	0x72, 0x6f, 0x61, 0x64, 0x63, 0x61, 0x73, 0x74, 0x43, 0x68, 0x61, 0x6e, 0x67, 0x65, 0x73, 0x12,
	0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46, 0x75, 0x6c, 0x63, 0x72, 0x75, 0x6d, 0x48,
	0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x1a, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x46,
	0x75, 0x6c, 0x63, 0x72, 0x75, 0x6d, 0x48, 0x69, 0x73, 0x74, 0x6f, 0x72, 0x79, 0x42, 0x09, 0x5a,
	0x07, 0x64, 0x69, 0x73, 0x74, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_pb_services_proto_rawDescOnce sync.Once
	file_pb_services_proto_rawDescData = file_pb_services_proto_rawDesc
)

func file_pb_services_proto_rawDescGZIP() []byte {
	file_pb_services_proto_rawDescOnce.Do(func() {
		file_pb_services_proto_rawDescData = protoimpl.X.CompressGZIP(file_pb_services_proto_rawDescData)
	})
	return file_pb_services_proto_rawDescData
}

var file_pb_services_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_pb_services_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_pb_services_proto_goTypes = []interface{}{
	(Command)(0),            // 0: proto.Command
	(*TimeVector)(nil),      // 1: proto.TimeVector
	(*CommandParams)(nil),   // 2: proto.CommandParams
	(*FulcrumHistory)(nil),  // 3: proto.FulcrumHistory
	(*FulcrumResponse)(nil), // 4: proto.FulcrumResponse
	(*Empty)(nil),           // 5: proto.Empty
}
var file_pb_services_proto_depIdxs = []int32{
	0, // 0: proto.CommandParams.command:type_name -> proto.Command
	1, // 1: proto.CommandParams.last_time_vector:type_name -> proto.TimeVector
	2, // 2: proto.FulcrumHistory.history:type_name -> proto.CommandParams
	1, // 3: proto.FulcrumResponse.single_time_vector:type_name -> proto.TimeVector
	1, // 4: proto.FulcrumResponse.mult_time_vectors:type_name -> proto.TimeVector
	2, // 5: proto.Communication.RunCommand:input_type -> proto.CommandParams
	3, // 6: proto.Communication.BroadcastChanges:input_type -> proto.FulcrumHistory
	4, // 7: proto.Communication.RunCommand:output_type -> proto.FulcrumResponse
	3, // 8: proto.Communication.BroadcastChanges:output_type -> proto.FulcrumHistory
	7, // [7:9] is the sub-list for method output_type
	5, // [5:7] is the sub-list for method input_type
	5, // [5:5] is the sub-list for extension type_name
	5, // [5:5] is the sub-list for extension extendee
	0, // [0:5] is the sub-list for field type_name
}

func init() { file_pb_services_proto_init() }
func file_pb_services_proto_init() {
	if File_pb_services_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pb_services_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TimeVector); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_services_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommandParams); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_services_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FulcrumHistory); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_services_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*FulcrumResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pb_services_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	file_pb_services_proto_msgTypes[1].OneofWrappers = []interface{}{}
	file_pb_services_proto_msgTypes[3].OneofWrappers = []interface{}{}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pb_services_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pb_services_proto_goTypes,
		DependencyIndexes: file_pb_services_proto_depIdxs,
		EnumInfos:         file_pb_services_proto_enumTypes,
		MessageInfos:      file_pb_services_proto_msgTypes,
	}.Build()
	File_pb_services_proto = out.File
	file_pb_services_proto_rawDesc = nil
	file_pb_services_proto_goTypes = nil
	file_pb_services_proto_depIdxs = nil
}
