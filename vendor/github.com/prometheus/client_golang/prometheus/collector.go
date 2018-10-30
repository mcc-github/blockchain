












package prometheus












type Collector interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	
	Describe(chan<- *Desc)
	
	
	
	
	
	
	
	
	
	
	
	
	Collect(chan<- Metric)
}























func DescribeByCollect(c Collector, descs chan<- *Desc) {
	metrics := make(chan Metric)
	go func() {
		c.Collect(metrics)
		close(metrics)
	}()
	for m := range metrics {
		descs <- m.Desc()
	}
}




type selfCollector struct {
	self Metric
}




func (c *selfCollector) init(self Metric) {
	c.self = self
}


func (c *selfCollector) Describe(ch chan<- *Desc) {
	ch <- c.self.Desc()
}


func (c *selfCollector) Collect(ch chan<- Metric) {
	ch <- c.self
}
