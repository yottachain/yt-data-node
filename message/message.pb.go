// Code generated by protoc-gen-go. DO NOT EDIT.
// source: message/message.proto

package message

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

//msgid=0xc487
type NodeCapacityRequest struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeCapacityRequest) Reset()         { *m = NodeCapacityRequest{} }
func (m *NodeCapacityRequest) String() string { return proto.CompactTextString(m) }
func (*NodeCapacityRequest) ProtoMessage()    {}
func (*NodeCapacityRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{0}
}

func (m *NodeCapacityRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeCapacityRequest.Unmarshal(m, b)
}
func (m *NodeCapacityRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeCapacityRequest.Marshal(b, m, deterministic)
}
func (m *NodeCapacityRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeCapacityRequest.Merge(m, src)
}
func (m *NodeCapacityRequest) XXX_Size() int {
	return xxx_messageInfo_NodeCapacityRequest.Size(m)
}
func (m *NodeCapacityRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeCapacityRequest.DiscardUnknown(m)
}

var xxx_messageInfo_NodeCapacityRequest proto.InternalMessageInfo

//msgid=0xe684
type NodeCapacityResponse struct {
	Writable             bool     `protobuf:"varint,1,opt,name=writable,proto3" json:"writable,omitempty"`
	AllocId              string   `protobuf:"bytes,2,opt,name=allocId,proto3" json:"allocId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeCapacityResponse) Reset()         { *m = NodeCapacityResponse{} }
func (m *NodeCapacityResponse) String() string { return proto.CompactTextString(m) }
func (*NodeCapacityResponse) ProtoMessage()    {}
func (*NodeCapacityResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{1}
}

func (m *NodeCapacityResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeCapacityResponse.Unmarshal(m, b)
}
func (m *NodeCapacityResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeCapacityResponse.Marshal(b, m, deterministic)
}
func (m *NodeCapacityResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeCapacityResponse.Merge(m, src)
}
func (m *NodeCapacityResponse) XXX_Size() int {
	return xxx_messageInfo_NodeCapacityResponse.Size(m)
}
func (m *NodeCapacityResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeCapacityResponse.DiscardUnknown(m)
}

var xxx_messageInfo_NodeCapacityResponse proto.InternalMessageInfo

func (m *NodeCapacityResponse) GetWritable() bool {
	if m != nil {
		return m.Writable
	}
	return false
}

func (m *NodeCapacityResponse) GetAllocId() string {
	if m != nil {
		return m.AllocId
	}
	return ""
}

//msgid=0xCB05
type UploadShardRequest struct {
	SHARDID              int32    `protobuf:"varint,1,opt,name=SHARDID,proto3" json:"SHARDID,omitempty"`
	BPDID                int32    `protobuf:"varint,2,opt,name=BPDID,proto3" json:"BPDID,omitempty"`
	VBI                  int64    `protobuf:"varint,3,opt,name=VBI,proto3" json:"VBI,omitempty"`
	BPDSIGN              []byte   `protobuf:"bytes,4,opt,name=BPDSIGN,proto3" json:"BPDSIGN,omitempty"`
	DAT                  []byte   `protobuf:"bytes,5,opt,name=DAT,proto3" json:"DAT,omitempty"`
	VHF                  []byte   `protobuf:"bytes,6,opt,name=VHF,proto3" json:"VHF,omitempty"`
	USERSIGN             []byte   `protobuf:"bytes,7,opt,name=USERSIGN,proto3" json:"USERSIGN,omitempty"`
	AllocId              string   `protobuf:"bytes,8,opt,name=allocId,proto3" json:"allocId,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadShardRequest) Reset()         { *m = UploadShardRequest{} }
func (m *UploadShardRequest) String() string { return proto.CompactTextString(m) }
func (*UploadShardRequest) ProtoMessage()    {}
func (*UploadShardRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{2}
}

func (m *UploadShardRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadShardRequest.Unmarshal(m, b)
}
func (m *UploadShardRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadShardRequest.Marshal(b, m, deterministic)
}
func (m *UploadShardRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadShardRequest.Merge(m, src)
}
func (m *UploadShardRequest) XXX_Size() int {
	return xxx_messageInfo_UploadShardRequest.Size(m)
}
func (m *UploadShardRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadShardRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UploadShardRequest proto.InternalMessageInfo

func (m *UploadShardRequest) GetSHARDID() int32 {
	if m != nil {
		return m.SHARDID
	}
	return 0
}

func (m *UploadShardRequest) GetBPDID() int32 {
	if m != nil {
		return m.BPDID
	}
	return 0
}

func (m *UploadShardRequest) GetVBI() int64 {
	if m != nil {
		return m.VBI
	}
	return 0
}

func (m *UploadShardRequest) GetBPDSIGN() []byte {
	if m != nil {
		return m.BPDSIGN
	}
	return nil
}

func (m *UploadShardRequest) GetDAT() []byte {
	if m != nil {
		return m.DAT
	}
	return nil
}

func (m *UploadShardRequest) GetVHF() []byte {
	if m != nil {
		return m.VHF
	}
	return nil
}

func (m *UploadShardRequest) GetUSERSIGN() []byte {
	if m != nil {
		return m.USERSIGN
	}
	return nil
}

func (m *UploadShardRequest) GetAllocId() string {
	if m != nil {
		return m.AllocId
	}
	return ""
}

type UploadShardResponse struct {
	RES                  int32    `protobuf:"varint,1,opt,name=RES,proto3" json:"RES,omitempty"`
	SHARDID              int32    `protobuf:"varint,2,opt,name=SHARDID,proto3" json:"SHARDID,omitempty"`
	VBI                  int64    `protobuf:"varint,3,opt,name=VBI,proto3" json:"VBI,omitempty"`
	VHF                  []byte   `protobuf:"bytes,4,opt,name=VHF,proto3" json:"VHF,omitempty"`
	USERSIGN             []byte   `protobuf:"bytes,5,opt,name=USERSIGN,proto3" json:"USERSIGN,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadShardResponse) Reset()         { *m = UploadShardResponse{} }
func (m *UploadShardResponse) String() string { return proto.CompactTextString(m) }
func (*UploadShardResponse) ProtoMessage()    {}
func (*UploadShardResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{3}
}

func (m *UploadShardResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadShardResponse.Unmarshal(m, b)
}
func (m *UploadShardResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadShardResponse.Marshal(b, m, deterministic)
}
func (m *UploadShardResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadShardResponse.Merge(m, src)
}
func (m *UploadShardResponse) XXX_Size() int {
	return xxx_messageInfo_UploadShardResponse.Size(m)
}
func (m *UploadShardResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadShardResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UploadShardResponse proto.InternalMessageInfo

func (m *UploadShardResponse) GetRES() int32 {
	if m != nil {
		return m.RES
	}
	return 0
}

func (m *UploadShardResponse) GetSHARDID() int32 {
	if m != nil {
		return m.SHARDID
	}
	return 0
}

func (m *UploadShardResponse) GetVBI() int64 {
	if m != nil {
		return m.VBI
	}
	return 0
}

func (m *UploadShardResponse) GetVHF() []byte {
	if m != nil {
		return m.VHF
	}
	return nil
}

func (m *UploadShardResponse) GetUSERSIGN() []byte {
	if m != nil {
		return m.USERSIGN
	}
	return nil
}

type VoidResponse struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VoidResponse) Reset()         { *m = VoidResponse{} }
func (m *VoidResponse) String() string { return proto.CompactTextString(m) }
func (*VoidResponse) ProtoMessage()    {}
func (*VoidResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{4}
}

func (m *VoidResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VoidResponse.Unmarshal(m, b)
}
func (m *VoidResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VoidResponse.Marshal(b, m, deterministic)
}
func (m *VoidResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VoidResponse.Merge(m, src)
}
func (m *VoidResponse) XXX_Size() int {
	return xxx_messageInfo_VoidResponse.Size(m)
}
func (m *VoidResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VoidResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VoidResponse proto.InternalMessageInfo

type UploadShard2CResponse struct {
	RES                  int32    `protobuf:"varint,1,opt,name=RES,proto3" json:"RES,omitempty"`
	DNSIGN               []byte   `protobuf:"bytes,2,opt,name=DNSIGN,proto3" json:"DNSIGN,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UploadShard2CResponse) Reset()         { *m = UploadShard2CResponse{} }
func (m *UploadShard2CResponse) String() string { return proto.CompactTextString(m) }
func (*UploadShard2CResponse) ProtoMessage()    {}
func (*UploadShard2CResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{5}
}

func (m *UploadShard2CResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UploadShard2CResponse.Unmarshal(m, b)
}
func (m *UploadShard2CResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UploadShard2CResponse.Marshal(b, m, deterministic)
}
func (m *UploadShard2CResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UploadShard2CResponse.Merge(m, src)
}
func (m *UploadShard2CResponse) XXX_Size() int {
	return xxx_messageInfo_UploadShard2CResponse.Size(m)
}
func (m *UploadShard2CResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UploadShard2CResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UploadShard2CResponse proto.InternalMessageInfo

func (m *UploadShard2CResponse) GetRES() int32 {
	if m != nil {
		return m.RES
	}
	return 0
}

func (m *UploadShard2CResponse) GetDNSIGN() []byte {
	if m != nil {
		return m.DNSIGN
	}
	return nil
}

type DownloadShardRequest struct {
	VHF                  []byte   `protobuf:"bytes,1,opt,name=VHF,proto3" json:"VHF,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DownloadShardRequest) Reset()         { *m = DownloadShardRequest{} }
func (m *DownloadShardRequest) String() string { return proto.CompactTextString(m) }
func (*DownloadShardRequest) ProtoMessage()    {}
func (*DownloadShardRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{6}
}

func (m *DownloadShardRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DownloadShardRequest.Unmarshal(m, b)
}
func (m *DownloadShardRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DownloadShardRequest.Marshal(b, m, deterministic)
}
func (m *DownloadShardRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DownloadShardRequest.Merge(m, src)
}
func (m *DownloadShardRequest) XXX_Size() int {
	return xxx_messageInfo_DownloadShardRequest.Size(m)
}
func (m *DownloadShardRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DownloadShardRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DownloadShardRequest proto.InternalMessageInfo

func (m *DownloadShardRequest) GetVHF() []byte {
	if m != nil {
		return m.VHF
	}
	return nil
}

type DownloadShardResponse struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DownloadShardResponse) Reset()         { *m = DownloadShardResponse{} }
func (m *DownloadShardResponse) String() string { return proto.CompactTextString(m) }
func (*DownloadShardResponse) ProtoMessage()    {}
func (*DownloadShardResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{7}
}

func (m *DownloadShardResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DownloadShardResponse.Unmarshal(m, b)
}
func (m *DownloadShardResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DownloadShardResponse.Marshal(b, m, deterministic)
}
func (m *DownloadShardResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DownloadShardResponse.Merge(m, src)
}
func (m *DownloadShardResponse) XXX_Size() int {
	return xxx_messageInfo_DownloadShardResponse.Size(m)
}
func (m *DownloadShardResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DownloadShardResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DownloadShardResponse proto.InternalMessageInfo

func (m *DownloadShardResponse) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

type NodeRegReq struct {
	Nodeid               string   `protobuf:"bytes,1,opt,name=nodeid,proto3" json:"nodeid,omitempty"`
	Owner                string   `protobuf:"bytes,2,opt,name=owner,proto3" json:"owner,omitempty"`
	MaxDataSpace         uint64   `protobuf:"varint,3,opt,name=maxDataSpace,proto3" json:"maxDataSpace,omitempty"`
	Addrs                []string `protobuf:"bytes,4,rep,name=addrs,proto3" json:"addrs,omitempty"`
	Relay                bool     `protobuf:"varint,5,opt,name=relay,proto3" json:"relay,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeRegReq) Reset()         { *m = NodeRegReq{} }
func (m *NodeRegReq) String() string { return proto.CompactTextString(m) }
func (*NodeRegReq) ProtoMessage()    {}
func (*NodeRegReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{8}
}

func (m *NodeRegReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeRegReq.Unmarshal(m, b)
}
func (m *NodeRegReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeRegReq.Marshal(b, m, deterministic)
}
func (m *NodeRegReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeRegReq.Merge(m, src)
}
func (m *NodeRegReq) XXX_Size() int {
	return xxx_messageInfo_NodeRegReq.Size(m)
}
func (m *NodeRegReq) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeRegReq.DiscardUnknown(m)
}

var xxx_messageInfo_NodeRegReq proto.InternalMessageInfo

func (m *NodeRegReq) GetNodeid() string {
	if m != nil {
		return m.Nodeid
	}
	return ""
}

func (m *NodeRegReq) GetOwner() string {
	if m != nil {
		return m.Owner
	}
	return ""
}

func (m *NodeRegReq) GetMaxDataSpace() uint64 {
	if m != nil {
		return m.MaxDataSpace
	}
	return 0
}

func (m *NodeRegReq) GetAddrs() []string {
	if m != nil {
		return m.Addrs
	}
	return nil
}

func (m *NodeRegReq) GetRelay() bool {
	if m != nil {
		return m.Relay
	}
	return false
}

type NodeRegResp struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	AssignedSpace        uint64   `protobuf:"varint,2,opt,name=assignedSpace,proto3" json:"assignedSpace,omitempty"`
	RelayUrl             string   `protobuf:"bytes,3,opt,name=relayUrl,proto3" json:"relayUrl,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *NodeRegResp) Reset()         { *m = NodeRegResp{} }
func (m *NodeRegResp) String() string { return proto.CompactTextString(m) }
func (*NodeRegResp) ProtoMessage()    {}
func (*NodeRegResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{9}
}

func (m *NodeRegResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_NodeRegResp.Unmarshal(m, b)
}
func (m *NodeRegResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_NodeRegResp.Marshal(b, m, deterministic)
}
func (m *NodeRegResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_NodeRegResp.Merge(m, src)
}
func (m *NodeRegResp) XXX_Size() int {
	return xxx_messageInfo_NodeRegResp.Size(m)
}
func (m *NodeRegResp) XXX_DiscardUnknown() {
	xxx_messageInfo_NodeRegResp.DiscardUnknown(m)
}

var xxx_messageInfo_NodeRegResp proto.InternalMessageInfo

func (m *NodeRegResp) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *NodeRegResp) GetAssignedSpace() uint64 {
	if m != nil {
		return m.AssignedSpace
	}
	return 0
}

func (m *NodeRegResp) GetRelayUrl() string {
	if m != nil {
		return m.RelayUrl
	}
	return ""
}

type StatusRepReq struct {
	Id                   uint32   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	Cpu                  uint32   `protobuf:"varint,2,opt,name=cpu,proto3" json:"cpu,omitempty"`
	Memory               uint32   `protobuf:"varint,3,opt,name=memory,proto3" json:"memory,omitempty"`
	Bandwidth            uint32   `protobuf:"varint,4,opt,name=bandwidth,proto3" json:"bandwidth,omitempty"`
	MaxDataSpace         uint64   `protobuf:"varint,5,opt,name=maxDataSpace,proto3" json:"maxDataSpace,omitempty"`
	AssignedSpace        uint64   `protobuf:"varint,6,opt,name=assignedSpace,proto3" json:"assignedSpace,omitempty"`
	UsedSpace            uint64   `protobuf:"varint,7,opt,name=usedSpace,proto3" json:"usedSpace,omitempty"`
	Addrs                []string `protobuf:"bytes,8,rep,name=addrs,proto3" json:"addrs,omitempty"`
	Relay                bool     `protobuf:"varint,9,opt,name=relay,proto3" json:"relay,omitempty"`
	Version              uint32   `protobuf:"varint,10,opt,name=version,proto3" json:"version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusRepReq) Reset()         { *m = StatusRepReq{} }
func (m *StatusRepReq) String() string { return proto.CompactTextString(m) }
func (*StatusRepReq) ProtoMessage()    {}
func (*StatusRepReq) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{10}
}

func (m *StatusRepReq) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRepReq.Unmarshal(m, b)
}
func (m *StatusRepReq) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRepReq.Marshal(b, m, deterministic)
}
func (m *StatusRepReq) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRepReq.Merge(m, src)
}
func (m *StatusRepReq) XXX_Size() int {
	return xxx_messageInfo_StatusRepReq.Size(m)
}
func (m *StatusRepReq) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRepReq.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRepReq proto.InternalMessageInfo

