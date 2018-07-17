package ifrit

import "os"


type Process interface {
	
	Ready() <-chan struct{}

	
	Wait() <-chan error

	
	Signal(os.Signal)
}


func Invoke(r Runner) Process {
	p := Background(r)

	select {
	case <-p.Ready():
	case <-p.Wait():
	}

	return p
}


func Envoke(r Runner) Process {
	return Invoke(r)
}


func Background(r Runner) Process {
	p := newProcess(r)
	go p.run()
	return p
}

type process struct {
	runner     Runner
	signals    chan os.Signal
	ready      chan struct{}
	exited     chan struct{}
	exitStatus error
}

func newProcess(runner Runner) *process {
	return &process{
		runner:  runner,
		signals: make(chan os.Signal),
		ready:   make(chan struct{}),
		exited:  make(chan struct{}),
	}
}

func (p *process) run() {
	p.exitStatus = p.runner.Run(p.signals, p.ready)
	close(p.exited)
}

func (p *process) Ready() <-chan struct{} {
	return p.ready
}

func (p *process) Wait() <-chan error {
	exitChan := make(chan error, 1)

	go func() {
		<-p.exited
		exitChan <- p.exitStatus
	}()

	return exitChan
}

func (p *process) Signal(signal os.Signal) {
	go func() {
		select {
		case p.signals <- signal:
		case <-p.exited:
		}
	}()
}
