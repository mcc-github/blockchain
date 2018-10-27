

package chaincode

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/mcc-github/blockchain/peer/common"
	pcommon "github.com/mcc-github/blockchain/protos/common"
)


func signExistingPackage(env *pcommon.Envelope, infile, outfile string) error {
	signer, err := common.GetDefaultSigner()
	if err != nil {
		return fmt.Errorf("Get default signer error: %v", err)
	}

	mockCF := &ChaincodeCmdFactory{Signer: signer}

	cmd := signpackageCmd(mockCF)
	addFlags(cmd)

	cmd.SetArgs([]string{infile, outfile})

	if err := cmd.Execute(); err != nil {
		return err
	}

	return nil
}


func TestSignExistingPackage(t *testing.T) {
	pdir := newTempDir()
	defer os.RemoveAll(pdir)

	ccpackfile := pdir + "/ccpack.file"
	err := createSignedCDSPackage([]string{"-n", "somecc", "-p", "some/go/package", "-v", "0", "-s", "-S", ccpackfile}, true)
	if err != nil {
		t.Fatalf("error creating signed :%v", err)
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

	signedfile := pdir + "/signed.file"
	err = signExistingPackage(e, ccpackfile, signedfile)
	if err != nil {
		t.Fatalf("could not sign envelope")
	}

	b, err = ioutil.ReadFile(signedfile)
	if err != nil {
		t.Fatalf("signed package file %s not created", signedfile)
	}

	e = &pcommon.Envelope{}
	err = proto.Unmarshal(b, e)
	if err != nil {
		t.Fatalf("could not unmarshall signed envelope")
	}

	_, p, err := extractSignedCCDepSpec(e)
	if err != nil {
		t.Fatalf("could not extract signed dep spec")
	}

	if p.OwnerEndorsements == nil {
		t.Fatalf("expected endorsements")
	}

	if len(p.OwnerEndorsements) != 2 {
		t.Fatalf("expected 2 endorserments but found %d", len(p.OwnerEndorsements))
	}
}


func TestFailSignUnsignedPackage(t *testing.T) {
	pdir := newTempDir()
	defer os.RemoveAll(pdir)

	ccpackfile := pdir + "/ccpack.file"
	
	err := createSignedCDSPackage([]string{"-n", "somecc", "-p", "some/go/package", "-v", "0", "-s", ccpackfile}, true)
	if err != nil {
		t.Fatalf("error creating signed :%v", err)
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

	signedfile := pdir + "/signed.file"
	err = signExistingPackage(e, ccpackfile, signedfile)
	if err == nil {
		t.Fatalf("expected signing a package that's not originally signed to fail")
	}
}
