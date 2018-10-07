



package docker

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/docker/docker/pkg/fileutils"
	"github.com/fsouza/go-dockerclient/internal/archive"
)

func createTarStream(srcPath, dockerfilePath string) (io.ReadCloser, error) {
	srcPath, err := filepath.Abs(srcPath)
	if err != nil {
		return nil, err
	}

	excludes, err := parseDockerignore(srcPath)
	if err != nil {
		return nil, err
	}

	includes := []string{"."}

	
	
	
	
	
	
	
	
	
	forceIncludeFiles := []string{".dockerignore", dockerfilePath}

	for _, includeFile := range forceIncludeFiles {
		if includeFile == "" {
			continue
		}
		keepThem, err := fileutils.Matches(includeFile, excludes)
		if err != nil {
			return nil, fmt.Errorf("cannot match .dockerfile: '%s', error: %s", includeFile, err)
		}
		if keepThem {
			includes = append(includes, includeFile)
		}
	}

	if err := validateContextDirectory(srcPath, excludes); err != nil {
		return nil, err
	}
	tarOpts := &archive.TarOptions{
		ExcludePatterns: excludes,
		IncludeFiles:    includes,
		Compression:     archive.Uncompressed,
		NoLchown:        true,
	}
	return archive.TarWithOptions(srcPath, tarOpts)
}




func validateContextDirectory(srcPath string, excludes []string) error {
	return filepath.Walk(filepath.Join(srcPath, "."), func(filePath string, f os.FileInfo, err error) error {
		
		if relFilePath, relErr := filepath.Rel(srcPath, filePath); relErr != nil {
			return relErr
		} else if skip, matchErr := fileutils.Matches(relFilePath, excludes); matchErr != nil {
			return matchErr
		} else if skip {
			if f.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		if err != nil {
			if os.IsPermission(err) {
				return fmt.Errorf("can't stat '%s'", filePath)
			}
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}

		
		
		if f.Mode()&(os.ModeSymlink|os.ModeNamedPipe) != 0 {
			return nil
		}

		if !f.IsDir() {
			currentFile, err := os.Open(filePath)
			if err != nil && os.IsPermission(err) {
				return fmt.Errorf("no permission to read from '%s'", filePath)
			}
			currentFile.Close()
		}
		return nil
	})
}

func parseDockerignore(root string) ([]string, error) {
	var excludes []string
	ignore, err := ioutil.ReadFile(path.Join(root, ".dockerignore"))
	if err != nil && !os.IsNotExist(err) {
		return excludes, fmt.Errorf("error reading .dockerignore: '%s'", err)
	}
	excludes = strings.Split(string(ignore), "\n")

	return excludes, nil
}
