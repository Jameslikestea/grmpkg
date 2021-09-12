package cmd

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/spf13/cobra"
	"golang.org/x/mod/module"
	"golang.org/x/mod/zip"
	"grmpkg.com/grmpkg/cmd/validation"
)

var (
	directory string
	version   string
)

var packageCmd = &cobra.Command{
	Use:   "package",
	Short: "Create a zip module file",
	Long:  `Uses the standard go module format for creating a zip file`,

	Run: packageCommand,
}

func init() {
	rootCmd.AddCommand(packageCmd)
	packageCmd.Flags().StringVarP(&directory, "directory", "d", ".", "Specify the directory of the go module (must contain a go.mod file)")
	packageCmd.Flags().StringVarP(&version, "version", "v", "", "Specify the version to package")
}

func packageCommand(command *cobra.Command, args []string) {
	if !validation.DirectoryExists(directory) {
		return
	}
	if !validation.FileExists(fmt.Sprintf("%s/go.mod", directory)) {
		return
	}
	if !validation.ValidVersion(version) {
		return
	}
	moduleValid, moduleName := validation.ValidateModule(fmt.Sprintf("%s/go.mod", directory))
	if !moduleValid {
		return
	}

	f, err := ioutil.TempFile("", "grmpkg_")
	if err != nil {
		log.Fatalf("Cannot create file: %v", err)
		return
	}

	log.Printf("Creating file at: %s", f.Name())

	err = zip.CreateFromDir(f, module.Version{Path: moduleName, Version: version}, directory)
	if err != nil {
		log.Fatalf("Could not write file: %v", err)
		return
	}
	f.Close()
}
