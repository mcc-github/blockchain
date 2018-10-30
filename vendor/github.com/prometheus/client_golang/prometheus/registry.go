












package prometheus

import (
	"bytes"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"sync"
	"unicode/utf8"

	"github.com/golang/protobuf/proto"

	dto "github.com/prometheus/client_model/go"

	"github.com/prometheus/client_golang/prometheus/internal"
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
	MustRegister(NewProcessCollector(ProcessCollectorOpts{}))
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






func Unregister(c Collector) bool {
	return DefaultRegisterer.Unregister(c)
}


type GathererFunc func() ([]*dto.MetricFamily, error)


func (gf GathererFunc) Gather() ([]*dto.MetricFamily, error) {
	return gf()
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


func (errs *MultiError) Append(err error) {
	if err != nil {
		*errs = append(*errs, err)
	}
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
	uncheckedCollectors   []Collector
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
	defer func() {
		
		for range descChan {
		}
		r.mtx.Unlock()
	}()
	
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
		r.uncheckedCollectors = append(r.uncheckedCollectors, c)
		return nil
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
		checkedMetricChan   = make(chan Metric, capMetricChan)
		uncheckedMetricChan = make(chan Metric, capMetricChan)
		metricHashes        = map[uint64]struct{}{}
		wg                  sync.WaitGroup
		errs                MultiError          
		registeredDescIDs   map[uint64]struct{} 
	)

	r.mtx.RLock()
	goroutineBudget := len(r.collectorsByID) + len(r.uncheckedCollectors)
	metricFamiliesByName := make(map[string]*dto.MetricFamily, len(r.dimHashesByName))
	checkedCollectors := make(chan Collector, len(r.collectorsByID))
	uncheckedCollectors := make(chan Collector, len(r.uncheckedCollectors))
	for _, collector := range r.collectorsByID {
		checkedCollectors <- collector
	}
	for _, collector := range r.uncheckedCollectors {
		uncheckedCollectors <- collector
	}
	
	
	if r.pedanticChecksEnabled {
		registeredDescIDs = make(map[uint64]struct{}, len(r.descIDs))
		for id := range r.descIDs {
			registeredDescIDs[id] = struct{}{}
		}
	}
	r.mtx.RUnlock()

	wg.Add(goroutineBudget)

	collectWorker := func() {
		for {
			select {
			case collector := <-checkedCollectors:
				collector.Collect(checkedMetricChan)
			case collector := <-uncheckedCollectors:
				collector.Collect(uncheckedMetricChan)
			default:
				return
			}
			wg.Done()
		}
	}

	
	go collectWorker()
	goroutineBudget--

	
	
	go func() {
		wg.Wait()
		close(checkedMetricChan)
		close(uncheckedMetricChan)
	}()

	
	defer func() {
		if checkedMetricChan != nil {
			for range checkedMetricChan {
			}
		}
		if uncheckedMetricChan != nil {
			for range uncheckedMetricChan {
			}
		}
	}()

	
	
	cmc := checkedMetricChan
	umc := uncheckedMetricChan

	for {
		select {
		case metric, ok := <-cmc:
			if !ok {
				cmc = nil
				break
			}
			errs.Append(processMetric(
				metric, metricFamiliesByName,
				metricHashes,
				registeredDescIDs,
			))
		case metric, ok := <-umc:
			if !ok {
				umc = nil
				break
			}
			errs.Append(processMetric(
				metric, metricFamiliesByName,
				metricHashes,
				nil,
			))
		default:
			if goroutineBudget <= 0 || len(checkedCollectors)+len(uncheckedCollectors) == 0 {
				
				
				
				
				select {
				case metric, ok := <-cmc:
					if !ok {
						cmc = nil
						break
					}
					errs.Append(processMetric(
						metric, metricFamiliesByName,
						metricHashes,
						registeredDescIDs,
					))
				case metric, ok := <-umc:
					if !ok {
						umc = nil
						break
					}
					errs.Append(processMetric(
						metric, metricFamiliesByName,
						metricHashes,
						nil,
					))
				}
				break
			}
			
			go collectWorker()
			goroutineBudget--
			runtime.Gosched()
		}
		
		
		
		if cmc == nil && umc == nil {
			break
		}
	}
	return internal.NormalizeMetricFamilies(metricFamiliesByName), errs.MaybeUnwrap()
}


func processMetric(
	metric Metric,
	metricFamiliesByName map[string]*dto.MetricFamily,
	metricHashes map[uint64]struct{},
	registeredDescIDs map[uint64]struct{},
) error {
	desc := metric.Desc()
	
	
	if desc.err != nil {
		return desc.err
	}
	dtoMetric := &dto.Metric{}
	if err := metric.Write(dtoMetric); err != nil {
		return fmt.Errorf("error collecting metric %v: %s", desc, err)
	}
	metricFamily, ok := metricFamiliesByName[desc.fqName]
	if ok { 
		if metricFamily.GetHelp() != desc.help {
			return fmt.Errorf(
				"collected metric %s %s has help %q but should have %q",
				desc.fqName, dtoMetric, desc.help, metricFamily.GetHelp(),
			)
		}
		
		switch metricFamily.GetType() {
		case dto.MetricType_COUNTER:
			if dtoMetric.Counter == nil {
				return fmt.Errorf(
					"collected metric %s %s should be a Counter",
					desc.fqName, dtoMetric,
				)
			}
		case dto.MetricType_GAUGE:
			if dtoMetric.Gauge == nil {
				return fmt.Errorf(
					"collected metric %s %s should be a Gauge",
					desc.fqName, dtoMetric,
				)
			}
		case dto.MetricType_SUMMARY:
			if dtoMetric.Summary == nil {
				return fmt.Errorf(
					"collected metric %s %s should be a Summary",
					desc.fqName, dtoMetric,
				)
			}
		case dto.MetricType_UNTYPED:
			if dtoMetric.Untyped == nil {
				return fmt.Errorf(
					"collected metric %s %s should be Untyped",
					desc.fqName, dtoMetric,
				)
			}
		case dto.MetricType_HISTOGRAM:
			if dtoMetric.Histogram == nil {
				return fmt.Errorf(
					"collected metric %s %s should be a Histogram",
					desc.fqName, dtoMetric,
				)
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
			return fmt.Errorf("empty metric collected: %s", dtoMetric)
		}
		if err := checkSuffixCollisions(metricFamily, metricFamiliesByName); err != nil {
			return err
		}
		metricFamiliesByName[desc.fqName] = metricFamily
	}
	if err := checkMetricConsistency(metricFamily, dtoMetric, metricHashes); err != nil {
		return err
	}
	if registeredDescIDs != nil {
		
		if _, exist := registeredDescIDs[desc.id]; !exist {
			return fmt.Errorf(
				"collected metric %s %s with unregistered descriptor %s",
				metricFamily.GetName(), dtoMetric, desc,
			)
		}
		if err := checkDescConsistency(metricFamily, dtoMetric, desc); err != nil {
			return err
		}
	}
	metricFamily.Metric = append(metricFamily.Metric, dtoMetric)
	return nil
}

















type Gatherers []Gatherer


func (gs Gatherers) Gather() ([]*dto.MetricFamily, error) {
	var (
		metricFamiliesByName = map[string]*dto.MetricFamily{}
		metricHashes         = map[uint64]struct{}{}
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
				if err := checkSuffixCollisions(existingMF, metricFamiliesByName); err != nil {
					errs = append(errs, err)
					continue
				}
				metricFamiliesByName[mf.GetName()] = existingMF
			}
			for _, m := range mf.Metric {
				if err := checkMetricConsistency(existingMF, m, metricHashes); err != nil {
					errs = append(errs, err)
					continue
				}
				existingMF.Metric = append(existingMF.Metric, m)
			}
		}
	}
	return internal.NormalizeMetricFamilies(metricFamiliesByName), errs.MaybeUnwrap()
}




