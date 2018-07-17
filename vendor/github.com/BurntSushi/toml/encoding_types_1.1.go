

package toml






type TextMarshaler interface {
	MarshalText() (text []byte, err error)
}



type TextUnmarshaler interface {
	UnmarshalText(text []byte) error
}
