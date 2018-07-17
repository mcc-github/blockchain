



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
func (FieldDescriptorProto_Type) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{3, 0} }

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
	return fileDescriptor0, []int{3, 1}
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
func (FileOptions_OptimizeMode) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{9, 0} }

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
func (FieldOptions_CType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{11, 0} }

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
func (FieldOptions_JSType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{11, 1} }




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
	return fileDescriptor0, []int{16, 0}
}



type FileDescriptorSet struct {
	File             []*FileDescriptorProto `protobuf:"bytes,1,rep,name=file" json:"file,omitempty"`
	XXX_unrecognized []byte                 `json:"-"`
}

func (m *FileDescriptorSet) Reset()                    { *m = FileDescriptorSet{} }
func (m *FileDescriptorSet) String() string            { return proto.CompactTextString(m) }
func (*FileDescriptorSet) ProtoMessage()               {}
func (*FileDescriptorSet) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

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
func (*FileDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

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
func (*DescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

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
	Start            *int32 `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	End              *int32 `protobuf:"varint,2,opt,name=end" json:"end,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *DescriptorProto_ExtensionRange) Reset()         { *m = DescriptorProto_ExtensionRange{} }
func (m *DescriptorProto_ExtensionRange) String() string { return proto.CompactTextString(m) }
func (*DescriptorProto_ExtensionRange) ProtoMessage()    {}
func (*DescriptorProto_ExtensionRange) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{2, 0}
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




type DescriptorProto_ReservedRange struct {
	Start            *int32 `protobuf:"varint,1,opt,name=start" json:"start,omitempty"`
	End              *int32 `protobuf:"varint,2,opt,name=end" json:"end,omitempty"`
	XXX_unrecognized []byte `json:"-"`
}

func (m *DescriptorProto_ReservedRange) Reset()         { *m = DescriptorProto_ReservedRange{} }
func (m *DescriptorProto_ReservedRange) String() string { return proto.CompactTextString(m) }
func (*DescriptorProto_ReservedRange) ProtoMessage()    {}
func (*DescriptorProto_ReservedRange) Descriptor() ([]byte, []int) {
	return fileDescriptor0, []int{2, 1}
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
func (*FieldDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

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
func (*OneofDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

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
	Name             *string                     `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Value            []*EnumValueDescriptorProto `protobuf:"bytes,2,rep,name=value" json:"value,omitempty"`
	Options          *EnumOptions                `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
	XXX_unrecognized []byte                      `json:"-"`
}

func (m *EnumDescriptorProto) Reset()                    { *m = EnumDescriptorProto{} }
func (m *EnumDescriptorProto) String() string            { return proto.CompactTextString(m) }
func (*EnumDescriptorProto) ProtoMessage()               {}
func (*EnumDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

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


type EnumValueDescriptorProto struct {
	Name             *string           `protobuf:"bytes,1,opt,name=name" json:"name,omitempty"`
	Number           *int32            `protobuf:"varint,2,opt,name=number" json:"number,omitempty"`
	Options          *EnumValueOptions `protobuf:"bytes,3,opt,name=options" json:"options,omitempty"`
	XXX_unrecognized []byte            `json:"-"`
}

func (m *EnumValueDescriptorProto) Reset()                    { *m = EnumValueDescriptorProto{} }
func (m *EnumValueDescriptorProto) String() string            { return proto.CompactTextString(m) }
func (*EnumValueDescriptorProto) ProtoMessage()               {}
func (*EnumValueDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

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
func (*ServiceDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

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
func (*MethodDescriptorProto) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

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
	
	
	
	
	Deprecated *bool `protobuf:"varint,23,opt,name=deprecated,def=0" json:"deprecated,omitempty"`
	
	
	CcEnableArenas *bool `protobuf:"varint,31,opt,name=cc_enable_arenas,json=ccEnableArenas,def=0" json:"cc_enable_arenas,omitempty"`
	
	
	ObjcClassPrefix *string `protobuf:"bytes,36,opt,name=objc_class_prefix,json=objcClassPrefix" json:"objc_class_prefix,omitempty"`
	
	CsharpNamespace *string `protobuf:"bytes,37,opt,name=csharp_namespace,json=csharpNamespace" json:"csharp_namespace,omitempty"`
	
	
	
	
	SwiftPrefix *string `protobuf:"bytes,39,opt,name=swift_prefix,json=swiftPrefix" json:"swift_prefix,omitempty"`
	
	
	PhpClassPrefix *string `protobuf:"bytes,40,opt,name=php_class_prefix,json=phpClassPrefix" json:"php_class_prefix,omitempty"`
	
	UninterpretedOption          []*UninterpretedOption `protobuf:"bytes,999,rep,name=uninterpreted_option,json=uninterpretedOption" json:"uninterpreted_option,omitempty"`
	proto.XXX_InternalExtensions `json:"-"`
	XXX_unrecognized             []byte `json:"-"`
}

func (m *FileOptions) Reset()                    { *m = FileOptions{} }
func (m *FileOptions) String() string            { return proto.CompactTextString(m) }
func (*FileOptions) ProtoMessage()               {}
func (*FileOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

var extRange_FileOptions = []proto.ExtensionRange{
	{1000, 536870911},
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
func (*MessageOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

var extRange_MessageOptions = []proto.ExtensionRange{
	{1000, 536870911},
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
func (*FieldOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

var extRange_FieldOptions = []proto.ExtensionRange{
	{1000, 536870911},
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
func (*OneofOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

var extRange_OneofOptions = []proto.ExtensionRange{
	{1000, 536870911},
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
func (*EnumOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

var extRange_EnumOptions = []proto.ExtensionRange{
	{1000, 536870911},
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
func (*EnumValueOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

var extRange_EnumValueOptions = []proto.ExtensionRange{
	{1000, 536870911},
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
func (*ServiceOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

var extRange_ServiceOptions = []proto.ExtensionRange{
	{1000, 536870911},
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
func (*MethodOptions) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

var extRange_MethodOptions = []proto.ExtensionRange{
	{1000, 536870911},
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
func (*UninterpretedOption) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

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
	return fileDescriptor0, []int{17, 0}
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
func (*SourceCodeInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

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

func (m *SourceCodeInfo_Location) Reset()                    { *m = SourceCodeInfo_Location{} }
func (m *SourceCodeInfo_Location) String() string            { return proto.CompactTextString(m) }
func (*SourceCodeInfo_Location) ProtoMessage()               {}
func (*SourceCodeInfo_Location) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18, 0} }

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
func (*GeneratedCodeInfo) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

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
	return fileDescriptor0, []int{19, 0}
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
	proto.RegisterType((*FieldDescriptorProto)(nil), "google.protobuf.FieldDescriptorProto")
	proto.RegisterType((*OneofDescriptorProto)(nil), "google.protobuf.OneofDescriptorProto")
	proto.RegisterType((*EnumDescriptorProto)(nil), "google.protobuf.EnumDescriptorProto")
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

func init() { proto.RegisterFile("google/protobuf/descriptor.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xc4, 0x59, 0x5b, 0x6f, 0xdb, 0xc8,
	0x15, 0x5e, 0x5d, 0x2d, 0x1d, 0xc9, 0xf2, 0x78, 0xec, 0x4d, 0x18, 0xef, 0x25, 0x8e, 0xf6, 0x12,
	0x6f, 0xd2, 0xc8, 0x0b, 0xe7, 0xb2, 0x59, 0xa7, 0x48, 0x21, 0x4b, 0x8c, 0x57, 0xa9, 0x2c, 0xa9,
	0x94, 0xdc, 0x4d, 0xf6, 0x85, 0x18, 0x93, 0x23, 0x99, 0x09, 0x45, 0x72, 0x49, 0x2a, 0x89, 0xf7,
	0x29, 0x40, 0x9f, 0x0a, 0xf4, 0x07, 0x14, 0x45, 0xd1, 0x87, 0x7d, 0x59, 0xa0, 0x3f, 0xa0, 0xcf,
	0xfd, 0x05, 0x05, 0xf6, 0xb9, 0x2f, 0x45, 0x51, 0xa0, 0xfd, 0x07, 0x7d, 0x2d, 0x66, 0x86, 0xa4,
	0x48, 0x5d, 0x12, 0x77, 0x81, 0xec, 0x3e, 0xd9, 0x73, 0xce, 0x77, 0x0e, 0xcf, 0x9c, 0xf9, 0x66,
	0xce, 0x99, 0x11, 0x6c, 0x8f, 0x6c, 0x7b, 0x64, 0xd2, 0x5d, 0xc7, 0xb5, 0x7d, 0xfb, 0x64, 0x32,
	0xdc, 0xd5, 0xa9, 0xa7, 0xb9, 0x86, 0xe3, 0xdb, 0x6e, 0x8d, 0xcb, 0xf0, 0x9a, 0x40, 0xd4, 0x42,
	0x44, 0xf5, 0x08, 0xd6, 0x1f, 0x18, 0x26, 0x6d, 0x46, 0xc0, 0x3e, 0xf5, 0xf1, 0x5d, 0xc8, 0x0e,
	0x0d, 0x93, 0x4a, 0xa9, 0xed, 0xcc, 0x4e, 0x69, 0xef, 0xc3, 0xda, 0x8c, 0x51, 0x2d, 0x69, 0xd1,
	0x63, 0x62, 0x85, 0x5b, 0x54, 0xff, 0x95, 0x85, 0x8d, 0x05, 0x5a, 0x8c, 0x21, 0x6b, 0x91, 0x31,
	0xf3, 0x98, 0xda, 0x29, 0x2a, 0xfc, 0x7f, 0x2c, 0xc1, 0x8a, 0x43, 0xb4, 0xa7, 0x64, 0x44, 0xa5,
	0x34, 0x17, 0x87, 0x43, 0xfc, 0x3e, 0x80, 0x4e, 0x1d, 0x6a, 0xe9, 0xd4, 0xd2, 0xce, 0xa4, 0xcc,
	0x76, 0x66, 0xa7, 0xa8, 0xc4, 0x24, 0xf8, 0x3a, 0xac, 0x3b, 0x93, 0x13, 0xd3, 0xd0, 0xd4, 0x18,
	0x0c, 0xb6, 0x33, 0x3b, 0x39, 0x05, 0x09, 0x45, 0x73, 0x0a, 0xbe, 0x0a, 0x6b, 0xcf, 0x29, 0x79,
	0x1a, 0x87, 0x96, 0x38, 0xb4, 0xc2, 0xc4, 0x31, 0x60, 0x03, 0xca, 0x63, 0xea, 0x79, 0x64, 0x44,
	0x55, 0xff, 0xcc, 0xa1, 0x52, 0x96, 0xcf, 0x7e, 0x7b, 0x6e, 0xf6, 0xb3, 0x33, 0x2f, 0x05, 0x56,
	0x83, 0x33, 0x87, 0xe2, 0x3a, 0x14, 0xa9, 0x35, 0x19, 0x0b, 0x0f, 0xb9, 0x25, 0xf9, 0x93, 0xad,
	0xc9, 0x78, 0xd6, 0x4b, 0x81, 0x99, 0x05, 0x2e, 0x56, 0x3c, 0xea, 0x3e, 0x33, 0x34, 0x2a, 0xe5,
	0xb9, 0x83, 0xab, 0x73, 0x0e, 0xfa, 0x42, 0x3f, 0xeb, 0x23, 0xb4, 0xc3, 0x0d, 0x28, 0xd2, 0x17,
	0x3e, 0xb5, 0x3c, 0xc3, 0xb6, 0xa4, 0x15, 0xee, 0xe4, 0xa3, 0x05, 0xab, 0x48, 0x4d, 0x7d, 0xd6,
	0xc5, 0xd4, 0x0e, 0xdf, 0x81, 0x15, 0xdb, 0xf1, 0x0d, 0xdb, 0xf2, 0xa4, 0xc2, 0x76, 0x6a, 0xa7,
	0xb4, 0xf7, 0xee, 0x42, 0x22, 0x74, 0x05, 0x46, 0x09, 0xc1, 0xb8, 0x05, 0xc8, 0xb3, 0x27, 0xae,
	0x46, 0x55, 0xcd, 0xd6, 0xa9, 0x6a, 0x58, 0x43, 0x5b, 0x2a, 0x72, 0x07, 0x97, 0xe7, 0x27, 0xc2,
	0x81, 0x0d, 0x5b, 0xa7, 0x2d, 0x6b, 0x68, 0x2b, 0x15, 0x2f, 0x31, 0xc6, 0x17, 0x20, 0xef, 0x9d,
	0x59, 0x3e, 0x79, 0x21, 0x95, 0x39, 0x43, 0x82, 0x51, 0xf5, 0xbf, 0x39, 0x58, 0x3b, 0x0f, 0xc5,
	0xee, 0x41, 0x6e, 0xc8, 0x66, 0x29, 0xa5, 0xff, 0x9f, 0x1c, 0x08, 0x9b, 0x64, 0x12, 0xf3, 0x3f,
	0x30, 0x89, 0x75, 0x28, 0x59, 0xd4, 0xf3, 0xa9, 0x2e, 0x18, 0x91, 0x39, 0x27, 0xa7, 0x40, 0x18,
	0xcd, 0x53, 0x2a, 0xfb, 0x83, 0x28, 0xf5, 0x08, 0xd6, 0xa2, 0x90, 0x54, 0x97, 0x58, 0xa3, 0x90,
	0x9b, 0xbb, 0xaf, 0x8b, 0xa4, 0x26, 0x87, 0x76, 0x0a, 0x33, 0x53, 0x2a, 0x34, 0x31, 0xc6, 0x4d,
	0x00, 0xdb, 0xa2, 0xf6, 0x50, 0xd5, 0xa9, 0x66, 0x4a, 0x85, 0x25, 0x59, 0xea, 0x32, 0xc8, 0x5c,
	0x96, 0x6c, 0x21, 0xd5, 0x4c, 0xfc, 0xf9, 0x94, 0x6a, 0x2b, 0x4b, 0x98, 0x72, 0x24, 0x36, 0xd9,
	0x1c, 0xdb, 0x8e, 0xa1, 0xe2, 0x52, 0xc6, 0x7b, 0xaa, 0x07, 0x33, 0x2b, 0xf2, 0x20, 0x6a, 0xaf,
	0x9d, 0x99, 0x12, 0x98, 0x89, 0x89, 0xad, 0xba, 0xf1, 0x21, 0xfe, 0x00, 0x22, 0x81, 0xca, 0x69,
	0x05, 0xfc, 0x14, 0x2a, 0x87, 0xc2, 0x0e, 0x19, 0xd3, 0xad, 0xbb, 0x50, 0x49, 0xa6, 0x07, 0x6f,
	0x42, 0xce, 0xf3, 0x89, 0xeb, 0x73, 0x16, 0xe6, 0x14, 0x31, 0xc0, 0x08, 0x32, 0xd4, 0xd2, 0xf9,
	0x29, 0x97, 0x53, 0xd8, 0xbf, 0x5b, 0x9f, 0xc1, 0x6a, 0xe2, 0xf3, 0xe7, 0x35, 0xac, 0xfe, 0x3e,
	0x0f, 0x9b, 0x8b, 0x38, 0xb7, 0x90, 0xfe, 0x17, 0x20, 0x6f, 0x4d, 0xc6, 0x27, 0xd4, 0x95, 0x32,
	0xdc, 0x43, 0x30, 0xc2, 0x75, 0xc8, 0x99, 0xe4, 0x84, 0x9a, 0x52, 0x76, 0x3b, 0xb5, 0x53, 0xd9,
	0xbb, 0x7e, 0x2e, 0x56, 0xd7, 0xda, 0xcc, 0x44, 0x11, 0x96, 0xf8, 0x3e, 0x64, 0x83, 0x23, 0x8e,
	0x79, 0xb8, 0x76, 0x3e, 0x0f, 0x8c, 0x8b, 0x0a, 0xb7, 0xc3, 0xef, 0x40, 0x91, 0xfd, 0x15, 0xb9,
	0xcd, 0xf3, 0x98, 0x0b, 0x4c, 0xc0, 0xf2, 0x8a, 0xb7, 0xa0, 0xc0, 0x69, 0xa6, 0xd3, 0xb0, 0x34,
	0x44, 0x63, 0xb6, 0x30, 0x3a, 0x1d, 0x92, 0x89, 0xe9, 0xab, 0xcf, 0x88, 0x39, 0xa1, 0x9c, 0x30,
	0x45, 0xa5, 0x1c, 0x08, 0x7f, 0xcd, 0x64, 0xf8, 0x32, 0x94, 0x04, 0x2b, 0x0d, 0x4b, 0xa7, 0x2f,
	0xf8, 0xe9, 0x93, 0x53, 0x04, 0x51, 0x5b, 0x4c, 0xc2, 0x3e, 0xff, 0xc4, 0xb3, 0xad, 0x70, 0x69,
	0xf9, 0x27, 0x98, 0x80, 0x7f, 0xfe, 0xb3, 0xd9, 0x83, 0xef, 0xbd, 0xc5, 0xd3, 0x9b, 0xe5, 0x62,
	0xf5, 0x2f, 0x69, 0xc8, 0xf2, 0xfd, 0xb6, 0x06, 0xa5, 0xc1, 0xe3, 0x9e, 0xac, 0x36, 0xbb, 0xc7,
	0x07, 0x6d, 0x19, 0xa5, 0x70, 0x05, 0x80, 0x0b, 0x1e, 0xb4, 0xbb, 0xf5, 0x01, 0x4a, 0x47, 0xe3,
	0x56, 0x67, 0x70, 0xe7, 0x16, 0xca, 0x44, 0x06, 0xc7, 0x42, 0x90, 0x8d, 0x03, 0x6e, 0xee, 0xa1,
	0x1c, 0x46, 0x50, 0x16, 0x0e, 0x5a, 0x8f, 0xe4, 0xe6, 0x9d, 0x5b, 0x28, 0x9f, 0x94, 0xdc, 0xdc,
	0x43, 0x2b, 0x78, 0x15, 0x8a, 0x5c, 0x72, 0xd0, 0xed, 0xb6, 0x51, 0x21, 0xf2, 0xd9, 0x1f, 0x28,
	0xad, 0xce, 0x21, 0x2a, 0x46, 0x3e, 0x0f, 0x95, 0xee, 0x71, 0x0f, 0x41, 0xe4, 0xe1, 0x48, 0xee,
	0xf7, 0xeb, 0x87, 0x32, 0x2a, 0x45, 0x88, 0x83, 0xc7, 0x03, 0xb9, 0x8f, 0xca, 0x89, 0xb0, 0x6e,
	0xee, 0xa1, 0xd5, 0xe8, 0x13, 0x72, 0xe7, 0xf8, 0x08, 0x55, 0xf0, 0x3a, 0xac, 0x8a, 0x4f, 0x84,
	0x41, 0xac, 0xcd, 0x88, 0xee, 0xdc, 0x42, 0x68, 0x1a, 0x88, 0xf0, 0xb2, 0x9e, 0x10, 0xdc, 0xb9,
	0x85, 0x70, 0xb5, 0x01, 0x39, 0xce, 0x2e, 0x8c, 0xa1, 0xd2, 0xae, 0x1f, 0xc8, 0x6d, 0xb5, 0xdb,
	0x1b, 0xb4, 0xba, 0x9d, 0x7a, 0x1b, 0xa5, 0xa6, 0x32, 0x45, 0xfe, 0xd5, 0x71, 0x4b, 0x91, 0x9b,
	0x28, 0x1d, 0x97, 0xf5, 0xe4, 0xfa, 0x40, 0x6e, 0xa2, 0x4c, 0x55, 0x83, 0xcd, 0x45, 0xe7, 0xcc,
	0xc2, 0x9d, 0x11, 0x5b, 0xe2, 0xf4, 0x92, 0x25, 0xe6, 0xbe, 0xe6, 0x96, 0xf8, 0xdb, 0x14, 0x6c,
	0x2c, 0x38, 0x6b, 0x17, 0x7e, 0xe4, 0x17, 0x90, 0x13, 0x14, 0x15, 0xd5, 0xe7, 0x93, 0x85, 0x87,
	0x36, 0x27, 0xec, 0x5c, 0x05, 0xe2, 0x76, 0xf1, 0x0a, 0x9c, 0x59, 0x52, 0x81, 0x99, 0x8b, 0xb9,
	0x20, 0x7f, 0x93, 0x02, 0x69, 0x99, 0xef, 0xd7, 0x1c, 0x14, 0xe9, 0xc4, 0x41, 0x71, 0x6f, 0x36,
	0x80, 0x2b, 0xcb, 0xe7, 0x30, 0x17, 0xc5, 0x77, 0x29, 0xb8, 0xb0, 0xb8, 0x51, 0x59, 0x18, 0xc3,
	0x7d, 0xc8, 0x8f, 0xa9, 0x7f, 0x6a, 0x87, 0xc5, 0xfa, 0xe3, 0x05, 0x25, 0x80, 0xa9, 0x67, 0x73,
	0x15, 0x58, 0xc5, 0x6b, 0x48, 0x66, 0x59, 0xb7, 0x21, 0xa2, 0x99, 0x8b, 0xf4, 0xb7, 0x69, 0x78,
	0x7b, 0xa1, 0xf3, 0x85, 0x81, 0xbe, 0x07, 0x60, 0x58, 0xce, 0xc4, 0x17, 0x05, 0x59, 0x9c, 0x4f,
	0x45, 0x2e, 0xe1, 0x7b, 0x9f, 0x9d, 0x3d, 0x13, 0x3f, 0xd2, 0x67, 0xb8, 0x1e, 0x84, 0x88, 0x03,
	0xee, 0x4e, 0x03, 0xcd, 0xf2, 0x40, 0xdf, 0x5f, 0x32, 0xd3, 0xb9, 0x5a, 0xf7, 0x29, 0x20, 0xcd,
	0x34, 0xa8, 0xe5, 0xab, 0x9e, 0xef, 0x52, 0x32, 0x36, 0xac, 0x11, 0x3f, 0x80, 0x0b, 0xfb, 0xb9,
	0x21, 0x31, 0x3d, 0xaa, 0xac, 0x09, 0x75, 0x3f, 0xd4, 0x32, 0x0b, 0x5e, 0x65, 0xdc, 0x98, 0x45,
	0x3e, 0x61, 0x21, 0xd4, 0x91, 0x45, 0xf5, 0xef, 0x2b, 0x50, 0x8a, 0xb5, 0x75, 0xf8, 0x0a, 0x94,
	0x9f, 0x90, 0x67, 0x44, 0x0d, 0x5b, 0x75, 0x91, 0x89, 0x12, 0x93, 0xf5, 0x82, 0x76, 0xfd, 0x53,
	0xd8, 0xe4, 0x10, 0x7b, 0xe2, 0x53, 0x57, 0xd5, 0x4c, 0xe2, 0x79, 0x3c, 0x69, 0x05, 0x0e, 0xc5,
	0x4c, 0xd7, 0x65, 0xaa, 0x46, 0xa8, 0xc1, 0xb7, 0x61, 0x83, 0x5b, 0x8c, 0x27, 0xa6, 0x6f, 0x38,
	0x26, 0x55, 0xd9, 0xe5, 0xc1, 0xe3, 0x07, 0x71, 0x14, 0xd9, 0x3a, 0x43, 0x1c, 0x05, 0x00, 0x16,
	0x91, 0x87, 0x9b, 0xf0, 0x1e, 0x37, 0x1b, 0x51, 0x8b, 0xba, 0xc4, 0xa7, 0x2a, 0xfd, 0x7a, 0x42,
	0x4c, 0x4f, 0x25, 0x96, 0xae, 0x9e, 0x12, 0xef, 0x54, 0xda, 0x64, 0x0e, 0x0e, 0xd2, 0x52, 0x4a,
	0xb9, 0xc4, 0x80, 0x87, 0x01, 0x4e, 0xe6, 0xb0, 0xba, 0xa5, 0x7f, 0x41, 0xbc, 0x53, 0xbc, 0x0f,
	0x17, 0xb8, 0x17, 0xcf, 0x77, 0x0d, 0x6b, 0xa4, 0x6a, 0xa7, 0x54, 0x7b, 0xaa, 0x4e, 0xfc, 0xe1,
	0x5d, 0xe9, 0x9d, 0xf8, 0xf7, 0x79, 0x84, 0x7d, 0x8e, 0x69, 0x30, 0xc8, 0xb1, 0x3f, 0xbc, 0x8b,
	0xfb, 0x50, 0x66, 0x8b, 0x31, 0x36, 0xbe, 0xa1, 0xea, 0xd0, 0x76, 0x79, 0x65, 0xa9, 0x2c, 0xd8,
	0xd9, 0xb1, 0x0c, 0xd6, 0xba, 0x81, 0xc1, 0x91, 0xad, 0xd3, 0xfd, 0x5c, 0xbf, 0x27, 0xcb, 0x4d,
	0xa5, 0x14, 0x7a, 0x79, 0x60, 0xbb, 0x8c, 0x50, 0x23, 0x3b, 0x4a, 0x70, 0x49, 0x10, 0x6a, 0x64,
	0x87, 0xe9, 0xbd, 0x0d, 0x1b, 0x9a, 0x26, 0xe6, 0x6c, 0x68, 0x6a, 0xd0, 0xe2, 0x7b, 0x12, 0x4a,
	0x24, 0x4b, 0xd3, 0x0e, 0x05, 0x20, 0xe0, 0xb8, 0x87, 0x3f, 0x87, 0xb7, 0xa7, 0xc9, 0x8a, 0x1b,
	0xae, 0xcf, 0xcd, 0x72, 0xd6, 0xf4, 0x36, 0x6c, 0x38, 0x67, 0xf3, 0x86, 0x38, 0xf1, 0x45, 0xe7,
	0x6c, 0xd6, 0xec, 0x23, 0x7e, 0x6d, 0x73, 0xa9, 0x46, 0x7c, 0xaa, 0x4b, 0x17, 0xe3, 0xe8, 0x98,
	0x02, 0xef, 0x02, 0xd2, 0x34, 0x95, 0x5a, 0xe4, 0xc4, 0xa4, 0x2a, 0x71, 0xa9, 0x45, 0x3c, 0xe9,
	0x72, 0x1c, 0x5c, 0xd1, 0x34, 0x99, 0x6b, 0xeb, 0x5c, 0x89, 0xaf, 0xc1, 0xba, 0x7d, 0xf2, 0x44,
	0x13, 0xcc, 0x52, 0x1d, 0x97, 0x0e, 0x8d, 0x17, 0xd2, 0x87, 0x3c, 0x4d, 0x6b, 0x4c, 0xc1, 0x79,
	0xd5, 0xe3, 0x62, 0xfc, 0x09, 0x20, 0xcd, 0x3b, 0x25, 0xae, 0xc3, 0x4b, 0xbb, 0xe7, 0x10, 0x8d,
	0x4a, 0x1f, 0x09, 0xa8, 0x90, 0x77, 0x42, 0x31, 0x63, 0xb6, 0xf7, 0xdc, 0x18, 0xfa, 0xa1, 0xc7,
	0xab, 0x82, 0xd9, 0x5c, 0x16, 0x78, 0xdb, 0x01, 0xe4, 0x9c, 0x3a, 0xc9, 0x0f, 0xef, 0x70, 0x58,
	0xc5, 0x39, 0x75, 0xe2, 0xdf, 0x7d, 0x04, 0x9b, 0x13, 0xcb, 0xb0, 0x7c, 0xea, 0x3a, 0x2e, 0x65,
	0xed, 0xbe, 0xd8, 0xb3, 0xd2, 0xbf, 0x57, 0x96, 0x34, 0xec, 0xc7, 0x71, 0xb4, 0xa0, 0x8a, 0xb2,
	0x31, 0x99, 0x17, 0x56, 0xf7, 0xa1, 0x1c, 0x67, 0x10, 0x2e, 0x82, 0xe0, 0x10, 0x4a, 0xb1, 0x6a,
	0xdc, 0xe8, 0x36, 0x59, 0x1d, 0xfd, 0x4a, 0x46, 0x69, 0x56, 0xcf, 0xdb, 0xad, 0x81, 0xac, 0x2a,
	0xc7, 0x9d, 0x41, 0xeb, 0x48, 0x46, 0x99, 0x6b, 0xc5, 0xc2, 0x7f, 0x56, 0xd0, 0xcb, 0x97, 0x2f,
	0x5f, 0xa6, 0x1f, 0x66, 0x0b, 0x1f, 0xa3, 0xab, 0xd5, 0xef, 0xd3, 0x50, 0x49, 0x76, 0xd2, 0xf8,
	0xe7, 0x70, 0x31, 0xbc, 0xf6, 0x7a, 0xd4, 0x57, 0x9f, 0x1b, 0x2e, 0xa7, 0xf6, 0x98, 0x88, 0x5e,
	0x34, 0x5a, 0x95, 0xcd, 0x00, 0xd5, 0xa7, 0xfe, 0x97, 0x86, 0xcb, 0x88, 0x3b, 0x26, 0x3e, 0x6e,
	0xc3, 0x65, 0xcb, 0x56, 0x3d, 0x9f, 0x58, 0x3a, 0x71, 0x75, 0x75, 0xfa, 0xe0, 0xa0, 0x12, 0x4d,
	0xa3, 0x9e, 0x67, 0x8b, 0x92, 0x12, 0x79, 0x79, 0xd7, 0xb2, 0xfb, 0x01, 0x78, 0x7a, 0xd6, 0xd6,
	0x03, 0xe8, 0x0c, 0x83, 0x32, 0xcb, 0x18, 0xf4, 0x0e, 0x14, 0xc7, 0xc4, 0x51, 0xa9, 0xe5, 0xbb,
	0x67, 0xbc, 0xff, 0x2b, 0x28, 0x85, 0x31, 0x71, 0x64, 0x36, 0x7e, 0x73, 0x2b, 0x91, 0xcc, 0x66,
	0x01, 0x15, 0x1f, 0x66, 0x0b, 0x45, 0x04, 0xd5, 0x7f, 0x66, 0xa0, 0x1c, 0xef, 0x07, 0x59, 0x7b,
	0xad, 0xf1, 0xb3, 0x3f, 0xc5, 0x4f, 0x87, 0x0f, 0x5e, 0xd9, 0x3d, 0xd6, 0x1a, 0xac, 0x28, 0xec,
	0xe7, 0x45, 0x97, 0xa6, 0x08, 0x4b, 0x56, 0x90, 0xd9, 0x79, 0x40, 0x45, 0xef, 0x5f, 0x50, 0x82,
	0x11, 0x3e, 0x84, 0xfc, 0x13, 0x8f, 0xfb, 0xce, 0x73, 0xdf, 0x1f, 0xbe, 0xda, 0xf7, 0xc3, 0x3e,
	0x77, 0x5e, 0x7c, 0xd8, 0x57, 0x3b, 0x5d, 0xe5, 0xa8, 0xde, 0x56, 0x02, 0x73, 0x7c, 0x09, 0xb2,
	0x26, 0xf9, 0xe6, 0x2c, 0x59, 0x3e, 0xb8, 0xe8, 0xbc, 0x8b, 0x70, 0x09, 0xb2, 0xcf, 0x29, 0x79,
	0x9a, 0x3c, 0xb4, 0xb9, 0xe8, 0x0d, 0x6e, 0x86, 0x5d, 0xc8, 0xf1, 0x7c, 0x61, 0x80, 0x20, 0x63,
	0xe8, 0x2d, 0x5c, 0x80, 0x6c, 0xa3, 0xab, 0xb0, 0x0d, 0x81, 0xa0, 0x2c, 0xa4, 0x6a, 0xaf, 0x25,
	0x37, 0x64, 0x94, 0xae, 0xde, 0x86, 0xbc, 0x48, 0x02, 0xdb, 0x2c, 0x51, 0x1a, 0xd0, 0x5b, 0xc1,
	0x30, 0xf0, 0x91, 0x0a, 0xb5, 0xc7, 0x47, 0x07, 0xb2, 0x82, 0xd2, 0xc9, 0xa5, 0xce, 0xa2, 0x5c,
	0xd5, 0x83, 0x72, 0xbc, 0x21, 0xfc, 0x51, 0x58, 0x56, 0xfd, 0x6b, 0x0a, 0x4a, 0xb1, 0x06, 0x8f,
	0xb5, 0x16, 0xc4, 0x34, 0xed, 0xe7, 0x2a, 0x31, 0x0d, 0xe2, 0x05, 0xd4, 0x00, 0x2e, 0xaa, 0x33,
	0xc9, 0x79, 0x97, 0xee, 0x47, 0xda, 0x22, 0x39, 0x94, 0xaf, 0xfe, 0x29, 0x05, 0x68, 0xb6, 0x45,
	0x9c, 0x09, 0x33, 0xf5, 0x53, 0x86, 0x59, 0xfd, 0x63, 0x0a, 0x2a, 0xc9, 0xbe, 0x70, 0x26, 0xbc,
	0x2b, 0x3f, 0x69, 0x78, 0xff, 0x48, 0xc3, 0x6a, 0xa2, 0x1b, 0x3c, 0x6f, 0x74, 0x5f, 0xc3, 0xba,
	0xa1, 0xd3, 0xb1, 0x63, 0xfb, 0xd4, 0xd2, 0xce, 0x54, 0x93, 0x3e, 0xa3, 0xa6, 0x54, 0xe5, 0x87,
	0xc6, 0xee, 0xab, 0xfb, 0xcd, 0x5a, 0x6b, 0x6a, 0xd7, 0x66, 0x66, 0xfb, 0x1b, 0xad, 0xa6, 0x7c,
	0xd4, 0xeb, 0x0e, 0xe4, 0x4e, 0xe3, 0xb1, 0x7a, 0xdc, 0xf9, 0x65, 0xa7, 0xfb, 0x65, 0x47, 0x41,
	0xc6, 0x0c, 0xec, 0x0d, 0x6e, 0xfb, 0x1e, 0xa0, 0xd9, 0xa0, 0xf0, 0x45, 0x58, 0x14, 0x16, 0x7a,
	0x0b, 0x6f, 0xc0, 0x5a, 0xa7, 0xab, 0xf6, 0x5b, 0x4d, 0x59, 0x95, 0x1f, 0x3c, 0x90, 0x1b, 0x83,
	0xbe, 0xb8, 0x80, 0x47, 0xe8, 0x41, 0x62, 0x83, 0x57, 0xff, 0x90, 0x81, 0x8d, 0x05, 0x91, 0xe0,
	0x7a, 0xd0, 0xfb, 0x8b, 0xeb, 0xc8, 0x8d, 0xf3, 0x44, 0x5f, 0x63, 0xdd, 0x45, 0x8f, 0xb8, 0x7e,
	0x70, 0x55, 0xf8, 0x04, 0x58, 0x96, 0x2c, 0xdf, 0x18, 0x1a, 0xd4, 0x0d, 0xde, 0x2b, 0xc4, 0x85,
	0x60, 0x6d, 0x2a, 0x17, 0x4f, 0x16, 0x3f, 0x03, 0xec, 0xd8, 0x9e, 0xe1, 0x1b, 0xcf, 0xa8, 0x6a,
	0x58, 0xe1, 0xe3, 0x06, 0xbb, 0x20, 0x64, 0x15, 0x14, 0x6a, 0x5a, 0x96, 0x1f, 0xa1, 0x2d, 0x3a,
	0x22, 0x33, 0x68, 0x76, 0x98, 0x67, 0x14, 0x14, 0x6a, 0x22, 0xf4, 0x15, 0x28, 0xeb, 0xf6, 0x84,
	0xb5, 0x5b, 0x02, 0xc7, 0x6a, 0x47, 0x4a, 0x29, 0x09, 0x59, 0x04, 0x09, 0xfa, 0xe1, 0xe9, 0xab,
	0x4a, 0x59, 0x29, 0x09, 0x99, 0x80, 0x5c, 0x85, 0x35, 0x32, 0x1a, 0xb9, 0xcc, 0x79, 0xe8, 0x48,
	0x74, 0xf8, 0x95, 0x48, 0xcc, 0x81, 0x5b, 0x0f, 0xa1, 0x10, 0xe6, 0x81, 0x95, 0x6a, 0x96, 0x09,
	0xd5, 0x11, 0x6f, 0x5b, 0xe9, 0x9d, 0xa2, 0x52, 0xb0, 0x42, 0xe5, 0x15, 0x28, 0x1b, 0x9e, 0x3a,
	0x7d, 0x64, 0x4d, 0x6f, 0xa7, 0x77, 0x0a, 0x4a, 0xc9, 0xf0, 0xa2, 0x57, 0xb5, 0xea, 0x77, 0x69,
	0xa8, 0x24, 0x1f, 0x89, 0x71, 0x13, 0x0a, 0xa6, 0xad, 0x11, 0x4e, 0x2d, 0xf1, 0x0b, 0xc5, 0xce,
	0x6b, 0xde, 0x95, 0x6b, 0xed, 0x00, 0xaf, 0x44, 0x96, 0x5b, 0x7f, 0x4b, 0x41, 0x21, 0x14, 0xe3,
	0x0b, 0x90, 0x75, 0x88, 0x7f, 0xca, 0xdd, 0xe5, 0x0e, 0xd2, 0x28, 0xa5, 0xf0, 0x31, 0x93, 0x7b,
	0x0e, 0xb1, 0x38, 0x05, 0x02, 0x39, 0x1b, 0xb3, 0x75, 0x35, 0x29, 0xd1, 0xf9, 0xf5, 0xc1, 0x1e,
	0x8f, 0xa9, 0xe5, 0x7b, 0xe1, 0xba, 0x06, 0xf2, 0x46, 0x20, 0xc6, 0xd7, 0x61, 0xdd, 0x77, 0x89,
	0x61, 0x26, 0xb0, 0x59, 0x8e, 0x45, 0xa1, 0x22, 0x02, 0xef, 0xc3, 0xa5, 0xd0, 0xaf, 0x4e, 0x7d,
	0xa2, 0x9d, 0x52, 0x7d, 0x6a, 0x94, 0xe7, 0x2f, 0x90, 0x17, 0x03, 0x40, 0x33, 0xd0, 0x87, 0xb6,
	0xd5, 0xef, 0x53, 0xb0, 0x1e, 0x5e, 0x78, 0xf4, 0x28, 0x59, 0x47, 0x00, 0xc4, 0xb2, 0x6c, 0x3f,
	0x9e, 0xae, 0x79, 0x2a, 0xcf, 0xd9, 0xd5, 0xea, 0x91, 0x91, 0x12, 0x73, 0xb0, 0x35, 0x06, 0x98,
	0x6a, 0x96, 0xa6, 0xed, 0x32, 0x94, 0x82, 0x5f, 0x00, 0xf8, 0xcf, 0x48, 0xe2, 0x8a, 0x0c, 0x42,
	0xc4, 0x6e, 0x46, 0x78, 0x13, 0x72, 0x27, 0x74, 0x64, 0x58, 0xc1, 0xbb, 0xa4, 0x18, 0x84, 0xaf,
	0x9d, 0xd9, 0xe8, 0xb5, 0xf3, 0xe0, 0x77, 0x29, 0xd8, 0xd0, 0xec, 0xf1, 0x6c, 0xbc, 0x07, 0x68,
	0xe6, 0x9e, 0xee, 0x7d, 0x91, 0xfa, 0xea, 0xfe, 0xc8, 0xf0, 0x4f, 0x27, 0x27, 0x35, 0xcd, 0x1e,
	0xef, 0x8e, 0x6c, 0x93, 0x58, 0xa3, 0xe9, 0xef, 0x60, 0xfc, 0x1f, 0xed, 0xc6, 0x88, 0x5a, 0x37,
	0x46, 0x76, 0xec, 0x57, 0xb1, 0x7b, 0xd3, 0x7f, 0xbf, 0x4d, 0x67, 0x0e, 0x7b, 0x07, 0x7f, 0x4e,
	0x6f, 0x1d, 0x8a, 0x6f, 0xf5, 0xc2, 0xdc, 0x28, 0x74, 0x68, 0x52, 0x8d, 0xcd, 0xf7, 0x7f, 0x01,
	0x00, 0x00, 0xff, 0xff, 0x8e, 0x54, 0xe7, 0xef, 0x60, 0x1b, 0x00, 0x00,
}
