
package mock

import (
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

type CCPackage struct {
	InitFromBufferStub        func(buf []byte) (*ccprovider.ChaincodeData, error)
	initFromBufferMutex       sync.RWMutex
	initFromBufferArgsForCall []struct {
		buf []byte
	}
	initFromBufferReturns struct {
		result1 *ccprovider.ChaincodeData
		result2 error
	}
	initFromBufferReturnsOnCall map[int]struct {
		result1 *ccprovider.ChaincodeData
		result2 error
	}
	InitFromFSStub        func(ccname string, ccversion string) ([]byte, *pb.ChaincodeDeploymentSpec, error)
	initFromFSMutex       sync.RWMutex
	initFromFSArgsForCall []struct {
		ccname    string
		ccversion string
	}
	initFromFSReturns struct {
		result1 []byte
		result2 *pb.ChaincodeDeploymentSpec
		result3 error
	}
	initFromFSReturnsOnCall map[int]struct {
		result1 []byte
		result2 *pb.ChaincodeDeploymentSpec
		result3 error
	}
	PutChaincodeToFSStub        func() error
	putChaincodeToFSMutex       sync.RWMutex
	putChaincodeToFSArgsForCall []struct{}
	putChaincodeToFSReturns     struct {
		result1 error
	}
	putChaincodeToFSReturnsOnCall map[int]struct {
		result1 error
	}
	GetDepSpecStub        func() *pb.ChaincodeDeploymentSpec
	getDepSpecMutex       sync.RWMutex
	getDepSpecArgsForCall []struct{}
	getDepSpecReturns     struct {
		result1 *pb.ChaincodeDeploymentSpec
	}
	getDepSpecReturnsOnCall map[int]struct {
		result1 *pb.ChaincodeDeploymentSpec
	}
	GetDepSpecBytesStub        func() []byte
	getDepSpecBytesMutex       sync.RWMutex
	getDepSpecBytesArgsForCall []struct{}
	getDepSpecBytesReturns     struct {
		result1 []byte
	}
	getDepSpecBytesReturnsOnCall map[int]struct {
		result1 []byte
	}
	ValidateCCStub        func(ccdata *ccprovider.ChaincodeData) error
	validateCCMutex       sync.RWMutex
	validateCCArgsForCall []struct {
		ccdata *ccprovider.ChaincodeData
	}
	validateCCReturns struct {
		result1 error
	}
	validateCCReturnsOnCall map[int]struct {
		result1 error
	}
	GetPackageObjectStub        func() proto.Message
	getPackageObjectMutex       sync.RWMutex
	getPackageObjectArgsForCall []struct{}
	getPackageObjectReturns     struct {
		result1 proto.Message
	}
	getPackageObjectReturnsOnCall map[int]struct {
		result1 proto.Message
	}
	GetChaincodeDataStub        func() *ccprovider.ChaincodeData
	getChaincodeDataMutex       sync.RWMutex
	getChaincodeDataArgsForCall []struct{}
	getChaincodeDataReturns     struct {
		result1 *ccprovider.ChaincodeData
	}
	getChaincodeDataReturnsOnCall map[int]struct {
		result1 *ccprovider.ChaincodeData
	}
	GetIdStub        func() []byte
	getIdMutex       sync.RWMutex
	getIdArgsForCall []struct{}
	getIdReturns     struct {
		result1 []byte
	}
	getIdReturnsOnCall map[int]struct {
		result1 []byte
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *CCPackage) InitFromBuffer(buf []byte) (*ccprovider.ChaincodeData, error) {
	var bufCopy []byte
	if buf != nil {
		bufCopy = make([]byte, len(buf))
		copy(bufCopy, buf)
	}
	fake.initFromBufferMutex.Lock()
	ret, specificReturn := fake.initFromBufferReturnsOnCall[len(fake.initFromBufferArgsForCall)]
	fake.initFromBufferArgsForCall = append(fake.initFromBufferArgsForCall, struct {
		buf []byte
	}{bufCopy})
	fake.recordInvocation("InitFromBuffer", []interface{}{bufCopy})
	fake.initFromBufferMutex.Unlock()
	if fake.InitFromBufferStub != nil {
		return fake.InitFromBufferStub(buf)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.initFromBufferReturns.result1, fake.initFromBufferReturns.result2
}

func (fake *CCPackage) InitFromBufferCallCount() int {
	fake.initFromBufferMutex.RLock()
	defer fake.initFromBufferMutex.RUnlock()
	return len(fake.initFromBufferArgsForCall)
}

func (fake *CCPackage) InitFromBufferArgsForCall(i int) []byte {
	fake.initFromBufferMutex.RLock()
	defer fake.initFromBufferMutex.RUnlock()
	return fake.initFromBufferArgsForCall[i].buf
}

func (fake *CCPackage) InitFromBufferReturns(result1 *ccprovider.ChaincodeData, result2 error) {
	fake.InitFromBufferStub = nil
	fake.initFromBufferReturns = struct {
		result1 *ccprovider.ChaincodeData
		result2 error
	}{result1, result2}
}

func (fake *CCPackage) InitFromBufferReturnsOnCall(i int, result1 *ccprovider.ChaincodeData, result2 error) {
	fake.InitFromBufferStub = nil
	if fake.initFromBufferReturnsOnCall == nil {
		fake.initFromBufferReturnsOnCall = make(map[int]struct {
			result1 *ccprovider.ChaincodeData
			result2 error
		})
	}
	fake.initFromBufferReturnsOnCall[i] = struct {
		result1 *ccprovider.ChaincodeData
		result2 error
	}{result1, result2}
}

func (fake *CCPackage) InitFromFS(ccname string, ccversion string) ([]byte, *pb.ChaincodeDeploymentSpec, error) {
	fake.initFromFSMutex.Lock()
	ret, specificReturn := fake.initFromFSReturnsOnCall[len(fake.initFromFSArgsForCall)]
	fake.initFromFSArgsForCall = append(fake.initFromFSArgsForCall, struct {
		ccname    string
		ccversion string
	}{ccname, ccversion})
	fake.recordInvocation("InitFromFS", []interface{}{ccname, ccversion})
	fake.initFromFSMutex.Unlock()
	if fake.InitFromFSStub != nil {
		return fake.InitFromFSStub(ccname, ccversion)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	return fake.initFromFSReturns.result1, fake.initFromFSReturns.result2, fake.initFromFSReturns.result3
}

func (fake *CCPackage) InitFromFSCallCount() int {
	fake.initFromFSMutex.RLock()
	defer fake.initFromFSMutex.RUnlock()
	return len(fake.initFromFSArgsForCall)
}

func (fake *CCPackage) InitFromFSArgsForCall(i int) (string, string) {
	fake.initFromFSMutex.RLock()
	defer fake.initFromFSMutex.RUnlock()
	return fake.initFromFSArgsForCall[i].ccname, fake.initFromFSArgsForCall[i].ccversion
}

func (fake *CCPackage) InitFromFSReturns(result1 []byte, result2 *pb.ChaincodeDeploymentSpec, result3 error) {
	fake.InitFromFSStub = nil
	fake.initFromFSReturns = struct {
		result1 []byte
		result2 *pb.ChaincodeDeploymentSpec
		result3 error
	}{result1, result2, result3}
}

func (fake *CCPackage) InitFromFSReturnsOnCall(i int, result1 []byte, result2 *pb.ChaincodeDeploymentSpec, result3 error) {
	fake.InitFromFSStub = nil
	if fake.initFromFSReturnsOnCall == nil {
		fake.initFromFSReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 *pb.ChaincodeDeploymentSpec
			result3 error
		})
	}
	fake.initFromFSReturnsOnCall[i] = struct {
		result1 []byte
		result2 *pb.ChaincodeDeploymentSpec
		result3 error
	}{result1, result2, result3}
}

func (fake *CCPackage) PutChaincodeToFS() error {
	fake.putChaincodeToFSMutex.Lock()
	ret, specificReturn := fake.putChaincodeToFSReturnsOnCall[len(fake.putChaincodeToFSArgsForCall)]
	fake.putChaincodeToFSArgsForCall = append(fake.putChaincodeToFSArgsForCall, struct{}{})
	fake.recordInvocation("PutChaincodeToFS", []interface{}{})
	fake.putChaincodeToFSMutex.Unlock()
	if fake.PutChaincodeToFSStub != nil {
		return fake.PutChaincodeToFSStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.putChaincodeToFSReturns.result1
}

func (fake *CCPackage) PutChaincodeToFSCallCount() int {
	fake.putChaincodeToFSMutex.RLock()
	defer fake.putChaincodeToFSMutex.RUnlock()
	return len(fake.putChaincodeToFSArgsForCall)
}

func (fake *CCPackage) PutChaincodeToFSReturns(result1 error) {
	fake.PutChaincodeToFSStub = nil
	fake.putChaincodeToFSReturns = struct {
		result1 error
	}{result1}
}

func (fake *CCPackage) PutChaincodeToFSReturnsOnCall(i int, result1 error) {
	fake.PutChaincodeToFSStub = nil
	if fake.putChaincodeToFSReturnsOnCall == nil {
		fake.putChaincodeToFSReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.putChaincodeToFSReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CCPackage) GetDepSpec() *pb.ChaincodeDeploymentSpec {
	fake.getDepSpecMutex.Lock()
	ret, specificReturn := fake.getDepSpecReturnsOnCall[len(fake.getDepSpecArgsForCall)]
	fake.getDepSpecArgsForCall = append(fake.getDepSpecArgsForCall, struct{}{})
	fake.recordInvocation("GetDepSpec", []interface{}{})
	fake.getDepSpecMutex.Unlock()
	if fake.GetDepSpecStub != nil {
		return fake.GetDepSpecStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getDepSpecReturns.result1
}

func (fake *CCPackage) GetDepSpecCallCount() int {
	fake.getDepSpecMutex.RLock()
	defer fake.getDepSpecMutex.RUnlock()
	return len(fake.getDepSpecArgsForCall)
}

func (fake *CCPackage) GetDepSpecReturns(result1 *pb.ChaincodeDeploymentSpec) {
	fake.GetDepSpecStub = nil
	fake.getDepSpecReturns = struct {
		result1 *pb.ChaincodeDeploymentSpec
	}{result1}
}

func (fake *CCPackage) GetDepSpecReturnsOnCall(i int, result1 *pb.ChaincodeDeploymentSpec) {
	fake.GetDepSpecStub = nil
	if fake.getDepSpecReturnsOnCall == nil {
		fake.getDepSpecReturnsOnCall = make(map[int]struct {
			result1 *pb.ChaincodeDeploymentSpec
		})
	}
	fake.getDepSpecReturnsOnCall[i] = struct {
		result1 *pb.ChaincodeDeploymentSpec
	}{result1}
}

func (fake *CCPackage) GetDepSpecBytes() []byte {
	fake.getDepSpecBytesMutex.Lock()
	ret, specificReturn := fake.getDepSpecBytesReturnsOnCall[len(fake.getDepSpecBytesArgsForCall)]
	fake.getDepSpecBytesArgsForCall = append(fake.getDepSpecBytesArgsForCall, struct{}{})
	fake.recordInvocation("GetDepSpecBytes", []interface{}{})
	fake.getDepSpecBytesMutex.Unlock()
	if fake.GetDepSpecBytesStub != nil {
		return fake.GetDepSpecBytesStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getDepSpecBytesReturns.result1
}

func (fake *CCPackage) GetDepSpecBytesCallCount() int {
	fake.getDepSpecBytesMutex.RLock()
	defer fake.getDepSpecBytesMutex.RUnlock()
	return len(fake.getDepSpecBytesArgsForCall)
}

func (fake *CCPackage) GetDepSpecBytesReturns(result1 []byte) {
	fake.GetDepSpecBytesStub = nil
	fake.getDepSpecBytesReturns = struct {
		result1 []byte
	}{result1}
}

func (fake *CCPackage) GetDepSpecBytesReturnsOnCall(i int, result1 []byte) {
	fake.GetDepSpecBytesStub = nil
	if fake.getDepSpecBytesReturnsOnCall == nil {
		fake.getDepSpecBytesReturnsOnCall = make(map[int]struct {
			result1 []byte
		})
	}
	fake.getDepSpecBytesReturnsOnCall[i] = struct {
		result1 []byte
	}{result1}
}

func (fake *CCPackage) ValidateCC(ccdata *ccprovider.ChaincodeData) error {
	fake.validateCCMutex.Lock()
	ret, specificReturn := fake.validateCCReturnsOnCall[len(fake.validateCCArgsForCall)]
	fake.validateCCArgsForCall = append(fake.validateCCArgsForCall, struct {
		ccdata *ccprovider.ChaincodeData
	}{ccdata})
	fake.recordInvocation("ValidateCC", []interface{}{ccdata})
	fake.validateCCMutex.Unlock()
	if fake.ValidateCCStub != nil {
		return fake.ValidateCCStub(ccdata)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.validateCCReturns.result1
}

func (fake *CCPackage) ValidateCCCallCount() int {
	fake.validateCCMutex.RLock()
	defer fake.validateCCMutex.RUnlock()
	return len(fake.validateCCArgsForCall)
}

func (fake *CCPackage) ValidateCCArgsForCall(i int) *ccprovider.ChaincodeData {
	fake.validateCCMutex.RLock()
	defer fake.validateCCMutex.RUnlock()
	return fake.validateCCArgsForCall[i].ccdata
}

func (fake *CCPackage) ValidateCCReturns(result1 error) {
	fake.ValidateCCStub = nil
	fake.validateCCReturns = struct {
		result1 error
	}{result1}
}

func (fake *CCPackage) ValidateCCReturnsOnCall(i int, result1 error) {
	fake.ValidateCCStub = nil
	if fake.validateCCReturnsOnCall == nil {
		fake.validateCCReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.validateCCReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *CCPackage) GetPackageObject() proto.Message {
	fake.getPackageObjectMutex.Lock()
	ret, specificReturn := fake.getPackageObjectReturnsOnCall[len(fake.getPackageObjectArgsForCall)]
	fake.getPackageObjectArgsForCall = append(fake.getPackageObjectArgsForCall, struct{}{})
	fake.recordInvocation("GetPackageObject", []interface{}{})
	fake.getPackageObjectMutex.Unlock()
	if fake.GetPackageObjectStub != nil {
		return fake.GetPackageObjectStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getPackageObjectReturns.result1
}

func (fake *CCPackage) GetPackageObjectCallCount() int {
	fake.getPackageObjectMutex.RLock()
	defer fake.getPackageObjectMutex.RUnlock()
	return len(fake.getPackageObjectArgsForCall)
}

func (fake *CCPackage) GetPackageObjectReturns(result1 proto.Message) {
	fake.GetPackageObjectStub = nil
	fake.getPackageObjectReturns = struct {
		result1 proto.Message
	}{result1}
}

func (fake *CCPackage) GetPackageObjectReturnsOnCall(i int, result1 proto.Message) {
	fake.GetPackageObjectStub = nil
	if fake.getPackageObjectReturnsOnCall == nil {
		fake.getPackageObjectReturnsOnCall = make(map[int]struct {
			result1 proto.Message
		})
	}
	fake.getPackageObjectReturnsOnCall[i] = struct {
		result1 proto.Message
	}{result1}
}

func (fake *CCPackage) GetChaincodeData() *ccprovider.ChaincodeData {
	fake.getChaincodeDataMutex.Lock()
	ret, specificReturn := fake.getChaincodeDataReturnsOnCall[len(fake.getChaincodeDataArgsForCall)]
	fake.getChaincodeDataArgsForCall = append(fake.getChaincodeDataArgsForCall, struct{}{})
	fake.recordInvocation("GetChaincodeData", []interface{}{})
	fake.getChaincodeDataMutex.Unlock()
	if fake.GetChaincodeDataStub != nil {
		return fake.GetChaincodeDataStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getChaincodeDataReturns.result1
}

func (fake *CCPackage) GetChaincodeDataCallCount() int {
	fake.getChaincodeDataMutex.RLock()
	defer fake.getChaincodeDataMutex.RUnlock()
	return len(fake.getChaincodeDataArgsForCall)
}

func (fake *CCPackage) GetChaincodeDataReturns(result1 *ccprovider.ChaincodeData) {
	fake.GetChaincodeDataStub = nil
	fake.getChaincodeDataReturns = struct {
		result1 *ccprovider.ChaincodeData
	}{result1}
}

func (fake *CCPackage) GetChaincodeDataReturnsOnCall(i int, result1 *ccprovider.ChaincodeData) {
	fake.GetChaincodeDataStub = nil
	if fake.getChaincodeDataReturnsOnCall == nil {
		fake.getChaincodeDataReturnsOnCall = make(map[int]struct {
			result1 *ccprovider.ChaincodeData
		})
	}
	fake.getChaincodeDataReturnsOnCall[i] = struct {
		result1 *ccprovider.ChaincodeData
	}{result1}
}

func (fake *CCPackage) GetId() []byte {
	fake.getIdMutex.Lock()
	ret, specificReturn := fake.getIdReturnsOnCall[len(fake.getIdArgsForCall)]
	fake.getIdArgsForCall = append(fake.getIdArgsForCall, struct{}{})
	fake.recordInvocation("GetId", []interface{}{})
	fake.getIdMutex.Unlock()
	if fake.GetIdStub != nil {
		return fake.GetIdStub()
	}
	if specificReturn {
		return ret.result1
	}
	return fake.getIdReturns.result1
}

func (fake *CCPackage) GetIdCallCount() int {
	fake.getIdMutex.RLock()
	defer fake.getIdMutex.RUnlock()
	return len(fake.getIdArgsForCall)
}

func (fake *CCPackage) GetIdReturns(result1 []byte) {
	fake.GetIdStub = nil
	fake.getIdReturns = struct {
		result1 []byte
	}{result1}
}

func (fake *CCPackage) GetIdReturnsOnCall(i int, result1 []byte) {
	fake.GetIdStub = nil
	if fake.getIdReturnsOnCall == nil {
		fake.getIdReturnsOnCall = make(map[int]struct {
			result1 []byte
		})
	}
	fake.getIdReturnsOnCall[i] = struct {
		result1 []byte
	}{result1}
}

func (fake *CCPackage) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.initFromBufferMutex.RLock()
	defer fake.initFromBufferMutex.RUnlock()
	fake.initFromFSMutex.RLock()
	defer fake.initFromFSMutex.RUnlock()
	fake.putChaincodeToFSMutex.RLock()
	defer fake.putChaincodeToFSMutex.RUnlock()
	fake.getDepSpecMutex.RLock()
	defer fake.getDepSpecMutex.RUnlock()
	fake.getDepSpecBytesMutex.RLock()
	defer fake.getDepSpecBytesMutex.RUnlock()
	fake.validateCCMutex.RLock()
	defer fake.validateCCMutex.RUnlock()
	fake.getPackageObjectMutex.RLock()
	defer fake.getPackageObjectMutex.RUnlock()
	fake.getChaincodeDataMutex.RLock()
	defer fake.getChaincodeDataMutex.RUnlock()
	fake.getIdMutex.RLock()
	defer fake.getIdMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *CCPackage) recordInvocation(key string, args []interface{}) {
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
