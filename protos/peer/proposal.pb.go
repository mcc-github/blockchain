


package peer 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

























type SignedProposal struct {
	
	ProposalBytes []byte `protobuf:"bytes,1,opt,name=proposal_bytes,json=proposalBytes,proto3" json:"proposal_bytes,omitempty"`
	
	
	
	Signature            []byte   `protobuf:"bytes,2,opt,name=signature,proto3" json:"signature,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *SignedProposal) Reset()         { *m = SignedProposal{} }
func (m *SignedProposal) String() string { return proto.CompactTextString(m) }
func (*SignedProposal) ProtoMessage()    {}
func (*SignedProposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_proposal_8e93dae9c9555ac0, []int{0}
}
func (m *SignedProposal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedProposal.Unmarshal(m, b)
}
func (m *SignedProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedProposal.Marshal(b, m, deterministic)
}
func (dst *SignedProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedProposal.Merge(dst, src)
}
func (m *SignedProposal) XXX_Size() int {
	return xxx_messageInfo_SignedProposal.Size(m)
}
func (m *SignedProposal) XXX_DiscardUnknown() {
	xxx_messageInfo_SignedProposal.DiscardUnknown(m)
}

var xxx_messageInfo_SignedProposal proto.InternalMessageInfo

func (m *SignedProposal) GetProposalBytes() []byte {
	if m != nil {
		return m.ProposalBytes
	}
	return nil
}

func (m *SignedProposal) GetSignature() []byte {
	if m != nil {
		return m.Signature
	}
	return nil
}




















type Proposal struct {
	
	Header []byte `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	
	
	Payload []byte `protobuf:"bytes,2,opt,name=payload,proto3" json:"payload,omitempty"`
	
	
	
	Extension            []byte   `protobuf:"bytes,3,opt,name=extension,proto3" json:"extension,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Proposal) Reset()         { *m = Proposal{} }
func (m *Proposal) String() string { return proto.CompactTextString(m) }
func (*Proposal) ProtoMessage()    {}
func (*Proposal) Descriptor() ([]byte, []int) {
	return fileDescriptor_proposal_8e93dae9c9555ac0, []int{1}
}
func (m *Proposal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Proposal.Unmarshal(m, b)
}
func (m *Proposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Proposal.Marshal(b, m, deterministic)
}
func (dst *Proposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Proposal.Merge(dst, src)
}
func (m *Proposal) XXX_Size() int {
	return xxx_messageInfo_Proposal.Size(m)
}
func (m *Proposal) XXX_DiscardUnknown() {
	xxx_messageInfo_Proposal.DiscardUnknown(m)
}

var xxx_messageInfo_Proposal proto.InternalMessageInfo

func (m *Proposal) GetHeader() []byte {
	if m != nil {
		return m.Header
	}
	return nil
}

func (m *Proposal) GetPayload() []byte {
	if m != nil {
		return m.Payload
	}
	return nil
}

func (m *Proposal) GetExtension() []byte {
	if m != nil {
		return m.Extension
	}
	return nil
}




type ChaincodeHeaderExtension struct {
	
	
	
	
	
	
	
	
	
	
	
	PayloadVisibility []byte `protobuf:"bytes,1,opt,name=payload_visibility,json=payloadVisibility,proto3" json:"payload_visibility,omitempty"`
	
	ChaincodeId          *ChaincodeID `protobuf:"bytes,2,opt,name=chaincode_id,json=chaincodeId,proto3" json:"chaincode_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ChaincodeHeaderExtension) Reset()         { *m = ChaincodeHeaderExtension{} }
func (m *ChaincodeHeaderExtension) String() string { return proto.CompactTextString(m) }
func (*ChaincodeHeaderExtension) ProtoMessage()    {}
func (*ChaincodeHeaderExtension) Descriptor() ([]byte, []int) {
	return fileDescriptor_proposal_8e93dae9c9555ac0, []int{2}
}
func (m *ChaincodeHeaderExtension) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeHeaderExtension.Unmarshal(m, b)
}
func (m *ChaincodeHeaderExtension) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeHeaderExtension.Marshal(b, m, deterministic)
}
func (dst *ChaincodeHeaderExtension) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeHeaderExtension.Merge(dst, src)
}
func (m *ChaincodeHeaderExtension) XXX_Size() int {
	return xxx_messageInfo_ChaincodeHeaderExtension.Size(m)
}
func (m *ChaincodeHeaderExtension) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeHeaderExtension.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeHeaderExtension proto.InternalMessageInfo

