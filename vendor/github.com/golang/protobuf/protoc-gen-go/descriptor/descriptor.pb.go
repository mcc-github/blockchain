


package descriptor 

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.ProtoPackageIsVersion2 

type FieldDescriptorProto_Type int32

const (
	
	
	FieldDescriptorProto_TYPE_DOUBLE FieldDescriptorProto_Type = 1
	FieldDescriptorProto_TYPE_FLOAT  FieldDescriptorProto_Type = 2
	
	
	FieldDescriptorProto_TYPE_INT64  FieldDescriptorProto_Type = 3
	FieldDescriptorProto_TYPE_UINT64 FieldDescriptorProto_Type = 4
	
	
	FieldDescriptorProto_TYPE_INT32   FieldDescriptorProto_Type = 5
	FieldDescriptorProto_TYPE_FIXED64 FieldDescriptorProto_Type = 6
	FieldDescriptorProto_TYPE_FIXED32 FieldDescriptorProto_Type = 7
	FieldDescriptorProto_TYPE_BOOL    FieldDescriptorProto_Type = 8
	FieldDescriptorProto_TYPE_STRING  FieldDescriptorProto_Type = 9
	
	
	
	
	FieldDescriptorProto_TYPE_GROUP   FieldDescriptorProto_Type = 10
	FieldDescriptorProto_TYPE_MESSAGE FieldDescriptorProto_Type = 11
	
	FieldDescriptorProto_TYPE_BYTES    FieldDescriptorProto_Type = 12
	FieldDescriptorProto_TYPE_UINT32   FieldDescriptorProto_Type = 13
	FieldDescriptorProto_TYPE_ENUM     FieldDescriptorProto_Type = 14
	FieldDescriptorProto_TYPE_SFIXED32 FieldDescriptorProto_Type = 15
	FieldDescriptorProto_TYPE_SFIXED64 FieldDescriptorProto_Type = 16
	FieldDescriptorProto_TYPE_SINT32   FieldDescriptorProto_Type = 17
	FieldDescriptorProto_TYPE_SINT64   FieldDescriptorProto_Type = 18
)

var FieldDescriptorProto_Type_name = map[int32]string{
	1:  "TYPE_DOUBLE",
	2:  "TYPE_FLOAT",
	3:  "TYPE_INT64",
	4:  "TYPE_UINT64",
	5:  "TYPE_INT32",
	6:  "TYPE_FIXED64",
	7:  "TYPE_FIXED32",
	8:  "TYPE_BOOL",
	9:  "TYPE_STRING",
	10: "TYPE_GROUP",
	11: "TYPE_MESSAGE",
	12: "TYPE_BYTES",
	13: "TYPE_UINT32",
	14: "TYPE_ENUM",
	15: "TYPE_SFIXED32",
	16: "TYPE_SFIXED64",
	17: "TYPE_SINT32",
	18: "TYPE_SINT64",
}
var FieldDescriptorProto_Type_value = map[string]int32{
	"TYPE_DOUBLE":   1,
	"TYPE_FLOAT":    2,
	"TYPE_INT64":    3,
	"TYPE_UINT64":   4,
	"TYPE_INT32":    5,
	"TYPE_FIXED64":  6,
	"TYPE_FIXED32":  7,
	"TYPE_BOOL":     8,
	"TYPE_STRING":   9,
	"TYPE_GROUP":    10,
	"TYPE_MESSAGE":  11,
	"TYPE_BYTES":    12,
	"TYPE_UINT32":   13,
	"TYPE_ENUM":     14,
	"TYPE_SFIXED32": 15,
	"TYPE_SFIXED64": 16,
	"TYPE_SINT32":   17,
	"TYPE_SINT64":   18,
}

func (x FieldDescriptorProto_Type) Enum() *FieldDescriptorProto_Type {
	p := new(FieldDescriptorProto_Type)
	*p = x
	return p
}
func (x FieldDescriptorProto_Type) String() string {
	return proto.EnumName(FieldDescriptorProto_Type_name, int32(x))
}
func (x *FieldDescriptorProto_Type) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(FieldDescriptorProto_Type_value, data, "FieldDescriptorProto_Type")
	if err != nil {
		return err
	}
	*x = FieldDescriptorProto_Type(value)
	return nil
}
func (FieldDescriptorProto_Type) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{4, 0}
}

type FieldDescriptorProto_Label int32

const (
	
	FieldDescriptorProto_LABEL_OPTIONAL FieldDescriptorProto_Label = 1
	FieldDescriptorProto_LABEL_REQUIRED FieldDescriptorProto_Label = 2
	FieldDescriptorProto_LABEL_REPEATED FieldDescriptorProto_Label = 3
)

var FieldDescriptorProto_Label_name = map[int32]string{
	1: "LABEL_OPTIONAL",
	2: "LABEL_REQUIRED",
	3: "LABEL_REPEATED",
}
var FieldDescriptorProto_Label_value = map[string]int32{
	"LABEL_OPTIONAL": 1,
	"LABEL_REQUIRED": 2,
	"LABEL_REPEATED": 3,
}

func (x FieldDescriptorProto_Label) Enum() *FieldDescriptorProto_Label {
	p := new(FieldDescriptorProto_Label)
	*p = x
	return p
}
func (x FieldDescriptorProto_Label) String() string {
	return proto.EnumName(FieldDescriptorProto_Label_name, int32(x))
}
func (x *FieldDescriptorProto_Label) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(FieldDescriptorProto_Label_value, data, "FieldDescriptorProto_Label")
	if err != nil {
		return err
	}
	*x = FieldDescriptorProto_Label(value)
	return nil
}
func (FieldDescriptorProto_Label) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{4, 1}
}


type FileOptions_OptimizeMode int32

const (
	FileOptions_SPEED FileOptions_OptimizeMode = 1
	
	FileOptions_CODE_SIZE    FileOptions_OptimizeMode = 2
	FileOptions_LITE_RUNTIME FileOptions_OptimizeMode = 3
)

var FileOptions_OptimizeMode_name = map[int32]string{
	1: "SPEED",
	2: "CODE_SIZE",
	3: "LITE_RUNTIME",
}
var FileOptions_OptimizeMode_value = map[string]int32{
	"SPEED":        1,
	"CODE_SIZE":    2,
	"LITE_RUNTIME": 3,
}

func (x FileOptions_OptimizeMode) Enum() *FileOptions_OptimizeMode {
	p := new(FileOptions_OptimizeMode)
	*p = x
	return p
}
func (x FileOptions_OptimizeMode) String() string {
	return proto.EnumName(FileOptions_OptimizeMode_name, int32(x))
}
func (x *FileOptions_OptimizeMode) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(FileOptions_OptimizeMode_value, data, "FileOptions_OptimizeMode")
	if err != nil {
		return err
	}
	*x = FileOptions_OptimizeMode(value)
	return nil
}
func (FileOptions_OptimizeMode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{10, 0}
}

type FieldOptions_CType int32

const (
	
	FieldOptions_STRING       FieldOptions_CType = 0
	FieldOptions_CORD         FieldOptions_CType = 1
	FieldOptions_STRING_PIECE FieldOptions_CType = 2
)

var FieldOptions_CType_name = map[int32]string{
	0: "STRING",
	1: "CORD",
	2: "STRING_PIECE",
}
var FieldOptions_CType_value = map[string]int32{
	"STRING":       0,
	"CORD":         1,
	"STRING_PIECE": 2,
}

func (x FieldOptions_CType) Enum() *FieldOptions_CType {
	p := new(FieldOptions_CType)
	*p = x
	return p
}
func (x FieldOptions_CType) String() string {
	return proto.EnumName(FieldOptions_CType_name, int32(x))
}
func (x *FieldOptions_CType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(FieldOptions_CType_value, data, "FieldOptions_CType")
	if err != nil {
		return err
	}
	*x = FieldOptions_CType(value)
	return nil
}
func (FieldOptions_CType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{12, 0}
}

type FieldOptions_JSType int32

const (
	
	FieldOptions_JS_NORMAL FieldOptions_JSType = 0
	
	FieldOptions_JS_STRING FieldOptions_JSType = 1
	
	FieldOptions_JS_NUMBER FieldOptions_JSType = 2
)

var FieldOptions_JSType_name = map[int32]string{
	0: "JS_NORMAL",
	1: "JS_STRING",
	2: "JS_NUMBER",
}
var FieldOptions_JSType_value = map[string]int32{
	"JS_NORMAL": 0,
	"JS_STRING": 1,
	"JS_NUMBER": 2,
}

func (x FieldOptions_JSType) Enum() *FieldOptions_JSType {
	p := new(FieldOptions_JSType)
	*p = x
	return p
}
func (x FieldOptions_JSType) String() string {
	return proto.EnumName(FieldOptions_JSType_name, int32(x))
}
func (x *FieldOptions_JSType) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(FieldOptions_JSType_value, data, "FieldOptions_JSType")
	if err != nil {
		return err
	}
	*x = FieldOptions_JSType(value)
	return nil
}
func (FieldOptions_JSType) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{12, 1}
}




type MethodOptions_IdempotencyLevel int32

const (
	MethodOptions_IDEMPOTENCY_UNKNOWN MethodOptions_IdempotencyLevel = 0
	MethodOptions_NO_SIDE_EFFECTS     MethodOptions_IdempotencyLevel = 1
	MethodOptions_IDEMPOTENT          MethodOptions_IdempotencyLevel = 2
)

var MethodOptions_IdempotencyLevel_name = map[int32]string{
	0: "IDEMPOTENCY_UNKNOWN",
	1: "NO_SIDE_EFFECTS",
	2: "IDEMPOTENT",
}
var MethodOptions_IdempotencyLevel_value = map[string]int32{
	"IDEMPOTENCY_UNKNOWN": 0,
	"NO_SIDE_EFFECTS":     1,
	"IDEMPOTENT":          2,
}

func (x MethodOptions_IdempotencyLevel) Enum() *MethodOptions_IdempotencyLevel {
	p := new(MethodOptions_IdempotencyLevel)
	*p = x
	return p
}
func (x MethodOptions_IdempotencyLevel) String() string {
	return proto.EnumName(MethodOptions_IdempotencyLevel_name, int32(x))
}
func (x *MethodOptions_IdempotencyLevel) UnmarshalJSON(data []byte) error {
	value, err := proto.UnmarshalJSONEnum(MethodOptions_IdempotencyLevel_value, data, "MethodOptions_IdempotencyLevel")
	if err != nil {
		return err
	}
	*x = MethodOptions_IdempotencyLevel(value)
	return nil
}
func (MethodOptions_IdempotencyLevel) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{17, 0}
}



type FileDescriptorSet struct {
	File                 []*FileDescriptorProto `protobuf:"bytes,1,rep,name=file" json:"file,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *FileDescriptorSet) Reset()         { *m = FileDescriptorSet{} }
func (m *FileDescriptorSet) String() string { return proto.CompactTextString(m) }
func (*FileDescriptorSet) ProtoMessage()    {}
func (*FileDescriptorSet) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{0}
}
func (m *FileDescriptorSet) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileDescriptorSet.Unmarshal(m, b)
}
func (m *FileDescriptorSet) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileDescriptorSet.Marshal(b, m, deterministic)
}
func (dst *FileDescriptorSet) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileDescriptorSet.Merge(dst, src)
}
func (m *FileDescriptorSet) XXX_Size() int {
	return xxx_messageInfo_FileDescriptorSet.Size(m)
}
func (m *FileDescriptorSet) XXX_DiscardUnknown() {
	xxx_messageInfo_FileDescriptorSet.DiscardUnknown(m)
}

var xxx_messageInfo_FileDescriptorSet proto.InternalMessageInfo

func (m *FileDescriptorSet) GetFile() []*FileDescriptorProto {
	if m != nil {
		return m.File
	}
	return nil
}


type FileDescriptorProto struct {
	Name    *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Package *string `protobuf:"bytes,2,opt,name=package" json:"package,omitempty"`
	
	Dependency []string `protobuf:"bytes,3,rep,name=dependency" json:"dependency,omitempty"`
	
	PublicDependency []int32 `protobuf:"varint,10,rep,name=public_dependency,json=publicDependency" json:"public_dependency,omitempty"`
	
	
	WeakDependency []int32 `protobuf:"varint,11,rep,name=weak_dependency,json=weakDependency" json:"weak_dependency,omitempty"`
	
	MessageType []*DescriptorProto        `protobuf:"bytes,4,rep,name=message_type,json=messageType" json:"message_type,omitempty"`
	EnumType    []*EnumDescriptorProto    `protobuf:"bytes,5,rep,name=enum_type,json=enumType" json:"enum_type,omitempty"`
	Service     []*ServiceDescriptorProto `protobuf:"bytes,6,rep,name=service" json:"service,omitempty"`
	Extension   []*FieldDescriptorProto   `protobuf:"bytes,7,rep,name=extension" json:"extension,omitempty"`
	Options     *FileOptions              `protobuf:"bytes,8,opt,name=options" json:"options,omitempty"`
	
	
	
	
	SourceCodeInfo *SourceCodeInfo `protobuf:"bytes,9,opt,name=source_code_info,json=sourceCodeInfo" json:"source_code_info,omitempty"`
	
	
	Syntax               *string  `protobuf:"bytes,12,opt,name=syntax" json:"syntax,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FileDescriptorProto) Reset()         { *m = FileDescriptorProto{} }
func (m *FileDescriptorProto) String() string { return proto.CompactTextString(m) }
func (*FileDescriptorProto) ProtoMessage()    {}
func (*FileDescriptorProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{1}
}
func (m *FileDescriptorProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileDescriptorProto.Unmarshal(m, b)
}
func (m *FileDescriptorProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileDescriptorProto.Marshal(b, m, deterministic)
}
func (dst *FileDescriptorProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileDescriptorProto.Merge(dst, src)
}
func (m *FileDescriptorProto) XXX_Size() int {
	return xxx_messageInfo_FileDescriptorProto.Size(m)
}
func (m *FileDescriptorProto) XXX_DiscardUnknown() {
	xxx_messageInfo_FileDescriptorProto.DiscardUnknown(m)
}

var xxx_messageInfo_FileDescriptorProto proto.InternalMessageInfo

func (m *FileDescriptorProto) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *FileDescriptorProto) GetPackage() string {
	if m != nil && m.Package != nil {
		return *m.Package
	}
	return ""
}