func (m *StatusRepReq) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *StatusRepReq) GetCpu() uint32 {
	if m != nil {
		return m.Cpu
	}
	return 0
}

func (m *StatusRepReq) GetMemory() uint32 {
	if m != nil {
		return m.Memory
	}
	return 0
}

func (m *StatusRepReq) GetBandwidth() uint32 {
	if m != nil {
		return m.Bandwidth
	}
	return 0
}

func (m *StatusRepReq) GetMaxDataSpace() uint64 {
	if m != nil {
		return m.MaxDataSpace
	}
	return 0
}

func (m *StatusRepReq) GetAssignedSpace() uint64 {
	if m != nil {
		return m.AssignedSpace
	}
	return 0
}

func (m *StatusRepReq) GetUsedSpace() uint64 {
	if m != nil {
		return m.UsedSpace
	}
	return 0
}

func (m *StatusRepReq) GetAddrs() []string {
	if m != nil {
		return m.Addrs
	}
	return nil
}

func (m *StatusRepReq) GetRelay() bool {
	if m != nil {
		return m.Relay
	}
	return false
}

func (m *StatusRepReq) GetVersion() uint32 {
	if m != nil {
		return m.Version
	}
	return 0
}

type StatusRepResp struct {
	ProductiveSpace      uint64   `protobuf:"varint,1,opt,name=productiveSpace,proto3" json:"productiveSpace,omitempty"`
	RelayUrl             string   `protobuf:"bytes,2,opt,name=relayUrl,proto3" json:"relayUrl,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StatusRepResp) Reset()         { *m = StatusRepResp{} }
func (m *StatusRepResp) String() string { return proto.CompactTextString(m) }
func (*StatusRepResp) ProtoMessage()    {}
func (*StatusRepResp) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{11}
}

func (m *StatusRepResp) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StatusRepResp.Unmarshal(m, b)
}
func (m *StatusRepResp) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StatusRepResp.Marshal(b, m, deterministic)
}
func (m *StatusRepResp) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StatusRepResp.Merge(m, src)
}
func (m *StatusRepResp) XXX_Size() int {
	return xxx_messageInfo_StatusRepResp.Size(m)
}
func (m *StatusRepResp) XXX_DiscardUnknown() {
	xxx_messageInfo_StatusRepResp.DiscardUnknown(m)
}

var xxx_messageInfo_StatusRepResp proto.InternalMessageInfo

func (m *StatusRepResp) GetProductiveSpace() uint64 {
	if m != nil {
		return m.ProductiveSpace
	}
	return 0
}

func (m *StatusRepResp) GetRelayUrl() string {
	if m != nil {
		return m.RelayUrl
	}
	return ""
}

type P2PLocation struct {
	NodeId               string   `protobuf:"bytes,1,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Addrs                []string `protobuf:"bytes,2,rep,name=addrs,proto3" json:"addrs,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *P2PLocation) Reset()         { *m = P2PLocation{} }
