



package descriptor

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.GoGoProtoPackageIsVersion2 

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
	return fileDescriptorDescriptor, []int{4, 0}
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
	return fileDescriptorDescriptor, []int{4, 1}
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
	return fileDescriptorDescriptor, []int{10, 0}
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
	return fileDescriptorDescriptor, []int{12, 0}
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
	return fileDescriptorDescriptor, []int{12, 1}
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
	return fileDescriptorDescriptor, []int{17, 0}
}



type FileDescriptorSet struct {
	File             []*FileDescriptorProto `protobuf:"bytes,1,rep,name=file" json:"file,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *FileDescriptorSet) Reset()                    { *m = FileDescriptorSet{} }
func (m *FileDescriptorSet) String() string            { return proto.CompactTextString(m) }
func (*FileDescriptorSet) ProtoMessage()               {}
func (*FileDescriptorSet) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{0} }

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
	
	
	Syntax           *string `protobuf:"bytes,12,opt,name=syntax" json:"syntax,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *FileDescriptorProto) Reset()                    { *m = FileDescriptorProto{} }
func (m *FileDescriptorProto) String() string            { return proto.CompactTextString(m) }
func (*FileDescriptorProto) ProtoMessage()               {}
func (*FileDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{1} }

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
	
	
	ReservedName     []string `protobuf:"bytes,10,rep,name=reserved_name,json=reservedName" json:"reserved_name,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *DescriptorProto) Reset()                    { *m = DescriptorProto{} }
func (m *DescriptorProto) String() string            { return proto.CompactTextString(m) }
func (*DescriptorProto) ProtoMessage()               {}
func (*DescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{2} }

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
	Start            *int32                 `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	End              *int32                 `protobuf:"varint,2,opt,name=end" json:"end,omitempty"`
	Options          *ExtensionRangeOptions `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *DescriptorProto_ExtensionRange) Reset()         { *m = DescriptorProto_ExtensionRange{} }
func (m *DescriptorProto_ExtensionRange) String() string { return proto.CompactTextString(m) }
func (*DescriptorProto_ExtensionRange) ProtoMessage()    {}
func (*DescriptorProto_ExtensionRange) Descriptor() ([]byte, []int) {
	return fileDescriptorDescriptor, []int{2, 0}
}

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
	Start            *int32 `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	End              *int32 `protobuf:"varint,2,opt,name=end" json:"end,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *DescriptorProto_ReservedRange) Reset()         { *m = DescriptorProto_ReservedRange{} }
func (m *DescriptorProto_ReservedRange) String() string { return proto.CompactTextString(m) }
func (*DescriptorProto_ReservedRange) ProtoMessage()    {}
func (*DescriptorProto_ReservedRange) Descriptor() ([]byte, []int) {
	return fileDescriptorDescriptor, []int{2, 1}
}

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
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *ExtensionRangeOptions) Reset()                    { *m = ExtensionRangeOptions{} }
func (m *ExtensionRangeOptions) String() string            { return proto.CompactTextString(m) }
func (*ExtensionRangeOptions) ProtoMessage()               {}
func (*ExtensionRangeOptions) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{3} }

var extRange_ExtensionRangeOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*ExtensionRangeOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_ExtensionRangeOptions
}

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
	
	
	
	
	JsonName         *string       `protobuf:"bytes,10,opt,name=json_name,json=jsonName" json:"json_name,omitempty"`
	Options          *FieldOptions `protobuf:"bytes,8,opt,name=options" json:"options,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *FieldDescriptorProto) Reset()                    { *m = FieldDescriptorProto{} }
func (m *FieldDescriptorProto) String() string            { return proto.CompactTextString(m) }
func (*FieldDescriptorProto) ProtoMessage()               {}
func (*FieldDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{4} }

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
	Name             *string       `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Options          *OneofOptions `protobuf:"bytes,2,opt,name=options" json:"options,omitempty"`
	XXX_unrecognized []byte        `json:"-"`
}

func (m *OneofDescriptorProto) Reset()                    { *m = OneofDescriptorProto{} }
func (m *OneofDescriptorProto) String() string            { return proto.CompactTextString(m) }
func (*OneofDescriptorProto) ProtoMessage()               {}
func (*OneofDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{5} }

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
	
	
	ReservedName     []string `protobuf:"bytes,5,rep,name=reserved_name,json=reservedName" json:"reserved_name,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *EnumDescriptorProto) Reset()                    { *m = EnumDescriptorProto{} }
func (m *EnumDescriptorProto) String() string            { return proto.CompactTextString(m) }
func (*EnumDescriptorProto) ProtoMessage()               {}
func (*EnumDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{6} }

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
	Start            *int32 `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	End              *int32 `protobuf:"varint,2,opt,name=end" json:"end,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *EnumDescriptorProto_EnumReservedRange) Reset()         { *m = EnumDescriptorProto_EnumReservedRange{} }
func (m *EnumDescriptorProto_EnumReservedRange) String() string { return proto.CompactTextString(m) }
func (*EnumDescriptorProto_EnumReservedRange) ProtoMessage()    {}
func (*EnumDescriptorProto_EnumReservedRange) Descriptor() ([]byte, []int) {
	return fileDescriptorDescriptor, []int{6, 0}
}

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
	Name             *string           `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Number           *int32            `protobuf:"varint,2,opt,name=number" json:"number,omitempty"`
	Options          *EnumValueOptions `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *EnumValueDescriptorProto) Reset()         { *m = EnumValueDescriptorProto{} }
func (m *EnumValueDescriptorProto) String() string { return proto.CompactTextString(m) }
func (*EnumValueDescriptorProto) ProtoMessage()    {}
func (*EnumValueDescriptorProto) Descriptor() ([]byte, []int) {
	return fileDescriptorDescriptor, []int{7}
}

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
	Name             *string                  `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Method           []*MethodDescriptorProto `protobuf:"bytes,2,rep,name=method" json:"method,omitempty"`
	Options          *ServiceOptions          `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
	XXX_unrecognized []byte                   `json:"-"`
}

