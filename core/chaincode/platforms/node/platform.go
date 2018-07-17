

package node

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/mcc-github/blockchain/common/flogging"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/ccmetadata"
	"github.com/mcc-github/blockchain/core/chaincode/platforms/util"
	cutil "github.com/mcc-github/blockchain/core/container/util"
	pb "github.com/mcc-github/blockchain/protos/peer"
)

var logger = flogging.MustGetLogger("node-platform")


type Platform struct {
}


func pathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}


func (nodePlatform *Platform) ValidateSpec(spec *pb.ChaincodeSpec) error {
	path, err := url.Parse(spec.ChaincodeId.Path)
	if err != nil || path == nil {
		return fmt.Errorf("invalid path: %s", err)
	}

	
	if path.Scheme == "" {
		pathToCheck, err := filepath.Abs(spec.ChaincodeId.Path)
		if err != nil {
			return fmt.Errorf("error obtaining absolute path of the chaincode: %s", err)
		}

		exists, err := pathExists(pathToCheck)
		if err != nil {
			return fmt.Errorf("error validating chaincode path: %s", err)
		}
		if !exists {
			return fmt.Errorf("path to chaincode does not exist: %s", spec.ChaincodeId.Path)
		}
	}
	return nil
}

func (nodePlatform *Platform) ValidateDeploymentSpec(cds *pb.ChaincodeDeploymentSpec) error {

	if cds.CodePackage == nil || len(cds.CodePackage) == 0 {
		
		return nil
	}

	
	
	
	
	
	
	
	re := regexp.MustCompile(`(/)?src/.*`)
	is := bytes.NewReader(cds.CodePackage)
	gr, err := gzip.NewReader(is)
	if err != nil {
		return fmt.Errorf("failure opening codepackage gzip stream: %s", err)
	}
	tr := tar.NewReader(gr)

	var foundPackageJson = false
	for {
		header, err := tr.Next()
		if err != nil {
			
			break
		}

		
		
		
		if !re.MatchString(header.Name) {
			return fmt.Errorf("illegal file detected in payload: \"%s\"", header.Name)
		}
		if header.Name == "src/package.json" {
			foundPackageJson = true
		}
		
		
		
		
		
		
		
		
		
		if header.Mode&^0100666 != 0 {
			return fmt.Errorf("illegal file mode detected for file %s: %o", header.Name, header.Mode)
		}
	}
	if !foundPackageJson {
		return fmt.Errorf("no package.json found at the root of the chaincode package")
	}

	return nil
}


func (nodePlatform *Platform) GetDeploymentPayload(spec *pb.ChaincodeSpec) ([]byte, error) {

	var err error

	
	
	
	payload := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(payload)
	tw := tar.NewWriter(gw)

	folder := spec.ChaincodeId.Path
	if folder == "" {
		return nil, errors.New("ChaincodeSpec's path cannot be empty")
	}

	
	if folder[len(folder)-1] == '/' {
		folder = folder[:len(folder)-1]
	}

	logger.Debugf("Packaging node.js project from path %s", folder)

	if err = cutil.WriteFolderToTarPackage(tw, folder, "node_modules", nil, nil); err != nil {

		logger.Errorf("Error writing folder to tar package %s", err)
		return nil, fmt.Errorf("Error writing Chaincode package contents: %s", err)
	}

	
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("Error writing Chaincode package contents: %s", err)
	}

	tw.Close()
	gw.Close()

	return payload.Bytes(), nil
}

func (nodePlatform *Platform) GenerateDockerfile(cds *pb.ChaincodeDeploymentSpec) (string, error) {

	var buf []string

	buf = append(buf, "FROM "+cutil.GetDockerfileFromConfig("chaincode.node.runtime"))
	buf = append(buf, "ADD binpackage.tar /usr/local/src")

	dockerFileContents := strings.Join(buf, "\n")

	return dockerFileContents, nil
}

func (nodePlatform *Platform) GenerateDockerBuild(cds *pb.ChaincodeDeploymentSpec, tw *tar.Writer) error {

	codepackage := bytes.NewReader(cds.CodePackage)
	binpackage := bytes.NewBuffer(nil)
	err := util.DockerBuild(util.DockerBuildOptions{
		Cmd:          fmt.Sprint("cp -R /chaincode/input/src/. /chaincode/output && cd /chaincode/output && npm install --production"),
		InputStream:  codepackage,
		OutputStream: binpackage,
	})
	if err != nil {
		return err
	}

	return cutil.WriteBytesToPackage("binpackage.tar", binpackage.Bytes(), tw)
}


func (nodePlatform *Platform) GetMetadataProvider(cds *pb.ChaincodeDeploymentSpec) ccmetadata.MetadataProvider {
	return &ccmetadata.TargzMetadataProvider{cds}
}
