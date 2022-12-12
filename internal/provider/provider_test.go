package provider

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"runtime"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func protoV5ProviderFactories() map[string]func() (tfprotov5.ProviderServer, error) {
	return map[string]func() (tfprotov5.ProviderServer, error){
		"local": providerserver.NewProtocol5WithError(New()),
	}
}

func providerVersion233() map[string]resource.ExternalProvider {
	return map[string]resource.ExternalProvider{
		"local": {
			VersionConstraint: "2.2.3",
			Source:            "hashicorp/local",
		},
	}
}

func checkFileDeleted(shouldNotExistFile string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if _, err := os.Stat(shouldNotExistFile); os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("file %s was not deleted", shouldNotExistFile)
	}
}

func checkFileCreation(resourceName, path string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		resultContent, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Error occured while reading file at path: %s\n, error: %s\n", path, err)
		}

		resource.TestCheckResourceAttr(resourceName, "content", string(resultContent))

		return nil
	}
}

func checkFilePermissions(destinationFilePath string, filePermission os.FileMode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		if runtime.GOOS == "windows" {
			// skip all checks if windows
			return nil
		}

		fileInfo, err := os.Stat(destinationFilePath)
		if err != nil {
			return fmt.Errorf("Error occured while retrieving file info at path: %s\n, error: %s\n",
				destinationFilePath, err)
		}

		if fileInfo.Mode() != filePermission {
			return fmt.Errorf(
				"File permission.\nexpected:%s\ngot: %s\n",
				filePermission, fileInfo.Mode())
		}

		return nil
	}
}

func checkDirectoryPermissions(destinationFilePath string, directoryPermission os.FileMode) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		dirInfo, _ := os.Stat(path.Dir(destinationFilePath))
		// we have to use FileMode.Perm() here, otherwise directory bit causes issues
		if dirInfo.Mode().Perm() != directoryPermission.Perm() {
			return fmt.Errorf(
				"Directory permission. \nexpected:%s\ngot: %s\n",
				directoryPermission, dirInfo.Mode().Perm())
		}

		return nil
	}
}

func createSourceFile(sourceContent string) error {
	return ioutil.WriteFile("source_file", []byte(sourceContent), 0644)
}

func checkDirExists(destinationFilePath string, isDirExist *bool) func() {
	return func() {
		// if directory already existed prior to check, skip check
		if _, err := os.Stat(path.Dir(destinationFilePath)); !os.IsNotExist(err) {
			*isDirExist = true
		}
	}
}

func skipTestsWindows() func() (bool, error) {
	return func() (bool, error) {
		if runtime.GOOS == "windows" {
			// skip all checks if windows
			return true, nil
		}
		return false, nil
	}
}