func (m *P2PLocation) String() string { return proto.CompactTextString(m) }
func (*P2PLocation) ProtoMessage()    {}
func (*P2PLocation) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{12}
}

func (m *P2PLocation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_P2PLocation.Unmarshal(m, b)
}
func (m *P2PLocation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_P2PLocation.Marshal(b, m, deterministic)
}
func (m *P2PLocation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_P2PLocation.Merge(m, src)
}
func (m *P2PLocation) XXX_Size() int {
	return xxx_messageInfo_P2PLocation.Size(m)
}
func (m *P2PLocation) XXX_DiscardUnknown() {
	xxx_messageInfo_P2PLocation.DiscardUnknown(m)
}

var xxx_messageInfo_P2PLocation proto.InternalMessageInfo

func (m *P2PLocation) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *P2PLocation) GetAddrs() []string {
	if m != nil {
		return m.Addrs
	}
	return nil
}

type TaskDescription struct {
	Id                   int64          `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	DataHash             [][]byte       `protobuf:"bytes,2,rep,name=dataHash,proto3" json:"dataHash,omitempty"`
	ParityHash           [][]byte       `protobuf:"bytes,3,rep,name=parityHash,proto3" json:"parityHash,omitempty"`
	Locations            []*P2PLocation `protobuf:"bytes,4,rep,name=locations,proto3" json:"locations,omitempty"`
	RecoverId            []int32        `protobuf:"varint,5,rep,packed,name=recoverId,proto3" json:"recoverId,omitempty"`
	XXX_NoUnkeyedLiteral struct{}       `json:"-"`
	XXX_unrecognized     []byte         `json:"-"`
	XXX_sizecache        int32          `json:"-"`
}

func (m *TaskDescription) Reset()         { *m = TaskDescription{} }
func (m *TaskDescription) String() string { return proto.CompactTextString(m) }
func (*TaskDescription) ProtoMessage()    {}
func (*TaskDescription) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{13}
}

func (m *TaskDescription) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskDescription.Unmarshal(m, b)
}
func (m *TaskDescription) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskDescription.Marshal(b, m, deterministic)
}
func (m *TaskDescription) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskDescription.Merge(m, src)
}
func (m *TaskDescription) XXX_Size() int {
	return xxx_messageInfo_TaskDescription.Size(m)
}
func (m *TaskDescription) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskDescription.DiscardUnknown(m)
}

var xxx_messageInfo_TaskDescription proto.InternalMessageInfo

func (m *TaskDescription) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *TaskDescription) GetDataHash() [][]byte {
	if m != nil {
		return m.DataHash
	}
	return nil
}

func (m *TaskDescription) GetParityHash() [][]byte {
	if m != nil {
		return m.ParityHash
	}
	return nil
}

func (m *TaskDescription) GetLocations() []*P2PLocation {
	if m != nil {
		return m.Locations
	}
	return nil
}

func (m *TaskDescription) GetRecoverId() []int32 {
	if m != nil {
		return m.RecoverId
	}
	return nil
}

type TaskOpResult struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	RES                  int32    `protobuf:"varint,2,opt,name=RES,proto3" json:"RES,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *TaskOpResult) Reset()         { *m = TaskOpResult{} }
func (m *TaskOpResult) String() string { return proto.CompactTextString(m) }
func (*TaskOpResult) ProtoMessage()    {}
func (*TaskOpResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{14}
}

func (m *TaskOpResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_TaskOpResult.Unmarshal(m, b)
}
func (m *TaskOpResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_TaskOpResult.Marshal(b, m, deterministic)
}
func (m *TaskOpResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TaskOpResult.Merge(m, src)
}
func (m *TaskOpResult) XXX_Size() int {
	return xxx_messageInfo_TaskOpResult.Size(m)
}
func (m *TaskOpResult) XXX_DiscardUnknown() {
	xxx_messageInfo_TaskOpResult.DiscardUnknown(m)
}

var xxx_messageInfo_TaskOpResult proto.InternalMessageInfo

func (m *TaskOpResult) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *TaskOpResult) GetRES() int32 {
	if m != nil {
		return m.RES
	}
	return 0
}

