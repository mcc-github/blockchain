


package peer

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion3 

























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
	return fileDescriptor_c4dbb4372a94bd5b, []int{0}
}

func (m *SignedProposal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SignedProposal.Unmarshal(m, b)
}
func (m *SignedProposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SignedProposal.Marshal(b, m, deterministic)
}
func (m *SignedProposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SignedProposal.Merge(m, src)
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
	return fileDescriptor_c4dbb4372a94bd5b, []int{1}
}

func (m *Proposal) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Proposal.Unmarshal(m, b)
}
func (m *Proposal) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Proposal.Marshal(b, m, deterministic)
}
func (m *Proposal) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Proposal.Merge(m, src)
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
	return fileDescriptor_c4dbb4372a94bd5b, []int{2}
}

func (m *ChaincodeHeaderExtension) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeHeaderExtension.Unmarshal(m, b)
}
func (m *ChaincodeHeaderExtension) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeHeaderExtension.Marshal(b, m, deterministic)
}
func (m *ChaincodeHeaderExtension) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeHeaderExtension.Merge(m, src)
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
	return fileDescriptor_c4dbb4372a94bd5b, []int{3}
}

func (m *ChaincodeProposalPayload) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeProposalPayload.Unmarshal(m, b)
}
func (m *ChaincodeProposalPayload) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeProposalPayload.Marshal(b, m, deterministic)
}
func (m *ChaincodeProposalPayload) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeProposalPayload.Merge(m, src)
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
	return fileDescriptor_c4dbb4372a94bd5b, []int{4}
}

func (m *ChaincodeAction) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ChaincodeAction.Unmarshal(m, b)
}
func (m *ChaincodeAction) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ChaincodeAction.Marshal(b, m, deterministic)
}
func (m *ChaincodeAction) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ChaincodeAction.Merge(m, src)
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

func init() { proto.RegisterFile("peer/proposal.proto", fileDescriptor_c4dbb4372a94bd5b) }

