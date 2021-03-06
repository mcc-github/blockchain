/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package externalbuilders

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	"github.com/mcc-github/blockchain/core/container/ccintf"

	"github.com/pkg/errors"
)

var logger = flogging.MustGetLogger("chaincode.externalbuilders")

const MetadataFile = "metadata.json"

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
			Logger:   logger.Named(filepath.Base(builderLocation)),
		}
		if builder.Detect(buildContext) {
			return builder
		}
	}

	return nil
}

func (d *Detector) Build(ccid string, md *persistence.ChaincodePackageMetadata, codeStream io.Reader) (*Instance, error) {
	if len(d.Builders) == 0 {
		
		
		
		return nil, errors.Errorf("no builders defined")
	}

	buildContext, err := NewBuildContext(string(ccid), md, codeStream)
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
	CCID        string
	Metadata    *persistence.ChaincodePackageMetadata
	ScratchDir  string
	SourceDir   string
	MetadataDir string
	OutputDir   string
	LaunchDir   string
}

var pkgIDreg = regexp.MustCompile("[^a-zA-Z0-9]+")

func NewBuildContext(ccid string, md *persistence.ChaincodePackageMetadata, codePackage io.Reader) (*BuildContext, error) {
	scratchDir, err := ioutil.TempDir("", "blockchain-"+pkgIDreg.ReplaceAllString(ccid, "-"))
	if err != nil {
		return nil, errors.WithMessage(err, "could not create temp dir")
	}

	sourceDir := filepath.Join(scratchDir, "src")
	err = os.Mkdir(sourceDir, 0700)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not create source dir")
	}

	metadataDir := filepath.Join(scratchDir, "metadata")
	err = os.Mkdir(metadataDir, 0700)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not create metadata dir")
	}

	outputDir := filepath.Join(scratchDir, "bld")
	err = os.Mkdir(outputDir, 0700)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not create build dir")
	}

	launchDir := filepath.Join(scratchDir, "run")
	err = os.MkdirAll(launchDir, 0700)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not create launch dir")
	}

	err = Untar(codePackage, sourceDir)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not untar source package")
	}

	err = writeMetadataFile(ccid, md, metadataDir)
	if err != nil {
		os.RemoveAll(scratchDir)
		return nil, errors.WithMessage(err, "could not write metadata file")
	}

	return &BuildContext{
		ScratchDir:  scratchDir,
		SourceDir:   sourceDir,
		MetadataDir: metadataDir,
		OutputDir:   outputDir,
		LaunchDir:   launchDir,
		Metadata:    md,
		CCID:        ccid,
	}, nil
}

type buildMetadata struct {
	Path      string `json:"path"`
	Type      string `json:"type"`
	PackageID string `json:"package_id"`
}

func writeMetadataFile(ccid string, md *persistence.ChaincodePackageMetadata, dst string) error {
	buildMetadata := &buildMetadata{
		Path:      md.Path,
		Type:      md.Type,
		PackageID: ccid,
	}
	mdBytes, err := json.Marshal(buildMetadata)
	if err != nil {
		return errors.Wrap(err, "failed to marshal build metadata into JSON")
	}

	return ioutil.WriteFile(filepath.Join(dst, MetadataFile), mdBytes, 0700)
}

func (bc *BuildContext) Cleanup() {
	os.RemoveAll(bc.ScratchDir)
}

type Builder struct {
	Location string
	Logger   *flogging.FabricLogger
}

func (b *Builder) Detect(buildContext *BuildContext) bool {
	detect := filepath.Join(b.Location, "bin", "detect")
	cmd := NewCommand(
		detect,
		buildContext.SourceDir,
		buildContext.MetadataDir,
	)

	err := RunCommand(b.Logger, cmd)
	if err != nil {
		logger.Debugf("Detection for builder '%s' failed: %s", b.Name(), err)
		
		
		return false
	}

	return true
}

func (b *Builder) Build(buildContext *BuildContext) error {
	build := filepath.Join(b.Location, "bin", "build")
	cmd := NewCommand(
		build,
		buildContext.SourceDir,
		buildContext.MetadataDir,
		buildContext.OutputDir,
	)

	err := RunCommand(b.Logger, cmd)
	if err != nil {
		return errors.Wrapf(err, "builder '%s' failed", b.Name())
	}

	return nil
}


type LaunchConfig struct {
	PeerAddress string `json:"peer_address"`
	ClientCert  []byte `json:"client_cert"`
	ClientKey   []byte `json:"client_key"`
	RootCert    []byte `json:"root_cert"`
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
	cmd := NewCommand(
		launch,
		buildContext.SourceDir,
		buildContext.MetadataDir,
		buildContext.OutputDir,
		buildContext.LaunchDir,
	)

	err = RunCommand(b.Logger, cmd)
	if err != nil {
		return errors.Wrapf(err, "builder '%s' failed", b.Name())
	}

	return nil
}



func (b *Builder) Name() string {
	return filepath.Base(b.Location)
}




func NewCommand(name string, args ...string) *exec.Cmd {
	cmd := exec.Command(name, args...)
	for _, key := range []string{"GOCACHE", "HOME", "LD_LIBRARY_PATH", "LIBPATH", "PATH", "TMPDIR"} {
		if val, ok := os.LookupEnv(key); ok {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, val))
		}
	}
	return cmd
}

func RunCommand(logger *flogging.FabricLogger, cmd *exec.Cmd) error {
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	is := bufio.NewReader(stderr)
	for done := false; !done; {
		
		line, err := is.ReadString('\n')
		switch err {
		case nil:
			logger.Info(strings.TrimSuffix(line, "\n"))
		case io.EOF:
			if len(line) > 0 {
				logger.Info(line)
			}
			done = true
		default:
			logger.Error("error reading command output", err)
			return err
		}
	}

	return cmd.Wait()
}
