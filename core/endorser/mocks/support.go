
package mocks

import (
	"sync"

	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	endorser_test "github.com/mcc-github/blockchain/core/endorser"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

type Support struct {
	SignStub        func(message []byte) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		message []byte
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SerializeStub        func() ([]byte, error)
	serializeMutex       sync.RWMutex
	serializeArgsForCall []struct{}
	serializeReturns     struct {
		result1 []byte
		result2 error
	}
	serializeReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	IsSysCCAndNotInvokableExternalStub        func(name string) bool
	isSysCCAndNotInvokableExternalMutex       sync.RWMutex
	isSysCCAndNotInvokableExternalArgsForCall []struct {
		name string
	}
	isSysCCAndNotInvokableExternalReturns struct {
		result1 bool
	}
	isSysCCAndNotInvokableExternalReturnsOnCall map[int]struct {
		result1 bool
	}
	GetTxSimulatorStub        func(ledgername string, txid string) (ledger.TxSimulator, error)
	getTxSimulatorMutex       sync.RWMutex
	getTxSimulatorArgsForCall []struct {
		ledgername string
		txid       string
	}
	getTxSimulatorReturns struct {
		result1 ledger.TxSimulator
		result2 error
	}
	getTxSimulatorReturnsOnCall map[int]struct {
		result1 ledger.TxSimulator
		result2 error
	}
	GetHistoryQueryExecutorStub        func(ledgername string) (ledger.HistoryQueryExecutor, error)
	getHistoryQueryExecutorMutex       sync.RWMutex
	getHistoryQueryExecutorArgsForCall []struct {
		ledgername string
	}
	getHistoryQueryExecutorReturns struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}
	getHistoryQueryExecutorReturnsOnCall map[int]struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}
	GetTransactionByIDStub        func(chid, txID string) (*pb.ProcessedTransaction, error)
	getTransactionByIDMutex       sync.RWMutex
	getTransactionByIDArgsForCall []struct {
		chid string
		txID string
	}
	getTransactionByIDReturns struct {
		result1 *pb.ProcessedTransaction
		result2 error
	}
	getTransactionByIDReturnsOnCall map[int]struct {
		result1 *pb.ProcessedTransaction
		result2 error
	}
	IsSysCCStub        func(name string) bool
	isSysCCMutex       sync.RWMutex
	isSysCCArgsForCall []struct {
		name string
	}
	isSysCCReturns struct {
		result1 bool
	}
	isSysCCReturnsOnCall map[int]struct {
		result1 bool
	}
	ExecuteStub        func(txParams *ccprovider.TransactionParams, cid, name, version, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error)
	executeMutex       sync.RWMutex
	executeArgsForCall []struct {
		txParams   *ccprovider.TransactionParams
		cid        string
		name       string
		version    string
		txid       string
		signedProp *pb.SignedProposal
		prop       *pb.Proposal
		input      *pb.ChaincodeInput
	}
	executeReturns struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}
	executeReturnsOnCall map[int]struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}
	ExecuteLegacyInitStub        func(txParams *ccprovider.TransactionParams, cid, name, version, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, spec *pb.ChaincodeDeploymentSpec) (*pb.Response, *pb.ChaincodeEvent, error)
	executeLegacyInitMutex       sync.RWMutex
	executeLegacyInitArgsForCall []struct {
		txParams   *ccprovider.TransactionParams
		cid        string
		name       string
		version    string
		txid       string
		signedProp *pb.SignedProposal
		prop       *pb.Proposal
		spec       *pb.ChaincodeDeploymentSpec
	}
	executeLegacyInitReturns struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}
	executeLegacyInitReturnsOnCall map[int]struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}
	GetChaincodeDefinitionStub        func(chaincodeID string, txsim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error)
	getChaincodeDefinitionMutex       sync.RWMutex
	getChaincodeDefinitionArgsForCall []struct {
		chaincodeID string
		txsim       ledger.QueryExecutor
	}
	getChaincodeDefinitionReturns struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}
	getChaincodeDefinitionReturnsOnCall map[int]struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}
	CheckACLStub        func(signedProp *pb.SignedProposal, chdr *common.ChannelHeader, shdr *common.SignatureHeader, hdrext *pb.ChaincodeHeaderExtension) error
	checkACLMutex       sync.RWMutex
	checkACLArgsForCall []struct {
		signedProp *pb.SignedProposal
		chdr       *common.ChannelHeader
		shdr       *common.SignatureHeader
		hdrext     *pb.ChaincodeHeaderExtension
	}
	checkACLReturns struct {
		result1 error
	}
	checkACLReturnsOnCall map[int]struct {
		result1 error
	}
	IsJavaCCStub        func(buf []byte) (bool, error)
	isJavaCCMutex       sync.RWMutex
	isJavaCCArgsForCall []struct {
		buf []byte
	}
	isJavaCCReturns struct {
		result1 bool
		result2 error
	}
	isJavaCCReturnsOnCall map[int]struct {
		result1 bool
		result2 error
	}
	CheckInstantiationPolicyStub        func(name, version string, cd ccprovider.ChaincodeDefinition) error
	checkInstantiationPolicyMutex       sync.RWMutex
	checkInstantiationPolicyArgsForCall []struct {
		name    string
		version string
		cd      ccprovider.ChaincodeDefinition
	}
	checkInstantiationPolicyReturns struct {
		result1 error
	}
	checkInstantiationPolicyReturnsOnCall map[int]struct {
		result1 error
	}
	GetChaincodeDeploymentSpecFSStub        func(cds *pb.ChaincodeDeploymentSpec) (*pb.ChaincodeDeploymentSpec, error)
	getChaincodeDeploymentSpecFSMutex       sync.RWMutex
	getChaincodeDeploymentSpecFSArgsForCall []struct {
		cds *pb.ChaincodeDeploymentSpec
	}
	getChaincodeDeploymentSpecFSReturns struct {
		result1 *pb.ChaincodeDeploymentSpec
		result2 error
	}
	getChaincodeDeploymentSpecFSReturnsOnCall map[int]struct {
		result1 *pb.ChaincodeDeploymentSpec
		result2 error
	}
	GetApplicationConfigStub        func(cid string) (channelconfig.Application, bool)
	getApplicationConfigMutex       sync.RWMutex
	getApplicationConfigArgsForCall []struct {
		cid string
	}
	getApplicationConfigReturns struct {
		result1 channelconfig.Application
		result2 bool
	}
	getApplicationConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Application
		result2 bool
	}
	NewQueryCreatorStub        func(channel string) (endorser_test.QueryCreator, error)
	newQueryCreatorMutex       sync.RWMutex
	newQueryCreatorArgsForCall []struct {
		channel string
	}
	newQueryCreatorReturns struct {
		result1 endorser_test.QueryCreator
		result2 error
	}
	newQueryCreatorReturnsOnCall map[int]struct {
		result1 endorser_test.QueryCreator
		result2 error
	}
	EndorseWithPluginStub        func(ctx endorser_test.Context) (*pb.ProposalResponse, error)
	endorseWithPluginMutex       sync.RWMutex
	endorseWithPluginArgsForCall []struct {
		ctx endorser_test.Context
	}
	endorseWithPluginReturns struct {
		result1 *pb.ProposalResponse
		result2 error
	}
	endorseWithPluginReturnsOnCall map[int]struct {
		result1 *pb.ProposalResponse
		result2 error
	}
	GetLedgerHeightStub        func(channelID string) (uint64, error)
	getLedgerHeightMutex       sync.RWMutex
	getLedgerHeightArgsForCall []struct {
		channelID string
	}
	getLedgerHeightReturns struct {
		result1 uint64
		result2 error
	}
	getLedgerHeightReturnsOnCall map[int]struct {
		result1 uint64
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *Support) Sign(message []byte) ([]byte, error) {
	var messageCopy []byte
	if message != nil {
		messageCopy = make([]byte, len(message))
		copy(messageCopy, message)
	}
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		message []byte
	}{messageCopy})
	fake.recordInvocation("Sign", []interface{}{messageCopy})
	fake.signMutex.Unlock()
	if fake.SignStub != nil {
		return fake.SignStub(message)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.signReturns.result1, fake.signReturns.result2
}

