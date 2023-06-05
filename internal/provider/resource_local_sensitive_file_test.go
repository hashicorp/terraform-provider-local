// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLocalSensitiveFile_Basic(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_sensitive_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccConfigLocalSensitiveFileContent("This is some sensitive content", f),
				Check:  checkFileCreation("local_sensitive_file_resource.test", f),
			},
			{
				Config: testAccConfigLocalSensitiveFileEncodedBase64Content("VGhpcyBpcyBzb21lIGJhc2U2NCBjb250ZW50", f),
				Check:  checkFileCreation("local_sensitive_file_resource.test", f),
			},
			{
				Config: testAccConfigLocalSensitiveFileDecodedBase64Content("This is some base64 content", f),
				Check:  checkFileCreation("local_sensitive_file_resource.test", f),
			},
		},
		CheckDestroy: checkFileDeleted(f),
	})
}

func TestLocalSensitiveFile_source(t *testing.T) {
	sourceDirPath := t.TempDir()
	sourceFilePath := filepath.Join(sourceDirPath, "source_file")
	sourceFilePath = strings.ReplaceAll(sourceFilePath, `\`, `\\`)
	// create a local file that will be used as the "source" file
	if err := createSourceFile(sourceFilePath, "local file content"); err != nil {
		t.Fatal(err)
	}

	destinationDirPath := t.TempDir()
	destinationFilePath := filepath.Join(destinationDirPath, "new_file")
	destinationFilePath = strings.ReplaceAll(destinationFilePath, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccConfigLocalSensitiveSourceFile(sourceFilePath, destinationFilePath),
				Check:  checkFileCreation("local_sensitive_file_resource.test", destinationFilePath),
			},
		},
		CheckDestroy: checkFileDeleted(destinationFilePath),
	})
}

func TestLocalSensitiveFile_Permissions(t *testing.T) {
	destinationDirPath := t.TempDir()
	destinationFilePath := filepath.Join(destinationDirPath, "local_sensitive_file")
	destinationFilePath = strings.ReplaceAll(destinationFilePath, `\`, `\\`)
	isDirExist := false

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				PreConfig: checkDirExists(destinationDirPath, &isDirExist),
				SkipFunc:  skipTestsWindows(),
				Config: fmt.Sprintf(`
					resource "local_sensitive_file" "file" {
						content              = "This is some content"
						filename             = "%s"
						file_permission      = "0600"
						directory_permission = "0700"
					}`, destinationFilePath,
				),
				Check: r.ComposeTestCheckFunc(
					checkFilePermissions(destinationFilePath),
					checkDirectoryPermissions(destinationFilePath),
				),
			},
		},
		ErrorCheck: func(err error) error {
			if match, _ := regexp.MatchString("Directory permission.", err.Error()); match && isDirExist {
				return nil
			}
			return err
		},
		CheckDestroy: checkFileDeleted(destinationFilePath),
	})
}

func TestLocalSensitiveFile_Validators(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "local_sensitive_file" "file" {
				  filename = "%s"
				}`, f),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				Config: fmt.Sprintf(`
				resource "local_sensitive_file" "file" {
                  content = "content"
				  content_base64 = "VGhpcyBpcyBzb21lIGJhc2U2NCBjb250ZW50"
				  filename = "%s"
				}`, f),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
		},
	})
}

func TestLocalSensitiveFile_Upgrade(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_sensitive_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.Test(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: providerVersion233(),
				Config:            testAccConfigLocalSensitiveFileContent("This is some content", f),
				Check:             checkFileCreation("local_sensitive_file_resource.test", f),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigLocalSensitiveFileContent("This is some content", f),
				PlanOnly:                 true,
			},
			{
				ExternalProviders: providerVersion233(),
				Config:            testAccConfigLocalSensitiveFileEncodedBase64Content("VGhpcyBpcyBzb21lIGJhc2U2NCBjb250ZW50", f),
				Check:             checkFileCreation("local_sensitive_file_resource.test", f),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigLocalSensitiveFileEncodedBase64Content("VGhpcyBpcyBzb21lIGJhc2U2NCBjb250ZW50", f),
				PlanOnly:                 true,
			},
			{
				ExternalProviders: providerVersion233(),
				Config:            testAccConfigLocalSensitiveFileDecodedBase64Content("This is some base64 content", f),
				Check:             checkFileCreation("local_sensitive_file_resource.test", f),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigLocalSensitiveFileDecodedBase64Content("This is some base64 content", f),
				PlanOnly:                 true,
			},
		},
		CheckDestroy: checkFileDeleted(f),
	})
}

