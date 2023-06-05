// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"path"
	"runtime"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/terraform"
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
		resultContent, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("Error occurred while reading file at path: %s\n, error: %s\n", path, err)
		}
		checkSums := genFileChecksums(resultContent)

		resource.TestCheckResourceAttr(resourceName, "content", string(resultContent))
		resource.TestCheckResourceAttr(resourceName, "content_md5", checkSums.md5Hex)
		resource.TestCheckResourceAttr(resourceName, "content_sha1", checkSums.sha1Hex)
		resource.TestCheckResourceAttr(resourceName, "content_sha256", checkSums.sha256Hex)
		resource.TestCheckResourceAttr(resourceName, "content_base64sha256", checkSums.sha256Base64)
		resource.TestCheckResourceAttr(resourceName, "content_sha512", checkSums.sha512Hex)
		resource.TestCheckResourceAttr(resourceName, "content_base64sha512", checkSums.sha512Base64)

		return nil
	}
}

func checkFilePermissions(destinationFilePath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		filePermission := os.FileMode(0600)
		fileInfo, err := os.Stat(destinationFilePath)
		if err != nil {
			return fmt.Errorf("Error occurred while retrieving file info at path: %s\n, error: %s\n",
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

func checkDirectoryPermissions(destinationFilePath string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		directoryPermission := os.FileMode(0700)
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

func createSourceFile(sourceFilePath, sourceContent string) error {
	return os.WriteFile(sourceFilePath, []byte(sourceContent), 0644)
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
