
package mock

import (
	"sync"

	"github.com/mcc-github/blockchain/core/common/ccprovider"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

type FileSystemSupport struct {
	PutChaincodeToLocalStorageStub        func(ccprovider.CCPackage) error
	putChaincodeToLocalStorageMutex       sync.RWMutex
	putChaincodeToLocalStorageArgsForCall []struct {
		arg1 ccprovider.CCPackage
	}
	putChaincodeToLocalStorageReturns struct {
		result1 error
	}
	putChaincodeToLocalStorageReturnsOnCall map[int]struct {
		result1 error
	}
	GetChaincodeFromLocalStorageStub        func(ccname string, ccversion string) (ccprovider.CCPackage, error)
	getChaincodeFromLocalStorageMutex       sync.RWMutex
	getChaincodeFromLocalStorageArgsForCall []struct {
		ccname    string
		ccversion string
	}
	getChaincodeFromLocalStorageReturns struct {
		result1 ccprovider.CCPackage
		result2 error
	}
	getChaincodeFromLocalStorageReturnsOnCall map[int]struct {
		result1 ccprovider.CCPackage
		result2 error
	}
	GetChaincodesFromLocalStorageStub        func() (*pb.ChaincodeQueryResponse, error)
	getChaincodesFromLocalStorageMutex       sync.RWMutex
	getChaincodesFromLocalStorageArgsForCall []struct{}
	getChaincodesFromLocalStorageReturns     struct {
		result1 *pb.ChaincodeQueryResponse
		result2 error
	}
	getChaincodesFromLocalStorageReturnsOnCall map[int]struct {
		result1 *pb.ChaincodeQueryResponse
		result2 error
	}
	GetInstantiationPolicyStub        func(channel string, ccpack ccprovider.CCPackage) ([]byte, error)
	getInstantiationPolicyMutex       sync.RWMutex
	getInstantiationPolicyArgsForCall []struct {
		channel string
		ccpack  ccprovider.CCPackage
	}
	getInstantiationPolicyReturns struct {
		result1 []byte
		result2 error
	}
	getInstantiationPolicyReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	CheckInstantiationPolicyStub        func(signedProposal *pb.SignedProposal, chainName string, instantiationPolicy []byte) error
	checkInstantiationPolicyMutex       sync.RWMutex
	checkInstantiationPolicyArgsForCall []struct {
		signedProposal      *pb.SignedProposal
		chainName           string
		instantiationPolicy []byte
	}
	checkInstantiationPolicyReturns struct {
		result1 error
	}
	checkInstantiationPolicyReturnsOnCall map[int]struct {
		result1 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FileSystemSupport) PutChaincodeToLocalStorage(arg1 ccprovider.CCPackage) error {
	fake.putChaincodeToLocalStorageMutex.Lock()
	ret, specificReturn := fake.putChaincodeToLocalStorageReturnsOnCall[len(fake.putChaincodeToLocalStorageArgsForCall)]
	fake.putChaincodeToLocalStorageArgsForCall = append(fake.putChaincodeToLocalStorageArgsForCall, struct {
		arg1 ccprovider.CCPackage
	}{arg1})
	fake.recordInvocation("PutChaincodeToLocalStorage", []interface{}{arg1})
	fake.putChaincodeToLocalStorageMutex.Unlock()
	if fake.PutChaincodeToLocalStorageStub != nil {
		return fake.PutChaincodeToLocalStorageStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.putChaincodeToLocalStorageReturns.result1
}

func (fake *FileSystemSupport) PutChaincodeToLocalStorageCallCount() int {
	fake.putChaincodeToLocalStorageMutex.RLock()
	defer fake.putChaincodeToLocalStorageMutex.RUnlock()
	return len(fake.putChaincodeToLocalStorageArgsForCall)
}

func (fake *FileSystemSupport) PutChaincodeToLocalStorageArgsForCall(i int) ccprovider.CCPackage {
	fake.putChaincodeToLocalStorageMutex.RLock()
	defer fake.putChaincodeToLocalStorageMutex.RUnlock()
	return fake.putChaincodeToLocalStorageArgsForCall[i].arg1
}

func (fake *FileSystemSupport) PutChaincodeToLocalStorageReturns(result1 error) {
	fake.PutChaincodeToLocalStorageStub = nil
	fake.putChaincodeToLocalStorageReturns = struct {
		result1 error
	}{result1}
}

func (fake *FileSystemSupport) PutChaincodeToLocalStorageReturnsOnCall(i int, result1 error) {
	fake.PutChaincodeToLocalStorageStub = nil
	if fake.putChaincodeToLocalStorageReturnsOnCall == nil {
		fake.putChaincodeToLocalStorageReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.putChaincodeToLocalStorageReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FileSystemSupport) GetChaincodeFromLocalStorage(ccname string, ccversion string) (ccprovider.CCPackage, error) {
	fake.getChaincodeFromLocalStorageMutex.Lock()
	ret, specificReturn := fake.getChaincodeFromLocalStorageReturnsOnCall[len(fake.getChaincodeFromLocalStorageArgsForCall)]
	fake.getChaincodeFromLocalStorageArgsForCall = append(fake.getChaincodeFromLocalStorageArgsForCall, struct {
		ccname    string
		ccversion string
	}{ccname, ccversion})
	fake.recordInvocation("GetChaincodeFromLocalStorage", []interface{}{ccname, ccversion})
	fake.getChaincodeFromLocalStorageMutex.Unlock()
	if fake.GetChaincodeFromLocalStorageStub != nil {
		return fake.GetChaincodeFromLocalStorageStub(ccname, ccversion)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getChaincodeFromLocalStorageReturns.result1, fake.getChaincodeFromLocalStorageReturns.result2
}

func (fake *FileSystemSupport) GetChaincodeFromLocalStorageCallCount() int {
	fake.getChaincodeFromLocalStorageMutex.RLock()
	defer fake.getChaincodeFromLocalStorageMutex.RUnlock()
	return len(fake.getChaincodeFromLocalStorageArgsForCall)
}

func (fake *FileSystemSupport) GetChaincodeFromLocalStorageArgsForCall(i int) (string, string) {
	fake.getChaincodeFromLocalStorageMutex.RLock()
	defer fake.getChaincodeFromLocalStorageMutex.RUnlock()
	return fake.getChaincodeFromLocalStorageArgsForCall[i].ccname, fake.getChaincodeFromLocalStorageArgsForCall[i].ccversion
}

func (fake *FileSystemSupport) GetChaincodeFromLocalStorageReturns(result1 ccprovider.CCPackage, result2 error) {
	fake.GetChaincodeFromLocalStorageStub = nil
	fake.getChaincodeFromLocalStorageReturns = struct {
		result1 ccprovider.CCPackage
		result2 error
	}{result1, result2}
}

func (fake *FileSystemSupport) GetChaincodeFromLocalStorageReturnsOnCall(i int, result1 ccprovider.CCPackage, result2 error) {
	fake.GetChaincodeFromLocalStorageStub = nil
	if fake.getChaincodeFromLocalStorageReturnsOnCall == nil {
		fake.getChaincodeFromLocalStorageReturnsOnCall = make(map[int]struct {
			result1 ccprovider.CCPackage
			result2 error
		})
	}
	fake.getChaincodeFromLocalStorageReturnsOnCall[i] = struct {
		result1 ccprovider.CCPackage
		result2 error
	}{result1, result2}
}

func (fake *FileSystemSupport) GetChaincodesFromLocalStorage() (*pb.ChaincodeQueryResponse, error) {
	fake.getChaincodesFromLocalStorageMutex.Lock()
	ret, specificReturn := fake.getChaincodesFromLocalStorageReturnsOnCall[len(fake.getChaincodesFromLocalStorageArgsForCall)]
	fake.getChaincodesFromLocalStorageArgsForCall = append(fake.getChaincodesFromLocalStorageArgsForCall, struct{}{})
	fake.recordInvocation("GetChaincodesFromLocalStorage", []interface{}{})
	fake.getChaincodesFromLocalStorageMutex.Unlock()
	if fake.GetChaincodesFromLocalStorageStub != nil {
		return fake.GetChaincodesFromLocalStorageStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getChaincodesFromLocalStorageReturns.result1, fake.getChaincodesFromLocalStorageReturns.result2
}

func (fake *FileSystemSupport) GetChaincodesFromLocalStorageCallCount() int {
	fake.getChaincodesFromLocalStorageMutex.RLock()
	defer fake.getChaincodesFromLocalStorageMutex.RUnlock()
	return len(fake.getChaincodesFromLocalStorageArgsForCall)
}

func (fake *FileSystemSupport) GetChaincodesFromLocalStorageReturns(result1 *pb.ChaincodeQueryResponse, result2 error) {
	fake.GetChaincodesFromLocalStorageStub = nil
	fake.getChaincodesFromLocalStorageReturns = struct {
		result1 *pb.ChaincodeQueryResponse
		result2 error
	}{result1, result2}
}

func (fake *FileSystemSupport) GetChaincodesFromLocalStorageReturnsOnCall(i int, result1 *pb.ChaincodeQueryResponse, result2 error) {
	fake.GetChaincodesFromLocalStorageStub = nil
	if fake.getChaincodesFromLocalStorageReturnsOnCall == nil {
		fake.getChaincodesFromLocalStorageReturnsOnCall = make(map[int]struct {
			result1 *pb.ChaincodeQueryResponse
			result2 error
		})
	}
	fake.getChaincodesFromLocalStorageReturnsOnCall[i] = struct {
		result1 *pb.ChaincodeQueryResponse
		result2 error
	}{result1, result2}
}

func (fake *FileSystemSupport) GetInstantiationPolicy(channel string, ccpack ccprovider.CCPackage) ([]byte, error) {
	fake.getInstantiationPolicyMutex.Lock()
	ret, specificReturn := fake.getInstantiationPolicyReturnsOnCall[len(fake.getInstantiationPolicyArgsForCall)]
	fake.getInstantiationPolicyArgsForCall = append(fake.getInstantiationPolicyArgsForCall, struct {
		channel string
		ccpack  ccprovider.CCPackage
	}{channel, ccpack})
	fake.recordInvocation("GetInstantiationPolicy", []interface{}{channel, ccpack})
	fake.getInstantiationPolicyMutex.Unlock()
	if fake.GetInstantiationPolicyStub != nil {
		return fake.GetInstantiationPolicyStub(channel, ccpack)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fake.getInstantiationPolicyReturns.result1, fake.getInstantiationPolicyReturns.result2
}

func (fake *FileSystemSupport) GetInstantiationPolicyCallCount() int {
	fake.getInstantiationPolicyMutex.RLock()
	defer fake.getInstantiationPolicyMutex.RUnlock()
	return len(fake.getInstantiationPolicyArgsForCall)
}

func (fake *FileSystemSupport) GetInstantiationPolicyArgsForCall(i int) (string, ccprovider.CCPackage) {
	fake.getInstantiationPolicyMutex.RLock()
	defer fake.getInstantiationPolicyMutex.RUnlock()
	return fake.getInstantiationPolicyArgsForCall[i].channel, fake.getInstantiationPolicyArgsForCall[i].ccpack
}

func (fake *FileSystemSupport) GetInstantiationPolicyReturns(result1 []byte, result2 error) {
	fake.GetInstantiationPolicyStub = nil
	fake.getInstantiationPolicyReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FileSystemSupport) GetInstantiationPolicyReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.GetInstantiationPolicyStub = nil
	if fake.getInstantiationPolicyReturnsOnCall == nil {
		fake.getInstantiationPolicyReturnsOnCall = make(map[int]struct {
			result1 []byte
			result2 error
		})
	}
	fake.getInstantiationPolicyReturnsOnCall[i] = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FileSystemSupport) CheckInstantiationPolicy(signedProposal *pb.SignedProposal, chainName string, instantiationPolicy []byte) error {
	var instantiationPolicyCopy []byte
	if instantiationPolicy != nil {
		instantiationPolicyCopy = make([]byte, len(instantiationPolicy))
		copy(instantiationPolicyCopy, instantiationPolicy)
	}
	fake.checkInstantiationPolicyMutex.Lock()
	ret, specificReturn := fake.checkInstantiationPolicyReturnsOnCall[len(fake.checkInstantiationPolicyArgsForCall)]
	fake.checkInstantiationPolicyArgsForCall = append(fake.checkInstantiationPolicyArgsForCall, struct {
		signedProposal      *pb.SignedProposal
		chainName           string
		instantiationPolicy []byte
	}{signedProposal, chainName, instantiationPolicyCopy})
	fake.recordInvocation("CheckInstantiationPolicy", []interface{}{signedProposal, chainName, instantiationPolicyCopy})
	fake.checkInstantiationPolicyMutex.Unlock()
	if fake.CheckInstantiationPolicyStub != nil {
		return fake.CheckInstantiationPolicyStub(signedProposal, chainName, instantiationPolicy)
	}
	if specificReturn {
		return ret.result1
	}
	return fake.checkInstantiationPolicyReturns.result1
}

func (fake *FileSystemSupport) CheckInstantiationPolicyCallCount() int {
	fake.checkInstantiationPolicyMutex.RLock()
	defer fake.checkInstantiationPolicyMutex.RUnlock()
	return len(fake.checkInstantiationPolicyArgsForCall)
}

func (fake *FileSystemSupport) CheckInstantiationPolicyArgsForCall(i int) (*pb.SignedProposal, string, []byte) {
	fake.checkInstantiationPolicyMutex.RLock()
	defer fake.checkInstantiationPolicyMutex.RUnlock()
	return fake.checkInstantiationPolicyArgsForCall[i].signedProposal, fake.checkInstantiationPolicyArgsForCall[i].chainName, fake.checkInstantiationPolicyArgsForCall[i].instantiationPolicy
}

func (fake *FileSystemSupport) CheckInstantiationPolicyReturns(result1 error) {
	fake.CheckInstantiationPolicyStub = nil
	fake.checkInstantiationPolicyReturns = struct {
		result1 error
	}{result1}
}

func (fake *FileSystemSupport) CheckInstantiationPolicyReturnsOnCall(i int, result1 error) {
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

func (fake *FileSystemSupport) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.putChaincodeToLocalStorageMutex.RLock()
	defer fake.putChaincodeToLocalStorageMutex.RUnlock()
	fake.getChaincodeFromLocalStorageMutex.RLock()
	defer fake.getChaincodeFromLocalStorageMutex.RUnlock()
	fake.getChaincodesFromLocalStorageMutex.RLock()
	defer fake.getChaincodesFromLocalStorageMutex.RUnlock()
	fake.getInstantiationPolicyMutex.RLock()
	defer fake.getInstantiationPolicyMutex.RUnlock()
	fake.checkInstantiationPolicyMutex.RLock()
	defer fake.checkInstantiationPolicyMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FileSystemSupport) recordInvocation(key string, args []interface{}) {
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
