/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package platforms

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"strings"

	docker "github.com/fsouza/go-dockerclient"
	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/common/metadata"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/golang"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/java"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/node"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/util"
	"github.com/pkg/errors"
)


var SupportedPlatforms = []Platform{
	&java.Platform{},
	&golang.Platform{},
	&node.Platform{},
}



type Platform interface {
	Name() string
	GenerateDockerfile() (string, error)
	DockerBuildOptions(path string) (util.DockerBuildOptions, error)
}

type PackageWriter interface {
	Write(name string, payload []byte, tw *tar.Writer) error
}

type PackageWriterWrapper func(name string, payload []byte, tw *tar.Writer) error

func (pw PackageWriterWrapper) Write(name string, payload []byte, tw *tar.Writer) error {
	return pw(name, payload, tw)
}

type Registry struct {
	Platforms     map[string]Platform
	PackageWriter PackageWriter
}

var logger = flogging.MustGetLogger("chaincode.platform")

func NewRegistry(platformTypes ...Platform) *Registry {
	platforms := make(map[string]Platform)
	for _, platform := range platformTypes {
		if _, ok := platforms[platform.Name()]; ok {
			logger.Panicf("Multiple platforms of the same name specified: %s", platform.Name())
		}
		platforms[platform.Name()] = platform
	}
	return &Registry{
		Platforms:     platforms,
		PackageWriter: PackageWriterWrapper(writeBytesToPackage),
	}
}

func (r *Registry) GenerateDockerfile(ccType string) (string, error) {
	platform, ok := r.Platforms[ccType]
	if !ok {
		return "", fmt.Errorf("Unknown chaincodeType: %s", ccType)
	}

	var buf []string

	
	
	
	base, err := platform.GenerateDockerfile()
	if err != nil {
		return "", fmt.Errorf("Failed to generate platform-specific Dockerfile: %s", err)
	}
	buf = append(buf, base)
	buf = append(buf, fmt.Sprintf(`LABEL %s.chaincode.type="%s" \`, metadata.BaseDockerLabel, ccType))
	buf = append(buf, fmt.Sprintf(`      %s.version="%s"`, metadata.BaseDockerLabel, metadata.Version))
	
	
	
	
	buf = append(buf, fmt.Sprintf("ENV CORE_CHAINCODE_BUILDLEVEL=%s", metadata.Version))

	
	
	
	contents := strings.Join(buf, "\n")
	logger.Debugf("\n%s", contents)

	return contents, nil
}

func (r *Registry) StreamDockerBuild(ccType, path string, codePackage io.Reader, inputFiles map[string][]byte, tw *tar.Writer, client *docker.Client) error {
	var err error

	
	
	
	platform, ok := r.Platforms[ccType]
	if !ok {
		return fmt.Errorf("could not find platform of type: %s", ccType)
	}

	
	
	
	for name, data := range inputFiles {
		err = r.PackageWriter.Write(name, data, tw)
		if err != nil {
			return fmt.Errorf(`Failed to inject "%s": %s`, name, err)
		}
	}

	buildOptions, err := platform.DockerBuildOptions(path)
	if err != nil {
		return errors.Wrap(err, "platform failed to create docker build options")
	}

	output := &bytes.Buffer{}
	buildOptions.InputStream = codePackage
	buildOptions.OutputStream = output

	err = util.DockerBuild(buildOptions, client)
	if err != nil {
		return errors.Wrap(err, "docker build failed")
	}

	return writeBytesToPackage("binpackage.tar", output.Bytes(), tw)
}

func writeBytesToPackage(name string, payload []byte, tw *tar.Writer) error {
	err := tw.WriteHeader(&tar.Header{
		Name: name,
		Size: int64(len(payload)),
		Mode: 0100644,
	})
	if err != nil {
		return err
	}

	_, err = tw.Write(payload)
	if err != nil {
		return err
	}

	return nil
}

func (r *Registry) GenerateDockerBuild(ccType, path string, codePackage io.Reader, client *docker.Client) (io.Reader, error) {
	inputFiles := make(map[string][]byte)

	
	
	
	dockerFile, err := r.GenerateDockerfile(ccType)
	if err != nil {
		return nil, fmt.Errorf("Failed to generate a Dockerfile: %s", err)
	}

	inputFiles["Dockerfile"] = []byte(dockerFile)

	
	
	
	input, output := io.Pipe()

	go func() {
		gw := gzip.NewWriter(output)
		tw := tar.NewWriter(gw)
		err := r.StreamDockerBuild(ccType, path, codePackage, inputFiles, tw, client)
		if err != nil {
			logger.Error(err)
		}

		tw.Close()
		gw.Close()
		output.CloseWithError(err)
	}()

	return input, nil
}
