

package chaincode

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"

	"github.com/mcc-github/blockchain/msp"
	msptesttools "github.com/mcc-github/blockchain/msp/mgmt/testtools"
	"github.com/mcc-github/blockchain/peer/common"
	pcommon "github.com/mcc-github/blockchain/protos/common"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

func TestMain(m *testing.M) {
	err := msptesttools.LoadMSPSetupForTesting()
	if err != nil {
		panic(fmt.Sprintf("Fatal error when reading MSP config: %s", err))
	}

	os.Exit(m.Run())
}

func newTempDir() string {
	tempDir, err := ioutil.TempDir("/tmp", "packagetest-")
	if err != nil {
		panic(err)
	}
	return tempDir
}

func mockCDSFactory(spec *pb.ChaincodeSpec) (*pb.ChaincodeDeploymentSpec, error) {
	return &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec, CodePackage: []byte("somecode")}, nil
}

func extractSignedCCDepSpec(env *pcommon.Envelope) (*pcommon.ChannelHeader, *pb.SignedChaincodeDeploymentSpec, error) {
	p := &pcommon.Payload{}
	err := proto.Unmarshal(env.Payload, p)
	if err != nil {
		return nil, nil, err
	}
	ch := &pcommon.ChannelHeader{}
	err = proto.Unmarshal(p.Header.ChannelHeader, ch)
	if err != nil {
		return nil, nil, err
	}

	sp := &pb.SignedChaincodeDeploymentSpec{}
	err = proto.Unmarshal(p.Data, sp)
	if err != nil {
		return nil, nil, err
	}

	return ch, sp, nil
}



func TestCDSPackage(t *testing.T) {
	pdir := newTempDir()
	defer os.RemoveAll(pdir)

	ccpackfile := pdir + "/ccpack.file"
	err := createSignedCDSPackage([]string{"-n", "somecc", "-p", "some/go/package", "-v", "0", ccpackfile}, false)
	if err != nil {
		t.Fatalf("Run chaincode package cmd error:%v", err)
	}

	b, err := ioutil.ReadFile(ccpackfile)
	if err != nil {
		t.Fatalf("package file %s not created", ccpackfile)
	}
	cds := &pb.ChaincodeDeploymentSpec{}
	err = proto.Unmarshal(b, cds)
	if err != nil {
		t.Fatalf("could not unmarshall package into CDS")
	}
}


func createSignedCDSPackage(args []string, sign bool) error {
	var signer msp.SigningIdentity
	var err error
	if sign {

		signer, err = common.GetDefaultSigner()
		if err != nil {
			return fmt.Errorf("Get default signer error: %v", err)
		}
	}

	mockCF := &ChaincodeCmdFactory{Signer: signer}

	cmd := packageCmd(mockCF, mockCDSFactory)
	addFlags(cmd)

	cmd.SetArgs(args)

	if err := cmd.Execute(); err != nil {
		return err
	}

	return nil
}



func TestSignedCDSPackage(t *testing.T) {
	pdir := newTempDir()
	defer os.RemoveAll(pdir)

	ccpackfile := pdir + "/ccpack.file"
	err := createSignedCDSPackage([]string{"-n", "somecc", "-p", "some/go/package", "-v", "0", "-s", ccpackfile}, false)
	if err != nil {
		t.Fatalf("could not create signed cds package %s", err)
	}

	b, err := ioutil.ReadFile(ccpackfile)
	if err != nil {
		t.Fatalf("package file %s not created", ccpackfile)
	}

	e := &pcommon.Envelope{}
	err = proto.Unmarshal(b, e)
	if err != nil {
		t.Fatalf("could not unmarshall envelope")
	}

	_, p, err := extractSignedCCDepSpec(e)
	if err != nil {
		t.Fatalf("could not extract signed dep spec")
	}

	if p.OwnerEndorsements != nil {
		t.Fatalf("expected no signtures but found endorsements")
	}
}



func TestSignedCDSPackageWithSignature(t *testing.T) {
	pdir := newTempDir()
	defer os.RemoveAll(pdir)

	ccpackfile := pdir + "/ccpack.file"
	err := createSignedCDSPackage([]string{"-n", "somecc", "-p", "some/go/package", "-v", "0", "-s", "-S", ccpackfile}, true)
	if err != nil {
		t.Fatalf("could not create signed cds package %s", err)
	}

	b, err := ioutil.ReadFile(ccpackfile)
	if err != nil {
		t.Fatalf("package file %s not created", ccpackfile)
	}
	e := &pcommon.Envelope{}
	err = proto.Unmarshal(b, e)
	if err != nil {
		t.Fatalf("could not unmarshall envelope")
	}

	_, p, err := extractSignedCCDepSpec(e)
	if err != nil {
		t.Fatalf("could not extract signed dep spec")
	}

	if p.OwnerEndorsements == nil {
		t.Fatalf("expected signtures and found nil")
	}
}

func TestNoOwnerToSign(t *testing.T) {
	pdir := newTempDir()
	defer os.RemoveAll(pdir)

	ccpackfile := pdir + "/ccpack.file"
	
	err := createSignedCDSPackage([]string{"-n", "somecc", "-p", "some/go/package", "-v", "0", "-s", "-S", ccpackfile}, false)

	if err == nil {
		t.Fatalf("Expected error with nil signer but succeeded")
	}
}

func TestInvalidPolicy(t *testing.T) {
	pdir := newTempDir()
	defer os.RemoveAll(pdir)

	ccpackfile := pdir + "/ccpack.file"
	err := createSignedCDSPackage([]string{"-n", "somecc", "-p", "some/go/package", "-v", "0", "-s", "-i", "AND('a bad policy')", ccpackfile}, false)

	if err == nil {
		t.Fatalf("Expected error with nil signer but succeeded")
	}
}
