



















package tally



type ObjectPool struct {
	values chan interface{}
	alloc  func() interface{}
}


func NewObjectPool(size int) *ObjectPool {
	return &ObjectPool{
		values: make(chan interface{}, size),
	}
}


func (p *ObjectPool) Init(alloc func() interface{}) {
	p.alloc = alloc

	for i := 0; i < cap(p.values); i++ {
		p.values <- p.alloc()
	}
}


func (p *ObjectPool) Get() interface{} {
	var v interface{}
	select {
	case v = <-p.values:
	default:
		v = p.alloc()
	}
	return v
}


func (p *ObjectPool) Put(obj interface{}) {
	select {
	case p.values <- obj:
	default:
	}
}