func (m *FileDescriptorProto) GetDependency() []string {
	if m != nil {
		return m.Dependency
	}
	return nil
}

func (m *FileDescriptorProto) GetPublicDependency() []int32 {
	if m != nil {
		return m.PublicDependency
	}
	return nil
}

func (m *FileDescriptorProto) GetWeakDependency() []int32 {
	if m != nil {
		return m.WeakDependency
	}
	return nil
}

func (m *FileDescriptorProto) GetMessageType() []*DescriptorProto {
	if m != nil {
		return m.MessageType
	}
	return nil
}

func (m *FileDescriptorProto) GetEnumType() []*EnumDescriptorProto {
	if m != nil {
		return m.EnumType
	}
	return nil
}

func (m *FileDescriptorProto) GetService() []*ServiceDescriptorProto {
	if m != nil {
		return m.Service
	}
	return nil
}

func (m *FileDescriptorProto) GetExtension() []*FieldDescriptorProto {
	if m != nil {
		return m.Extension
	}
	return nil
}

func (m *FileDescriptorProto) GetOptions() *FileOptions {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *FileDescriptorProto) GetSourceCodeInfo() *SourceCodeInfo {
	if m != nil {
		return m.SourceCodeInfo
	}
	return nil
}

func (m *FileDescriptorProto) GetSyntax() string {
	if m != nil && m.Syntax != nil {
		return *m.Syntax
	}
	return ""
}


type DescriptorProto struct {
	Name           *string                           `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Field          []*FieldDescriptorProto           `protobuf:"bytes,2,rep,name=field" json:"field,omitempty"`
	Extension      []*FieldDescriptorProto           `protobuf:"bytes,6,rep,name=extension" json:"extension,omitempty"`
	NestedType     []*DescriptorProto                `protobuf:"bytes,3,rep,name=nested_type,json=nestedType" json:"nested_type,omitempty"`
	EnumType       []*EnumDescriptorProto            `protobuf:"bytes,4,rep,name=enum_type,json=enumType" json:"enum_type,omitempty"`
	ExtensionRange []*DescriptorProto_ExtensionRange `protobuf:"bytes,5,rep,name=extension_range,json=extensionRange" json:"extension_range,omitempty"`
	OneofDecl      []*OneofDescriptorProto           `protobuf:"bytes,8,rep,name=oneof_decl,json=oneofDecl" json:"oneof_decl,omitempty"`
	Options        *MessageOptions                   `protobuf:"bytes,7,opt,name=options" json:"options,omitempty"`
	ReservedRange  []*DescriptorProto_ReservedRange  `protobuf:"bytes,9,rep,name=reserved_range,json=reservedRange" json:"reserved_range,omitempty"`
	
	
	ReservedName         []string `protobuf:"bytes,10,rep,name=reserved_name,json=reservedName" json:"reserved_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DescriptorProto) Reset()         { *m = DescriptorProto{} }
func (m *DescriptorProto) String() string { return proto.CompactTextString(m) }
func (*DescriptorProto) ProtoMessage()    {}
func (*DescriptorProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{2}
}
func (m *DescriptorProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DescriptorProto.Unmarshal(m, b)
}
func (m *DescriptorProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DescriptorProto.Marshal(b, m, deterministic)
}
func (dst *DescriptorProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DescriptorProto.Merge(dst, src)
}
func (m *DescriptorProto) XXX_Size() int {
	return xxx_messageInfo_DescriptorProto.Size(m)
}
func (m *DescriptorProto) XXX_DiscardUnknown() {
	xxx_messageInfo_DescriptorProto.DiscardUnknown(m)
}

var xxx_messageInfo_DescriptorProto proto.InternalMessageInfo

func (m *DescriptorProto) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *DescriptorProto) GetField() []*FieldDescriptorProto {
	if m != nil {
		return m.Field
	}
	return nil
}

func (m *DescriptorProto) GetExtension() []*FieldDescriptorProto {
	if m != nil {
		return m.Extension
	}
	return nil
}

func (m *DescriptorProto) GetNestedType() []*DescriptorProto {
	if m != nil {
		return m.NestedType
	}
	return nil
}

func (m *DescriptorProto) GetEnumType() []*EnumDescriptorProto {
	if m != nil {
		return m.EnumType
	}
	return nil
}

func (m *DescriptorProto) GetExtensionRange() []*DescriptorProto_ExtensionRange {
	if m != nil {
		return m.ExtensionRange
	}
	return nil
}

func (m *DescriptorProto) GetOneofDecl() []*OneofDescriptorProto {
	if m != nil {
		return m.OneofDecl
	}
	return nil
}

func (m *DescriptorProto) GetOptions() *MessageOptions {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *DescriptorProto) GetReservedRange() []*DescriptorProto_ReservedRange {
	if m != nil {
		return m.ReservedRange
	}
	return nil
}

func (m *DescriptorProto) GetReservedName() []string {
	if m != nil {
		return m.ReservedName
	}
	return nil
}

type DescriptorProto_ExtensionRange struct {
	Start                *int32                 `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	End                  *int32                 `protobuf:"varint,2,opt,name=end" json:"end,omitempty"`
	Options              *ExtensionRangeOptions `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}               `json:"-"`
	XXX_unrecognized     []byte                 `json:"-"`
	XXX_sizecache        int32                  `json:"-"`
}

func (m *DescriptorProto_ExtensionRange) Reset()         { *m = DescriptorProto_ExtensionRange{} }
func (m *DescriptorProto_ExtensionRange) String() string { return proto.CompactTextString(m) }
func (*DescriptorProto_ExtensionRange) ProtoMessage()    {}
func (*DescriptorProto_ExtensionRange) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{2, 0}
}
func (m *DescriptorProto_ExtensionRange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DescriptorProto_ExtensionRange.Unmarshal(m, b)
}
func (m *DescriptorProto_ExtensionRange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DescriptorProto_ExtensionRange.Marshal(b, m, deterministic)
}
func (dst *DescriptorProto_ExtensionRange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DescriptorProto_ExtensionRange.Merge(dst, src)
}
func (m *DescriptorProto_ExtensionRange) XXX_Size() int {
	return xxx_messageInfo_DescriptorProto_ExtensionRange.Size(m)
}
func (m *DescriptorProto_ExtensionRange) XXX_DiscardUnknown() {
	xxx_messageInfo_DescriptorProto_ExtensionRange.DiscardUnknown(m)
}

var xxx_messageInfo_DescriptorProto_ExtensionRange proto.InternalMessageInfo

func (m *DescriptorProto_ExtensionRange) GetStart() int32 {
	if m != nil && m.Start != nil {
		return *m.Start
	}
	return 0
}

func (m *DescriptorProto_ExtensionRange) GetEnd() int32 {
	if m != nil && m.End != nil {
		return *m.End
	}
	return 0
}

func (m *DescriptorProto_ExtensionRange) GetOptions() *ExtensionRangeOptions {
	if m != nil {
		return m.Options
	}
	return nil
}




type DescriptorProto_ReservedRange struct {
	Start                *int32   `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	End                  *int32   `protobuf:"varint,2,opt,name=end" json:"end,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DescriptorProto_ReservedRange) Reset()         { *m = DescriptorProto_ReservedRange{} }
func (m *DescriptorProto_ReservedRange) String() string { return proto.CompactTextString(m) }
func (*DescriptorProto_ReservedRange) ProtoMessage()    {}
func (*DescriptorProto_ReservedRange) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{2, 1}
}
func (m *DescriptorProto_ReservedRange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DescriptorProto_ReservedRange.Unmarshal(m, b)
}
func (m *DescriptorProto_ReservedRange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DescriptorProto_ReservedRange.Marshal(b, m, deterministic)
}
func (dst *DescriptorProto_ReservedRange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DescriptorProto_ReservedRange.Merge(dst, src)
}
func (m *DescriptorProto_ReservedRange) XXX_Size() int {
	return xxx_messageInfo_DescriptorProto_ReservedRange.Size(m)
}
func (m *DescriptorProto_ReservedRange) XXX_DiscardUnknown() {
	xxx_messageInfo_DescriptorProto_ReservedRange.DiscardUnknown(m)
}

var xxx_messageInfo_DescriptorProto_ReservedRange proto.InternalMessageInfo

func (m *DescriptorProto_ReservedRange) GetStart() int32 {
	if m != nil && m.Start != nil {
		return *m.Start
	}
	return 0
}

func (m *DescriptorProto_ReservedRange) GetEnd() int32 {
	if m != nil && m.End != nil {
		return *m.End
	}
	return 0
}

type ExtensionRangeOptions struct {
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}               `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *ExtensionRangeOptions) Reset()         { *m = ExtensionRangeOptions{} }
func (m *ExtensionRangeOptions) String() string { return proto.CompactTextString(m) }
func (*ExtensionRangeOptions) ProtoMessage()    {}
func (*ExtensionRangeOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{3}
}

var extRange_ExtensionRangeOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*ExtensionRangeOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_ExtensionRangeOptions
}
func (m *ExtensionRangeOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ExtensionRangeOptions.Unmarshal(m, b)
}
func (m *ExtensionRangeOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ExtensionRangeOptions.Marshal(b, m, deterministic)
}
func (dst *ExtensionRangeOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ExtensionRangeOptions.Merge(dst, src)
}
func (m *ExtensionRangeOptions) XXX_Size() int {
	return xxx_messageInfo_ExtensionRangeOptions.Size(m)
}
func (m *ExtensionRangeOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_ExtensionRangeOptions.DiscardUnknown(m)
}

var xxx_messageInfo_ExtensionRangeOptions proto.InternalMessageInfo

func (m *ExtensionRangeOptions) GetUninterpretedOption() []*UninterpretedOption {
	if m != nil {
		return m.UninterpretedOption
	}
	return nil
}


type FieldDescriptorProto struct {
	Name   *string                     `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Number *int32                      `protobuf:"varint,3,opt,name=number" json:"number,omitempty"`
	Label  *FieldDescriptorProto_Label `protobuf:"varint,4,opt,name=label,enum=google.protobuf.FieldDescriptorProto_Label" json:"label,omitempty"`
	
	
	Type *FieldDescriptorProto_Type `protobuf:"varint,5,opt,name=type,enum=google.protobuf.FieldDescriptorProto_Type" json:"type,omitempty"`
	
	
	
	
	
	TypeName *string `protobuf:"bytes,6,opt,name=type_name,json=typeName" json:"type_name,omitempty"`
	
	
	Extendee *string `protobuf:"bytes,2,opt,name=extendee" json:"extendee,omitempty"`
	
	
	
	
	
	DefaultValue *string `protobuf:"bytes,7,opt,name=default_value,json=defaultValue" json:"default_value,omitempty"`
	
	
	OneofIndex *int32 `protobuf:"varint,9,opt,name=oneof_index,json=oneofIndex" json:"oneof_index,omitempty"`
	
	
	
	
	JsonName             *string       `protobuf:"bytes,10,opt,name=json_name,json=jsonName" json:"json_name,omitempty"`
	Options              *FieldOptions `protobuf:"bytes,8,opt,name=options" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *FieldDescriptorProto) Reset()         { *m = FieldDescriptorProto{} }
func (m *FieldDescriptorProto) String() string { return proto.CompactTextString(m) }
func (*FieldDescriptorProto) ProtoMessage()    {}
func (*FieldDescriptorProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{4}
}
func (m *FieldDescriptorProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FieldDescriptorProto.Unmarshal(m, b)
}
func (m *FieldDescriptorProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FieldDescriptorProto.Marshal(b, m, deterministic)
}
func (dst *FieldDescriptorProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FieldDescriptorProto.Merge(dst, src)
}
func (m *FieldDescriptorProto) XXX_Size() int {
	return xxx_messageInfo_FieldDescriptorProto.Size(m)
}
func (m *FieldDescriptorProto) XXX_DiscardUnknown() {
	xxx_messageInfo_FieldDescriptorProto.DiscardUnknown(m)
}

var xxx_messageInfo_FieldDescriptorProto proto.InternalMessageInfo

func (m *FieldDescriptorProto) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *FieldDescriptorProto) GetNumber() int32 {
	if m != nil && m.Number != nil {
		return *m.Number
	}
	return 0
}

func (m *FieldDescriptorProto) GetLabel() FieldDescriptorProto_Label {
	if m != nil && m.Label != nil {
		return *m.Label
	}
	return FieldDescriptorProto_LABEL_OPTIONAL
}

func (m *FieldDescriptorProto) GetType() FieldDescriptorProto_Type {
	if m != nil && m.Type != nil {
		return *m.Type
	}
	return FieldDescriptorProto_TYPE_DOUBLE
}

func (m *FieldDescriptorProto) GetTypeName() string {
	if m != nil && m.TypeName != nil {
		return *m.TypeName
	}
	return ""
}

func (m *FieldDescriptorProto) GetExtendee() string {
	if m != nil && m.Extendee != nil {
		return *m.Extendee
	}
	return ""
}

func (m *FieldDescriptorProto) GetDefaultValue() string {
	if m != nil && m.DefaultValue != nil {
		return *m.DefaultValue
	}
	return ""
}

func (m *FieldDescriptorProto) GetOneofIndex() int32 {
	if m != nil && m.OneofIndex != nil {
		return *m.OneofIndex
	}
	return 0
}

func (m *FieldDescriptorProto) GetJsonName() string {
	if m != nil && m.JsonName != nil {
		return *m.JsonName
	}
	return ""
}

func (m *FieldDescriptorProto) GetOptions() *FieldOptions {
	if m != nil {
		return m.Options
	}
	return nil
}


type OneofDescriptorProto struct {
	Name                 *string       `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Options              *OneofOptions `protobuf:"bytes,2,opt,name=options" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *OneofDescriptorProto) Reset()         { *m = OneofDescriptorProto{} }
func (m *OneofDescriptorProto) String() string { return proto.CompactTextString(m) }
func (*OneofDescriptorProto) ProtoMessage()    {}
func (*OneofDescriptorProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{5}
}
func (m *OneofDescriptorProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OneofDescriptorProto.Unmarshal(m, b)
}
func (m *OneofDescriptorProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OneofDescriptorProto.Marshal(b, m, deterministic)
}
func (dst *OneofDescriptorProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OneofDescriptorProto.Merge(dst, src)
}
func (m *OneofDescriptorProto) XXX_Size() int {
	return xxx_messageInfo_OneofDescriptorProto.Size(m)
}
func (m *OneofDescriptorProto) XXX_DiscardUnknown() {
	xxx_messageInfo_OneofDescriptorProto.DiscardUnknown(m)
}

var xxx_messageInfo_OneofDescriptorProto proto.InternalMessageInfo

func (m *OneofDescriptorProto) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *OneofDescriptorProto) GetOptions() *OneofOptions {
	if m != nil {
		return m.Options
	}
	return nil
}


type EnumDescriptorProto struct {
	Name    *string                     `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value   []*EnumValueDescriptorProto `protobuf:"bytes,2,rep,name=value" json:"value,omitempty"`
	Options *EnumOptions                `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
	
	
	
	ReservedRange []*EnumDescriptorProto_EnumReservedRange `protobuf:"bytes,4,rep,name=reserved_range,json=reservedRange" json:"reserved_range,omitempty"`
	
	
	ReservedName         []string `protobuf:"bytes,5,rep,name=reserved_name,json=reservedName" json:"reserved_name,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnumDescriptorProto) Reset()         { *m = EnumDescriptorProto{} }
