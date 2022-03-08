package provider

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"testing"

	r "github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestLocalSensitiveFile_Basic(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_sensitive_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	var cases = []struct {
		path    string
		content string
		config  string
	}{
		{
			f,
			"This is some sensitive content", fmt.Sprintf(`
				resource "local_sensitive_file" "file" {
				  content  = "This is some sensitive content"
				  filename = "%s"
				}`, f,
			),
		},
		{
			f,
			"This is some sensitive base64 content", fmt.Sprintf(`
				resource "local_sensitive_file" "file" {
				  content_base64 = "VGhpcyBpcyBzb21lIHNlbnNpdGl2ZSBiYXNlNjQgY29udGVudA=="
				  filename = "%s"
				}`, f,
			),
		},
		{
			f,
			"This is some sensitive base64 content", fmt.Sprintf(`
				resource "local_sensitive_file" "file" {
				  content_base64 = base64encode("This is some sensitive base64 content")
				  filename = "%s"
				}`, f,
			),
		},
	}

	for i, tt := range cases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			r.UnitTest(t, r.TestCase{
				Providers: testProviders,
				Steps: []r.TestStep{
					{
						Config: tt.config,
						Check: func(s *terraform.State) error {
							content, err := ioutil.ReadFile(tt.path)
							if err != nil {
								return fmt.Errorf("config:\n%s\n,got: %s\n", tt.config, err)
							}
							if string(content) != tt.content {
								return fmt.Errorf("config:\n%s\ngot:\n%s\nwant:\n%s\n", tt.config, content, tt.content)
							}
							return nil
						},
					},
				},
				CheckDestroy: checkFileDeleted(tt.path),
			})
		})
	}
}

func TestLocalSensitiveFile_source(t *testing.T) {
	// create a local file that will be used as the "source" file
	source_content := "local file content"
	if err := ioutil.WriteFile("source_file", []byte(source_content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("source_file")

	config := `
		resource "local_sensitive_file" "file" {
		  source = "source_file"
		  filename = "new_file"
		}
	`

	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					content, err := ioutil.ReadFile("new_file")
					if err != nil {
						return fmt.Errorf("config:\n%s\n,got: %s\n", config, err)
					}
					if string(content) != source_content {
						return fmt.Errorf("config:\n%s\ngot:\n%s\nwant:\n%s\n", config, content, source_content)
					}
					return nil
				},
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
	skipDirCheck := false
	config := fmt.Sprintf(`
		resource "local_sensitive_file" "file" {
			content              = "This is some content"
			filename             = "%s"
			file_permission      = "0600"
			directory_permission = "0700"
		}`, destinationFilePath,
	)

	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: config,
				PreConfig: func() {
					// if directory already existed prior to check, skip check
					if _, err := os.Stat(path.Dir(destinationFilePath)); !os.IsNotExist(err) {
						skipDirCheck = true
					}
				},
				Check: func(s *terraform.State) error {
					if runtime.GOOS == "windows" {
						// skip all checks if windows
						return nil
					}

					fileInfo, err := os.Stat(destinationFilePath)
					if err != nil {
						return fmt.Errorf("config:\n%s\ngot:%s\n", config, err)
					}

					if fileInfo.Mode() != filePermission {
						return fmt.Errorf(
							"File permission.\nconfig:\n%s\nexpected:%s\ngot: %s\n",
							config, filePermission, fileInfo.Mode())
					}

					if !skipDirCheck {
						dirInfo, _ := os.Stat(path.Dir(destinationFilePath))
						// we have to use FileMode.Perm() here, otherwise directory bit causes issues
						if dirInfo.Mode().Perm() != directoryPermission.Perm() {
							return fmt.Errorf(
								"Directory permission.\nconfig:\n%s\nexpected:%s\ngot: %s\n",
								config, directoryPermission, dirInfo.Mode().Perm())
						}
					}

					return nil
				},
			},
		},
		CheckDestroy: checkFileDeleted(destinationFilePath),
	})
}