func TestLocalSensitiveFile_Source_Upgrade(t *testing.T) {
	// create a local file that will be used as the "source" file
	if err := os.WriteFile("./testdata/source_file", []byte("sourceContent"), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("./testdata/source_file")

	r.Test(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: providerVersion233(),
				Config: `
					resource "local_sensitive_file" "file" {
					  source = "./testdata/source_file"
					  filename = "./testdata/new_file"
					}
				`,
				Check: checkFileCreation("local_sensitive_file_resource.test", "./testdata/new_file"),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config: `
					resource "local_sensitive_file" "file" {
					  source = "./testdata/source_file"
					  filename = "./testdata/new_file"
					}
				`,
				PlanOnly: true,
			},
		},
		CheckDestroy: checkFileDeleted("new_file"),
	})
}

func TestLocalSensitiveFile_Permissions_Upgrade(t *testing.T) {
	destinationDirPath := t.TempDir()
	destinationFilePath := filepath.Join(destinationDirPath, "local_sensitive_file")
	destinationFilePath = strings.ReplaceAll(destinationFilePath, `\`, `\\`)
	isDirExist := false

	r.Test(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: providerVersion233(),
				SkipFunc:          skipTestsWindows(),
				PreConfig:         checkDirExists(destinationDirPath, &isDirExist),
				Config: fmt.Sprintf(`
					resource "local_sensitive_file" "file" {
						content              = "This is some content"
						filename             = "%s"
						file_permission      = "0600"
						directory_permission = "0700"
					}`, destinationFilePath,
				),
				Check: r.ComposeTestCheckFunc(
					checkFilePermissions(destinationFilePath),
					checkDirectoryPermissions(destinationFilePath),
				),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				SkipFunc:                 skipTestsWindows(),
				Config: fmt.Sprintf(`
					resource "local_sensitive_file" "file" {
						content              = "This is some content"
						filename             = "%s"
						file_permission      = "0600"
						directory_permission = "0700"
					}`, destinationFilePath,
				),
				PlanOnly: true,
			},
		},
		ErrorCheck: func(err error) error {
			if match, _ := regexp.MatchString("Directory permission.", err.Error()); match && isDirExist {
				return nil
			}
			return err
		},
		CheckDestroy: checkFileDeleted(destinationFilePath),
	})
}

func testAccConfigLocalSensitiveSourceFile(source, filename string) string {
	return fmt.Sprintf(`
				resource "local_sensitive_file" "file" {
				  source  = %[1]q
				  filename = %[2]q
				}`, source, filename)
}

func testAccConfigLocalSensitiveFileContent(content, filename string) string {
	return fmt.Sprintf(`
				resource "local_sensitive_file" "file" {
				  content  = %[1]q
				  filename = %[2]q
				}`, content, filename)
}

func testAccConfigLocalSensitiveFileEncodedBase64Content(content, filename string) string {
	return fmt.Sprintf(`
				resource "local_sensitive_file" "file" {
				  content_base64  = %[1]q
				  filename = %[2]q
				}`, content, filename)
}

func testAccConfigLocalSensitiveFileDecodedBase64Content(content, filename string) string {
	return fmt.Sprintf(`
				resource "local_sensitive_file" "file" {
				  content_base64  = base64encode(%[1]q)
				  filename = %[2]q
				}`, content, filename)
}
