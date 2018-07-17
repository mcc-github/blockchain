
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/ledger"
)

type StateListener struct {
	InterestedInNamespacesStub        func() []string
	interestedInNamespacesMutex       sync.RWMutex
	interestedInNamespacesArgsForCall []struct{}
	interestedInNamespacesReturns     struct {
		result1 []string
	}
	interestedInNamespacesReturnsOnCall map[int]struct {
		result1 []string
	}
	HandleStateUpdatesStub        func(ledgerID string, stateUpdates ledger.StateUpdates, committingBlockNum uint64) error
	handleStateUpdatesMutex       sync.RWMutex
	handleStateUpdatesArgsForCall []struct {
		ledgerID           string
		stateUpdates       ledger.StateUpdates
		committingBlockNum uint64
	}
	handleStateUpdatesReturns struct {
		result1 error
	}
	handleStateUpdatesReturnsOnCall map[int]struct {
		result1 error
	}
	StateCommitDoneStub        func(channelID string)
	stateCommitDoneMutex       sync.RWMutex
	stateCommitDoneArgsForCall []struct {
		channelID string
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *StateListener) InterestedInNamespaces() []string {
	fake.interestedInNamespacesMutex.Lock()
	ret, specificReturn := fake.interestedInNamespacesReturnsOnCall[len(fake.interestedInNamespacesArgsForCall)]
	fake.interestedInNamespacesArgsForCall = append(fake.interestedInNamespacesArgsForCall, struct{}{})
	fake.recordInvocation("InterestedInNamespaces", []interface{}{})
	fake.interestedInNamespacesMutex.Unlock()
	if fake.InterestedInNamespacesStub != nil {
		return fake.InterestedInNamespacesStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.interestedInNamespacesReturns.result1
}

func (fake *StateListener) InterestedInNamespacesCallCount() int {
	fake.interestedInNamespacesMutex.RLock()
	defer fake.interestedInNamespacesMutex.RUnlock()
	return len(fake.interestedInNamespacesArgsForCall)
}

func (fake *StateListener) InterestedInNamespacesReturns(result1 []string) {
	fake.InterestedInNamespacesStub = nil
	fake.interestedInNamespacesReturns = struct {
		result1 []string
	}{result1}
}

func (fake *StateListener) InterestedInNamespacesReturnsOnCall(i int, result1 []string) {
	fake.InterestedInNamespacesStub = nil
	if fake.interestedInNamespacesReturnsOnCall == nil {
		fake.interestedInNamespacesReturnsOnCall = make(map[int]struct {
			result1 []string
		})
	}
	fake.interestedInNamespacesReturnsOnCall[i] = struct {
		result1 []string
	}{result1}
}

func (fake *StateListener) HandleStateUpdates(ledgerID string, stateUpdates ledger.StateUpdates, committingBlockNum uint64) error {
	fake.handleStateUpdatesMutex.Lock()
	ret, specificReturn := fake.handleStateUpdatesReturnsOnCall[len(fake.handleStateUpdatesArgsForCall)]
	fake.handleStateUpdatesArgsForCall = append(fake.handleStateUpdatesArgsForCall, struct {
		ledgerID           string
		stateUpdates       ledger.StateUpdates
		committingBlockNum uint64
	}{ledgerID, stateUpdates, committingBlockNum})
	fake.recordInvocation("HandleStateUpdates", []interface{}{ledgerID, stateUpdates, committingBlockNum})
	fake.handleStateUpdatesMutex.Unlock()
	if fake.HandleStateUpdatesStub != nil {
		return fake.HandleStateUpdatesStub(ledgerID, stateUpdates, committingBlockNum)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.handleStateUpdatesReturns.result1
}

func (fake *StateListener) HandleStateUpdatesCallCount() int {
	fake.handleStateUpdatesMutex.RLock()
	defer fake.handleStateUpdatesMutex.RUnlock()
	return len(fake.handleStateUpdatesArgsForCall)
}

func (fake *StateListener) HandleStateUpdatesArgsForCall(i int) (string, ledger.StateUpdates, uint64) {
	fake.handleStateUpdatesMutex.RLock()
	defer fake.handleStateUpdatesMutex.RUnlock()
	return fake.handleStateUpdatesArgsForCall[i].ledgerID, fake.handleStateUpdatesArgsForCall[i].stateUpdates, fake.handleStateUpdatesArgsForCall[i].committingBlockNum
}

func (fake *StateListener) HandleStateUpdatesReturns(result1 error) {
	fake.HandleStateUpdatesStub = nil
	fake.handleStateUpdatesReturns = struct {
		result1 error
	}{result1}
}

func (fake *StateListener) HandleStateUpdatesReturnsOnCall(i int, result1 error) {
	fake.HandleStateUpdatesStub = nil
	if fake.handleStateUpdatesReturnsOnCall == nil {
		fake.handleStateUpdatesReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.handleStateUpdatesReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *StateListener) StateCommitDone(channelID string) {
	fake.stateCommitDoneMutex.Lock()
	fake.stateCommitDoneArgsForCall = append(fake.stateCommitDoneArgsForCall, struct {
		channelID string
	}{channelID})
	fake.recordInvocation("StateCommitDone", []interface{}{channelID})
	fake.stateCommitDoneMutex.Unlock()
	if fake.StateCommitDoneStub != nil {
		fake.StateCommitDoneStub(channelID)
	}
}

func (fake *StateListener) StateCommitDoneCallCount() int {
	fake.stateCommitDoneMutex.RLock()
	defer fake.stateCommitDoneMutex.RUnlock()
	return len(fake.stateCommitDoneArgsForCall)
}

func (fake *StateListener) StateCommitDoneArgsForCall(i int) string {
	fake.stateCommitDoneMutex.RLock()
	defer fake.stateCommitDoneMutex.RUnlock()
	return fake.stateCommitDoneArgsForCall[i].channelID
}

func (fake *StateListener) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.interestedInNamespacesMutex.RLock()
	defer fake.interestedInNamespacesMutex.RUnlock()
	fake.handleStateUpdatesMutex.RLock()
	defer fake.handleStateUpdatesMutex.RUnlock()
	fake.stateCommitDoneMutex.RLock()
	defer fake.stateCommitDoneMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *StateListener) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}

var _ ledger.StateListener = new(StateListener)