type StringMsg struct {
	Msg                  string   `protobuf:"bytes,1,opt,name=msg,proto3" json:"msg,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *StringMsg) Reset()         { *m = StringMsg{} }
func (m *StringMsg) String() string { return proto.CompactTextString(m) }
func (*StringMsg) ProtoMessage()    {}
func (*StringMsg) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{15}
}

func (m *StringMsg) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StringMsg.Unmarshal(m, b)
}
func (m *StringMsg) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StringMsg.Marshal(b, m, deterministic)
}
func (m *StringMsg) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StringMsg.Merge(m, src)
}
func (m *StringMsg) XXX_Size() int {
	return xxx_messageInfo_StringMsg.Size(m)
}
func (m *StringMsg) XXX_DiscardUnknown() {
	xxx_messageInfo_StringMsg.DiscardUnknown(m)
}

var xxx_messageInfo_StringMsg proto.InternalMessageInfo

func (m *StringMsg) GetMsg() string {
	if m != nil {
		return m.Msg
	}
	return ""
}

type SpotCheckTaskList struct {
	TaskId               string           `protobuf:"bytes,1,opt,name=taskId,proto3" json:"taskId,omitempty"`
	Snid                 int32            `protobuf:"varint,2,opt,name=snid,proto3" json:"snid,omitempty"`
	TaskList             []*SpotCheckTask `protobuf:"bytes,3,rep,name=taskList,proto3" json:"taskList,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *SpotCheckTaskList) Reset()         { *m = SpotCheckTaskList{} }
func (m *SpotCheckTaskList) String() string { return proto.CompactTextString(m) }
func (*SpotCheckTaskList) ProtoMessage()    {}
func (*SpotCheckTaskList) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{16}
}