func (m *ServiceDescriptorProto) Reset()                    { *m = ServiceDescriptorProto{} }
func (m *ServiceDescriptorProto) String() string            { return proto.CompactTextString(m) }
func (*ServiceDescriptorProto) ProtoMessage()               {}
func (*ServiceDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{8} }

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
	
	ServerStreaming  *bool  `protobuf:"varint,6,opt,name=server_streaming,json=serverStreaming,def=0" json:"server_streaming,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *MethodDescriptorProto) Reset()                    { *m = MethodDescriptorProto{} }
func (m *MethodDescriptorProto) String() string            { return proto.CompactTextString(m) }
func (*MethodDescriptorProto) ProtoMessage()               {}
func (*MethodDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{9} }

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
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *FileOptions) Reset()                    { *m = FileOptions{} }
func (m *FileOptions) String() string            { return proto.CompactTextString(m) }
func (*FileOptions) ProtoMessage()               {}
func (*FileOptions) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{10} }

var extRange_FileOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*FileOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_FileOptions
}

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
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *MessageOptions) Reset()                    { *m = MessageOptions{} }
func (m *MessageOptions) String() string            { return proto.CompactTextString(m) }
func (*MessageOptions) ProtoMessage()               {}
func (*MessageOptions) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{11} }

var extRange_MessageOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*MessageOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_MessageOptions
}

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
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *FieldOptions) Reset()                    { *m = FieldOptions{} }
func (m *FieldOptions) String() string            { return proto.CompactTextString(m) }
func (*FieldOptions) ProtoMessage()               {}
func (*FieldOptions) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{12} }

var extRange_FieldOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*FieldOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_FieldOptions
}

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
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *OneofOptions) Reset()                    { *m = OneofOptions{} }
func (m *OneofOptions) String() string            { return proto.CompactTextString(m) }
func (*OneofOptions) ProtoMessage()               {}
func (*OneofOptions) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{13} }

var extRange_OneofOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*OneofOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_OneofOptions
}

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
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *EnumOptions) Reset()                    { *m = EnumOptions{} }
func (m *EnumOptions) String() string            { return proto.CompactTextString(m) }
func (*EnumOptions) ProtoMessage()               {}
func (*EnumOptions) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{14} }

var extRange_EnumOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*EnumOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_EnumOptions
}

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
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *EnumValueOptions) Reset()                    { *m = EnumValueOptions{} }
func (m *EnumValueOptions) String() string            { return proto.CompactTextString(m) }
func (*EnumValueOptions) ProtoMessage()               {}
func (*EnumValueOptions) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{15} }

var extRange_EnumValueOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*EnumValueOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_EnumValueOptions
}

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
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *ServiceOptions) Reset()                    { *m = ServiceOptions{} }
func (m *ServiceOptions) String() string            { return proto.CompactTextString(m) }
func (*ServiceOptions) ProtoMessage()               {}
func (*ServiceOptions) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{16} }

var extRange_ServiceOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*ServiceOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_ServiceOptions
}

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
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *MethodOptions) Reset()                    { *m = MethodOptions{} }
func (m *MethodOptions) String() string            { return proto.CompactTextString(m) }
func (*MethodOptions) ProtoMessage()               {}
func (*MethodOptions) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{17} }

var extRange_MethodOptions = []proto.ExtensionRange{
	{Start: 1000, End: 536870911},
}

func (*MethodOptions) ExtensionRangeArray() []proto.ExtensionRange {
	return extRange_MethodOptions
}

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
	
	
	IdentifierValue  *string  `protobuf:"bytes,3,opt,name=identifier_value,json=identifierValue" json:"identifier_value,omitempty"`
	PositiveIntValue *uint64  `protobuf:"varint,4,opt,name=positive_int_value,json=positiveIntValue" json:"positive_int_value,omitempty"`
	NegativeIntValue *int64   `protobuf:"varint,5,opt,name=negative_int_value,json=negativeIntValue" json:"negative_int_value,omitempty"`
	DoubleValue      *float64 `protobuf:"fixed64,6,opt,name=double_value,json=doubleValue" json:"double_value,omitempty"`
	StringValue      []byte   `protobuf:"bytes,7,opt,name=string_value,json=stringValue" json:"string_value,omitempty"`
	AggregateValue   *string  `protobuf:"bytes,8,opt,name=aggregate_value,json=aggregateValue" json:"aggregate_value,omitempty"`
	XXX_unrecognized []byte   `json:"-"`
}

func (m *UninterpretedOption) Reset()                    { *m = UninterpretedOption{} }
func (m *UninterpretedOption) String() string            { return proto.CompactTextString(m) }
func (*UninterpretedOption) ProtoMessage()               {}
func (*UninterpretedOption) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{18} }

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
	NamePart         *string `protobuf:"bytes,1,req,name=name_part,json=namePart" json:"name_part,omitempty"`
	IsExtension      *bool   `protobuf:"varint,2,req,name=is_extension,json=isExtension" json:"is_extension,omitempty"`
	XXX_unrecognized []byte  `json:"-"`
}

func (m *UninterpretedOption_NamePart) Reset()         { *m = UninterpretedOption_NamePart{} }
func (m *UninterpretedOption_NamePart) String() string { return proto.CompactTextString(m) }
func (*UninterpretedOption_NamePart) ProtoMessage()    {}
func (*UninterpretedOption_NamePart) Descriptor() ([]byte, []int) {
	return fileDescriptorDescriptor, []int{18, 0}
}

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
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Location         []*SourceCodeInfo_Location `protobuf:"bytes,1,rep,name=location" json:"location,omitempty"`
	XXX_unrecognized []byte                     `json:"-"`
}

func (m *SourceCodeInfo) Reset()                    { *m = SourceCodeInfo{} }
func (m *SourceCodeInfo) String() string            { return proto.CompactTextString(m) }
func (*SourceCodeInfo) ProtoMessage()               {}
func (*SourceCodeInfo) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{19} }

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
	XXX_unrecognized        []byte   `json:"-"`
}

func (m *SourceCodeInfo_Location) Reset()         { *m = SourceCodeInfo_Location{} }
func (m *SourceCodeInfo_Location) String() string { return proto.CompactTextString(m) }
func (*SourceCodeInfo_Location) ProtoMessage()    {}
func (*SourceCodeInfo_Location) Descriptor() ([]byte, []int) {
	return fileDescriptorDescriptor, []int{19, 0}
}

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
	
	
	Annotation       []*GeneratedCodeInfo_Annotation `protobuf:"bytes,1,rep,name=annotation" json:"annotation,omitempty"`
	XXX_unrecognized []byte                          `json:"-"`
}

func (m *GeneratedCodeInfo) Reset()                    { *m = GeneratedCodeInfo{} }
func (m *GeneratedCodeInfo) String() string            { return proto.CompactTextString(m) }
func (*GeneratedCodeInfo) ProtoMessage()               {}
func (*GeneratedCodeInfo) Descriptor() ([]byte, []int) { return fileDescriptorDescriptor, []int{20} }

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
	
	
	
	End              *int32 `protobuf:"varint,4,opt,name=end" json:"end,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *GeneratedCodeInfo_Annotation) Reset()         { *m = GeneratedCodeInfo_Annotation{} }
