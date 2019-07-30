

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
	"github.com/mcc-github/blockchain/core/common/ccprovider"
	cutil "github.com/mcc-github/blockchain/core/container/util"
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
		PackageWriter: PackageWriterWrapper(cutil.WriteBytesToPackage),
	}
}

func (r *Registry) GenerateDockerfile(ccType, name, version string) (string, error) {
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

	
	
	
	
	
	buf = append(buf, fmt.Sprintf(`LABEL %s.chaincode.id.name="%s" \`, metadata.BaseDockerLabel, name))
	 buf = append(buf, fmt.Sprintf(`      %s.chaincode.id.version="%s" \`, metadata.BaseDockerLabel, version))

	buf = append(buf, fmt.Sprintf(`      %s.chaincode.type="%s" \`, metadata.BaseDockerLabel, ccType))
	buf = append(buf, fmt.Sprintf(`      %s.version="%s"`, metadata.BaseDockerLabel, metadata.Version))
	
	
	
	
	buf = append(buf, fmt.Sprintf("ENV CORE_CHAINCODE_BUILDLEVEL=%s", metadata.Version))

	
	
	
	contents := strings.Join(buf, "\n")
	logger.Debugf("\n%s", contents)

	return contents, nil
}

func (r *Registry) StreamDockerBuild(ccType, path string, codePackage []byte, inputFiles map[string][]byte, tw *tar.Writer, client *docker.Client) error {
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
	buildOptions.InputStream = bytes.NewReader(codePackage)
	buildOptions.OutputStream = output

	err = util.DockerBuild(buildOptions, client)
	if err != nil {
		return errors.Wrap(err, "docker build failed")
	}

	return cutil.WriteBytesToPackage("binpackage.tar", output.Bytes(), tw)
}

func (r *Registry) GenerateDockerBuild(ccType, path, name, version string, codePackage []byte, client *docker.Client) (io.Reader, error) {
	inputFiles := make(map[string][]byte)

	
	
	
	dockerFile, err := r.GenerateDockerfile(ccType, name, version)
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

type Builder struct {
	Registry *Registry
	Client   *docker.Client
}

func (b *Builder) GenerateDockerBuild(ccci *ccprovider.ChaincodeContainerInfo, codePackage []byte) (io.Reader, error) {
	return b.Registry.GenerateDockerBuild(ccci.Type, ccci.Path, ccci.Name, ccci.Version, codePackage, b.Client)
}
