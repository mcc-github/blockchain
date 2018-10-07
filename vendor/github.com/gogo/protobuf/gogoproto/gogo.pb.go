


package gogoproto 

import proto "github.com/gogo/protobuf/proto"
import fmt "fmt"
import math "math"
import descriptor "github.com/gogo/protobuf/protoc-gen-gogo/descriptor"


var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf





const _ = proto.GoGoProtoPackageIsVersion2 

var E_GoprotoEnumPrefix = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         62001,
	Name:          "gogoproto.goproto_enum_prefix",
	Tag:           "varint,62001,opt,name=goproto_enum_prefix,json=goprotoEnumPrefix",
	Filename:      "gogo.proto",
}

var E_GoprotoEnumStringer = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         62021,
	Name:          "gogoproto.goproto_enum_stringer",
	Tag:           "varint,62021,opt,name=goproto_enum_stringer,json=goprotoEnumStringer",
	Filename:      "gogo.proto",
}

var E_EnumStringer = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         62022,
	Name:          "gogoproto.enum_stringer",
	Tag:           "varint,62022,opt,name=enum_stringer,json=enumStringer",
	Filename:      "gogo.proto",
}

var E_EnumCustomname = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         62023,
	Name:          "gogoproto.enum_customname",
	Tag:           "bytes,62023,opt,name=enum_customname,json=enumCustomname",
	Filename:      "gogo.proto",
}

var E_Enumdecl = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         62024,
	Name:          "gogoproto.enumdecl",
	Tag:           "varint,62024,opt,name=enumdecl",
	Filename:      "gogo.proto",
}

var E_EnumvalueCustomname = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.EnumValueOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         66001,
	Name:          "gogoproto.enumvalue_customname",
	Tag:           "bytes,66001,opt,name=enumvalue_customname,json=enumvalueCustomname",
	Filename:      "gogo.proto",
}

var E_GoprotoGettersAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63001,
	Name:          "gogoproto.goproto_getters_all",
	Tag:           "varint,63001,opt,name=goproto_getters_all,json=goprotoGettersAll",
	Filename:      "gogo.proto",
}

var E_GoprotoEnumPrefixAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63002,
	Name:          "gogoproto.goproto_enum_prefix_all",
	Tag:           "varint,63002,opt,name=goproto_enum_prefix_all,json=goprotoEnumPrefixAll",
	Filename:      "gogo.proto",
}

var E_GoprotoStringerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63003,
	Name:          "gogoproto.goproto_stringer_all",
	Tag:           "varint,63003,opt,name=goproto_stringer_all,json=goprotoStringerAll",
	Filename:      "gogo.proto",
}

var E_VerboseEqualAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63004,
	Name:          "gogoproto.verbose_equal_all",
	Tag:           "varint,63004,opt,name=verbose_equal_all,json=verboseEqualAll",
	Filename:      "gogo.proto",
}

var E_FaceAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63005,
	Name:          "gogoproto.face_all",
	Tag:           "varint,63005,opt,name=face_all,json=faceAll",
	Filename:      "gogo.proto",
}

var E_GostringAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63006,
	Name:          "gogoproto.gostring_all",
	Tag:           "varint,63006,opt,name=gostring_all,json=gostringAll",
	Filename:      "gogo.proto",
}

var E_PopulateAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63007,
	Name:          "gogoproto.populate_all",
	Tag:           "varint,63007,opt,name=populate_all,json=populateAll",
	Filename:      "gogo.proto",
}

var E_StringerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63008,
	Name:          "gogoproto.stringer_all",
	Tag:           "varint,63008,opt,name=stringer_all,json=stringerAll",
	Filename:      "gogo.proto",
}

var E_OnlyoneAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63009,
	Name:          "gogoproto.onlyone_all",
	Tag:           "varint,63009,opt,name=onlyone_all,json=onlyoneAll",
	Filename:      "gogo.proto",
}

var E_EqualAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63013,
	Name:          "gogoproto.equal_all",
	Tag:           "varint,63013,opt,name=equal_all,json=equalAll",
	Filename:      "gogo.proto",
}

var E_DescriptionAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63014,
	Name:          "gogoproto.description_all",
	Tag:           "varint,63014,opt,name=description_all,json=descriptionAll",
	Filename:      "gogo.proto",
}

var E_TestgenAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63015,
	Name:          "gogoproto.testgen_all",
	Tag:           "varint,63015,opt,name=testgen_all,json=testgenAll",
	Filename:      "gogo.proto",
}

var E_BenchgenAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63016,
	Name:          "gogoproto.benchgen_all",
	Tag:           "varint,63016,opt,name=benchgen_all,json=benchgenAll",
	Filename:      "gogo.proto",
}

