












package prometheus


type UntypedOpts Opts






type UntypedFunc interface {
	Metric
	Collector
}







func NewUntypedFunc(opts UntypedOpts, function func() float64) UntypedFunc {
	return newValueFunc(NewDesc(
		BuildFQName(opts.Namespace, opts.Subsystem, opts.Name),
		opts.Help,
		nil,
		opts.ConstLabels,
	), UntypedValue, function)
}
