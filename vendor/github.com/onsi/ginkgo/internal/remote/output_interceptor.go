package remote


type OutputInterceptor interface {
	StartInterceptingOutput() error
	StopInterceptingAndReturnOutput() (string, error)
}