var E_MarshalerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63017,
	Name:          "gogoproto.marshaler_all",
	Tag:           "varint,63017,opt,name=marshaler_all,json=marshalerAll",
	Filename:      "gogo.proto",
}

var E_UnmarshalerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63018,
	Name:          "gogoproto.unmarshaler_all",
	Tag:           "varint,63018,opt,name=unmarshaler_all,json=unmarshalerAll",
	Filename:      "gogo.proto",
}

var E_StableMarshalerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63019,
	Name:          "gogoproto.stable_marshaler_all",
	Tag:           "varint,63019,opt,name=stable_marshaler_all,json=stableMarshalerAll",
	Filename:      "gogo.proto",
}

var E_SizerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63020,
	Name:          "gogoproto.sizer_all",
	Tag:           "varint,63020,opt,name=sizer_all,json=sizerAll",
	Filename:      "gogo.proto",
}

var E_GoprotoEnumStringerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63021,
	Name:          "gogoproto.goproto_enum_stringer_all",
	Tag:           "varint,63021,opt,name=goproto_enum_stringer_all,json=goprotoEnumStringerAll",
	Filename:      "gogo.proto",
}

var E_EnumStringerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63022,
	Name:          "gogoproto.enum_stringer_all",
	Tag:           "varint,63022,opt,name=enum_stringer_all,json=enumStringerAll",
	Filename:      "gogo.proto",
}

var E_UnsafeMarshalerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63023,
	Name:          "gogoproto.unsafe_marshaler_all",
	Tag:           "varint,63023,opt,name=unsafe_marshaler_all,json=unsafeMarshalerAll",
	Filename:      "gogo.proto",
}

var E_UnsafeUnmarshalerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63024,
	Name:          "gogoproto.unsafe_unmarshaler_all",
	Tag:           "varint,63024,opt,name=unsafe_unmarshaler_all,json=unsafeUnmarshalerAll",
	Filename:      "gogo.proto",
}

var E_GoprotoExtensionsMapAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63025,
	Name:          "gogoproto.goproto_extensions_map_all",
	Tag:           "varint,63025,opt,name=goproto_extensions_map_all,json=goprotoExtensionsMapAll",
	Filename:      "gogo.proto",
}

var E_GoprotoUnrecognizedAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63026,
	Name:          "gogoproto.goproto_unrecognized_all",
	Tag:           "varint,63026,opt,name=goproto_unrecognized_all,json=goprotoUnrecognizedAll",
	Filename:      "gogo.proto",
}

var E_GogoprotoImport = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63027,
	Name:          "gogoproto.gogoproto_import",
	Tag:           "varint,63027,opt,name=gogoproto_import,json=gogoprotoImport",
	Filename:      "gogo.proto",
}

var E_ProtosizerAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63028,
	Name:          "gogoproto.protosizer_all",
	Tag:           "varint,63028,opt,name=protosizer_all,json=protosizerAll",
	Filename:      "gogo.proto",
}

var E_CompareAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63029,
	Name:          "gogoproto.compare_all",
	Tag:           "varint,63029,opt,name=compare_all,json=compareAll",
	Filename:      "gogo.proto",
}

var E_TypedeclAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63030,
	Name:          "gogoproto.typedecl_all",
	Tag:           "varint,63030,opt,name=typedecl_all,json=typedeclAll",
	Filename:      "gogo.proto",
}

var E_EnumdeclAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63031,
	Name:          "gogoproto.enumdecl_all",
	Tag:           "varint,63031,opt,name=enumdecl_all,json=enumdeclAll",
	Filename:      "gogo.proto",
}

var E_GoprotoRegistration = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63032,
	Name:          "gogoproto.goproto_registration",
	Tag:           "varint,63032,opt,name=goproto_registration,json=goprotoRegistration",
	Filename:      "gogo.proto",
}

var E_MessagenameAll = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FileOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         63033,
	Name:          "gogoproto.messagename_all",
	Tag:           "varint,63033,opt,name=messagename_all,json=messagenameAll",
	Filename:      "gogo.proto",
}

var E_GoprotoGetters = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64001,
	Name:          "gogoproto.goproto_getters",
	Tag:           "varint,64001,opt,name=goproto_getters,json=goprotoGetters",
	Filename:      "gogo.proto",
}

var E_GoprotoStringer = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64003,
	Name:          "gogoproto.goproto_stringer",
	Tag:           "varint,64003,opt,name=goproto_stringer,json=goprotoStringer",
	Filename:      "gogo.proto",
}

var E_VerboseEqual = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64004,
	Name:          "gogoproto.verbose_equal",
	Tag:           "varint,64004,opt,name=verbose_equal,json=verboseEqual",
	Filename:      "gogo.proto",
}

