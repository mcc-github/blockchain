/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package golang

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/mcc-github/blockchain/core/chaincode/platforms/util"
	cutil "github.com/mcc-github/blockchain/core/container/util"
	"github.com/mcc-github/blockchain/internal/ccmetadata"
	pb "github.com/mcc-github/blockchain/protos/peer"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
)


type Platform struct{}


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

func decodeUrl(path string) (string, error) {
	var urlLocation string
	if strings.HasPrefix(path, "http://") {
		urlLocation = path[7:]
	} else if strings.HasPrefix(path, "https://") {
		urlLocation = path[8:]
	} else {
		urlLocation = path
	}

	if len(urlLocation) < 2 {
		return "", errors.New("ChaincodeSpec's path/URL invalid")
	}

	if strings.LastIndex(urlLocation, "/") == len(urlLocation)-1 {
		urlLocation = urlLocation[:len(urlLocation)-1]
	}

	return urlLocation, nil
}

func getGopath() (string, error) {
	env, err := getGoEnv()
	if err != nil {
		return "", err
	}
	
	splitGoPath := filepath.SplitList(env["GOPATH"])
	if len(splitGoPath) == 0 {
		return "", fmt.Errorf("invalid GOPATH environment variable value: %s", env["GOPATH"])
	}
	return splitGoPath[0], nil
}

func filter(vs []string, f func(string) bool) []string {
	vsf := make([]string, 0)
	for _, v := range vs {
		if f(v) {
			vsf = append(vsf, v)
		}
	}
	return vsf
}


func (p *Platform) Name() string {
	return pb.ChaincodeSpec_GOLANG.String()
}


func (p *Platform) ValidatePath(rawPath string) error {
	path, err := url.Parse(rawPath)
	if err != nil || path == nil {
		return fmt.Errorf("invalid path: %s", err)
	}

	
	
	
	if path.Scheme == "" {
		gopath, err := getGopath()
		if err != nil {
			return err
		}
		pathToCheck := filepath.Join(gopath, "src", rawPath)
		exists, err := pathExists(pathToCheck)
		if err != nil {
			return fmt.Errorf("error validating chaincode path: %s", err)
		}
		if !exists {
			return fmt.Errorf("path to chaincode does not exist: %s", pathToCheck)
		}
	}
	return nil
}

func (p *Platform) ValidateCodePackage(code []byte) error {
	
	
	
	
	
	
	
	
	
	
	re := regexp.MustCompile(`^(/)?(src|META-INF)/.*`)
	is := bytes.NewReader(code)
	gr, err := gzip.NewReader(is)
	if err != nil {
		return fmt.Errorf("failure opening codepackage gzip stream: %s", err)
	}
	tr := tar.NewReader(gr)

	for {
		header, err := tr.Next()
		if err != nil {
			
			break
		}

		
		
		
		if !re.MatchString(header.Name) {
			return fmt.Errorf("illegal file detected in payload: \"%s\"", header.Name)
		}

		
		
		
		
		
		
		
		
		
		if header.Mode&^0100666 != 0 {
			return fmt.Errorf("illegal file mode detected for file %s: %o", header.Name, header.Mode)
		}
	}

	return nil
}

type Sources []SourceDescriptor

func (s Sources) Len() int {
	return len(s)
}

func (s Sources) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s Sources) Less(i, j int) bool {
	return strings.Compare(s[i].Name, s[j].Name) < 0
}






func vendorDependencies(pkg string, files Sources) {

	exclusions := make([]string, 0)
	elements := strings.Split(pkg, "/")

	
	
	
	
	
	
	
	
	
	prev := filepath.Join("src")
	for _, element := range elements {
		curr := filepath.Join(prev, element)
		vendor := filepath.Join(curr, "vendor")
		exclusions = append(exclusions, vendor)
		prev = curr
	}

	
	
	
	exclusions = append(exclusions, filepath.Join("src", pkg))

	count := len(files)
	sem := make(chan bool, count)

	
	
	
	
	
	vendorPath := filepath.Join("src", pkg, "vendor")
	for i, file := range files {
		go func(i int, file SourceDescriptor) {
			excluded := false

			for _, exclusion := range exclusions {
				if strings.HasPrefix(file.Name, exclusion) == true {
					excluded = true
					break
				}
			}

			if excluded == false {
				origName := file.Name
				file.Name = strings.Replace(origName, "src", vendorPath, 1)
				logger.Debugf("vendoring %s -> %s", origName, file.Name)
			}

			files[i] = file
			sem <- true
		}(i, file)
	}

	for i := 0; i < count; i++ {
		<-sem
	}
}


