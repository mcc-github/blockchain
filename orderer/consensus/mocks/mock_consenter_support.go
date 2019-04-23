
package mocks

import (
	sync "sync"

	channelconfig "github.com/mcc-github/blockchain/common/channelconfig"
	blockcutter "github.com/mcc-github/blockchain/orderer/common/blockcutter"
	msgprocessor "github.com/mcc-github/blockchain/orderer/common/msgprocessor"
	consensus "github.com/mcc-github/blockchain/orderer/consensus"
	common "github.com/mcc-github/blockchain/protos/common"
	protoutil "github.com/mcc-github/blockchain/protoutil"
)

type FakeConsenterSupport struct {
	AppendStub        func(*common.Block) error
	appendMutex       sync.RWMutex
	appendArgsForCall []struct {
		arg1 *common.Block
	}
	appendReturns struct {
		result1 error
	}
	appendReturnsOnCall map[int]struct {
		result1 error
	}
	BlockStub        func(uint64) *common.Block
	blockMutex       sync.RWMutex
	blockArgsForCall []struct {
		arg1 uint64
	}
	blockReturns struct {
		result1 *common.Block
	}
	blockReturnsOnCall map[int]struct {
		result1 *common.Block
	}
	BlockCutterStub        func() blockcutter.Receiver
	blockCutterMutex       sync.RWMutex
	blockCutterArgsForCall []struct {
	}
	blockCutterReturns struct {
		result1 blockcutter.Receiver
	}
	blockCutterReturnsOnCall map[int]struct {
		result1 blockcutter.Receiver
	}
	ChainIDStub        func() string
	chainIDMutex       sync.RWMutex
	chainIDArgsForCall []struct {
	}
	chainIDReturns struct {
		result1 string
	}
	chainIDReturnsOnCall map[int]struct {
		result1 string
	}
	ChannelConfigStub        func() channelconfig.Channel
	channelConfigMutex       sync.RWMutex
	channelConfigArgsForCall []struct {
	}
	channelConfigReturns struct {
		result1 channelconfig.Channel
	}
	channelConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Channel
	}
	ClassifyMsgStub        func(*common.ChannelHeader) msgprocessor.Classification
	classifyMsgMutex       sync.RWMutex
	classifyMsgArgsForCall []struct {
		arg1 *common.ChannelHeader
	}
	classifyMsgReturns struct {
		result1 msgprocessor.Classification
	}
	classifyMsgReturnsOnCall map[int]struct {
		result1 msgprocessor.Classification
	}
	CreateNextBlockStub        func([]*common.Envelope) *common.Block
	createNextBlockMutex       sync.RWMutex
	createNextBlockArgsForCall []struct {
		arg1 []*common.Envelope
	}
	createNextBlockReturns struct {
		result1 *common.Block
	}
	createNextBlockReturnsOnCall map[int]struct {
		result1 *common.Block
	}
	HeightStub        func() uint64
	heightMutex       sync.RWMutex
	heightArgsForCall []struct {
	}
	heightReturns struct {
		result1 uint64
	}
	heightReturnsOnCall map[int]struct {
		result1 uint64
	}
	ProcessConfigMsgStub        func(*common.Envelope) (*common.Envelope, uint64, error)
	processConfigMsgMutex       sync.RWMutex
	processConfigMsgArgsForCall []struct {
		arg1 *common.Envelope
	}
	processConfigMsgReturns struct {
		result1 *common.Envelope
		result2 uint64
		result3 error
	}
	processConfigMsgReturnsOnCall map[int]struct {
		result1 *common.Envelope
		result2 uint64
		result3 error
	}
	ProcessConfigUpdateMsgStub        func(*common.Envelope) (*common.Envelope, uint64, error)
	processConfigUpdateMsgMutex       sync.RWMutex
	processConfigUpdateMsgArgsForCall []struct {
		arg1 *common.Envelope
	}
	processConfigUpdateMsgReturns struct {
		result1 *common.Envelope
		result2 uint64
		result3 error
	}
	processConfigUpdateMsgReturnsOnCall map[int]struct {
		result1 *common.Envelope
		result2 uint64
		result3 error
	}
	ProcessNormalMsgStub        func(*common.Envelope) (uint64, error)
	processNormalMsgMutex       sync.RWMutex
	processNormalMsgArgsForCall []struct {
		arg1 *common.Envelope
	}
	processNormalMsgReturns struct {
		result1 uint64
		result2 error
	}
	processNormalMsgReturnsOnCall map[int]struct {
		result1 uint64
		result2 error
	}
	SequenceStub        func() uint64
	sequenceMutex       sync.RWMutex
	sequenceArgsForCall []struct {
	}
	sequenceReturns struct {
		result1 uint64
	}
	sequenceReturnsOnCall map[int]struct {
		result1 uint64
	}
	SerializeStub        func() ([]byte, error)
	serializeMutex       sync.RWMutex
	serializeArgsForCall []struct {
	}
	serializeReturns struct {
		result1 []byte
		result2 error
	}
	serializeReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	SharedConfigStub        func() channelconfig.Orderer
	sharedConfigMutex       sync.RWMutex
	sharedConfigArgsForCall []struct {
	}
	sharedConfigReturns struct {
		result1 channelconfig.Orderer
	}
	sharedConfigReturnsOnCall map[int]struct {
		result1 channelconfig.Orderer
	}
	SignStub        func([]byte) ([]byte, error)
	signMutex       sync.RWMutex
	signArgsForCall []struct {
		arg1 []byte
	}
	signReturns struct {
		result1 []byte
		result2 error
	}
	signReturnsOnCall map[int]struct {
		result1 []byte
		result2 error
	}
	VerifyBlockSignatureStub        func([]*protoutil.SignedData, *common.ConfigEnvelope) error
	verifyBlockSignatureMutex       sync.RWMutex
	verifyBlockSignatureArgsForCall []struct {
		arg1 []*protoutil.SignedData
		arg2 *common.ConfigEnvelope
	}
	verifyBlockSignatureReturns struct {
		result1 error
	}
	verifyBlockSignatureReturnsOnCall map[int]struct {
		result1 error
	}
	WriteBlockStub        func(*common.Block, []byte)
	writeBlockMutex       sync.RWMutex
	writeBlockArgsForCall []struct {
		arg1 *common.Block
		arg2 []byte
	}
	WriteConfigBlockStub        func(*common.Block, []byte)
	writeConfigBlockMutex       sync.RWMutex
	writeConfigBlockArgsForCall []struct {
		arg1 *common.Block
		arg2 []byte
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakeConsenterSupport) Append(arg1 *common.Block) error {
	fake.appendMutex.Lock()
	ret, specificReturn := fake.appendReturnsOnCall[len(fake.appendArgsForCall)]
	fake.appendArgsForCall = append(fake.appendArgsForCall, struct {
		arg1 *common.Block
	}{arg1})
	fake.recordInvocation("Append", []interface{}{arg1})
	fake.appendMutex.Unlock()
	if fake.AppendStub != nil {
		return fake.AppendStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.appendReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) AppendCallCount() int {
	fake.appendMutex.RLock()
	defer fake.appendMutex.RUnlock()
	return len(fake.appendArgsForCall)
}

func (fake *FakeConsenterSupport) AppendCalls(stub func(*common.Block) error) {
	fake.appendMutex.Lock()
	defer fake.appendMutex.Unlock()
	fake.AppendStub = stub
}

func (fake *FakeConsenterSupport) AppendArgsForCall(i int) *common.Block {
	fake.appendMutex.RLock()
	defer fake.appendMutex.RUnlock()
	argsForCall := fake.appendArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsenterSupport) AppendReturns(result1 error) {
	fake.appendMutex.Lock()
	defer fake.appendMutex.Unlock()
	fake.AppendStub = nil
	fake.appendReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsenterSupport) AppendReturnsOnCall(i int, result1 error) {
	fake.appendMutex.Lock()
	defer fake.appendMutex.Unlock()
	fake.AppendStub = nil
	if fake.appendReturnsOnCall == nil {
		fake.appendReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.appendReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsenterSupport) Block(arg1 uint64) *common.Block {
	fake.blockMutex.Lock()
	ret, specificReturn := fake.blockReturnsOnCall[len(fake.blockArgsForCall)]
	fake.blockArgsForCall = append(fake.blockArgsForCall, struct {
		arg1 uint64
	}{arg1})
	fake.recordInvocation("Block", []interface{}{arg1})
	fake.blockMutex.Unlock()
	if fake.BlockStub != nil {
		return fake.BlockStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.blockReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) BlockCallCount() int {
	fake.blockMutex.RLock()
	defer fake.blockMutex.RUnlock()
	return len(fake.blockArgsForCall)
}

func (fake *FakeConsenterSupport) BlockCalls(stub func(uint64) *common.Block) {
	fake.blockMutex.Lock()
	defer fake.blockMutex.Unlock()
	fake.BlockStub = stub
}

func (fake *FakeConsenterSupport) BlockArgsForCall(i int) uint64 {
	fake.blockMutex.RLock()
	defer fake.blockMutex.RUnlock()
	argsForCall := fake.blockArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsenterSupport) BlockReturns(result1 *common.Block) {
	fake.blockMutex.Lock()
	defer fake.blockMutex.Unlock()
	fake.BlockStub = nil
	fake.blockReturns = struct {
		result1 *common.Block
	}{result1}
}

func (fake *FakeConsenterSupport) BlockReturnsOnCall(i int, result1 *common.Block) {
	fake.blockMutex.Lock()
	defer fake.blockMutex.Unlock()
	fake.BlockStub = nil
	if fake.blockReturnsOnCall == nil {
		fake.blockReturnsOnCall = make(map[int]struct {
			result1 *common.Block
		})
	}
	fake.blockReturnsOnCall[i] = struct {
		result1 *common.Block
	}{result1}
}

func (fake *FakeConsenterSupport) BlockCutter() blockcutter.Receiver {
	fake.blockCutterMutex.Lock()
	ret, specificReturn := fake.blockCutterReturnsOnCall[len(fake.blockCutterArgsForCall)]
	fake.blockCutterArgsForCall = append(fake.blockCutterArgsForCall, struct {
	}{})
	fake.recordInvocation("BlockCutter", []interface{}{})
	fake.blockCutterMutex.Unlock()
	if fake.BlockCutterStub != nil {
		return fake.BlockCutterStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.blockCutterReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) BlockCutterCallCount() int {
	fake.blockCutterMutex.RLock()
	defer fake.blockCutterMutex.RUnlock()
	return len(fake.blockCutterArgsForCall)
}

func (fake *FakeConsenterSupport) BlockCutterCalls(stub func() blockcutter.Receiver) {
	fake.blockCutterMutex.Lock()
	defer fake.blockCutterMutex.Unlock()
	fake.BlockCutterStub = stub
}

func (fake *FakeConsenterSupport) BlockCutterReturns(result1 blockcutter.Receiver) {
	fake.blockCutterMutex.Lock()
	defer fake.blockCutterMutex.Unlock()
	fake.BlockCutterStub = nil
	fake.blockCutterReturns = struct {
		result1 blockcutter.Receiver
	}{result1}
}

func (fake *FakeConsenterSupport) BlockCutterReturnsOnCall(i int, result1 blockcutter.Receiver) {
	fake.blockCutterMutex.Lock()
	defer fake.blockCutterMutex.Unlock()
	fake.BlockCutterStub = nil
	if fake.blockCutterReturnsOnCall == nil {
		fake.blockCutterReturnsOnCall = make(map[int]struct {
			result1 blockcutter.Receiver
		})
	}
	fake.blockCutterReturnsOnCall[i] = struct {
		result1 blockcutter.Receiver
	}{result1}
}

func (fake *FakeConsenterSupport) ChainID() string {
	fake.chainIDMutex.Lock()
	ret, specificReturn := fake.chainIDReturnsOnCall[len(fake.chainIDArgsForCall)]
	fake.chainIDArgsForCall = append(fake.chainIDArgsForCall, struct {
	}{})
	fake.recordInvocation("ChainID", []interface{}{})
	fake.chainIDMutex.Unlock()
	if fake.ChainIDStub != nil {
		return fake.ChainIDStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.chainIDReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) ChainIDCallCount() int {
	fake.chainIDMutex.RLock()
	defer fake.chainIDMutex.RUnlock()
	return len(fake.chainIDArgsForCall)
}

func (fake *FakeConsenterSupport) ChainIDCalls(stub func() string) {
	fake.chainIDMutex.Lock()
	defer fake.chainIDMutex.Unlock()
	fake.ChainIDStub = stub
}

func (fake *FakeConsenterSupport) ChainIDReturns(result1 string) {
	fake.chainIDMutex.Lock()
	defer fake.chainIDMutex.Unlock()
	fake.ChainIDStub = nil
	fake.chainIDReturns = struct {
		result1 string
	}{result1}
}

func (fake *FakeConsenterSupport) ChainIDReturnsOnCall(i int, result1 string) {
	fake.chainIDMutex.Lock()
	defer fake.chainIDMutex.Unlock()
	fake.ChainIDStub = nil
	if fake.chainIDReturnsOnCall == nil {
		fake.chainIDReturnsOnCall = make(map[int]struct {
			result1 string
		})
	}
	fake.chainIDReturnsOnCall[i] = struct {
		result1 string
	}{result1}
}

func (fake *FakeConsenterSupport) ChannelConfig() channelconfig.Channel {
	fake.channelConfigMutex.Lock()
	ret, specificReturn := fake.channelConfigReturnsOnCall[len(fake.channelConfigArgsForCall)]
	fake.channelConfigArgsForCall = append(fake.channelConfigArgsForCall, struct {
	}{})
	fake.recordInvocation("ChannelConfig", []interface{}{})
	fake.channelConfigMutex.Unlock()
	if fake.ChannelConfigStub != nil {
		return fake.ChannelConfigStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.channelConfigReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) ChannelConfigCallCount() int {
	fake.channelConfigMutex.RLock()
	defer fake.channelConfigMutex.RUnlock()
	return len(fake.channelConfigArgsForCall)
}

func (fake *FakeConsenterSupport) ChannelConfigCalls(stub func() channelconfig.Channel) {
	fake.channelConfigMutex.Lock()
	defer fake.channelConfigMutex.Unlock()
	fake.ChannelConfigStub = stub
}

func (fake *FakeConsenterSupport) ChannelConfigReturns(result1 channelconfig.Channel) {
	fake.channelConfigMutex.Lock()
	defer fake.channelConfigMutex.Unlock()
	fake.ChannelConfigStub = nil
	fake.channelConfigReturns = struct {
		result1 channelconfig.Channel
	}{result1}
}

func (fake *FakeConsenterSupport) ChannelConfigReturnsOnCall(i int, result1 channelconfig.Channel) {
	fake.channelConfigMutex.Lock()
	defer fake.channelConfigMutex.Unlock()
	fake.ChannelConfigStub = nil
	if fake.channelConfigReturnsOnCall == nil {
		fake.channelConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Channel
		})
	}
	fake.channelConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Channel
	}{result1}
}

func (fake *FakeConsenterSupport) ClassifyMsg(arg1 *common.ChannelHeader) msgprocessor.Classification {
	fake.classifyMsgMutex.Lock()
	ret, specificReturn := fake.classifyMsgReturnsOnCall[len(fake.classifyMsgArgsForCall)]
	fake.classifyMsgArgsForCall = append(fake.classifyMsgArgsForCall, struct {
		arg1 *common.ChannelHeader
	}{arg1})
	fake.recordInvocation("ClassifyMsg", []interface{}{arg1})
	fake.classifyMsgMutex.Unlock()
	if fake.ClassifyMsgStub != nil {
		return fake.ClassifyMsgStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.classifyMsgReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) ClassifyMsgCallCount() int {
	fake.classifyMsgMutex.RLock()
	defer fake.classifyMsgMutex.RUnlock()
	return len(fake.classifyMsgArgsForCall)
}

func (fake *FakeConsenterSupport) ClassifyMsgCalls(stub func(*common.ChannelHeader) msgprocessor.Classification) {
	fake.classifyMsgMutex.Lock()
	defer fake.classifyMsgMutex.Unlock()
	fake.ClassifyMsgStub = stub
}

func (fake *FakeConsenterSupport) ClassifyMsgArgsForCall(i int) *common.ChannelHeader {
	fake.classifyMsgMutex.RLock()
	defer fake.classifyMsgMutex.RUnlock()
	argsForCall := fake.classifyMsgArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsenterSupport) ClassifyMsgReturns(result1 msgprocessor.Classification) {
	fake.classifyMsgMutex.Lock()
	defer fake.classifyMsgMutex.Unlock()
	fake.ClassifyMsgStub = nil
	fake.classifyMsgReturns = struct {
		result1 msgprocessor.Classification
	}{result1}
}

func (fake *FakeConsenterSupport) ClassifyMsgReturnsOnCall(i int, result1 msgprocessor.Classification) {
	fake.classifyMsgMutex.Lock()
	defer fake.classifyMsgMutex.Unlock()
	fake.ClassifyMsgStub = nil
	if fake.classifyMsgReturnsOnCall == nil {
		fake.classifyMsgReturnsOnCall = make(map[int]struct {
			result1 msgprocessor.Classification
		})
	}
	fake.classifyMsgReturnsOnCall[i] = struct {
		result1 msgprocessor.Classification
	}{result1}
}

func (fake *FakeConsenterSupport) CreateNextBlock(arg1 []*common.Envelope) *common.Block {
	var arg1Copy []*common.Envelope
	if arg1 != nil {
		arg1Copy = make([]*common.Envelope, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.createNextBlockMutex.Lock()
	ret, specificReturn := fake.createNextBlockReturnsOnCall[len(fake.createNextBlockArgsForCall)]
	fake.createNextBlockArgsForCall = append(fake.createNextBlockArgsForCall, struct {
		arg1 []*common.Envelope
	}{arg1Copy})
	fake.recordInvocation("CreateNextBlock", []interface{}{arg1Copy})
	fake.createNextBlockMutex.Unlock()
	if fake.CreateNextBlockStub != nil {
		return fake.CreateNextBlockStub(arg1)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.createNextBlockReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) CreateNextBlockCallCount() int {
	fake.createNextBlockMutex.RLock()
	defer fake.createNextBlockMutex.RUnlock()
	return len(fake.createNextBlockArgsForCall)
}

func (fake *FakeConsenterSupport) CreateNextBlockCalls(stub func([]*common.Envelope) *common.Block) {
	fake.createNextBlockMutex.Lock()
	defer fake.createNextBlockMutex.Unlock()
	fake.CreateNextBlockStub = stub
}

func (fake *FakeConsenterSupport) CreateNextBlockArgsForCall(i int) []*common.Envelope {
	fake.createNextBlockMutex.RLock()
	defer fake.createNextBlockMutex.RUnlock()
	argsForCall := fake.createNextBlockArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsenterSupport) CreateNextBlockReturns(result1 *common.Block) {
	fake.createNextBlockMutex.Lock()
	defer fake.createNextBlockMutex.Unlock()
	fake.CreateNextBlockStub = nil
	fake.createNextBlockReturns = struct {
		result1 *common.Block
	}{result1}
}

func (fake *FakeConsenterSupport) CreateNextBlockReturnsOnCall(i int, result1 *common.Block) {
	fake.createNextBlockMutex.Lock()
	defer fake.createNextBlockMutex.Unlock()
	fake.CreateNextBlockStub = nil
	if fake.createNextBlockReturnsOnCall == nil {
		fake.createNextBlockReturnsOnCall = make(map[int]struct {
			result1 *common.Block
		})
	}
	fake.createNextBlockReturnsOnCall[i] = struct {
		result1 *common.Block
	}{result1}
}

func (fake *FakeConsenterSupport) Height() uint64 {
	fake.heightMutex.Lock()
	ret, specificReturn := fake.heightReturnsOnCall[len(fake.heightArgsForCall)]
	fake.heightArgsForCall = append(fake.heightArgsForCall, struct {
	}{})
	fake.recordInvocation("Height", []interface{}{})
	fake.heightMutex.Unlock()
	if fake.HeightStub != nil {
		return fake.HeightStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.heightReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) HeightCallCount() int {
	fake.heightMutex.RLock()
	defer fake.heightMutex.RUnlock()
	return len(fake.heightArgsForCall)
}

func (fake *FakeConsenterSupport) HeightCalls(stub func() uint64) {
	fake.heightMutex.Lock()
	defer fake.heightMutex.Unlock()
	fake.HeightStub = stub
}

func (fake *FakeConsenterSupport) HeightReturns(result1 uint64) {
	fake.heightMutex.Lock()
	defer fake.heightMutex.Unlock()
	fake.HeightStub = nil
	fake.heightReturns = struct {
		result1 uint64
	}{result1}
}

func (fake *FakeConsenterSupport) HeightReturnsOnCall(i int, result1 uint64) {
	fake.heightMutex.Lock()
	defer fake.heightMutex.Unlock()
	fake.HeightStub = nil
	if fake.heightReturnsOnCall == nil {
		fake.heightReturnsOnCall = make(map[int]struct {
			result1 uint64
		})
	}
	fake.heightReturnsOnCall[i] = struct {
		result1 uint64
	}{result1}
}

func (fake *FakeConsenterSupport) ProcessConfigMsg(arg1 *common.Envelope) (*common.Envelope, uint64, error) {
	fake.processConfigMsgMutex.Lock()
	ret, specificReturn := fake.processConfigMsgReturnsOnCall[len(fake.processConfigMsgArgsForCall)]
	fake.processConfigMsgArgsForCall = append(fake.processConfigMsgArgsForCall, struct {
		arg1 *common.Envelope
	}{arg1})
	fake.recordInvocation("ProcessConfigMsg", []interface{}{arg1})
	fake.processConfigMsgMutex.Unlock()
	if fake.ProcessConfigMsgStub != nil {
		return fake.ProcessConfigMsgStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.processConfigMsgReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeConsenterSupport) ProcessConfigMsgCallCount() int {
	fake.processConfigMsgMutex.RLock()
	defer fake.processConfigMsgMutex.RUnlock()
	return len(fake.processConfigMsgArgsForCall)
}

func (fake *FakeConsenterSupport) ProcessConfigMsgCalls(stub func(*common.Envelope) (*common.Envelope, uint64, error)) {
	fake.processConfigMsgMutex.Lock()
	defer fake.processConfigMsgMutex.Unlock()
	fake.ProcessConfigMsgStub = stub
}

func (fake *FakeConsenterSupport) ProcessConfigMsgArgsForCall(i int) *common.Envelope {
	fake.processConfigMsgMutex.RLock()
	defer fake.processConfigMsgMutex.RUnlock()
	argsForCall := fake.processConfigMsgArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsenterSupport) ProcessConfigMsgReturns(result1 *common.Envelope, result2 uint64, result3 error) {
	fake.processConfigMsgMutex.Lock()
	defer fake.processConfigMsgMutex.Unlock()
	fake.ProcessConfigMsgStub = nil
	fake.processConfigMsgReturns = struct {
		result1 *common.Envelope
		result2 uint64
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeConsenterSupport) ProcessConfigMsgReturnsOnCall(i int, result1 *common.Envelope, result2 uint64, result3 error) {
	fake.processConfigMsgMutex.Lock()
	defer fake.processConfigMsgMutex.Unlock()
	fake.ProcessConfigMsgStub = nil
	if fake.processConfigMsgReturnsOnCall == nil {
		fake.processConfigMsgReturnsOnCall = make(map[int]struct {
			result1 *common.Envelope
			result2 uint64
			result3 error
		})
	}
	fake.processConfigMsgReturnsOnCall[i] = struct {
		result1 *common.Envelope
		result2 uint64
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeConsenterSupport) ProcessConfigUpdateMsg(arg1 *common.Envelope) (*common.Envelope, uint64, error) {
	fake.processConfigUpdateMsgMutex.Lock()
	ret, specificReturn := fake.processConfigUpdateMsgReturnsOnCall[len(fake.processConfigUpdateMsgArgsForCall)]
	fake.processConfigUpdateMsgArgsForCall = append(fake.processConfigUpdateMsgArgsForCall, struct {
		arg1 *common.Envelope
	}{arg1})
	fake.recordInvocation("ProcessConfigUpdateMsg", []interface{}{arg1})
	fake.processConfigUpdateMsgMutex.Unlock()
	if fake.ProcessConfigUpdateMsgStub != nil {
		return fake.ProcessConfigUpdateMsgStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3
	}
	fakeReturns := fake.processConfigUpdateMsgReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3
}

func (fake *FakeConsenterSupport) ProcessConfigUpdateMsgCallCount() int {
	fake.processConfigUpdateMsgMutex.RLock()
	defer fake.processConfigUpdateMsgMutex.RUnlock()
	return len(fake.processConfigUpdateMsgArgsForCall)
}

func (fake *FakeConsenterSupport) ProcessConfigUpdateMsgCalls(stub func(*common.Envelope) (*common.Envelope, uint64, error)) {
	fake.processConfigUpdateMsgMutex.Lock()
	defer fake.processConfigUpdateMsgMutex.Unlock()
	fake.ProcessConfigUpdateMsgStub = stub
}

func (fake *FakeConsenterSupport) ProcessConfigUpdateMsgArgsForCall(i int) *common.Envelope {
	fake.processConfigUpdateMsgMutex.RLock()
	defer fake.processConfigUpdateMsgMutex.RUnlock()
	argsForCall := fake.processConfigUpdateMsgArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsenterSupport) ProcessConfigUpdateMsgReturns(result1 *common.Envelope, result2 uint64, result3 error) {
	fake.processConfigUpdateMsgMutex.Lock()
	defer fake.processConfigUpdateMsgMutex.Unlock()
	fake.ProcessConfigUpdateMsgStub = nil
	fake.processConfigUpdateMsgReturns = struct {
		result1 *common.Envelope
		result2 uint64
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeConsenterSupport) ProcessConfigUpdateMsgReturnsOnCall(i int, result1 *common.Envelope, result2 uint64, result3 error) {
	fake.processConfigUpdateMsgMutex.Lock()
	defer fake.processConfigUpdateMsgMutex.Unlock()
	fake.ProcessConfigUpdateMsgStub = nil
	if fake.processConfigUpdateMsgReturnsOnCall == nil {
		fake.processConfigUpdateMsgReturnsOnCall = make(map[int]struct {
			result1 *common.Envelope
			result2 uint64
			result3 error
		})
	}
	fake.processConfigUpdateMsgReturnsOnCall[i] = struct {
		result1 *common.Envelope
		result2 uint64
		result3 error
	}{result1, result2, result3}
}

func (fake *FakeConsenterSupport) ProcessNormalMsg(arg1 *common.Envelope) (uint64, error) {
	fake.processNormalMsgMutex.Lock()
	ret, specificReturn := fake.processNormalMsgReturnsOnCall[len(fake.processNormalMsgArgsForCall)]
	fake.processNormalMsgArgsForCall = append(fake.processNormalMsgArgsForCall, struct {
		arg1 *common.Envelope
	}{arg1})
	fake.recordInvocation("ProcessNormalMsg", []interface{}{arg1})
	fake.processNormalMsgMutex.Unlock()
	if fake.ProcessNormalMsgStub != nil {
		return fake.ProcessNormalMsgStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.processNormalMsgReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeConsenterSupport) ProcessNormalMsgCallCount() int {
	fake.processNormalMsgMutex.RLock()
	defer fake.processNormalMsgMutex.RUnlock()
	return len(fake.processNormalMsgArgsForCall)
}

func (fake *FakeConsenterSupport) ProcessNormalMsgCalls(stub func(*common.Envelope) (uint64, error)) {
	fake.processNormalMsgMutex.Lock()
	defer fake.processNormalMsgMutex.Unlock()
	fake.ProcessNormalMsgStub = stub
}

func (fake *FakeConsenterSupport) ProcessNormalMsgArgsForCall(i int) *common.Envelope {
	fake.processNormalMsgMutex.RLock()
	defer fake.processNormalMsgMutex.RUnlock()
	argsForCall := fake.processNormalMsgArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsenterSupport) ProcessNormalMsgReturns(result1 uint64, result2 error) {
	fake.processNormalMsgMutex.Lock()
	defer fake.processNormalMsgMutex.Unlock()
	fake.ProcessNormalMsgStub = nil
	fake.processNormalMsgReturns = struct {
		result1 uint64
		result2 error
	}{result1, result2}
}

func (fake *FakeConsenterSupport) ProcessNormalMsgReturnsOnCall(i int, result1 uint64, result2 error) {
	fake.processNormalMsgMutex.Lock()
	defer fake.processNormalMsgMutex.Unlock()
	fake.ProcessNormalMsgStub = nil
	if fake.processNormalMsgReturnsOnCall == nil {
		fake.processNormalMsgReturnsOnCall = make(map[int]struct {
			result1 uint64
			result2 error
		})
	}
	fake.processNormalMsgReturnsOnCall[i] = struct {
		result1 uint64
		result2 error
	}{result1, result2}
}

func (fake *FakeConsenterSupport) Sequence() uint64 {
	fake.sequenceMutex.Lock()
	ret, specificReturn := fake.sequenceReturnsOnCall[len(fake.sequenceArgsForCall)]
	fake.sequenceArgsForCall = append(fake.sequenceArgsForCall, struct {
	}{})
	fake.recordInvocation("Sequence", []interface{}{})
	fake.sequenceMutex.Unlock()
	if fake.SequenceStub != nil {
		return fake.SequenceStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sequenceReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) SequenceCallCount() int {
	fake.sequenceMutex.RLock()
	defer fake.sequenceMutex.RUnlock()
	return len(fake.sequenceArgsForCall)
}

func (fake *FakeConsenterSupport) SequenceCalls(stub func() uint64) {
	fake.sequenceMutex.Lock()
	defer fake.sequenceMutex.Unlock()
	fake.SequenceStub = stub
}

func (fake *FakeConsenterSupport) SequenceReturns(result1 uint64) {
	fake.sequenceMutex.Lock()
	defer fake.sequenceMutex.Unlock()
	fake.SequenceStub = nil
	fake.sequenceReturns = struct {
		result1 uint64
	}{result1}
}

func (fake *FakeConsenterSupport) SequenceReturnsOnCall(i int, result1 uint64) {
	fake.sequenceMutex.Lock()
	defer fake.sequenceMutex.Unlock()
	fake.SequenceStub = nil
	if fake.sequenceReturnsOnCall == nil {
		fake.sequenceReturnsOnCall = make(map[int]struct {
			result1 uint64
		})
	}
	fake.sequenceReturnsOnCall[i] = struct {
		result1 uint64
	}{result1}
}

func (fake *FakeConsenterSupport) Serialize() ([]byte, error) {
	fake.serializeMutex.Lock()
	ret, specificReturn := fake.serializeReturnsOnCall[len(fake.serializeArgsForCall)]
	fake.serializeArgsForCall = append(fake.serializeArgsForCall, struct {
	}{})
	fake.recordInvocation("Serialize", []interface{}{})
	fake.serializeMutex.Unlock()
	if fake.SerializeStub != nil {
		return fake.SerializeStub()
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.serializeReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeConsenterSupport) SerializeCallCount() int {
	fake.serializeMutex.RLock()
	defer fake.serializeMutex.RUnlock()
	return len(fake.serializeArgsForCall)
}

func (fake *FakeConsenterSupport) SerializeCalls(stub func() ([]byte, error)) {
	fake.serializeMutex.Lock()
	defer fake.serializeMutex.Unlock()
	fake.SerializeStub = stub
}

func (fake *FakeConsenterSupport) SerializeReturns(result1 []byte, result2 error) {
	fake.serializeMutex.Lock()
	defer fake.serializeMutex.Unlock()
	fake.SerializeStub = nil
	fake.serializeReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeConsenterSupport) SerializeReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.serializeMutex.Lock()
	defer fake.serializeMutex.Unlock()
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

func (fake *FakeConsenterSupport) SharedConfig() channelconfig.Orderer {
	fake.sharedConfigMutex.Lock()
	ret, specificReturn := fake.sharedConfigReturnsOnCall[len(fake.sharedConfigArgsForCall)]
	fake.sharedConfigArgsForCall = append(fake.sharedConfigArgsForCall, struct {
	}{})
	fake.recordInvocation("SharedConfig", []interface{}{})
	fake.sharedConfigMutex.Unlock()
	if fake.SharedConfigStub != nil {
		return fake.SharedConfigStub()
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.sharedConfigReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) SharedConfigCallCount() int {
	fake.sharedConfigMutex.RLock()
	defer fake.sharedConfigMutex.RUnlock()
	return len(fake.sharedConfigArgsForCall)
}

func (fake *FakeConsenterSupport) SharedConfigCalls(stub func() channelconfig.Orderer) {
	fake.sharedConfigMutex.Lock()
	defer fake.sharedConfigMutex.Unlock()
	fake.SharedConfigStub = stub
}

func (fake *FakeConsenterSupport) SharedConfigReturns(result1 channelconfig.Orderer) {
	fake.sharedConfigMutex.Lock()
	defer fake.sharedConfigMutex.Unlock()
	fake.SharedConfigStub = nil
	fake.sharedConfigReturns = struct {
		result1 channelconfig.Orderer
	}{result1}
}

func (fake *FakeConsenterSupport) SharedConfigReturnsOnCall(i int, result1 channelconfig.Orderer) {
	fake.sharedConfigMutex.Lock()
	defer fake.sharedConfigMutex.Unlock()
	fake.SharedConfigStub = nil
	if fake.sharedConfigReturnsOnCall == nil {
		fake.sharedConfigReturnsOnCall = make(map[int]struct {
			result1 channelconfig.Orderer
		})
	}
	fake.sharedConfigReturnsOnCall[i] = struct {
		result1 channelconfig.Orderer
	}{result1}
}

func (fake *FakeConsenterSupport) Sign(arg1 []byte) ([]byte, error) {
	var arg1Copy []byte
	if arg1 != nil {
		arg1Copy = make([]byte, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.signMutex.Lock()
	ret, specificReturn := fake.signReturnsOnCall[len(fake.signArgsForCall)]
	fake.signArgsForCall = append(fake.signArgsForCall, struct {
		arg1 []byte
	}{arg1Copy})
	fake.recordInvocation("Sign", []interface{}{arg1Copy})
	fake.signMutex.Unlock()
	if fake.SignStub != nil {
		return fake.SignStub(arg1)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	fakeReturns := fake.signReturns
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakeConsenterSupport) SignCallCount() int {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	return len(fake.signArgsForCall)
}

func (fake *FakeConsenterSupport) SignCalls(stub func([]byte) ([]byte, error)) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
	fake.SignStub = stub
}

func (fake *FakeConsenterSupport) SignArgsForCall(i int) []byte {
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	argsForCall := fake.signArgsForCall[i]
	return argsForCall.arg1
}

func (fake *FakeConsenterSupport) SignReturns(result1 []byte, result2 error) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
	fake.SignStub = nil
	fake.signReturns = struct {
		result1 []byte
		result2 error
	}{result1, result2}
}

func (fake *FakeConsenterSupport) SignReturnsOnCall(i int, result1 []byte, result2 error) {
	fake.signMutex.Lock()
	defer fake.signMutex.Unlock()
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

func (fake *FakeConsenterSupport) VerifyBlockSignature(arg1 []*protoutil.SignedData, arg2 *common.ConfigEnvelope) error {
	var arg1Copy []*protoutil.SignedData
	if arg1 != nil {
		arg1Copy = make([]*protoutil.SignedData, len(arg1))
		copy(arg1Copy, arg1)
	}
	fake.verifyBlockSignatureMutex.Lock()
	ret, specificReturn := fake.verifyBlockSignatureReturnsOnCall[len(fake.verifyBlockSignatureArgsForCall)]
	fake.verifyBlockSignatureArgsForCall = append(fake.verifyBlockSignatureArgsForCall, struct {
		arg1 []*protoutil.SignedData
		arg2 *common.ConfigEnvelope
	}{arg1Copy, arg2})
	fake.recordInvocation("VerifyBlockSignature", []interface{}{arg1Copy, arg2})
	fake.verifyBlockSignatureMutex.Unlock()
	if fake.VerifyBlockSignatureStub != nil {
		return fake.VerifyBlockSignatureStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1
	}
	fakeReturns := fake.verifyBlockSignatureReturns
	return fakeReturns.result1
}

func (fake *FakeConsenterSupport) VerifyBlockSignatureCallCount() int {
	fake.verifyBlockSignatureMutex.RLock()
	defer fake.verifyBlockSignatureMutex.RUnlock()
	return len(fake.verifyBlockSignatureArgsForCall)
}

func (fake *FakeConsenterSupport) VerifyBlockSignatureCalls(stub func([]*protoutil.SignedData, *common.ConfigEnvelope) error) {
	fake.verifyBlockSignatureMutex.Lock()
	defer fake.verifyBlockSignatureMutex.Unlock()
	fake.VerifyBlockSignatureStub = stub
}

func (fake *FakeConsenterSupport) VerifyBlockSignatureArgsForCall(i int) ([]*protoutil.SignedData, *common.ConfigEnvelope) {
	fake.verifyBlockSignatureMutex.RLock()
	defer fake.verifyBlockSignatureMutex.RUnlock()
	argsForCall := fake.verifyBlockSignatureArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeConsenterSupport) VerifyBlockSignatureReturns(result1 error) {
	fake.verifyBlockSignatureMutex.Lock()
	defer fake.verifyBlockSignatureMutex.Unlock()
	fake.VerifyBlockSignatureStub = nil
	fake.verifyBlockSignatureReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsenterSupport) VerifyBlockSignatureReturnsOnCall(i int, result1 error) {
	fake.verifyBlockSignatureMutex.Lock()
	defer fake.verifyBlockSignatureMutex.Unlock()
	fake.VerifyBlockSignatureStub = nil
	if fake.verifyBlockSignatureReturnsOnCall == nil {
		fake.verifyBlockSignatureReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.verifyBlockSignatureReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakeConsenterSupport) WriteBlock(arg1 *common.Block, arg2 []byte) {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.writeBlockMutex.Lock()
	fake.writeBlockArgsForCall = append(fake.writeBlockArgsForCall, struct {
		arg1 *common.Block
		arg2 []byte
	}{arg1, arg2Copy})
	fake.recordInvocation("WriteBlock", []interface{}{arg1, arg2Copy})
	fake.writeBlockMutex.Unlock()
	if fake.WriteBlockStub != nil {
		fake.WriteBlockStub(arg1, arg2)
	}
}

func (fake *FakeConsenterSupport) WriteBlockCallCount() int {
	fake.writeBlockMutex.RLock()
	defer fake.writeBlockMutex.RUnlock()
	return len(fake.writeBlockArgsForCall)
}

func (fake *FakeConsenterSupport) WriteBlockCalls(stub func(*common.Block, []byte)) {
	fake.writeBlockMutex.Lock()
	defer fake.writeBlockMutex.Unlock()
	fake.WriteBlockStub = stub
}

func (fake *FakeConsenterSupport) WriteBlockArgsForCall(i int) (*common.Block, []byte) {
	fake.writeBlockMutex.RLock()
	defer fake.writeBlockMutex.RUnlock()
	argsForCall := fake.writeBlockArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeConsenterSupport) WriteConfigBlock(arg1 *common.Block, arg2 []byte) {
	var arg2Copy []byte
	if arg2 != nil {
		arg2Copy = make([]byte, len(arg2))
		copy(arg2Copy, arg2)
	}
	fake.writeConfigBlockMutex.Lock()
	fake.writeConfigBlockArgsForCall = append(fake.writeConfigBlockArgsForCall, struct {
		arg1 *common.Block
		arg2 []byte
	}{arg1, arg2Copy})
	fake.recordInvocation("WriteConfigBlock", []interface{}{arg1, arg2Copy})
	fake.writeConfigBlockMutex.Unlock()
	if fake.WriteConfigBlockStub != nil {
		fake.WriteConfigBlockStub(arg1, arg2)
	}
}

func (fake *FakeConsenterSupport) WriteConfigBlockCallCount() int {
	fake.writeConfigBlockMutex.RLock()
	defer fake.writeConfigBlockMutex.RUnlock()
	return len(fake.writeConfigBlockArgsForCall)
}

func (fake *FakeConsenterSupport) WriteConfigBlockCalls(stub func(*common.Block, []byte)) {
	fake.writeConfigBlockMutex.Lock()
	defer fake.writeConfigBlockMutex.Unlock()
	fake.WriteConfigBlockStub = stub
}

func (fake *FakeConsenterSupport) WriteConfigBlockArgsForCall(i int) (*common.Block, []byte) {
	fake.writeConfigBlockMutex.RLock()
	defer fake.writeConfigBlockMutex.RUnlock()
	argsForCall := fake.writeConfigBlockArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *FakeConsenterSupport) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.appendMutex.RLock()
	defer fake.appendMutex.RUnlock()
	fake.blockMutex.RLock()
	defer fake.blockMutex.RUnlock()
	fake.blockCutterMutex.RLock()
	defer fake.blockCutterMutex.RUnlock()
	fake.chainIDMutex.RLock()
	defer fake.chainIDMutex.RUnlock()
	fake.channelConfigMutex.RLock()
	defer fake.channelConfigMutex.RUnlock()
	fake.classifyMsgMutex.RLock()
	defer fake.classifyMsgMutex.RUnlock()
	fake.createNextBlockMutex.RLock()
	defer fake.createNextBlockMutex.RUnlock()
	fake.heightMutex.RLock()
	defer fake.heightMutex.RUnlock()
	fake.processConfigMsgMutex.RLock()
	defer fake.processConfigMsgMutex.RUnlock()
	fake.processConfigUpdateMsgMutex.RLock()
	defer fake.processConfigUpdateMsgMutex.RUnlock()
	fake.processNormalMsgMutex.RLock()
	defer fake.processNormalMsgMutex.RUnlock()
	fake.sequenceMutex.RLock()
	defer fake.sequenceMutex.RUnlock()
	fake.serializeMutex.RLock()
	defer fake.serializeMutex.RUnlock()
	fake.sharedConfigMutex.RLock()
	defer fake.sharedConfigMutex.RUnlock()
	fake.signMutex.RLock()
	defer fake.signMutex.RUnlock()
	fake.verifyBlockSignatureMutex.RLock()
	defer fake.verifyBlockSignatureMutex.RUnlock()
	fake.writeBlockMutex.RLock()
	defer fake.writeBlockMutex.RUnlock()
	fake.writeConfigBlockMutex.RLock()
	defer fake.writeConfigBlockMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakeConsenterSupport) recordInvocation(key string, args []interface{}) {
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

var _ consensus.ConsenterSupport = new(FakeConsenterSupport)
