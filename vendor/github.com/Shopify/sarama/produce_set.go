package sarama

import (
	"encoding/binary"
	"time"
)

type partitionSet struct {
	msgs          []*ProducerMessage
	recordsToSend Records
	bufferBytes   int
}

type produceSet struct {
	parent *asyncProducer
	msgs   map[string]map[int32]*partitionSet

	bufferBytes int
	bufferCount int
}

func newProduceSet(parent *asyncProducer) *produceSet {
	return &produceSet{
		msgs:   make(map[string]map[int32]*partitionSet),
		parent: parent,
	}
}

func (ps *produceSet) add(msg *ProducerMessage) error {
	var err error
	var key, val []byte

	if msg.Key != nil {
		if key, err = msg.Key.Encode(); err != nil {
			return err
		}
	}

	if msg.Value != nil {
		if val, err = msg.Value.Encode(); err != nil {
			return err
		}
	}

	timestamp := msg.Timestamp
	if msg.Timestamp.IsZero() {
		timestamp = time.Now()
	}

	partitions := ps.msgs[msg.Topic]
	if partitions == nil {
		partitions = make(map[int32]*partitionSet)
		ps.msgs[msg.Topic] = partitions
	}

	var size int

	set := partitions[msg.Partition]
	if set == nil {
		if ps.parent.conf.Version.IsAtLeast(V0_11_0_0) {
			batch := &RecordBatch{
				FirstTimestamp:   timestamp,
				Version:          2,
				ProducerID:       -1, 
				Codec:            ps.parent.conf.Producer.Compression,
				CompressionLevel: ps.parent.conf.Producer.CompressionLevel,
			}
			set = &partitionSet{recordsToSend: newDefaultRecords(batch)}
			size = recordBatchOverhead
		} else {
			set = &partitionSet{recordsToSend: newLegacyRecords(new(MessageSet))}
		}
		partitions[msg.Partition] = set
	}

	set.msgs = append(set.msgs, msg)
	if ps.parent.conf.Version.IsAtLeast(V0_11_0_0) {
		
		size += maximumRecordOverhead
		rec := &Record{
			Key:            key,
			Value:          val,
			TimestampDelta: timestamp.Sub(set.recordsToSend.RecordBatch.FirstTimestamp),
		}
		size += len(key) + len(val)
		if len(msg.Headers) > 0 {
			rec.Headers = make([]*RecordHeader, len(msg.Headers))
			for i := range msg.Headers {
				rec.Headers[i] = &msg.Headers[i]
				size += len(rec.Headers[i].Key) + len(rec.Headers[i].Value) + 2*binary.MaxVarintLen32
			}
		}
		set.recordsToSend.RecordBatch.addRecord(rec)
	} else {
		msgToSend := &Message{Codec: CompressionNone, Key: key, Value: val}
		if ps.parent.conf.Version.IsAtLeast(V0_10_0_0) {
			msgToSend.Timestamp = timestamp
			msgToSend.Version = 1
		}
		set.recordsToSend.MsgSet.addMessage(msgToSend)
		size = producerMessageOverhead + len(key) + len(val)
	}

	set.bufferBytes += size
	ps.bufferBytes += size
	ps.bufferCount++

	return nil
}

