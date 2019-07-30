/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package chaincode

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/mcc-github/blockchain/core/chaincode/persistence"
	cutil "github.com/mcc-github/blockchain/core/container/util"
	"github.com/mcc-github/blockchain/internal/peer/packaging"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)



type PlatformRegistry interface {
	GetDeploymentPayload(ccType, path string) ([]byte, error)
}



type Packager struct {
	Command          *cobra.Command
	Input            *PackageInput
	PlatformRegistry PlatformRegistry
	Writer           Writer
}



type PackageInput struct {
	OutputFile string
	Path       string
	Type       string
	Label      string
}


func (p *PackageInput) Validate() error {
	if p.Path == "" {
		return errors.New("chaincode path must be specified")
	}
	if p.Type == "" {
		return errors.New("chaincode language must be specified")
	}
	if p.OutputFile == "" {
		return errors.New("output file must be specified")
	}
	if p.Label == "" {
		return errors.New("package label must be specified")
	}

	return nil
}


func PackageCmd(p *Packager) *cobra.Command {
	chaincodePackageCmd := &cobra.Command{
		Use:       "package [outputfile]",
		Short:     "Package a chaincode",
		Long:      "Package a chaincode and write the package to a file.",
		ValidArgs: []string{"1"},
		RunE: func(cmd *cobra.Command, args []string) error {
			if p == nil {
				pr := packaging.NewRegistry(packaging.SupportedPlatforms...)

				p = &Packager{
					PlatformRegistry: pr,
					Writer:           &persistence.FilesystemIO{},
				}
			}
			p.Command = cmd

			return p.PackageChaincode(args)
		},
	}
	flagList := []string{
		"label",
		"lang",
		"path",
		"peerAddresses",
		"tlsRootCertFiles",
		"connectionProfile",
	}
	attachFlags(chaincodePackageCmd, flagList)

	return chaincodePackageCmd
}


func (p *Packager) PackageChaincode(args []string) error {
	if p.Command != nil {
		
		p.Command.SilenceUsage = true
	}

	if len(args) != 1 {
		return errors.New("invalid number of args. expected only the output file")
	}
	p.setInput(args[0])

	return p.Package()
}

func (p *Packager) setInput(outputFile string) {
	p.Input = &PackageInput{
		OutputFile: outputFile,
		Path:       chaincodePath,
		Type:       chaincodeLang,
		Label:      packageLabel,
	}
}



func (p *Packager) Package() error {
	err := p.Input.Validate()
	if err != nil {
		return err
	}

	pkgTarGzBytes, err := p.getTarGzBytes()
	if err != nil {
		return err
	}

	dir, name := filepath.Split(p.Input.OutputFile)
	
	
	if dir, err = filepath.Abs(dir); err != nil {
		return err
	}
	err = p.Writer.WriteFile(dir, name, pkgTarGzBytes)
	if err != nil {
		err = errors.Wrapf(err, "error writing chaincode package to %s", p.Input.OutputFile)
		logger.Error(err.Error())
		return err
	}

	return nil
}

func (p *Packager) getTarGzBytes() ([]byte, error) {
	payload := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(payload)
	tw := tar.NewWriter(gw)

	metadataBytes, err := toJSON(p.Input.Path, p.Input.Type, p.Input.Label)
	if err != nil {
		return nil, err
	}
	err = cutil.WriteBytesToPackage("Chaincode-Package-Metadata.json", metadataBytes, tw)
	if err != nil {
		return nil, errors.Wrap(err, "error writing package metadata to tar")
	}

	codeBytes, err := p.PlatformRegistry.GetDeploymentPayload(strings.ToUpper(p.Input.Type), p.Input.Path)
	if err != nil {
		return nil, errors.WithMessage(err, "error getting chaincode bytes")
	}

	codePackageName := "Code-Package.tar.gz"
	if strings.ToLower(p.Input.Type) == "car" {
		codePackageName = "Code-Package.car"
	}

	err = cutil.WriteBytesToPackage(codePackageName, codeBytes, tw)
	if err != nil {
		return nil, errors.Wrap(err, "error writing package code bytes to tar")
	}

	err = tw.Close()
	if err == nil {
		err = gw.Close()
	}
	if err != nil {
		return nil, errors.Wrapf(err, "failed to create tar for chaincode")
	}

	return payload.Bytes(), nil
}


type PackageMetadata struct {
	Path  string `json:"Path"`
	Type  string `json:"Type"`
	Label string `json:"Label"`
}

func toJSON(path, ccType, label string) ([]byte, error) {
	metadata := &PackageMetadata{
		Path:  path,
		Type:  ccType,
		Label: label,
	}

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		return nil, errors.Wrap(err, "failed to marshal chaincode package metadata into JSON")
	}

	return metadataBytes, nil
}