func checkSuffixCollisions(mf *dto.MetricFamily, mfs map[string]*dto.MetricFamily) error {
	var (
		newName              = mf.GetName()
		newType              = mf.GetType()
		newNameWithoutSuffix = ""
	)
	switch {
	case strings.HasSuffix(newName, "_count"):
		newNameWithoutSuffix = newName[:len(newName)-6]
	case strings.HasSuffix(newName, "_sum"):
		newNameWithoutSuffix = newName[:len(newName)-4]
	case strings.HasSuffix(newName, "_bucket"):
		newNameWithoutSuffix = newName[:len(newName)-7]
	}
	if newNameWithoutSuffix != "" {
		if existingMF, ok := mfs[newNameWithoutSuffix]; ok {
			switch existingMF.GetType() {
			case dto.MetricType_SUMMARY:
				if !strings.HasSuffix(newName, "_bucket") {
					return fmt.Errorf(
						"collected metric named %q collides with previously collected summary named %q",
						newName, newNameWithoutSuffix,
					)
				}
			case dto.MetricType_HISTOGRAM:
				return fmt.Errorf(
					"collected metric named %q collides with previously collected histogram named %q",
					newName, newNameWithoutSuffix,
				)
			}
		}
	}
	if newType == dto.MetricType_SUMMARY || newType == dto.MetricType_HISTOGRAM {
		if _, ok := mfs[newName+"_count"]; ok {
			return fmt.Errorf(
				"collected histogram or summary named %q collides with previously collected metric named %q",
				newName, newName+"_count",
			)
		}
		if _, ok := mfs[newName+"_sum"]; ok {
			return fmt.Errorf(
				"collected histogram or summary named %q collides with previously collected metric named %q",
				newName, newName+"_sum",
			)
		}
	}
	if newType == dto.MetricType_HISTOGRAM {
		if _, ok := mfs[newName+"_bucket"]; ok {
			return fmt.Errorf(
				"collected histogram named %q collides with previously collected metric named %q",
				newName, newName+"_bucket",
			)
		}
	}
	return nil
}





