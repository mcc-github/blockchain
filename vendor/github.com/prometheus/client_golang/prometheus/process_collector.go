












package prometheus

import "github.com/prometheus/procfs"

type processCollector struct {
	pid             int
	collectFn       func(chan<- Metric)
	pidFn           func() (int, error)
	cpuTotal        Counter
	openFDs, maxFDs Gauge
	vsize, rss      Gauge
	startTime       Gauge
}




func NewProcessCollector(pid int, namespace string) Collector {
	return NewProcessCollectorPIDFn(
		func() (int, error) { return pid, nil },
		namespace,
	)
}






func NewProcessCollectorPIDFn(
	pidFn func() (int, error),
	namespace string,
) Collector {
	c := processCollector{
		pidFn:     pidFn,
		collectFn: func(chan<- Metric) {},

		cpuTotal: NewCounter(CounterOpts{
			Namespace: namespace,
			Name:      "process_cpu_seconds_total",
			Help:      "Total user and system CPU time spent in seconds.",
		}),
		openFDs: NewGauge(GaugeOpts{
			Namespace: namespace,
			Name:      "process_open_fds",
			Help:      "Number of open file descriptors.",
		}),
		maxFDs: NewGauge(GaugeOpts{
			Namespace: namespace,
			Name:      "process_max_fds",
			Help:      "Maximum number of open file descriptors.",
		}),
		vsize: NewGauge(GaugeOpts{
			Namespace: namespace,
			Name:      "process_virtual_memory_bytes",
			Help:      "Virtual memory size in bytes.",
		}),
		rss: NewGauge(GaugeOpts{
			Namespace: namespace,
			Name:      "process_resident_memory_bytes",
			Help:      "Resident memory size in bytes.",
		}),
		startTime: NewGauge(GaugeOpts{
			Namespace: namespace,
			Name:      "process_start_time_seconds",
			Help:      "Start time of the process since unix epoch in seconds.",
		}),
	}

	
	if _, err := procfs.NewStat(); err == nil {
		c.collectFn = c.processCollect
	}

	return &c
}


func (c *processCollector) Describe(ch chan<- *Desc) {
	ch <- c.cpuTotal.Desc()
	ch <- c.openFDs.Desc()
	ch <- c.maxFDs.Desc()
	ch <- c.vsize.Desc()
	ch <- c.rss.Desc()
	ch <- c.startTime.Desc()
}


func (c *processCollector) Collect(ch chan<- Metric) {
	c.collectFn(ch)
}



func (c *processCollector) processCollect(ch chan<- Metric) {
	pid, err := c.pidFn()
	if err != nil {
		return
	}

	p, err := procfs.NewProc(pid)
	if err != nil {
		return
	}

	if stat, err := p.NewStat(); err == nil {
		c.cpuTotal.Set(stat.CPUTime())
		ch <- c.cpuTotal
		c.vsize.Set(float64(stat.VirtualMemory()))
		ch <- c.vsize
		c.rss.Set(float64(stat.ResidentMemory()))
		ch <- c.rss

		if startTime, err := stat.StartTime(); err == nil {
			c.startTime.Set(startTime)
			ch <- c.startTime
		}
	}

	if fds, err := p.FileDescriptorsLen(); err == nil {
		c.openFDs.Set(float64(fds))
		ch <- c.openFDs
	}

	if limits, err := p.NewLimits(); err == nil {
		c.maxFDs.Set(float64(limits.OpenFiles))
		ch <- c.maxFDs
	}
}
