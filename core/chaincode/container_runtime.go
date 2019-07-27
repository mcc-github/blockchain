/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"bytes"
	"fmt"
	"sort"
	"strings"

	"github.com/mcc-github/blockchain/core/chaincode/accesscontrol"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/ccintf"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
)


type Processor interface {
	Process(vmtype string, req container.VMCReq) error
}


type CertGenerator interface {
	
	
	Generate(ccName string) (*accesscontrol.CertAndPrivKeyPair, error)
}


type ContainerRuntime struct {
	CertGenerator CertGenerator
	Processor     Processor
	CACert        []byte
	CommonEnv     []string
	PeerAddress   string
}


func (c *ContainerRuntime) Start(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) error {
	packageID := ccci.PackageID.String()

	lc, err := c.LaunchConfig(packageID, ccci.Type)
	if err != nil {
		return err
	}

	
	bcr := container.BuildReq{
		CCID:        ccintf.New(ccci.PackageID),
		Type:        ccci.Type,
		Name:        ccci.Name,
		Version:     ccci.Version,
		Path:        ccci.Path,
		CodePackage: codePackage,
	}
	if err := c.Processor.Process(ccci.ContainerType, bcr); err != nil {
		return errors.WithMessage(err, "error building image")
	}

	chaincodeLogger.Debugf("start container: %s", packageID)
	chaincodeLogger.Debugf("start container with args: %s", strings.Join(lc.Args, " "))
	chaincodeLogger.Debugf("start container with env:\n\t%s", strings.Join(lc.Envs, "\n\t"))

	scr := container.StartContainerReq{
		Args:          lc.Args,
		Env:           lc.Envs,
		FilesToUpload: lc.Files,
		CCID:          ccintf.New(ccci.PackageID),
	}

	if err := c.Processor.Process(ccci.ContainerType, scr); err != nil {
		return errors.WithMessage(err, "error starting container")
	}

	return nil
}


func (c *ContainerRuntime) Stop(ccci *ccprovider.ChaincodeContainerInfo) error {
	scr := container.StopContainerReq{
		CCID:       ccintf.New(ccci.PackageID),
		Timeout:    0,
		Dontremove: false,
	}

	if err := c.Processor.Process(ccci.ContainerType, scr); err != nil {
		return errors.WithMessage(err, "error stopping container")
	}

	return nil
}


func (c *ContainerRuntime) Wait(ccci *ccprovider.ChaincodeContainerInfo) (int, error) {
	type result struct {
		exitCode int
		err      error
	}

	resultCh := make(chan result, 1)
	wcr := container.WaitContainerReq{
		CCID: ccintf.New(ccci.PackageID),
		Exited: func(exitCode int, err error) {
			resultCh <- result{exitCode: exitCode, err: err}
			close(resultCh)
		},
	}

	if err := c.Processor.Process(ccci.ContainerType, wcr); err != nil {
		return -1, err
	}
	r := <-resultCh

	return r.exitCode, r.err
}

const (
	
	TLSClientKeyPath      string = "/etc/mcc-github/blockchain/client.key"
	TLSClientCertPath     string = "/etc/mcc-github/blockchain/client.crt"
	TLSClientRootCertPath string = "/etc/mcc-github/blockchain/peer.crt"
)

func (c *ContainerRuntime) getTLSFiles(keyPair *accesscontrol.CertAndPrivKeyPair) map[string][]byte {
	if keyPair == nil {
		return nil
	}

	return map[string][]byte{
		TLSClientKeyPath:      []byte(keyPair.Key),
		TLSClientCertPath:     []byte(keyPair.Cert),
		TLSClientRootCertPath: c.CACert,
	}
}


type LaunchConfig struct {
	Args  []string
	Envs  []string
	Files map[string][]byte
}


func (c *ContainerRuntime) LaunchConfig(packageID string, ccType string) (*LaunchConfig, error) {
	var lc LaunchConfig

	
	
	
	
	
	
	lc.Envs = append(c.CommonEnv, "CORE_CHAINCODE_ID_NAME="+packageID)

	
	switch ccType {
	case pb.ChaincodeSpec_GOLANG.String(), pb.ChaincodeSpec_CAR.String():
		lc.Args = []string{"chaincode", fmt.Sprintf("-peer.address=%s", c.PeerAddress)}
	case pb.ChaincodeSpec_JAVA.String():
		lc.Args = []string{"/root/chaincode-java/start", "--peerAddress", c.PeerAddress}
	case pb.ChaincodeSpec_NODE.String():
		lc.Args = []string{"/bin/sh", "-c", fmt.Sprintf("cd /usr/local/src; npm start -- --peer.address %s", c.PeerAddress)}
	default:
		return nil, errors.Errorf("unknown chaincodeType: %s", ccType)
	}

	
	if c.CertGenerator != nil {
		certKeyPair, err := c.CertGenerator.Generate(packageID)
		if err != nil {
			return nil, errors.WithMessagef(err, "failed to generate TLS certificates for %s", packageID)
		}
		lc.Files = c.getTLSFiles(certKeyPair)
		if lc.Files == nil {
			return nil, errors.Errorf("failed to acquire TLS certificates for %s", packageID)
		}

		lc.Envs = append(lc.Envs, "CORE_PEER_TLS_ENABLED=true")
		lc.Envs = append(lc.Envs, fmt.Sprintf("CORE_TLS_CLIENT_KEY_PATH=%s", TLSClientKeyPath))
		lc.Envs = append(lc.Envs, fmt.Sprintf("CORE_TLS_CLIENT_CERT_PATH=%s", TLSClientCertPath))
		lc.Envs = append(lc.Envs, fmt.Sprintf("CORE_PEER_TLS_ROOTCERT_FILE=%s", TLSClientRootCertPath))
	} else {
		lc.Envs = append(lc.Envs, "CORE_PEER_TLS_ENABLED=false")
	}

	chaincodeLogger.Debugf("launchConfig: %s", lc.String())

	return &lc, nil
}

func (lc *LaunchConfig) String() string {
	buf := &bytes.Buffer{}
	if len(lc.Args) > 0 {
		fmt.Fprintf(buf, "executable:%q,", lc.Args[0])
	}

	fileNames := []string{}
	for k := range lc.Files {
		fileNames = append(fileNames, k)
	}
	sort.Strings(fileNames)

	fmt.Fprintf(buf, "Args:[%s],", strings.Join(lc.Args, ","))
	fmt.Fprintf(buf, "Envs:[%s],", strings.Join(lc.Envs, ","))
	fmt.Fprintf(buf, "Files:%v", fileNames)
	return buf.String()
}