func (m *ChaincodeHeaderExtension) GetPayloadVisibility() []byte {
	if m != nil {
		return m.PayloadVisibility
	}
	return nil
}

func (m *ChaincodeHeaderExtension) GetChaincodeId() *ChaincodeID {
	if m != nil {
		return m.ChaincodeId
	}
	return nil
}




type ChaincodeProposalPayload struct {
	
	
	
	Input []byte `protobuf:"bytes,1,opt,name=input,proto3" json:"input,omitempty"`
	
	
	
	
	TransientMap         map[string][]byte `protobuf:"bytes,2,rep,name=TransientMap,proto3" json:"TransientMap,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *ChaincodeProposalPayload) Reset()         { *m = ChaincodeProposalPayload{} }
func (m *ChaincodeProposalPayload) String() string { return proto.CompactTextString(m) }
func (*ChaincodeProposalPayload) ProtoMessage()    {}
func (*ChaincodeProposalPayload) Descriptor() ([]byte, []int) {
	return fileDescriptor_proposal_8e93dae9c9555ac0, []int{3}
}
func (m *ChaincodeProposalPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeProposalPayload.Unmarshal(m, b)
}
func (m *ChaincodeProposalPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeProposalPayload.Marshal(b, m, deterministic)
}
func (dst *ChaincodeProposalPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeProposalPayload.Merge(dst, src)
}
func (m *ChaincodeProposalPayload) XXX_Size() int {
	return xxx_messageInfo_ChaincodeProposalPayload.Size(m)
}
func (m *ChaincodeProposalPayload) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeProposalPayload.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeProposalPayload proto.InternalMessageInfo

func (m *ChaincodeProposalPayload) GetInput() []byte {
	if m != nil {
		return m.Input
	}
	return nil
}

func (m *ChaincodeProposalPayload) GetTransientMap() map[string][]byte {
	if m != nil {
		return m.TransientMap
	}
	return nil
}



type ChaincodeAction struct {
	
	
	Results []byte `protobuf:"bytes,1,opt,name=results,proto3" json:"results,omitempty"`
	
	
	Events []byte `protobuf:"bytes,2,opt,name=events,proto3" json:"events,omitempty"`
	
	Response *Response `protobuf:"bytes,3,opt,name=response,proto3" json:"response,omitempty"`
	
	
	
	
	
	ChaincodeId          *ChaincodeID `protobuf:"bytes,4,opt,name=chaincode_id,json=chaincodeId,proto3" json:"chaincode_id,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ChaincodeAction) Reset()         { *m = ChaincodeAction{} }
func (m *ChaincodeAction) String() string { return proto.CompactTextString(m) }
func (*ChaincodeAction) ProtoMessage()    {}
func (*ChaincodeAction) Descriptor() ([]byte, []int) {
	return fileDescriptor_proposal_8e93dae9c9555ac0, []int{4}
}
func (m *ChaincodeAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeAction.Unmarshal(m, b)
}
func (m *ChaincodeAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeAction.Marshal(b, m, deterministic)
}
func (dst *ChaincodeAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeAction.Merge(dst, src)
}
func (m *ChaincodeAction) XXX_Size() int {
	return xxx_messageInfo_ChaincodeAction.Size(m)
}
func (m *ChaincodeAction) XXX_DiscardUnknown() {
	xxx_messageInfo_ChaincodeAction.DiscardUnknown(m)
}

var xxx_messageInfo_ChaincodeAction proto.InternalMessageInfo

func (m *ChaincodeAction) GetResults() []byte {
	if m != nil {
		return m.Results
	}
	return nil
}

func (m *ChaincodeAction) GetEvents() []byte {
	if m != nil {
		return m.Events
	}
	return nil
}

func (m *ChaincodeAction) GetResponse() *Response {
	if m != nil {
		return m.Response
	}
	return nil
}

func (m *ChaincodeAction) GetChaincodeId() *ChaincodeID {
	if m != nil {
		return m.ChaincodeId
	}
	return nil
}

func init() {
	proto.RegisterType((*SignedProposal)(nil), "protos.SignedProposal")
	proto.RegisterType((*Proposal)(nil), "protos.Proposal")
	proto.RegisterType((*ChaincodeHeaderExtension)(nil), "protos.ChaincodeHeaderExtension")
	proto.RegisterType((*ChaincodeProposalPayload)(nil), "protos.ChaincodeProposalPayload")
	proto.RegisterMapType((map[string][]byte)(nil), "protos.ChaincodeProposalPayload.TransientMapEntry")
	proto.RegisterType((*ChaincodeAction)(nil), "protos.ChaincodeAction")
}

func init() { proto.RegisterFile("peer/proposal.proto", fileDescriptor_proposal_8e93dae9c9555ac0) }

var fileDescriptor_proposal_8e93dae9c9555ac0 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x4d, 0x6f, 0xd3, 0x4c,
	0x10, 0x96, 0x93, 0xf7, 0xed, 0xc7, 0x24, 0xf4, 0x63, 0x5b, 0x21, 0x2b, 0xea, 0xa1, 0xb2, 0x84,
	0x54, 0x24, 0xb0, 0xa5, 0x20, 0x21, 0xc4, 0x05, 0x11, 0xa8, 0x44, 0x0f, 0x48, 0x95, 0x81, 0x1e,
	0x7a, 0x09, 0x6b, 0x7b, 0x70, 0x56, 0x35, 0xbb, 0xab, 0xdd, 0x75, 0x84, 0x8f, 0xfc, 0x1c, 0x7e,
	0x0a, 0xff, 0x0a, 0xd9, 0xfb, 0xd1, 0x94, 0x5c, 0x38, 0x25, 0x33, 0xf3, 0x3c, 0xcf, 0xcc, 0x3c,
	0xb3, 0x86, 0x13, 0x89, 0xa8, 0x32, 0xa9, 0x84, 0x14, 0x9a, 0x36, 0xa9, 0x54, 0xc2, 0x08, 0xb2,
	0x33, 0xfc, 0xe8, 0xd9, 0xe9, 0x50, 0x2c, 0x57, 0x94, 0xf1, 0x52, 0x54, 0x68, 0xab, 0xb3, 0xb3,
	0x07, 0x94, 0xa5, 0x42, 0x2d, 0x05, 0xd7, 0xae, 0x9a, 0x7c, 0x81, 0x83, 0x4f, 0xac, 0xe6, 0x58,
	0x5d, 0x3b, 0x00, 0x79, 0x02, 0x07, 0x01, 0x5c, 0x74, 0x06, 0x75, 0x1c, 0x9d, 0x47, 0x17, 0xd3,
	0xfc, 0x91, 0xcf, 0x2e, 0xfa, 0x24, 0x39, 0x83, 0x7d, 0xcd, 0x6a, 0x4e, 0x4d, 0xab, 0x30, 0x1e,
	0x0d, 0x88, 0xfb, 0x44, 0x72, 0x0b, 0x7b, 0x41, 0xf0, 0x31, 0xec, 0xac, 0x90, 0x56, 0xa8, 0x9c,
	0x90, 0x8b, 0x48, 0x0c, 0xbb, 0x92, 0x76, 0x8d, 0xa0, 0x95, 0xe3, 0xfb, 0xb0, 0xd7, 0xc6, 0x1f,
	0x06, 0xb9, 0x66, 0x82, 0xc7, 0x63, 0xab, 0x1d, 0x12, 0xc9, 0xcf, 0x08, 0xe2, 0x77, 0x7e, 0xc9,
	0x0f, 0x83, 0xd6, 0xa5, 0x2f, 0x92, 0xe7, 0x40, 0x9c, 0xca, 0x72, 0xcd, 0x34, 0x2b, 0x58, 0xc3,
	0x4c, 0xe7, 0x1a, 0x1f, 0xbb, 0xca, 0x4d, 0x28, 0x90, 0x97, 0x30, 0x0d, 0x7e, 0x2d, 0x99, 0x1d,
	0x64, 0x32, 0x3f, 0xb1, 0xe6, 0xe8, 0x34, 0xb4, 0xb9, 0x7a, 0x9f, 0x4f, 0x02, 0xf0, 0xaa, 0x4a,
	0x7e, 0x6f, 0xce, 0xe0, 0x37, 0xbd, 0x76, 0xe3, 0x9f, 0xc2, 0xff, 0x8c, 0xcb, 0xd6, 0xb8, 0xb6,
	0x36, 0x20, 0x37, 0x30, 0xfd, 0xac, 0x28, 0xd7, 0x0c, 0xb9, 0xf9, 0x48, 0x65, 0x3c, 0x3a, 0x1f,
	0x5f, 0x4c, 0xe6, 0xf3, 0xad, 0x56, 0x7f, 0xa9, 0xa5, 0x9b, 0xa4, 0x4b, 0x6e, 0x54, 0x97, 0x3f,
	0xd0, 0x99, 0xbd, 0x81, 0xe3, 0x2d, 0x08, 0x39, 0x82, 0xf1, 0x1d, 0xda, 0xbd, 0xf7, 0xf3, 0xfe,
	0x6f, 0x3f, 0xd4, 0x9a, 0x36, 0xad, 0xbf, 0x95, 0x0d, 0x5e, 0x8f, 0x5e, 0x45, 0xc9, 0xaf, 0x08,
	0x0e, 0x43, 0xf7, 0xb7, 0xa5, 0xe9, 0x6d, 0x8c, 0x61, 0x57, 0xa1, 0x6e, 0x1b, 0xe3, 0xaf, 0xef,
	0xc3, 0xfe, 0x9a, 0xb8, 0x46, 0x6e, 0xb4, 0x13, 0x72, 0x11, 0x79, 0x06, 0x7b, 0xfe, 0x69, 0x0d,
	0x27, 0x9b, 0xcc, 0x8f, 0xfc, 0x6a, 0xb9, 0xcb, 0xe7, 0x01, 0xb1, 0xe5, 0xfb, 0x7f, 0xff, 0xe6,
	0xfb, 0xe2, 0x2b, 0x24, 0x42, 0xd5, 0xe9, 0xaa, 0x93, 0xa8, 0x1a, 0xac, 0x6a, 0x54, 0xe9, 0x37,
	0x5a, 0x28, 0x56, 0x7a, 0x66, 0xff, 0xd8, 0x17, 0x87, 0xf7, 0x1e, 0x96, 0x77, 0xb4, 0xc6, 0xdb,
	0xa7, 0x35, 0x33, 0xab, 0xb6, 0x48, 0x4b, 0xf1, 0x3d, 0xdb, 0xe0, 0x66, 0x96, 0x9b, 0x59, 0x6e,
	0xd6, 0x73, 0x0b, 0xfb, 0x31, 0xbd, 0xf8, 0x13, 0x00, 0x00, 0xff, 0xff, 0x12, 0x75, 0xb6, 0xaf,
	0x6a, 0x03, 0x00, 0x00,
}
