












package prometheus

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"sort"
	"sync"

	"github.com/golang/protobuf/proto"

	dto "github.com/prometheus/client_model/go"
)

const (
	
	capMetricChan = 1000
	capDescChan   = 10
)











var (
	defaultRegistry              = NewRegistry()
	DefaultRegisterer Registerer = defaultRegistry
	DefaultGatherer   Gatherer   = defaultRegistry
)

func init() {
	MustRegister(NewProcessCollector(os.Getpid(), ""))
	MustRegister(NewGoCollector())
}



func NewRegistry() *Registry {
	return &Registry{
		collectorsByID:  map[uint64]Collector{},
		descIDs:         map[uint64]struct{}{},
		dimHashesByName: map[string]uint64{},
	}
}










func NewPedanticRegistry() *Registry {
	r := NewRegistry()
	r.pedanticChecksEnabled = true
	return r
}






type Registerer interface {
	
	
	
	
	
	
	
	
	
	
	
	
	
	Register(Collector) error
	
	
	
	MustRegister(...Collector)
	
	
	
	
	
	
	
	
	
	
	
	Unregister(Collector) bool
}





type Gatherer interface {
	
	
	
	
	
	
	
	
	
	
	
	Gather() ([]*dto.MetricFamily, error)
}





func Register(c Collector) error {
	return DefaultRegisterer.Register(c)
}






func MustRegister(cs ...Collector) {
	DefaultRegisterer.MustRegister(cs...)
}










func RegisterOrGet(c Collector) (Collector, error) {
	if err := Register(c); err != nil {
		if are, ok := err.(AlreadyRegisteredError); ok {
			return are.ExistingCollector, nil
		}
		return nil, err
	}
	return c, nil
}






func MustRegisterOrGet(c Collector) Collector {
	c, err := RegisterOrGet(c)
	if err != nil {
		panic(err)
	}
	return c
}






func Unregister(c Collector) bool {
	return DefaultRegisterer.Unregister(c)
}


type GathererFunc func() ([]*dto.MetricFamily, error)


func (gf GathererFunc) Gather() ([]*dto.MetricFamily, error) {
	return gf()
}













func SetMetricFamilyInjectionHook(hook func() []*dto.MetricFamily) {
	DefaultGatherer = Gatherers{
		DefaultGatherer,
		GathererFunc(func() ([]*dto.MetricFamily, error) { return hook(), nil }),
	}
}









type AlreadyRegisteredError struct {
	ExistingCollector, NewCollector Collector
}

func (err AlreadyRegisteredError) Error() string {
	return "duplicate metrics collector registration attempted"
}



type MultiError []error

func (errs MultiError) Error() string {
	if len(errs) == 0 {
		return ""
	}
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, "%d error(s) occurred:", len(errs))
	for _, err := range errs {
		fmt.Fprintf(buf, "\n* %s", err)
	}
	return buf.String()
}





func (errs MultiError) MaybeUnwrap() error {
	switch len(errs) {
	case 0:
		return nil
	case 1:
		return errs[0]
	default:
		return errs
	}
}





type Registry struct {
	mtx                   sync.RWMutex
	collectorsByID        map[uint64]Collector 
	descIDs               map[uint64]struct{}
	dimHashesByName       map[string]uint64
	pedanticChecksEnabled bool
}


func (r *Registry) Register(c Collector) error {
	var (
		descChan           = make(chan *Desc, capDescChan)
		newDescIDs         = map[uint64]struct{}{}
		newDimHashesByName = map[string]uint64{}
		collectorID        uint64 
		duplicateDescErr   error
	)
	go func() {
		c.Describe(descChan)
		close(descChan)
	}()
	r.mtx.Lock()
	defer r.mtx.Unlock()
	
	for desc := range descChan {

		
		if desc.err != nil {
			return fmt.Errorf("descriptor %s is invalid: %s", desc, desc.err)
		}

		
		
		if _, exists := r.descIDs[desc.id]; exists {
			duplicateDescErr = fmt.Errorf("descriptor %s already exists with the same fully-qualified name and const label values", desc)
		}
		
		
		
		if _, exists := newDescIDs[desc.id]; !exists {
			newDescIDs[desc.id] = struct{}{}
			collectorID += desc.id
		}

		
		
		
		if dimHash, exists := r.dimHashesByName[desc.fqName]; exists {
			if dimHash != desc.dimHash {
				return fmt.Errorf("a previously registered descriptor with the same fully-qualified name as %s has different label names or a different help string", desc)
			}
		} else {
			
			if dimHash, exists := newDimHashesByName[desc.fqName]; exists {
				if dimHash != desc.dimHash {
					return fmt.Errorf("descriptors reported by collector have inconsistent label names or help strings for the same fully-qualified name, offender is %s", desc)
				}
			} else {
				newDimHashesByName[desc.fqName] = desc.dimHash
			}
		}
	}
	
	if len(newDescIDs) == 0 {
		return errors.New("collector has no descriptors")
	}
	if existing, exists := r.collectorsByID[collectorID]; exists {
		return AlreadyRegisteredError{
			ExistingCollector: existing,
			NewCollector:      c,
		}
	}
	
	
	if duplicateDescErr != nil {
		return duplicateDescErr
	}

	
	r.collectorsByID[collectorID] = c
	for hash := range newDescIDs {
		r.descIDs[hash] = struct{}{}
	}
	for name, dimHash := range newDimHashesByName {
		r.dimHashesByName[name] = dimHash
	}
	return nil
}