func (m *EnumDescriptorProto) String() string { return proto.CompactTextString(m) }
func (*EnumDescriptorProto) ProtoMessage()    {}
func (*EnumDescriptorProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{6}
}
func (m *EnumDescriptorProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnumDescriptorProto.Unmarshal(m, b)
}
func (m *EnumDescriptorProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnumDescriptorProto.Marshal(b, m, deterministic)
}
func (dst *EnumDescriptorProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnumDescriptorProto.Merge(dst, src)
}
func (m *EnumDescriptorProto) XXX_Size() int {
	return xxx_messageInfo_EnumDescriptorProto.Size(m)
}
func (m *EnumDescriptorProto) XXX_DiscardUnknown() {
	xxx_messageInfo_EnumDescriptorProto.DiscardUnknown(m)
}

var xxx_messageInfo_EnumDescriptorProto proto.InternalMessageInfo

func (m *EnumDescriptorProto) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *EnumDescriptorProto) GetValue() []*EnumValueDescriptorProto {
	if m != nil {
		return m.Value
	}
	return nil
}

func (m *EnumDescriptorProto) GetOptions() *EnumOptions {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *EnumDescriptorProto) GetReservedRange() []*EnumDescriptorProto_EnumReservedRange {
	if m != nil {
		return m.ReservedRange
	}
	return nil
}

func (m *EnumDescriptorProto) GetReservedName() []string {
	if m != nil {
		return m.ReservedName
	}
	return nil
}







type EnumDescriptorProto_EnumReservedRange struct {
	Start                *int32   `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	End                  *int32   `protobuf:"varint,2,opt,name=end" json:"end,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *EnumDescriptorProto_EnumReservedRange) Reset()         { *m = EnumDescriptorProto_EnumReservedRange{} }
func (m *EnumDescriptorProto_EnumReservedRange) String() string { return proto.CompactTextString(m) }
func (*EnumDescriptorProto_EnumReservedRange) ProtoMessage()    {}
func (*EnumDescriptorProto_EnumReservedRange) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{6, 0}
}
func (m *EnumDescriptorProto_EnumReservedRange) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnumDescriptorProto_EnumReservedRange.Unmarshal(m, b)
}
func (m *EnumDescriptorProto_EnumReservedRange) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnumDescriptorProto_EnumReservedRange.Marshal(b, m, deterministic)
}
func (dst *EnumDescriptorProto_EnumReservedRange) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnumDescriptorProto_EnumReservedRange.Merge(dst, src)
}
func (m *EnumDescriptorProto_EnumReservedRange) XXX_Size() int {
	return xxx_messageInfo_EnumDescriptorProto_EnumReservedRange.Size(m)
}
func (m *EnumDescriptorProto_EnumReservedRange) XXX_DiscardUnknown() {
	xxx_messageInfo_EnumDescriptorProto_EnumReservedRange.DiscardUnknown(m)
}

var xxx_messageInfo_EnumDescriptorProto_EnumReservedRange proto.InternalMessageInfo

func (m *EnumDescriptorProto_EnumReservedRange) GetStart() int32 {
	if m != nil && m.Start != nil {
		return *m.Start
	}
	return 0
}

func (m *EnumDescriptorProto_EnumReservedRange) GetEnd() int32 {
	if m != nil && m.End != nil {
		return *m.End
	}
	return 0
}


type EnumValueDescriptorProto struct {
	Name                 *string           `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Number               *int32            `protobuf:"varint,2,opt,name=number" json:"number,omitempty"`
	Options              *EnumValueOptions `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}          `json:"-"`
	XXX_unrecognized     []byte            `json:"-"`
	XXX_sizecache        int32             `json:"-"`
}

func (m *EnumValueDescriptorProto) Reset()         { *m = EnumValueDescriptorProto{} }
func (m *EnumValueDescriptorProto) String() string { return proto.CompactTextString(m) }
func (*EnumValueDescriptorProto) ProtoMessage()    {}
func (*EnumValueDescriptorProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{7}
}
func (m *EnumValueDescriptorProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnumValueDescriptorProto.Unmarshal(m, b)
}
func (m *EnumValueDescriptorProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnumValueDescriptorProto.Marshal(b, m, deterministic)
}
func (dst *EnumValueDescriptorProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnumValueDescriptorProto.Merge(dst, src)
}
func (m *EnumValueDescriptorProto) XXX_Size() int {
	return xxx_messageInfo_EnumValueDescriptorProto.Size(m)
}
func (m *EnumValueDescriptorProto) XXX_DiscardUnknown() {
	xxx_messageInfo_EnumValueDescriptorProto.DiscardUnknown(m)
}

var xxx_messageInfo_EnumValueDescriptorProto proto.InternalMessageInfo

func (m *EnumValueDescriptorProto) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *EnumValueDescriptorProto) GetNumber() int32 {
	if m != nil && m.Number != nil {
		return *m.Number
	}
	return 0
}

func (m *EnumValueDescriptorProto) GetOptions() *EnumValueOptions {
	if m != nil {
		return m.Options
	}
	return nil
}


type ServiceDescriptorProto struct {
	Name                 *string                  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Method               []*MethodDescriptorProto `protobuf:"bytes,2,rep,name=method" json:"method,omitempty"`
	Options              *ServiceOptions          `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                 `json:"-"`
	XXX_unrecognized     []byte                   `json:"-"`
	XXX_sizecache        int32                    `json:"-"`
}

func (m *ServiceDescriptorProto) Reset()         { *m = ServiceDescriptorProto{} }
func (m *ServiceDescriptorProto) String() string { return proto.CompactTextString(m) }
func (*ServiceDescriptorProto) ProtoMessage()    {}
func (*ServiceDescriptorProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{8}
}
func (m *ServiceDescriptorProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceDescriptorProto.Unmarshal(m, b)
}
func (m *ServiceDescriptorProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceDescriptorProto.Marshal(b, m, deterministic)
}
func (dst *ServiceDescriptorProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceDescriptorProto.Merge(dst, src)
}
func (m *ServiceDescriptorProto) XXX_Size() int {
	return xxx_messageInfo_ServiceDescriptorProto.Size(m)
}
func (m *ServiceDescriptorProto) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceDescriptorProto.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceDescriptorProto proto.InternalMessageInfo

func (m *ServiceDescriptorProto) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *ServiceDescriptorProto) GetMethod() []*MethodDescriptorProto {
	if m != nil {
		return m.Method
	}
	return nil
}

func (m *ServiceDescriptorProto) GetOptions() *ServiceOptions {
	if m != nil {
		return m.Options
	}
	return nil
}


type MethodDescriptorProto struct {
	Name *string `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	
	
	InputType  *string        `protobuf:"bytes,2,opt,name=input_type,json=inputType" json:"input_type,omitempty"`
	OutputType *string        `protobuf:"bytes,3,opt,name=output_type,json=outputType" json:"output_type,omitempty"`
	Options    *MethodOptions `protobuf:"bytes,4,opt,name=options" json:"options,omitempty"`
	
	ClientStreaming *bool `protobuf:"varint,5,opt,name=client_streaming,json=clientStreaming,def=0" json:"client_streaming,omitempty"`
	
	ServerStreaming      *bool    `protobuf:"varint,6,opt,name=server_streaming,json=serverStreaming,def=0" json:"server_streaming,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *MethodDescriptorProto) Reset()         { *m = MethodDescriptorProto{} }
func (m *MethodDescriptorProto) String() string { return proto.CompactTextString(m) }
func (*MethodDescriptorProto) ProtoMessage()    {}
func (*MethodDescriptorProto) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{9}
}
func (m *MethodDescriptorProto) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MethodDescriptorProto.Unmarshal(m, b)
}
func (m *MethodDescriptorProto) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MethodDescriptorProto.Marshal(b, m, deterministic)
}
func (dst *MethodDescriptorProto) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MethodDescriptorProto.Merge(dst, src)
}
func (m *MethodDescriptorProto) XXX_Size() int {
	return xxx_messageInfo_MethodDescriptorProto.Size(m)
}
func (m *MethodDescriptorProto) XXX_DiscardUnknown() {
	xxx_messageInfo_MethodDescriptorProto.DiscardUnknown(m)
}

var xxx_messageInfo_MethodDescriptorProto proto.InternalMessageInfo

const Default_MethodDescriptorProto_ClientStreaming bool = false
const Default_MethodDescriptorProto_ServerStreaming bool = false

func (m *MethodDescriptorProto) GetName() string {
	if m != nil && m.Name != nil {
		return *m.Name
	}
	return ""
}

func (m *MethodDescriptorProto) GetInputType() string {
	if m != nil && m.InputType != nil {
		return *m.InputType
	}
	return ""
}

func (m *MethodDescriptorProto) GetOutputType() string {
	if m != nil && m.OutputType != nil {
		return *m.OutputType
	}
	return ""
}

func (m *MethodDescriptorProto) GetOptions() *MethodOptions {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *MethodDescriptorProto) GetClientStreaming() bool {
	if m != nil && m.ClientStreaming != nil {
		return *m.ClientStreaming
	}
	return Default_MethodDescriptorProto_ClientStreaming
}

func (m *MethodDescriptorProto) GetServerStreaming() bool {
	if m != nil && m.ServerStreaming != nil {
		return *m.ServerStreaming
	}
	return Default_MethodDescriptorProto_ServerStreaming
}

type FileOptions struct {
	
	
	
	
	JavaPackage *string `protobuf:"bytes,1,opt,name=java_package,json=javaPackage" json:"java_package,omitempty"`
	
	
	
	
	
	JavaOuterClassname *string `protobuf:"bytes,8,opt,name=java_outer_classname,json=javaOuterClassname" json:"java_outer_classname,omitempty"`
	
	
	
	
	
	
	JavaMultipleFiles *bool `protobuf:"varint,10,opt,name=java_multiple_files,json=javaMultipleFiles,def=0" json:"java_multiple_files,omitempty"`
	
	JavaGenerateEqualsAndHash *bool `protobuf:"varint,20,opt,name=java_generate_equals_and_hash,json=javaGenerateEqualsAndHash" json:"java_generate_equals_and_hash,omitempty"` 
	
	
	
	
	
	
	JavaStringCheckUtf8 *bool                     `protobuf:"varint,27,opt,name=java_string_check_utf8,json=javaStringCheckUtf8,def=0" json:"java_string_check_utf8,omitempty"`
	OptimizeFor         *FileOptions_OptimizeMode `protobuf:"varint,9,opt,name=optimize_for,json=optimizeFor,enum=google.protobuf.FileOptions_OptimizeMode,def=1" json:"optimize_for,omitempty"`
	
	
	
	
	
	GoPackage *string `protobuf:"bytes,11,opt,name=go_package,json=goPackage" json:"go_package,omitempty"`
	
	
	
	
	
	
	
	
	
	
	CcGenericServices   *bool `protobuf:"varint,16,opt,name=cc_generic_services,json=ccGenericServices,def=0" json:"cc_generic_services,omitempty"`
	JavaGenericServices *bool `protobuf:"varint,17,opt,name=java_generic_services,json=javaGenericServices,def=0" json:"java_generic_services,omitempty"`
	PyGenericServices   *bool `protobuf:"varint,18,opt,name=py_generic_services,json=pyGenericServices,def=0" json:"py_generic_services,omitempty"`
	PhpGenericServices  *bool `protobuf:"varint,42,opt,name=php_generic_services,json=phpGenericServices,def=0" json:"php_generic_services,omitempty"`
	
	
	
	
	Deprecated *bool `protobuf:"varint,23,opt,name=deprecated,def=0" json:"deprecated,omitempty"`
	
	
	CcEnableArenas *bool `protobuf:"varint,31,opt,name=cc_enable_arenas,json=ccEnableArenas,def=0" json:"cc_enable_arenas,omitempty"`
	
	
	ObjcClassPrefix *string `protobuf:"bytes,36,opt,name=objc_class_prefix,json=objcClassPrefix" json:"objc_class_prefix,omitempty"`
	
	CsharpNamespace *string `protobuf:"bytes,37,opt,name=csharp_namespace,json=csharpNamespace" json:"csharp_namespace,omitempty"`
	
	
	
	
	SwiftPrefix *string `protobuf:"bytes,39,opt,name=swift_prefix,json=swiftPrefix" json:"swift_prefix,omitempty"`
	
	
	PhpClassPrefix *string `protobuf:"bytes,40,opt,name=php_class_prefix,json=phpClassPrefix" json:"php_class_prefix,omitempty"`
	
	
	
	PhpNamespace *string `protobuf:"bytes,41,opt,name=php_namespace,json=phpNamespace" json:"php_namespace,omitempty"`
	
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}               `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *FileOptions) Reset()         { *m = FileOptions{} }
func (m *FileOptions) String() string { return proto.CompactTextString(m) }
func (*FileOptions) ProtoMessage()    {}
func (*FileOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{10}
}

