// Code generated by protoc-gen-gogo.
// source: raft/raftpb/message_type.proto
// DO NOT EDIT!

package raftpb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/gogo/protobuf/gogoproto"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// (etcd raft.raftpb.MessageType)
type MESSAGE_TYPE int32

const (
	MESSAGE_TYPE_INTERNAL_TRIGGER_CAMPAIGN                  MESSAGE_TYPE = 0
	MESSAGE_TYPE_INTERNAL_TRIGGER_LEADER_HEARTBEAT          MESSAGE_TYPE = 1
	MESSAGE_TYPE_INTERNAL_TRIGGER_CHECK_QUORUM              MESSAGE_TYPE = 2
	MESSAGE_TYPE_INTERNAL_LEADER_CANNOT_CONNECT_TO_FOLLOWER MESSAGE_TYPE = 3
	MESSAGE_TYPE_LEADER_HEARTBEAT                           MESSAGE_TYPE = 4
	MESSAGE_TYPE_RESPONSE_TO_LEADER_HEARTBEAT               MESSAGE_TYPE = 5
	MESSAGE_TYPE_CANDIDATE_REQUEST_VOTE                     MESSAGE_TYPE = 6
	MESSAGE_TYPE_RESPONSE_TO_CANDIDATE_REQUEST_VOTE         MESSAGE_TYPE = 7
	MESSAGE_TYPE_PROPOSAL_TO_LEADER                         MESSAGE_TYPE = 8
	MESSAGE_TYPE_LEADER_APPEND                              MESSAGE_TYPE = 9
	MESSAGE_TYPE_RESPONSE_TO_LEADER_APPEND                  MESSAGE_TYPE = 10
	MESSAGE_TYPE_LEADER_SNAPSHOT                            MESSAGE_TYPE = 11
	MESSAGE_TYPE_INTERNAL_RESPONSE_TO_LEADER_SNAPSHOT       MESSAGE_TYPE = 12
	MESSAGE_TYPE_TRANSFER_LEADER                            MESSAGE_TYPE = 13
	MESSAGE_TYPE_FORCE_ELECTION_TIMEOUT                     MESSAGE_TYPE = 14
	MESSAGE_TYPE_TRIGGER_READ_INDEX                         MESSAGE_TYPE = 15
	MESSAGE_TYPE_READ_INDEX_DATA                            MESSAGE_TYPE = 16
)

var MESSAGE_TYPE_name = map[int32]string{
	0:  "INTERNAL_TRIGGER_CAMPAIGN",
	1:  "INTERNAL_TRIGGER_LEADER_HEARTBEAT",
	2:  "INTERNAL_TRIGGER_CHECK_QUORUM",
	3:  "INTERNAL_LEADER_CANNOT_CONNECT_TO_FOLLOWER",
	4:  "LEADER_HEARTBEAT",
	5:  "RESPONSE_TO_LEADER_HEARTBEAT",
	6:  "CANDIDATE_REQUEST_VOTE",
	7:  "RESPONSE_TO_CANDIDATE_REQUEST_VOTE",
	8:  "PROPOSAL_TO_LEADER",
	9:  "LEADER_APPEND",
	10: "RESPONSE_TO_LEADER_APPEND",
	11: "LEADER_SNAPSHOT",
	12: "INTERNAL_RESPONSE_TO_LEADER_SNAPSHOT",
	13: "TRANSFER_LEADER",
	14: "FORCE_ELECTION_TIMEOUT",
	15: "TRIGGER_READ_INDEX",
	16: "READ_INDEX_DATA",
}
var MESSAGE_TYPE_value = map[string]int32{
	"INTERNAL_TRIGGER_CAMPAIGN":                  0,
	"INTERNAL_TRIGGER_LEADER_HEARTBEAT":          1,
	"INTERNAL_TRIGGER_CHECK_QUORUM":              2,
	"INTERNAL_LEADER_CANNOT_CONNECT_TO_FOLLOWER": 3,
	"LEADER_HEARTBEAT":                           4,
	"RESPONSE_TO_LEADER_HEARTBEAT":               5,
	"CANDIDATE_REQUEST_VOTE":                     6,
	"RESPONSE_TO_CANDIDATE_REQUEST_VOTE":         7,
	"PROPOSAL_TO_LEADER":                         8,
	"LEADER_APPEND":                              9,
	"RESPONSE_TO_LEADER_APPEND":                  10,
	"LEADER_SNAPSHOT":                            11,
	"INTERNAL_RESPONSE_TO_LEADER_SNAPSHOT":       12,
	"TRANSFER_LEADER":                            13,
	"FORCE_ELECTION_TIMEOUT":                     14,
	"TRIGGER_READ_INDEX":                         15,
	"READ_INDEX_DATA":                            16,
}

func (x MESSAGE_TYPE) String() string {
	return proto.EnumName(MESSAGE_TYPE_name, int32(x))
}
func (MESSAGE_TYPE) EnumDescriptor() ([]byte, []int) { return fileDescriptorMessageType, []int{0} }

func init() {
	proto.RegisterEnum("raftpb.MESSAGE_TYPE", MESSAGE_TYPE_name, MESSAGE_TYPE_value)
}

func init() { proto.RegisterFile("raft/raftpb/message_type.proto", fileDescriptorMessageType) }