func (r *Registry) Unregister(c Collector) bool {
	var (
		descChan    = make(chan *Desc, capDescChan)
		descIDs     = map[uint64]struct{}{}
		collectorID uint64 
	)
	go func() {
		c.Describe(descChan)
		close(descChan)
	}()
	for desc := range descChan {
		if _, exists := descIDs[desc.id]; !exists {
			collectorID += desc.id
			descIDs[desc.id] = struct{}{}
		}
	}

	r.mtx.RLock()
	if _, exists := r.collectorsByID[collectorID]; !exists {
		r.mtx.RUnlock()
		return false
	}
	r.mtx.RUnlock()

	r.mtx.Lock()
	defer r.mtx.Unlock()

	delete(r.collectorsByID, collectorID)
	for id := range descIDs {
		delete(r.descIDs, id)
	}
	
	
	return true
}


func (r *Registry) MustRegister(cs ...Collector) {
	for _, c := range cs {
		if err := r.Register(c); err != nil {
			panic(err)
		}
	}
}


func (r *Registry) Gather() ([]*dto.MetricFamily, error) {
	var (
		metricChan        = make(chan Metric, capMetricChan)
		metricHashes      = map[uint64]struct{}{}
		dimHashes         = map[string]uint64{}
		wg                sync.WaitGroup
		errs              MultiError          
		registeredDescIDs map[uint64]struct{} 
	)

	r.mtx.RLock()
	metricFamiliesByName := make(map[string]*dto.MetricFamily, len(r.dimHashesByName))

	
	
	wg.Add(len(r.collectorsByID))
	go func() {
		wg.Wait()
		close(metricChan)
	}()
	for _, collector := range r.collectorsByID {
		go func(collector Collector) {
			defer wg.Done()
			collector.Collect(metricChan)
		}(collector)
	}

	
	
	if r.pedanticChecksEnabled {
		registeredDescIDs = make(map[uint64]struct{}, len(r.descIDs))
		for id := range r.descIDs {
			registeredDescIDs[id] = struct{}{}
		}
	}

	r.mtx.RUnlock()

	
	defer func() {
		for _ = range metricChan {
		}
	}()

	
	for metric := range metricChan {
		
		
		
		desc := metric.Desc()
		dtoMetric := &dto.Metric{}
		if err := metric.Write(dtoMetric); err != nil {
			errs = append(errs, fmt.Errorf(
				"error collecting metric %v: %s", desc, err,
			))
			continue
		}
		metricFamily, ok := metricFamiliesByName[desc.fqName]
		if ok {
			if metricFamily.GetHelp() != desc.help {
				errs = append(errs, fmt.Errorf(
					"collected metric %s %s has help %q but should have %q",
					desc.fqName, dtoMetric, desc.help, metricFamily.GetHelp(),
				))
				continue
			}
			
			switch metricFamily.GetType() {
			case dto.MetricType_COUNTER:
				if dtoMetric.Counter == nil {
					errs = append(errs, fmt.Errorf(
						"collected metric %s %s should be a Counter",
						desc.fqName, dtoMetric,
					))
					continue
				}
			case dto.MetricType_GAUGE:
				if dtoMetric.Gauge == nil {
					errs = append(errs, fmt.Errorf(
						"collected metric %s %s should be a Gauge",
						desc.fqName, dtoMetric,
					))
					continue
				}
			case dto.MetricType_SUMMARY:
				if dtoMetric.Summary == nil {
					errs = append(errs, fmt.Errorf(
						"collected metric %s %s should be a Summary",
						desc.fqName, dtoMetric,
					))
					continue
				}
			case dto.MetricType_UNTYPED:
				if dtoMetric.Untyped == nil {
					errs = append(errs, fmt.Errorf(
						"collected metric %s %s should be Untyped",
						desc.fqName, dtoMetric,
					))
					continue
				}
			case dto.MetricType_HISTOGRAM:
				if dtoMetric.Histogram == nil {
					errs = append(errs, fmt.Errorf(
						"collected metric %s %s should be a Histogram",
						desc.fqName, dtoMetric,
					))
					continue
				}
			default:
				panic("encountered MetricFamily with invalid type")
			}
		} else {
			metricFamily = &dto.MetricFamily{}
			metricFamily.Name = proto.String(desc.fqName)
			metricFamily.Help = proto.String(desc.help)
			
			switch {
			case dtoMetric.Gauge != nil:
				metricFamily.Type = dto.MetricType_GAUGE.Enum()
			case dtoMetric.Counter != nil:
				metricFamily.Type = dto.MetricType_COUNTER.Enum()
			case dtoMetric.Summary != nil:
				metricFamily.Type = dto.MetricType_SUMMARY.Enum()
			case dtoMetric.Untyped != nil:
				metricFamily.Type = dto.MetricType_UNTYPED.Enum()
			case dtoMetric.Histogram != nil:
				metricFamily.Type = dto.MetricType_HISTOGRAM.Enum()
			default:
				errs = append(errs, fmt.Errorf(
					"empty metric collected: %s", dtoMetric,
				))
				continue
			}
			metricFamiliesByName[desc.fqName] = metricFamily
		}
		if err := checkMetricConsistency(metricFamily, dtoMetric, metricHashes, dimHashes); err != nil {
			errs = append(errs, err)
			continue
		}
		if r.pedanticChecksEnabled {
			
			if _, exist := registeredDescIDs[desc.id]; !exist {
				errs = append(errs, fmt.Errorf(
					"collected metric %s %s with unregistered descriptor %s",
					metricFamily.GetName(), dtoMetric, desc,
				))
				continue
			}
			if err := checkDescConsistency(metricFamily, dtoMetric, desc); err != nil {
				errs = append(errs, err)
				continue
			}
		}
		metricFamily.Metric = append(metricFamily.Metric, dtoMetric)
	}
	return normalizeMetricFamilies(metricFamiliesByName), errs.MaybeUnwrap()
}

















