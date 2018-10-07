



























package proto

func NewRequiredNotSetError(field string) *RequiredNotSetError {
	return &RequiredNotSetError{field}
}
