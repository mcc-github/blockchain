/*
Copyright IBM Corp. 2016 All Rights Reserved.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

		 http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package chaincode

import (
	"fmt"
	"sync"
	"testing"

	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	pb "github.com/mcc-github/blockchain/protos/peer"

	"golang.org/x/net/context"
)



func TestExecuteConcurrentInvokes(t *testing.T) {
	
	

	t.Skip()
	chainID := util.GetTestChainID()

	_, chaincodeSupport, cleanup, err := initPeer(chainID)
	if err != nil {
		t.Fail()
		t.Logf("Error creating peer: %s", err)
	}

	defer cleanup()

	var ctxt = context.Background()

	url := "github.com/mcc-github/blockchain/examples/ccchecker/chaincodes/newkeyperinvoke"

	chaincodeID := &pb.ChaincodeID{Name: "nkpi", Path: url, Version: "0"}

	args := util.ToChaincodeArgs("init", "")

	spec := &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: args}}

	cccid := ccprovider.NewCCContext(chainID, "nkpi", "0", "", false, nil, nil)

	defer chaincodeSupport.Stop(ctxt, cccid, &pb.ChaincodeDeploymentSpec{ChaincodeSpec: spec})

	var nextBlockNumber uint64
	_, err = deploy(ctxt, cccid, spec, nextBlockNumber, chaincodeSupport)
	nextBlockNumber++
	if err != nil {
		t.Fail()
		t.Logf("Error initializing chaincode %s(%s)", chaincodeID, err)
		return
	}

	var wg sync.WaitGroup

	
	numTrans := 100

	results := make([][]byte, numTrans)
	errs := make([]error, numTrans)

	e := func(inv bool, qnum int) {
		defer wg.Done()

		newkey := fmt.Sprintf("%d", qnum)

		var args [][]byte
		if inv {
			args = util.ToChaincodeArgs("put", newkey, newkey)
		} else {
			args = util.ToChaincodeArgs("get", newkey)
		}

		spec = &pb.ChaincodeSpec{Type: 1, ChaincodeId: chaincodeID, Input: &pb.ChaincodeInput{Args: args}}

		
		_, _, results[qnum], err = invoke(context.Background(), chainID, spec, nextBlockNumber, nil, chaincodeSupport)

		if err != nil {
			errs[qnum] = fmt.Errorf("Error executing <%s>: %s", chaincodeID.Name, err)
			return
		}
	}

	wg.Add(numTrans)

	
	for i := 0; i < numTrans; i++ {
		go e(true, i)
	}

	wg.Wait()

	for i := 0; i < numTrans; i++ {
		if errs[i] != nil {
			t.Fail()
			t.Logf("Error invoking chaincode iter %d %s(%s)", i, chaincodeID.Name, errs[i])
		}
		if results[i] == nil || string(results[i]) != "OK" {
			t.Fail()
			t.Logf("Error concurrent invoke %d %s", i, chaincodeID.Name)
			return
		}
	}

	wg.Add(numTrans)

	
	for i := 0; i < numTrans; i++ {
		go e(false, i)
	}

	wg.Wait()

	for i := 0; i < numTrans; i++ {
		if errs[i] != nil {
			t.Fail()
			t.Logf("Error querying chaincode iter %d %s(%s)", i, chaincodeID.Name, errs[i])
			return
		}
		if results[i] == nil || string(results[i]) != fmt.Sprintf("%d", i) {
			t.Fail()
			if results[i] == nil {
				t.Logf("Error concurrent query %d(%s)", i, chaincodeID.Name)
			} else {
				t.Logf("Error concurrent query %d(%s, %s, %v)", i, chaincodeID.Name, string(results[i]), results[i])
			}
			return
		}
	}
}