func (fake *Support) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *Support) SignArgsForCall(i int) []byte {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return fake.signArgsForCall[i].message
}

func (fake *Support) SignReturns(result1 []byte, result2 error) {
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Support) SignReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.SignStub = nil
	if fake.signReturnsOnCall == nil {
		fake.signReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.signReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Support) Serialize() ([]byte, error) {
	fake.serializeMutex.Lock()
	ret, specificReturn := fake.serializeReturnsOnCall[len(fake.serializeArgsForCall)]
	fake.serializeArgsForCall = append(fake.serializeArgsForCall, struct{}{})
	fake.recordInvocation("Serialize", []interface{}{})
	fake.serializeMutex.Unlock()
	if fake.SerializeStub != nil {
		return fake.SerializeStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.serializeReturns.result1, fake.serializeReturns.result2
}

func (fake *Support) SerializeCallCount() int {
	fake.serializeMutex.RLock()
	defer fake.serializeMutex.RUnlock()
	return len(fake.serializeArgsForCall)
}

func (fake *Support) SerializeReturns(result1 []byte, result2 error) {
	fake.SerializeStub = nil
	fake.serializeReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Support) SerializeReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.SerializeStub = nil
	if fake.serializeReturnsOnCall == nil {
		fake.serializeReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.serializeReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *Support) IsSysCCAndNotInvokableExternal(name string) bool {
	fake.isSysCCAndNotInvokableExternalMutex.Lock()
	ret, specificReturn := fake.isSysCCAndNotInvokableExternalReturnsOnCall[len(fake.isSysCCAndNotInvokableExternalArgsForCall)]
	fake.isSysCCAndNotInvokableExternalArgsForCall = append(fake.isSysCCAndNotInvokableExternalArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("IsSysCCAndNotInvokableExternal", []interface{}{name})
	fake.isSysCCAndNotInvokableExternalMutex.Unlock()
	if fake.IsSysCCAndNotInvokableExternalStub != nil {
		return fake.IsSysCCAndNotInvokableExternalStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isSysCCAndNotInvokableExternalReturns.result1
}

func (fake *Support) IsSysCCAndNotInvokableExternalCallCount() int {
	fake.isSysCCAndNotInvokableExternalMutex.RLock()
	defer fake.isSysCCAndNotInvokableExternalMutex.RUnlock()
	return len(fake.isSysCCAndNotInvokableExternalArgsForCall)
}

func (fake *Support) IsSysCCAndNotInvokableExternalArgsForCall(i int) string {
	fake.isSysCCAndNotInvokableExternalMutex.RLock()
	defer fake.isSysCCAndNotInvokableExternalMutex.RUnlock()
	return fake.isSysCCAndNotInvokableExternalArgsForCall[i].name
}

func (fake *Support) IsSysCCAndNotInvokableExternalReturns(result1 bool) {
	fake.IsSysCCAndNotInvokableExternalStub = nil
	fake.isSysCCAndNotInvokableExternalReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Support) IsSysCCAndNotInvokableExternalReturnsOnCall(i int, result1 bool) {
	fake.IsSysCCAndNotInvokableExternalStub = nil
	if fake.isSysCCAndNotInvokableExternalReturnsOnCall == nil {
		fake.isSysCCAndNotInvokableExternalReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isSysCCAndNotInvokableExternalReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Support) GetTxSimulator(ledgername string, txid string) (ledger.TxSimulator, error) {
	fake.getTxSimulatorMutex.Lock()
	ret, specificReturn := fake.getTxSimulatorReturnsOnCall[len(fake.getTxSimulatorArgsForCall)]
	fake.getTxSimulatorArgsForCall = append(fake.getTxSimulatorArgsForCall, struct {
		ledgername string
		txid       string
	}{ledgername, txid})
	fake.recordInvocation("GetTxSimulator", []interface{}{ledgername, txid})
	fake.getTxSimulatorMutex.Unlock()
	if fake.GetTxSimulatorStub != nil {
		return fake.GetTxSimulatorStub(ledgername, txid)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTxSimulatorReturns.result1, fake.getTxSimulatorReturns.result2
}

func (fake *Support) GetTxSimulatorCallCount() int {
	fake.getTxSimulatorMutex.RLock()
	defer fake.getTxSimulatorMutex.RUnlock()
	return len(fake.getTxSimulatorArgsForCall)
}

func (fake *Support) GetTxSimulatorArgsForCall(i int) (string, string) {
	fake.getTxSimulatorMutex.RLock()
	defer fake.getTxSimulatorMutex.RUnlock()
	return fake.getTxSimulatorArgsForCall[i].ledgername, fake.getTxSimulatorArgsForCall[i].txid
}

func (fake *Support) GetTxSimulatorReturns(result1 ledger.TxSimulator, result2 error) {
	fake.GetTxSimulatorStub = nil
	fake.getTxSimulatorReturns = struct {
		result1 ledger.TxSimulator
		result2 error
	}{result1, result2}
}

func (fake *Support) GetTxSimulatorReturnsOnCall(i int, result1 ledger.TxSimulator, result2 error) {
	fake.GetTxSimulatorStub = nil
	if fake.getTxSimulatorReturnsOnCall == nil {
		fake.getTxSimulatorReturnsOnCall = make(map[int]struct {
			result1 ledger.TxSimulator
			result2 error
		})
	}
	fake.getTxSimulatorReturnsOnCall[i] = struct {
		result1 ledger.TxSimulator
		result2 error
	}{result1, result2}
}

func (fake *Support) GetHistoryQueryExecutor(ledgername string) (ledger.HistoryQueryExecutor, error) {
	fake.getHistoryQueryExecutorMutex.Lock()
	ret, specificReturn := fake.getHistoryQueryExecutorReturnsOnCall[len(fake.getHistoryQueryExecutorArgsForCall)]
	fake.getHistoryQueryExecutorArgsForCall = append(fake.getHistoryQueryExecutorArgsForCall, struct {
		ledgername string
	}{ledgername})
	fake.recordInvocation("GetHistoryQueryExecutor", []interface{}{ledgername})
	fake.getHistoryQueryExecutorMutex.Unlock()
	if fake.GetHistoryQueryExecutorStub != nil {
		return fake.GetHistoryQueryExecutorStub(ledgername)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getHistoryQueryExecutorReturns.result1, fake.getHistoryQueryExecutorReturns.result2
}

func (fake *Support) GetHistoryQueryExecutorCallCount() int {
	fake.getHistoryQueryExecutorMutex.RLock()
	defer fake.getHistoryQueryExecutorMutex.RUnlock()
	return len(fake.getHistoryQueryExecutorArgsForCall)
}

func (fake *Support) GetHistoryQueryExecutorArgsForCall(i int) string {
	fake.getHistoryQueryExecutorMutex.RLock()
	defer fake.getHistoryQueryExecutorMutex.RUnlock()
	return fake.getHistoryQueryExecutorArgsForCall[i].ledgername
}

func (fake *Support) GetHistoryQueryExecutorReturns(result1 ledger.HistoryQueryExecutor, result2 error) {
	fake.GetHistoryQueryExecutorStub = nil
	fake.getHistoryQueryExecutorReturns = struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *Support) GetHistoryQueryExecutorReturnsOnCall(i int, result1 ledger.HistoryQueryExecutor, result2 error) {
	fake.GetHistoryQueryExecutorStub = nil
	if fake.getHistoryQueryExecutorReturnsOnCall == nil {
		fake.getHistoryQueryExecutorReturnsOnCall = make(map[int]struct {
			result1 ledger.HistoryQueryExecutor
			result2 error
		})
	}
	fake.getHistoryQueryExecutorReturnsOnCall[i] = struct {
		result1 ledger.HistoryQueryExecutor
		result2 error
	}{result1, result2}
}

func (fake *Support) GetTransactionByID(chid string, txID string) (*pb.ProcessedTransaction, error) {
	fake.getTransactionByIDMutex.Lock()
	ret, specificReturn := fake.getTransactionByIDReturnsOnCall[len(fake.getTransactionByIDArgsForCall)]
	fake.getTransactionByIDArgsForCall = append(fake.getTransactionByIDArgsForCall, struct {
		chid string
		txID string
	}{chid, txID})
	fake.recordInvocation("GetTransactionByID", []interface{}{chid, txID})
	fake.getTransactionByIDMutex.Unlock()
	if fake.GetTransactionByIDStub != nil {
		return fake.GetTransactionByIDStub(chid, txID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTransactionByIDReturns.result1, fake.getTransactionByIDReturns.result2
}

func (fake *Support) GetTransactionByIDCallCount() int {
	fake.getTransactionByIDMutex.RLock()
	defer fake.getTransactionByIDMutex.RUnlock()
	return len(fake.getTransactionByIDArgsForCall)
}

func (fake *Support) GetTransactionByIDArgsForCall(i int) (string, string) {
	fake.getTransactionByIDMutex.RLock()
	defer fake.getTransactionByIDMutex.RUnlock()
	return fake.getTransactionByIDArgsForCall[i].chid, fake.getTransactionByIDArgsForCall[i].txID
}

func (fake *Support) GetTransactionByIDReturns(result1 *pb.ProcessedTransaction, result2 error) {
	fake.GetTransactionByIDStub = nil
	fake.getTransactionByIDReturns = struct {
		result1 *pb.ProcessedTransaction
		result2 error
	}{result1, result2}
}

func (fake *Support) GetTransactionByIDReturnsOnCall(i int, result1 *pb.ProcessedTransaction, result2 error) {
	fake.GetTransactionByIDStub = nil
	if fake.getTransactionByIDReturnsOnCall == nil {
		fake.getTransactionByIDReturnsOnCall = make(map[int]struct {
			result1 *pb.ProcessedTransaction
			result2 error
		})
	}
	fake.getTransactionByIDReturnsOnCall[i] = struct {
		result1 *pb.ProcessedTransaction
		result2 error
	}{result1, result2}
}

func (fake *Support) IsSysCC(name string) bool {
	fake.isSysCCMutex.Lock()
	ret, specificReturn := fake.isSysCCReturnsOnCall[len(fake.isSysCCArgsForCall)]
	fake.isSysCCArgsForCall = append(fake.isSysCCArgsForCall, struct {
		name string
	}{name})
	fake.recordInvocation("IsSysCC", []interface{}{name})
	fake.isSysCCMutex.Unlock()
	if fake.IsSysCCStub != nil {
		return fake.IsSysCCStub(name)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.isSysCCReturns.result1
}

func (fake *Support) IsSysCCCallCount() int {
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	return len(fake.isSysCCArgsForCall)
}

func (fake *Support) IsSysCCArgsForCall(i int) string {
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	return fake.isSysCCArgsForCall[i].name
}

func (fake *Support) IsSysCCReturns(result1 bool) {
	fake.IsSysCCStub = nil
	fake.isSysCCReturns = struct {
		result1 bool
	}{result1}
}

func (fake *Support) IsSysCCReturnsOnCall(i int, result1 bool) {
	fake.IsSysCCStub = nil
	if fake.isSysCCReturnsOnCall == nil {
		fake.isSysCCReturnsOnCall = make(map[int]struct {
			result1 bool
		})
	}
	fake.isSysCCReturnsOnCall[i] = struct {
		result1 bool
	}{result1}
}

func (fake *Support) Execute(txParams *ccprovider.TransactionParams, cid string, name string, version string, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, input *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error) {
	fake.executeMutex.Lock()
	ret, specificReturn := fake.executeReturnsOnCall[len(fake.executeArgsForCall)]
	fake.executeArgsForCall = append(fake.executeArgsForCall, struct {
		txParams   *ccprovider.TransactionParams
		cid        string
		name       string
		version    string
		txid       string
		signedProp *pb.SignedProposal
		prop       *pb.Proposal
		input      *pb.ChaincodeInput
	}{txParams, cid, name, version, txid, signedProp, prop, input})
	fake.recordInvocation("Execute", []interface{}{txParams, cid, name, version, txid, signedProp, prop, input})
	fake.executeMutex.Unlock()
	if fake.ExecuteStub != nil {
		return fake.ExecuteStub(txParams, cid, name, version, txid, signedProp, prop, input)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.executeReturns.result1, fake.executeReturns.result2, fake.executeReturns.result3
}

func (fake *Support) ExecuteCallCount() int {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return len(fake.executeArgsForCall)
}

func (fake *Support) ExecuteArgsForCall(i int) (*ccprovider.TransactionParams, string, string, string, string, *pb.SignedProposal, *pb.Proposal, *pb.ChaincodeInput) {
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	return fake.executeArgsForCall[i].txParams, fake.executeArgsForCall[i].cid, fake.executeArgsForCall[i].name, fake.executeArgsForCall[i].version, fake.executeArgsForCall[i].txid, fake.executeArgsForCall[i].signedProp, fake.executeArgsForCall[i].prop, fake.executeArgsForCall[i].input
}

func (fake *Support) ExecuteReturns(result1 *pb.Response, result2 *pb.ChaincodeEvent, result3 error) {
	fake.ExecuteStub = nil
	fake.executeReturns = struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}{result1, result2, result3}
}

func (fake *Support) ExecuteReturnsOnCall(i int, result1 *pb.Response, result2 *pb.ChaincodeEvent, result3 error) {
	fake.ExecuteStub = nil
	if fake.executeReturnsOnCall == nil {
		fake.executeReturnsOnCall = make(map[int]struct {
			result1 *pb.Response
			result2 *pb.ChaincodeEvent
			result3 error
		})
	}
	fake.executeReturnsOnCall[i] = struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}{result1, result2, result3}
}

func (fake *Support) ExecuteLegacyInit(txParams *ccprovider.TransactionParams, cid string, name string, version string, txid string, signedProp *pb.SignedProposal, prop *pb.Proposal, spec *pb.ChaincodeDeploymentSpec) (*pb.Response, *pb.ChaincodeEvent, error) {
	fake.executeLegacyInitMutex.Lock()
	ret, specificReturn := fake.executeLegacyInitReturnsOnCall[len(fake.executeLegacyInitArgsForCall)]
	fake.executeLegacyInitArgsForCall = append(fake.executeLegacyInitArgsForCall, struct {
		txParams   *ccprovider.TransactionParams
		cid        string
		name       string
		version    string
		txid       string
		signedProp *pb.SignedProposal
		prop       *pb.Proposal
		spec       *pb.ChaincodeDeploymentSpec
	}{txParams, cid, name, version, txid, signedProp, prop, spec})
	fake.recordInvocation("ExecuteLegacyInit", []interface{}{txParams, cid, name, version, txid, signedProp, prop, spec})
	fake.executeLegacyInitMutex.Unlock()
	if fake.ExecuteLegacyInitStub != nil {
		return fake.ExecuteLegacyInitStub(txParams, cid, name, version, txid, signedProp, prop, spec)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.executeLegacyInitReturns.result1, fake.executeLegacyInitReturns.result2, fake.executeLegacyInitReturns.result3
}

func (fake *Support) ExecuteLegacyInitCallCount() int {
	fake.executeLegacyInitMutex.RLock()
	defer fake.executeLegacyInitMutex.RUnlock()
	return len(fake.executeLegacyInitArgsForCall)
}

func (fake *Support) ExecuteLegacyInitArgsForCall(i int) (*ccprovider.TransactionParams, string, string, string, string, *pb.SignedProposal, *pb.Proposal, *pb.ChaincodeDeploymentSpec) {
	fake.executeLegacyInitMutex.RLock()
	defer fake.executeLegacyInitMutex.RUnlock()
	return fake.executeLegacyInitArgsForCall[i].txParams, fake.executeLegacyInitArgsForCall[i].cid, fake.executeLegacyInitArgsForCall[i].name, fake.executeLegacyInitArgsForCall[i].version, fake.executeLegacyInitArgsForCall[i].txid, fake.executeLegacyInitArgsForCall[i].signedProp, fake.executeLegacyInitArgsForCall[i].prop, fake.executeLegacyInitArgsForCall[i].spec
}

func (fake *Support) ExecuteLegacyInitReturns(result1 *pb.Response, result2 *pb.ChaincodeEvent, result3 error) {
	fake.ExecuteLegacyInitStub = nil
	fake.executeLegacyInitReturns = struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}{result1, result2, result3}
}

func (fake *Support) ExecuteLegacyInitReturnsOnCall(i int, result1 *pb.Response, result2 *pb.ChaincodeEvent, result3 error) {
	fake.ExecuteLegacyInitStub = nil
	if fake.executeLegacyInitReturnsOnCall == nil {
		fake.executeLegacyInitReturnsOnCall = make(map[int]struct {
			result1 *pb.Response
			result2 *pb.ChaincodeEvent
			result3 error
		})
	}
	fake.executeLegacyInitReturnsOnCall[i] = struct {
		result1 *pb.Response
		result2 *pb.ChaincodeEvent
		result3 error
	}{result1, result2, result3}
}

func (fake *Support) GetChaincodeDefinition(chaincodeID string, txsim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	fake.getChaincodeDefinitionMutex.Lock()
	ret, specificReturn := fake.getChaincodeDefinitionReturnsOnCall[len(fake.getChaincodeDefinitionArgsForCall)]
	fake.getChaincodeDefinitionArgsForCall = append(fake.getChaincodeDefinitionArgsForCall, struct {
		chaincodeID string
		txsim       ledger.QueryExecutor
	}{chaincodeID, txsim})
	fake.recordInvocation("GetChaincodeDefinition", []interface{}{chaincodeID, txsim})
	fake.getChaincodeDefinitionMutex.Unlock()
	if fake.GetChaincodeDefinitionStub != nil {
		return fake.GetChaincodeDefinitionStub(chaincodeID, txsim)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getChaincodeDefinitionReturns.result1, fake.getChaincodeDefinitionReturns.result2
}

func (fake *Support) GetChaincodeDefinitionCallCount() int {
	fake.getChaincodeDefinitionMutex.RLock()
	defer fake.getChaincodeDefinitionMutex.RUnlock()
	return len(fake.getChaincodeDefinitionArgsForCall)
}

func (fake *Support) GetChaincodeDefinitionArgsForCall(i int) (string, ledger.QueryExecutor) {
	fake.getChaincodeDefinitionMutex.RLock()
	defer fake.getChaincodeDefinitionMutex.RUnlock()
	return fake.getChaincodeDefinitionArgsForCall[i].chaincodeID, fake.getChaincodeDefinitionArgsForCall[i].txsim
}

func (fake *Support) GetChaincodeDefinitionReturns(result1 ccprovider.ChaincodeDefinition, result2 error) {
	fake.GetChaincodeDefinitionStub = nil
	fake.getChaincodeDefinitionReturns = struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}{result1, result2}
}

func (fake *Support) GetChaincodeDefinitionReturnsOnCall(i int, result1 ccprovider.ChaincodeDefinition, result2 error) {
	fake.GetChaincodeDefinitionStub = nil
	if fake.getChaincodeDefinitionReturnsOnCall == nil {
		fake.getChaincodeDefinitionReturnsOnCall = make(map[int]struct {
			result1 ccprovider.ChaincodeDefinition
			result2 error
		})
	}
	fake.getChaincodeDefinitionReturnsOnCall[i] = struct {
		result1 ccprovider.ChaincodeDefinition
		result2 error
	}{result1, result2}
}

func (fake *Support) CheckACL(signedProp *pb.SignedProposal, chdr *common.ChannelHeader, shdr *common.SignatureHeader, hdrext *pb.ChaincodeHeaderExtension) error {
	fake.checkACLMutex.Lock()
	ret, specificReturn := fake.checkACLReturnsOnCall[len(fake.checkACLArgsForCall)]
	fake.checkACLArgsForCall = append(fake.checkACLArgsForCall, struct {
		signedProp *pb.SignedProposal
		chdr       *common.ChannelHeader
		shdr       *common.SignatureHeader
		hdrext     *pb.ChaincodeHeaderExtension
	}{signedProp, chdr, shdr, hdrext})
	fake.recordInvocation("CheckACL", []interface{}{signedProp, chdr, shdr, hdrext})
	fake.checkACLMutex.Unlock()
	if fake.CheckACLStub != nil {
		return fake.CheckACLStub(signedProp, chdr, shdr, hdrext)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.checkACLReturns.result1
}

func (fake *Support) CheckACLCallCount() int {
	fake.checkACLMutex.RLock()
	defer fake.checkACLMutex.RUnlock()
	return len(fake.checkACLArgsForCall)
}

func (fake *Support) CheckACLArgsForCall(i int) (*pb.SignedProposal, *common.ChannelHeader, *common.SignatureHeader, *pb.ChaincodeHeaderExtension) {
	fake.checkACLMutex.RLock()
	defer fake.checkACLMutex.RUnlock()
	return fake.checkACLArgsForCall[i].signedProp, fake.checkACLArgsForCall[i].chdr, fake.checkACLArgsForCall[i].shdr, fake.checkACLArgsForCall[i].hdrext
}

func (fake *Support) CheckACLReturns(result1 error) {
	fake.CheckACLStub = nil
	fake.checkACLReturns = struct {
		result1 error
	}{result1}
}

func (fake *Support) CheckACLReturnsOnCall(i int, result1 error) {
	fake.CheckACLStub = nil
	if fake.checkACLReturnsOnCall == nil {
		fake.checkACLReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkACLReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Support) IsJavaCC(buf []byte) (bool, error) {
	var bufCopy []byte
	if buf != nil {
		bufCopy = make([]byte, len(buf))
		copy(bufCopy, buf)
	}
	fake.isJavaCCMutex.Lock()
	ret, specificReturn := fake.isJavaCCReturnsOnCall[len(fake.isJavaCCArgsForCall)]
	fake.isJavaCCArgsForCall = append(fake.isJavaCCArgsForCall, struct {
		buf []byte
	}{bufCopy})
	fake.recordInvocation("IsJavaCC", []interface{}{bufCopy})
	fake.isJavaCCMutex.Unlock()
	if fake.IsJavaCCStub != nil {
		return fake.IsJavaCCStub(buf)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.isJavaCCReturns.result1, fake.isJavaCCReturns.result2
}

func (fake *Support) IsJavaCCCallCount() int {
	fake.isJavaCCMutex.RLock()
	defer fake.isJavaCCMutex.RUnlock()
	return len(fake.isJavaCCArgsForCall)
}

func (fake *Support) IsJavaCCArgsForCall(i int) []byte {
	fake.isJavaCCMutex.RLock()
	defer fake.isJavaCCMutex.RUnlock()
	return fake.isJavaCCArgsForCall[i].buf
}

func (fake *Support) IsJavaCCReturns(result1 bool, result2 error) {
	fake.IsJavaCCStub = nil
	fake.isJavaCCReturns = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *Support) IsJavaCCReturnsOnCall(i int, result1 bool, result2 error) {
	fake.IsJavaCCStub = nil
	if fake.isJavaCCReturnsOnCall == nil {
		fake.isJavaCCReturnsOnCall = make(map[int]struct {
			result1 bool
			result2 error
		})
	}
	fake.isJavaCCReturnsOnCall[i] = struct {
		result1 bool
		result2 error
	}{result1, result2}
}

func (fake *Support) CheckInstantiationPolicy(name string, version string, cd ccprovider.ChaincodeDefinition) error {
	fake.checkInstantiationPolicyMutex.Lock()
	ret, specificReturn := fake.checkInstantiationPolicyReturnsOnCall[len(fake.checkInstantiationPolicyArgsForCall)]
	fake.checkInstantiationPolicyArgsForCall = append(fake.checkInstantiationPolicyArgsForCall, struct {
		name    string
		version string
		cd      ccprovider.ChaincodeDefinition
	}{name, version, cd})
	fake.recordInvocation("CheckInstantiationPolicy", []interface{}{name, version, cd})
	fake.checkInstantiationPolicyMutex.Unlock()
	if fake.CheckInstantiationPolicyStub != nil {
		return fake.CheckInstantiationPolicyStub(name, version, cd)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.checkInstantiationPolicyReturns.result1
}

func (fake *Support) CheckInstantiationPolicyCallCount() int {
	fake.checkInstantiationPolicyMutex.RLock()
	defer fake.checkInstantiationPolicyMutex.RUnlock()
	return len(fake.checkInstantiationPolicyArgsForCall)
}

func (fake *Support) CheckInstantiationPolicyArgsForCall(i int) (string, string, ccprovider.ChaincodeDefinition) {
	fake.checkInstantiationPolicyMutex.RLock()
	defer fake.checkInstantiationPolicyMutex.RUnlock()
	return fake.checkInstantiationPolicyArgsForCall[i].name, fake.checkInstantiationPolicyArgsForCall[i].version, fake.checkInstantiationPolicyArgsForCall[i].cd
}

func (fake *Support) CheckInstantiationPolicyReturns(result1 error) {
	fake.CheckInstantiationPolicyStub = nil
	fake.checkInstantiationPolicyReturns = struct {
		result1 error
	}{result1}
}

func (fake *Support) CheckInstantiationPolicyReturnsOnCall(i int, result1 error) {
	fake.CheckInstantiationPolicyStub = nil
	if fake.checkInstantiationPolicyReturnsOnCall == nil {
		fake.checkInstantiationPolicyReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.checkInstantiationPolicyReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *Support) GetChaincodeDeploymentSpecFS(cds *pb.ChaincodeDeploymentSpec) (*pb.ChaincodeDeploymentSpec, error) {
	fake.getChaincodeDeploymentSpecFSMutex.Lock()
	ret, specificReturn := fake.getChaincodeDeploymentSpecFSReturnsOnCall[len(fake.getChaincodeDeploymentSpecFSArgsForCall)]
	fake.getChaincodeDeploymentSpecFSArgsForCall = append(fake.getChaincodeDeploymentSpecFSArgsForCall, struct {
		cds *pb.ChaincodeDeploymentSpec
	}{cds})
	fake.recordInvocation("GetChaincodeDeploymentSpecFS", []interface{}{cds})
	fake.getChaincodeDeploymentSpecFSMutex.Unlock()
	if fake.GetChaincodeDeploymentSpecFSStub != nil {
		return fake.GetChaincodeDeploymentSpecFSStub(cds)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getChaincodeDeploymentSpecFSReturns.result1, fake.getChaincodeDeploymentSpecFSReturns.result2
}

func (fake *Support) GetChaincodeDeploymentSpecFSCallCount() int {
	fake.getChaincodeDeploymentSpecFSMutex.RLock()
	defer fake.getChaincodeDeploymentSpecFSMutex.RUnlock()
	return len(fake.getChaincodeDeploymentSpecFSArgsForCall)
}

func (fake *Support) GetChaincodeDeploymentSpecFSArgsForCall(i int) *pb.ChaincodeDeploymentSpec {
	fake.getChaincodeDeploymentSpecFSMutex.RLock()
	defer fake.getChaincodeDeploymentSpecFSMutex.RUnlock()
	return fake.getChaincodeDeploymentSpecFSArgsForCall[i].cds
}

func (fake *Support) GetChaincodeDeploymentSpecFSReturns(result1 *pb.ChaincodeDeploymentSpec, result2 error) {
	fake.GetChaincodeDeploymentSpecFSStub = nil
	fake.getChaincodeDeploymentSpecFSReturns = struct {
		result1 *pb.ChaincodeDeploymentSpec
		result2 error
	}{result1, result2}
}

func (fake *Support) GetChaincodeDeploymentSpecFSReturnsOnCall(i int, result1 *pb.ChaincodeDeploymentSpec, result2 error) {
	fake.GetChaincodeDeploymentSpecFSStub = nil
	if fake.getChaincodeDeploymentSpecFSReturnsOnCall == nil {
		fake.getChaincodeDeploymentSpecFSReturnsOnCall = make(map[int]struct {
			result1 *pb.ChaincodeDeploymentSpec
			result2 error
		})
	}
	fake.getChaincodeDeploymentSpecFSReturnsOnCall[i] = struct {
		result1 *pb.ChaincodeDeploymentSpec
		result2 error
	}{result1, result2}
}

func (fake *Support) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	fake.getApplicationConfigMutex.Lock()
	ret, specificReturn := fake.getApplicationConfigReturnsOnCall[len(fake.getApplicationConfigArgsForCall)]
	fake.getApplicationConfigArgsForCall = append(fake.getApplicationConfigArgsForCall, struct {
		cid string
	}{cid})
	fake.recordInvocation("GetApplicationConfig", []interface{}{cid})
	fake.getApplicationConfigMutex.Unlock()
	if fake.GetApplicationConfigStub != nil {
		return fake.GetApplicationConfigStub(cid)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getApplicationConfigReturns.result1, fake.getApplicationConfigReturns.result2
}

func (fake *Support) GetApplicationConfigCallCount() int {
	fake.getApplicationConfigMutex.RLock()
	defer fake.getApplicationConfigMutex.RUnlock()
	return len(fake.getApplicationConfigArgsForCall)
}

func (fake *Support) GetApplicationConfigArgsForCall(i int) string {
	fake.getApplicationConfigMutex.RLock()
	defer fake.getApplicationConfigMutex.RUnlock()
	return fake.getApplicationConfigArgsForCall[i].cid
}

func (fake *Support) GetApplicationConfigReturns(result1 channelconfig.Application, result2 bool) {
	fake.GetApplicationConfigStub = nil
	fake.getApplicationConfigReturns = struct {
		result1 channelconfig.Application
		result2 bool
	}{result1, result2}
}

func (fake *Support) GetApplicationConfigReturnsOnCall(i int, result1 channelconfig.Application, result2 bool) {
	fake.GetApplicationConfigStub = nil
	if fake.getApplicationConfigReturnsOnCall == nil {
		fake.getApplicationConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Application
			result2 bool
		})
	}
	fake.getApplicationConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Application
		result2 bool
	}{result1, result2}
}

func (fake *Support) NewQueryCreator(channel string) (endorser_test.QueryCreator, error) {
	fake.newQueryCreatorMutex.Lock()
	ret, specificReturn := fake.newQueryCreatorReturnsOnCall[len(fake.newQueryCreatorArgsForCall)]
	fake.newQueryCreatorArgsForCall = append(fake.newQueryCreatorArgsForCall, struct {
		channel string
	}{channel})
	fake.recordInvocation("NewQueryCreator", []interface{}{channel})
	fake.newQueryCreatorMutex.Unlock()
	if fake.NewQueryCreatorStub != nil {
		return fake.NewQueryCreatorStub(channel)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.newQueryCreatorReturns.result1, fake.newQueryCreatorReturns.result2
}

func (fake *Support) NewQueryCreatorCallCount() int {
	fake.newQueryCreatorMutex.RLock()
	defer fake.newQueryCreatorMutex.RUnlock()
	return len(fake.newQueryCreatorArgsForCall)
}

func (fake *Support) NewQueryCreatorArgsForCall(i int) string {
	fake.newQueryCreatorMutex.RLock()
	defer fake.newQueryCreatorMutex.RUnlock()
	return fake.newQueryCreatorArgsForCall[i].channel
}

func (fake *Support) NewQueryCreatorReturns(result1 endorser_test.QueryCreator, result2 error) {
	fake.NewQueryCreatorStub = nil
	fake.newQueryCreatorReturns = struct {
		result1 endorser_test.QueryCreator
		result2 error
	}{result1, result2}
}

func (fake *Support) NewQueryCreatorReturnsOnCall(i int, result1 endorser_test.QueryCreator, result2 error) {
	fake.NewQueryCreatorStub = nil
	if fake.newQueryCreatorReturnsOnCall == nil {
		fake.newQueryCreatorReturnsOnCall = make(map[int]struct {
			result1 endorser_test.QueryCreator
			result2 error
		})
	}
	fake.newQueryCreatorReturnsOnCall[i] = struct {
		result1 endorser_test.QueryCreator
		result2 error
	}{result1, result2}
}

func (fake *Support) EndorseWithPlugin(ctx endorser_test.Context) (*pb.ProposalResponse, error) {
	fake.endorseWithPluginMutex.Lock()
	ret, specificReturn := fake.endorseWithPluginReturnsOnCall[len(fake.endorseWithPluginArgsForCall)]
	fake.endorseWithPluginArgsForCall = append(fake.endorseWithPluginArgsForCall, struct {
		ctx endorser_test.Context
	}{ctx})
	fake.recordInvocation("EndorseWithPlugin", []interface{}{ctx})
	fake.endorseWithPluginMutex.Unlock()
	if fake.EndorseWithPluginStub != nil {
		return fake.EndorseWithPluginStub(ctx)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.endorseWithPluginReturns.result1, fake.endorseWithPluginReturns.result2
}

func (fake *Support) EndorseWithPluginCallCount() int {
	fake.endorseWithPluginMutex.RLock()
	defer fake.endorseWithPluginMutex.RUnlock()
	return len(fake.endorseWithPluginArgsForCall)
}

func (fake *Support) EndorseWithPluginArgsForCall(i int) endorser_test.Context {
	fake.endorseWithPluginMutex.RLock()
	defer fake.endorseWithPluginMutex.RUnlock()
	return fake.endorseWithPluginArgsForCall[i].ctx
}

func (fake *Support) EndorseWithPluginReturns(result1 *pb.ProposalResponse, result2 error) {
	fake.EndorseWithPluginStub = nil
	fake.endorseWithPluginReturns = struct {
		result1 *pb.ProposalResponse
		result2 error
	}{result1, result2}
}

func (fake *Support) EndorseWithPluginReturnsOnCall(i int, result1 *pb.ProposalResponse, result2 error) {
	fake.EndorseWithPluginStub = nil
	if fake.endorseWithPluginReturnsOnCall == nil {
		fake.endorseWithPluginReturnsOnCall = make(map[int]struct {
			result1 *pb.ProposalResponse
			result2 error
		})
	}
	fake.endorseWithPluginReturnsOnCall[i] = struct {
		result1 *pb.ProposalResponse
		result2 error
	}{result1, result2}
}

func (fake *Support) GetLedgerHeight(channelID string) (uint64, error) {
	fake.getLedgerHeightMutex.Lock()
	ret, specificReturn := fake.getLedgerHeightReturnsOnCall[len(fake.getLedgerHeightArgsForCall)]
	fake.getLedgerHeightArgsForCall = append(fake.getLedgerHeightArgsForCall, struct {
		channelID string
	}{channelID})
	fake.recordInvocation("GetLedgerHeight", []interface{}{channelID})
	fake.getLedgerHeightMutex.Unlock()
	if fake.GetLedgerHeightStub != nil {
		return fake.GetLedgerHeightStub(channelID)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getLedgerHeightReturns.result1, fake.getLedgerHeightReturns.result2
}

func (fake *Support) GetLedgerHeightCallCount() int {
	fake.getLedgerHeightMutex.RLock()
	defer fake.getLedgerHeightMutex.RUnlock()
	return len(fake.getLedgerHeightArgsForCall)
}

func (fake *Support) GetLedgerHeightArgsForCall(i int) string {
	fake.getLedgerHeightMutex.RLock()
	defer fake.getLedgerHeightMutex.RUnlock()
	return fake.getLedgerHeightArgsForCall[i].channelID
}

func (fake *Support) GetLedgerHeightReturns(result1 uint64, result2 error) {
	fake.GetLedgerHeightStub = nil
	fake.getLedgerHeightReturns = struct {
		result1 uint64
		result2 error
	}{result1, result2}
}

func (fake *Support) GetLedgerHeightReturnsOnCall(i int, result1 uint64, result2 error) {
	fake.GetLedgerHeightStub = nil
	if fake.getLedgerHeightReturnsOnCall == nil {
		fake.getLedgerHeightReturnsOnCall = make(map[int]struct {
			result1 uint64
			result2 error
		})
	}
	fake.getLedgerHeightReturnsOnCall[i] = struct {
		result1 uint64
		result2 error
	}{result1, result2}
}

func (fake *Support) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	fake.serializeMutex.RLock()
	defer fake.serializeMutex.RUnlock()
	fake.isSysCCAndNotInvokableExternalMutex.RLock()
	defer fake.isSysCCAndNotInvokableExternalMutex.RUnlock()
	fake.getTxSimulatorMutex.RLock()
	defer fake.getTxSimulatorMutex.RUnlock()
	fake.getHistoryQueryExecutorMutex.RLock()
	defer fake.getHistoryQueryExecutorMutex.RUnlock()
	fake.getTransactionByIDMutex.RLock()
	defer fake.getTransactionByIDMutex.RUnlock()
	fake.isSysCCMutex.RLock()
	defer fake.isSysCCMutex.RUnlock()
	fake.executeMutex.RLock()
	defer fake.executeMutex.RUnlock()
	fake.executeLegacyInitMutex.RLock()
	defer fake.executeLegacyInitMutex.RUnlock()
	fake.getChaincodeDefinitionMutex.RLock()
	defer fake.getChaincodeDefinitionMutex.RUnlock()
	fake.checkACLMutex.RLock()
	defer fake.checkACLMutex.RUnlock()
	fake.isJavaCCMutex.RLock()
	defer fake.isJavaCCMutex.RUnlock()
	fake.checkInstantiationPolicyMutex.RLock()
	defer fake.checkInstantiationPolicyMutex.RUnlock()
	fake.getChaincodeDeploymentSpecFSMutex.RLock()
	defer fake.getChaincodeDeploymentSpecFSMutex.RUnlock()
	fake.getApplicationConfigMutex.RLock()
	defer fake.getApplicationConfigMutex.RUnlock()
	fake.newQueryCreatorMutex.RLock()
	defer fake.newQueryCreatorMutex.RUnlock()
	fake.endorseWithPluginMutex.RLock()
	defer fake.endorseWithPluginMutex.RUnlock()
	fake.getLedgerHeightMutex.RLock()
	defer fake.getLedgerHeightMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *Support) recordInvocation(key string, args []interface{}) {
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