var extRange_FileOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*FileOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_FileOptions
}
func (m *FileOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FileOptions.Unmarshal(m, b)
}
func (m *FileOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FileOptions.Marshal(b, m, deterministic)
}
func (dst *FileOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FileOptions.Merge(dst, src)
}
func (m *FileOptions) XXX_Size() int {
	return xxx_messageInfo_FileOptions.Size(m)
}
func (m *FileOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_FileOptions.DiscardUnknown(m)
}

var xxx_messageInfo_FileOptions proto.InternalMessageInfo

const Default_FileOptions_JavaMultipleFiles bool = false
const Default_FileOptions_JavaStringCheckUtf8 bool = false
const Default_FileOptions_OptimizeFor FileOptions_OptimizeMode = FileOptions_SPEED
const Default_FileOptions_CcGenericServices bool = false
const Default_FileOptions_JavaGenericServices bool = false
const Default_FileOptions_PyGenericServices bool = false
const Default_FileOptions_PhpGenericServices bool = false
const Default_FileOptions_Deprecated bool = false
const Default_FileOptions_CcEnableArenas bool = false

func (m *FileOptions) GetJavaPackage() string {
	if m != nil && m.JavaPackage != nil {
		return *m.JavaPackage
	}
	return ""
}

func (m *FileOptions) GetJavaOuterClassname() string {
	if m != nil && m.JavaOuterClassname != nil {
		return *m.JavaOuterClassname
	}
	return ""
}

func (m *FileOptions) GetJavaMultipleFiles() bool {
	if m != nil && m.JavaMultipleFiles != nil {
		return *m.JavaMultipleFiles
	}
	return Default_FileOptions_JavaMultipleFiles
}


func (m *FileOptions) GetJavaGenerateEqualsAndHash() bool {
	if m != nil && m.JavaGenerateEqualsAndHash != nil {
		return *m.JavaGenerateEqualsAndHash
	}
	return false
}

func (m *FileOptions) GetJavaStringCheckUtf8() bool {
	if m != nil && m.JavaStringCheckUtf8 != nil {
		return *m.JavaStringCheckUtf8
	}
	return Default_FileOptions_JavaStringCheckUtf8
}

func (m *FileOptions) GetOptimizeFor() FileOptions_OptimizeMode {
	if m != nil && m.OptimizeFor != nil {
		return *m.OptimizeFor
	}
	return Default_FileOptions_OptimizeFor
}

func (m *FileOptions) GetGoPackage() string {
	if m != nil && m.GoPackage != nil {
		return *m.GoPackage
	}
	return ""
}

func (m *FileOptions) GetCcGenericServices() bool {
	if m != nil && m.CcGenericServices != nil {
		return *m.CcGenericServices
	}
	return Default_FileOptions_CcGenericServices
}

func (m *FileOptions) GetJavaGenericServices() bool {
	if m != nil && m.JavaGenericServices != nil {
		return *m.JavaGenericServices
	}
	return Default_FileOptions_JavaGenericServices
}

func (m *FileOptions) GetPyGenericServices() bool {
	if m != nil && m.PyGenericServices != nil {
		return *m.PyGenericServices
	}
	return Default_FileOptions_PyGenericServices
}

func (m *FileOptions) GetPhpGenericServices() bool {
	if m != nil && m.PhpGenericServices != nil {
		return *m.PhpGenericServices
	}
	return Default_FileOptions_PhpGenericServices
}

func (m *FileOptions) GetDeprecated() bool {
	if m != nil && m.Deprecated != nil {
		return *m.Deprecated
	}
	return Default_FileOptions_Deprecated
}

func (m *FileOptions) GetCcEnableArenas() bool {
	if m != nil && m.CcEnableArenas != nil {
		return *m.CcEnableArenas
	}
	return Default_FileOptions_CcEnableArenas
}

func (m *FileOptions) GetObjcClassPrefix() string {
	if m != nil && m.ObjcClassPrefix != nil {
		return *m.ObjcClassPrefix
	}
	return ""
}

func (m *FileOptions) GetCsharpNamespace() string {
	if m != nil && m.CsharpNamespace != nil {
		return *m.CsharpNamespace
	}
	return ""
}

func (m *FileOptions) GetSwiftPrefix() string {
	if m != nil && m.SwiftPrefix != nil {
		return *m.SwiftPrefix
	}
	return ""
}

func (m *FileOptions) GetPhpClassPrefix() string {
	if m != nil && m.PhpClassPrefix != nil {
		return *m.PhpClassPrefix
	}
	return ""
}

func (m *FileOptions) GetPhpNamespace() string {
	if m != nil && m.PhpNamespace != nil {
		return *m.PhpNamespace
	}
	return ""
}

func (m *FileOptions) GetUninterpretedOption() []*UninterpretedOption {
	if m != nil {
		return m.UninterpretedOption
	}
	return nil
}

type MessageOptions struct {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	MessageSetWireFormat *bool `protobuf:"varint,1,opt,name=message_set_wire_format,json=messageSetWireFormat,def=0" json:"message_set_wire_format,omitempty"`
	
	
	
	NoStandardDescriptorAccessor *bool `protobuf:"varint,2,opt,name=no_standard_descriptor_accessor,json=noStandardDescriptorAccessor,def=0" json:"no_standard_descriptor_accessor,omitempty"`
	
	
	
	
	Deprecated *bool `protobuf:"varint,3,opt,name=deprecated,def=0" json:"deprecated,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	MapEntry *bool `protobuf:"varint,7,opt,name=map_entry,json=mapEntry" json:"map_entry,omitempty"`
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}               `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *MessageOptions) Reset()         { *m = MessageOptions{} }
func (m *MessageOptions) String() string { return proto.CompactTextString(m) }
func (*MessageOptions) ProtoMessage()    {}
func (*MessageOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{11}
}

var extRange_MessageOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*MessageOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_MessageOptions
}
func (m *MessageOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MessageOptions.Unmarshal(m, b)
}
func (m *MessageOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MessageOptions.Marshal(b, m, deterministic)
}
func (dst *MessageOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MessageOptions.Merge(dst, src)
}
func (m *MessageOptions) XXX_Size() int {
	return xxx_messageInfo_MessageOptions.Size(m)
}
func (m *MessageOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_MessageOptions.DiscardUnknown(m)
}

var xxx_messageInfo_MessageOptions proto.InternalMessageInfo

const Default_MessageOptions_MessageSetWireFormat bool = false
const Default_MessageOptions_NoStandardDescriptorAccessor bool = false
const Default_MessageOptions_Deprecated bool = false

func (m *MessageOptions) GetMessageSetWireFormat() bool {
	if m != nil && m.MessageSetWireFormat != nil {
		return *m.MessageSetWireFormat
	}
	return Default_MessageOptions_MessageSetWireFormat
}

func (m *MessageOptions) GetNoStandardDescriptorAccessor() bool {
	if m != nil && m.NoStandardDescriptorAccessor != nil {
		return *m.NoStandardDescriptorAccessor
	}
	return Default_MessageOptions_NoStandardDescriptorAccessor
}

func (m *MessageOptions) GetDeprecated() bool {
	if m != nil && m.Deprecated != nil {
		return *m.Deprecated
	}
	return Default_MessageOptions_Deprecated
}

func (m *MessageOptions) GetMapEntry() bool {
	if m != nil && m.MapEntry != nil {
		return *m.MapEntry
	}
	return false
}

func (m *MessageOptions) GetUninterpretedOption() []*UninterpretedOption {
	if m != nil {
		return m.UninterpretedOption
	}
	return nil
}

type FieldOptions struct {
	
	
	
	
	Ctype *FieldOptions_CType `protobuf:"varint,1,opt,name=ctype,enum=google.protobuf.FieldOptions_CType,def=0" json:"ctype,omitempty"`
	
	
	
	
	
	Packed *bool `protobuf:"varint,2,opt,name=packed" json:"packed,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	Jstype *FieldOptions_JSType `protobuf:"varint,6,opt,name=jstype,enum=google.protobuf.FieldOptions_JSType,def=0" json:"jstype,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Lazy *bool `protobuf:"varint,5,opt,name=lazy,def=0" json:"lazy,omitempty"`
	
	
	
	
	Deprecated *bool `protobuf:"varint,3,opt,name=deprecated,def=0" json:"deprecated,omitempty"`
	
	Weak *bool `protobuf:"varint,10,opt,name=weak,def=0" json:"weak,omitempty"`
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}               `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *FieldOptions) Reset()         { *m = FieldOptions{} }
func (m *FieldOptions) String() string { return proto.CompactTextString(m) }
func (*FieldOptions) ProtoMessage()    {}
func (*FieldOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{12}
}

var extRange_FieldOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*FieldOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_FieldOptions
}
func (m *FieldOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FieldOptions.Unmarshal(m, b)
}
func (m *FieldOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FieldOptions.Marshal(b, m, deterministic)
}
func (dst *FieldOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FieldOptions.Merge(dst, src)
}
func (m *FieldOptions) XXX_Size() int {
	return xxx_messageInfo_FieldOptions.Size(m)
}
func (m *FieldOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_FieldOptions.DiscardUnknown(m)
}

var xxx_messageInfo_FieldOptions proto.InternalMessageInfo

const Default_FieldOptions_Ctype FieldOptions_CType = FieldOptions_STRING
const Default_FieldOptions_Jstype FieldOptions_JSType = FieldOptions_JS_NORMAL
const Default_FieldOptions_Lazy bool = false
const Default_FieldOptions_Deprecated bool = false
const Default_FieldOptions_Weak bool = false

func (m *FieldOptions) GetCtype() FieldOptions_CType {
	if m != nil && m.Ctype != nil {
		return *m.Ctype
	}
	return Default_FieldOptions_Ctype
}

func (m *FieldOptions) GetPacked() bool {
	if m != nil && m.Packed != nil {
		return *m.Packed
	}
	return false
}

func (m *FieldOptions) GetJstype() FieldOptions_JSType {
	if m != nil && m.Jstype != nil {
		return *m.Jstype
	}
	return Default_FieldOptions_Jstype
}

func (m *FieldOptions) GetLazy() bool {
	if m != nil && m.Lazy != nil {
		return *m.Lazy
	}
	return Default_FieldOptions_Lazy
}

func (m *FieldOptions) GetDeprecated() bool {
	if m != nil && m.Deprecated != nil {
		return *m.Deprecated
	}
	return Default_FieldOptions_Deprecated
}

func (m *FieldOptions) GetWeak() bool {
	if m != nil && m.Weak != nil {
		return *m.Weak
	}
	return Default_FieldOptions_Weak
}

func (m *FieldOptions) GetUninterpretedOption() []*UninterpretedOption {
	if m != nil {
		return m.UninterpretedOption
	}
	return nil
}

type OneofOptions struct {
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}               `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *OneofOptions) Reset()         { *m = OneofOptions{} }
func (m *OneofOptions) String() string { return proto.CompactTextString(m) }
func (*OneofOptions) ProtoMessage()    {}
func (*OneofOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{13}
}

var extRange_OneofOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*OneofOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_OneofOptions
}
func (m *OneofOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_OneofOptions.Unmarshal(m, b)
}
func (m *OneofOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_OneofOptions.Marshal(b, m, deterministic)
}
func (dst *OneofOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_OneofOptions.Merge(dst, src)
}
func (m *OneofOptions) XXX_Size() int {
	return xxx_messageInfo_OneofOptions.Size(m)
}
func (m *OneofOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_OneofOptions.DiscardUnknown(m)
}

var xxx_messageInfo_OneofOptions proto.InternalMessageInfo

func (m *OneofOptions) GetUninterpretedOption() []*UninterpretedOption {
	if m != nil {
		return m.UninterpretedOption
	}
	return nil
}

type EnumOptions struct {
	
	
	AllowAlias *bool `protobuf:"varint,2,opt,name=allow_alias,json=allowAlias" json:"allow_alias,omitempty"`
	
	
	
	
	Deprecated *bool `protobuf:"varint,3,opt,name=deprecated,def=0" json:"deprecated,omitempty"`
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}               `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *EnumOptions) Reset()         { *m = EnumOptions{} }
func (m *EnumOptions) String() string { return proto.CompactTextString(m) }
func (*EnumOptions) ProtoMessage()    {}
func (*EnumOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{14}
}

var extRange_EnumOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*EnumOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_EnumOptions
}
func (m *EnumOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnumOptions.Unmarshal(m, b)
}
func (m *EnumOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnumOptions.Marshal(b, m, deterministic)
}
func (dst *EnumOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnumOptions.Merge(dst, src)
}
func (m *EnumOptions) XXX_Size() int {
	return xxx_messageInfo_EnumOptions.Size(m)
}
func (m *EnumOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_EnumOptions.DiscardUnknown(m)
}

var xxx_messageInfo_EnumOptions proto.InternalMessageInfo

const Default_EnumOptions_Deprecated bool = false

func (m *EnumOptions) GetAllowAlias() bool {
	if m != nil && m.AllowAlias != nil {
		return *m.AllowAlias
	}
	return false
}

func (m *EnumOptions) GetDeprecated() bool {
	if m != nil && m.Deprecated != nil {
		return *m.Deprecated
	}
	return Default_EnumOptions_Deprecated
}

func (m *EnumOptions) GetUninterpretedOption() []*UninterpretedOption {
	if m != nil {
		return m.UninterpretedOption
	}
	return nil
}

type EnumValueOptions struct {
	
	
	
	
	Deprecated *bool `protobuf:"varint,1,opt,name=deprecated,def=0" json:"deprecated,omitempty"`
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}               `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *EnumValueOptions) Reset()         { *m = EnumValueOptions{} }
func (m *EnumValueOptions) String() string { return proto.CompactTextString(m) }
func (*EnumValueOptions) ProtoMessage()    {}
func (*EnumValueOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{15}
}

