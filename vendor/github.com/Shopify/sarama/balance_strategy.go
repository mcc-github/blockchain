package sarama

import (
	"math"
	"sort"
)




type BalanceStrategyPlan map[string]map[string][]int32


func (p BalanceStrategyPlan) Add(memberID, topic string, partitions ...int32) {
	if len(partitions) == 0 {
		return
	}
	if _, ok := p[memberID]; !ok {
		p[memberID] = make(map[string][]int32, 1)
	}
	p[memberID][topic] = append(p[memberID][topic], partitions...)
}





type BalanceStrategy interface {
	
	Name() string

	
	
	Plan(members map[string]ConsumerGroupMemberMetadata, topics map[string][]int32) (BalanceStrategyPlan, error)
}







var BalanceStrategyRange = &balanceStrategy{
	name: "range",
	coreFn: func(plan BalanceStrategyPlan, memberIDs []string, topic string, partitions []int32) {
		step := float64(len(partitions)) / float64(len(memberIDs))

		for i, memberID := range memberIDs {
			pos := float64(i)
			min := int(math.Floor(pos*step + 0.5))
			max := int(math.Floor((pos+1)*step + 0.5))
			plan.Add(memberID, topic, partitions[min:max]...)
		}
	},
}





var BalanceStrategyRoundRobin = &balanceStrategy{
	name: "roundrobin",
	coreFn: func(plan BalanceStrategyPlan, memberIDs []string, topic string, partitions []int32) {
		for i, part := range partitions {
			memberID := memberIDs[i%len(memberIDs)]
			plan.Add(memberID, topic, part)
		}
	},
}



type balanceStrategy struct {
	name   string
	coreFn func(plan BalanceStrategyPlan, memberIDs []string, topic string, partitions []int32)
}


func (s *balanceStrategy) Name() string { return s.name }


func (s *balanceStrategy) Plan(members map[string]ConsumerGroupMemberMetadata, topics map[string][]int32) (BalanceStrategyPlan, error) {
	
	mbt := make(map[string][]string)
	for memberID, meta := range members {
		for _, topic := range meta.Topics {
			mbt[topic] = append(mbt[topic], memberID)
		}
	}

	
	for topic, memberIDs := range mbt {
		sort.Sort(&balanceStrategySortable{
			topic:     topic,
			memberIDs: memberIDs,
		})
	}

	
	plan := make(BalanceStrategyPlan, len(members))
	for topic, memberIDs := range mbt {
		s.coreFn(plan, memberIDs, topic, topics[topic])
	}
	return plan, nil
}

type balanceStrategySortable struct {
	topic     string
	memberIDs []string
}

func (p balanceStrategySortable) Len() int { return len(p.memberIDs) }
func (p balanceStrategySortable) Swap(i, j int) {
	p.memberIDs[i], p.memberIDs[j] = p.memberIDs[j], p.memberIDs[i]
}
func (p balanceStrategySortable) Less(i, j int) bool {
	return balanceStrategyHashValue(p.topic, p.memberIDs[i]) < balanceStrategyHashValue(p.topic, p.memberIDs[j])
}

func balanceStrategyHashValue(vv ...string) uint32 {
	h := uint32(2166136261)
	for _, s := range vv {
		for _, c := range s {
			h ^= uint32(c)
			h *= 16777619
		}
	}
	return h
}
