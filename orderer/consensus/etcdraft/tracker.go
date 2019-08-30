

package etcdraft

import (
	"sync/atomic"

	"github.com/mcc-github/blockchain-protos-go/orderer/etcdraft"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metrics"
	"github.com/mcc-github/blockchain/protoutil"
	"go.etcd.io/etcd/raft"
)



type Tracker struct {
	id     uint64
	sender *Disseminator
	gauge  metrics.Gauge
	active *atomic.Value

	counter int

	logger *flogging.FabricLogger
}

func (t *Tracker) Check(status *raft.Status) {
	
	if status.Lead == raft.None {
		t.gauge.Set(0)
		t.active.Store([]uint64{})
		return
	}

	
	if status.RaftState == raft.StateFollower {
		return
	}

	
	current := []uint64{t.id}
	for id, progress := range status.Progress {

		if id == t.id {
			
			
			
			continue
		}

		if progress.RecentActive {
			current = append(current, id)
		}
	}

	last := t.active.Load().([]uint64)
	t.active.Store(current)

	if len(current) != len(last) {
		t.counter = 0
		return
	}

	
	
	if t.counter < 3 {
		t.counter++
		return
	}

	t.counter = 0
	t.logger.Debugf("Current active nodes in cluster are: %+v", current)

	t.gauge.Set(float64(len(current)))
	metadata := protoutil.MarshalOrPanic(&etcdraft.ClusterMetadata{ActiveNodes: current})
	t.sender.UpdateMetadata(metadata)
}