var E_Face = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64005,
	Name:          "gogoproto.face",
	Tag:           "varint,64005,opt,name=face",
	Filename:      "gogo.proto",
}

var E_Gostring = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64006,
	Name:          "gogoproto.gostring",
	Tag:           "varint,64006,opt,name=gostring",
	Filename:      "gogo.proto",
}

var E_Populate = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64007,
	Name:          "gogoproto.populate",
	Tag:           "varint,64007,opt,name=populate",
	Filename:      "gogo.proto",
}

var E_Stringer = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         67008,
	Name:          "gogoproto.stringer",
	Tag:           "varint,67008,opt,name=stringer",
	Filename:      "gogo.proto",
}

var E_Onlyone = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64009,
	Name:          "gogoproto.onlyone",
	Tag:           "varint,64009,opt,name=onlyone",
	Filename:      "gogo.proto",
}

var E_Equal = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64013,
	Name:          "gogoproto.equal",
	Tag:           "varint,64013,opt,name=equal",
	Filename:      "gogo.proto",
}

var E_Description = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64014,
	Name:          "gogoproto.description",
	Tag:           "varint,64014,opt,name=description",
	Filename:      "gogo.proto",
}

var E_Testgen = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64015,
	Name:          "gogoproto.testgen",
	Tag:           "varint,64015,opt,name=testgen",
	Filename:      "gogo.proto",
}

var E_Benchgen = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64016,
	Name:          "gogoproto.benchgen",
	Tag:           "varint,64016,opt,name=benchgen",
	Filename:      "gogo.proto",
}

var E_Marshaler = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64017,
	Name:          "gogoproto.marshaler",
	Tag:           "varint,64017,opt,name=marshaler",
	Filename:      "gogo.proto",
}

var E_Unmarshaler = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64018,
	Name:          "gogoproto.unmarshaler",
	Tag:           "varint,64018,opt,name=unmarshaler",
	Filename:      "gogo.proto",
}

var E_StableMarshaler = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64019,
	Name:          "gogoproto.stable_marshaler",
	Tag:           "varint,64019,opt,name=stable_marshaler,json=stableMarshaler",
	Filename:      "gogo.proto",
}

var E_Sizer = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64020,
	Name:          "gogoproto.sizer",
	Tag:           "varint,64020,opt,name=sizer",
	Filename:      "gogo.proto",
}

var E_UnsafeMarshaler = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64023,
	Name:          "gogoproto.unsafe_marshaler",
	Tag:           "varint,64023,opt,name=unsafe_marshaler,json=unsafeMarshaler",
	Filename:      "gogo.proto",
}

var E_UnsafeUnmarshaler = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64024,
	Name:          "gogoproto.unsafe_unmarshaler",
	Tag:           "varint,64024,opt,name=unsafe_unmarshaler,json=unsafeUnmarshaler",
	Filename:      "gogo.proto",
}

var E_GoprotoExtensionsMap = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64025,
	Name:          "gogoproto.goproto_extensions_map",
	Tag:           "varint,64025,opt,name=goproto_extensions_map,json=goprotoExtensionsMap",
	Filename:      "gogo.proto",
}

var E_GoprotoUnrecognized = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64026,
	Name:          "gogoproto.goproto_unrecognized",
	Tag:           "varint,64026,opt,name=goproto_unrecognized,json=goprotoUnrecognized",
	Filename:      "gogo.proto",
}

var E_Protosizer = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64028,
	Name:          "gogoproto.protosizer",
	Tag:           "varint,64028,opt,name=protosizer",
	Filename:      "gogo.proto",
}

var E_Compare = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64029,
	Name:          "gogoproto.compare",
	Tag:           "varint,64029,opt,name=compare",
	Filename:      "gogo.proto",
}

var E_Typedecl = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64030,
	Name:          "gogoproto.typedecl",
	Tag:           "varint,64030,opt,name=typedecl",
	Filename:      "gogo.proto",
}

var E_Messagename = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.MessageOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         64033,
	Name:          "gogoproto.messagename",
	Tag:           "varint,64033,opt,name=messagename",
	Filename:      "gogo.proto",
}

var E_Nullable = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         65001,
	Name:          "gogoproto.nullable",
	Tag:           "varint,65001,opt,name=nullable",
	Filename:      "gogo.proto",
}

var E_Embed = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         65002,
	Name:          "gogoproto.embed",
	Tag:           "varint,65002,opt,name=embed",
	Filename:      "gogo.proto",
}

var E_Customtype = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         65003,
	Name:          "gogoproto.customtype",
	Tag:           "bytes,65003,opt,name=customtype",
	Filename:      "gogo.proto",
}

