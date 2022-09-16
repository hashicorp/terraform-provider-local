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

func TestLocalFile_Basic(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	var cases = []struct {
		path    string
		content string
		config  string
	}{
		{
			f,
			"This is some content", fmt.Sprintf(`
				resource "local_file" "file" {
				  content  = "This is some content"
				  filename = "%s"
				}`, f,
			),
		},
		{
			f,
			"This is some sensitive content", fmt.Sprintf(`
				resource "local_file" "file" {
				  sensitive_content = "This is some sensitive content"
				  filename = "%s"
				}`, f,
			),
		},
		{
			f,
			"This is some base64 content", fmt.Sprintf(`
				resource "local_file" "file" {
				  content_base64 = "VGhpcyBpcyBzb21lIGJhc2U2NCBjb250ZW50"
				  filename = "%s"
				}`, f,
			),
		},
		{
			f,
			"This is some base64 content", fmt.Sprintf(`
				resource "local_file" "file" {
				  content_base64 = base64encode("This is some base64 content")
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

func TestLocalFile_source(t *testing.T) {
	// create a local file that will be used as the "source" file
	source_content := "local file content"
	if err := ioutil.WriteFile("source_file", []byte(source_content), 0644); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("source_file")

	config := `
		resource "local_file" "file" {
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

func TestLocalFile_Permissions(t *testing.T) {
	destinationDirPath := t.TempDir()
	destinationFilePath := filepath.Join(destinationDirPath, "local_file")
	destinationFilePath = strings.ReplaceAll(destinationFilePath, `\`, `\\`)
	filePermission := os.FileMode(0600)
	directoryPermission := os.FileMode(0700)
	skipDirCheck := false
	config := fmt.Sprintf(`
		resource "local_file" "file" {
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

func TestLocalFile_checksums(t *testing.T) {
	content := "This is some content"
	filename := filepath.Join(t.TempDir(), "local_file")
	filename = strings.ReplaceAll(filename, `\`, `\\`)

	config := fmt.Sprintf(`
	resource "local_file" "file" {
	  content  = "%s"
	  filename = "%s"
	}`, content, filename)

	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: config,
				Check: r.ComposeAggregateTestCheckFunc(
					r.TestCheckResourceAttr("local_file.file", "content_md5", "ee428920507e39e8d89c2cabe6641b67"),
					r.TestCheckResourceAttr("local_file.file", "content_sha1", "f3705a38abd5d2bd1f4fecda606d216216c536b1"),
					r.TestCheckResourceAttr("local_file.file", "content_sha256", "d68e560efbe6f20c31504b2fc1c6d3afa1f58b8ee293ad3311939a5fd5059a12"),
					r.TestCheckResourceAttr("local_file.file", "content_base64sha256", "1o5WDvvm8gwxUEsvwcbTr6H1i47ik60zEZOaX9UFmhI="),
					r.TestCheckResourceAttr("local_file.file", "content_sha512", "217150cec0dac8ba2d640eeb80f12407c3b9362650e716bc568fcd2cca0fd951db25fd4aa0aefa6454803697ecb74fd3dc8b36bd2c2e5a3a3ac2456e3017728d"),
					r.TestCheckResourceAttr("local_file.file", "content_base64sha512", "IXFQzsDayLotZA7rgPEkB8O5NiZQ5xa8Vo/NLMoP2VHbJf1KoK76ZFSANpfst0/T3Is2vSwuWjo6wkVuMBdyjQ=="),
				),
			},
		},
		CheckDestroy: checkFileDeleted(filename),
	})
}