func (p *Platform) GetDeploymentPayload(path string) ([]byte, error) {

	var err error

	
	
	
	code, err := getCode(path)
	if err != nil {
		return nil, err
	}
	if code.Cleanup != nil {
		defer code.Cleanup()
	}

	
	
	
	env, err := getGoEnv()
	if err != nil {
		return nil, err
	}
	gopaths := splitEnvPaths(env["GOPATH"])
	goroots := splitEnvPaths(env["GOROOT"])
	gopaths[code.Gopath] = true
	env["GOPATH"] = flattenEnvPaths(gopaths)

	
	
	
	imports, err := listImports(env, code.Pkg)
	if err != nil {
		return nil, fmt.Errorf("Error obtaining imports: %s", err)
	}

	
	
	
	var provided = map[string]bool{
		"github.com/mcc-github/blockchain/core/chaincode/shim": true,
		"github.com/mcc-github/blockchain/protos/peer":         true,
	}

	
	var pseudo = map[string]bool{
		"C": true,
	}

	imports = filter(imports, func(pkg string) bool {
		
		if _, ok := provided[pkg]; ok == true {
			logger.Debugf("Discarding provided package %s", pkg)
			return false
		}

		
		if _, ok := pseudo[pkg]; ok == true {
			logger.Debugf("Discarding pseudo-package %s", pkg)
			return false
		}

		
		for goroot := range goroots {
			fqp := filepath.Join(goroot, "src", pkg)
			exists, err := pathExists(fqp)
			if err == nil && exists {
				logger.Debugf("Discarding GOROOT package %s", pkg)
				return false
			}
		}

		
		logger.Debugf("Accepting import: %s", pkg)
		return true
	})

	
	
	
	
	deps := make(map[string]bool)

	for _, pkg := range imports {
		
		
		
		transitives, err := listDeps(env, pkg)
		if err != nil {
			return nil, fmt.Errorf("Error obtaining dependencies for %s: %s", pkg, err)
		}

		
		
		

		
		deps[pkg] = true

		
		for _, dep := range transitives {
			deps[dep] = true
		}
	}

	
	delete(deps, "")

	
	
	
	fileMap, err := findSource(code.Gopath, code.Pkg)
	if err != nil {
		return nil, err
	}

	
	
	
	
	for dep := range deps {
		logger.Debugf("processing dep: %s", dep)

		
		
		
		
		for gopath := range gopaths {
			fqp := filepath.Join(gopath, "src", dep)
			exists, err := pathExists(fqp)

			logger.Debugf("checking: %s exists: %v", fqp, exists)

			if err == nil && exists {

				
				files, err := findSource(gopath, dep)
				if err != nil {
					return nil, err
				}

				
				for _, file := range files {
					fileMap[file.Name] = file
				}
			}
		}
	}

	logger.Debugf("done")

	
	
	
	files := make(Sources, 0)
	for _, file := range fileMap {
		files = append(files, file)
	}

	
	
	
	vendorDependencies(code.Pkg, files)

	
	
	
	sort.Sort(files)

	
	
	
	payload := bytes.NewBuffer(nil)
	gw := gzip.NewWriter(payload)
	tw := tar.NewWriter(gw)

	for _, file := range files {
		
		
		
		if file.IsMetadata {
			file.Name, err = filepath.Rel(filepath.Join("src", code.Pkg), file.Name)
			if err != nil {
				return nil, fmt.Errorf("This error was caused by bad packaging of the metadata.  The file [%s] is marked as MetaFile, however not located under META-INF   Error:[%s]", file.Name, err)
			}

			
			_, filename := filepath.Split(file.Name)

			
			
			if strings.HasPrefix(filename, ".") {
				logger.Warningf("Ignoring hidden file in metadata directory: %s", file.Name)
				continue
			}

			fileBytes, err := ioutil.ReadFile(file.Path)
			if err != nil {
				return nil, err
			}

			
			
			err = ccmetadata.ValidateMetadataFile(file.Name, fileBytes)
			if err != nil {
				return nil, err
			}
		}

		err = cutil.WriteFileToPackage(file.Path, file.Name, tw)
		if err != nil {
			return nil, fmt.Errorf("Error writing %s to tar: %s", file.Name, err)
		}
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

func (p *Platform) GenerateDockerfile() (string, error) {
	var buf []string
	buf = append(buf, "FROM "+util.GetDockerfileFromConfig("chaincode.golang.runtime"))
	buf = append(buf, "ADD binpackage.tar /usr/local/bin")

	return strings.Join(buf, "\n"), nil
}

const staticLDFlagsOpts = "-ldflags \"-linkmode external -extldflags '-static'\""
const dynamicLDFlagsOpts = ""

func getLDFlagsOpts() string {
	if viper.GetBool("chaincode.golang.dynamicLink") {
		return dynamicLDFlagsOpts
	}
	return staticLDFlagsOpts
}

func (p *Platform) DockerBuildOptions(path string) (util.DockerBuildOptions, error) {
	pkgname, err := decodeUrl(path)
	if err != nil {
		return util.DockerBuildOptions{}, fmt.Errorf("could not decode url: %s", err)
	}

	ldflagsOpt := getLDFlagsOpts()
	return util.DockerBuildOptions{
		Cmd: fmt.Sprintf("GOPATH=/chaincode/input:$GOPATH go build  %s -o /chaincode/output/chaincode %s", ldflagsOpt, pkgname),
	}, nil
}
