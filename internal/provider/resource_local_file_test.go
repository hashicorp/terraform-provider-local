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

func TestLocalFile_Basic(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")

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

func TestLocalFile_source(t *testing.T) {
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
					resource "local_file" "file" {
					  source = "source_file"
					  filename = "new_file"
					}
				`,
				Check: checkFileCreation("local_file_resource.test", "new_file"),
			},
		},
		CheckDestroy: checkFileDeleted("new_file"),
	})
}

func TestLocalFile_Permissions(t *testing.T) {
	destinationDirPath := t.TempDir()
	destinationFilePath := filepath.Join(destinationDirPath, "local_file")
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
					resource "local_file" "file" {
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

func TestLocalFile_Validators(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")

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
