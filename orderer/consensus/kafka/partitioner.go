/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package kafka

import "github.com/Shopify/sarama"

type staticPartitioner struct {
	partitionID int32
}



func newStaticPartitioner(partition int32) sarama.PartitionerConstructor {
	return func(topic string) sarama.Partitioner {
		return &staticPartitioner{partition}
	}
}


func (prt *staticPartitioner) Partition(message *sarama.ProducerMessage, numPartitions int32) (int32, error) {
	return prt.partitionID, nil
}



func (prt *staticPartitioner) RequiresConsistency() bool {
	return true
}
