












package prometheus












type Collector interface {
	
	
	
	
	
	
	
	
	
	
	
	Describe(chan<- *Desc)
	
	
	
	
	
	
	
	
	
	
	Collect(chan<- Metric)
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