func checkMetricConsistency(
	metricFamily *dto.MetricFamily,
	dtoMetric *dto.Metric,
	metricHashes map[uint64]struct{},
) error {
	name := metricFamily.GetName()

	
	if metricFamily.GetType() == dto.MetricType_GAUGE && dtoMetric.Gauge == nil ||
		metricFamily.GetType() == dto.MetricType_COUNTER && dtoMetric.Counter == nil ||
		metricFamily.GetType() == dto.MetricType_SUMMARY && dtoMetric.Summary == nil ||
		metricFamily.GetType() == dto.MetricType_HISTOGRAM && dtoMetric.Histogram == nil ||
		metricFamily.GetType() == dto.MetricType_UNTYPED && dtoMetric.Untyped == nil {
		return fmt.Errorf(
			"collected metric %q { %s} is not a %s",
			name, dtoMetric, metricFamily.GetType(),
		)
	}

	previousLabelName := ""
	for _, labelPair := range dtoMetric.GetLabel() {
		labelName := labelPair.GetName()
		if labelName == previousLabelName {
			return fmt.Errorf(
				"collected metric %q { %s} has two or more labels with the same name: %s",
				name, dtoMetric, labelName,
			)
		}
		if !checkLabelName(labelName) {
			return fmt.Errorf(
				"collected metric %q { %s} has a label with an invalid name: %s",
				name, dtoMetric, labelName,
			)
		}
		if dtoMetric.Summary != nil && labelName == quantileLabel {
			return fmt.Errorf(
				"collected metric %q { %s} must not have an explicit %q label",
				name, dtoMetric, quantileLabel,
			)
		}
		if !utf8.ValidString(labelPair.GetValue()) {
			return fmt.Errorf(
				"collected metric %q { %s} has a label named %q whose value is not utf8: %#v",
				name, dtoMetric, labelName, labelPair.GetValue())
		}
		previousLabelName = labelName
	}

	
	h := hashNew()
	h = hashAdd(h, name)
	h = hashAddByte(h, separatorByte)
	
	
	sort.Sort(labelPairSorter(dtoMetric.Label))
	for _, lp := range dtoMetric.Label {
		h = hashAdd(h, lp.GetName())
		h = hashAddByte(h, separatorByte)
		h = hashAdd(h, lp.GetValue())
		h = hashAddByte(h, separatorByte)
	}
	if _, exists := metricHashes[h]; exists {
		return fmt.Errorf(
			"collected metric %q { %s} was collected before with the same name and label values",
			name, dtoMetric,
		)
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
	sort.Sort(labelPairSorter(lpsFromDesc))
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
