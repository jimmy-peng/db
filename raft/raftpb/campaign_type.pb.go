// Code generated by protoc-gen-gogo.
// source: raft/raftpb/campaign_type.proto
// DO NOT EDIT!

/*
	Package raftpb is a generated protocol buffer package.

	It is generated from these files:
		raft/raftpb/campaign_type.proto
		raft/raftpb/config_change.proto
		raft/raftpb/entry.proto
		raft/raftpb/hard_state.proto
		raft/raftpb/message.proto
		raft/raftpb/message_type.proto
		raft/raftpb/node_state.proto
		raft/raftpb/progress_state.proto
		raft/raftpb/snapshot.proto
		raft/raftpb/soft_state.proto

	It has these top-level messages:
		ConfigChange
		Entry
		HardState
		Message
		ConfigState
		SnapshotMetadata
		Snapshot
		SoftState
*/
package raftpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

// (etcd raft.CampaignType)
type CAMPAIGN_TYPE int32

const (
	CAMPAIGN_TYPE_LEADER_ELECTION CAMPAIGN_TYPE = 0
	CAMPAIGN_TYPE_LEADER_TRANSFER CAMPAIGN_TYPE = 1
)

var CAMPAIGN_TYPE_name = map[int32]string{
	0: "LEADER_ELECTION",
	1: "LEADER_TRANSFER",
}
var CAMPAIGN_TYPE_value = map[string]int32{
	"LEADER_ELECTION": 0,
	"LEADER_TRANSFER": 1,
}

func (x CAMPAIGN_TYPE) String() string {
	return proto.EnumName(CAMPAIGN_TYPE_name, int32(x))
}
func (CAMPAIGN_TYPE) EnumDescriptor() ([]byte, []int) { return fileDescriptorCampaignType, []int{0} }

func init() {
	proto.RegisterEnum("raftpb.CAMPAIGN_TYPE", CAMPAIGN_TYPE_name, CAMPAIGN_TYPE_value)
}

func init() { proto.RegisterFile("raft/raftpb/campaign_type.proto", fileDescriptorCampaignType) }

var fileDescriptorCampaignType = []byte{
	// 167 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xe2, 0x92, 0x2f, 0x4a, 0x4c, 0x2b,
	0xd1, 0x07, 0x11, 0x05, 0x49, 0xfa, 0xc9, 0x89, 0xb9, 0x05, 0x89, 0x99, 0xe9, 0x79, 0xf1, 0x25,
	0x95, 0x05, 0xa9, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0x6c, 0x10, 0x39, 0x29, 0xdd, 0xf4,
	0xcc, 0x92, 0x8c, 0xd2, 0x24, 0xbd, 0xe4, 0xfc, 0x5c, 0xfd, 0xf4, 0xfc, 0xf4, 0x7c, 0x7d, 0xb0,
	0x74, 0x52, 0x69, 0x1a, 0x98, 0x07, 0xe6, 0x80, 0x59, 0x10, 0x6d, 0x5a, 0x96, 0x5c, 0xbc, 0xce,
	0x8e, 0xbe, 0x01, 0x8e, 0x9e, 0xee, 0x7e, 0xf1, 0x21, 0x91, 0x01, 0xae, 0x42, 0xc2, 0x5c, 0xfc,
	0x3e, 0xae, 0x8e, 0x2e, 0xae, 0x41, 0xf1, 0xae, 0x3e, 0xae, 0xce, 0x21, 0x9e, 0xfe, 0x7e, 0x02,
	0x0c, 0x48, 0x82, 0x21, 0x41, 0x8e, 0x7e, 0xc1, 0x6e, 0xae, 0x41, 0x02, 0x8c, 0x4e, 0x22, 0x27,
	0x1e, 0xca, 0x31, 0x9c, 0x78, 0x24, 0xc7, 0x78, 0xe1, 0x91, 0x1c, 0xe3, 0x83, 0x47, 0x72, 0x8c,
	0x33, 0x1e, 0xcb, 0x31, 0x24, 0xb1, 0x81, 0xcd, 0x35, 0x06, 0x04, 0x00, 0x00, 0xff, 0xff, 0xfb,
	0x02, 0x6d, 0x0e, 0xb1, 0x00, 0x00, 0x00,
}