func (m *SpotCheckTaskList) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpotCheckTaskList.Unmarshal(m, b)
}
func (m *SpotCheckTaskList) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpotCheckTaskList.Marshal(b, m, deterministic)
}
func (m *SpotCheckTaskList) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpotCheckTaskList.Merge(m, src)
}
func (m *SpotCheckTaskList) XXX_Size() int {
	return xxx_messageInfo_SpotCheckTaskList.Size(m)
}
func (m *SpotCheckTaskList) XXX_DiscardUnknown() {
	xxx_messageInfo_SpotCheckTaskList.DiscardUnknown(m)
}

var xxx_messageInfo_SpotCheckTaskList proto.InternalMessageInfo

func (m *SpotCheckTaskList) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *SpotCheckTaskList) GetSnid() int32 {
	if m != nil {
		return m.Snid
	}
	return 0
}

func (m *SpotCheckTaskList) GetTaskList() []*SpotCheckTask {
	if m != nil {
		return m.TaskList
	}
	return nil
}

type SpotCheckTask struct {
	Id                   int32    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	NodeId               string   `protobuf:"bytes,2,opt,name=nodeId,proto3" json:"nodeId,omitempty"`
	Addr                 string   `protobuf:"bytes,3,opt,name=addr,proto3" json:"addr,omitempty"`
	VHF                  []byte   `protobuf:"bytes,4,opt,name=VHF,proto3" json:"VHF,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SpotCheckTask) Reset()         { *m = SpotCheckTask{} }
func (m *SpotCheckTask) String() string { return proto.CompactTextString(m) }
func (*SpotCheckTask) ProtoMessage()    {}
func (*SpotCheckTask) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{17}
}

func (m *SpotCheckTask) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpotCheckTask.Unmarshal(m, b)
}
func (m *SpotCheckTask) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpotCheckTask.Marshal(b, m, deterministic)
}
func (m *SpotCheckTask) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpotCheckTask.Merge(m, src)
}
func (m *SpotCheckTask) XXX_Size() int {
	return xxx_messageInfo_SpotCheckTask.Size(m)
}
func (m *SpotCheckTask) XXX_DiscardUnknown() {
	xxx_messageInfo_SpotCheckTask.DiscardUnknown(m)
}

var xxx_messageInfo_SpotCheckTask proto.InternalMessageInfo

func (m *SpotCheckTask) GetId() int32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *SpotCheckTask) GetNodeId() string {
	if m != nil {
		return m.NodeId
	}
	return ""
}

func (m *SpotCheckTask) GetAddr() string {
	if m != nil {
		return m.Addr
	}
	return ""
}

func (m *SpotCheckTask) GetVHF() []byte {
	if m != nil {
		return m.VHF
	}
	return nil
}

type SpotCheckStatus struct {
	TaskId               string   `protobuf:"bytes,1,opt,name=taskId,proto3" json:"taskId,omitempty"`
	InvalidNodeList      []int32  `protobuf:"varint,2,rep,packed,name=invalidNodeList,proto3" json:"invalidNodeList,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SpotCheckStatus) Reset()         { *m = SpotCheckStatus{} }
func (m *SpotCheckStatus) String() string { return proto.CompactTextString(m) }
func (*SpotCheckStatus) ProtoMessage()    {}
func (*SpotCheckStatus) Descriptor() ([]byte, []int) {
	return fileDescriptor_ebceca9e8703e37f, []int{18}
}

func (m *SpotCheckStatus) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SpotCheckStatus.Unmarshal(m, b)
}
func (m *SpotCheckStatus) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SpotCheckStatus.Marshal(b, m, deterministic)
}
func (m *SpotCheckStatus) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SpotCheckStatus.Merge(m, src)
}
func (m *SpotCheckStatus) XXX_Size() int {
	return xxx_messageInfo_SpotCheckStatus.Size(m)
}
func (m *SpotCheckStatus) XXX_DiscardUnknown() {
	xxx_messageInfo_SpotCheckStatus.DiscardUnknown(m)
}