var E_Customname = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         65004,
	Name:          "gogoproto.customname",
	Tag:           "bytes,65004,opt,name=customname",
	Filename:      "gogo.proto",
}

var E_Jsontag = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         65005,
	Name:          "gogoproto.jsontag",
	Tag:           "bytes,65005,opt,name=jsontag",
	Filename:      "gogo.proto",
}

var E_Moretags = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         65006,
	Name:          "gogoproto.moretags",
	Tag:           "bytes,65006,opt,name=moretags",
	Filename:      "gogo.proto",
}

var E_Casttype = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         65007,
	Name:          "gogoproto.casttype",
	Tag:           "bytes,65007,opt,name=casttype",
	Filename:      "gogo.proto",
}

var E_Castkey = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         65008,
	Name:          "gogoproto.castkey",
	Tag:           "bytes,65008,opt,name=castkey",
	Filename:      "gogo.proto",
}

var E_Castvalue = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*string)(nil),
	Field:         65009,
	Name:          "gogoproto.castvalue",
	Tag:           "bytes,65009,opt,name=castvalue",
	Filename:      "gogo.proto",
}

var E_Stdtime = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         65010,
	Name:          "gogoproto.stdtime",
	Tag:           "varint,65010,opt,name=stdtime",
	Filename:      "gogo.proto",
}

var E_Stdduration = &proto.ExtensionDesc{
	ExtendedType:  (*descriptor.FieldOptions)(nil),
	ExtensionType: (*bool)(nil),
	Field:         65011,
	Name:          "gogoproto.stdduration",
	Tag:           "varint,65011,opt,name=stdduration",
	Filename:      "gogo.proto",
}

func init() {
	proto.RegisterExtension(E_GoprotoEnumPrefix)
	proto.RegisterExtension(E_GoprotoEnumStringer)
	proto.RegisterExtension(E_EnumStringer)
	proto.RegisterExtension(E_EnumCustomname)
	proto.RegisterExtension(E_Enumdecl)
	proto.RegisterExtension(E_EnumvalueCustomname)
	proto.RegisterExtension(E_GoprotoGettersAll)
	proto.RegisterExtension(E_GoprotoEnumPrefixAll)
	proto.RegisterExtension(E_GoprotoStringerAll)
	proto.RegisterExtension(E_VerboseEqualAll)
	proto.RegisterExtension(E_FaceAll)
	proto.RegisterExtension(E_GostringAll)
	proto.RegisterExtension(E_PopulateAll)
	proto.RegisterExtension(E_StringerAll)
	proto.RegisterExtension(E_OnlyoneAll)
	proto.RegisterExtension(E_EqualAll)
	proto.RegisterExtension(E_DescriptionAll)
	proto.RegisterExtension(E_TestgenAll)
	proto.RegisterExtension(E_BenchgenAll)
	proto.RegisterExtension(E_MarshalerAll)
	proto.RegisterExtension(E_UnmarshalerAll)
	proto.RegisterExtension(E_StableMarshalerAll)
	proto.RegisterExtension(E_SizerAll)
	proto.RegisterExtension(E_GoprotoEnumStringerAll)
	proto.RegisterExtension(E_EnumStringerAll)
	proto.RegisterExtension(E_UnsafeMarshalerAll)
	proto.RegisterExtension(E_UnsafeUnmarshalerAll)
	proto.RegisterExtension(E_GoprotoExtensionsMapAll)
	proto.RegisterExtension(E_GoprotoUnrecognizedAll)
	proto.RegisterExtension(E_GogoprotoImport)
	proto.RegisterExtension(E_ProtosizerAll)
	proto.RegisterExtension(E_CompareAll)
	proto.RegisterExtension(E_TypedeclAll)
	proto.RegisterExtension(E_EnumdeclAll)
	proto.RegisterExtension(E_GoprotoRegistration)
	proto.RegisterExtension(E_MessagenameAll)
	proto.RegisterExtension(E_GoprotoGetters)
	proto.RegisterExtension(E_GoprotoStringer)
	proto.RegisterExtension(E_VerboseEqual)
	proto.RegisterExtension(E_Face)
	proto.RegisterExtension(E_Gostring)
	proto.RegisterExtension(E_Populate)
	proto.RegisterExtension(E_Stringer)
	proto.RegisterExtension(E_Onlyone)
	proto.RegisterExtension(E_Equal)
	proto.RegisterExtension(E_Description)
	proto.RegisterExtension(E_Testgen)
	proto.RegisterExtension(E_Benchgen)
	proto.RegisterExtension(E_Marshaler)
	proto.RegisterExtension(E_Unmarshaler)
	proto.RegisterExtension(E_StableMarshaler)
	proto.RegisterExtension(E_Sizer)
	proto.RegisterExtension(E_UnsafeMarshaler)
	proto.RegisterExtension(E_UnsafeUnmarshaler)
	proto.RegisterExtension(E_GoprotoExtensionsMap)
	proto.RegisterExtension(E_GoprotoUnrecognized)
	proto.RegisterExtension(E_Protosizer)
	proto.RegisterExtension(E_Compare)
	proto.RegisterExtension(E_Typedecl)
	proto.RegisterExtension(E_Messagename)
	proto.RegisterExtension(E_Nullable)
	proto.RegisterExtension(E_Embed)
	proto.RegisterExtension(E_Customtype)
	proto.RegisterExtension(E_Customname)
	proto.RegisterExtension(E_Jsontag)
	proto.RegisterExtension(E_Moretags)
	proto.RegisterExtension(E_Casttype)
	proto.RegisterExtension(E_Castkey)
	proto.RegisterExtension(E_Castvalue)
	proto.RegisterExtension(E_Stdtime)
	proto.RegisterExtension(E_Stdduration)
}