type Gatherers []Gatherer


func (gs Gatherers) Gather() ([]*dto.MetricFamily, error) {
	var (
		metricFamiliesByName = map[string]*dto.MetricFamily{}
		metricHashes         = map[uint64]struct{}{}
		dimHashes            = map[string]uint64{}
		errs                 MultiError 
	)

	for i, g := range gs {
		mfs, err := g.Gather()
		if err != nil {
			if multiErr, ok := err.(MultiError); ok {
				for _, err := range multiErr {
					errs = append(errs, fmt.Errorf("[from Gatherer #%d] %s", i+1, err))
				}
			} else {
				errs = append(errs, fmt.Errorf("[from Gatherer #%d] %s", i+1, err))
			}
		}
		for _, mf := range mfs {
			existingMF, exists := metricFamiliesByName[mf.GetName()]
			if exists {
				if existingMF.GetHelp() != mf.GetHelp() {
					errs = append(errs, fmt.Errorf(
						"gathered metric family %s has help %q but should have %q",
						mf.GetName(), mf.GetHelp(), existingMF.GetHelp(),
					))
					continue
				}
				if existingMF.GetType() != mf.GetType() {
					errs = append(errs, fmt.Errorf(
						"gathered metric family %s has type %s but should have %s",
						mf.GetName(), mf.GetType(), existingMF.GetType(),
					))
					continue
				}
			} else {
				existingMF = &dto.MetricFamily{}
				existingMF.Name = mf.Name
				existingMF.Help = mf.Help
				existingMF.Type = mf.Type
				metricFamiliesByName[mf.GetName()] = existingMF
			}
			for _, m := range mf.Metric {
				if err := checkMetricConsistency(existingMF, m, metricHashes, dimHashes); err != nil {
					errs = append(errs, err)
					continue
				}
				existingMF.Metric = append(existingMF.Metric, m)
			}
		}
	}
	return normalizeMetricFamilies(metricFamiliesByName), errs.MaybeUnwrap()
}


type metricSorter []*dto.Metric

func (s metricSorter) Len() int {
	return len(s)
}

