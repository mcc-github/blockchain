



package plugin_go

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import google_protobuf "github.com/golang/protobuf/protoc-gen-go/descriptor"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 


type Version struct {
	Major *int32 `protobuf:"varint,1,opt,name=major" json:"major,omitempty"`
	Minor *int32 `protobuf:"varint,2,opt,name=minor" json:"minor,omitempty"`
	Patch *int32 `protobuf:"varint,3,opt,name=patch" json:"patch,omitempty"`
	
	
	Suffix               *string  `protobuf:"bytes,4,opt,name=suffix" json:"suffix,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Version) Reset()                    { *m = Version{} }
func (m *Version) String() string            { return proto.CompactTextString(m) }
func (*Version) ProtoMessage()               {}
func (*Version) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }
func (m *Version) Unmarshal(b []byte) error {
	return xxx_messageInfo_Version.Unmarshal(m, b)
}
func (m *Version) Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Version.Marshal(b, m, deterministic)
}
func (dst *Version) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Version.Merge(dst, src)
}
func (m *Version) XXX_Size() int {
	return xxx_messageInfo_Version.Size(m)
}
func (m *Version) XXX_DiscardUnknown() {
	xxx_messageInfo_Version.DiscardUnknown(m)
}

var xxx_messageInfo_Version proto.InternalMessageInfo

func (m *Version) GetMajor() int32 {
	if m != nil && m.Major != nil {
		return *m.Major
	}
	return 0
}

func (m *Version) GetMinor() int32 {
	if m != nil && m.Minor != nil {
		return *m.Minor
	}
	return 0
}

func (m *Version) GetPatch() int32 {
	if m != nil && m.Patch != nil {
		return *m.Patch
	}
	return 0
}

func (m *Version) GetSuffix() string {
	if m != nil && m.Suffix != nil {
		return *m.Suffix
	}
	return ""
}


type CodeGeneratorRequest struct {
	
	
	
	FileToGenerate []string `protobuf:"bytes,1,rep,name=file_to_generate,json=fileToGenerate" json:"file_to_generate,omitempty"`
	
	Parameter *string `protobuf:"bytes,2,opt,name=parameter" json:"parameter,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	ProtoFile []*google_protobuf.FileDescriptorProto `protobuf:"bytes,15,rep,name=proto_file,json=protoFile" json:"proto_file,omitempty"`
	
	CompilerVersion      *Version `protobuf:"bytes,3,opt,name=compiler_version,json=compilerVersion" json:"compiler_version,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CodeGeneratorRequest) Reset()                    { *m = CodeGeneratorRequest{} }
func (m *CodeGeneratorRequest) String() string            { return proto.CompactTextString(m) }
func (*CodeGeneratorRequest) ProtoMessage()               {}
func (*CodeGeneratorRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }
func (m *CodeGeneratorRequest) Unmarshal(b []byte) error {
	return xxx_messageInfo_CodeGeneratorRequest.Unmarshal(m, b)
}
func (m *CodeGeneratorRequest) Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CodeGeneratorRequest.Marshal(b, m, deterministic)
}
func (dst *CodeGeneratorRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CodeGeneratorRequest.Merge(dst, src)
}
func (m *CodeGeneratorRequest) XXX_Size() int {
	return xxx_messageInfo_CodeGeneratorRequest.Size(m)
}
func (m *CodeGeneratorRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CodeGeneratorRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CodeGeneratorRequest proto.InternalMessageInfo

func (m *CodeGeneratorRequest) GetFileToGenerate() []string {
	if m != nil {
		return m.FileToGenerate
	}
	return nil
}

func (m *CodeGeneratorRequest) GetParameter() string {
	if m != nil && m.Parameter != nil {
		return *m.Parameter
	}
	return ""
}

func (m *CodeGeneratorRequest) GetProtoFile() []*google_protobuf.FileDescriptorProto {
	if m != nil {
		return m.ProtoFile
	}
	return nil
}

func (m *CodeGeneratorRequest) GetCompilerVersion() *Version {
	if m != nil {
		return m.CompilerVersion
	}
	return nil
}


type CodeGeneratorResponse struct {
	
	
	
	
	
	
	
	
	Error                *string                       `protobuf:"bytes,1,opt,name=error" json:"error,omitempty"`
	File                 []*CodeGeneratorResponse_File `protobuf:"bytes,15,rep,name=file" json:"file,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                      `json:"-"`
	XXX_unrecognized     []byte                        `json:"-"`
	XXX_sizecache        int32                         `json:"-"`
}