func init() { proto.RegisterFile("gogo.proto", fileDescriptor_gogo_68790841c0f79064) }

var fileDescriptor_gogo_68790841c0f79064 = []byte{
	
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x98, 0x49, 0x6f, 0x1c, 0x45,
	0x14, 0x80, 0x85, 0x70, 0x64, 0xcf, 0xf3, 0x86, 0xc7, 0xc6, 0x84, 0x08, 0x44, 0xe0, 0xc4, 0xc9,
	0x3e, 0x45, 0x28, 0x65, 0x45, 0x96, 0x63, 0x39, 0x56, 0x10, 0x06, 0x63, 0xe2, 0xb0, 0x1d, 0x46,
	0x3d, 0x33, 0xe5, 0x76, 0x43, 0x77, 0xd7, 0xd0, 0x5d, 0x1d, 0xc5, 0xb9, 0xa1, 0xb0, 0x08, 0x21,
	0x76, 0x24, 0x48, 0x48, 0x02, 0x39, 0xb0, 0xaf, 0x61, 0xe7, 0xc6, 0x85, 0xe5, 0xca, 0x7f, 0xe0,
	0x02, 0x98, 0xdd, 0x37, 0x5f, 0xa2, 0xd7, 0xfd, 0x5e, 0x4f, 0xcd, 0x78, 0xa4, 0xaa, 0xb9, 0xb5,
	0xed, 0xfa, 0x3e, 0x57, 0xbf, 0x57, 0xf5, 0xde, 0x9b, 0x01, 0xf0, 0x95, 0xaf, 0x66, 0x5a, 0x89,
	0xd2, 0xaa, 0x5a, 0xc1, 0xe7, 0xfc, 0xf1, 0xc0, 0x41, 0x5f, 0x29, 0x3f, 0x94, 0xb3, 0xf9, 0x4f,
	0xf5, 0x6c, 0x63, 0xb6, 0x29, 0xd3, 0x46, 0x12, 0xb4, 0xb4, 0x4a, 0x8a, 0xc5, 0xe2, 0x6e, 0x98,
	0xa4, 0xc5, 0x35, 0x19, 0x67, 0x51, 0xad, 0x95, 0xc8, 0x8d, 0xe0, 0x74, 0xf5, 0xa6, 0x99, 0x82,
	0x9c, 0x61, 0x72, 0x66, 0x29, 0xce, 0xa2, 0x7b, 0x5a, 0x3a, 0x50, 0x71, 0xba, 0xff, 0xca, 0xaf,
	0xd7, 0x1e, 0xbc, 0xe6, 0xf6, 0xa1, 0xb5, 0x09, 0x42, 0xf1, 0x6f, 0xab, 0x39, 0x28, 0xd6, 0xe0,
	0xfa, 0x0e, 0x5f, 0xaa, 0x93, 0x20, 0xf6, 0x65, 0x62, 0x31, 0xfe, 0x40, 0xc6, 0x49, 0xc3, 0x78,
	0x1f, 0xa1, 0x62, 0x11, 0x46, 0xfb, 0x71, 0xfd, 0x48, 0xae, 0x11, 0x69, 0x4a, 0x96, 0x61, 0x3c,
	0x97, 0x34, 0xb2, 0x54, 0xab, 0x28, 0xf6, 0x22, 0x69, 0xd1, 0xfc, 0x94, 0x6b, 0x2a, 0x6b, 0x63,
	0x88, 0x2d, 0x96, 0x94, 0x10, 0x30, 0x84, 0xbf, 0x69, 0xca, 0x46, 0x68, 0x31, 0xfc, 0x4c, 0x1b,
	0x29, 0xd7, 0x8b, 0x93, 0x30, 0x85, 0xcf, 0xa7, 0xbc, 0x30, 0x93, 0xe6, 0x4e, 0x6e, 0xed, 0xe9,
	0x39, 0x89, 0xcb, 0x58, 0xf6, 0xcb, 0xd9, 0x81, 0x7c, 0x3b, 0x93, 0xa5, 0xc0, 0xd8, 0x93, 0x91,
	0x45, 0x5f, 0x6a, 0x2d, 0x93, 0xb4, 0xe6, 0x85, 0xbd, 0xb6, 0x77, 0x2c, 0x08, 0x4b, 0xe3, 0xb9,
	0xed, 0xce, 0x2c, 0x2e, 0x17, 0xe4, 0x42, 0x18, 0x8a, 0x75, 0xb8, 0xa1, 0xc7, 0xa9, 0x70, 0x70,
	0x9e, 0x27, 0xe7, 0xd4, 0x9e, 0x93, 0x81, 0xda, 0x55, 0xe0, 0xdf, 0x97, 0xb9, 0x74, 0x70, 0xbe,
	0x41, 0xce, 0x2a, 0xb1, 0x9c, 0x52, 0x34, 0xde, 0x09, 0x13, 0xa7, 0x64, 0x52, 0x57, 0xa9, 0xac,
	0xc9, 0xc7, 0x32, 0x2f, 0x74, 0xd0, 0x5d, 0x20, 0xdd, 0x38, 0x81, 0x4b, 0xc8, 0xa1, 0xeb, 0x30,
	0x0c, 0x6d, 0x78, 0x0d, 0xe9, 0xa0, 0xb8, 0x48, 0x8a, 0x41, 0x5c, 0x8f, 0xe8, 0x02, 0x8c, 0xf8,
	0xaa, 0x78, 0x25, 0x07, 0xfc, 0x12, 0xe1, 0xc3, 0xcc, 0x90, 0xa2, 0xa5, 0x5a, 0x59, 0xe8, 0x69,
	0x97, 0x1d, 0xbc, 0xc9, 0x0a, 0x66, 0x48, 0xd1, 0x47, 0x58, 0xdf, 0x62, 0x45, 0x6a, 0xc4, 0x73,
	0x1e, 0x86, 0x55, 0x1c, 0x6e, 0xa9, 0xd8, 0x65, 0x13, 0x97, 0xc9, 0x00, 0x84, 0xa0, 0x60, 0x0e,
	0x2a, 0xae, 0x89, 0x78, 0x7b, 0x9b, 0xaf, 0x07, 0x67, 0x60, 0x19, 0xc6, 0xb9, 0x40, 0x05, 0x2a,
	0x76, 0x50, 0xbc, 0x43, 0x8a, 0x31, 0x03, 0xa3, 0xd7, 0xd0, 0x32, 0xd5, 0xbe, 0x74, 0x91, 0xbc,
	0xcb, 0xaf, 0x41, 0x08, 0x85, 0xb2, 0x2e, 0xe3, 0xc6, 0xa6, 0x9b, 0xe1, 0x3d, 0x0e, 0x25, 0x33,
	0xa8, 0x58, 0x84, 0xd1, 0xc8, 0x4b, 0xd2, 0x4d, 0x2f, 0x74, 0x4a, 0xc7, 0xfb, 0xe4, 0x18, 0x29,
	0x21, 0x8a, 0x48, 0x16, 0xf7, 0xa3, 0xf9, 0x80, 0x23, 0x62, 0x60, 0x74, 0xf5, 0x52, 0xed, 0xd5,
	0x43, 0x59, 0xeb, 0xc7, 0xf6, 0x21, 0x5f, 0xbd, 0x82, 0x5d, 0x31, 0x8d, 0x73, 0x50, 0x49, 0x83,
	0x33, 0x4e, 0x9a, 0x8f, 0x38, 0xd3, 0x39, 0x80, 0xf0, 0x83, 0x70, 0x63, 0xcf, 0x36, 0xe1, 0x20,
	0xfb, 0x98, 0x64, 0xd3, 0x3d, 0x5a, 0x05, 0x95, 0x84, 0x7e, 0x95, 0x9f, 0x70, 0x49, 0x90, 0x5d,
	0xae, 0x55, 0x98, 0xca, 0xe2, 0xd4, 0xdb, 0xe8, 0x2f, 0x6a, 0x9f, 0x72, 0xd4, 0x0a, 0xb6, 0x23,
	0x6a, 0x27, 0x60, 0x9a, 0x8c, 0xfd, 0xe5, 0xf5, 0x33, 0x2e, 0xac, 0x05, 0xbd, 0xde, 0x99, 0xdd,
	0x87, 0xe1, 0x40, 0x19, 0xce, 0xd3, 0x5a, 0xc6, 0x29, 0x32, 0xb5, 0xc8, 0x6b, 0x39, 0x98, 0xaf,
	0x90, 0x99, 0x2b, 0xfe, 0x52, 0x29, 0x58, 0xf1, 0x5a, 0x28, 0x7f, 0x00, 0xf6, 0xb3, 0x3c, 0x8b,
	0x13, 0xd9, 0x50, 0x7e, 0x1c, 0x9c, 0x91, 0x4d, 0x07, 0xf5, 0xe7, 0x5d, 0xa9, 0x5a, 0x37, 0x70,
	0x34, 0x1f, 0x87, 0xeb, 0xca, 0x59, 0xa5, 0x16, 0x44, 0x2d, 0x95, 0x68, 0x8b, 0xf1, 0x0b, 0xce,
	0x54, 0xc9, 0x1d, 0xcf, 0x31, 0xb1, 0x04, 0x63, 0xf9, 0x8f, 0xae, 0x47, 0xf2, 0x4b, 0x12, 0x8d,
	0xb6, 0x29, 0x2a, 0x1c, 0x0d, 0x15, 0xb5, 0xbc, 0xc4, 0xa5, 0xfe, 0x7d, 0xc5, 0x85, 0x83, 0x10,
	0x2a, 0x1c, 0x7a, 0xab, 0x25, 0xb1, 0xdb, 0x3b, 0x18, 0xbe, 0xe6, 0xc2, 0xc1, 0x0c, 0x29, 0x78,
	0x60, 0x70, 0x50, 0x7c, 0xc3, 0x0a, 0x66, 0x50, 0x71, 0x6f, 0xbb, 0xd1, 0x26, 0xd2, 0x0f, 0x52,
	0x9d, 0x78, 0xb8, 0xda, 0xa2, 0xfa, 0x76, 0xbb, 0x73, 0x08, 0x5b, 0x33, 0x50, 0xac, 0x44, 0x91,
	0x4c, 0x53, 0xcf, 0x97, 0x38, 0x71, 0x38, 0x6c, 0xec, 0x3b, 0xae, 0x44, 0x06, 0x56, 0xdc, 0xcf,
	0xf1, 0xae, 0x59, 0xa5, 0x7a, 0xcb, 0x1e, 0xd1, 0x4a, 0xc1, 0xb0, 0xeb, 0xf1, 0x1d, 0x72, 0x75,
	0x8e, 0x2a, 0xe2, 0x2e, 0x3c, 0x40, 0x9d, 0x03, 0x85, 0x5d, 0x76, 0x76, 0xa7, 0x3c, 0x43, 0x1d,
	0xf3, 0x84, 0x38, 0x06, 0xa3, 0x1d, 0xc3, 0x84, 0x5d, 0xf5, 0x04, 0xa9, 0x46, 0xcc, 0x59, 0x42,
	0x1c, 0x82, 0x01, 0x1c, 0x0c, 0xec, 0xf8, 0x93, 0x84, 0xe7, 0xcb, 0xc5, 0x11, 0x18, 0xe2, 0x81,
	0xc0, 0x8e, 0x3e, 0x45, 0x68, 0x89, 0x20, 0xce, 0xc3, 0x80, 0x1d, 0x7f, 0x9a, 0x71, 0x46, 0x10,
	0x77, 0x0f, 0xe1, 0xf7, 0xcf, 0x0e, 0x50, 0x41, 0xe7, 0xd8, 0xcd, 0xc1, 0x20, 0x4d, 0x01, 0x76,
	0xfa, 0x19, 0xfa, 0xe7, 0x4c, 0x88, 0x3b, 0x60, 0x9f, 0x63, 0xc0, 0x9f, 0x23, 0xb4, 0x58, 0x2f,
	0x16, 0x61, 0xd8, 0xe8, 0xfc, 0x76, 0xfc, 0x79, 0xc2, 0x4d, 0x0a, 0xb7, 0x4e, 0x9d, 0xdf, 0x2e,
	0x78, 0x81, 0xb7, 0x4e, 0x04, 0x86, 0x8d, 0x9b, 0xbe, 0x9d, 0x7e, 0x91, 0xa3, 0xce, 0x88, 0x98,
	0x87, 0x4a, 0x59, 0xc8, 0xed, 0xfc, 0x4b, 0xc4, 0xb7, 0x19, 0x8c, 0x80, 0xd1, 0x48, 0xec, 0x8a,
	0x97, 0x39, 0x02, 0x06, 0x85, 0xd7, 0xa8, 0x7b, 0x38, 0xb0, 0x9b, 0x5e, 0xe1, 0x6b, 0xd4, 0x35,
	0x1b, 0x60, 0x36, 0xf3, 0x7a, 0x6a, 0x57, 0xbc, 0xca, 0xd9, 0xcc, 0xd7, 0xe3, 0x36, 0xba, 0xbb,
	0xad, 0xdd, 0xf1, 0x1a, 0x6f, 0xa3, 0xab, 0xd9, 0x8a, 0x55, 0xa8, 0xee, 0xed, 0xb4, 0x76, 0xdf,
	0xeb, 0xe4, 0x9b, 0xd8, 0xd3, 0x68, 0xc5, 0xfd, 0x30, 0xdd, 0xbb, 0xcb, 0xda, 0xad, 0xe7, 0x76,
	0xba, 0x3e, 0x17, 0x99, 0x4d, 0x56, 0x9c, 0x68, 0x97, 0x6b, 0xb3, 0xc3, 0xda, 0xb5, 0xe7, 0x77,
	0x3a, 0x2b, 0xb6, 0xd9, 0x60, 0xc5, 0x02, 0x40, 0xbb, 0xb9, 0xd9, 0x5d, 0x17, 0xc8, 0x65, 0x40,
	0x78, 0x35, 0xa8, 0xb7, 0xd9, 0xf9, 0x8b, 0x7c, 0x35, 0x88, 0xc0, 0xab, 0xc1, 0x6d, 0xcd, 0x4e,
	0x5f, 0xe2, 0xab, 0xc1, 0x08, 0x9e, 0x6c, 0xa3, 0x73, 0xd8, 0x0d, 0x97, 0xf9, 0x64, 0x1b, 0x94,
	0x98, 0x83, 0xa1, 0x38, 0x0b, 0x43, 0x3c, 0xa0, 0xd5, 0x9b, 0x7b, 0xb4, 0x2b, 0x19, 0x36, 0x99,
	0xff, 0x6d, 0x97, 0x76, 0xc0, 0x80, 0x38, 0x04, 0xfb, 0x64, 0x54, 0x97, 0x4d, 0x1b, 0xf9, 0xfb,
	0x2e, 0x17, 0x25, 0x5c, 0x2d, 0xe6, 0x01, 0x8a, 0x8f, 0xf6, 0xf8, 0x2a, 0x36, 0xf6, 0x8f, 0xdd,
	0xe2, 0x5b, 0x06, 0x03, 0x69, 0x0b, 0xf2, 0x17, 0xb7, 0x08, 0xb6, 0x3b, 0x05, 0xf9, 0x5b, 0x1f,
	0x86, 0xc1, 0x47, 0x52, 0x15, 0x6b, 0xcf, 0xb7, 0xd1, 0x7f, 0x12, 0xcd, 0xeb, 0x31, 0x60, 0x91,
	0x4a, 0xa4, 0xf6, 0xfc, 0xd4, 0xc6, 0xfe, 0x45, 0x6c, 0x09, 0x20, 0xdc, 0xf0, 0x52, 0xed, 0xf2,
	0xde, 0x7f, 0x33, 0xcc, 0x00, 0x6e, 0x1a, 0x9f, 0x1f, 0x95, 0x5b, 0x36, 0xf6, 0x1f, 0xde, 0x34,
	0xad, 0x17, 0x47, 0xa0, 0x82, 0x8f, 0xf9, 0xb7, 0x22, 0x36, 0xf8, 0x5f, 0x82, 0xdb, 0x04, 0xfe,
	0xe7, 0x54, 0x37, 0x75, 0x60, 0x0f, 0xf6, 0x7f, 0x94, 0x69, 0x5e, 0x2f, 0x16, 0x60, 0x38, 0xd5,
	0xcd, 0x66, 0x46, 0xf3, 0x95, 0x05, 0xff, 0x7f, 0xb7, 0xfc, 0xc8, 0x5d, 0x32, 0x47, 0x97, 0x60,
	0xb2, 0xa1, 0xa2, 0x6e, 0xf0, 0x28, 0x2c, 0xab, 0x65, 0xb5, 0x9a, 0x5f, 0xc5, 0x87, 0x6e, 0xf3,
	0x03, 0xbd, 0x99, 0xd5, 0x67, 0x1a, 0x2a, 0x9a, 0xc5, 0xc1, 0xb7, 0xfd, 0x7d, 0x5e, 0x39, 0x06,
	0x5f, 0x0d, 0x00, 0x00, 0xff, 0xff, 0x51, 0xf0, 0xa5, 0x95, 0x02, 0x14, 0x00, 0x00,
}