var extRange_EnumValueOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*EnumValueOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_EnumValueOptions
}
func (m *EnumValueOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_EnumValueOptions.Unmarshal(m, b)
}
func (m *EnumValueOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_EnumValueOptions.Marshal(b, m, deterministic)
}
func (dst *EnumValueOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_EnumValueOptions.Merge(dst, src)
}
func (m *EnumValueOptions) XXX_Size() int {
	return xxx_messageInfo_EnumValueOptions.Size(m)
}
func (m *EnumValueOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_EnumValueOptions.DiscardUnknown(m)
}

var xxx_messageInfo_EnumValueOptions proto.InternalMessageInfo

const Default_EnumValueOptions_Deprecated bool = false

func (m *EnumValueOptions) GetDeprecated() bool {
	if m != nil && m.Deprecated != nil {
		return *m.Deprecated
	}
	return Default_EnumValueOptions_Deprecated
}

func (m *EnumValueOptions) GetUninterpretedOption() []*UninterpretedOption {
	if m != nil {
		return m.UninterpretedOption
	}
	return nil
}

type ServiceOptions struct {
	
	
	
	
	Deprecated *bool `protobuf:"varint,33,opt,name=deprecated,def=0" json:"deprecated,omitempty"`
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}               `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *ServiceOptions) Reset()         { *m = ServiceOptions{} }
func (m *ServiceOptions) String() string { return proto.CompactTextString(m) }
func (*ServiceOptions) ProtoMessage()    {}
func (*ServiceOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{16}
}

var extRange_ServiceOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*ServiceOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_ServiceOptions
}
func (m *ServiceOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ServiceOptions.Unmarshal(m, b)
}
func (m *ServiceOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ServiceOptions.Marshal(b, m, deterministic)
}
func (dst *ServiceOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ServiceOptions.Merge(dst, src)
}
func (m *ServiceOptions) XXX_Size() int {
	return xxx_messageInfo_ServiceOptions.Size(m)
}
func (m *ServiceOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_ServiceOptions.DiscardUnknown(m)
}

var xxx_messageInfo_ServiceOptions proto.InternalMessageInfo

const Default_ServiceOptions_Deprecated bool = false

func (m *ServiceOptions) GetDeprecated() bool {
	if m != nil && m.Deprecated != nil {
		return *m.Deprecated
	}
	return Default_ServiceOptions_Deprecated
}

func (m *ServiceOptions) GetUninterpretedOption() []*UninterpretedOption {
	if m != nil {
		return m.UninterpretedOption
	}
	return nil
}

type MethodOptions struct {
	
	
	
	
	Deprecated       *bool                           `protobuf:"varint,33,opt,name=deprecated,def=0" json:"deprecated,omitempty"`
	IdempotencyLevel *MethodOptions_IdempotencyLevel `protobuf:"varint,34,opt,name=idempotency_level,json=idempotencyLevel,enum=google.protobuf.MethodOptions_IdempotencyLevel,def=0" json:"idempotency_level,omitempty"`
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	XXX_NoUnkeyedLiteral         struct{}               `json:"-"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
	XXX_sizecache                int32  `json:"-"`
}

func (m *MethodOptions) Reset()         { *m = MethodOptions{} }
func (m *MethodOptions) String() string { return proto.CompactTextString(m) }
func (*MethodOptions) ProtoMessage()    {}
func (*MethodOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{17}
}

var extRange_MethodOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*MethodOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_MethodOptions
}
func (m *MethodOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_MethodOptions.Unmarshal(m, b)
}
func (m *MethodOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_MethodOptions.Marshal(b, m, deterministic)
}
func (dst *MethodOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_MethodOptions.Merge(dst, src)
}
func (m *MethodOptions) XXX_Size() int {
	return xxx_messageInfo_MethodOptions.Size(m)
}
func (m *MethodOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_MethodOptions.DiscardUnknown(m)
}

var xxx_messageInfo_MethodOptions proto.InternalMessageInfo

const Default_MethodOptions_Deprecated bool = false
const Default_MethodOptions_IdempotencyLevel MethodOptions_IdempotencyLevel = MethodOptions_IDEMPOTENCY_UNKNOWN

func (m *MethodOptions) GetDeprecated() bool {
	if m != nil && m.Deprecated != nil {
		return *m.Deprecated
	}
	return Default_MethodOptions_Deprecated
}

func (m *MethodOptions) GetIdempotencyLevel() MethodOptions_IdempotencyLevel {
	if m != nil && m.IdempotencyLevel != nil {
		return *m.IdempotencyLevel
	}
	return Default_MethodOptions_IdempotencyLevel
}

func (m *MethodOptions) GetUninterpretedOption() []*UninterpretedOption {
	if m != nil {
		return m.UninterpretedOption
	}
	return nil
}







type UninterpretedOption struct {
	Name []*UninterpretedOption_NamePart `protobuf:"bytes,2,rep,name=name" json:"name,omitempty"`
	
	
	IdentifierValue      *string  `protobuf:"bytes,3,opt,name=identifier_value,json=identifierValue" json:"identifier_value,omitempty"`
	PositiveIntValue     *uint64  `protobuf:"varint,4,opt,name=positive_int_value,json=positiveIntValue" json:"positive_int_value,omitempty"`
	NegativeIntValue     *int64   `protobuf:"varint,5,opt,name=negative_int_value,json=negativeIntValue" json:"negative_int_value,omitempty"`
	DoubleValue          *float64 `protobuf:"fixed64,6,opt,name=double_value,json=doubleValue" json:"double_value,omitempty"`
	StringValue          []byte   `protobuf:"bytes,7,opt,name=string_value,json=stringValue" json:"string_value,omitempty"`
	AggregateValue       *string  `protobuf:"bytes,8,opt,name=aggregate_value,json=aggregateValue" json:"aggregate_value,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UninterpretedOption) Reset()         { *m = UninterpretedOption{} }
func (m *UninterpretedOption) String() string { return proto.CompactTextString(m) }
func (*UninterpretedOption) ProtoMessage()    {}
func (*UninterpretedOption) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{18}
}
func (m *UninterpretedOption) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UninterpretedOption.Unmarshal(m, b)
}
func (m *UninterpretedOption) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UninterpretedOption.Marshal(b, m, deterministic)
}
func (dst *UninterpretedOption) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UninterpretedOption.Merge(dst, src)
}
func (m *UninterpretedOption) XXX_Size() int {
	return xxx_messageInfo_UninterpretedOption.Size(m)
}
func (m *UninterpretedOption) XXX_DiscardUnknown() {
	xxx_messageInfo_UninterpretedOption.DiscardUnknown(m)
}

var xxx_messageInfo_UninterpretedOption proto.InternalMessageInfo

func (m *UninterpretedOption) GetName() []*UninterpretedOption_NamePart {
	if m != nil {
		return m.Name
	}
	return nil
}

func (m *UninterpretedOption) GetIdentifierValue() string {
	if m != nil && m.IdentifierValue != nil {
		return *m.IdentifierValue
	}
	return ""
}

func (m *UninterpretedOption) GetPositiveIntValue() uint64 {
	if m != nil && m.PositiveIntValue != nil {
		return *m.PositiveIntValue
	}
	return 0
}

func (m *UninterpretedOption) GetNegativeIntValue() int64 {
	if m != nil && m.NegativeIntValue != nil {
		return *m.NegativeIntValue
	}
	return 0
}

func (m *UninterpretedOption) GetDoubleValue() float64 {
	if m != nil && m.DoubleValue != nil {
		return *m.DoubleValue
	}
	return 0
}

func (m *UninterpretedOption) GetStringValue() []byte {
	if m != nil {
		return m.StringValue
	}
	return nil
}

func (m *UninterpretedOption) GetAggregateValue() string {
	if m != nil && m.AggregateValue != nil {
		return *m.AggregateValue
	}
	return ""
}






type UninterpretedOption_NamePart struct {
	NamePart             *string  `protobuf:"bytes,1,req,name=name_part,json=namePart" json:"name_part,omitempty"`
	IsExtension          *bool    `protobuf:"varint,2,req,name=is_extension,json=isExtension" json:"is_extension,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UninterpretedOption_NamePart) Reset()         { *m = UninterpretedOption_NamePart{} }
func (m *UninterpretedOption_NamePart) String() string { return proto.CompactTextString(m) }
func (*UninterpretedOption_NamePart) ProtoMessage()    {}
func (*UninterpretedOption_NamePart) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{18, 0}
}
func (m *UninterpretedOption_NamePart) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UninterpretedOption_NamePart.Unmarshal(m, b)
}
func (m *UninterpretedOption_NamePart) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UninterpretedOption_NamePart.Marshal(b, m, deterministic)
}
func (dst *UninterpretedOption_NamePart) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UninterpretedOption_NamePart.Merge(dst, src)
}
func (m *UninterpretedOption_NamePart) XXX_Size() int {
	return xxx_messageInfo_UninterpretedOption_NamePart.Size(m)
}
func (m *UninterpretedOption_NamePart) XXX_DiscardUnknown() {
	xxx_messageInfo_UninterpretedOption_NamePart.DiscardUnknown(m)
}

var xxx_messageInfo_UninterpretedOption_NamePart proto.InternalMessageInfo

func (m *UninterpretedOption_NamePart) GetNamePart() string {
	if m != nil && m.NamePart != nil {
		return *m.NamePart
	}
	return ""
}

func (m *UninterpretedOption_NamePart) GetIsExtension() bool {
	if m != nil && m.IsExtension != nil {
		return *m.IsExtension
	}
	return false
}



type SourceCodeInfo struct {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Location             []*SourceCodeInfo_Location `protobuf:"bytes,1,rep,name=location" json:"location,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                   `json:"-"`
	XXX_unrecognized     []byte                     `json:"-"`
	XXX_sizecache        int32                      `json:"-"`
}

func (m *SourceCodeInfo) Reset()         { *m = SourceCodeInfo{} }
func (m *SourceCodeInfo) String() string { return proto.CompactTextString(m) }
func (*SourceCodeInfo) ProtoMessage()    {}
func (*SourceCodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{19}
}
func (m *SourceCodeInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SourceCodeInfo.Unmarshal(m, b)
}
func (m *SourceCodeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SourceCodeInfo.Marshal(b, m, deterministic)
}
func (dst *SourceCodeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SourceCodeInfo.Merge(dst, src)
}
func (m *SourceCodeInfo) XXX_Size() int {
	return xxx_messageInfo_SourceCodeInfo.Size(m)
}
func (m *SourceCodeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_SourceCodeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_SourceCodeInfo proto.InternalMessageInfo

func (m *SourceCodeInfo) GetLocation() []*SourceCodeInfo_Location {
	if m != nil {
		return m.Location
	}
	return nil
}

type SourceCodeInfo_Location struct {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Path []int32 `protobuf:"varint,1,rep,packed,name=path" json:"path,omitempty"`
	
	
	
	
	
	Span []int32 `protobuf:"varint,2,rep,packed,name=span" json:"span,omitempty"`
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	LeadingComments         *string  `protobuf:"bytes,3,opt,name=leading_comments,json=leadingComments" json:"leading_comments,omitempty"`
	TrailingComments        *string  `protobuf:"bytes,4,opt,name=trailing_comments,json=trailingComments" json:"trailing_comments,omitempty"`
	LeadingDetachedComments []string `protobuf:"bytes,6,rep,name=leading_detached_comments,json=leadingDetachedComments" json:"leading_detached_comments,omitempty"`
	XXX_NoUnkeyedLiteral    struct{} `json:"-"`
	XXX_unrecognized        []byte   `json:"-"`
	XXX_sizecache           int32    `json:"-"`
}

func (m *SourceCodeInfo_Location) Reset()         { *m = SourceCodeInfo_Location{} }
func (m *SourceCodeInfo_Location) String() string { return proto.CompactTextString(m) }
func (*SourceCodeInfo_Location) ProtoMessage()    {}
func (*SourceCodeInfo_Location) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{19, 0}
}
func (m *SourceCodeInfo_Location) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_SourceCodeInfo_Location.Unmarshal(m, b)
}
func (m *SourceCodeInfo_Location) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_SourceCodeInfo_Location.Marshal(b, m, deterministic)
}
func (dst *SourceCodeInfo_Location) XXX_Merge(src proto.Message) {
	xxx_messageInfo_SourceCodeInfo_Location.Merge(dst, src)
}
func (m *SourceCodeInfo_Location) XXX_Size() int {
	return xxx_messageInfo_SourceCodeInfo_Location.Size(m)
}
func (m *SourceCodeInfo_Location) XXX_DiscardUnknown() {
	xxx_messageInfo_SourceCodeInfo_Location.DiscardUnknown(m)
}

var xxx_messageInfo_SourceCodeInfo_Location proto.InternalMessageInfo

func (m *SourceCodeInfo_Location) GetPath() []int32 {
	if m != nil {
		return m.Path
	}
	return nil
}

func (m *SourceCodeInfo_Location) GetSpan() []int32 {
	if m != nil {
		return m.Span
	}
	return nil
}

func (m *SourceCodeInfo_Location) GetLeadingComments() string {
	if m != nil && m.LeadingComments != nil {
		return *m.LeadingComments
	}
	return ""
}

func (m *SourceCodeInfo_Location) GetTrailingComments() string {
	if m != nil && m.TrailingComments != nil {
		return *m.TrailingComments
	}
	return ""
}

func (m *SourceCodeInfo_Location) GetLeadingDetachedComments() []string {
	if m != nil {
		return m.LeadingDetachedComments
	}
	return nil
}




type GeneratedCodeInfo struct {
	
	
	Annotation           []*GeneratedCodeInfo_Annotation `protobuf:"bytes,1,rep,name=annotation" json:"annotation,omitempty"`
	XXX_NoUnkeyedLiteral struct{}                        `json:"-"`
	XXX_unrecognized     []byte                          `json:"-"`
	XXX_sizecache        int32                           `json:"-"`
}

func (m *GeneratedCodeInfo) Reset()         { *m = GeneratedCodeInfo{} }
func (m *GeneratedCodeInfo) String() string { return proto.CompactTextString(m) }
func (*GeneratedCodeInfo) ProtoMessage()    {}
func (*GeneratedCodeInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{20}
}
func (m *GeneratedCodeInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeneratedCodeInfo.Unmarshal(m, b)
}
func (m *GeneratedCodeInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeneratedCodeInfo.Marshal(b, m, deterministic)
}
func (dst *GeneratedCodeInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeneratedCodeInfo.Merge(dst, src)
}
func (m *GeneratedCodeInfo) XXX_Size() int {
	return xxx_messageInfo_GeneratedCodeInfo.Size(m)
}
func (m *GeneratedCodeInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_GeneratedCodeInfo.DiscardUnknown(m)
}

var xxx_messageInfo_GeneratedCodeInfo proto.InternalMessageInfo

func (m *GeneratedCodeInfo) GetAnnotation() []*GeneratedCodeInfo_Annotation {
	if m != nil {
		return m.Annotation
	}
	return nil
}

type GeneratedCodeInfo_Annotation struct {
	
	
	Path []int32 `protobuf:"varint,1,rep,packed,name=path" json:"path,omitempty"`
	
	SourceFile *string `protobuf:"bytes,2,opt,name=source_file,json=sourceFile" json:"source_file,omitempty"`
	
	
	Begin *int32 `protobuf:"varint,3,opt,name=begin" json:"begin,omitempty"`
	
	
	
	End                  *int32   `protobuf:"varint,4,opt,name=end" json:"end,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GeneratedCodeInfo_Annotation) Reset()         { *m = GeneratedCodeInfo_Annotation{} }
func (m *GeneratedCodeInfo_Annotation) String() string { return proto.CompactTextString(m) }
func (*GeneratedCodeInfo_Annotation) ProtoMessage()    {}
func (*GeneratedCodeInfo_Annotation) Descriptor() ([]byte, []int) {
	return fileDescriptor_descriptor_4df4cb5f42392df6, []int{20, 0}
}
func (m *GeneratedCodeInfo_Annotation) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GeneratedCodeInfo_Annotation.Unmarshal(m, b)
}
func (m *GeneratedCodeInfo_Annotation) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GeneratedCodeInfo_Annotation.Marshal(b, m, deterministic)
}
func (dst *GeneratedCodeInfo_Annotation) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GeneratedCodeInfo_Annotation.Merge(dst, src)
}
func (m *GeneratedCodeInfo_Annotation) XXX_Size() int {
	return xxx_messageInfo_GeneratedCodeInfo_Annotation.Size(m)
}
func (m *GeneratedCodeInfo_Annotation) XXX_DiscardUnknown() {
	xxx_messageInfo_GeneratedCodeInfo_Annotation.DiscardUnknown(m)
}