var xxx_messageInfo_SpotCheckStatus proto.InternalMessageInfo

func (m *SpotCheckStatus) GetTaskId() string {
	if m != nil {
		return m.TaskId
	}
	return ""
}

func (m *SpotCheckStatus) GetInvalidNodeList() []int32 {
	if m != nil {
		return m.InvalidNodeList
	}
	return nil
}

func init() {
	proto.RegisterType((*NodeCapacityRequest)(nil), "message.NodeCapacityRequest")
	proto.RegisterType((*NodeCapacityResponse)(nil), "message.NodeCapacityResponse")
	proto.RegisterType((*UploadShardRequest)(nil), "message.UploadShardRequest")
	proto.RegisterType((*UploadShardResponse)(nil), "message.UploadShardResponse")
	proto.RegisterType((*VoidResponse)(nil), "message.VoidResponse")
	proto.RegisterType((*UploadShard2CResponse)(nil), "message.UploadShard2CResponse")
	proto.RegisterType((*DownloadShardRequest)(nil), "message.DownloadShardRequest")
	proto.RegisterType((*DownloadShardResponse)(nil), "message.DownloadShardResponse")
	proto.RegisterType((*NodeRegReq)(nil), "message.NodeRegReq")
	proto.RegisterType((*NodeRegResp)(nil), "message.NodeRegResp")
	proto.RegisterType((*StatusRepReq)(nil), "message.StatusRepReq")
	proto.RegisterType((*StatusRepResp)(nil), "message.StatusRepResp")
	proto.RegisterType((*P2PLocation)(nil), "message.P2PLocation")
	proto.RegisterType((*TaskDescription)(nil), "message.TaskDescription")
	proto.RegisterType((*TaskOpResult)(nil), "message.TaskOpResult")
	proto.RegisterType((*StringMsg)(nil), "message.StringMsg")
	proto.RegisterType((*SpotCheckTaskList)(nil), "message.SpotCheckTaskList")
	proto.RegisterType((*SpotCheckTask)(nil), "message.SpotCheckTask")
	proto.RegisterType((*SpotCheckStatus)(nil), "message.SpotCheckStatus")
}

func init() { proto.RegisterFile("message/message.proto", fileDescriptor_ebceca9e8703e37f) }