var fileDescriptorMessageType = []byte{
	// 410 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x74, 0x52, 0x5d, 0x6e, 0xd4, 0x30,
	0x10, 0xde, 0x85, 0xb2, 0x80, 0x69, 0xa9, 0x31, 0x55, 0x25, 0x2a, 0x1a, 0x51, 0x04, 0x08, 0x55,
	0x62, 0xf7, 0x81, 0x13, 0x4c, 0x9d, 0xd9, 0xdd, 0x88, 0xc4, 0x76, 0xed, 0x09, 0x3f, 0x4f, 0x56,
	0x83, 0xb6, 0x81, 0x87, 0x2a, 0xab, 0x36, 0x7d, 0xe0, 0x0a, 0x9c, 0x80, 0x23, 0xf5, 0x91, 0x23,
	0xc0, 0x72, 0x11, 0xe4, 0x34, 0x0d, 0x48, 0x69, 0x5f, 0x2c, 0x8f, 0xbf, 0xef, 0x9b, 0x6f, 0xe6,
	0x93, 0x59, 0x74, 0x7a, 0x74, 0x5c, 0x4f, 0xc2, 0xb1, 0x2c, 0x26, 0x27, 0x8b, 0xb3, 0xb3, 0xa3,
	0x72, 0xe1, 0xeb, 0x6f, 0xcb, 0xc5, 0x78, 0x79, 0x5a, 0xd5, 0x95, 0x18, 0x5d, 0x42, 0x3b, 0x6f,
	0xca, 0xaf, 0xf5, 0x97, 0xf3, 0x62, 0xfc, 0xb9, 0x3a, 0x99, 0x94, 0x55, 0x59, 0x4d, 0x1a, 0xb8,
	0x38, 0x3f, 0x6e, 0xaa, 0xa6, 0x68, 0x6e, 0x97, 0xb2, 0xfd, 0xef, 0x6b, 0x6c, 0x3d, 0x43, 0xe7,
	0x60, 0x86, 0x9e, 0x3e, 0x19, 0x14, 0xbb, 0xec, 0x49, 0xa2, 0x08, 0xad, 0x82, 0xd4, 0x93, 0x4d,
	0x66, 0x33, 0xb4, 0x5e, 0x42, 0x66, 0x20, 0x99, 0x29, 0x3e, 0x10, 0x2f, 0xd9, 0x5e, 0x0f, 0x4e,
	0x11, 0x62, 0xb4, 0x7e, 0x8e, 0x60, 0xe9, 0x00, 0x81, 0xf8, 0x50, 0xec, 0xb1, 0xdd, 0x7e, 0x97,
	0x39, 0xca, 0x77, 0xfe, 0x30, 0xd7, 0x36, 0xcf, 0xf8, 0x2d, 0x31, 0x66, 0xfb, 0x1d, 0xa5, 0xed,
	0x20, 0x41, 0x29, 0x4d, 0x5e, 0x6a, 0xa5, 0x50, 0x92, 0x27, 0xed, 0xa7, 0x3a, 0x4d, 0xf5, 0x07,
	0xb4, 0xfc, 0xb6, 0xd8, 0x62, 0xbc, 0x67, 0xb4, 0x26, 0x9e, 0xb1, 0xa7, 0x16, 0x9d, 0xd1, 0xca,
	0x61, 0xe0, 0xf7, 0x18, 0x77, 0xc4, 0x0e, 0xdb, 0x96, 0xa0, 0xe2, 0x24, 0x06, 0x42, 0x6f, 0xf1,
	0x30, 0x47, 0x47, 0xfe, 0xbd, 0x26, 0xe4, 0x23, 0xf1, 0x8a, 0x3d, 0xff, 0x5f, 0x7d, 0x03, 0xef,
	0xae, 0xd8, 0x66, 0xc2, 0x58, 0x6d, 0xb4, 0x0b, 0xeb, 0x5c, 0xb9, 0xf0, 0x7b, 0xe2, 0x11, 0xdb,
	0x68, 0x1d, 0xc1, 0x18, 0x54, 0x31, 0xbf, 0x1f, 0xf2, 0xbb, 0x66, 0xa0, 0x16, 0x66, 0xe2, 0x31,
	0xdb, 0x6c, 0x9f, 0x9c, 0x02, 0xe3, 0xe6, 0x9a, 0xf8, 0x03, 0xf1, 0x9a, 0xbd, 0xe8, 0xa2, 0xb8,
	0x46, 0xdc, 0x31, 0xd7, 0x83, 0x9c, 0x2c, 0x28, 0x37, 0xed, 0x62, 0xe7, 0x1b, 0x61, 0xc3, 0xa9,
	0xb6, 0x12, 0x3d, 0xa6, 0x28, 0x29, 0xd1, 0xca, 0x53, 0x92, 0xa1, 0xce, 0x89, 0x3f, 0x0c, 0x93,
	0x5f, 0xe5, 0x6f, 0x11, 0x62, 0x9f, 0xa8, 0x18, 0x3f, 0xf2, 0xcd, 0xd0, 0xe8, 0x5f, 0xed, 0x63,
	0x20, 0xe0, 0xfc, 0x60, 0xeb, 0xe2, 0x77, 0x34, 0xb8, 0x58, 0x45, 0xc3, 0x9f, 0xab, 0x68, 0xf8,
	0x6b, 0x15, 0x0d, 0x7f, 0xfc, 0x89, 0x06, 0xc5, 0xa8, 0xf9, 0x29, 0x6f, 0xff, 0x06, 0x00, 0x00,
	0xff, 0xff, 0x85, 0x46, 0x70, 0xd8, 0x82, 0x02, 0x00, 0x00,
}
