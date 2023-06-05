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

func TestLocalFile_Basic(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccConfigLocalFileContent("This is some content", f),
				Check:  checkFileCreation("local_file_resource.test", f),
			},
			{
				Config: testAccConfigLocalFileSensitiveContent("This is some sensitive content", f),
				Check:  checkFileCreation("local_file_resource.test", f),
			},
			{
				Config: testAccConfigLocalFileEncodedBase64Content("VGhpcyBpcyBzb21lIGJhc2U2NCBjb250ZW50", f),
				Check:  checkFileCreation("local_file_resource.test", f),
			},
			{
				Config: testAccConfigLocalFileDecodedBase64Content("This is some base64 content", f),
				Check:  checkFileCreation("local_file_resource.test", f),
			},
		},
		CheckDestroy: checkFileDeleted(f),
	})
}

func TestLocalFile_Source(t *testing.T) {
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
				Config: testAccConfigLocalSourceFile(sourceFilePath, destinationFilePath),
				Check:  checkFileCreation("local_file_resource.test", destinationFilePath),
			},
		},
		CheckDestroy: checkFileDeleted(destinationFilePath),
	})
}

func TestLocalFile_Permissions(t *testing.T) {
	destinationDirPath := t.TempDir()
	destinationFilePath := filepath.Join(destinationDirPath, "local_file")
	destinationFilePath = strings.ReplaceAll(destinationFilePath, `\`, `\\`)
	isDirExist := false

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				PreConfig: checkDirExists(destinationDirPath, &isDirExist),
				SkipFunc:  skipTestsWindows(),
				Config: fmt.Sprintf(`
					resource "local_file" "file" {
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

func TestLocalFile_Validators(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		CheckDestroy:             nil,
		Steps: []r.TestStep{
			{
				Config: fmt.Sprintf(`
				resource "local_file" "file" {
				  filename = "%s"
				}`, f),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
			{
				Config: fmt.Sprintf(`
				resource "local_file" "file" {
                  content = "content"
				  sensitive_content = "sensitive_content"
				  filename = "%s"
				}`, f),
				ExpectError: regexp.MustCompile(`.*Error: Invalid Attribute Combination`),
			},
		},
	})
}

func TestLocalFile_Upgrade(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.Test(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: providerVersion233(),
				Config:            testAccConfigLocalFileContent("This is some content", f),
				Check:             checkFileCreation("local_file_resource.test", f),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigLocalFileContent("This is some content", f),
				PlanOnly:                 true,
			},
			{
				ExternalProviders: providerVersion233(),
				Config:            testAccConfigLocalFileSensitiveContent("This is some sensitive content", f),
				Check:             checkFileCreation("local_file_resource.test", f),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigLocalFileSensitiveContent("This is some sensitive content", f),
				PlanOnly:                 true,
			},
			{
				ExternalProviders: providerVersion233(),
				Config:            testAccConfigLocalFileEncodedBase64Content("VGhpcyBpcyBzb21lIGJhc2U2NCBjb250ZW50", f),
				Check:             checkFileCreation("local_file_resource.test", f),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigLocalFileEncodedBase64Content("VGhpcyBpcyBzb21lIGJhc2U2NCBjb250ZW50", f),
				PlanOnly:                 true,
			},
			{
				ExternalProviders: providerVersion233(),
				Config:            testAccConfigLocalFileDecodedBase64Content("This is some base64 content", f),
				Check:             checkFileCreation("local_file_resource.test", f),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config:                   testAccConfigLocalFileDecodedBase64Content("This is some base64 content", f),
				PlanOnly:                 true,
			},
		},
		CheckDestroy: checkFileDeleted(f),
	})
}

func TestLocalFile_Source_Upgrade(t *testing.T) {
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
					resource "local_file" "file" {
					  source = "./testdata/source_file"
					  filename = "./testdata/new_file"
					}
				`,
				Check: checkFileCreation("local_file_resource.test", "./testdata/new_file"),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				Config: `
					resource "local_file" "file" {
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

func TestLocalFile_Permissions_Upgrade(t *testing.T) {
	destinationDirPath := t.TempDir()
	destinationFilePath := filepath.Join(destinationDirPath, "local_file")
	destinationFilePath = strings.ReplaceAll(destinationFilePath, `\`, `\\`)
	isDirExist := false

	r.Test(t, r.TestCase{
		Steps: []r.TestStep{
			{
				ExternalProviders: providerVersion233(),
				SkipFunc:          skipTestsWindows(),
				PreConfig:         checkDirExists(destinationDirPath, &isDirExist),
				Config: fmt.Sprintf(`
					resource "local_file" "file" {
						content              = "This is some content"
						filename             = "%s"
						file_permission      = "0600"
						directory_permission = "0700"
					}`, destinationFilePath,
				),
				Check: r.ComposeTestCheckFunc(
					checkFilePermissions(destinationFilePath),
					checkDirectoryPermissions(destinationFilePath)),
			},
			{
				ProtoV5ProviderFactories: protoV5ProviderFactories(),
				SkipFunc:                 skipTestsWindows(),
				Config: fmt.Sprintf(`
					resource "local_file" "file" {
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

func testAccConfigLocalSourceFile(source, filename string) string {
	return fmt.Sprintf(`
				resource "local_file" "file" {
				  source  = %[1]q
				  filename = %[2]q
				}`, source, filename)
}

func testAccConfigLocalFileContent(content, filename string) string {
	return fmt.Sprintf(`
				resource "local_file" "file" {
				  content  = %[1]q
				  filename = %[2]q
				}`, content, filename)
}

func testAccConfigLocalFileSensitiveContent(content, filename string) string {
	return fmt.Sprintf(`
				resource "local_file" "file" {
				  sensitive_content  = %[1]q
				  filename = %[2]q
				}`, content, filename)
}

func testAccConfigLocalFileEncodedBase64Content(content, filename string) string {
	return fmt.Sprintf(`
				resource "local_file" "file" {
				  content_base64  = %[1]q
				  filename = %[2]q
				}`, content, filename)
}

func testAccConfigLocalFileDecodedBase64Content(content, filename string) string {
	return fmt.Sprintf(`
				resource "local_file" "file" {
				  content_base64  = base64encode(%[1]q)
				  filename = %[2]q
				}`, content, filename)
}
