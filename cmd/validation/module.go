package validation

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"golang.org/x/mod/modfile"
	"golang.org/x/mod/semver"
)

func DirectoryExists(directory string) bool {
	folderInfo, err := os.Stat(directory)
	if os.IsNotExist(err) {
		fmt.Printf("ERROR: No Such Directory %v\n", directory)
		return false
	} else if err != nil {
		fmt.Printf("ERROR: %v", err)
		return false
	}

	if !folderInfo.IsDir() {
		fmt.Printf("ERROR: Not a directory %v", directory)
		return false
	}
	return true
}

func FileExists(file string) bool {
	fileInfo, err := os.Stat(file)
	if os.IsNotExist(err) {
		fmt.Printf("ERROR: No Such File %v\n", file)
		return false
	} else if err != nil {
		fmt.Printf("ERROR: %v", err)
		return false
	}

	if fileInfo.IsDir() {
		fmt.Printf("ERROR: Not a file %v\n", file)
		return false
	}
	return true
}

func ValidVersion(version string) bool {
	if !semver.IsValid(version) {
		fmt.Println("ERROR: Not a valid SemVer version")
		return false
	}
	return true
}

func ValidateModule(file string) (bool, string) {
	ofile, err := os.Open(file)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return false, ""
	}

	bt, err := ioutil.ReadAll(ofile)
	if err != nil {
		fmt.Printf("ERROR: %v", err)
		return false, ""
	}

	module := modfile.ModulePath(bt)

	log.Println(module)

	return true, module
}
