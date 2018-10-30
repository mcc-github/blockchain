












package prometheus

import (
	"errors"
	"os"

	"github.com/prometheus/procfs"
)

type processCollector struct {
	collectFn       func(chan<- Metric)
	pidFn           func() (int, error)
	reportErrors    bool
	cpuTotal        *Desc
	openFDs, maxFDs *Desc
	vsize, maxVsize *Desc
	rss             *Desc
	startTime       *Desc
}



type ProcessCollectorOpts struct {
	
	
	
	
	PidFn func() (int, error)
	
	
	Namespace string
	
	
	
	
	
	
	
	ReportErrors bool
}






















func NewProcessCollector(opts ProcessCollectorOpts) Collector {
	ns := ""
	if len(opts.Namespace) > 0 {
		ns = opts.Namespace + "_"
	}

	c := &processCollector{
		reportErrors: opts.ReportErrors,
		cpuTotal: NewDesc(
			ns+"process_cpu_seconds_total",
			"Total user and system CPU time spent in seconds.",
			nil, nil,
		),
		openFDs: NewDesc(
			ns+"process_open_fds",
			"Number of open file descriptors.",
			nil, nil,
		),
		maxFDs: NewDesc(
			ns+"process_max_fds",
			"Maximum number of open file descriptors.",
			nil, nil,
		),
		vsize: NewDesc(
			ns+"process_virtual_memory_bytes",
			"Virtual memory size in bytes.",
			nil, nil,
		),
		maxVsize: NewDesc(
			ns+"process_virtual_memory_max_bytes",
			"Maximum amount of virtual memory available in bytes.",
			nil, nil,
		),
		rss: NewDesc(
			ns+"process_resident_memory_bytes",
			"Resident memory size in bytes.",
			nil, nil,
		),
		startTime: NewDesc(
			ns+"process_start_time_seconds",
			"Start time of the process since unix epoch in seconds.",
			nil, nil,
		),
	}

	if opts.PidFn == nil {
		pid := os.Getpid()
		c.pidFn = func() (int, error) { return pid, nil }
	} else {
		c.pidFn = opts.PidFn
	}

	
	if _, err := procfs.NewStat(); err == nil {
		c.collectFn = c.processCollect
	} else {
		c.collectFn = func(ch chan<- Metric) {
			c.reportError(ch, nil, errors.New("process metrics not supported on this platform"))
		}
	}

	return c
}


func (c *processCollector) Describe(ch chan<- *Desc) {
	ch <- c.cpuTotal
	ch <- c.openFDs
	ch <- c.maxFDs
	ch <- c.vsize
	ch <- c.maxVsize
	ch <- c.rss
	ch <- c.startTime
}


func (c *processCollector) Collect(ch chan<- Metric) {
	c.collectFn(ch)
}

func (c *processCollector) processCollect(ch chan<- Metric) {
	pid, err := c.pidFn()
	if err != nil {
		c.reportError(ch, nil, err)
		return
	}

	p, err := procfs.NewProc(pid)
	if err != nil {
		c.reportError(ch, nil, err)
		return
	}

	if stat, err := p.NewStat(); err == nil {
		ch <- MustNewConstMetric(c.cpuTotal, CounterValue, stat.CPUTime())
		ch <- MustNewConstMetric(c.vsize, GaugeValue, float64(stat.VirtualMemory()))
		ch <- MustNewConstMetric(c.rss, GaugeValue, float64(stat.ResidentMemory()))
		if startTime, err := stat.StartTime(); err == nil {
			ch <- MustNewConstMetric(c.startTime, GaugeValue, startTime)
		} else {
			c.reportError(ch, c.startTime, err)
		}
	} else {
		c.reportError(ch, nil, err)
	}

	if fds, err := p.FileDescriptorsLen(); err == nil {
		ch <- MustNewConstMetric(c.openFDs, GaugeValue, float64(fds))
	} else {
		c.reportError(ch, c.openFDs, err)
	}

	if limits, err := p.NewLimits(); err == nil {
		ch <- MustNewConstMetric(c.maxFDs, GaugeValue, float64(limits.OpenFiles))
		ch <- MustNewConstMetric(c.maxVsize, GaugeValue, float64(limits.AddressSpace))
	} else {
		c.reportError(ch, nil, err)
	}
}

func (c *processCollector) reportError(ch chan<- Metric, desc *Desc, err error) {
	if !c.reportErrors {
		return
	}
	if desc == nil {
		desc = NewInvalidDesc(err)
	}
	ch <- NewInvalidMetric(desc, err)
}