func (m *GeneratedCodeInfo_Annotation) String() string { return proto.CompactTextString(m) }
func (*GeneratedCodeInfo_Annotation) ProtoMessage()    {}
func (*GeneratedCodeInfo_Annotation) Descriptor() ([]byte, []int) {
	return fileDescriptorDescriptor, []int{20, 0}
}

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

func init() { proto.RegisterFile("descriptor.proto", fileDescriptorDescriptor) }

var fileDescriptorDescriptor = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x59, 0xcd, 0x6f, 0xdb, 0xc8,
	0x15, 0x5f, 0x7d, 0x5a, 0x7a, 0x92, 0xe5, 0xf1, 0xd8, 0x9b, 0x30, 0xde, 0x8f, 0x38, 0xda, 0x8f,
	0x38, 0x49, 0xab, 0x2c, 0x9c, 0xc4, 0xc9, 0x3a, 0xc5, 0xb6, 0xb2, 0xc4, 0x78, 0x95, 0xca, 0x92,
	0x4a, 0xc9, 0xdd, 0x64, 0x8b, 0x82, 0x18, 0x93, 0x23, 0x89, 0x09, 0x45, 0x72, 0x49, 0x2a, 0x89,
	0x83, 0x1e, 0x02, 0xf4, 0xd4, 0xff, 0xa0, 0x28, 0x8a, 0x1e, 0x7a, 0x59, 0xa0, 0xd7, 0x02, 0x05,
	0xda, 0x7b, 0xaf, 0x05, 0x7a, 0xef, 0xa1, 0x40, 0x0b, 0xb4, 0x7f, 0x42, 0x8f, 0xc5, 0xcc, 0x90,
	0x14, 0xf5, 0x95, 0x78, 0x17, 0x48, 0xf6, 0x64, 0xcf, 0xef, 0xfd, 0xde, 0xe3, 0x9b, 0x37, 0x6f,
	0xde, 0xbc, 0x19, 0x01, 0xd2, 0xa9, 0xa7, 0xb9, 0x86, 0xe3, 0xdb, 0x6e, 0xc5, 0x71, 0x6d, 0xdf,
	0xc6, 0x6b, 0x03, 0xdb, 0x1e, 0x98, 0x54, 0x8c, 0x4e, 0xc6, 0xfd, 0xf2, 0x11, 0xac, 0xdf, 0x33,
	0x4c, 0x5a, 0x8f, 0x88, 0x5d, 0xea, 0xe3, 0x3b, 0x90, 0xee, 0x1b, 0x26, 0x95, 0x12, 0xdb, 0xa9,
	0x9d, 0xc2, 0xee, 0x87, 0x95, 0x19, 0xa5, 0xca, 0xb4, 0x46, 0x87, 0xc1, 0x0a, 0xd7, 0x28, 0xff,
	0x3b, 0x0d, 0x1b, 0x0b, 0xa4, 0x18, 0x43, 0xda, 0x22, 0x23, 0x66, 0x31, 0xb1, 0x93, 0x57, 0xf8,
	0xff, 0x58, 0x82, 0x15, 0x87, 0x68, 0x8f, 0xc9, 0x80, 0x4a, 0x49, 0x0e, 0x87, 0x43, 0xfc, 0x3e,
	0x80, 0x4e, 0x1d, 0x6a, 0xe9, 0xd4, 0xd2, 0x4e, 0xa5, 0xd4, 0x76, 0x6a, 0x27, 0xaf, 0xc4, 0x10,
	0x7c, 0x0d, 0xd6, 0x9d, 0xf1, 0x89, 0x69, 0x68, 0x6a, 0x8c, 0x06, 0xdb, 0xa9, 0x9d, 0x8c, 0x82,
	0x84, 0xa0, 0x3e, 0x21, 0x5f, 0x86, 0xb5, 0xa7, 0x94, 0x3c, 0x8e, 0x53, 0x0b, 0x9c, 0x5a, 0x62,
	0x70, 0x8c, 0x58, 0x83, 0xe2, 0x88, 0x7a, 0x1e, 0x19, 0x50, 0xd5, 0x3f, 0x75, 0xa8, 0x94, 0xe6,
	0xb3, 0xdf, 0x9e, 0x9b, 0xfd, 0xec, 0xcc, 0x0b, 0x81, 0x56, 0xef, 0xd4, 0xa1, 0xb8, 0x0a, 0x79,
	0x6a, 0x8d, 0x47, 0xc2, 0x42, 0x66, 0x49, 0xfc, 0x64, 0x6b, 0x3c, 0x9a, 0xb5, 0x92, 0x63, 0x6a,
	0x81, 0x89, 0x15, 0x8f, 0xba, 0x4f, 0x0c, 0x8d, 0x4a, 0x59, 0x6e, 0xe0, 0xf2, 0x9c, 0x81, 0xae,
	0x90, 0xcf, 0xda, 0x08, 0xf5, 0x70, 0x0d, 0xf2, 0xf4, 0x99, 0x4f, 0x2d, 0xcf, 0xb0, 0x2d, 0x69,
	0x85, 0x1b, 0xf9, 0x68, 0xc1, 0x2a, 0x52, 0x53, 0x9f, 0x35, 0x31, 0xd1, 0xc3, 0x7b, 0xb0, 0x62,
	0x3b, 0xbe, 0x61, 0x5b, 0x9e, 0x94, 0xdb, 0x4e, 0xec, 0x14, 0x76, 0xdf, 0x5d, 0x98, 0x08, 0x6d,
	0xc1, 0x51, 0x42, 0x32, 0x6e, 0x00, 0xf2, 0xec, 0xb1, 0xab, 0x51, 0x55, 0xb3, 0x75, 0xaa, 0x1a,
	0x56, 0xdf, 0x96, 0xf2, 0xdc, 0xc0, 0xc5, 0xf9, 0x89, 0x70, 0x62, 0xcd, 0xd6, 0x69, 0xc3, 0xea,
	0xdb, 0x4a, 0xc9, 0x9b, 0x1a, 0xe3, 0x73, 0x90, 0xf5, 0x4e, 0x2d, 0x9f, 0x3c, 0x93, 0x8a, 0x3c,
	0x43, 0x82, 0x51, 0xf9, 0xcf, 0x59, 0x58, 0x3b, 0x4b, 0x8a, 0xdd, 0x85, 0x4c, 0x9f, 0xcd, 0x52,
	0x4a, 0x7e, 0x93, 0x18, 0x08, 0x9d, 0xe9, 0x20, 0x66, 0xbf, 0x65, 0x10, 0xab, 0x50, 0xb0, 0xa8,
	0xe7, 0x53, 0x5d, 0x64, 0x44, 0xea, 0x8c, 0x39, 0x05, 0x42, 0x69, 0x3e, 0xa5, 0xd2, 0xdf, 0x2a,
	0xa5, 0x1e, 0xc0, 0x5a, 0xe4, 0x92, 0xea, 0x12, 0x6b, 0x10, 0xe6, 0xe6, 0xf5, 0x57, 0x79, 0x52,
	0x91, 0x43, 0x3d, 0x85, 0xa9, 0x29, 0x25, 0x3a, 0x35, 0xc6, 0x75, 0x00, 0xdb, 0xa2, 0x76, 0x5f,
	0xd5, 0xa9, 0x66, 0x4a, 0xb9, 0x25, 0x51, 0x6a, 0x33, 0xca, 0x5c, 0x94, 0x6c, 0x81, 0x6a, 0x26,
	0xfe, 0x74, 0x92, 0x6a, 0x2b, 0x4b, 0x32, 0xe5, 0x48, 0x6c, 0xb2, 0xb9, 0x6c, 0x3b, 0x86, 0x92,
	0x4b, 0x59, 0xde, 0x53, 0x3d, 0x98, 0x59, 0x9e, 0x3b, 0x51, 0x79, 0xe5, 0xcc, 0x94, 0x40, 0x4d,
	0x4c, 0x6c, 0xd5, 0x8d, 0x0f, 0xf1, 0x07, 0x10, 0x01, 0x2a, 0x4f, 0x2b, 0xe0, 0x55, 0xa8, 0x18,
	0x82, 0x2d, 0x32, 0xa2, 0x5b, 0xcf, 0xa1, 0x34, 0x1d, 0x1e, 0xbc, 0x09, 0x19, 0xcf, 0x27, 0xae,
	0xcf, 0xb3, 0x30, 0xa3, 0x88, 0x01, 0x46, 0x90, 0xa2, 0x96, 0xce, 0xab, 0x5c, 0x46, 0x61, 0xff,
	0xe2, 0x1f, 0x4d, 0x26, 0x9c, 0xe2, 0x13, 0xfe, 0x78, 0x7e, 0x45, 0xa7, 0x2c, 0xcf, 0xce, 0x7b,
	0xeb, 0x36, 0xac, 0x4e, 0x4d, 0xe0, 0xac, 0x9f, 0x2e, 0xff, 0x02, 0xde, 0x5e, 0x68, 0x1a, 0x3f,
	0x80, 0xcd, 0xb1, 0x65, 0x58, 0x3e, 0x75, 0x1d, 0x97, 0xb2, 0x8c, 0x15, 0x9f, 0x92, 0xfe, 0xb3,
	0xb2, 0x24, 0xe7, 0x8e, 0xe3, 0x6c, 0x61, 0x45, 0xd9, 0x18, 0xcf, 0x83, 0x57, 0xf3, 0xb9, 0xff,
	0xae, 0xa0, 0x17, 0x2f, 0x5e, 0xbc, 0x48, 0x96, 0x7f, 0x9d, 0x85, 0xcd, 0x45, 0x7b, 0x66, 0xe1,
	0xf6, 0x3d, 0x07, 0x59, 0x6b, 0x3c, 0x3a, 0xa1, 0x2e, 0x0f, 0x52, 0x46, 0x09, 0x46, 0xb8, 0x0a,
	0x19, 0x93, 0x9c, 0x50, 0x53, 0x4a, 0x6f, 0x27, 0x76, 0x4a, 0xbb, 0xd7, 0xce, 0xb4, 0x2b, 0x2b,
	0x4d, 0xa6, 0xa2, 0x08, 0x4d, 0xfc, 0x19, 0xa4, 0x83, 0x12, 0xcd, 0x2c, 0x5c, 0x3d, 0x9b, 0x05,
	0xb6, 0x97, 0x14, 0xae, 0x87, 0xdf, 0x81, 0x3c, 0xfb, 0x2b, 0x72, 0x23, 0xcb, 0x7d, 0xce, 0x31,
	0x80, 0xe5, 0x05, 0xde, 0x82, 0x1c, 0xdf, 0x26, 0x3a, 0x0d, 0x8f, 0xb6, 0x68, 0xcc, 0x12, 0x4b,
	0xa7, 0x7d, 0x32, 0x36, 0x7d, 0xf5, 0x09, 0x31, 0xc7, 0x94, 0x27, 0x7c, 0x5e, 0x29, 0x06, 0xe0,
	0x4f, 0x19, 0x86, 0x2f, 0x42, 0x41, 0xec, 0x2a, 0xc3, 0xd2, 0xe9, 0x33, 0x5e, 0x3d, 0x33, 0x8a,
	0xd8, 0x68, 0x0d, 0x86, 0xb0, 0xcf, 0x3f, 0xf2, 0x6c, 0x2b, 0x4c, 0x4d, 0xfe, 0x09, 0x06, 0xf0,
	0xcf, 0xdf, 0x9e, 0x2d, 0xdc, 0xef, 0x2d, 0x9e, 0xde, 0x6c, 0x4e, 0x95, 0xff, 0x94, 0x84, 0x34,
	0xaf, 0x17, 0x6b, 0x50, 0xe8, 0x3d, 0xec, 0xc8, 0x6a, 0xbd, 0x7d, 0x7c, 0xd0, 0x94, 0x51, 0x02,
	0x97, 0x00, 0x38, 0x70, 0xaf, 0xd9, 0xae, 0xf6, 0x50, 0x32, 0x1a, 0x37, 0x5a, 0xbd, 0xbd, 0x9b,
	0x28, 0x15, 0x29, 0x1c, 0x0b, 0x20, 0x1d, 0x27, 0xdc, 0xd8, 0x45, 0x19, 0x8c, 0xa0, 0x28, 0x0c,
	0x34, 0x1e, 0xc8, 0xf5, 0xbd, 0x9b, 0x28, 0x3b, 0x8d, 0xdc, 0xd8, 0x45, 0x2b, 0x78, 0x15, 0xf2,
	0x1c, 0x39, 0x68, 0xb7, 0x9b, 0x28, 0x17, 0xd9, 0xec, 0xf6, 0x94, 0x46, 0xeb, 0x10, 0xe5, 0x23,
	0x9b, 0x87, 0x4a, 0xfb, 0xb8, 0x83, 0x20, 0xb2, 0x70, 0x24, 0x77, 0xbb, 0xd5, 0x43, 0x19, 0x15,
	0x22, 0xc6, 0xc1, 0xc3, 0x9e, 0xdc, 0x45, 0xc5, 0x29, 0xb7, 0x6e, 0xec, 0xa2, 0xd5, 0xe8, 0x13,
	0x72, 0xeb, 0xf8, 0x08, 0x95, 0xf0, 0x3a, 0xac, 0x8a, 0x4f, 0x84, 0x4e, 0xac, 0xcd, 0x40, 0x7b,
	0x37, 0x11, 0x9a, 0x38, 0x22, 0xac, 0xac, 0x4f, 0x01, 0x7b, 0x37, 0x11, 0x2e, 0xd7, 0x20, 0xc3,
	0xb3, 0x0b, 0x63, 0x28, 0x35, 0xab, 0x07, 0x72, 0x53, 0x6d, 0x77, 0x7a, 0x8d, 0x76, 0xab, 0xda,
	0x44, 0x89, 0x09, 0xa6, 0xc8, 0x3f, 0x39, 0x6e, 0x28, 0x72, 0x1d, 0x25, 0xe3, 0x58, 0x47, 0xae,
	0xf6, 0xe4, 0x3a, 0x4a, 0x95, 0x35, 0xd8, 0x5c, 0x54, 0x27, 0x17, 0xee, 0x8c, 0xd8, 0x12, 0x27,
	0x97, 0x2c, 0x31, 0xb7, 0x35, 0xb7, 0xc4, 0xff, 0x4a, 0xc2, 0xc6, 0x82, 0xb3, 0x62, 0xe1, 0x47,
	0x7e, 0x08, 0x19, 0x91, 0xa2, 0xe2, 0xf4, 0xbc, 0xb2, 0xf0, 0xd0, 0xe1, 0x09, 0x3b, 0x77, 0x82,
	0x72, 0xbd, 0x78, 0x07, 0x91, 0x5a, 0xd2, 0x41, 0x30, 0x13, 0x73, 0x35, 0xfd, 0xe7, 0x73, 0x35,
	0x5d, 0x1c, 0x7b, 0x7b, 0x67, 0x39, 0xf6, 0x38, 0xf6, 0xcd, 0x6a, 0x7b, 0x66, 0x41, 0x6d, 0xbf,
	0x0b, 0xeb, 0x73, 0x86, 0xce, 0x5c, 0x63, 0x7f, 0x99, 0x00, 0x69, 0x59, 0x70, 0x5e, 0x51, 0xe9,
	0x92, 0x53, 0x95, 0xee, 0xee, 0x6c, 0x04, 0x2f, 0x2d, 0x5f, 0x84, 0xb9, 0xb5, 0xfe, 0x3a, 0x01,
	0xe7, 0x16, 0x77, 0x8a, 0x0b, 0x7d, 0xf8, 0x0c, 0xb2, 0x23, 0xea, 0x0f, 0xed, 0xb0, 0x5b, 0xfa,
	0x78, 0xc1, 0x19, 0xcc, 0xc4, 0xb3, 0x8b, 0x1d, 0x68, 0xc5, 0x0f, 0xf1, 0xd4, 0xb2, 0x76, 0x4f,
	0x78, 0x33, 0xe7, 0xe9, 0xaf, 0x92, 0xf0, 0xf6, 0x42, 0xe3, 0x0b, 0x1d, 0x7d, 0x0f, 0xc0, 0xb0,
	0x9c, 0xb1, 0x2f, 0x3a, 0x22, 0x51, 0x60, 0xf3, 0x1c, 0xe1, 0xc5, 0x8b, 0x15, 0xcf, 0xb1, 0x1f,
	0xc9, 0x53, 0x5c, 0x0e, 0x02, 0xe2, 0x84, 0x3b, 0x13, 0x47, 0xd3, 0xdc, 0xd1, 0xf7, 0x97, 0xcc,
	0x74, 0x2e, 0x31, 0x3f, 0x01, 0xa4, 0x99, 0x06, 0xb5, 0x7c, 0xd5, 0xf3, 0x5d, 0x4a, 0x46, 0x86,
	0x35, 0xe0, 0x27, 0x48, 0x6e, 0x3f, 0xd3, 0x27, 0xa6, 0x47, 0x95, 0x35, 0x21, 0xee, 0x86, 0x52,
	0xa6, 0xc1, 0x13, 0xc8, 0x8d, 0x69, 0x64, 0xa7, 0x34, 0x84, 0x38, 0xd2, 0x28, 0xff, 0x31, 0x07,
	0x85, 0x58, 0x5f, 0x8d, 0x2f, 0x41, 0xf1, 0x11, 0x79, 0x42, 0xd4, 0xf0, 0xae, 0x24, 0x22, 0x51,
	0x60, 0x58, 0x27, 0xb8, 0x2f, 0x7d, 0x02, 0x9b, 0x9c, 0x62, 0x8f, 0x7d, 0xea, 0xaa, 0x9a, 0x49,
	0x3c, 0x8f, 0x07, 0x2d, 0xc7, 0xa9, 0x98, 0xc9, 0xda, 0x4c, 0x54, 0x0b, 0x25, 0xf8, 0x16, 0x6c,
	0x70, 0x8d, 0xd1, 0xd8, 0xf4, 0x0d, 0xc7, 0xa4, 0x2a, 0xbb, 0xbd, 0x79, 0xfc, 0x24, 0x89, 0x3c,
	0x5b, 0x67, 0x8c, 0xa3, 0x80, 0xc0, 0x3c, 0xf2, 0x70, 0x1d, 0xde, 0xe3, 0x6a, 0x03, 0x6a, 0x51,
	0x97, 0xf8, 0x54, 0xa5, 0x5f, 0x8d, 0x89, 0xe9, 0xa9, 0xc4, 0xd2, 0xd5, 0x21, 0xf1, 0x86, 0xd2,
	0x26, 0x33, 0x70, 0x90, 0x94, 0x12, 0xca, 0x05, 0x46, 0x3c, 0x0c, 0x78, 0x32, 0xa7, 0x55, 0x2d,
	0xfd, 0x73, 0xe2, 0x0d, 0xf1, 0x3e, 0x9c, 0xe3, 0x56, 0x3c, 0xdf, 0x35, 0xac, 0x81, 0xaa, 0x0d,
	0xa9, 0xf6, 0x58, 0x1d, 0xfb, 0xfd, 0x3b, 0xd2, 0x3b, 0xf1, 0xef, 0x73, 0x0f, 0xbb, 0x9c, 0x53,
	0x63, 0x94, 0x63, 0xbf, 0x7f, 0x07, 0x77, 0xa1, 0xc8, 0x16, 0x63, 0x64, 0x3c, 0xa7, 0x6a, 0xdf,
	0x76, 0xf9, 0xd1, 0x58, 0x5a, 0x50, 0x9a, 0x62, 0x11, 0xac, 0xb4, 0x03, 0x85, 0x23, 0x5b, 0xa7,
	0xfb, 0x99, 0x6e, 0x47, 0x96, 0xeb, 0x4a, 0x21, 0xb4, 0x72, 0xcf, 0x76, 0x59, 0x42, 0x0d, 0xec,
	0x28, 0xc0, 0x05, 0x91, 0x50, 0x03, 0x3b, 0x0c, 0xef, 0x2d, 0xd8, 0xd0, 0x34, 0x31, 0x67, 0x43,
	0x53, 0x83, 0x3b, 0x96, 0x27, 0xa1, 0xa9, 0x60, 0x69, 0xda, 0xa1, 0x20, 0x04, 0x39, 0xee, 0xe1,
	0x4f, 0xe1, 0xed, 0x49, 0xb0, 0xe2, 0x8a, 0xeb, 0x73, 0xb3, 0x9c, 0x55, 0xbd, 0x05, 0x1b, 0xce,
	0xe9, 0xbc, 0x22, 0x9e, 0xfa, 0xa2, 0x73, 0x3a, 0xab, 0x76, 0x1b, 0x36, 0x9d, 0xa1, 0x33, 0xaf,
	0x77, 0x35, 0xae, 0x87, 0x9d, 0xa1, 0x33, 0xab, 0xf8, 0x11, 0xbf, 0x70, 0xbb, 0x54, 0x23, 0x3e,
	0xd5, 0xa5, 0xf3, 0x71, 0x7a, 0x4c, 0x80, 0xaf, 0x03, 0xd2, 0x34, 0x95, 0x5a, 0xe4, 0xc4, 0xa4,
	0x2a, 0x71, 0xa9, 0x45, 0x3c, 0xe9, 0x62, 0x9c, 0x5c, 0xd2, 0x34, 0x99, 0x4b, 0xab, 0x5c, 0x88,
	0xaf, 0xc2, 0xba, 0x7d, 0xf2, 0x48, 0x13, 0x29, 0xa9, 0x3a, 0x2e, 0xed, 0x1b, 0xcf, 0xa4, 0x0f,
	0x79, 0x7c, 0xd7, 0x98, 0x80, 0x27, 0x64, 0x87, 0xc3, 0xf8, 0x0a, 0x20, 0xcd, 0x1b, 0x12, 0xd7,
	0xe1, 0x35, 0xd9, 0x73, 0x88, 0x46, 0xa5, 0x8f, 0x04, 0x55, 0xe0, 0xad, 0x10, 0x66, 0x5b, 0xc2,
	0x7b, 0x6a, 0xf4, 0xfd, 0xd0, 0xe2, 0x65, 0xb1, 0x25, 0x38, 0x16, 0x58, 0xdb, 0x01, 0xc4, 0x42,
	0x31, 0xf5, 0xe1, 0x1d, 0x4e, 0x2b, 0x39, 0x43, 0x27, 0xfe, 0xdd, 0x0f, 0x60, 0x95, 0x31, 0x27,
	0x1f, 0xbd, 0x22, 0x1a, 0x32, 0x67, 0x18, 0xfb, 0xe2, 0x6b, 0xeb, 0x8d, 0xcb, 0xfb, 0x50, 0x8c,
	0xe7, 0x27, 0xce, 0x83, 0xc8, 0x50, 0x94, 0x60, 0xcd, 0x4a, 0xad, 0x5d, 0x67, 0x6d, 0xc6, 0x97,
	0x32, 0x4a, 0xb2, 0x76, 0xa7, 0xd9, 0xe8, 0xc9, 0xaa, 0x72, 0xdc, 0xea, 0x35, 0x8e, 0x64, 0x94,
	0x8a, 0xf7, 0xd5, 0x7f, 0x4d, 0x42, 0x69, 0xfa, 0x8a, 0x84, 0x7f, 0x00, 0xe7, 0xc3, 0xf7, 0x0c,
	0x8f, 0xfa, 0xea, 0x53, 0xc3, 0xe5, 0x5b, 0x66, 0x44, 0xc4, 0xf1, 0x15, 0x2d, 0xda, 0x66, 0xc0,
	0xea, 0x52, 0xff, 0x0b, 0xc3, 0x65, 0x1b, 0x62, 0x44, 0x7c, 0xdc, 0x84, 0x8b, 0x96, 0xad, 0x7a,
	0x3e, 0xb1, 0x74, 0xe2, 0xea, 0xea, 0xe4, 0x25, 0x49, 0x25, 0x9a, 0x46, 0x3d, 0xcf, 0x16, 0x47,
	0x55, 0x64, 0xe5, 0x5d, 0xcb, 0xee, 0x06, 0xe4, 0x49, 0x0d, 0xaf, 0x06, 0xd4, 0x99, 0x04, 0x4b,
	0x2d, 0x4b, 0xb0, 0x77, 0x20, 0x3f, 0x22, 0x8e, 0x4a, 0x2d, 0xdf, 0x3d, 0xe5, 0x8d, 0x71, 0x4e,
	0xc9, 0x8d, 0x88, 0x23, 0xb3, 0xf1, 0x9b, 0xb9, 0x9f, 0xfc, 0x23, 0x05, 0xc5, 0x78, 0x73, 0xcc,
	0xee, 0x1a, 0x1a, 0x3f, 0x47, 0x12, 0xbc, 0xd2, 0x7c, 0xf0, 0xd2, 0x56, 0xba, 0x52, 0x63, 0x07,
	0xcc, 0x7e, 0x56, 0xb4, 0xac, 0x8a, 0xd0, 0x64, 0x87, 0x3b, 0xab, 0x2d, 0x54, 0xb4, 0x08, 0x39,
	0x25, 0x18, 0xe1, 0x43, 0xc8, 0x3e, 0xf2, 0xb8, 0xed, 0x2c, 0xb7, 0xfd, 0xe1, 0xcb, 0x6d, 0xdf,
	0xef, 0x72, 0xe3, 0xf9, 0xfb, 0x5d, 0xb5, 0xd5, 0x56, 0x8e, 0xaa, 0x4d, 0x25, 0x50, 0xc7, 0x17,
	0x20, 0x6d, 0x92, 0xe7, 0xa7, 0xd3, 0x47, 0x11, 0x87, 0xce, 0x1a, 0xf8, 0x0b, 0x90, 0x7e, 0x4a,
	0xc9, 0xe3, 0xe9, 0x03, 0x80, 0x43, 0xaf, 0x31, 0xf5, 0xaf, 0x43, 0x86, 0xc7, 0x0b, 0x03, 0x04,
	0x11, 0x43, 0x6f, 0xe1, 0x1c, 0xa4, 0x6b, 0x6d, 0x85, 0xa5, 0x3f, 0x82, 0xa2, 0x40, 0xd5, 0x4e,
	0x43, 0xae, 0xc9, 0x28, 0x59, 0xbe, 0x05, 0x59, 0x11, 0x04, 0xb6, 0x35, 0xa2, 0x30, 0xa0, 0xb7,
	0x82, 0x61, 0x60, 0x23, 0x11, 0x4a, 0x8f, 0x8f, 0x0e, 0x64, 0x05, 0x25, 0xe3, 0xcb, 0xeb, 0x41,
	0x31, 0xde, 0x17, 0xbf, 0x99, 0x9c, 0xfa, 0x4b, 0x02, 0x0a, 0xb1, 0x3e, 0x97, 0x35, 0x28, 0xc4,
	0x34, 0xed, 0xa7, 0x2a, 0x31, 0x0d, 0xe2, 0x05, 0x49, 0x01, 0x1c, 0xaa, 0x32, 0xe4, 0xac, 0x8b,
	0xf6, 0x46, 0x9c, 0xff, 0x5d, 0x02, 0xd0, 0x6c, 0x8b, 0x39, 0xe3, 0x60, 0xe2, 0x3b, 0x75, 0xf0,
	0xb7, 0x09, 0x28, 0x4d, 0xf7, 0x95, 0x33, 0xee, 0x5d, 0xfa, 0x4e, 0xdd, 0xfb, 0x67, 0x12, 0x56,
	0xa7, 0xba, 0xc9, 0xb3, 0x7a, 0xf7, 0x15, 0xac, 0x1b, 0x3a, 0x1d, 0x39, 0xb6, 0x4f, 0x2d, 0xed,
	0x54, 0x35, 0xe9, 0x13, 0x6a, 0x4a, 0x65, 0x5e, 0x28, 0xae, 0xbf, 0xbc, 0x5f, 0xad, 0x34, 0x26,
	0x7a, 0x4d, 0xa6, 0xb6, 0xbf, 0xd1, 0xa8, 0xcb, 0x47, 0x9d, 0x76, 0x4f, 0x6e, 0xd5, 0x1e, 0xaa,
	0xc7, 0xad, 0x1f, 0xb7, 0xda, 0x5f, 0xb4, 0x14, 0x64, 0xcc, 0xd0, 0x5e, 0xe3, 0x56, 0xef, 0x00,
	0x9a, 0x75, 0x0a, 0x9f, 0x87, 0x45, 0x6e, 0xa1, 0xb7, 0xf0, 0x06, 0xac, 0xb5, 0xda, 0x6a, 0xb7,
	0x51, 0x97, 0x55, 0xf9, 0xde, 0x3d, 0xb9, 0xd6, 0xeb, 0x8a, 0x17, 0x88, 0x88, 0xdd, 0x9b, 0xde,
	0xd4, 0xbf, 0x49, 0xc1, 0xc6, 0x02, 0x4f, 0x70, 0x35, 0xb8, 0x3b, 0x88, 0xeb, 0xcc, 0xf7, 0xcf,
	0xe2, 0x7d, 0x85, 0x1d, 0xf9, 0x1d, 0xe2, 0xfa, 0xc1, 0x55, 0xe3, 0x0a, 0xb0, 0x28, 0x59, 0xbe,
	0xd1, 0x37, 0xa8, 0x1b, 0x3c, 0xd8, 0x88, 0x0b, 0xc5, 0xda, 0x04, 0x17, 0x6f, 0x36, 0xdf, 0x03,
	0xec, 0xd8, 0x9e, 0xe1, 0x1b, 0x4f, 0xa8, 0x6a, 0x58, 0xe1, 0xeb, 0x0e, 0xbb, 0x60, 0xa4, 0x15,
	0x14, 0x4a, 0x1a, 0x96, 0x1f, 0xb1, 0x2d, 0x3a, 0x20, 0x33, 0x6c, 0x56, 0xc0, 0x53, 0x0a, 0x0a,
	0x25, 0x11, 0xfb, 0x12, 0x14, 0x75, 0x7b, 0xcc, 0xba, 0x2e, 0xc1, 0x63, 0xe7, 0x45, 0x42, 0x29,
	0x08, 0x2c, 0xa2, 0x04, 0xfd, 0xf4, 0xe4, 0x59, 0xa9, 0xa8, 0x14, 0x04, 0x26, 0x28, 0x97, 0x61,
	0x8d, 0x0c, 0x06, 0x2e, 0x33, 0x1e, 0x1a, 0x12, 0x37, 0x84, 0x52, 0x04, 0x73, 0xe2, 0xd6, 0x7d,
	0xc8, 0x85, 0x71, 0x60, 0x47, 0x32, 0x8b, 0x84, 0xea, 0x88, 0x6b, 0x6f, 0x72, 0x27, 0xaf, 0xe4,
	0xac, 0x50, 0x78, 0x09, 0x8a, 0x86, 0xa7, 0x4e, 0x5e, 0xc9, 0x93, 0xdb, 0xc9, 0x9d, 0x9c, 0x52,
	0x30, 0xbc, 0xe8, 0x85, 0xb1, 0xfc, 0x75, 0x12, 0x4a, 0xd3, 0xaf, 0xfc, 0xb8, 0x0e, 0x39, 0xd3,
	0xd6, 0x08, 0x4f, 0x2d, 0xf1, 0x13, 0xd3, 0xce, 0x2b, 0x7e, 0x18, 0xa8, 0x34, 0x03, 0xbe, 0x12,
	0x69, 0x6e, 0xfd, 0x2d, 0x01, 0xb9, 0x10, 0xc6, 0xe7, 0x20, 0xed, 0x10, 0x7f, 0xc8, 0xcd, 0x65,
	0x0e, 0x92, 0x28, 0xa1, 0xf0, 0x31, 0xc3, 0x3d, 0x87, 0x58, 0x3c, 0x05, 0x02, 0x9c, 0x8d, 0xd9,
	0xba, 0x9a, 0x94, 0xe8, 0xfc, 0xfa, 0x61, 0x8f, 0x46, 0xd4, 0xf2, 0xbd, 0x70, 0x5d, 0x03, 0xbc,
	0x16, 0xc0, 0xf8, 0x1a, 0xac, 0xfb, 0x2e, 0x31, 0xcc, 0x29, 0x6e, 0x9a, 0x73, 0x51, 0x28, 0x88,
	0xc8, 0xfb, 0x70, 0x21, 0xb4, 0xab, 0x53, 0x9f, 0x68, 0x43, 0xaa, 0x4f, 0x94, 0xb2, 0xfc, 0x99,
	0xe1, 0x7c, 0x40, 0xa8, 0x07, 0xf2, 0x50, 0xb7, 0xfc, 0xf7, 0x04, 0xac, 0x87, 0x17, 0x26, 0x3d,
	0x0a, 0xd6, 0x11, 0x00, 0xb1, 0x2c, 0xdb, 0x8f, 0x87, 0x6b, 0x3e, 0x95, 0xe7, 0xf4, 0x2a, 0xd5,
	0x48, 0x49, 0x89, 0x19, 0xd8, 0x1a, 0x01, 0x4c, 0x24, 0x4b, 0xc3, 0x76, 0x11, 0x0a, 0xc1, 0x4f,
	0x38, 0xfc, 0x77, 0x40, 0x71, 0xc5, 0x06, 0x01, 0xb1, 0x9b, 0x15, 0xde, 0x84, 0xcc, 0x09, 0x1d,
	0x18, 0x56, 0xf0, 0x30, 0x2b, 0x06, 0xe1, 0x43, 0x48, 0x3a, 0x7a, 0x08, 0x39, 0xf8, 0x19, 0x6c,
	0x68, 0xf6, 0x68, 0xd6, 0xdd, 0x03, 0x34, 0x73, 0xcd, 0xf7, 0x3e, 0x4f, 0x7c, 0x09, 0x93, 0x16,
	0xf3, 0x7f, 0x89, 0xc4, 0xef, 0x93, 0xa9, 0xc3, 0xce, 0xc1, 0x1f, 0x92, 0x5b, 0x87, 0x42, 0xb5,
	0x13, 0xce, 0x54, 0xa1, 0x7d, 0x93, 0x6a, 0xcc, 0xfb, 0xff, 0x07, 0x00, 0x00, 0xff, 0xff, 0xa3,
	0x58, 0x22, 0x30, 0xdf, 0x1c, 0x00, 0x00,
}
