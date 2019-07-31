/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package externalbuilders

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	"github.com/mcc-github/blockchain/core/container"
	"github.com/mcc-github/blockchain/core/container/ccintf"

	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("chaincode.externalbuilders")

type Instance struct {
	BuildContext *BuildContext
	Builder      *Builder
}

func (i *Instance) Start(peerConnection *ccintf.PeerConnection) error {
	return i.Builder.Launch(i.BuildContext, peerConnection)
}

func (i *Instance) Stop() error {
	
	return errors.Errorf("stop is not implemented for external builders yet")
}

func (i *Instance) Wait() (int, error) {
	
	
	select {}
}

type Detector struct {
	Builders []string
}

func (d *Detector) Detect(buildContext *BuildContext) *Builder {
	for _, builderLocation := range d.Builders {
		builder := &Builder{
			Location: builderLocation,
		}
		if builder.Detect(buildContext) {
			return builder
		}
	}

	return nil
}

func (d *Detector) Build(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) (container.Instance, error) {
	if len(d.Builders) == 0 {
		
		
		
		return nil, errors.Errorf("no builders defined")
	}

	buildContext, err := NewBuildContext(ccci, bytes.NewBuffer(codePackage))
	if err != nil {
		return nil, errors.WithMessage(err, "could not create build context")
	}

	builder := d.Detect(buildContext)
	if builder == nil {
		buildContext.Cleanup()
		return nil, errors.Errorf("no builder found")
	}

	if err := builder.Build(buildContext); err != nil {
		buildContext.Cleanup()
		return nil, errors.WithMessage(err, "external builder failed to build")
	}

	return &Instance{
		BuildContext: buildContext,
		Builder:      builder,
	}, nil
}

type BuildContext struct {
	CCCI       *ccprovider.ChaincodeContainerInfo
	ScratchDir string
	SourceDir  string
	OutputDir  string
	LaunchDir  string
}

func NewBuildContext(ccci *ccprovider.ChaincodeContainerInfo, codePackage io.Reader) (*BuildContext, error) {
	
	scratchDir, err := ioutil.TempDir("", "blockchain-"+strings.ReplaceAll(string(ccci.PackageID), string(os.PathListSeparator), "-"))
	if err != nil {
		return nil, errors.WithMessage(err, "could not create temp dir")
	}

	sourceDir := filepath.Join(scratchDir, "src")
	err = os.Mkdir(sourceDir, 0700)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not create source dir")
	}

	outputDir := filepath.Join(scratchDir, "bld")
	err = os.Mkdir(outputDir, 0700)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not create source dir")
	}

	launchDir := filepath.Join(scratchDir, "run")
	err = os.MkdirAll(launchDir, 0700)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not create source dir")
	}

	err = Untar(codePackage, sourceDir)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not untar source package")
	}

	return &BuildContext{
		ScratchDir: scratchDir,
		SourceDir:  sourceDir,
		OutputDir:  outputDir,
		LaunchDir:  launchDir,
		CCCI:       ccci,
	}, nil
}

func (bc *BuildContext) Cleanup() {
	os.RemoveAll(bc.ScratchDir)
}

type Builder struct {
	Location string
}

func (b *Builder) Detect(buildContext *BuildContext) bool {
	detect := filepath.Join(b.Location, "bin", "detect")
	cmd := exec.Cmd{
		Path: detect,
		Args: []string{
			"detect", 
			"--package-id", string(buildContext.CCCI.PackageID),
			
			
			
			"--path", buildContext.CCCI.Path,
			"--type", buildContext.CCCI.Type,
			"--source", buildContext.SourceDir,
		},

		
		
	}

	err := cmd.Run()
	if err != nil {
		logger.Debugf("Detection for builder '%s' failed: %s", b.Name(), err)
		
		
		return false
	}

	return true
}

func (b *Builder) Build(buildContext *BuildContext) error {
	build := filepath.Join(b.Location, "bin", "build")
	cmd := exec.Cmd{
		Path: build,
		Args: []string{
			"build", 
			"--package-id", string(buildContext.CCCI.PackageID),
			"--path", buildContext.CCCI.Path,
			"--type", buildContext.CCCI.Type,
			"--source", buildContext.SourceDir,
			"--output", buildContext.OutputDir,
		},
	}

	err := cmd.Run()
	if err != nil {
		return errors.Errorf("builder '%s' failed: %s", b.Name(), err)
	}

	return nil
}


type LaunchConfig struct {
	PeerAddress string `json:"PeerAddress"`
	ClientCert  []byte `json:"ClientCert"`
	ClientKey   []byte `json:"ClientKey"`
	RootCert    []byte `json:"RootCert"`
}

func (b *Builder) Launch(buildContext *BuildContext, peerConnection *ccintf.PeerConnection) error {
	lc := &LaunchConfig{
		PeerAddress: peerConnection.Address,
	}

	if peerConnection.TLSConfig != nil {
		lc.ClientCert = peerConnection.TLSConfig.ClientCert
		lc.ClientKey = peerConnection.TLSConfig.ClientKey
		lc.RootCert = peerConnection.TLSConfig.RootCert
	}

	marshaledLC, err := json.Marshal(lc)
	if err != nil {
		return errors.WithMessage(err, "could not marshal launch config")
	}

	if err := ioutil.WriteFile(filepath.Join(buildContext.LaunchDir, "chaincode.json"), marshaledLC, 0600); err != nil {
		return errors.WithMessage(err, "could not write root cert")
	}

	launch := filepath.Join(b.Location, "bin", "launch")
	cmd := exec.Cmd{
		Path: launch,
		Args: []string{
			"launch", 
			"--package-id", string(buildContext.CCCI.PackageID),
			"--path", buildContext.CCCI.Path,
			"--type", buildContext.CCCI.Type,
			"--source", buildContext.SourceDir,
			"--output", buildContext.OutputDir,
			"--artifacts", buildContext.LaunchDir,
		},
	}

	err = cmd.Run()
	if err != nil {
		return errors.Errorf("builder '%s' failed: %s", b.Name(), err)
	}

	return nil
}



func (b *Builder) Name() string {
	return filepath.Base(b.Location)
}