var fileDescriptor_ebceca9e8703e37f = []byte{
	// 823 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x7c, 0x55, 0xcd, 0x6e, 0xe3, 0x36,
	0x10, 0x86, 0x24, 0x3b, 0xb1, 0x27, 0xf6, 0x26, 0xe5, 0x26, 0x0b, 0x61, 0xb1, 0x2d, 0x0c, 0xa2,
	0x07, 0x03, 0x05, 0xb6, 0x85, 0x7b, 0xec, 0x29, 0x89, 0x77, 0x1b, 0x03, 0x69, 0x1a, 0x50, 0x9b,
	0xdc, 0x8a, 0x82, 0x11, 0x09, 0x87, 0x88, 0x2c, 0x6a, 0x45, 0x3a, 0x69, 0x8e, 0xbd, 0xf6, 0xdc,
	0xf7, 0xe8, 0x7b, 0xf4, 0xa9, 0x0a, 0x8e, 0x68, 0xfd, 0xd8, 0xe9, 0x9e, 0x3c, 0xdf, 0x70, 0x38,
	0x33, 0xdf, 0x37, 0x43, 0x19, 0x4e, 0x56, 0xd2, 0x18, 0xbe, 0x94, 0xdf, 0xfb, 0xdf, 0xf7, 0x45,
	0xa9, 0xad, 0x26, 0xfb, 0x1e, 0xd2, 0x13, 0x78, 0x7d, 0xa5, 0x85, 0x3c, 0xe7, 0x05, 0x4f, 0x95,
	0x7d, 0x66, 0xf2, 0xf3, 0x5a, 0x1a, 0x4b, 0x2f, 0xe1, 0xb8, 0xeb, 0x36, 0x85, 0xce, 0x8d, 0x24,
	0x6f, 0x61, 0xf0, 0x54, 0x2a, 0xcb, 0xef, 0x32, 0x19, 0x07, 0x93, 0x60, 0x3a, 0x60, 0x35, 0x26,
	0x31, 0xec, 0xf3, 0x2c, 0xd3, 0xe9, 0x42, 0xc4, 0xe1, 0x24, 0x98, 0x0e, 0xd9, 0x06, 0xd2, 0x7f,
	0x03, 0x20, 0x37, 0x45, 0xa6, 0xb9, 0x48, 0xee, 0x79, 0x29, 0x7c, 0x11, 0x77, 0x21, 0xb9, 0x38,
	0x65, 0xf3, 0xc5, 0x1c, 0x73, 0xf5, 0xd9, 0x06, 0x92, 0x63, 0xe8, 0x9f, 0x5d, 0x3b, 0x7f, 0x88,
	0xfe, 0x0a, 0x90, 0x23, 0x88, 0x6e, 0xcf, 0x16, 0x71, 0x34, 0x09, 0xa6, 0x11, 0x73, 0xa6, 0xcb,
	0x70, 0x76, 0x3d, 0x4f, 0x16, 0x3f, 0x5f, 0xc5, 0xbd, 0x49, 0x30, 0x1d, 0xb1, 0x0d, 0x74, 0xb1,
	0xf3, 0xd3, 0x4f, 0x71, 0x1f, 0xbd, 0xce, 0xc4, 0xdb, 0x17, 0x1f, 0xe3, 0xbd, 0xca, 0x73, 0x7b,
	0xf1, 0xd1, 0x91, 0xb9, 0x49, 0x3e, 0x30, 0xbc, 0xbe, 0x8f, 0xee, 0x1a, 0xb7, 0xc9, 0x0c, 0xba,
	0x64, 0xfe, 0x0c, 0xe0, 0x75, 0x87, 0x8c, 0x97, 0xe6, 0x08, 0x22, 0xf6, 0x21, 0xf1, 0x4c, 0x9c,
	0xd9, 0xe6, 0x17, 0x76, 0xf9, 0xed, 0x32, 0xf1, 0xdd, 0xf5, 0x5e, 0xee, 0xae, 0xdf, 0xed, 0x8e,
	0xbe, 0x82, 0xd1, 0xad, 0x56, 0x75, 0x6d, 0x7a, 0x0a, 0x27, 0xad, 0x96, 0x66, 0xe7, 0x5f, 0x68,
	0xea, 0x0d, 0xec, 0xcd, 0xaf, 0x30, 0x69, 0x88, 0x49, 0x3d, 0xa2, 0x53, 0x38, 0x9e, 0xeb, 0xa7,
	0x7c, 0x67, 0x48, 0xbe, 0xb1, 0xa0, 0x6e, 0x8c, 0x7e, 0x07, 0x27, 0x5b, 0x91, 0xbe, 0x18, 0x81,
	0x9e, 0xe0, 0x96, 0xfb, 0x58, 0xb4, 0xe9, 0x5f, 0x01, 0x80, 0xdb, 0x24, 0x26, 0x97, 0x4c, 0x7e,
	0x76, 0xd5, 0x73, 0x2d, 0xa4, 0x12, 0x18, 0x34, 0x64, 0x1e, 0xb9, 0x81, 0xeb, 0xa7, 0x5c, 0x96,
	0x7e, 0x73, 0x2a, 0x40, 0x28, 0x8c, 0x56, 0xfc, 0x8f, 0x39, 0xb7, 0x3c, 0x29, 0x78, 0x2a, 0x51,
	0xaf, 0x1e, 0xeb, 0xf8, 0xdc, 0x4d, 0x2e, 0x44, 0x69, 0xe2, 0xde, 0x24, 0x72, 0x37, 0x11, 0x38,
	0x6f, 0x29, 0x33, 0xfe, 0x8c, 0xca, 0x0d, 0x58, 0x05, 0xe8, 0xef, 0x70, 0x50, 0xf7, 0x62, 0x0a,
	0xf2, 0x0a, 0x42, 0xdf, 0xc8, 0x98, 0x85, 0x4a, 0x90, 0x6f, 0x61, 0xcc, 0x8d, 0x51, 0xcb, 0x5c,
	0x8a, 0xaa, 0x5e, 0x88, 0xf5, 0xba, 0x4e, 0x37, 0x17, 0xcc, 0x76, 0x53, 0x66, 0xd8, 0xd0, 0x90,
	0xd5, 0x98, 0xfe, 0x1d, 0xc2, 0x28, 0xb1, 0xdc, 0xae, 0x0d, 0x93, 0x85, 0xe3, 0xbb, 0x5d, 0xe2,
	0x08, 0xa2, 0xb4, 0x58, 0x63, 0xe2, 0x31, 0x73, 0xa6, 0x53, 0x64, 0x25, 0x57, 0xba, 0x7c, 0xc6,
	0x64, 0x63, 0xe6, 0x11, 0x79, 0x07, 0xc3, 0x3b, 0x9e, 0x8b, 0x27, 0x25, 0xec, 0x3d, 0xae, 0xc5,
	0x98, 0x35, 0x8e, 0x1d, 0x65, 0xfa, 0x2f, 0x28, 0xb3, 0x43, 0x67, 0xef, 0x25, 0x3a, 0xef, 0x60,
	0xb8, 0x36, 0x9b, 0x88, 0x7d, 0x8c, 0x68, 0x1c, 0x8d, 0xba, 0x83, 0x17, 0xd5, 0x1d, 0xb6, 0xd4,
	0x75, 0xeb, 0xfe, 0x28, 0x4b, 0xa3, 0x74, 0x1e, 0x03, 0xf6, 0xbb, 0x81, 0xf4, 0x06, 0xc6, 0x2d,
	0x55, 0x4c, 0x41, 0xa6, 0x70, 0x58, 0x94, 0x5a, 0xac, 0x53, 0xab, 0x1e, 0x65, 0x55, 0x3a, 0xc0,
	0xd2, 0xdb, 0xee, 0x8e, 0xda, 0xe1, 0x96, 0xda, 0x3f, 0xc1, 0xc1, 0xf5, 0xec, 0xfa, 0x52, 0xa7,
	0xdc, 0x2a, 0x9d, 0x6f, 0x76, 0x6b, 0xd1, 0xd9, 0xad, 0x85, 0x68, 0x38, 0x84, 0x2d, 0x0e, 0xf4,
	0x9f, 0x00, 0x0e, 0x3f, 0x71, 0xf3, 0x30, 0x97, 0x26, 0x2d, 0x55, 0x81, 0x19, 0x9a, 0x69, 0x45,
	0x38, 0xad, 0xb7, 0x30, 0x70, 0x4b, 0x7c, 0xc1, 0xcd, 0x3d, 0x5e, 0x1e, 0xb1, 0x1a, 0x93, 0x6f,
	0x00, 0x0a, 0x5e, 0x2a, 0xfb, 0x8c, 0xa7, 0x11, 0x9e, 0xb6, 0x3c, 0x64, 0x06, 0xc3, 0xcc, 0x77,
	0x56, 0xed, 0xe6, 0xc1, 0xec, 0xf8, 0xfd, 0xe6, 0x23, 0xdc, 0x6a, 0x9b, 0x35, 0x61, 0x6e, 0x16,
	0xa5, 0x4c, 0xf5, 0xa3, 0x2c, 0x17, 0x22, 0xee, 0x4f, 0xa2, 0x69, 0x9f, 0x35, 0x0e, 0xfa, 0x03,
	0x8c, 0x5c, 0xc3, 0xbf, 0x3a, 0x09, 0xd7, 0x99, 0xdd, 0xe9, 0xd6, 0xbf, 0xf5, 0xb0, 0x7e, 0xeb,
	0xf4, 0x6b, 0x18, 0x26, 0xb6, 0x54, 0xf9, 0xf2, 0x17, 0xb3, 0x74, 0xc7, 0x2b, 0xb3, 0xf4, 0xda,
	0x38, 0x93, 0x1a, 0xf8, 0x2a, 0x29, 0xb4, 0x3d, 0xbf, 0x97, 0xe9, 0x83, 0xcb, 0x7c, 0xa9, 0x8c,
	0x75, 0x2a, 0x5a, 0x6e, 0x1e, 0x1a, 0x15, 0x2b, 0xe4, 0x1e, 0xb7, 0xc9, 0x95, 0xf0, 0xe9, 0xd1,
	0x26, 0x33, 0x18, 0x58, 0x7f, 0x0f, 0x15, 0x38, 0x98, 0xbd, 0xa9, 0x29, 0x76, 0x32, 0xb3, 0x3a,
	0x8e, 0xfe, 0x06, 0xe3, 0xce, 0x51, 0x8b, 0x46, 0x1f, 0x69, 0x34, 0x63, 0x0c, 0x3b, 0x63, 0x24,
	0xd0, 0x73, 0x93, 0xf3, 0x6f, 0x0e, 0xed, 0xdd, 0xaf, 0x26, 0x4d, 0xe0, 0xb0, 0x4e, 0x5f, 0xed,
	0xdc, 0xff, 0x32, 0x9a, 0xc2, 0xa1, 0xca, 0x1f, 0x79, 0xa6, 0x84, 0xfb, 0x28, 0x20, 0x89, 0x10,
	0x35, 0xdf, 0x76, 0xdf, 0xed, 0xe1, 0x9f, 0xe6, 0x8f, 0xff, 0x05, 0x00, 0x00, 0xff, 0xff, 0xa5,
	0x53, 0x63, 0xf6, 0x4d, 0x07, 0x00, 0x00,
}
