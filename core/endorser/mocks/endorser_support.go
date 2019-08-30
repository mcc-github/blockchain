/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package mocks

import (
	pb "github.com/mcc-github/blockchain-protos-go/peer"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/endorser"
	"github.com/mcc-github/blockchain/core/ledger"
	"github.com/stretchr/testify/mock"
)

type MockSupport struct {
	*mock.Mock
	*endorser.PluginEndorser
	IsSysCCAndNotInvokableExternalRv bool
	IsSysCCRv                        bool
	ExecuteCDSResp                   *pb.Response
	ExecuteCDSEvent                  *pb.ChaincodeEvent
	ExecuteCDSError                  error
	ExecuteResp                      *pb.Response
	ExecuteEvent                     *pb.ChaincodeEvent
	ExecuteError                     error
	ChaincodeDefinitionRv            ccprovider.ChaincodeDefinition
	ChaincodeDefinitionError         error
	GetTxSimulatorRv                 *MockTxSim
	GetTxSimulatorErr                error
	CheckInstantiationPolicyError    error
	GetTransactionByIDErr            error
	CheckACLErr                      error
	SysCCMap                         map[string]struct{}
	IsJavaRV                         bool
	IsJavaErr                        error
	GetApplicationConfigRv           channelconfig.Application
	GetApplicationConfigBoolRv       bool
	DeployedCCInfoProvider           ledger.DeployedChaincodeInfoProvider
}

func (s *MockSupport) Serialize() ([]byte, error) {
	args := s.Called()
	return args.Get(0).([]byte), args.Error(1)
}

func (s *MockSupport) NewQueryCreator(channel string) (endorser.QueryCreator, error) {
	panic("implement me")
}

func (s *MockSupport) Sign(message []byte) ([]byte, error) {
	args := s.Called(message)
	return args.Get(0).([]byte), args.Error(1)
}

func (s *MockSupport) ChannelState(channel string) (endorser.QueryCreator, error) {
	panic("implement me")
}

func (s *MockSupport) IsSysCCAndNotInvokableExternal(name string) bool {
	return s.IsSysCCAndNotInvokableExternalRv
}

func (s *MockSupport) GetTxSimulator(ledgername string, txid string) (ledger.TxSimulator, error) {
	if s.Mock == nil {
		return s.GetTxSimulatorRv, s.GetTxSimulatorErr
	}

	args := s.Called(ledgername, txid)
	return args.Get(0).(ledger.TxSimulator), args.Error(1)
}

func (s *MockSupport) GetHistoryQueryExecutor(ledgername string) (ledger.HistoryQueryExecutor, error) {
	return nil, nil
}

func (s *MockSupport) GetTransactionByID(chid, txID string) (*pb.ProcessedTransaction, error) {
	return nil, s.GetTransactionByIDErr
}

func (s *MockSupport) GetLedgerHeight(channelID string) (uint64, error) {
	args := s.Called(channelID)
	return args.Get(0).(uint64), args.Error(1)
}

func (s *MockSupport) IsSysCC(name string) bool {
	if s.SysCCMap != nil {
		_, in := s.SysCCMap[name]
		return in
	}
	return s.IsSysCCRv
}

func (s *MockSupport) ExecuteLegacyInit(txParams *ccprovider.TransactionParams, name, version string, spec *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error) {
	return s.ExecuteCDSResp, s.ExecuteCDSEvent, s.ExecuteCDSError
}

func (s *MockSupport) Execute(txParams *ccprovider.TransactionParams, name string, spec *pb.ChaincodeInput) (*pb.Response, *pb.ChaincodeEvent, error) {
	return s.ExecuteResp, s.ExecuteEvent, s.ExecuteError
}

func (s *MockSupport) GetChaincodeDeploymentSpecFS(cds *pb.ChaincodeDeploymentSpec) (*pb.ChaincodeDeploymentSpec, error) {
	return cds, nil
}

func (s *MockSupport) GetChaincodeDefinition(channelID, chaincodeName string, txsim ledger.QueryExecutor) (ccprovider.ChaincodeDefinition, error) {
	return s.ChaincodeDefinitionRv, s.ChaincodeDefinitionError
}

func (s *MockSupport) CheckACL(channelID string, signedProp *pb.SignedProposal) error {
	return s.CheckACLErr
}

func (s *MockSupport) CheckInstantiationPolicy(nameVersion string, cd ccprovider.ChaincodeDefinition) error {
	return s.CheckInstantiationPolicyError
}

func (s *MockSupport) GetApplicationConfig(cid string) (channelconfig.Application, bool) {
	return s.GetApplicationConfigRv, s.GetApplicationConfigBoolRv
}

func (s *MockSupport) GetDeployedCCInfoProvider() ledger.DeployedChaincodeInfoProvider {
	return s.DeployedCCInfoProvider
}
