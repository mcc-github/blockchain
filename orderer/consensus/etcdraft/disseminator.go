

package etcdraft

import (
	"sync"

	"github.com/mcc-github/blockchain-protos-go/orderer"
)


type Disseminator struct {
	RPC

	l        sync.Mutex
	sent     map[uint64]bool
	metadata []byte
}

func (d *Disseminator) SendConsensus(dest uint64, msg *orderer.ConsensusRequest) error {
	d.l.Lock()
	defer d.l.Unlock()

	if !d.sent[dest] && len(d.metadata) != 0 {
		msg.Metadata = d.metadata
		d.sent[dest] = true
	}

	return d.RPC.SendConsensus(dest, msg)
}

func (d *Disseminator) UpdateMetadata(m []byte) {
	d.l.Lock()
	defer d.l.Unlock()

	d.sent = make(map[uint64]bool)
	d.metadata = m
}
