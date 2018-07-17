package convert

import (
	"fmt"
	"go/build"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)


func RewritePackage(packageName string) {
	pkg, err := packageWithName(packageName)
	if err != nil {
		panic(fmt.Sprintf("unexpected error reading package: '%s'\n%s\n", packageName, err.Error()))
	}

	for _, filename := range findTestsInPackage(pkg) {
		rewriteTestsInFile(filename)
	}
	return
}


func findTestsInPackage(pkg *build.Package) (testfiles []string) {
	for _, file := range append(pkg.TestGoFiles, pkg.XTestGoFiles...) {
		testfiles = append(testfiles, filepath.Join(pkg.Dir, file))
	}

	dirFiles, err := ioutil.ReadDir(pkg.Dir)
	if err != nil {
		panic(fmt.Sprintf("unexpected error reading dir: '%s'\n%s\n", pkg.Dir, err.Error()))
	}

	re := regexp.MustCompile(`^[._]`)

	for _, file := range dirFiles {
		if !file.IsDir() {
			continue
		}

		if re.Match([]byte(file.Name())) {
			continue
		}

		packageName := filepath.Join(pkg.ImportPath, file.Name())
		subPackage, err := packageWithName(packageName)
		if err != nil {
			panic(fmt.Sprintf("unexpected error reading package: '%s'\n%s\n", packageName, err.Error()))
		}

		testfiles = append(testfiles, findTestsInPackage(subPackage)...)
	}

	addGinkgoSuiteForPackage(pkg)
	goFmtPackage(pkg)
	return
}


func addGinkgoSuiteForPackage(pkg *build.Package) {
	originalDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	suite_test_file := filepath.Join(pkg.Dir, pkg.Name+"_suite_test.go")

	_, err = os.Stat(suite_test_file)
	if err == nil {
		return 
	}

	err = os.Chdir(pkg.Dir)
	if err != nil {
		panic(err)
	}

	output, err := exec.Command("ginkgo", "bootstrap").Output()

	if err != nil {
		panic(fmt.Sprintf("error running 'ginkgo bootstrap'.\nstdout: %s\n%s\n", output, err.Error()))
	}

	err = os.Chdir(originalDir)
	if err != nil {
		panic(err)
	}
}


func goFmtPackage(pkg *build.Package) {
	output, err := exec.Command("go", "fmt", pkg.ImportPath).Output()

	if err != nil {
		fmt.Printf("Warning: Error running 'go fmt %s'.\nstdout: %s\n%s\n", pkg.ImportPath, output, err.Error())
	}
}


func packageWithName(name string) (pkg *build.Package, err error) {
	pkg, err = build.Default.Import(name, ".", build.ImportMode(0))
	if err == nil {
		return
	}

	pkg, err = build.Default.Import(name, ".", build.ImportMode(1))
	return
}