var xxx_messageInfo_GeneratedCodeInfo_Annotation proto.InternalMessageInfo

func (m *GeneratedCodeInfo_Annotation) GetPath() []int32 {
	if m != nil {
		return m.Path
	}
	return nil
}

func (m *GeneratedCodeInfo_Annotation) GetSourceFile() string {
	if m != nil && m.SourceFile != nil {
		return *m.SourceFile
	}
	return ""
}

func (m *GeneratedCodeInfo_Annotation) GetBegin() int32 {
	if m != nil && m.Begin != nil {
		return *m.Begin
	}
	return 0
}

func (m *GeneratedCodeInfo_Annotation) GetEnd() int32 {
	if m != nil && m.End != nil {
		return *m.End
	}
	return 0
}

func init() {
	proto.RegisterType((*FileDescriptorSet)(nil), "google.protobuf.FileDescriptorSet")
	proto.RegisterType((*FileDescriptorProto)(nil), "google.protobuf.FileDescriptorProto")
	proto.RegisterType((*DescriptorProto)(nil), "google.protobuf.DescriptorProto")
	proto.RegisterType((*DescriptorProto_ExtensionRange)(nil), "google.protobuf.DescriptorProto.ExtensionRange")
	proto.RegisterType((*DescriptorProto_ReservedRange)(nil), "google.protobuf.DescriptorProto.ReservedRange")
	proto.RegisterType((*ExtensionRangeOptions)(nil), "google.protobuf.ExtensionRangeOptions")
	proto.RegisterType((*FieldDescriptorProto)(nil), "google.protobuf.FieldDescriptorProto")
	proto.RegisterType((*OneofDescriptorProto)(nil), "google.protobuf.OneofDescriptorProto")
	proto.RegisterType((*EnumDescriptorProto)(nil), "google.protobuf.EnumDescriptorProto")
	proto.RegisterType((*EnumDescriptorProto_EnumReservedRange)(nil), "google.protobuf.EnumDescriptorProto.EnumReservedRange")
	proto.RegisterType((*EnumValueDescriptorProto)(nil), "google.protobuf.EnumValueDescriptorProto")
	proto.RegisterType((*ServiceDescriptorProto)(nil), "google.protobuf.ServiceDescriptorProto")
	proto.RegisterType((*MethodDescriptorProto)(nil), "google.protobuf.MethodDescriptorProto")
	proto.RegisterType((*FileOptions)(nil), "google.protobuf.FileOptions")
	proto.RegisterType((*MessageOptions)(nil), "google.protobuf.MessageOptions")
	proto.RegisterType((*FieldOptions)(nil), "google.protobuf.FieldOptions")
	proto.RegisterType((*OneofOptions)(nil), "google.protobuf.OneofOptions")
	proto.RegisterType((*EnumOptions)(nil), "google.protobuf.EnumOptions")
	proto.RegisterType((*EnumValueOptions)(nil), "google.protobuf.EnumValueOptions")
	proto.RegisterType((*ServiceOptions)(nil), "google.protobuf.ServiceOptions")
	proto.RegisterType((*MethodOptions)(nil), "google.protobuf.MethodOptions")
	proto.RegisterType((*UninterpretedOption)(nil), "google.protobuf.UninterpretedOption")
	proto.RegisterType((*UninterpretedOption_NamePart)(nil), "google.protobuf.UninterpretedOption.NamePart")
	proto.RegisterType((*SourceCodeInfo)(nil), "google.protobuf.SourceCodeInfo")
	proto.RegisterType((*SourceCodeInfo_Location)(nil), "google.protobuf.SourceCodeInfo.Location")
	proto.RegisterType((*GeneratedCodeInfo)(nil), "google.protobuf.GeneratedCodeInfo")
	proto.RegisterType((*GeneratedCodeInfo_Annotation)(nil), "google.protobuf.GeneratedCodeInfo.Annotation")
	proto.RegisterEnum("google.protobuf.FieldDescriptorProto_Type", FieldDescriptorProto_Type_name, FieldDescriptorProto_Type_value)
	proto.RegisterEnum("google.protobuf.FieldDescriptorProto_Label", FieldDescriptorProto_Label_name, FieldDescriptorProto_Label_value)
	proto.RegisterEnum("google.protobuf.FileOptions_OptimizeMode", FileOptions_OptimizeMode_name, FileOptions_OptimizeMode_value)
	proto.RegisterEnum("google.protobuf.FieldOptions_CType", FieldOptions_CType_name, FieldOptions_CType_value)
	proto.RegisterEnum("google.protobuf.FieldOptions_JSType", FieldOptions_JSType_name, FieldOptions_JSType_value)
	proto.RegisterEnum("google.protobuf.MethodOptions_IdempotencyLevel", MethodOptions_IdempotencyLevel_name, MethodOptions_IdempotencyLevel_value)
}

func init() {
	proto.RegisterFile("google/protobuf/descriptor.proto", fileDescriptor_descriptor_4df4cb5f42392df6)
}

