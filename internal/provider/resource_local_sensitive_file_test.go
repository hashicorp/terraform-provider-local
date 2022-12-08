package provider

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
	// create a local file that will be used as the "source" file
	if err := createSourceFile("local file content"); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("source_file")

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: `
					resource "local_sensitive_file" "file" {
					  source = "source_file"
					  filename = "new_file"
					}
				`,
				Check: checkFileCreation("local_sensitive_file_resource.test", "new_file"),
			},
		},
		CheckDestroy: checkFileDeleted("new_file"),
	})
}

func TestLocalSensitiveFile_Permissions(t *testing.T) {
	destinationDirPath := t.TempDir()
	destinationFilePath := filepath.Join(destinationDirPath, "local_sensitive_file")
	destinationFilePath = strings.ReplaceAll(destinationFilePath, `\`, `\\`)
	filePermission := os.FileMode(0600)
	directoryPermission := os.FileMode(0700)
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
					checkFilePermissions(destinationFilePath, filePermission),
					checkDirectoryPermissions(destinationFilePath, directoryPermission)),
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
