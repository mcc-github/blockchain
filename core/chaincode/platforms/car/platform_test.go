/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package car_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/mcc-github/blockchain/common/util"
	"github.com/mcc-github/blockchain/core/chaincode/platforms"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/car"
	"github.com/mcc-github/blockchain/core/testutil"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

var _ = platforms.Platform(&car.Platform{})

func TestMain(m *testing.M) {
	testutil.SetupTestConfig()
	os.Exit(m.Run())
}

func TestCar_BuildImage(t *testing.T) {
	vm, err := NewVM()
	if err != nil {
		t.Errorf("Error getting VM: %s", err)
		return
	}

	chaincodePath := filepath.Join("testdata", "/org.mcc-github.chaincode.example02-0.1-SNAPSHOT.car")
	spec := &pb.ChaincodeSpec{
		Type: pb.ChaincodeSpec_CAR,
		ChaincodeId: &pb.ChaincodeID{
			Name: "cartest",
			Path: chaincodePath,
		},
		Input: &pb.ChaincodeInput{
			Args: util.ToChaincodeArgs("f"),
		},
	}
	if err := vm.BuildChaincodeContainer(spec); err != nil {
		t.Error(err)
	}
}
