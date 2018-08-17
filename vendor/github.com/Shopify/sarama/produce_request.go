package sarama

import "github.com/rcrowley/go-metrics"






type RequiredAcks int16

const (
	
	NoResponse RequiredAcks = 0
	
	WaitForLocal RequiredAcks = 1
	
	
	
	WaitForAll RequiredAcks = -1
)

type ProduceRequest struct {
	TransactionalID *string
	RequiredAcks    RequiredAcks
	Timeout         int32
	Version         int16 
	records         map[string]map[int32]Records
}

func updateMsgSetMetrics(msgSet *MessageSet, compressionRatioMetric metrics.Histogram,
	topicCompressionRatioMetric metrics.Histogram) int64 {
	var topicRecordCount int64
	for _, messageBlock := range msgSet.Messages {
		
		if messageBlock.Msg.Set != nil {
			topicRecordCount += int64(len(messageBlock.Msg.Set.Messages))
		} else {
			
			topicRecordCount++
		}
		
		if messageBlock.Msg.compressedSize != 0 {
			compressionRatio := float64(len(messageBlock.Msg.Value)) /
				float64(messageBlock.Msg.compressedSize)
			
			intCompressionRatio := int64(100 * compressionRatio)
			compressionRatioMetric.Update(intCompressionRatio)
			topicCompressionRatioMetric.Update(intCompressionRatio)
		}
	}
	return topicRecordCount
}

func updateBatchMetrics(recordBatch *RecordBatch, compressionRatioMetric metrics.Histogram,
	topicCompressionRatioMetric metrics.Histogram) int64 {
	if recordBatch.compressedRecords != nil {
		compressionRatio := int64(float64(recordBatch.recordsLen) / float64(len(recordBatch.compressedRecords)) * 100)
		compressionRatioMetric.Update(compressionRatio)
		topicCompressionRatioMetric.Update(compressionRatio)
	}

	return int64(len(recordBatch.Records))
}

func (r *ProduceRequest) encode(pe packetEncoder) error {
	if r.Version >= 3 {
		if err := pe.putNullableString(r.TransactionalID); err != nil {
			return err
		}
	}
	pe.putInt16(int16(r.RequiredAcks))
	pe.putInt32(r.Timeout)
	metricRegistry := pe.metricRegistry()
	var batchSizeMetric metrics.Histogram
	var compressionRatioMetric metrics.Histogram
	if metricRegistry != nil {
		batchSizeMetric = getOrRegisterHistogram("batch-size", metricRegistry)
		compressionRatioMetric = getOrRegisterHistogram("compression-ratio", metricRegistry)
	}
	totalRecordCount := int64(0)

	err := pe.putArrayLength(len(r.records))
	if err != nil {
		return err
	}

	for topic, partitions := range r.records {
		err = pe.putString(topic)
		if err != nil {
			return err
		}
		err = pe.putArrayLength(len(partitions))
		if err != nil {
			return err
		}
		topicRecordCount := int64(0)
		var topicCompressionRatioMetric metrics.Histogram
		if metricRegistry != nil {
			topicCompressionRatioMetric = getOrRegisterTopicHistogram("compression-ratio", topic, metricRegistry)
		}
		for id, records := range partitions {
			startOffset := pe.offset()
			pe.putInt32(id)
			pe.push(&lengthField{})
			err = records.encode(pe)
			if err != nil {
				return err
			}
			err = pe.pop()
			if err != nil {
				return err
			}
			if metricRegistry != nil {
				if r.Version >= 3 {
					topicRecordCount += updateBatchMetrics(records.RecordBatch, compressionRatioMetric, topicCompressionRatioMetric)
				} else {
					topicRecordCount += updateMsgSetMetrics(records.MsgSet, compressionRatioMetric, topicCompressionRatioMetric)
				}
				batchSize := int64(pe.offset() - startOffset)
				batchSizeMetric.Update(batchSize)
				getOrRegisterTopicHistogram("batch-size", topic, metricRegistry).Update(batchSize)
			}
		}
		if topicRecordCount > 0 {
			getOrRegisterTopicMeter("record-send-rate", topic, metricRegistry).Mark(topicRecordCount)
			getOrRegisterTopicHistogram("records-per-request", topic, metricRegistry).Update(topicRecordCount)
			totalRecordCount += topicRecordCount
		}
	}
	if totalRecordCount > 0 {
		metrics.GetOrRegisterMeter("record-send-rate", metricRegistry).Mark(totalRecordCount)
		getOrRegisterHistogram("records-per-request", metricRegistry).Update(totalRecordCount)
	}

	return nil
}

func (r *ProduceRequest) decode(pd packetDecoder, version int16) error {
	r.Version = version

	if version >= 3 {
		id, err := pd.getNullableString()
		if err != nil {
			return err
		}
		r.TransactionalID = id
	}
	requiredAcks, err := pd.getInt16()
	if err != nil {
		return err
	}
	r.RequiredAcks = RequiredAcks(requiredAcks)
	if r.Timeout, err = pd.getInt32(); err != nil {
		return err
	}
	topicCount, err := pd.getArrayLength()
	if err != nil {
		return err
	}
	if topicCount == 0 {
		return nil
	}

	r.records = make(map[string]map[int32]Records)
	for i := 0; i < topicCount; i++ {
		topic, err := pd.getString()
		if err != nil {
			return err
		}
		partitionCount, err := pd.getArrayLength()
		if err != nil {
			return err
		}
		r.records[topic] = make(map[int32]Records)

		for j := 0; j < partitionCount; j++ {
			partition, err := pd.getInt32()
			if err != nil {
				return err
			}
			size, err := pd.getInt32()
			if err != nil {
				return err
			}
			recordsDecoder, err := pd.getSubset(int(size))
			if err != nil {
				return err
			}
			var records Records
			if err := records.decode(recordsDecoder); err != nil {
				return err
			}
			r.records[topic][partition] = records
		}
	}

	return nil
}

func (r *ProduceRequest) key() int16 {
	return 0
}

func (r *ProduceRequest) version() int16 {
	return r.Version
}

func (r *ProduceRequest) requiredVersion() KafkaVersion {
	switch r.Version {
	case 1:
		return V0_9_0_0
	case 2:
		return V0_10_0_0
	case 3:
		return V0_11_0_0
	default:
		return MinVersion
	}
}

func (r *ProduceRequest) ensureRecords(topic string, partition int32) {
	if r.records == nil {
		r.records = make(map[string]map[int32]Records)
	}

	if r.records[topic] == nil {
		r.records[topic] = make(map[int32]Records)
	}
}

func (r *ProduceRequest) AddMessage(topic string, partition int32, msg *Message) {
	r.ensureRecords(topic, partition)
	set := r.records[topic][partition].MsgSet

	if set == nil {
		set = new(MessageSet)
		r.records[topic][partition] = newLegacyRecords(set)
	}

	set.addMessage(msg)
}

func (r *ProduceRequest) AddSet(topic string, partition int32, set *MessageSet) {
	r.ensureRecords(topic, partition)
	r.records[topic][partition] = newLegacyRecords(set)
}

func (r *ProduceRequest) AddBatch(topic string, partition int32, batch *RecordBatch) {
	r.ensureRecords(topic, partition)
	r.records[topic][partition] = newDefaultRecords(batch)
}