func (ps *produceSet) buildRequest() *ProduceRequest {
	req := &ProduceRequest{
		RequiredAcks: ps.parent.conf.Producer.RequiredAcks,
		Timeout:      int32(ps.parent.conf.Producer.Timeout / time.Millisecond),
	}
	if ps.parent.conf.Version.IsAtLeast(V0_10_0_0) {
		req.Version = 2
	}
	if ps.parent.conf.Version.IsAtLeast(V0_11_0_0) {
		req.Version = 3
	}

	for topic, partitionSet := range ps.msgs {
		for partition, set := range partitionSet {
			if req.Version >= 3 {
				
				
				
				
				
				
				
				rb := set.recordsToSend.RecordBatch
				if len(rb.Records) > 0 {
					rb.LastOffsetDelta = int32(len(rb.Records) - 1)
					for i, record := range rb.Records {
						record.OffsetDelta = int64(i)
					}
				}

				req.AddBatch(topic, partition, rb)
				continue
			}
			if ps.parent.conf.Producer.Compression == CompressionNone {
				req.AddSet(topic, partition, set.recordsToSend.MsgSet)
			} else {
				
				
				
				

				if ps.parent.conf.Version.IsAtLeast(V0_10_0_0) {
					
					
					
					
					
					for i, msg := range set.recordsToSend.MsgSet.Messages {
						msg.Offset = int64(i)
					}
				}
				payload, err := encode(set.recordsToSend.MsgSet, ps.parent.conf.MetricRegistry)
				if err != nil {
					Logger.Println(err) 
					panic(err)
				}
				compMsg := &Message{
					Codec:            ps.parent.conf.Producer.Compression,
					CompressionLevel: ps.parent.conf.Producer.CompressionLevel,
					Key:              nil,
					Value:            payload,
					Set:              set.recordsToSend.MsgSet, 
				}
				if ps.parent.conf.Version.IsAtLeast(V0_10_0_0) {
					compMsg.Version = 1
					compMsg.Timestamp = set.recordsToSend.MsgSet.Messages[0].Msg.Timestamp
				}
				req.AddMessage(topic, partition, compMsg)
			}
		}
	}

	return req
}

func (ps *produceSet) eachPartition(cb func(topic string, partition int32, msgs []*ProducerMessage)) {
	for topic, partitionSet := range ps.msgs {
		for partition, set := range partitionSet {
			cb(topic, partition, set.msgs)
		}
	}
}

func (ps *produceSet) dropPartition(topic string, partition int32) []*ProducerMessage {
	if ps.msgs[topic] == nil {
		return nil
	}
	set := ps.msgs[topic][partition]
	if set == nil {
		return nil
	}
	ps.bufferBytes -= set.bufferBytes
	ps.bufferCount -= len(set.msgs)
	delete(ps.msgs[topic], partition)
	return set.msgs
}

func (ps *produceSet) wouldOverflow(msg *ProducerMessage) bool {
	version := 1
	if ps.parent.conf.Version.IsAtLeast(V0_11_0_0) {
		version = 2
	}

	switch {
	
	case ps.bufferBytes+msg.byteSize(version) >= int(MaxRequestSize-(10*1024)):
		return true
	
	case ps.parent.conf.Producer.Compression != CompressionNone &&
		ps.msgs[msg.Topic] != nil && ps.msgs[msg.Topic][msg.Partition] != nil &&
		ps.msgs[msg.Topic][msg.Partition].bufferBytes+msg.byteSize(version) >= ps.parent.conf.Producer.MaxMessageBytes:
		return true
	
	case ps.parent.conf.Producer.Flush.MaxMessages > 0 && ps.bufferCount >= ps.parent.conf.Producer.Flush.MaxMessages:
		return true
	default:
		return false
	}
}

func (ps *produceSet) readyToFlush() bool {
	switch {
	
	case ps.empty():
		return false
	
	case ps.parent.conf.Producer.Flush.Frequency == 0 && ps.parent.conf.Producer.Flush.Bytes == 0 && ps.parent.conf.Producer.Flush.Messages == 0:
		return true
	
	case ps.parent.conf.Producer.Flush.Messages > 0 && ps.bufferCount >= ps.parent.conf.Producer.Flush.Messages:
		return true
	
	case ps.parent.conf.Producer.Flush.Bytes > 0 && ps.bufferBytes >= ps.parent.conf.Producer.Flush.Bytes:
		return true
	default:
		return false
	}
}

func (ps *produceSet) empty() bool {
	return ps.bufferCount == 0
}
