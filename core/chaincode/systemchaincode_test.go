/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/mcc-github/blockchain/common/crypto/tlsgen"
	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/accesscontrol"
	"github.com/mcc-github/blockchain/core/chaincode/shim"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/dockercontroller"
	"github.com/mcc-github/blockchain/core/container/inproccontroller"
	"github.com/mcc-github/blockchain/core/peer"
	"github.com/mcc-github/blockchain/core/scc"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/spf13/viper"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type oldSysCCInfo struct {
	origSystemCC       []*scc.SystemChaincode
	origSysCCWhitelist map[string]string
}

func (osyscc *oldSysCCInfo) reset() {
	viper.Set("chaincode.system", osyscc.origSysCCWhitelist)
}

type SampleSysCC struct{}

func (t *SampleSysCC) Init(stub shim.ChaincodeStubInterface) pb.Response {
	return shim.Success(nil)
}

func (t *SampleSysCC) Invoke(stub shim.ChaincodeStubInterface) pb.Response {
	f, args := stub.GetFunctionAndParameters()

	switch f {
	case "putval":
		if len(args) != 2 {
			return shim.Error("need 2 args (key and a value)")
		}

		
		key := args[0]
		val := args[1]

		_, err := stub.GetState(key)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get val for " + key + "\"}"
			return shim.Error(jsonResp)
		}

		
		err = stub.PutState(key, []byte(val))
		if err != nil {
			return shim.Error(err.Error())
		}

		return shim.Success(nil)
	case "getval":
		var err error

		if len(args) != 1 {
			return shim.Error("Incorrect number of arguments. Expecting key to query")
		}

		key := args[0]

		
		valbytes, err := stub.GetState(key)
		if err != nil {
			jsonResp := "{\"Error\":\"Failed to get state for " + key + "\"}"
			return shim.Error(jsonResp)
		}

		if valbytes == nil {
			jsonResp := "{\"Error\":\"Nil val for " + key + "\"}"
			return shim.Error(jsonResp)
		}

		return shim.Success(valbytes)
	default:
		jsonResp := "{\"Error\":\"Unknown function " + f + "\"}"
		return shim.Error(jsonResp)
	}
}

func initSysCCTests() (*oldSysCCInfo, net.Listener, *ChaincodeSupport, error) {
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	viper.Set("peer.fileSystemPath", "/tmp/mcc-github/test/tmpdb")
	defer os.RemoveAll("/tmp/mcc-github/test/tmpdb")

	peer.MockInitialize()

	mspGetter := func(cid string) []string {
		return []string{"SampleOrg"}
	}

	peer.MockSetMSPIDGetter(mspGetter)

	
	
	
	peerAddress := "0.0.0.0:21726"
	lis, err := net.Listen("tcp", peerAddress)
	if err != nil {
		return nil, nil, nil, err
	}

	ca, _ := tlsgen.NewCA()
	certGenerator := accesscontrol.NewAuthenticator(ca)
	config := GlobalConfig()
	config.ExecuteTimeout = 5 * time.Second
	ipRegistry := inproccontroller.NewRegistry()
	sccp := &scc.Provider{Peer: peer.Default, PeerSupport: peer.DefaultSupport, Registrar: ipRegistry}
	chaincodeSupport := NewChaincodeSupport(
		config,
		peerAddress,
		false,
		ca.CertBytes(),
		certGenerator,
		&ccprovider.CCInfoFSImpl{},
		mockAclProvider,
		container.NewVMController(
			map[string]container.VMProvider{
				dockercontroller.ContainerType: dockercontroller.NewProvider("", ""),
				inproccontroller.ContainerType: ipRegistry,
			},
		),
		sccp,
	)
	pb.RegisterChaincodeSupportServer(grpcServer, chaincodeSupport)

	go grpcServer.Serve(lis)

	
	syscc := &scc.SystemChaincode{
		Enabled:   true,
		Name:      "sample_syscc",
		Path:      "github.com/mcc-github/blockchain/core/scc/samplesyscc",
		InitArgs:  [][]byte{},
		Chaincode: &SampleSysCC{},
	}

	sysccinfo := &oldSysCCInfo{origSysCCWhitelist: viper.GetStringMapString("chaincode.system")}

	
	viper.Set("chaincode.system", map[string]string{"sample_syscc": "true"})

	sccp.RegisterSysCC(syscc)

	
	return sysccinfo, lis, chaincodeSupport, nil
}

