/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package common

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/mcc-github/blockchain/bccsp/factory"
	"github.com/mcc-github/blockchain/common/channelconfig"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/comm"
	"github.com/mcc-github/blockchain/core/config"
	"github.com/mcc-github/blockchain/core/scc/cscc"
	"github.com/mcc-github/blockchain/msp"
	mspmgmt "github.com/mcc-github/blockchain/msp/mgmt"
	pcommon "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/mcc-github/blockchain/protoutil"
	"github.com/mitchellh/mapstructure"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)


const UndefinedParamValue = ""
const CmdRoot = "core"

var mainLogger = flogging.MustGetLogger("main")
var logOutput = os.Stderr

var (
	defaultConnTimeout = 3 * time.Second
	
	

	
	
	
	GetEndorserClientFnc func(address, tlsRootCertFile string) (pb.EndorserClient, error)

	
	
	
	GetPeerDeliverClientFnc func(address, tlsRootCertFile string) (pb.DeliverClient, error)

	
	
	
	GetDeliverClientFnc func(address, tlsRootCertFile string) (pb.Deliver_DeliverClient, error)

	
	
	GetDefaultSignerFnc func() (msp.SigningIdentity, error)

	
	
	GetBroadcastClientFnc func() (BroadcastClient, error)

	
	
	GetOrdererEndpointOfChainFnc func(chainID string, signer Signer,
		endorserClient pb.EndorserClient) ([]string, error)

	
	GetCertificateFnc func() (tls.Certificate, error)
)

type CommonClient struct {
	*comm.GRPCClient
	Address string
	sn      string
}

func init() {
	GetEndorserClientFnc = GetEndorserClient
	GetDefaultSignerFnc = GetDefaultSigner
	GetBroadcastClientFnc = GetBroadcastClient
	GetOrdererEndpointOfChainFnc = GetOrdererEndpointOfChain
	GetDeliverClientFnc = GetDeliverClient
	GetPeerDeliverClientFnc = GetPeerDeliverClient
	GetCertificateFnc = GetCertificate
}


func InitConfig(cmdRoot string) error {

	err := config.InitViper(nil, cmdRoot)
	if err != nil {
		return err
	}

	err = viper.ReadInConfig() 
	if err != nil {            
		
		
		if strings.Contains(fmt.Sprint(err), "Unsupported Config Type") {
			return errors.New(fmt.Sprintf("Could not find config file. "+
				"Please make sure that FABRIC_CFG_PATH is set to a path "+
				"which contains %s.yaml", cmdRoot))
		} else {
			return errors.WithMessagef(err, "error when reading %s config file", cmdRoot)
		}
	}

	return nil
}


func InitCrypto(mspMgrConfigDir, localMSPID, localMSPType string) error {
	
	fi, err := os.Stat(mspMgrConfigDir)
	if os.IsNotExist(err) || !fi.IsDir() {
		
		return errors.Errorf("cannot init crypto, folder \"%s\" does not exist", mspMgrConfigDir)
	}
	
	if localMSPID == "" {
		return errors.New("the local MSP must have an ID")
	}

	
	SetBCCSPKeystorePath()
	bccspConfig := factory.GetDefaultOpts()
	if config := viper.Get("peer.BCCSP"); config != nil {
		err = mapstructure.Decode(config, bccspConfig)
		if err != nil {
			return errors.WithMessage(err, "could not decode peer BCCSP configuration")
		}
	}

	err = mspmgmt.LoadLocalMspWithType(mspMgrConfigDir, bccspConfig, localMSPID, localMSPType)
	if err != nil {
		return errors.WithMessagef(err, "error when setting up MSP of type %s from directory %s", localMSPType, mspMgrConfigDir)
	}

	return nil
}



func SetBCCSPKeystorePath() {
	viper.Set("peer.BCCSP.SW.FileKeyStore.KeyStore",
		config.GetPath("peer.BCCSP.SW.FileKeyStore.KeyStore"))
}


func GetDefaultSigner() (msp.SigningIdentity, error) {
	signer, err := mspmgmt.GetLocalMSP().GetDefaultSigningIdentity()
	if err != nil {
		return nil, errors.WithMessage(err, "error obtaining the default signing identity")
	}

	return signer, err
}


type Signer interface {
	Sign(msg []byte) ([]byte, error)
	Serialize() ([]byte, error)
}


