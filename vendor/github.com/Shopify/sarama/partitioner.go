package sarama

import (
	"hash"
	"hash/fnv"
	"math/rand"
	"time"
)




type Partitioner interface {
	
	Partition(message *ProducerMessage, numPartitions int32) (int32, error)

	
	
	
	
	
	RequiresConsistency() bool
}


type PartitionerConstructor func(topic string) Partitioner

type manualPartitioner struct{}



func NewManualPartitioner(topic string) Partitioner {
	return new(manualPartitioner)
}

func (p *manualPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) {
	return message.Partition, nil
}

func (p *manualPartitioner) RequiresConsistency() bool {
	return true
}

type randomPartitioner struct {
	generator *rand.Rand
}


func NewRandomPartitioner(topic string) Partitioner {
	p := new(randomPartitioner)
	p.generator = rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	return p
}

func (p *randomPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) {
	return int32(p.generator.Intn(int(numPartitions))), nil
}

func (p *randomPartitioner) RequiresConsistency() bool {
	return false
}

type roundRobinPartitioner struct {
	partition int32
}


func NewRoundRobinPartitioner(topic string) Partitioner {
	return &roundRobinPartitioner{}
}

func (p *roundRobinPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) {
	if p.partition >= numPartitions {
		p.partition = 0
	}
	ret := p.partition
	p.partition++
	return ret, nil
}

func (p *roundRobinPartitioner) RequiresConsistency() bool {
	return false
}

type hashPartitioner struct {
	random Partitioner
	hasher hash.Hash32
}




func NewCustomHashPartitioner(hasher func() hash.Hash32) PartitionerConstructor {
	return func(topic string) Partitioner {
		p := new(hashPartitioner)
		p.random = NewRandomPartitioner(topic)
		p.hasher = hasher()
		return p
	}
}





func NewHashPartitioner(topic string) Partitioner {
	p := new(hashPartitioner)
	p.random = NewRandomPartitioner(topic)
	p.hasher = fnv.New32a()
	return p
}

func (p *hashPartitioner) Partition(message *ProducerMessage, numPartitions int32) (int32, error) {
	if message.Key == nil {
		return p.random.Partition(message, numPartitions)
	}
	bytes, err := message.Key.Encode()
	if err != nil {
		return -1, err
	}
	p.hasher.Reset()
	_, err = p.hasher.Write(bytes)
	if err != nil {
		return -1, err
	}
	partition := int32(p.hasher.Sum32()) % numPartitions
	if partition < 0 {
		partition = -partition
	}
	return partition, nil
}

func (p *hashPartitioner) RequiresConsistency() bool {
	return true
}