func deploySampleSysCC(t *testing.T, ctxt context.Context, chainID string, chaincodeSupport *ChaincodeSupport) error {
	ccp := &CCProviderImpl{cs: chaincodeSupport}
	chaincodeSupport.sccp.(*scc.Provider).DeploySysCCs(chainID, ccp)

	defer chaincodeSupport.sccp.(*scc.Provider).DeDeploySysCCs(chainID, ccp)

	url := "github.com/mcc-github/blockchain/core/scc/sample_syscc"

	sysCCVers := util.GetSysCCVersion()

	f := "putval"
	args := util.ToChaincodeArgs(f, "greeting", "hey there")

	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "sample_syscc", Path: url, Version: sysCCVers}, Input: &pb.ChaincodeInput{Args: args}}
	
	var nextBlockNumber uint64 = 1
	_, _, _, err := invokeWithVersion(ctxt, chainID, sysCCVers, spec, nextBlockNumber, nil, chaincodeSupport)
	nextBlockNumber++

	cccid := ccprovider.NewCCContext(chainID, "sample_syscc", sysCCVers, "", true, nil, nil)
	cdsforStop := &pb.ChaincodeDeploymentSpec{ExecEnv: 1, ChaincodeSpec: spec}
	if err != nil {
		chaincodeSupport.Stop(ctxt, cccid, cdsforStop)
		t.Logf("Error invoking sample_syscc: %s", err)
		return err
	}

	f = "getval"
	args = util.ToChaincodeArgs(f, "greeting")
	spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: &pb.ChaincodeID{Name: "sample_syscc", Path: url, Version: sysCCVers}, Input: &pb.ChaincodeInput{Args: args}}
	_, _, _, err = invokeWithVersion(ctxt, chainID, sysCCVers, spec, nextBlockNumber, nil, chaincodeSupport)
	if err != nil {
		chaincodeSupport.Stop(ctxt, cccid, cdsforStop)
		t.Logf("Error invoking sample_syscc: %s", err)
		return err
	}

	chaincodeSupport.Stop(ctxt, cccid, cdsforStop)

	return nil
}


func TestExecuteDeploySysChaincode(t *testing.T) {
	testForSkip(t)
	sysccinfo, lis, chaincodeSupport, err := initSysCCTests()
	if err != nil {
		t.Fail()
		return
	}

	defer func() {
		sysccinfo.reset()
	}()

	chainID := util.GetTestChainID()

	if err = peer.MockCreateChain(chainID); err != nil {
		closeListenerAndSleep(lis)
		return
	}

	var ctxt = context.Background()

	err = deploySampleSysCC(t, ctxt, chainID, chaincodeSupport)
	if err != nil {
		closeListenerAndSleep(lis)
		t.Fail()
		return
	}

	closeListenerAndSleep(lis)
}


func TestMultichains(t *testing.T) {
	testForSkip(t)
	sysccinfo, lis, chaincodeSupport, err := initSysCCTests()
	if err != nil {
		t.Fail()
		return
	}

	defer func() {
		sysccinfo.reset()
	}()

	chainID := "chain1"

	if err = peer.MockCreateChain(chainID); err != nil {
		closeListenerAndSleep(lis)
		return
	}

	var ctxt = context.Background()

	err = deploySampleSysCC(t, ctxt, chainID, chaincodeSupport)
	if err != nil {
		closeListenerAndSleep(lis)
		t.Fail()
		return
	}

	chainID = "chain2"

	if err = peer.MockCreateChain(chainID); err != nil {
		closeListenerAndSleep(lis)
		return
	}

	err = deploySampleSysCC(t, ctxt, chainID, chaincodeSupport)
	if err != nil {
		closeListenerAndSleep(lis)
		t.Fail()
		return
	}

	closeListenerAndSleep(lis)
}
