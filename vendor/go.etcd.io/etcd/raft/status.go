













package raft

import (
	"fmt"

	pb "go.etcd.io/etcd/raft/raftpb"
)

type Status struct {
	ID uint64

	pb.HardState
	SoftState

	Applied  uint64
	Progress map[uint64]Progress

	LeadTransferee uint64
}

func getProgressCopy(r *raft) map[uint64]Progress {
	prs := make(map[uint64]Progress)
	for id, p := range r.prs {
		prs[id] = *p
	}

	for id, p := range r.learnerPrs {
		prs[id] = *p
	}
	return prs
}

func getStatusWithoutProgress(r *raft) Status {
	s := Status{
		ID:             r.id,
		LeadTransferee: r.leadTransferee,
	}
	s.HardState = r.hardState()
	s.SoftState = *r.softState()
	s.Applied = r.raftLog.applied
	return s
}


func getStatus(r *raft) Status {
	s := getStatusWithoutProgress(r)
	if s.RaftState == StateLeader {
		s.Progress = getProgressCopy(r)
	}
	return s
}



func (s Status) MarshalJSON() ([]byte, error) {
	j := fmt.Sprintf(`{"id":"%x","term":%d,"vote":"%x","commit":%d,"lead":"%x","raftState":%q,"applied":%d,"progress":{`,
		s.ID, s.Term, s.Vote, s.Commit, s.Lead, s.RaftState, s.Applied)

	if len(s.Progress) == 0 {
		j += "},"
	} else {
		for k, v := range s.Progress {
			subj := fmt.Sprintf(`"%x":{"match":%d,"next":%d,"state":%q},`, k, v.Match, v.Next, v.State)
			j += subj
		}
		
		j = j[:len(j)-1] + "},"
	}

	j += fmt.Sprintf(`"leadtransferee":"%x"}`, s.LeadTransferee)
	return []byte(j), nil
}

func (s Status) String() string {
	b, err := s.MarshalJSON()
	if err != nil {
		raftLogger.Panicf("unexpected error: %v", err)
	}
	return string(b)
}