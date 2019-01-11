
package mock

import (
	"sync"

	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

type ChaincodeStub struct {
	GetArgsStub        func() [][]byte
	getArgsMutex       sync.RWMutex
	getArgsArgsForCall []struct{}
	getArgsReturns     struct {
		result1 [][]byte
	}
	getArgsReturnsOnCall map[int]struct {
		result1 [][]byte
	}
	GetStringArgsStub        func() []string
	getStringArgsMutex       sync.RWMutex
	getStringArgsArgsForCall []struct{}
	getStringArgsReturns     struct {
		result1 []string
	}
	getStringArgsReturnsOnCall map[int]struct {
		result1 []string
	}
	GetFunctionAndParametersStub        func() (string, []string)
	getFunctionAndParametersMutex       sync.RWMutex
	getFunctionAndParametersArgsForCall []struct{}
	getFunctionAndParametersReturns     struct {
		result1 string
		result2 []string
	}
	getFunctionAndParametersReturnsOnCall map[int]struct {
		result1 string
		result2 []string
	}
	GetArgsSliceStub        func() ([]byte, error)
	getArgsSliceMutex       sync.RWMutex
	getArgsSliceArgsForCall []struct{}
	getArgsSliceReturns     struct {
		result1 []byte
		result2 error
	}
	getArgsSliceReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetTxIDStub        func() string
	getTxIDMutex       sync.RWMutex
	getTxIDArgsForCall []struct{}
	getTxIDReturns     struct {
		result1 string
	}
	getTxIDReturnsOnCall map[int]struct {
		result1 string
	}
	GetChannelIDStub        func() string
	getChannelIDMutex       sync.RWMutex
	getChannelIDArgsForCall []struct{}
	getChannelIDReturns     struct {
		result1 string
	}
	getChannelIDReturnsOnCall map[int]struct {
		result1 string
	}
	InvokeChaincodeStub        func(chaincodeName string, args [][]byte, channel string) pb.Response
	invokeChaincodeMutex       sync.RWMutex
	invokeChaincodeArgsForCall []struct {
		chaincodeName string
		args          [][]byte
		channel       string
	}
	invokeChaincodeReturns struct {
		result1 pb.Response
	}
	invokeChaincodeReturnsOnCall map[int]struct {
		result1 pb.Response
	}
	GetStateStub        func(key string) ([]byte, error)
	getStateMutex       sync.RWMutex
	getStateArgsForCall []struct {
		key string
	}
	getStateReturns struct {
		result1 []byte
		result2 error
	}
	getStateReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	PutStateStub        func(key string, value []byte) error
	putStateMutex       sync.RWMutex
	putStateArgsForCall []struct {
		key   string
		value []byte
	}
	putStateReturns struct {
		result1 error
	}
	putStateReturnsOnCall map[int]struct {
		result1 error
	}
	DelStateStub        func(key string) error
	delStateMutex       sync.RWMutex
	delStateArgsForCall []struct {
		key string
	}
	delStateReturns struct {
		result1 error
	}
	delStateReturnsOnCall map[int]struct {
		result1 error
	}
	SetStateValidationParameterStub        func(key string, ep []byte) error
	setStateValidationParameterMutex       sync.RWMutex
	setStateValidationParameterArgsForCall []struct {
		key string
		ep  []byte
	}
	setStateValidationParameterReturns struct {
		result1 error
	}
	setStateValidationParameterReturnsOnCall map[int]struct {
		result1 error
	}
	GetStateValidationParameterStub        func(key string) ([]byte, error)
	getStateValidationParameterMutex       sync.RWMutex
	getStateValidationParameterArgsForCall []struct {
		key string
	}
	getStateValidationParameterReturns struct {
		result1 []byte
		result2 error
	}
	getStateValidationParameterReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetStateByRangeStub        func(startKey, endKey string) (shim.StateQueryIteratorInterface, error)
	getStateByRangeMutex       sync.RWMutex
	getStateByRangeArgsForCall []struct {
		startKey string
		endKey   string
	}
	getStateByRangeReturns struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	getStateByRangeReturnsOnCall map[int]struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	GetStateByRangeWithPaginationStub        func(startKey, endKey string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)
	getStateByRangeWithPaginationMutex       sync.RWMutex
	getStateByRangeWithPaginationArgsForCall []struct {
		startKey string
		endKey   string
		pageSize int32
		bookmark string
	}
	getStateByRangeWithPaginationReturns struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}
	getStateByRangeWithPaginationReturnsOnCall map[int]struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}
	GetStateByPartialCompositeKeyStub        func(objectType string, keys []string) (shim.StateQueryIteratorInterface, error)
	getStateByPartialCompositeKeyMutex       sync.RWMutex
	getStateByPartialCompositeKeyArgsForCall []struct {
		objectType string
		keys       []string
	}
	getStateByPartialCompositeKeyReturns struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	getStateByPartialCompositeKeyReturnsOnCall map[int]struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	GetStateByPartialCompositeKeyWithPaginationStub        func(objectType string, keys []string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)
	getStateByPartialCompositeKeyWithPaginationMutex       sync.RWMutex
	getStateByPartialCompositeKeyWithPaginationArgsForCall []struct {
		objectType string
		keys       []string
		pageSize   int32
		bookmark   string
	}
	getStateByPartialCompositeKeyWithPaginationReturns struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}
	getStateByPartialCompositeKeyWithPaginationReturnsOnCall map[int]struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}
	CreateCompositeKeyStub        func(objectType string, attributes []string) (string, error)
	createCompositeKeyMutex       sync.RWMutex
	createCompositeKeyArgsForCall []struct {
		objectType string
		attributes []string
	}
	createCompositeKeyReturns struct {
		result1 string
		result2 error
	}
	createCompositeKeyReturnsOnCall map[int]struct {
		result1 string
		result2 error
	}
	SplitCompositeKeyStub        func(compositeKey string) (string, []string, error)
	splitCompositeKeyMutex       sync.RWMutex
	splitCompositeKeyArgsForCall []struct {
		compositeKey string
	}
	splitCompositeKeyReturns struct {
		result1 string
		result2 []string
		result3 error
	}
	splitCompositeKeyReturnsOnCall map[int]struct {
		result1 string
		result2 []string
		result3 error
	}
	GetQueryResultStub        func(query string) (shim.StateQueryIteratorInterface, error)
	getQueryResultMutex       sync.RWMutex
	getQueryResultArgsForCall []struct {
		query string
	}
	getQueryResultReturns struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	getQueryResultReturnsOnCall map[int]struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	GetQueryResultWithPaginationStub        func(query string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error)
	getQueryResultWithPaginationMutex       sync.RWMutex
	getQueryResultWithPaginationArgsForCall []struct {
		query    string
		pageSize int32
		bookmark string
	}
	getQueryResultWithPaginationReturns struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}
	getQueryResultWithPaginationReturnsOnCall map[int]struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}
	GetHistoryForKeyStub        func(key string) (shim.HistoryQueryIteratorInterface, error)
	getHistoryForKeyMutex       sync.RWMutex
	getHistoryForKeyArgsForCall []struct {
		key string
	}
	getHistoryForKeyReturns struct {
		result1 shim.HistoryQueryIteratorInterface
		result2 error
	}
	getHistoryForKeyReturnsOnCall map[int]struct {
		result1 shim.HistoryQueryIteratorInterface
		result2 error
	}
	GetPrivateDataStub        func(collection, key string) ([]byte, error)
	getPrivateDataMutex       sync.RWMutex
	getPrivateDataArgsForCall []struct {
		collection string
		key        string
	}
	getPrivateDataReturns struct {
		result1 []byte
		result2 error
	}
	getPrivateDataReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetPrivateDataHashStub        func(collection, key string) ([]byte, error)
	getPrivateDataHashMutex       sync.RWMutex
	getPrivateDataHashArgsForCall []struct {
		collection string
		key        string
	}
	getPrivateDataHashReturns struct {
		result1 []byte
		result2 error
	}
	getPrivateDataHashReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	PutPrivateDataStub        func(collection string, key string, value []byte) error
	putPrivateDataMutex       sync.RWMutex
	putPrivateDataArgsForCall []struct {
		collection string
		key        string
		value      []byte
	}
	putPrivateDataReturns struct {
		result1 error
	}
	putPrivateDataReturnsOnCall map[int]struct {
		result1 error
	}
	DelPrivateDataStub        func(collection, key string) error
	delPrivateDataMutex       sync.RWMutex
	delPrivateDataArgsForCall []struct {
		collection string
		key        string
	}
	delPrivateDataReturns struct {
		result1 error
	}
	delPrivateDataReturnsOnCall map[int]struct {
		result1 error
	}
	SetPrivateDataValidationParameterStub        func(collection, key string, ep []byte) error
	setPrivateDataValidationParameterMutex       sync.RWMutex
	setPrivateDataValidationParameterArgsForCall []struct {
		collection string
		key        string
		ep         []byte
	}
	setPrivateDataValidationParameterReturns struct {
		result1 error
	}
	setPrivateDataValidationParameterReturnsOnCall map[int]struct {
		result1 error
	}
	GetPrivateDataValidationParameterStub        func(collection, key string) ([]byte, error)
	getPrivateDataValidationParameterMutex       sync.RWMutex
	getPrivateDataValidationParameterArgsForCall []struct {
		collection string
		key        string
	}
	getPrivateDataValidationParameterReturns struct {
		result1 []byte
		result2 error
	}
	getPrivateDataValidationParameterReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetPrivateDataByRangeStub        func(collection, startKey, endKey string) (shim.StateQueryIteratorInterface, error)
	getPrivateDataByRangeMutex       sync.RWMutex
	getPrivateDataByRangeArgsForCall []struct {
		collection string
		startKey   string
		endKey     string
	}
	getPrivateDataByRangeReturns struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	getPrivateDataByRangeReturnsOnCall map[int]struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	GetPrivateDataByPartialCompositeKeyStub        func(collection, objectType string, keys []string) (shim.StateQueryIteratorInterface, error)
	getPrivateDataByPartialCompositeKeyMutex       sync.RWMutex
	getPrivateDataByPartialCompositeKeyArgsForCall []struct {
		collection string
		objectType string
		keys       []string
	}
	getPrivateDataByPartialCompositeKeyReturns struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	getPrivateDataByPartialCompositeKeyReturnsOnCall map[int]struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	GetPrivateDataQueryResultStub        func(collection, query string) (shim.StateQueryIteratorInterface, error)
	getPrivateDataQueryResultMutex       sync.RWMutex
	getPrivateDataQueryResultArgsForCall []struct {
		collection string
		query      string
	}
	getPrivateDataQueryResultReturns struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	getPrivateDataQueryResultReturnsOnCall map[int]struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}
	GetCreatorStub        func() ([]byte, error)
	getCreatorMutex       sync.RWMutex
	getCreatorArgsForCall []struct{}
	getCreatorReturns     struct {
		result1 []byte
		result2 error
	}
	getCreatorReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetTransientStub        func() (map[string][]byte, error)
	getTransientMutex       sync.RWMutex
	getTransientArgsForCall []struct{}
	getTransientReturns     struct {
		result1 map[string][]byte
		result2 error
	}
	getTransientReturnsOnCall map[int]struct {
		result1 map[string][]byte
		result2 error
	}
	GetBindingStub        func() ([]byte, error)
	getBindingMutex       sync.RWMutex
	getBindingArgsForCall []struct{}
	getBindingReturns     struct {
		result1 []byte
		result2 error
	}
	getBindingReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	GetDecorationsStub        func() map[string][]byte
	getDecorationsMutex       sync.RWMutex
	getDecorationsArgsForCall []struct{}
	getDecorationsReturns     struct {
		result1 map[string][]byte
	}
	getDecorationsReturnsOnCall map[int]struct {
		result1 map[string][]byte
	}
	GetSignedProposalStub        func() (*pb.SignedProposal, error)
	getSignedProposalMutex       sync.RWMutex
	getSignedProposalArgsForCall []struct{}
	getSignedProposalReturns     struct {
		result1 *pb.SignedProposal
		result2 error
	}
	getSignedProposalReturnsOnCall map[int]struct {
		result1 *pb.SignedProposal
		result2 error
	}
	GetTxTimestampStub        func() (*timestamp.Timestamp, error)
	getTxTimestampMutex       sync.RWMutex
	getTxTimestampArgsForCall []struct{}
	getTxTimestampReturns     struct {
		result1 *timestamp.Timestamp
		result2 error
	}
	getTxTimestampReturnsOnCall map[int]struct {
		result1 *timestamp.Timestamp
		result2 error
	}
	SetEventStub        func(name string, payload []byte) error
	setEventMutex       sync.RWMutex
	setEventArgsForCall []struct {
		name    string
		payload []byte
	}
	setEventReturns struct {
		result1 error
	}
	setEventReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ChaincodeStub) GetArgs() [][]byte {
	fake.getArgsMutex.Lock()
	ret, specificReturn := fake.getArgsReturnsOnCall[len(fake.getArgsArgsForCall)]
	fake.getArgsArgsForCall = append(fake.getArgsArgsForCall, struct{}{})
	fake.recordInvocation("GetArgs", []interface{}{})
	fake.getArgsMutex.Unlock()
	if fake.GetArgsStub != nil {
		return fake.GetArgsStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getArgsReturns.result1
}

func (fake *ChaincodeStub) GetArgsCallCount() int {
	fake.getArgsMutex.RLock()
	defer fake.getArgsMutex.RUnlock()
	return len(fake.getArgsArgsForCall)
}

func (fake *ChaincodeStub) GetArgsReturns(result1 [][]byte) {
	fake.GetArgsStub = nil
	fake.getArgsReturns = struct {
		result1 [][]byte
	}{result1}
}

func (fake *ChaincodeStub) GetArgsReturnsOnCall(i int, result1 [][]byte) {
	fake.GetArgsStub = nil
	if fake.getArgsReturnsOnCall == nil {
		fake.getArgsReturnsOnCall = make(map[int]struct {
			result1 [][]byte
		})
	}
	fake.getArgsReturnsOnCall[i] = struct {
		result1 [][]byte
	}{result1}
}

func (fake *ChaincodeStub) GetStringArgs() []string {
	fake.getStringArgsMutex.Lock()
	ret, specificReturn := fake.getStringArgsReturnsOnCall[len(fake.getStringArgsArgsForCall)]
	fake.getStringArgsArgsForCall = append(fake.getStringArgsArgsForCall, struct{}{})
	fake.recordInvocation("GetStringArgs", []interface{}{})
	fake.getStringArgsMutex.Unlock()
	if fake.GetStringArgsStub != nil {
		return fake.GetStringArgsStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getStringArgsReturns.result1
}

func (fake *ChaincodeStub) GetStringArgsCallCount() int {
	fake.getStringArgsMutex.RLock()
	defer fake.getStringArgsMutex.RUnlock()
	return len(fake.getStringArgsArgsForCall)
}

func (fake *ChaincodeStub) GetStringArgsReturns(result1 []string) {
	fake.GetStringArgsStub = nil
	fake.getStringArgsReturns = struct {
		result1 []string
	}{result1}
}

func (fake *ChaincodeStub) GetStringArgsReturnsOnCall(i int, result1 []string) {
	fake.GetStringArgsStub = nil
	if fake.getStringArgsReturnsOnCall == nil {
		fake.getStringArgsReturnsOnCall = make(map[int]struct {
			result1 []string
		})
	}
	fake.getStringArgsReturnsOnCall[i] = struct {
		result1 []string
	}{result1}
}

func (fake *ChaincodeStub) GetFunctionAndParameters() (string, []string) {
	fake.getFunctionAndParametersMutex.Lock()
	ret, specificReturn := fake.getFunctionAndParametersReturnsOnCall[len(fake.getFunctionAndParametersArgsForCall)]
	fake.getFunctionAndParametersArgsForCall = append(fake.getFunctionAndParametersArgsForCall, struct{}{})
	fake.recordInvocation("GetFunctionAndParameters", []interface{}{})
	fake.getFunctionAndParametersMutex.Unlock()
	if fake.GetFunctionAndParametersStub != nil {
		return fake.GetFunctionAndParametersStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getFunctionAndParametersReturns.result1, fake.getFunctionAndParametersReturns.result2
}

func (fake *ChaincodeStub) GetFunctionAndParametersCallCount() int {
	fake.getFunctionAndParametersMutex.RLock()
	defer fake.getFunctionAndParametersMutex.RUnlock()
	return len(fake.getFunctionAndParametersArgsForCall)
}

func (fake *ChaincodeStub) GetFunctionAndParametersReturns(result1 string, result2 []string) {
	fake.GetFunctionAndParametersStub = nil
	fake.getFunctionAndParametersReturns = struct {
		result1 string
		result2 []string
	}{result1, result2}
}

func (fake *ChaincodeStub) GetFunctionAndParametersReturnsOnCall(i int, result1 string, result2 []string) {
	fake.GetFunctionAndParametersStub = nil
	if fake.getFunctionAndParametersReturnsOnCall == nil {
		fake.getFunctionAndParametersReturnsOnCall = make(map[int]struct {
			result1 string
			result2 []string
		})
	}
	fake.getFunctionAndParametersReturnsOnCall[i] = struct {
		result1 string
		result2 []string
	}{result1, result2}
}

func (fake *ChaincodeStub) GetArgsSlice() ([]byte, error) {
	fake.getArgsSliceMutex.Lock()
	ret, specificReturn := fake.getArgsSliceReturnsOnCall[len(fake.getArgsSliceArgsForCall)]
	fake.getArgsSliceArgsForCall = append(fake.getArgsSliceArgsForCall, struct{}{})
	fake.recordInvocation("GetArgsSlice", []interface{}{})
	fake.getArgsSliceMutex.Unlock()
	if fake.GetArgsSliceStub != nil {
		return fake.GetArgsSliceStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getArgsSliceReturns.result1, fake.getArgsSliceReturns.result2
}

func (fake *ChaincodeStub) GetArgsSliceCallCount() int {
	fake.getArgsSliceMutex.RLock()
	defer fake.getArgsSliceMutex.RUnlock()
	return len(fake.getArgsSliceArgsForCall)
}

func (fake *ChaincodeStub) GetArgsSliceReturns(result1 []byte, result2 error) {
	fake.GetArgsSliceStub = nil
	fake.getArgsSliceReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetArgsSliceReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetArgsSliceStub = nil
	if fake.getArgsSliceReturnsOnCall == nil {
		fake.getArgsSliceReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getArgsSliceReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetTxID() string {
	fake.getTxIDMutex.Lock()
	ret, specificReturn := fake.getTxIDReturnsOnCall[len(fake.getTxIDArgsForCall)]
	fake.getTxIDArgsForCall = append(fake.getTxIDArgsForCall, struct{}{})
	fake.recordInvocation("GetTxID", []interface{}{})
	fake.getTxIDMutex.Unlock()
	if fake.GetTxIDStub != nil {
		return fake.GetTxIDStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getTxIDReturns.result1
}

func (fake *ChaincodeStub) GetTxIDCallCount() int {
	fake.getTxIDMutex.RLock()
	defer fake.getTxIDMutex.RUnlock()
	return len(fake.getTxIDArgsForCall)
}

func (fake *ChaincodeStub) GetTxIDReturns(result1 string) {
	fake.GetTxIDStub = nil
	fake.getTxIDReturns = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeStub) GetTxIDReturnsOnCall(i int, result1 string) {
	fake.GetTxIDStub = nil
	if fake.getTxIDReturnsOnCall == nil {
		fake.getTxIDReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getTxIDReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeStub) GetChannelID() string {
	fake.getChannelIDMutex.Lock()
	ret, specificReturn := fake.getChannelIDReturnsOnCall[len(fake.getChannelIDArgsForCall)]
	fake.getChannelIDArgsForCall = append(fake.getChannelIDArgsForCall, struct{}{})
	fake.recordInvocation("GetChannelID", []interface{}{})
	fake.getChannelIDMutex.Unlock()
	if fake.GetChannelIDStub != nil {
		return fake.GetChannelIDStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getChannelIDReturns.result1
}

func (fake *ChaincodeStub) GetChannelIDCallCount() int {
	fake.getChannelIDMutex.RLock()
	defer fake.getChannelIDMutex.RUnlock()
	return len(fake.getChannelIDArgsForCall)
}

func (fake *ChaincodeStub) GetChannelIDReturns(result1 string) {
	fake.GetChannelIDStub = nil
	fake.getChannelIDReturns = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeStub) GetChannelIDReturnsOnCall(i int, result1 string) {
	fake.GetChannelIDStub = nil
	if fake.getChannelIDReturnsOnCall == nil {
		fake.getChannelIDReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.getChannelIDReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *ChaincodeStub) InvokeChaincode(chaincodeName string, args [][]byte, channel string) pb.Response {
	var argsCopy [][]byte
	if args != nil {
		argsCopy = make([][]byte, len(args))
		copy(argsCopy, args)
	}
	fake.invokeChaincodeMutex.Lock()
	ret, specificReturn := fake.invokeChaincodeReturnsOnCall[len(fake.invokeChaincodeArgsForCall)]
	fake.invokeChaincodeArgsForCall = append(fake.invokeChaincodeArgsForCall, struct {
		chaincodeName string
		args          [][]byte
		channel       string
	}{chaincodeName, argsCopy, channel})
	fake.recordInvocation("InvokeChaincode", []interface{}{chaincodeName, argsCopy, channel})
	fake.invokeChaincodeMutex.Unlock()
	if fake.InvokeChaincodeStub != nil {
		return fake.InvokeChaincodeStub(chaincodeName, args, channel)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.invokeChaincodeReturns.result1
}

func (fake *ChaincodeStub) InvokeChaincodeCallCount() int {
	fake.invokeChaincodeMutex.RLock()
	defer fake.invokeChaincodeMutex.RUnlock()
	return len(fake.invokeChaincodeArgsForCall)
}

func (fake *ChaincodeStub) InvokeChaincodeArgsForCall(i int) (string, [][]byte, string) {
	fake.invokeChaincodeMutex.RLock()
	defer fake.invokeChaincodeMutex.RUnlock()
	return fake.invokeChaincodeArgsForCall[i].chaincodeName, fake.invokeChaincodeArgsForCall[i].args, fake.invokeChaincodeArgsForCall[i].channel
}

func (fake *ChaincodeStub) InvokeChaincodeReturns(result1 pb.Response) {
	fake.InvokeChaincodeStub = nil
	fake.invokeChaincodeReturns = struct {
		result1 pb.Response
	}{result1}
}

func (fake *ChaincodeStub) InvokeChaincodeReturnsOnCall(i int, result1 pb.Response) {
	fake.InvokeChaincodeStub = nil
	if fake.invokeChaincodeReturnsOnCall == nil {
		fake.invokeChaincodeReturnsOnCall = make(map[int]struct {
			result1 pb.Response
		})
	}
	fake.invokeChaincodeReturnsOnCall[i] = struct {
		result1 pb.Response
	}{result1}
}

func (fake *ChaincodeStub) GetState(key string) ([]byte, error) {
	fake.getStateMutex.Lock()
	ret, specificReturn := fake.getStateReturnsOnCall[len(fake.getStateArgsForCall)]
	fake.getStateArgsForCall = append(fake.getStateArgsForCall, struct {
		key string
	}{key})
	fake.recordInvocation("GetState", []interface{}{key})
	fake.getStateMutex.Unlock()
	if fake.GetStateStub != nil {
		return fake.GetStateStub(key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateReturns.result1, fake.getStateReturns.result2
}

func (fake *ChaincodeStub) GetStateCallCount() int {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return len(fake.getStateArgsForCall)
}

func (fake *ChaincodeStub) GetStateArgsForCall(i int) string {
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	return fake.getStateArgsForCall[i].key
}

func (fake *ChaincodeStub) GetStateReturns(result1 []byte, result2 error) {
	fake.GetStateStub = nil
	fake.getStateReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetStateReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetStateStub = nil
	if fake.getStateReturnsOnCall == nil {
		fake.getStateReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getStateReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) PutState(key string, value []byte) error {
	var valueCopy []byte
	if value != nil {
		valueCopy = make([]byte, len(value))
		copy(valueCopy, value)
	}
	fake.putStateMutex.Lock()
	ret, specificReturn := fake.putStateReturnsOnCall[len(fake.putStateArgsForCall)]
	fake.putStateArgsForCall = append(fake.putStateArgsForCall, struct {
		key   string
		value []byte
	}{key, valueCopy})
	fake.recordInvocation("PutState", []interface{}{key, valueCopy})
	fake.putStateMutex.Unlock()
	if fake.PutStateStub != nil {
		return fake.PutStateStub(key, value)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.putStateReturns.result1
}

func (fake *ChaincodeStub) PutStateCallCount() int {
	fake.putStateMutex.RLock()
	defer fake.putStateMutex.RUnlock()
	return len(fake.putStateArgsForCall)
}

func (fake *ChaincodeStub) PutStateArgsForCall(i int) (string, []byte) {
	fake.putStateMutex.RLock()
	defer fake.putStateMutex.RUnlock()
	return fake.putStateArgsForCall[i].key, fake.putStateArgsForCall[i].value
}

func (fake *ChaincodeStub) PutStateReturns(result1 error) {
	fake.PutStateStub = nil
	fake.putStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) PutStateReturnsOnCall(i int, result1 error) {
	fake.PutStateStub = nil
	if fake.putStateReturnsOnCall == nil {
		fake.putStateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.putStateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) DelState(key string) error {
	fake.delStateMutex.Lock()
	ret, specificReturn := fake.delStateReturnsOnCall[len(fake.delStateArgsForCall)]
	fake.delStateArgsForCall = append(fake.delStateArgsForCall, struct {
		key string
	}{key})
	fake.recordInvocation("DelState", []interface{}{key})
	fake.delStateMutex.Unlock()
	if fake.DelStateStub != nil {
		return fake.DelStateStub(key)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.delStateReturns.result1
}

func (fake *ChaincodeStub) DelStateCallCount() int {
	fake.delStateMutex.RLock()
	defer fake.delStateMutex.RUnlock()
	return len(fake.delStateArgsForCall)
}

func (fake *ChaincodeStub) DelStateArgsForCall(i int) string {
	fake.delStateMutex.RLock()
	defer fake.delStateMutex.RUnlock()
	return fake.delStateArgsForCall[i].key
}

func (fake *ChaincodeStub) DelStateReturns(result1 error) {
	fake.DelStateStub = nil
	fake.delStateReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) DelStateReturnsOnCall(i int, result1 error) {
	fake.DelStateStub = nil
	if fake.delStateReturnsOnCall == nil {
		fake.delStateReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.delStateReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) SetStateValidationParameter(key string, ep []byte) error {
	var epCopy []byte
	if ep != nil {
		epCopy = make([]byte, len(ep))
		copy(epCopy, ep)
	}
	fake.setStateValidationParameterMutex.Lock()
	ret, specificReturn := fake.setStateValidationParameterReturnsOnCall[len(fake.setStateValidationParameterArgsForCall)]
	fake.setStateValidationParameterArgsForCall = append(fake.setStateValidationParameterArgsForCall, struct {
		key string
		ep  []byte
	}{key, epCopy})
	fake.recordInvocation("SetStateValidationParameter", []interface{}{key, epCopy})
	fake.setStateValidationParameterMutex.Unlock()
	if fake.SetStateValidationParameterStub != nil {
		return fake.SetStateValidationParameterStub(key, ep)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setStateValidationParameterReturns.result1
}

func (fake *ChaincodeStub) SetStateValidationParameterCallCount() int {
	fake.setStateValidationParameterMutex.RLock()
	defer fake.setStateValidationParameterMutex.RUnlock()
	return len(fake.setStateValidationParameterArgsForCall)
}

func (fake *ChaincodeStub) SetStateValidationParameterArgsForCall(i int) (string, []byte) {
	fake.setStateValidationParameterMutex.RLock()
	defer fake.setStateValidationParameterMutex.RUnlock()
	return fake.setStateValidationParameterArgsForCall[i].key, fake.setStateValidationParameterArgsForCall[i].ep
}

func (fake *ChaincodeStub) SetStateValidationParameterReturns(result1 error) {
	fake.SetStateValidationParameterStub = nil
	fake.setStateValidationParameterReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) SetStateValidationParameterReturnsOnCall(i int, result1 error) {
	fake.SetStateValidationParameterStub = nil
	if fake.setStateValidationParameterReturnsOnCall == nil {
		fake.setStateValidationParameterReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setStateValidationParameterReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) GetStateValidationParameter(key string) ([]byte, error) {
	fake.getStateValidationParameterMutex.Lock()
	ret, specificReturn := fake.getStateValidationParameterReturnsOnCall[len(fake.getStateValidationParameterArgsForCall)]
	fake.getStateValidationParameterArgsForCall = append(fake.getStateValidationParameterArgsForCall, struct {
		key string
	}{key})
	fake.recordInvocation("GetStateValidationParameter", []interface{}{key})
	fake.getStateValidationParameterMutex.Unlock()
	if fake.GetStateValidationParameterStub != nil {
		return fake.GetStateValidationParameterStub(key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateValidationParameterReturns.result1, fake.getStateValidationParameterReturns.result2
}

func (fake *ChaincodeStub) GetStateValidationParameterCallCount() int {
	fake.getStateValidationParameterMutex.RLock()
	defer fake.getStateValidationParameterMutex.RUnlock()
	return len(fake.getStateValidationParameterArgsForCall)
}

func (fake *ChaincodeStub) GetStateValidationParameterArgsForCall(i int) string {
	fake.getStateValidationParameterMutex.RLock()
	defer fake.getStateValidationParameterMutex.RUnlock()
	return fake.getStateValidationParameterArgsForCall[i].key
}

func (fake *ChaincodeStub) GetStateValidationParameterReturns(result1 []byte, result2 error) {
	fake.GetStateValidationParameterStub = nil
	fake.getStateValidationParameterReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetStateValidationParameterReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetStateValidationParameterStub = nil
	if fake.getStateValidationParameterReturnsOnCall == nil {
		fake.getStateValidationParameterReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getStateValidationParameterReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetStateByRange(startKey string, endKey string) (shim.StateQueryIteratorInterface, error) {
	fake.getStateByRangeMutex.Lock()
	ret, specificReturn := fake.getStateByRangeReturnsOnCall[len(fake.getStateByRangeArgsForCall)]
	fake.getStateByRangeArgsForCall = append(fake.getStateByRangeArgsForCall, struct {
		startKey string
		endKey   string
	}{startKey, endKey})
	fake.recordInvocation("GetStateByRange", []interface{}{startKey, endKey})
	fake.getStateByRangeMutex.Unlock()
	if fake.GetStateByRangeStub != nil {
		return fake.GetStateByRangeStub(startKey, endKey)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateByRangeReturns.result1, fake.getStateByRangeReturns.result2
}

func (fake *ChaincodeStub) GetStateByRangeCallCount() int {
	fake.getStateByRangeMutex.RLock()
	defer fake.getStateByRangeMutex.RUnlock()
	return len(fake.getStateByRangeArgsForCall)
}

func (fake *ChaincodeStub) GetStateByRangeArgsForCall(i int) (string, string) {
	fake.getStateByRangeMutex.RLock()
	defer fake.getStateByRangeMutex.RUnlock()
	return fake.getStateByRangeArgsForCall[i].startKey, fake.getStateByRangeArgsForCall[i].endKey
}

func (fake *ChaincodeStub) GetStateByRangeReturns(result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetStateByRangeStub = nil
	fake.getStateByRangeReturns = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetStateByRangeReturnsOnCall(i int, result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetStateByRangeStub = nil
	if fake.getStateByRangeReturnsOnCall == nil {
		fake.getStateByRangeReturnsOnCall = make(map[int]struct {
			result1 shim.StateQueryIteratorInterface
			result2 error
		})
	}
	fake.getStateByRangeReturnsOnCall[i] = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetStateByRangeWithPagination(startKey string, endKey string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	fake.getStateByRangeWithPaginationMutex.Lock()
	ret, specificReturn := fake.getStateByRangeWithPaginationReturnsOnCall[len(fake.getStateByRangeWithPaginationArgsForCall)]
	fake.getStateByRangeWithPaginationArgsForCall = append(fake.getStateByRangeWithPaginationArgsForCall, struct {
		startKey string
		endKey   string
		pageSize int32
		bookmark string
	}{startKey, endKey, pageSize, bookmark})
	fake.recordInvocation("GetStateByRangeWithPagination", []interface{}{startKey, endKey, pageSize, bookmark})
	fake.getStateByRangeWithPaginationMutex.Unlock()
	if fake.GetStateByRangeWithPaginationStub != nil {
		return fake.GetStateByRangeWithPaginationStub(startKey, endKey, pageSize, bookmark)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.getStateByRangeWithPaginationReturns.result1, fake.getStateByRangeWithPaginationReturns.result2, fake.getStateByRangeWithPaginationReturns.result3
}

func (fake *ChaincodeStub) GetStateByRangeWithPaginationCallCount() int {
	fake.getStateByRangeWithPaginationMutex.RLock()
	defer fake.getStateByRangeWithPaginationMutex.RUnlock()
	return len(fake.getStateByRangeWithPaginationArgsForCall)
}

func (fake *ChaincodeStub) GetStateByRangeWithPaginationArgsForCall(i int) (string, string, int32, string) {
	fake.getStateByRangeWithPaginationMutex.RLock()
	defer fake.getStateByRangeWithPaginationMutex.RUnlock()
	return fake.getStateByRangeWithPaginationArgsForCall[i].startKey, fake.getStateByRangeWithPaginationArgsForCall[i].endKey, fake.getStateByRangeWithPaginationArgsForCall[i].pageSize, fake.getStateByRangeWithPaginationArgsForCall[i].bookmark
}

func (fake *ChaincodeStub) GetStateByRangeWithPaginationReturns(result1 shim.StateQueryIteratorInterface, result2 *pb.QueryResponseMetadata, result3 error) {
	fake.GetStateByRangeWithPaginationStub = nil
	fake.getStateByRangeWithPaginationReturns = struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeStub) GetStateByRangeWithPaginationReturnsOnCall(i int, result1 shim.StateQueryIteratorInterface, result2 *pb.QueryResponseMetadata, result3 error) {
	fake.GetStateByRangeWithPaginationStub = nil
	if fake.getStateByRangeWithPaginationReturnsOnCall == nil {
		fake.getStateByRangeWithPaginationReturnsOnCall = make(map[int]struct {
			result1 shim.StateQueryIteratorInterface
			result2 *pb.QueryResponseMetadata
			result3 error
		})
	}
	fake.getStateByRangeWithPaginationReturnsOnCall[i] = struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKey(objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	var keysCopy []string
	if keys != nil {
		keysCopy = make([]string, len(keys))
		copy(keysCopy, keys)
	}
	fake.getStateByPartialCompositeKeyMutex.Lock()
	ret, specificReturn := fake.getStateByPartialCompositeKeyReturnsOnCall[len(fake.getStateByPartialCompositeKeyArgsForCall)]
	fake.getStateByPartialCompositeKeyArgsForCall = append(fake.getStateByPartialCompositeKeyArgsForCall, struct {
		objectType string
		keys       []string
	}{objectType, keysCopy})
	fake.recordInvocation("GetStateByPartialCompositeKey", []interface{}{objectType, keysCopy})
	fake.getStateByPartialCompositeKeyMutex.Unlock()
	if fake.GetStateByPartialCompositeKeyStub != nil {
		return fake.GetStateByPartialCompositeKeyStub(objectType, keys)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getStateByPartialCompositeKeyReturns.result1, fake.getStateByPartialCompositeKeyReturns.result2
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKeyCallCount() int {
	fake.getStateByPartialCompositeKeyMutex.RLock()
	defer fake.getStateByPartialCompositeKeyMutex.RUnlock()
	return len(fake.getStateByPartialCompositeKeyArgsForCall)
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKeyArgsForCall(i int) (string, []string) {
	fake.getStateByPartialCompositeKeyMutex.RLock()
	defer fake.getStateByPartialCompositeKeyMutex.RUnlock()
	return fake.getStateByPartialCompositeKeyArgsForCall[i].objectType, fake.getStateByPartialCompositeKeyArgsForCall[i].keys
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKeyReturns(result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetStateByPartialCompositeKeyStub = nil
	fake.getStateByPartialCompositeKeyReturns = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKeyReturnsOnCall(i int, result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetStateByPartialCompositeKeyStub = nil
	if fake.getStateByPartialCompositeKeyReturnsOnCall == nil {
		fake.getStateByPartialCompositeKeyReturnsOnCall = make(map[int]struct {
			result1 shim.StateQueryIteratorInterface
			result2 error
		})
	}
	fake.getStateByPartialCompositeKeyReturnsOnCall[i] = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKeyWithPagination(objectType string, keys []string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	var keysCopy []string
	if keys != nil {
		keysCopy = make([]string, len(keys))
		copy(keysCopy, keys)
	}
	fake.getStateByPartialCompositeKeyWithPaginationMutex.Lock()
	ret, specificReturn := fake.getStateByPartialCompositeKeyWithPaginationReturnsOnCall[len(fake.getStateByPartialCompositeKeyWithPaginationArgsForCall)]
	fake.getStateByPartialCompositeKeyWithPaginationArgsForCall = append(fake.getStateByPartialCompositeKeyWithPaginationArgsForCall, struct {
		objectType string
		keys       []string
		pageSize   int32
		bookmark   string
	}{objectType, keysCopy, pageSize, bookmark})
	fake.recordInvocation("GetStateByPartialCompositeKeyWithPagination", []interface{}{objectType, keysCopy, pageSize, bookmark})
	fake.getStateByPartialCompositeKeyWithPaginationMutex.Unlock()
	if fake.GetStateByPartialCompositeKeyWithPaginationStub != nil {
		return fake.GetStateByPartialCompositeKeyWithPaginationStub(objectType, keys, pageSize, bookmark)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.getStateByPartialCompositeKeyWithPaginationReturns.result1, fake.getStateByPartialCompositeKeyWithPaginationReturns.result2, fake.getStateByPartialCompositeKeyWithPaginationReturns.result3
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKeyWithPaginationCallCount() int {
	fake.getStateByPartialCompositeKeyWithPaginationMutex.RLock()
	defer fake.getStateByPartialCompositeKeyWithPaginationMutex.RUnlock()
	return len(fake.getStateByPartialCompositeKeyWithPaginationArgsForCall)
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKeyWithPaginationArgsForCall(i int) (string, []string, int32, string) {
	fake.getStateByPartialCompositeKeyWithPaginationMutex.RLock()
	defer fake.getStateByPartialCompositeKeyWithPaginationMutex.RUnlock()
	return fake.getStateByPartialCompositeKeyWithPaginationArgsForCall[i].objectType, fake.getStateByPartialCompositeKeyWithPaginationArgsForCall[i].keys, fake.getStateByPartialCompositeKeyWithPaginationArgsForCall[i].pageSize, fake.getStateByPartialCompositeKeyWithPaginationArgsForCall[i].bookmark
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKeyWithPaginationReturns(result1 shim.StateQueryIteratorInterface, result2 *pb.QueryResponseMetadata, result3 error) {
	fake.GetStateByPartialCompositeKeyWithPaginationStub = nil
	fake.getStateByPartialCompositeKeyWithPaginationReturns = struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeStub) GetStateByPartialCompositeKeyWithPaginationReturnsOnCall(i int, result1 shim.StateQueryIteratorInterface, result2 *pb.QueryResponseMetadata, result3 error) {
	fake.GetStateByPartialCompositeKeyWithPaginationStub = nil
	if fake.getStateByPartialCompositeKeyWithPaginationReturnsOnCall == nil {
		fake.getStateByPartialCompositeKeyWithPaginationReturnsOnCall = make(map[int]struct {
			result1 shim.StateQueryIteratorInterface
			result2 *pb.QueryResponseMetadata
			result3 error
		})
	}
	fake.getStateByPartialCompositeKeyWithPaginationReturnsOnCall[i] = struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeStub) CreateCompositeKey(objectType string, attributes []string) (string, error) {
	var attributesCopy []string
	if attributes != nil {
		attributesCopy = make([]string, len(attributes))
		copy(attributesCopy, attributes)
	}
	fake.createCompositeKeyMutex.Lock()
	ret, specificReturn := fake.createCompositeKeyReturnsOnCall[len(fake.createCompositeKeyArgsForCall)]
	fake.createCompositeKeyArgsForCall = append(fake.createCompositeKeyArgsForCall, struct {
		objectType string
		attributes []string
	}{objectType, attributesCopy})
	fake.recordInvocation("CreateCompositeKey", []interface{}{objectType, attributesCopy})
	fake.createCompositeKeyMutex.Unlock()
	if fake.CreateCompositeKeyStub != nil {
		return fake.CreateCompositeKeyStub(objectType, attributes)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.createCompositeKeyReturns.result1, fake.createCompositeKeyReturns.result2
}

func (fake *ChaincodeStub) CreateCompositeKeyCallCount() int {
	fake.createCompositeKeyMutex.RLock()
	defer fake.createCompositeKeyMutex.RUnlock()
	return len(fake.createCompositeKeyArgsForCall)
}

func (fake *ChaincodeStub) CreateCompositeKeyArgsForCall(i int) (string, []string) {
	fake.createCompositeKeyMutex.RLock()
	defer fake.createCompositeKeyMutex.RUnlock()
	return fake.createCompositeKeyArgsForCall[i].objectType, fake.createCompositeKeyArgsForCall[i].attributes
}

func (fake *ChaincodeStub) CreateCompositeKeyReturns(result1 string, result2 error) {
	fake.CreateCompositeKeyStub = nil
	fake.createCompositeKeyReturns = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) CreateCompositeKeyReturnsOnCall(i int, result1 string, result2 error) {
	fake.CreateCompositeKeyStub = nil
	if fake.createCompositeKeyReturnsOnCall == nil {
		fake.createCompositeKeyReturnsOnCall = make(map[int]struct {
			result1 string
			result2 error
		})
	}
	fake.createCompositeKeyReturnsOnCall[i] = struct {
		result1 string
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) SplitCompositeKey(compositeKey string) (string, []string, error) {
	fake.splitCompositeKeyMutex.Lock()
	ret, specificReturn := fake.splitCompositeKeyReturnsOnCall[len(fake.splitCompositeKeyArgsForCall)]
	fake.splitCompositeKeyArgsForCall = append(fake.splitCompositeKeyArgsForCall, struct {
		compositeKey string
	}{compositeKey})
	fake.recordInvocation("SplitCompositeKey", []interface{}{compositeKey})
	fake.splitCompositeKeyMutex.Unlock()
	if fake.SplitCompositeKeyStub != nil {
		return fake.SplitCompositeKeyStub(compositeKey)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.splitCompositeKeyReturns.result1, fake.splitCompositeKeyReturns.result2, fake.splitCompositeKeyReturns.result3
}

func (fake *ChaincodeStub) SplitCompositeKeyCallCount() int {
	fake.splitCompositeKeyMutex.RLock()
	defer fake.splitCompositeKeyMutex.RUnlock()
	return len(fake.splitCompositeKeyArgsForCall)
}

func (fake *ChaincodeStub) SplitCompositeKeyArgsForCall(i int) string {
	fake.splitCompositeKeyMutex.RLock()
	defer fake.splitCompositeKeyMutex.RUnlock()
	return fake.splitCompositeKeyArgsForCall[i].compositeKey
}

func (fake *ChaincodeStub) SplitCompositeKeyReturns(result1 string, result2 []string, result3 error) {
	fake.SplitCompositeKeyStub = nil
	fake.splitCompositeKeyReturns = struct {
		result1 string
		result2 []string
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeStub) SplitCompositeKeyReturnsOnCall(i int, result1 string, result2 []string, result3 error) {
	fake.SplitCompositeKeyStub = nil
	if fake.splitCompositeKeyReturnsOnCall == nil {
		fake.splitCompositeKeyReturnsOnCall = make(map[int]struct {
			result1 string
			result2 []string
			result3 error
		})
	}
	fake.splitCompositeKeyReturnsOnCall[i] = struct {
		result1 string
		result2 []string
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeStub) GetQueryResult(query string) (shim.StateQueryIteratorInterface, error) {
	fake.getQueryResultMutex.Lock()
	ret, specificReturn := fake.getQueryResultReturnsOnCall[len(fake.getQueryResultArgsForCall)]
	fake.getQueryResultArgsForCall = append(fake.getQueryResultArgsForCall, struct {
		query string
	}{query})
	fake.recordInvocation("GetQueryResult", []interface{}{query})
	fake.getQueryResultMutex.Unlock()
	if fake.GetQueryResultStub != nil {
		return fake.GetQueryResultStub(query)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getQueryResultReturns.result1, fake.getQueryResultReturns.result2
}

func (fake *ChaincodeStub) GetQueryResultCallCount() int {
	fake.getQueryResultMutex.RLock()
	defer fake.getQueryResultMutex.RUnlock()
	return len(fake.getQueryResultArgsForCall)
}

func (fake *ChaincodeStub) GetQueryResultArgsForCall(i int) string {
	fake.getQueryResultMutex.RLock()
	defer fake.getQueryResultMutex.RUnlock()
	return fake.getQueryResultArgsForCall[i].query
}

func (fake *ChaincodeStub) GetQueryResultReturns(result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetQueryResultStub = nil
	fake.getQueryResultReturns = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetQueryResultReturnsOnCall(i int, result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetQueryResultStub = nil
	if fake.getQueryResultReturnsOnCall == nil {
		fake.getQueryResultReturnsOnCall = make(map[int]struct {
			result1 shim.StateQueryIteratorInterface
			result2 error
		})
	}
	fake.getQueryResultReturnsOnCall[i] = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetQueryResultWithPagination(query string, pageSize int32, bookmark string) (shim.StateQueryIteratorInterface, *pb.QueryResponseMetadata, error) {
	fake.getQueryResultWithPaginationMutex.Lock()
	ret, specificReturn := fake.getQueryResultWithPaginationReturnsOnCall[len(fake.getQueryResultWithPaginationArgsForCall)]
	fake.getQueryResultWithPaginationArgsForCall = append(fake.getQueryResultWithPaginationArgsForCall, struct {
		query    string
		pageSize int32
		bookmark string
	}{query, pageSize, bookmark})
	fake.recordInvocation("GetQueryResultWithPagination", []interface{}{query, pageSize, bookmark})
	fake.getQueryResultWithPaginationMutex.Unlock()
	if fake.GetQueryResultWithPaginationStub != nil {
		return fake.GetQueryResultWithPaginationStub(query, pageSize, bookmark)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.getQueryResultWithPaginationReturns.result1, fake.getQueryResultWithPaginationReturns.result2, fake.getQueryResultWithPaginationReturns.result3
}

func (fake *ChaincodeStub) GetQueryResultWithPaginationCallCount() int {
	fake.getQueryResultWithPaginationMutex.RLock()
	defer fake.getQueryResultWithPaginationMutex.RUnlock()
	return len(fake.getQueryResultWithPaginationArgsForCall)
}

func (fake *ChaincodeStub) GetQueryResultWithPaginationArgsForCall(i int) (string, int32, string) {
	fake.getQueryResultWithPaginationMutex.RLock()
	defer fake.getQueryResultWithPaginationMutex.RUnlock()
	return fake.getQueryResultWithPaginationArgsForCall[i].query, fake.getQueryResultWithPaginationArgsForCall[i].pageSize, fake.getQueryResultWithPaginationArgsForCall[i].bookmark
}

func (fake *ChaincodeStub) GetQueryResultWithPaginationReturns(result1 shim.StateQueryIteratorInterface, result2 *pb.QueryResponseMetadata, result3 error) {
	fake.GetQueryResultWithPaginationStub = nil
	fake.getQueryResultWithPaginationReturns = struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeStub) GetQueryResultWithPaginationReturnsOnCall(i int, result1 shim.StateQueryIteratorInterface, result2 *pb.QueryResponseMetadata, result3 error) {
	fake.GetQueryResultWithPaginationStub = nil
	if fake.getQueryResultWithPaginationReturnsOnCall == nil {
		fake.getQueryResultWithPaginationReturnsOnCall = make(map[int]struct {
			result1 shim.StateQueryIteratorInterface
			result2 *pb.QueryResponseMetadata
			result3 error
		})
	}
	fake.getQueryResultWithPaginationReturnsOnCall[i] = struct {
		result1 shim.StateQueryIteratorInterface
		result2 *pb.QueryResponseMetadata
		result3 error
	}{result1, result2, result3}
}

func (fake *ChaincodeStub) GetHistoryForKey(key string) (shim.HistoryQueryIteratorInterface, error) {
	fake.getHistoryForKeyMutex.Lock()
	ret, specificReturn := fake.getHistoryForKeyReturnsOnCall[len(fake.getHistoryForKeyArgsForCall)]
	fake.getHistoryForKeyArgsForCall = append(fake.getHistoryForKeyArgsForCall, struct {
		key string
	}{key})
	fake.recordInvocation("GetHistoryForKey", []interface{}{key})
	fake.getHistoryForKeyMutex.Unlock()
	if fake.GetHistoryForKeyStub != nil {
		return fake.GetHistoryForKeyStub(key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getHistoryForKeyReturns.result1, fake.getHistoryForKeyReturns.result2
}

func (fake *ChaincodeStub) GetHistoryForKeyCallCount() int {
	fake.getHistoryForKeyMutex.RLock()
	defer fake.getHistoryForKeyMutex.RUnlock()
	return len(fake.getHistoryForKeyArgsForCall)
}

func (fake *ChaincodeStub) GetHistoryForKeyArgsForCall(i int) string {
	fake.getHistoryForKeyMutex.RLock()
	defer fake.getHistoryForKeyMutex.RUnlock()
	return fake.getHistoryForKeyArgsForCall[i].key
}

func (fake *ChaincodeStub) GetHistoryForKeyReturns(result1 shim.HistoryQueryIteratorInterface, result2 error) {
	fake.GetHistoryForKeyStub = nil
	fake.getHistoryForKeyReturns = struct {
		result1 shim.HistoryQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetHistoryForKeyReturnsOnCall(i int, result1 shim.HistoryQueryIteratorInterface, result2 error) {
	fake.GetHistoryForKeyStub = nil
	if fake.getHistoryForKeyReturnsOnCall == nil {
		fake.getHistoryForKeyReturnsOnCall = make(map[int]struct {
			result1 shim.HistoryQueryIteratorInterface
			result2 error
		})
	}
	fake.getHistoryForKeyReturnsOnCall[i] = struct {
		result1 shim.HistoryQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateData(collection string, key string) ([]byte, error) {
	fake.getPrivateDataMutex.Lock()
	ret, specificReturn := fake.getPrivateDataReturnsOnCall[len(fake.getPrivateDataArgsForCall)]
	fake.getPrivateDataArgsForCall = append(fake.getPrivateDataArgsForCall, struct {
		collection string
		key        string
	}{collection, key})
	fake.recordInvocation("GetPrivateData", []interface{}{collection, key})
	fake.getPrivateDataMutex.Unlock()
	if fake.GetPrivateDataStub != nil {
		return fake.GetPrivateDataStub(collection, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataReturns.result1, fake.getPrivateDataReturns.result2
}

func (fake *ChaincodeStub) GetPrivateDataCallCount() int {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	return len(fake.getPrivateDataArgsForCall)
}

func (fake *ChaincodeStub) GetPrivateDataArgsForCall(i int) (string, string) {
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	return fake.getPrivateDataArgsForCall[i].collection, fake.getPrivateDataArgsForCall[i].key
}

func (fake *ChaincodeStub) GetPrivateDataReturns(result1 []byte, result2 error) {
	fake.GetPrivateDataStub = nil
	fake.getPrivateDataReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetPrivateDataStub = nil
	if fake.getPrivateDataReturnsOnCall == nil {
		fake.getPrivateDataReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getPrivateDataReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataHash(collection string, key string) ([]byte, error) {
	fake.getPrivateDataHashMutex.Lock()
	ret, specificReturn := fake.getPrivateDataHashReturnsOnCall[len(fake.getPrivateDataHashArgsForCall)]
	fake.getPrivateDataHashArgsForCall = append(fake.getPrivateDataHashArgsForCall, struct {
		collection string
		key        string
	}{collection, key})
	fake.recordInvocation("GetPrivateDataHash", []interface{}{collection, key})
	fake.getPrivateDataHashMutex.Unlock()
	if fake.GetPrivateDataHashStub != nil {
		return fake.GetPrivateDataHashStub(collection, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataHashReturns.result1, fake.getPrivateDataHashReturns.result2
}

func (fake *ChaincodeStub) GetPrivateDataHashCallCount() int {
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	return len(fake.getPrivateDataHashArgsForCall)
}

func (fake *ChaincodeStub) GetPrivateDataHashArgsForCall(i int) (string, string) {
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	return fake.getPrivateDataHashArgsForCall[i].collection, fake.getPrivateDataHashArgsForCall[i].key
}

func (fake *ChaincodeStub) GetPrivateDataHashReturns(result1 []byte, result2 error) {
	fake.GetPrivateDataHashStub = nil
	fake.getPrivateDataHashReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataHashReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetPrivateDataHashStub = nil
	if fake.getPrivateDataHashReturnsOnCall == nil {
		fake.getPrivateDataHashReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getPrivateDataHashReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) PutPrivateData(collection string, key string, value []byte) error {
	var valueCopy []byte
	if value != nil {
		valueCopy = make([]byte, len(value))
		copy(valueCopy, value)
	}
	fake.putPrivateDataMutex.Lock()
	ret, specificReturn := fake.putPrivateDataReturnsOnCall[len(fake.putPrivateDataArgsForCall)]
	fake.putPrivateDataArgsForCall = append(fake.putPrivateDataArgsForCall, struct {
		collection string
		key        string
		value      []byte
	}{collection, key, valueCopy})
	fake.recordInvocation("PutPrivateData", []interface{}{collection, key, valueCopy})
	fake.putPrivateDataMutex.Unlock()
	if fake.PutPrivateDataStub != nil {
		return fake.PutPrivateDataStub(collection, key, value)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.putPrivateDataReturns.result1
}

func (fake *ChaincodeStub) PutPrivateDataCallCount() int {
	fake.putPrivateDataMutex.RLock()
	defer fake.putPrivateDataMutex.RUnlock()
	return len(fake.putPrivateDataArgsForCall)
}

func (fake *ChaincodeStub) PutPrivateDataArgsForCall(i int) (string, string, []byte) {
	fake.putPrivateDataMutex.RLock()
	defer fake.putPrivateDataMutex.RUnlock()
	return fake.putPrivateDataArgsForCall[i].collection, fake.putPrivateDataArgsForCall[i].key, fake.putPrivateDataArgsForCall[i].value
}

func (fake *ChaincodeStub) PutPrivateDataReturns(result1 error) {
	fake.PutPrivateDataStub = nil
	fake.putPrivateDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) PutPrivateDataReturnsOnCall(i int, result1 error) {
	fake.PutPrivateDataStub = nil
	if fake.putPrivateDataReturnsOnCall == nil {
		fake.putPrivateDataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.putPrivateDataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) DelPrivateData(collection string, key string) error {
	fake.delPrivateDataMutex.Lock()
	ret, specificReturn := fake.delPrivateDataReturnsOnCall[len(fake.delPrivateDataArgsForCall)]
	fake.delPrivateDataArgsForCall = append(fake.delPrivateDataArgsForCall, struct {
		collection string
		key        string
	}{collection, key})
	fake.recordInvocation("DelPrivateData", []interface{}{collection, key})
	fake.delPrivateDataMutex.Unlock()
	if fake.DelPrivateDataStub != nil {
		return fake.DelPrivateDataStub(collection, key)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.delPrivateDataReturns.result1
}

func (fake *ChaincodeStub) DelPrivateDataCallCount() int {
	fake.delPrivateDataMutex.RLock()
	defer fake.delPrivateDataMutex.RUnlock()
	return len(fake.delPrivateDataArgsForCall)
}

func (fake *ChaincodeStub) DelPrivateDataArgsForCall(i int) (string, string) {
	fake.delPrivateDataMutex.RLock()
	defer fake.delPrivateDataMutex.RUnlock()
	return fake.delPrivateDataArgsForCall[i].collection, fake.delPrivateDataArgsForCall[i].key
}

func (fake *ChaincodeStub) DelPrivateDataReturns(result1 error) {
	fake.DelPrivateDataStub = nil
	fake.delPrivateDataReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) DelPrivateDataReturnsOnCall(i int, result1 error) {
	fake.DelPrivateDataStub = nil
	if fake.delPrivateDataReturnsOnCall == nil {
		fake.delPrivateDataReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.delPrivateDataReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) SetPrivateDataValidationParameter(collection string, key string, ep []byte) error {
	var epCopy []byte
	if ep != nil {
		epCopy = make([]byte, len(ep))
		copy(epCopy, ep)
	}
	fake.setPrivateDataValidationParameterMutex.Lock()
	ret, specificReturn := fake.setPrivateDataValidationParameterReturnsOnCall[len(fake.setPrivateDataValidationParameterArgsForCall)]
	fake.setPrivateDataValidationParameterArgsForCall = append(fake.setPrivateDataValidationParameterArgsForCall, struct {
		collection string
		key        string
		ep         []byte
	}{collection, key, epCopy})
	fake.recordInvocation("SetPrivateDataValidationParameter", []interface{}{collection, key, epCopy})
	fake.setPrivateDataValidationParameterMutex.Unlock()
	if fake.SetPrivateDataValidationParameterStub != nil {
		return fake.SetPrivateDataValidationParameterStub(collection, key, ep)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setPrivateDataValidationParameterReturns.result1
}

func (fake *ChaincodeStub) SetPrivateDataValidationParameterCallCount() int {
	fake.setPrivateDataValidationParameterMutex.RLock()
	defer fake.setPrivateDataValidationParameterMutex.RUnlock()
	return len(fake.setPrivateDataValidationParameterArgsForCall)
}

func (fake *ChaincodeStub) SetPrivateDataValidationParameterArgsForCall(i int) (string, string, []byte) {
	fake.setPrivateDataValidationParameterMutex.RLock()
	defer fake.setPrivateDataValidationParameterMutex.RUnlock()
	return fake.setPrivateDataValidationParameterArgsForCall[i].collection, fake.setPrivateDataValidationParameterArgsForCall[i].key, fake.setPrivateDataValidationParameterArgsForCall[i].ep
}

func (fake *ChaincodeStub) SetPrivateDataValidationParameterReturns(result1 error) {
	fake.SetPrivateDataValidationParameterStub = nil
	fake.setPrivateDataValidationParameterReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) SetPrivateDataValidationParameterReturnsOnCall(i int, result1 error) {
	fake.SetPrivateDataValidationParameterStub = nil
	if fake.setPrivateDataValidationParameterReturnsOnCall == nil {
		fake.setPrivateDataValidationParameterReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setPrivateDataValidationParameterReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) GetPrivateDataValidationParameter(collection string, key string) ([]byte, error) {
	fake.getPrivateDataValidationParameterMutex.Lock()
	ret, specificReturn := fake.getPrivateDataValidationParameterReturnsOnCall[len(fake.getPrivateDataValidationParameterArgsForCall)]
	fake.getPrivateDataValidationParameterArgsForCall = append(fake.getPrivateDataValidationParameterArgsForCall, struct {
		collection string
		key        string
	}{collection, key})
	fake.recordInvocation("GetPrivateDataValidationParameter", []interface{}{collection, key})
	fake.getPrivateDataValidationParameterMutex.Unlock()
	if fake.GetPrivateDataValidationParameterStub != nil {
		return fake.GetPrivateDataValidationParameterStub(collection, key)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataValidationParameterReturns.result1, fake.getPrivateDataValidationParameterReturns.result2
}

func (fake *ChaincodeStub) GetPrivateDataValidationParameterCallCount() int {
	fake.getPrivateDataValidationParameterMutex.RLock()
	defer fake.getPrivateDataValidationParameterMutex.RUnlock()
	return len(fake.getPrivateDataValidationParameterArgsForCall)
}

func (fake *ChaincodeStub) GetPrivateDataValidationParameterArgsForCall(i int) (string, string) {
	fake.getPrivateDataValidationParameterMutex.RLock()
	defer fake.getPrivateDataValidationParameterMutex.RUnlock()
	return fake.getPrivateDataValidationParameterArgsForCall[i].collection, fake.getPrivateDataValidationParameterArgsForCall[i].key
}

func (fake *ChaincodeStub) GetPrivateDataValidationParameterReturns(result1 []byte, result2 error) {
	fake.GetPrivateDataValidationParameterStub = nil
	fake.getPrivateDataValidationParameterReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataValidationParameterReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetPrivateDataValidationParameterStub = nil
	if fake.getPrivateDataValidationParameterReturnsOnCall == nil {
		fake.getPrivateDataValidationParameterReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getPrivateDataValidationParameterReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataByRange(collection string, startKey string, endKey string) (shim.StateQueryIteratorInterface, error) {
	fake.getPrivateDataByRangeMutex.Lock()
	ret, specificReturn := fake.getPrivateDataByRangeReturnsOnCall[len(fake.getPrivateDataByRangeArgsForCall)]
	fake.getPrivateDataByRangeArgsForCall = append(fake.getPrivateDataByRangeArgsForCall, struct {
		collection string
		startKey   string
		endKey     string
	}{collection, startKey, endKey})
	fake.recordInvocation("GetPrivateDataByRange", []interface{}{collection, startKey, endKey})
	fake.getPrivateDataByRangeMutex.Unlock()
	if fake.GetPrivateDataByRangeStub != nil {
		return fake.GetPrivateDataByRangeStub(collection, startKey, endKey)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataByRangeReturns.result1, fake.getPrivateDataByRangeReturns.result2
}

func (fake *ChaincodeStub) GetPrivateDataByRangeCallCount() int {
	fake.getPrivateDataByRangeMutex.RLock()
	defer fake.getPrivateDataByRangeMutex.RUnlock()
	return len(fake.getPrivateDataByRangeArgsForCall)
}

func (fake *ChaincodeStub) GetPrivateDataByRangeArgsForCall(i int) (string, string, string) {
	fake.getPrivateDataByRangeMutex.RLock()
	defer fake.getPrivateDataByRangeMutex.RUnlock()
	return fake.getPrivateDataByRangeArgsForCall[i].collection, fake.getPrivateDataByRangeArgsForCall[i].startKey, fake.getPrivateDataByRangeArgsForCall[i].endKey
}

func (fake *ChaincodeStub) GetPrivateDataByRangeReturns(result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetPrivateDataByRangeStub = nil
	fake.getPrivateDataByRangeReturns = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataByRangeReturnsOnCall(i int, result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetPrivateDataByRangeStub = nil
	if fake.getPrivateDataByRangeReturnsOnCall == nil {
		fake.getPrivateDataByRangeReturnsOnCall = make(map[int]struct {
			result1 shim.StateQueryIteratorInterface
			result2 error
		})
	}
	fake.getPrivateDataByRangeReturnsOnCall[i] = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataByPartialCompositeKey(collection string, objectType string, keys []string) (shim.StateQueryIteratorInterface, error) {
	var keysCopy []string
	if keys != nil {
		keysCopy = make([]string, len(keys))
		copy(keysCopy, keys)
	}
	fake.getPrivateDataByPartialCompositeKeyMutex.Lock()
	ret, specificReturn := fake.getPrivateDataByPartialCompositeKeyReturnsOnCall[len(fake.getPrivateDataByPartialCompositeKeyArgsForCall)]
	fake.getPrivateDataByPartialCompositeKeyArgsForCall = append(fake.getPrivateDataByPartialCompositeKeyArgsForCall, struct {
		collection string
		objectType string
		keys       []string
	}{collection, objectType, keysCopy})
	fake.recordInvocation("GetPrivateDataByPartialCompositeKey", []interface{}{collection, objectType, keysCopy})
	fake.getPrivateDataByPartialCompositeKeyMutex.Unlock()
	if fake.GetPrivateDataByPartialCompositeKeyStub != nil {
		return fake.GetPrivateDataByPartialCompositeKeyStub(collection, objectType, keys)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataByPartialCompositeKeyReturns.result1, fake.getPrivateDataByPartialCompositeKeyReturns.result2
}

func (fake *ChaincodeStub) GetPrivateDataByPartialCompositeKeyCallCount() int {
	fake.getPrivateDataByPartialCompositeKeyMutex.RLock()
	defer fake.getPrivateDataByPartialCompositeKeyMutex.RUnlock()
	return len(fake.getPrivateDataByPartialCompositeKeyArgsForCall)
}

func (fake *ChaincodeStub) GetPrivateDataByPartialCompositeKeyArgsForCall(i int) (string, string, []string) {
	fake.getPrivateDataByPartialCompositeKeyMutex.RLock()
	defer fake.getPrivateDataByPartialCompositeKeyMutex.RUnlock()
	return fake.getPrivateDataByPartialCompositeKeyArgsForCall[i].collection, fake.getPrivateDataByPartialCompositeKeyArgsForCall[i].objectType, fake.getPrivateDataByPartialCompositeKeyArgsForCall[i].keys
}

func (fake *ChaincodeStub) GetPrivateDataByPartialCompositeKeyReturns(result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetPrivateDataByPartialCompositeKeyStub = nil
	fake.getPrivateDataByPartialCompositeKeyReturns = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataByPartialCompositeKeyReturnsOnCall(i int, result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetPrivateDataByPartialCompositeKeyStub = nil
	if fake.getPrivateDataByPartialCompositeKeyReturnsOnCall == nil {
		fake.getPrivateDataByPartialCompositeKeyReturnsOnCall = make(map[int]struct {
			result1 shim.StateQueryIteratorInterface
			result2 error
		})
	}
	fake.getPrivateDataByPartialCompositeKeyReturnsOnCall[i] = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataQueryResult(collection string, query string) (shim.StateQueryIteratorInterface, error) {
	fake.getPrivateDataQueryResultMutex.Lock()
	ret, specificReturn := fake.getPrivateDataQueryResultReturnsOnCall[len(fake.getPrivateDataQueryResultArgsForCall)]
	fake.getPrivateDataQueryResultArgsForCall = append(fake.getPrivateDataQueryResultArgsForCall, struct {
		collection string
		query      string
	}{collection, query})
	fake.recordInvocation("GetPrivateDataQueryResult", []interface{}{collection, query})
	fake.getPrivateDataQueryResultMutex.Unlock()
	if fake.GetPrivateDataQueryResultStub != nil {
		return fake.GetPrivateDataQueryResultStub(collection, query)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getPrivateDataQueryResultReturns.result1, fake.getPrivateDataQueryResultReturns.result2
}

func (fake *ChaincodeStub) GetPrivateDataQueryResultCallCount() int {
	fake.getPrivateDataQueryResultMutex.RLock()
	defer fake.getPrivateDataQueryResultMutex.RUnlock()
	return len(fake.getPrivateDataQueryResultArgsForCall)
}

func (fake *ChaincodeStub) GetPrivateDataQueryResultArgsForCall(i int) (string, string) {
	fake.getPrivateDataQueryResultMutex.RLock()
	defer fake.getPrivateDataQueryResultMutex.RUnlock()
	return fake.getPrivateDataQueryResultArgsForCall[i].collection, fake.getPrivateDataQueryResultArgsForCall[i].query
}

func (fake *ChaincodeStub) GetPrivateDataQueryResultReturns(result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetPrivateDataQueryResultStub = nil
	fake.getPrivateDataQueryResultReturns = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetPrivateDataQueryResultReturnsOnCall(i int, result1 shim.StateQueryIteratorInterface, result2 error) {
	fake.GetPrivateDataQueryResultStub = nil
	if fake.getPrivateDataQueryResultReturnsOnCall == nil {
		fake.getPrivateDataQueryResultReturnsOnCall = make(map[int]struct {
			result1 shim.StateQueryIteratorInterface
			result2 error
		})
	}
	fake.getPrivateDataQueryResultReturnsOnCall[i] = struct {
		result1 shim.StateQueryIteratorInterface
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetCreator() ([]byte, error) {
	fake.getCreatorMutex.Lock()
	ret, specificReturn := fake.getCreatorReturnsOnCall[len(fake.getCreatorArgsForCall)]
	fake.getCreatorArgsForCall = append(fake.getCreatorArgsForCall, struct{}{})
	fake.recordInvocation("GetCreator", []interface{}{})
	fake.getCreatorMutex.Unlock()
	if fake.GetCreatorStub != nil {
		return fake.GetCreatorStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getCreatorReturns.result1, fake.getCreatorReturns.result2
}

func (fake *ChaincodeStub) GetCreatorCallCount() int {
	fake.getCreatorMutex.RLock()
	defer fake.getCreatorMutex.RUnlock()
	return len(fake.getCreatorArgsForCall)
}

func (fake *ChaincodeStub) GetCreatorReturns(result1 []byte, result2 error) {
	fake.GetCreatorStub = nil
	fake.getCreatorReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetCreatorReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetCreatorStub = nil
	if fake.getCreatorReturnsOnCall == nil {
		fake.getCreatorReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getCreatorReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetTransient() (map[string][]byte, error) {
	fake.getTransientMutex.Lock()
	ret, specificReturn := fake.getTransientReturnsOnCall[len(fake.getTransientArgsForCall)]
	fake.getTransientArgsForCall = append(fake.getTransientArgsForCall, struct{}{})
	fake.recordInvocation("GetTransient", []interface{}{})
	fake.getTransientMutex.Unlock()
	if fake.GetTransientStub != nil {
		return fake.GetTransientStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTransientReturns.result1, fake.getTransientReturns.result2
}

func (fake *ChaincodeStub) GetTransientCallCount() int {
	fake.getTransientMutex.RLock()
	defer fake.getTransientMutex.RUnlock()
	return len(fake.getTransientArgsForCall)
}

func (fake *ChaincodeStub) GetTransientReturns(result1 map[string][]byte, result2 error) {
	fake.GetTransientStub = nil
	fake.getTransientReturns = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetTransientReturnsOnCall(i int, result1 map[string][]byte, result2 error) {
	fake.GetTransientStub = nil
	if fake.getTransientReturnsOnCall == nil {
		fake.getTransientReturnsOnCall = make(map[int]struct {
			result1 map[string][]byte
			result2 error
		})
	}
	fake.getTransientReturnsOnCall[i] = struct {
		result1 map[string][]byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetBinding() ([]byte, error) {
	fake.getBindingMutex.Lock()
	ret, specificReturn := fake.getBindingReturnsOnCall[len(fake.getBindingArgsForCall)]
	fake.getBindingArgsForCall = append(fake.getBindingArgsForCall, struct{}{})
	fake.recordInvocation("GetBinding", []interface{}{})
	fake.getBindingMutex.Unlock()
	if fake.GetBindingStub != nil {
		return fake.GetBindingStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getBindingReturns.result1, fake.getBindingReturns.result2
}

func (fake *ChaincodeStub) GetBindingCallCount() int {
	fake.getBindingMutex.RLock()
	defer fake.getBindingMutex.RUnlock()
	return len(fake.getBindingArgsForCall)
}

func (fake *ChaincodeStub) GetBindingReturns(result1 []byte, result2 error) {
	fake.GetBindingStub = nil
	fake.getBindingReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetBindingReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetBindingStub = nil
	if fake.getBindingReturnsOnCall == nil {
		fake.getBindingReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getBindingReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetDecorations() map[string][]byte {
	fake.getDecorationsMutex.Lock()
	ret, specificReturn := fake.getDecorationsReturnsOnCall[len(fake.getDecorationsArgsForCall)]
	fake.getDecorationsArgsForCall = append(fake.getDecorationsArgsForCall, struct{}{})
	fake.recordInvocation("GetDecorations", []interface{}{})
	fake.getDecorationsMutex.Unlock()
	if fake.GetDecorationsStub != nil {
		return fake.GetDecorationsStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getDecorationsReturns.result1
}

func (fake *ChaincodeStub) GetDecorationsCallCount() int {
	fake.getDecorationsMutex.RLock()
	defer fake.getDecorationsMutex.RUnlock()
	return len(fake.getDecorationsArgsForCall)
}

func (fake *ChaincodeStub) GetDecorationsReturns(result1 map[string][]byte) {
	fake.GetDecorationsStub = nil
	fake.getDecorationsReturns = struct {
		result1 map[string][]byte
	}{result1}
}

func (fake *ChaincodeStub) GetDecorationsReturnsOnCall(i int, result1 map[string][]byte) {
	fake.GetDecorationsStub = nil
	if fake.getDecorationsReturnsOnCall == nil {
		fake.getDecorationsReturnsOnCall = make(map[int]struct {
			result1 map[string][]byte
		})
	}
	fake.getDecorationsReturnsOnCall[i] = struct {
		result1 map[string][]byte
	}{result1}
}

func (fake *ChaincodeStub) GetSignedProposal() (*pb.SignedProposal, error) {
	fake.getSignedProposalMutex.Lock()
	ret, specificReturn := fake.getSignedProposalReturnsOnCall[len(fake.getSignedProposalArgsForCall)]
	fake.getSignedProposalArgsForCall = append(fake.getSignedProposalArgsForCall, struct{}{})
	fake.recordInvocation("GetSignedProposal", []interface{}{})
	fake.getSignedProposalMutex.Unlock()
	if fake.GetSignedProposalStub != nil {
		return fake.GetSignedProposalStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getSignedProposalReturns.result1, fake.getSignedProposalReturns.result2
}

func (fake *ChaincodeStub) GetSignedProposalCallCount() int {
	fake.getSignedProposalMutex.RLock()
	defer fake.getSignedProposalMutex.RUnlock()
	return len(fake.getSignedProposalArgsForCall)
}

func (fake *ChaincodeStub) GetSignedProposalReturns(result1 *pb.SignedProposal, result2 error) {
	fake.GetSignedProposalStub = nil
	fake.getSignedProposalReturns = struct {
		result1 *pb.SignedProposal
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetSignedProposalReturnsOnCall(i int, result1 *pb.SignedProposal, result2 error) {
	fake.GetSignedProposalStub = nil
	if fake.getSignedProposalReturnsOnCall == nil {
		fake.getSignedProposalReturnsOnCall = make(map[int]struct {
			result1 *pb.SignedProposal
			result2 error
		})
	}
	fake.getSignedProposalReturnsOnCall[i] = struct {
		result1 *pb.SignedProposal
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetTxTimestamp() (*timestamp.Timestamp, error) {
	fake.getTxTimestampMutex.Lock()
	ret, specificReturn := fake.getTxTimestampReturnsOnCall[len(fake.getTxTimestampArgsForCall)]
	fake.getTxTimestampArgsForCall = append(fake.getTxTimestampArgsForCall, struct{}{})
	fake.recordInvocation("GetTxTimestamp", []interface{}{})
	fake.getTxTimestampMutex.Unlock()
	if fake.GetTxTimestampStub != nil {
		return fake.GetTxTimestampStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getTxTimestampReturns.result1, fake.getTxTimestampReturns.result2
}

func (fake *ChaincodeStub) GetTxTimestampCallCount() int {
	fake.getTxTimestampMutex.RLock()
	defer fake.getTxTimestampMutex.RUnlock()
	return len(fake.getTxTimestampArgsForCall)
}

func (fake *ChaincodeStub) GetTxTimestampReturns(result1 *timestamp.Timestamp, result2 error) {
	fake.GetTxTimestampStub = nil
	fake.getTxTimestampReturns = struct {
		result1 *timestamp.Timestamp
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) GetTxTimestampReturnsOnCall(i int, result1 *timestamp.Timestamp, result2 error) {
	fake.GetTxTimestampStub = nil
	if fake.getTxTimestampReturnsOnCall == nil {
		fake.getTxTimestampReturnsOnCall = make(map[int]struct {
			result1 *timestamp.Timestamp
			result2 error
		})
	}
	fake.getTxTimestampReturnsOnCall[i] = struct {
		result1 *timestamp.Timestamp
		result2 error
	}{result1, result2}
}

func (fake *ChaincodeStub) SetEvent(name string, payload []byte) error {
	var payloadCopy []byte
	if payload != nil {
		payloadCopy = make([]byte, len(payload))
		copy(payloadCopy, payload)
	}
	fake.setEventMutex.Lock()
	ret, specificReturn := fake.setEventReturnsOnCall[len(fake.setEventArgsForCall)]
	fake.setEventArgsForCall = append(fake.setEventArgsForCall, struct {
		name    string
		payload []byte
	}{name, payloadCopy})
	fake.recordInvocation("SetEvent", []interface{}{name, payloadCopy})
	fake.setEventMutex.Unlock()
	if fake.SetEventStub != nil {
		return fake.SetEventStub(name, payload)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.setEventReturns.result1
}

func (fake *ChaincodeStub) SetEventCallCount() int {
	fake.setEventMutex.RLock()
	defer fake.setEventMutex.RUnlock()
	return len(fake.setEventArgsForCall)
}

func (fake *ChaincodeStub) SetEventArgsForCall(i int) (string, []byte) {
	fake.setEventMutex.RLock()
	defer fake.setEventMutex.RUnlock()
	return fake.setEventArgsForCall[i].name, fake.setEventArgsForCall[i].payload
}

func (fake *ChaincodeStub) SetEventReturns(result1 error) {
	fake.SetEventStub = nil
	fake.setEventReturns = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) SetEventReturnsOnCall(i int, result1 error) {
	fake.SetEventStub = nil
	if fake.setEventReturnsOnCall == nil {
		fake.setEventReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.setEventReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *ChaincodeStub) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.getArgsMutex.RLock()
	defer fake.getArgsMutex.RUnlock()
	fake.getStringArgsMutex.RLock()
	defer fake.getStringArgsMutex.RUnlock()
	fake.getFunctionAndParametersMutex.RLock()
	defer fake.getFunctionAndParametersMutex.RUnlock()
	fake.getArgsSliceMutex.RLock()
	defer fake.getArgsSliceMutex.RUnlock()
	fake.getTxIDMutex.RLock()
	defer fake.getTxIDMutex.RUnlock()
	fake.getChannelIDMutex.RLock()
	defer fake.getChannelIDMutex.RUnlock()
	fake.invokeChaincodeMutex.RLock()
	defer fake.invokeChaincodeMutex.RUnlock()
	fake.getStateMutex.RLock()
	defer fake.getStateMutex.RUnlock()
	fake.putStateMutex.RLock()
	defer fake.putStateMutex.RUnlock()
	fake.delStateMutex.RLock()
	defer fake.delStateMutex.RUnlock()
	fake.setStateValidationParameterMutex.RLock()
	defer fake.setStateValidationParameterMutex.RUnlock()
	fake.getStateValidationParameterMutex.RLock()
	defer fake.getStateValidationParameterMutex.RUnlock()
	fake.getStateByRangeMutex.RLock()
	defer fake.getStateByRangeMutex.RUnlock()
	fake.getStateByRangeWithPaginationMutex.RLock()
	defer fake.getStateByRangeWithPaginationMutex.RUnlock()
	fake.getStateByPartialCompositeKeyMutex.RLock()
	defer fake.getStateByPartialCompositeKeyMutex.RUnlock()
	fake.getStateByPartialCompositeKeyWithPaginationMutex.RLock()
	defer fake.getStateByPartialCompositeKeyWithPaginationMutex.RUnlock()
	fake.createCompositeKeyMutex.RLock()
	defer fake.createCompositeKeyMutex.RUnlock()
	fake.splitCompositeKeyMutex.RLock()
	defer fake.splitCompositeKeyMutex.RUnlock()
	fake.getQueryResultMutex.RLock()
	defer fake.getQueryResultMutex.RUnlock()
	fake.getQueryResultWithPaginationMutex.RLock()
	defer fake.getQueryResultWithPaginationMutex.RUnlock()
	fake.getHistoryForKeyMutex.RLock()
	defer fake.getHistoryForKeyMutex.RUnlock()
	fake.getPrivateDataMutex.RLock()
	defer fake.getPrivateDataMutex.RUnlock()
	fake.getPrivateDataHashMutex.RLock()
	defer fake.getPrivateDataHashMutex.RUnlock()
	fake.putPrivateDataMutex.RLock()
	defer fake.putPrivateDataMutex.RUnlock()
	fake.delPrivateDataMutex.RLock()
	defer fake.delPrivateDataMutex.RUnlock()
	fake.setPrivateDataValidationParameterMutex.RLock()
	defer fake.setPrivateDataValidationParameterMutex.RUnlock()
	fake.getPrivateDataValidationParameterMutex.RLock()
	defer fake.getPrivateDataValidationParameterMutex.RUnlock()
	fake.getPrivateDataByRangeMutex.RLock()
	defer fake.getPrivateDataByRangeMutex.RUnlock()
	fake.getPrivateDataByPartialCompositeKeyMutex.RLock()
	defer fake.getPrivateDataByPartialCompositeKeyMutex.RUnlock()
	fake.getPrivateDataQueryResultMutex.RLock()
	defer fake.getPrivateDataQueryResultMutex.RUnlock()
	fake.getCreatorMutex.RLock()
	defer fake.getCreatorMutex.RUnlock()
	fake.getTransientMutex.RLock()
	defer fake.getTransientMutex.RUnlock()
	fake.getBindingMutex.RLock()
	defer fake.getBindingMutex.RUnlock()
	fake.getDecorationsMutex.RLock()
	defer fake.getDecorationsMutex.RUnlock()
	fake.getSignedProposalMutex.RLock()
	defer fake.getSignedProposalMutex.RUnlock()
	fake.getTxTimestampMutex.RLock()
	defer fake.getTxTimestampMutex.RUnlock()
	fake.setEventMutex.RLock()
	defer fake.setEventMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ChaincodeStub) recordInvocation(key string, args []interface{}) {
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
