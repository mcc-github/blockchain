package gexec

import (
	"errors"
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

var (
	mu     sync.Mutex
	tmpDir string
)


func Build(packagePath string, args ...string) (compiledPath string, err error) {
	return doBuild(build.Default.GOPATH, packagePath, nil, args...)
}


func BuildWithEnvironment(packagePath string, env []string, args ...string) (compiledPath string, err error) {
	return doBuild(build.Default.GOPATH, packagePath, env, args...)
}


func BuildIn(gopath string, packagePath string, args ...string) (compiledPath string, err error) {
	return doBuild(gopath, packagePath, nil, args...)
}

func replaceGoPath(environ []string, newGoPath string) []string {
	newEnviron := []string{}
	for _, v := range environ {
		if !strings.HasPrefix(v, "GOPATH=") {
			newEnviron = append(newEnviron, v)
		}
	}
	return append(newEnviron, "GOPATH="+newGoPath)
}

func doBuild(gopath, packagePath string, env []string, args ...string) (compiledPath string, err error) {
	tmpDir, err := temporaryDirectory()
	if err != nil {
		return "", err
	}

	if len(gopath) == 0 {
		return "", errors.New("$GOPATH not provided when building " + packagePath)
	}

	executable := filepath.Join(tmpDir, path.Base(packagePath))
	if runtime.GOOS == "windows" {
		executable = executable + ".exe"
	}

	cmdArgs := append([]string{"build"}, args...)
	cmdArgs = append(cmdArgs, "-o", executable, packagePath)

	build := exec.Command("go", cmdArgs...)
	build.Env = replaceGoPath(os.Environ(), gopath)
	build.Env = append(build.Env, env...)

	output, err := build.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("Failed to build %s:\n\nError:\n%s\n\nOutput:\n%s", packagePath, err, string(output))
	}

	return executable, nil
}


func CleanupBuildArtifacts() {
	mu.Lock()
	defer mu.Unlock()
	if tmpDir != "" {
		os.RemoveAll(tmpDir)
		tmpDir = ""
	}
}

func temporaryDirectory() (string, error) {
	var err error
	mu.Lock()
	defer mu.Unlock()
	if tmpDir == "" {
		tmpDir, err = ioutil.TempDir("", "gexec_artifacts")
		if err != nil {
			return "", err
		}
	}

	return ioutil.TempDir(tmpDir, "g")
}