func (m *CodeGeneratorResponse) Reset()                    { *m = CodeGeneratorResponse{} }
func (m *CodeGeneratorResponse) String() string            { return proto.CompactTextString(m) }
func (*CodeGeneratorResponse) ProtoMessage()               {}
func (*CodeGeneratorResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }
func (m *CodeGeneratorResponse) Unmarshal(b []byte) error {
	return xxx_messageInfo_CodeGeneratorResponse.Unmarshal(m, b)
}
func (m *CodeGeneratorResponse) Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CodeGeneratorResponse.Marshal(b, m, deterministic)
}
func (dst *CodeGeneratorResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CodeGeneratorResponse.Merge(dst, src)
}
func (m *CodeGeneratorResponse) XXX_Size() int {
	return xxx_messageInfo_CodeGeneratorResponse.Size(m)
}
func (m *CodeGeneratorResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_CodeGeneratorResponse.DiscardUnknown(m)
}

var xxx_messageInfo_CodeGeneratorResponse proto.InternalMessageInfo

func (m *CodeGeneratorResponse) GetError() string {
	if m != nil && m.Error != nil {
		return *m.Error
	}
	return ""
}

func (m *CodeGeneratorResponse) GetFile() []*CodeGeneratorResponse_File {
	if m != nil {
		return m.File
	}
	return nil
}