var fileDescriptor_descriptor_4df4cb5f42392df6 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x59, 0xdd, 0x6e, 0x1b, 0xc7,
	0xf5, 0xcf, 0xf2, 0x4b, 0xe4, 0x21, 0x45, 0x8d, 0x46, 0x8a, 0xbd, 0x56, 0x3e, 0x2c, 0x33, 0x1f,
	0x96, 0x9d, 0x7f, 0xa8, 0xc0, 0xb1, 0x1d, 0x47, 0xfe, 0x23, 0x2d, 0x45, 0xae, 0x15, 0xaa, 0x12,
	0xc9, 0x2e, 0xa9, 0xe6, 0x03, 0x28, 0x16, 0xa3, 0xdd, 0x21, 0xb9, 0xf6, 0x72, 0x77, 0xb3, 0xbb,
	0xb4, 0xad, 0xa0, 0x17, 0x06, 0x7a, 0xd5, 0xab, 0xde, 0x16, 0x45, 0xd1, 0x8b, 0xde, 0x04, 0xe8,
	0x03, 0x14, 0xc8, 0x5d, 0x9f, 0xa0, 0x40, 0xde, 0xa0, 0x68, 0x0b, 0xb4, 0x8f, 0xd0, 0xcb, 0x62,
	0x66, 0x76, 0x97, 0xbb, 0x24, 0x15, 0x2b, 0x01, 0xe2, 0x5c, 0x91, 0xf3, 0x9b, 0xdf, 0x39, 0x73,
	0xe6, 0xcc, 0x99, 0x33, 0x67, 0x66, 0x61, 0x7b, 0xe4, 0x38, 0x23, 0x8b, 0xee, 0xba, 0x9e, 0x13,
	0x38, 0xa7, 0xd3, 0xe1, 0xae, 0x41, 0x7d, 0xdd, 0x33, 0xdd, 0xc0, 0xf1, 0xea, 0x1c, 0xc3, 0x6b,
	0x82, 0x51, 0x8f, 0x18, 0xb5, 0x63, 0x58, 0x7f, 0x60, 0x5a, 0xb4, 0x15, 0x13, 0xfb, 0x34, 0xc0,
	0xf7, 0x20, 0x37, 0x34, 0x2d, 0x2a, 0x4b, 0xdb, 0xd9, 0x9d, 0xf2, 0xad, 0x37, 0xeb, 0x73, 0x42,
	0xf5, 0xb4, 0x44, 0x8f, 0xc1, 0x2a, 0x97, 0xa8, 0xfd, 0x2b, 0x07, 0x1b, 0x4b, 0x7a, 0x31, 0x86,
	0x9c, 0x4d, 0x26, 0x4c, 0xa3, 0xb4, 0x53, 0x52, 0xf9, 0x7f, 0x2c, 0xc3, 0x8a, 0x4b, 0xf4, 0x47,
	0x64, 0x44, 0xe5, 0x0c, 0x87, 0xa3, 0x26, 0x7e, 0x1d, 0xc0, 0xa0, 0x2e, 0xb5, 0x0d, 0x6a, 0xeb,
	0x67, 0x72, 0x76, 0x3b, 0xbb, 0x53, 0x52, 0x13, 0x08, 0x7e, 0x07, 0xd6, 0xdd, 0xe9, 0xa9, 0x65,
	0xea, 0x5a, 0x82, 0x06, 0xdb, 0xd9, 0x9d, 0xbc, 0x8a, 0x44, 0x47, 0x6b, 0x46, 0xbe, 0x0e, 0x6b,
	0x4f, 0x28, 0x79, 0x94, 0xa4, 0x96, 0x39, 0xb5, 0xca, 0xe0, 0x04, 0xb1, 0x09, 0x95, 0x09, 0xf5,
	0x7d, 0x32, 0xa2, 0x5a, 0x70, 0xe6, 0x52, 0x39, 0xc7, 0x67, 0xbf, 0xbd, 0x30, 0xfb, 0xf9, 0x99,
	0x97, 0x43, 0xa9, 0xc1, 0x99, 0x4b, 0x71, 0x03, 0x4a, 0xd4, 0x9e, 0x4e, 0x84, 0x86, 0xfc, 0x39,
	0xfe, 0x53, 0xec, 0xe9, 0x64, 0x5e, 0x4b, 0x91, 0x89, 0x85, 0x2a, 0x56, 0x7c, 0xea, 0x3d, 0x36,
	0x75, 0x2a, 0x17, 0xb8, 0x82, 0xeb, 0x0b, 0x0a, 0xfa, 0xa2, 0x7f, 0x5e, 0x47, 0x24, 0x87, 0x9b,
	0x50, 0xa2, 0x4f, 0x03, 0x6a, 0xfb, 0xa6, 0x63, 0xcb, 0x2b, 0x5c, 0xc9, 0x5b, 0x4b, 0x56, 0x91,
	0x5a, 0xc6, 0xbc, 0x8a, 0x99, 0x1c, 0xbe, 0x0b, 0x2b, 0x8e, 0x1b, 0x98, 0x8e, 0xed, 0xcb, 0xc5,
	0x6d, 0x69, 0xa7, 0x7c, 0xeb, 0xd5, 0xa5, 0x81, 0xd0, 0x15, 0x1c, 0x35, 0x22, 0xe3, 0x36, 0x20,
	0xdf, 0x99, 0x7a, 0x3a, 0xd5, 0x74, 0xc7, 0xa0, 0x9a, 0x69, 0x0f, 0x1d, 0xb9, 0xc4, 0x15, 0x5c,
	0x5d, 0x9c, 0x08, 0x27, 0x36, 0x1d, 0x83, 0xb6, 0xed, 0xa1, 0xa3, 0x56, 0xfd, 0x54, 0x1b, 0x5f,
	0x82, 0x82, 0x7f, 0x66, 0x07, 0xe4, 0xa9, 0x5c, 0xe1, 0x11, 0x12, 0xb6, 0x6a, 0x5f, 0x17, 0x60,
	0xed, 0x22, 0x21, 0x76, 0x1f, 0xf2, 0x43, 0x36, 0x4b, 0x39, 0xf3, 0x5d, 0x7c, 0x20, 0x64, 0xd2,
	0x4e, 0x2c, 0x7c, 0x4f, 0x27, 0x36, 0xa0, 0x6c, 0x53, 0x3f, 0xa0, 0x86, 0x88, 0x88, 0xec, 0x05,
	0x63, 0x0a, 0x84, 0xd0, 0x62, 0x48, 0xe5, 0xbe, 0x57, 0x48, 0x7d, 0x0a, 0x6b, 0xb1, 0x49, 0x9a,
	0x47, 0xec, 0x51, 0x14, 0x9b, 0xbb, 0xcf, 0xb3, 0xa4, 0xae, 0x44, 0x72, 0x2a, 0x13, 0x53, 0xab,
	0x34, 0xd5, 0xc6, 0x2d, 0x00, 0xc7, 0xa6, 0xce, 0x50, 0x33, 0xa8, 0x6e, 0xc9, 0xc5, 0x73, 0xbc,
	0xd4, 0x65, 0x94, 0x05, 0x2f, 0x39, 0x02, 0xd5, 0x2d, 0xfc, 0xe1, 0x2c, 0xd4, 0x56, 0xce, 0x89,
	0x94, 0x63, 0xb1, 0xc9, 0x16, 0xa2, 0xed, 0x04, 0xaa, 0x1e, 0x65, 0x71, 0x4f, 0x8d, 0x70, 0x66,
	0x25, 0x6e, 0x44, 0xfd, 0xb9, 0x33, 0x53, 0x43, 0x31, 0x31, 0xb1, 0x55, 0x2f, 0xd9, 0xc4, 0x6f,
	0x40, 0x0c, 0x68, 0x3c, 0xac, 0x80, 0x67, 0xa1, 0x4a, 0x04, 0x76, 0xc8, 0x84, 0x6e, 0x7d, 0x09,
	0xd5, 0xb4, 0x7b, 0xf0, 0x26, 0xe4, 0xfd, 0x80, 0x78, 0x01, 0x8f, 0xc2, 0xbc, 0x2a, 0x1a, 0x18,
	0x41, 0x96, 0xda, 0x06, 0xcf, 0x72, 0x79, 0x95, 0xfd, 0xc5, 0x3f, 0x9d, 0x4d, 0x38, 0xcb, 0x27,
	0xfc, 0xf6, 0xe2, 0x8a, 0xa6, 0x34, 0xcf, 0xcf, 0x7b, 0xeb, 0x03, 0x58, 0x4d, 0x4d, 0xe0, 0xa2,
	0x43, 0xd7, 0x7e, 0x05, 0x2f, 0x2f, 0x55, 0x8d, 0x3f, 0x85, 0xcd, 0xa9, 0x6d, 0xda, 0x01, 0xf5,
	0x5c, 0x8f, 0xb2, 0x88, 0x15, 0x43, 0xc9, 0xff, 0x5e, 0x39, 0x27, 0xe6, 0x4e, 0x92, 0x6c, 0xa1,
	0x45, 0xdd, 0x98, 0x2e, 0x82, 0x37, 0x4b, 0xc5, 0xff, 0xac, 0xa0, 0x67, 0xcf, 0x9e, 0x3d, 0xcb,
	0xd4, 0x7e, 0x57, 0x80, 0xcd, 0x65, 0x7b, 0x66, 0xe9, 0xf6, 0xbd, 0x04, 0x05, 0x7b, 0x3a, 0x39,
	0xa5, 0x1e, 0x77, 0x52, 0x5e, 0x0d, 0x5b, 0xb8, 0x01, 0x79, 0x8b, 0x9c, 0x52, 0x4b, 0xce, 0x6d,
	0x4b, 0x3b, 0xd5, 0x5b, 0xef, 0x5c, 0x68, 0x57, 0xd6, 0x8f, 0x98, 0x88, 0x2a, 0x24, 0xf1, 0x47,
	0x90, 0x0b, 0x53, 0x34, 0xd3, 0x70, 0xf3, 0x62, 0x1a, 0xd8, 0x5e, 0x52, 0xb9, 0x1c, 0x7e, 0x05,
	0x4a, 0xec, 0x57, 0xc4, 0x46, 0x81, 0xdb, 0x5c, 0x64, 0x00, 0x8b, 0x0b, 0xbc, 0x05, 0x45, 0xbe,
	0x4d, 0x0c, 0x1a, 0x1d, 0x6d, 0x71, 0x9b, 0x05, 0x96, 0x41, 0x87, 0x64, 0x6a, 0x05, 0xda, 0x63,
	0x62, 0x4d, 0x29, 0x0f, 0xf8, 0x92, 0x5a, 0x09, 0xc1, 0x5f, 0x30, 0x0c, 0x5f, 0x85, 0xb2, 0xd8,
	0x55, 0xa6, 0x6d, 0xd0, 0xa7, 0x3c, 0x7b, 0xe6, 0x55, 0xb1, 0xd1, 0xda, 0x0c, 0x61, 0xc3, 0x3f,
	0xf4, 0x1d, 0x3b, 0x0a, 0x4d, 0x3e, 0x04, 0x03, 0xf8, 0xf0, 0x1f, 0xcc, 0x27, 0xee, 0xd7, 0x96,
	0x4f, 0x6f, 0x3e, 0xa6, 0x6a, 0x7f, 0xc9, 0x40, 0x8e, 0xe7, 0x8b, 0x35, 0x28, 0x0f, 0x3e, 0xeb,
	0x29, 0x5a, 0xab, 0x7b, 0xb2, 0x7f, 0xa4, 0x20, 0x09, 0x57, 0x01, 0x38, 0xf0, 0xe0, 0xa8, 0xdb,
	0x18, 0xa0, 0x4c, 0xdc, 0x6e, 0x77, 0x06, 0x77, 0x6f, 0xa3, 0x6c, 0x2c, 0x70, 0x22, 0x80, 0x5c,
	0x92, 0xf0, 0xfe, 0x2d, 0x94, 0xc7, 0x08, 0x2a, 0x42, 0x41, 0xfb, 0x53, 0xa5, 0x75, 0xf7, 0x36,
	0x2a, 0xa4, 0x91, 0xf7, 0x6f, 0xa1, 0x15, 0xbc, 0x0a, 0x25, 0x8e, 0xec, 0x77, 0xbb, 0x47, 0xa8,
	0x18, 0xeb, 0xec, 0x0f, 0xd4, 0x76, 0xe7, 0x00, 0x95, 0x62, 0x9d, 0x07, 0x6a, 0xf7, 0xa4, 0x87,
	0x20, 0xd6, 0x70, 0xac, 0xf4, 0xfb, 0x8d, 0x03, 0x05, 0x95, 0x63, 0xc6, 0xfe, 0x67, 0x03, 0xa5,
	0x8f, 0x2a, 0x29, 0xb3, 0xde, 0xbf, 0x85, 0x56, 0xe3, 0x21, 0x94, 0xce, 0xc9, 0x31, 0xaa, 0xe2,
	0x75, 0x58, 0x15, 0x43, 0x44, 0x46, 0xac, 0xcd, 0x41, 0x77, 0x6f, 0x23, 0x34, 0x33, 0x44, 0x68,
	0x59, 0x4f, 0x01, 0x77, 0x6f, 0x23, 0x5c, 0x6b, 0x42, 0x9e, 0x47, 0x17, 0xc6, 0x50, 0x3d, 0x6a,
	0xec, 0x2b, 0x47, 0x5a, 0xb7, 0x37, 0x68, 0x77, 0x3b, 0x8d, 0x23, 0x24, 0xcd, 0x30, 0x55, 0xf9,
	0xf9, 0x49, 0x5b, 0x55, 0x5a, 0x28, 0x93, 0xc4, 0x7a, 0x4a, 0x63, 0xa0, 0xb4, 0x50, 0xb6, 0xa6,
	0xc3, 0xe6, 0xb2, 0x3c, 0xb9, 0x74, 0x67, 0x24, 0x96, 0x38, 0x73, 0xce, 0x12, 0x73, 0x5d, 0x0b,
	0x4b, 0xfc, 0xcf, 0x0c, 0x6c, 0x2c, 0x39, 0x2b, 0x96, 0x0e, 0xf2, 0x13, 0xc8, 0x8b, 0x10, 0x15,
	0xa7, 0xe7, 0x8d, 0xa5, 0x87, 0x0e, 0x0f, 0xd8, 0x85, 0x13, 0x94, 0xcb, 0x25, 0x2b, 0x88, 0xec,
	0x39, 0x15, 0x04, 0x53, 0xb1, 0x90, 0xd3, 0x7f, 0xb9, 0x90, 0xd3, 0xc5, 0xb1, 0x77, 0xf7, 0x22,
	0xc7, 0x1e, 0xc7, 0xbe, 0x5b, 0x6e, 0xcf, 0x2f, 0xc9, 0xed, 0xf7, 0x61, 0x7d, 0x41, 0xd1, 0x85,
	0x73, 0xec, 0xaf, 0x25, 0x90, 0xcf, 0x73, 0xce, 0x73, 0x32, 0x5d, 0x26, 0x95, 0xe9, 0xee, 0xcf,
	0x7b, 0xf0, 0xda, 0xf9, 0x8b, 0xb0, 0xb0, 0xd6, 0x5f, 0x49, 0x70, 0x69, 0x79, 0xa5, 0xb8, 0xd4,
	0x86, 0x8f, 0xa0, 0x30, 0xa1, 0xc1, 0xd8, 0x89, 0xaa, 0xa5, 0xb7, 0x97, 0x9c, 0xc1, 0xac, 0x7b,
	0x7e, 0xb1, 0x43, 0xa9, 0xe4, 0x21, 0x9e, 0x3d, 0xaf, 0xdc, 0x13, 0xd6, 0x2c, 0x58, 0xfa, 0x9b,
	0x0c, 0xbc, 0xbc, 0x54, 0xf9, 0x52, 0x43, 0x5f, 0x03, 0x30, 0x6d, 0x77, 0x1a, 0x88, 0x8a, 0x48,
	0x24, 0xd8, 0x12, 0x47, 0x78, 0xf2, 0x62, 0xc9, 0x73, 0x1a, 0xc4, 0xfd, 0x59, 0xde, 0x0f, 0x02,
	0xe2, 0x84, 0x7b, 0x33, 0x43, 0x73, 0xdc, 0xd0, 0xd7, 0xcf, 0x99, 0xe9, 0x42, 0x60, 0xbe, 0x07,
	0x48, 0xb7, 0x4c, 0x6a, 0x07, 0x9a, 0x1f, 0x78, 0x94, 0x4c, 0x4c, 0x7b, 0xc4, 0x4f, 0x90, 0xe2,
	0x5e, 0x7e, 0x48, 0x2c, 0x9f, 0xaa, 0x6b, 0xa2, 0xbb, 0x1f, 0xf5, 0x32, 0x09, 0x1e, 0x40, 0x5e,
	0x42, 0xa2, 0x90, 0x92, 0x10, 0xdd, 0xb1, 0x44, 0xed, 0xeb, 0x22, 0x94, 0x13, 0x75, 0x35, 0xbe,
	0x06, 0x95, 0x87, 0xe4, 0x31, 0xd1, 0xa2, 0xbb, 0x92, 0xf0, 0x44, 0x99, 0x61, 0xbd, 0xf0, 0xbe,
	0xf4, 0x1e, 0x6c, 0x72, 0x8a, 0x33, 0x0d, 0xa8, 0xa7, 0xe9, 0x16, 0xf1, 0x7d, 0xee, 0xb4, 0x22,
	0xa7, 0x62, 0xd6, 0xd7, 0x65, 0x5d, 0xcd, 0xa8, 0x07, 0xdf, 0x81, 0x0d, 0x2e, 0x31, 0x99, 0x5a,
	0x81, 0xe9, 0x5a, 0x54, 0x63, 0xb7, 0x37, 0x9f, 0x9f, 0x24, 0xb1, 0x65, 0xeb, 0x8c, 0x71, 0x1c,
	0x12, 0x98, 0x45, 0x3e, 0x6e, 0xc1, 0x6b, 0x5c, 0x6c, 0x44, 0x6d, 0xea, 0x91, 0x80, 0x6a, 0xf4,
	0x8b, 0x29, 0xb1, 0x7c, 0x8d, 0xd8, 0x86, 0x36, 0x26, 0xfe, 0x58, 0xde, 0x64, 0x0a, 0xf6, 0x33,
	0xb2, 0xa4, 0x5e, 0x61, 0xc4, 0x83, 0x90, 0xa7, 0x70, 0x5a, 0xc3, 0x36, 0x3e, 0x26, 0xfe, 0x18,
	0xef, 0xc1, 0x25, 0xae, 0xc5, 0x0f, 0x3c, 0xd3, 0x1e, 0x69, 0xfa, 0x98, 0xea, 0x8f, 0xb4, 0x69,
	0x30, 0xbc, 0x27, 0xbf, 0x92, 0x1c, 0x9f, 0x5b, 0xd8, 0xe7, 0x9c, 0x26, 0xa3, 0x9c, 0x04, 0xc3,
	0x7b, 0xb8, 0x0f, 0x15, 0xb6, 0x18, 0x13, 0xf3, 0x4b, 0xaa, 0x0d, 0x1d, 0x8f, 0x1f, 0x8d, 0xd5,
	0x25, 0xa9, 0x29, 0xe1, 0xc1, 0x7a, 0x37, 0x14, 0x38, 0x76, 0x0c, 0xba, 0x97, 0xef, 0xf7, 0x14,
	0xa5, 0xa5, 0x96, 0x23, 0x2d, 0x0f, 0x1c, 0x8f, 0x05, 0xd4, 0xc8, 0x89, 0x1d, 0x5c, 0x16, 0x01,
	0x35, 0x72, 0x22, 0xf7, 0xde, 0x81, 0x0d, 0x5d, 0x17, 0x73, 0x36, 0x75, 0x2d, 0xbc, 0x63, 0xf9,
	0x32, 0x4a, 0x39, 0x4b, 0xd7, 0x0f, 0x04, 0x21, 0x8c, 0x71, 0x1f, 0x7f, 0x08, 0x2f, 0xcf, 0x9c,
	0x95, 0x14, 0x5c, 0x5f, 0x98, 0xe5, 0xbc, 0xe8, 0x1d, 0xd8, 0x70, 0xcf, 0x16, 0x05, 0x71, 0x6a,
	0x44, 0xf7, 0x6c, 0x5e, 0xec, 0x03, 0xd8, 0x74, 0xc7, 0xee, 0xa2, 0xdc, 0xcd, 0xa4, 0x1c, 0x76,
	0xc7, 0xee, 0xbc, 0xe0, 0x5b, 0xfc, 0xc2, 0xed, 0x51, 0x9d, 0x04, 0xd4, 0x90, 0x2f, 0x27, 0xe9,
	0x89, 0x0e, 0xbc, 0x0b, 0x48, 0xd7, 0x35, 0x6a, 0x93, 0x53, 0x8b, 0x6a, 0xc4, 0xa3, 0x36, 0xf1,
	0xe5, 0xab, 0x49, 0x72, 0x55, 0xd7, 0x15, 0xde, 0xdb, 0xe0, 0x9d, 0xf8, 0x26, 0xac, 0x3b, 0xa7,
	0x0f, 0x75, 0x11, 0x92, 0x9a, 0xeb, 0xd1, 0xa1, 0xf9, 0x54, 0x7e, 0x93, 0xfb, 0x77, 0x8d, 0x75,
	0xf0, 0x80, 0xec, 0x71, 0x18, 0xdf, 0x00, 0xa4, 0xfb, 0x63, 0xe2, 0xb9, 0x3c, 0x27, 0xfb, 0x2e,
	0xd1, 0xa9, 0xfc, 0x96, 0xa0, 0x0a, 0xbc, 0x13, 0xc1, 0x6c, 0x4b, 0xf8, 0x4f, 0xcc, 0x61, 0x10,
	0x69, 0xbc, 0x2e, 0xb6, 0x04, 0xc7, 0x42, 0x6d, 0x3b, 0x80, 0x98, 0x2b, 0x52, 0x03, 0xef, 0x70,
	0x5a, 0xd5, 0x1d, 0xbb, 0xc9, 0x71, 0xdf, 0x80, 0x55, 0xc6, 0x9c, 0x0d, 0x7a, 0x43, 0x14, 0x64,
	0xee, 0x38, 0x31, 0xe2, 0x0f, 0x56, 0x1b, 0xd7, 0xf6, 0xa0, 0x92, 0x8c, 0x4f, 0x5c, 0x02, 0x11,
	0xa1, 0x48, 0x62, 0xc5, 0x4a, 0xb3, 0xdb, 0x62, 0x65, 0xc6, 0xe7, 0x0a, 0xca, 0xb0, 0x72, 0xe7,
	0xa8, 0x3d, 0x50, 0x34, 0xf5, 0xa4, 0x33, 0x68, 0x1f, 0x2b, 0x28, 0x9b, 0xa8, 0xab, 0x0f, 0x73,
	0xc5, 0xb7, 0xd1, 0xf5, 0xda, 0x37, 0x19, 0xa8, 0xa6, 0x2f, 0x4a, 0xf8, 0xff, 0xe1, 0x72, 0xf4,
	0xaa, 0xe1, 0xd3, 0x40, 0x7b, 0x62, 0x7a, 0x7c, 0xe3, 0x4c, 0x88, 0x38, 0xc4, 0xe2, 0xa5, 0xdb,
	0x0c, 0x59, 0x7d, 0x1a, 0x7c, 0x62, 0x7a, 0x6c, 0x5b, 0x4c, 0x48, 0x80, 0x8f, 0xe0, 0xaa, 0xed,
	0x68, 0x7e, 0x40, 0x6c, 0x83, 0x78, 0x86, 0x36, 0x7b, 0x4f, 0xd2, 0x88, 0xae, 0x53, 0xdf, 0x77,
	0xc4, 0x81, 0x15, 0x6b, 0x79, 0xd5, 0x76, 0xfa, 0x21, 0x79, 0x96, 0xc9, 0x1b, 0x21, 0x75, 0x2e,
	0xcc, 0xb2, 0xe7, 0x85, 0xd9, 0x2b, 0x50, 0x9a, 0x10, 0x57, 0xa3, 0x76, 0xe0, 0x9d, 0xf1, 0xf2,
	0xb8, 0xa8, 0x16, 0x27, 0xc4, 0x55, 0x58, 0xfb, 0x85, 0xdc, 0x52, 0x0e, 0x73, 0xc5, 0x22, 0x2a,
	0x1d, 0xe6, 0x8a, 0x25, 0x04, 0xb5, 0x7f, 0x64, 0xa1, 0x92, 0x2c, 0x97, 0xd9, 0xed, 0x43, 0xe7,
	0x27, 0x8b, 0xc4, 0x73, 0xcf, 0x1b, 0xdf, 0x5a, 0x5c, 0xd7, 0x9b, 0xec, 0xc8, 0xd9, 0x2b, 0x88,
	0x22, 0x56, 0x15, 0x92, 0xec, 0xb8, 0x67, 0xd9, 0x86, 0x8a, 0xa2, 0xa1, 0xa8, 0x86, 0x2d, 0x7c,
	0x00, 0x85, 0x87, 0x3e, 0xd7, 0x5d, 0xe0, 0xba, 0xdf, 0xfc, 0x76, 0xdd, 0x87, 0x7d, 0xae, 0xbc,
	0x74, 0xd8, 0xd7, 0x3a, 0x5d, 0xf5, 0xb8, 0x71, 0xa4, 0x86, 0xe2, 0xf8, 0x0a, 0xe4, 0x2c, 0xf2,
	0xe5, 0x59, 0xfa, 0x70, 0xe2, 0xd0, 0x45, 0x17, 0xe1, 0x0a, 0xe4, 0x9e, 0x50, 0xf2, 0x28, 0x7d,
	0x24, 0x70, 0xe8, 0x07, 0xdc, 0x0c, 0xbb, 0x90, 0xe7, 0xfe, 0xc2, 0x00, 0xa1, 0xc7, 0xd0, 0x4b,
	0xb8, 0x08, 0xb9, 0x66, 0x57, 0x65, 0x1b, 0x02, 0x41, 0x45, 0xa0, 0x5a, 0xaf, 0xad, 0x34, 0x15,
	0x94, 0xa9, 0xdd, 0x81, 0x82, 0x70, 0x02, 0xdb, 0x2c, 0xb1, 0x1b, 0xd0, 0x4b, 0x61, 0x33, 0xd4,
	0x21, 0x45, 0xbd, 0x27, 0xc7, 0xfb, 0x8a, 0x8a, 0x32, 0xe9, 0xa5, 0xce, 0xa1, 0x7c, 0xcd, 0x87,
	0x4a, 0xb2, 0x5e, 0x7e, 0x31, 0x77, 0xe1, 0xbf, 0x4a, 0x50, 0x4e, 0xd4, 0xbf, 0xac, 0x70, 0x21,
	0x96, 0xe5, 0x3c, 0xd1, 0x88, 0x65, 0x12, 0x3f, 0x0c, 0x0d, 0xe0, 0x50, 0x83, 0x21, 0x17, 0x5d,
	0xba, 0x17, 0xb4, 0x45, 0xf2, 0xa8, 0x50, 0xfb, 0xa3, 0x04, 0x68, 0xbe, 0x00, 0x9d, 0x33, 0x53,
	0xfa, 0x31, 0xcd, 0xac, 0xfd, 0x41, 0x82, 0x6a, 0xba, 0xea, 0x9c, 0x33, 0xef, 0xda, 0x8f, 0x6a,
	0xde, 0xdf, 0x33, 0xb0, 0x9a, 0xaa, 0x35, 0x2f, 0x6a, 0xdd, 0x17, 0xb0, 0x6e, 0x1a, 0x74, 0xe2,
	0x3a, 0x01, 0xb5, 0xf5, 0x33, 0xcd, 0xa2, 0x8f, 0xa9, 0x25, 0xd7, 0x78, 0xd2, 0xd8, 0xfd, 0xf6,
	0x6a, 0xb6, 0xde, 0x9e, 0xc9, 0x1d, 0x31, 0xb1, 0xbd, 0x8d, 0x76, 0x4b, 0x39, 0xee, 0x75, 0x07,
	0x4a, 0xa7, 0xf9, 0x99, 0x76, 0xd2, 0xf9, 0x59, 0xa7, 0xfb, 0x49, 0x47, 0x45, 0xe6, 0x1c, 0xed,
	0x07, 0xdc, 0xf6, 0x3d, 0x40, 0xf3, 0x46, 0xe1, 0xcb, 0xb0, 0xcc, 0x2c, 0xf4, 0x12, 0xde, 0x80,
	0xb5, 0x4e, 0x57, 0xeb, 0xb7, 0x5b, 0x8a, 0xa6, 0x3c, 0x78, 0xa0, 0x34, 0x07, 0x7d, 0xf1, 0x3e,
	0x11, 0xb3, 0x07, 0xa9, 0x0d, 0x5e, 0xfb, 0x7d, 0x16, 0x36, 0x96, 0x58, 0x82, 0x1b, 0xe1, 0xcd,
	0x42, 0x5c, 0x76, 0xde, 0xbd, 0x88, 0xf5, 0x75, 0x56, 0x10, 0xf4, 0x88, 0x17, 0x84, 0x17, 0x91,
	0x1b, 0xc0, 0xbc, 0x64, 0x07, 0xe6, 0xd0, 0xa4, 0x5e, 0xf8, 0x9c, 0x23, 0xae, 0x1b, 0x6b, 0x33,
	0x5c, 0xbc, 0xe8, 0xfc, 0x1f, 0x60, 0xd7, 0xf1, 0xcd, 0xc0, 0x7c, 0x4c, 0x35, 0xd3, 0x8e, 0xde,
	0x7e, 0xd8, 0xf5, 0x23, 0xa7, 0xa2, 0xa8, 0xa7, 0x6d, 0x07, 0x31, 0xdb, 0xa6, 0x23, 0x32, 0xc7,
	0x66, 0xc9, 0x3c, 0xab, 0xa2, 0xa8, 0x27, 0x66, 0x5f, 0x83, 0x8a, 0xe1, 0x4c, 0x59, 0x4d, 0x26,
	0x78, 0xec, 0xec, 0x90, 0xd4, 0xb2, 0xc0, 0x62, 0x4a, 0x58, 0x6d, 0xcf, 0x1e, 0x9d, 0x2a, 0x6a,
	0x59, 0x60, 0x82, 0x72, 0x1d, 0xd6, 0xc8, 0x68, 0xe4, 0x31, 0xe5, 0x91, 0x22, 0x71, 0x7f, 0xa8,
	0xc6, 0x30, 0x27, 0x6e, 0x1d, 0x42, 0x31, 0xf2, 0x03, 0x3b, 0xaa, 0x99, 0x27, 0x34, 0x57, 0x5c,
	0x8a, 0x33, 0x3b, 0x25, 0xb5, 0x68, 0x47, 0x9d, 0xd7, 0xa0, 0x62, 0xfa, 0xda, 0xec, 0x0d, 0x3d,
	0xb3, 0x9d, 0xd9, 0x29, 0xaa, 0x65, 0xd3, 0x8f, 0xdf, 0x1f, 0x6b, 0x5f, 0x65, 0xa0, 0x9a, 0xfe,
	0x06, 0x80, 0x5b, 0x50, 0xb4, 0x1c, 0x9d, 0xf0, 0xd0, 0x12, 0x1f, 0xa0, 0x76, 0x9e, 0xf3, 0xd9,
	0xa0, 0x7e, 0x14, 0xf2, 0xd5, 0x58, 0x72, 0xeb, 0x6f, 0x12, 0x14, 0x23, 0x18, 0x5f, 0x82, 0x9c,
	0x4b, 0x82, 0x31, 0x57, 0x97, 0xdf, 0xcf, 0x20, 0x49, 0xe5, 0x6d, 0x86, 0xfb, 0x2e, 0xb1, 0x79,
	0x08, 0x84, 0x38, 0x6b, 0xb3, 0x75, 0xb5, 0x28, 0x31, 0xf8, 0xe5, 0xc4, 0x99, 0x4c, 0xa8, 0x1d,
	0xf8, 0xd1, 0xba, 0x86, 0x78, 0x33, 0x84, 0xf1, 0x3b, 0xb0, 0x1e, 0x78, 0xc4, 0xb4, 0x52, 0xdc,
	0x1c, 0xe7, 0xa2, 0xa8, 0x23, 0x26, 0xef, 0xc1, 0x95, 0x48, 0xaf, 0x41, 0x03, 0xa2, 0x8f, 0xa9,
	0x31, 0x13, 0x2a, 0xf0, 0x47, 0x88, 0xcb, 0x21, 0xa1, 0x15, 0xf6, 0x47, 0xb2, 0xb5, 0x6f, 0x24,
	0x58, 0x8f, 0xae, 0x53, 0x46, 0xec, 0xac, 0x63, 0x00, 0x62, 0xdb, 0x4e, 0x90, 0x74, 0xd7, 0x62,
	0x28, 0x2f, 0xc8, 0xd5, 0x1b, 0xb1, 0x90, 0x9a, 0x50, 0xb0, 0x35, 0x01, 0x98, 0xf5, 0x9c, 0xeb,
	0xb6, 0xab, 0x50, 0x0e, 0x3f, 0xf0, 0xf0, 0xaf, 0x84, 0xe2, 0x02, 0x0e, 0x02, 0x62, 0xf7, 0x2e,
	0xbc, 0x09, 0xf9, 0x53, 0x3a, 0x32, 0xed, 0xf0, 0xd9, 0x56, 0x34, 0xa2, 0x67, 0x92, 0x5c, 0xfc,
	0x4c, 0xb2, 0xff, 0x5b, 0x09, 0x36, 0x74, 0x67, 0x32, 0x6f, 0xef, 0x3e, 0x9a, 0x7b, 0x05, 0xf0,
	0x3f, 0x96, 0x3e, 0xff, 0x68, 0x64, 0x06, 0xe3, 0xe9, 0x69, 0x5d, 0x77, 0x26, 0xbb, 0x23, 0xc7,
	0x22, 0xf6, 0x68, 0xf6, 0x99, 0x93, 0xff, 0xd1, 0xdf, 0x1d, 0x51, 0xfb, 0xdd, 0x91, 0x93, 0xf8,
	0xe8, 0x79, 0x7f, 0xf6, 0xf7, 0xbf, 0x92, 0xf4, 0xa7, 0x4c, 0xf6, 0xa0, 0xb7, 0xff, 0xe7, 0xcc,
	0xd6, 0x81, 0x18, 0xae, 0x17, 0xb9, 0x47, 0xa5, 0x43, 0x8b, 0xea, 0x6c, 0xca, 0xff, 0x0b, 0x00,
	0x00, 0xff, 0xff, 0x1a, 0x28, 0x25, 0x79, 0x42, 0x1d, 0x00, 0x00,
}