func (s metricSorter) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s metricSorter) Less(i, j int) bool {
	if len(s[i].Label) != len(s[j].Label) {
		
		
		
		
		
		
		return len(s[i].Label) < len(s[j].Label)
	}
	for n, lp := range s[i].Label {
		vi := lp.GetValue()
		vj := s[j].Label[n].GetValue()
		if vi != vj {
			return vi < vj
		}
	}

	
	
	
	
	
	
	if s[i].TimestampMs == nil {
		return false
	}
	if s[j].TimestampMs == nil {
		return true
	}
	return s[i].GetTimestampMs() < s[j].GetTimestampMs()
}




func normalizeMetricFamilies(metricFamiliesByName map[string]*dto.MetricFamily) []*dto.MetricFamily {
	for _, mf := range metricFamiliesByName {
		sort.Sort(metricSorter(mf.Metric))
	}
	names := make([]string, 0, len(metricFamiliesByName))
	for name, mf := range metricFamiliesByName {
		if len(mf.Metric) > 0 {
			names = append(names, name)
		}
	}
	sort.Strings(names)
	result := make([]*dto.MetricFamily, 0, len(names))
	for _, name := range names {
		result = append(result, metricFamiliesByName[name])
	}
	return result
}









func checkMetricConsistency(
	metricFamily *dto.MetricFamily,
	dtoMetric *dto.Metric,
	metricHashes map[uint64]struct{},
	dimHashes map[string]uint64,
) error {
	
	if metricFamily.GetType() == dto.MetricType_GAUGE && dtoMetric.Gauge == nil ||
		metricFamily.GetType() == dto.MetricType_COUNTER && dtoMetric.Counter == nil ||
		metricFamily.GetType() == dto.MetricType_SUMMARY && dtoMetric.Summary == nil ||
		metricFamily.GetType() == dto.MetricType_HISTOGRAM && dtoMetric.Histogram == nil ||
		metricFamily.GetType() == dto.MetricType_UNTYPED && dtoMetric.Untyped == nil {
		return fmt.Errorf(
			"collected metric %s %s is not a %s",
			metricFamily.GetName(), dtoMetric, metricFamily.GetType(),
		)
	}

	
	h := hashNew()
	h = hashAdd(h, metricFamily.GetName())
	h = hashAddByte(h, separatorByte)
	dh := hashNew()
	
	
	sort.Sort(LabelPairSorter(dtoMetric.Label))
	for _, lp := range dtoMetric.Label {
		h = hashAdd(h, lp.GetValue())
		h = hashAddByte(h, separatorByte)
		dh = hashAdd(dh, lp.GetName())
		dh = hashAddByte(dh, separatorByte)
	}
	if _, exists := metricHashes[h]; exists {
		return fmt.Errorf(
			"collected metric %s %s was collected before with the same name and label values",
			metricFamily.GetName(), dtoMetric,
		)
	}
	if dimHash, ok := dimHashes[metricFamily.GetName()]; ok {
		if dimHash != dh {
			return fmt.Errorf(
				"collected metric %s %s has label dimensions inconsistent with previously collected metrics in the same metric family",
				metricFamily.GetName(), dtoMetric,
			)
		}
	} else {
		dimHashes[metricFamily.GetName()] = dh
	}
	metricHashes[h] = struct{}{}
	return nil
}

func checkDescConsistency(
	metricFamily *dto.MetricFamily,
	dtoMetric *dto.Metric,
	desc *Desc,
) error {
	
	if metricFamily.GetHelp() != desc.help {
		return fmt.Errorf(
			"collected metric %s %s has help %q but should have %q",
			metricFamily.GetName(), dtoMetric, metricFamily.GetHelp(), desc.help,
		)
	}

	
	lpsFromDesc := make([]*dto.LabelPair, 0, len(dtoMetric.Label))
	lpsFromDesc = append(lpsFromDesc, desc.constLabelPairs...)
	for _, l := range desc.variableLabels {
		lpsFromDesc = append(lpsFromDesc, &dto.LabelPair{
			Name: proto.String(l),
		})
	}
	if len(lpsFromDesc) != len(dtoMetric.Label) {
		return fmt.Errorf(
			"labels in collected metric %s %s are inconsistent with descriptor %s",
			metricFamily.GetName(), dtoMetric, desc,
		)
	}
	sort.Sort(LabelPairSorter(lpsFromDesc))
	for i, lpFromDesc := range lpsFromDesc {
		lpFromMetric := dtoMetric.Label[i]
		if lpFromDesc.GetName() != lpFromMetric.GetName() ||
			lpFromDesc.Value != nil && lpFromDesc.GetValue() != lpFromMetric.GetValue() {
			return fmt.Errorf(
				"labels in collected metric %s %s are inconsistent with descriptor %s",
				metricFamily.GetName(), dtoMetric, desc,
			)
		}
	}
	return nil
}
