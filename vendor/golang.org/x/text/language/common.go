

package language




type langAliasType int8

const (
	langDeprecated langAliasType = iota
	langMacro
	langLegacy

	langAliasTypeUnknown langAliasType = -1
)