type CodeGeneratorResponse_File struct {
	
	
	
	
	
	
	
	
	
	
	
	Name *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	InsertionPoint *string `protobuf:"bytes,2,opt,name=insertion_point,json=insertionPoint" json:"insertion_point,omitempty"`
	
	Content              *string  `protobuf:"bytes,15,opt,name=content" json:"content,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CodeGeneratorResponse_File) Reset()                    { *m = CodeGeneratorResponse_File{} }
func (m *CodeGeneratorResponse_File) String() string            { return proto.CompactTextString(m) }
func (*CodeGeneratorResponse_File) ProtoMessage()               {}
func (*CodeGeneratorResponse_File) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2, 0} }
func (m *CodeGeneratorResponse_File) Unmarshal(b []byte) error {
	return xxx_messageInfo_CodeGeneratorResponse_File.Unmarshal(m, b)
}
func (m *CodeGeneratorResponse_File) Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CodeGeneratorResponse_File.Marshal(b, m, deterministic)
}
func (dst *CodeGeneratorResponse_File) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CodeGeneratorResponse_File.Merge(dst, src)
}
func (m *CodeGeneratorResponse_File) XXX_Size() int {
	return xxx_messageInfo_CodeGeneratorResponse_File.Size(m)
}
func (m *CodeGeneratorResponse_File) XXX_DiscardUnknown() {
	xxx_messageInfo_CodeGeneratorResponse_File.DiscardUnknown(m)
}

var xxx_messageInfo_CodeGeneratorResponse_File proto.InternalMessageInfo

func (m *CodeGeneratorResponse_File) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *CodeGeneratorResponse_File) GetInsertionPoint() string {
	if m != nil && m.InsertionPoint != nil {
		return *m.InsertionPoint
	}
	return ""
}

func (m *CodeGeneratorResponse_File) GetContent() string {
	if m != nil && m.Content != nil {
		return *m.Content
	}
	return ""
}

func init() {
	proto.RegisterType((*Version)(nil), "google.protobuf.compiler.Version")
	proto.RegisterType((*CodeGeneratorRequest)(nil), "google.protobuf.compiler.CodeGeneratorRequest")
	proto.RegisterType((*CodeGeneratorResponse)(nil), "google.protobuf.compiler.CodeGeneratorResponse")
	proto.RegisterType((*CodeGeneratorResponse_File)(nil), "google.protobuf.compiler.CodeGeneratorResponse.File")
}

func init() { proto.RegisterFile("google/protobuf/compiler/plugin.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x74, 0x92, 0xcf, 0x6a, 0x14, 0x41,
	0x10, 0xc6, 0x19, 0x77, 0x63, 0x98, 0x8a, 0x64, 0x43, 0x13, 0xa5, 0x09, 0x39, 0x8c, 0x8b, 0xe2,
	0x5c, 0x32, 0x0b, 0xc1, 0x8b, 0x78, 0x4b, 0x44, 0x3d, 0x78, 0x58, 0x1a, 0xf1, 0x20, 0xc8, 0x30,
	0x99, 0xd4, 0x74, 0x5a, 0x66, 0xba, 0xc6, 0xee, 0x1e, 0xf1, 0x49, 0x7d, 0x0f, 0xdf, 0x40, 0xfa,
	0xcf, 0x24, 0xb2, 0xb8, 0xa7, 0xee, 0xef, 0x57, 0xd5, 0xd5, 0x55, 0x1f, 0x05, 0x2f, 0x25, 0x91,
	0xec, 0x71, 0x33, 0x1a, 0x72, 0x74, 0x33, 0x75, 0x9b, 0x96, 0x86, 0x51, 0xf5, 0x68, 0x36, 0x63,
	0x3f, 0x49, 0xa5, 0xab, 0x10, 0x60, 0x3c, 0xa6, 0x55, 0x73, 0x5a, 0x35, 0xa7, 0x9d, 0x15, 0xbb,
	0x05, 0x6e, 0xd1, 0xb6, 0x46, 0x8d, 0x8e, 0x4c, 0xcc, 0x5e, 0xb7, 0x70, 0xf8, 0x05, 0x8d, 0x55,
	0xa4, 0xd9, 0x29, 0x1c, 0x0c, 0xcd, 0x77, 0x32, 0x3c, 0x2b, 0xb2, 0xf2, 0x40, 0x44, 0x11, 0xa8,
	0xd2, 0x64, 0xf8, 0xa3, 0x44, 0xbd, 0xf0, 0x74, 0x6c, 0x5c, 0x7b, 0xc7, 0x17, 0x91, 0x06, 0xc1,
	0x9e, 0xc1, 0x63, 0x3b, 0x75, 0x9d, 0xfa, 0xc5, 0x97, 0x45, 0x56, 0xe6, 0x22, 0xa9, 0xf5, 0x9f,
	0x0c, 0x4e, 0xaf, 0xe9, 0x16, 0x3f, 0xa0, 0x46, 0xd3, 0x38, 0x32, 0x02, 0x7f, 0x4c, 0x68, 0x1d,
	0x2b, 0xe1, 0xa4, 0x53, 0x3d, 0xd6, 0x8e, 0x6a, 0x19, 0x63, 0xc8, 0xb3, 0x62, 0x51, 0xe6, 0xe2,
	0xd8, 0xf3, 0xcf, 0x94, 0x5e, 0x20, 0x3b, 0x87, 0x7c, 0x6c, 0x4c, 0x33, 0xa0, 0xc3, 0xd8, 0x4a,
	0x2e, 0x1e, 0x00, 0xbb, 0x06, 0x08, 0xe3, 0xd4, 0xfe, 0x15, 0x5f, 0x15, 0x8b, 0xf2, 0xe8, 0xf2,
	0x45, 0xb5, 0x6b, 0xcb, 0x7b, 0xd5, 0xe3, 0xbb, 0x7b, 0x03, 0xb6, 0x1e, 0x8b, 0x3c, 0x44, 0x7d,
	0x84, 0x7d, 0x82, 0x93, 0xd9, 0xb8, 0xfa, 0x67, 0xf4, 0x24, 0x8c, 0x77, 0x74, 0xf9, 0xbc, 0xda,
	0xe7, 0x70, 0x95, 0xcc, 0x13, 0xab, 0x99, 0x24, 0xb0, 0xfe, 0x9d, 0xc1, 0xd3, 0x9d, 0x99, 0xed,
	0x48, 0xda, 0xa2, 0xf7, 0x0e, 0x8d, 0x49, 0x3e, 0xe7, 0x22, 0x0a, 0xf6, 0x11, 0x96, 0xff, 0x34,
	0xff, 0x7a, 0xff, 0x8f, 0xff, 0x2d, 0x1a, 0x66, 0x13, 0xa1, 0xc2, 0xd9, 0x37, 0x58, 0x86, 0x79,
	0x18, 0x2c, 0x75, 0x33, 0x60, 0xfa, 0x26, 0xdc, 0xd9, 0x2b, 0x58, 0x29, 0x6d, 0xd1, 0x38, 0x45,
	0xba, 0x1e, 0x49, 0x69, 0x97, 0xcc, 0x3c, 0xbe, 0xc7, 0x5b, 0x4f, 0x19, 0x87, 0xc3, 0x96, 0xb4,
	0x43, 0xed, 0xf8, 0x2a, 0x24, 0xcc, 0xf2, 0x4a, 0xc2, 0x79, 0x4b, 0xc3, 0xde, 0xfe, 0xae, 0x9e,
	0x6c, 0xc3, 0x6e, 0x06, 0x7b, 0xed, 0xd7, 0x37, 0x52, 0xb9, 0xbb, 0xe9, 0xc6, 0x87, 0x37, 0x92,
	0xfa, 0x46, 0xcb, 0x87, 0x65, 0x0c, 0x97, 0xf6, 0x42, 0xa2, 0xbe, 0x90, 0x94, 0x56, 0xfa, 0x6d,
	0x3c, 0x6a, 0x49, 0x7f, 0x03, 0x00, 0x00, 0xff, 0xff, 0xf7, 0x15, 0x40, 0xc5, 0xfe, 0x02, 0x00,
	0x00,
}