func GetOrdererEndpointOfChain(chainID string, signer Signer, endorserClient pb.EndorserClient) ([]string, error) {
	
	invocation := &pb.ChaincodeInvocationSpec{
		ChaincodeSpec: &pb.ChaincodeSpec{
			Type:        pb.ChaincodeSpec_Type(pb.ChaincodeSpec_Type_value["GOLANG"]),
			ChaincodeId: &pb.ChaincodeID{Name: "cscc"},
			Input:       &pb.ChaincodeInput{Args: [][]byte{[]byte(cscc.GetConfigBlock), []byte(chainID)}},
		},
	}

	creator, err := signer.Serialize()
	if err != nil {
		return nil, errors.WithMessage(err, "error serializing identity for signer")
	}

	prop, _, err := protoutil.CreateProposalFromCIS(pcommon.HeaderType_CONFIG, "", invocation, creator)
	if err != nil {
		return nil, errors.WithMessage(err, "error creating GetConfigBlock proposal")
	}

	signedProp, err := protoutil.GetSignedProposal(prop, signer)
	if err != nil {
		return nil, errors.WithMessage(err, "error creating signed GetConfigBlock proposal")
	}

	proposalResp, err := endorserClient.ProcessProposal(context.Background(), signedProp)
	if err != nil {
		return nil, errors.WithMessage(err, "error endorsing GetConfigBlock")
	}

	if proposalResp == nil {
		return nil, errors.WithMessage(err, "error nil proposal response")
	}

	if proposalResp.Response.Status != 0 && proposalResp.Response.Status != 200 {
		return nil, errors.Errorf("error bad proposal response %d: %s", proposalResp.Response.Status, proposalResp.Response.Message)
	}

	
	block, err := protoutil.UnmarshalBlock(proposalResp.Response.Payload)
	if err != nil {
		return nil, errors.WithMessage(err, "error unmarshaling config block")
	}

	envelopeConfig, err := protoutil.ExtractEnvelope(block, 0)
	if err != nil {
		return nil, errors.WithMessage(err, "error extracting config block envelope")
	}
	bundle, err := channelconfig.NewBundleFromEnvelope(envelopeConfig, factory.GetDefault())
	if err != nil {
		return nil, errors.WithMessage(err, "error loading config block")
	}

	return bundle.ChannelConfig().OrdererAddresses(), nil
}


func CheckLogLevel(level string) error {
	if !flogging.IsValidLevel(level) {
		return errors.Errorf("invalid log level provided - %s", level)
	}
	return nil
}

func configFromEnv(prefix string) (address, override string, clientConfig comm.ClientConfig, err error) {
	address = viper.GetString(prefix + ".address")
	override = viper.GetString(prefix + ".tls.serverhostoverride")
	clientConfig = comm.ClientConfig{}
	connTimeout := viper.GetDuration(prefix + ".client.connTimeout")
	if connTimeout == time.Duration(0) {
		connTimeout = defaultConnTimeout
	}
	clientConfig.Timeout = connTimeout
	secOpts := comm.SecureOptions{
		UseTLS:            viper.GetBool(prefix + ".tls.enabled"),
		RequireClientCert: viper.GetBool(prefix + ".tls.clientAuthRequired")}
	if secOpts.UseTLS {
		caPEM, res := ioutil.ReadFile(config.GetPath(prefix + ".tls.rootcert.file"))
		if res != nil {
			err = errors.WithMessage(res,
				fmt.Sprintf("unable to load %s.tls.rootcert.file", prefix))
			return
		}
		secOpts.ServerRootCAs = [][]byte{caPEM}
	}
	if secOpts.RequireClientCert {
		keyPEM, res := ioutil.ReadFile(config.GetPath(prefix + ".tls.clientKey.file"))
		if res != nil {
			err = errors.WithMessage(res,
				fmt.Sprintf("unable to load %s.tls.clientKey.file", prefix))
			return
		}
		secOpts.Key = keyPEM
		certPEM, res := ioutil.ReadFile(config.GetPath(prefix + ".tls.clientCert.file"))
		if res != nil {
			err = errors.WithMessage(res,
				fmt.Sprintf("unable to load %s.tls.clientCert.file", prefix))
			return
		}
		secOpts.Certificate = certPEM
	}
	clientConfig.SecOpts = secOpts
	return
}

func InitCmd(cmd *cobra.Command, args []string) {
	err := InitConfig(CmdRoot)
	if err != nil { 
		mainLogger.Errorf("Fatal error when initializing %s config : %s", CmdRoot, err)
		os.Exit(1)
	}

	
	
	var loggingLevel string
	if viper.GetString("logging_level") != "" {
		loggingLevel = viper.GetString("logging_level")
	} else {
		loggingLevel = viper.GetString("logging.level")
	}
	if loggingLevel != "" {
		mainLogger.Warning("CORE_LOGGING_LEVEL is no longer supported, please use the FABRIC_LOGGING_SPEC environment variable")
	}

	loggingSpec := os.Getenv("FABRIC_LOGGING_SPEC")
	loggingFormat := os.Getenv("FABRIC_LOGGING_FORMAT")

	flogging.Init(flogging.Config{
		Format:  loggingFormat,
		Writer:  logOutput,
		LogSpec: loggingSpec,
	})

	
	var mspMgrConfigDir = config.GetPath("peer.mspConfigPath")
	var mspID = viper.GetString("peer.localMspId")
	var mspType = viper.GetString("peer.localMspType")
	if mspType == "" {
		mspType = msp.ProviderTypeToString(msp.FABRIC)
	}
	err = InitCrypto(mspMgrConfigDir, mspID, mspType)
	if err != nil { 
		mainLogger.Errorf("Cannot run peer because %s", err.Error())
		os.Exit(1)
	}

	runtime.GOMAXPROCS(viper.GetInt("peer.gomaxprocs"))
}