var fileDescriptor_c4dbb4372a94bd5b = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x8c, 0x53, 0x51, 0x6b, 0x13, 0x41,
	0x10, 0x26, 0x49, 0x9b, 0xa6, 0x93, 0xd8, 0xa6, 0xdb, 0x22, 0x47, 0xe8, 0x43, 0x39, 0x10, 0x2a,
	0xe8, 0x1d, 0x44, 0x10, 0xf1, 0x45, 0x8c, 0x16, 0xac, 0x20, 0x94, 0x53, 0xfb, 0xd0, 0x97, 0xb8,
	0x77, 0x37, 0x5e, 0x96, 0x9c, 0xbb, 0xcb, 0xee, 0x5e, 0xf0, 0x1e, 0xfd, 0x69, 0x3e, 0xfa, 0xaf,
	0xe4, 0x6e, 0x77, 0xaf, 0xa9, 0x79, 0xf1, 0x29, 0x99, 0x99, 0xef, 0xfb, 0x66, 0xe6, 0x9b, 0x5b,
	0x38, 0x95, 0x88, 0x2a, 0x96, 0x4a, 0x48, 0xa1, 0x69, 0x19, 0x49, 0x25, 0x8c, 0x20, 0xc3, 0xf6,
	0x47, 0xcf, 0xce, 0xda, 0x62, 0xb6, 0xa2, 0x8c, 0x67, 0x22, 0x47, 0x5b, 0x9d, 0x9d, 0x3f, 0xa0,
	0x2c, 0x15, 0x6a, 0x29, 0xb8, 0x76, 0xd5, 0xf0, 0x2b, 0x1c, 0x7d, 0x66, 0x05, 0xc7, 0xfc, 0xc6,
	0x01, 0xc8, 0x13, 0x38, 0xea, 0xc0, 0x69, 0x6d, 0x50, 0x07, 0xbd, 0x8b, 0xde, 0xe5, 0x24, 0x79,
	0xe4, 0xb3, 0x8b, 0x26, 0x49, 0xce, 0xe1, 0x50, 0xb3, 0x82, 0x53, 0x53, 0x29, 0x0c, 0xfa, 0x2d,
	0xe2, 0x3e, 0x11, 0xde, 0xc1, 0xa8, 0x13, 0x7c, 0x0c, 0xc3, 0x15, 0xd2, 0x1c, 0x95, 0x13, 0x72,
	0x11, 0x09, 0xe0, 0x40, 0xd2, 0xba, 0x14, 0x34, 0x77, 0x7c, 0x1f, 0x36, 0xda, 0xf8, 0xd3, 0x20,
	0xd7, 0x4c, 0xf0, 0x60, 0x60, 0xb5, 0xbb, 0x44, 0xf8, 0xab, 0x07, 0xc1, 0x3b, 0xbf, 0xe4, 0x87,
	0x56, 0xeb, 0xca, 0x17, 0xc9, 0x73, 0x20, 0x4e, 0x65, 0xb9, 0x61, 0x9a, 0xa5, 0xac, 0x64, 0xa6,
	0x76, 0x8d, 0x4f, 0x5c, 0xe5, 0xb6, 0x2b, 0x90, 0x97, 0x30, 0xe9, 0xfc, 0x5a, 0x32, 0x3b, 0xc8,
	0x78, 0x7e, 0x6a, 0xcd, 0xd1, 0x51, 0xd7, 0xe6, 0xfa, 0x7d, 0x32, 0xee, 0x80, 0xd7, 0x79, 0xf8,
	0x67, 0x7b, 0x06, 0xbf, 0xe9, 0x8d, 0x1b, 0xff, 0x0c, 0xf6, 0x19, 0x97, 0x95, 0x71, 0x6d, 0x6d,
	0x40, 0x6e, 0x61, 0xf2, 0x45, 0x51, 0xae, 0x19, 0x72, 0xf3, 0x89, 0xca, 0xa0, 0x7f, 0x31, 0xb8,
	0x1c, 0xcf, 0xe7, 0x3b, 0xad, 0xfe, 0x51, 0x8b, 0xb6, 0x49, 0x57, 0xdc, 0xa8, 0x3a, 0x79, 0xa0,
	0x33, 0x7b, 0x03, 0x27, 0x3b, 0x10, 0x32, 0x85, 0xc1, 0x1a, 0xed, 0xde, 0x87, 0x49, 0xf3, 0xb7,
	0x19, 0x6a, 0x43, 0xcb, 0xca, 0xdf, 0xca, 0x06, 0xaf, 0xfb, 0xaf, 0x7a, 0xe1, 0xef, 0x1e, 0x1c,
	0x77, 0xdd, 0xdf, 0x66, 0xa6, 0xb1, 0x31, 0x80, 0x03, 0x85, 0xba, 0x2a, 0x8d, 0xbf, 0xbe, 0x0f,
	0x9b, 0x6b, 0xe2, 0x06, 0xb9, 0xd1, 0x4e, 0xc8, 0x45, 0xe4, 0x19, 0x8c, 0xfc, 0xa7, 0xd5, 0x9e,
	0x6c, 0x3c, 0x9f, 0xfa, 0xd5, 0x12, 0x97, 0x4f, 0x3a, 0xc4, 0x8e, 0xef, 0x7b, 0xff, 0xe7, 0xfb,
	0xc7, 0xbd, 0xd1, 0xfe, 0x74, 0x98, 0x4c, 0x8d, 0x58, 0x23, 0x5f, 0x0a, 0x89, 0x8a, 0x36, 0xe3,
	0xea, 0xc5, 0x37, 0x08, 0x85, 0x2a, 0xa2, 0x55, 0x2d, 0x51, 0x95, 0x98, 0x17, 0xa8, 0xa2, 0xef,
	0x34, 0x55, 0x2c, 0xf3, 0x8a, 0xcd, 0x23, 0x58, 0x1c, 0xdf, 0x7b, 0x9b, 0xad, 0x69, 0x81, 0x77,
	0x4f, 0x0b, 0x66, 0x56, 0x55, 0x1a, 0x65, 0xe2, 0x47, 0xbc, 0xc5, 0x8d, 0x2d, 0x37, 0xb6, 0xdc,
	0xb8, 0xe1, 0xa6, 0xf6, 0x91, 0xbd, 0xf8, 0x1b, 0x00, 0x00, 0xff, 0xff, 0x68, 0x52, 0x92, 0xde,
	0x82, 0x03, 0x00, 0x00,
}
