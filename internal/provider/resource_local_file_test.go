package provider

import (
	"errors"
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
	td := testTempDir(t)
	defer os.RemoveAll(td)

	f := filepath.Join(td, "local_file")
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
}`, f),
		},
		{
			f,
			"This is some sensitive content", fmt.Sprintf(`
resource "local_file" "file" {
  sensitive_content = "This is some sensitive content"
  filename = "%s"
}`, f),
		},
		{
			f,
			"This is some sensitive content", fmt.Sprintf(`
resource "local_file" "file" {
  content_base64 = "VGhpcyBpcyBzb21lIHNlbnNpdGl2ZSBjb250ZW50"
  filename = "%s"
}`, f),
		},
		{
			f,
			"This is some sensitive content", fmt.Sprintf(`
resource "local_file" "file" {
  content_base64 = base64encode("This is some sensitive content")
  filename = "%s"
}`, f),
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
				CheckDestroy: func(*terraform.State) error {
					if _, err := os.Stat(tt.path); os.IsNotExist(err) {
						return nil
					}
					return errors.New("local_file did not get destroyed")
				},
			})
		})
	}
}

func TestLocalFile_source(t *testing.T) {
	tmp, err := ioutil.TempDir("", "local-file-source")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)

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
}`

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
		CheckDestroy: func(*terraform.State) error {
			if _, err := os.Stat("new_file"); os.IsNotExist(err) {
				return nil
			}
			return errors.New("local_file did not get destroyed")
		},
	})
}

func TestLocalFile_Permissions(t *testing.T) {
	td := testTempDir(t)
	defer os.RemoveAll(td)

	destinationDirPath := td
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
}`, destinationFilePath)

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
		CheckDestroy: func(*terraform.State) error {
			if _, err := os.Stat(destinationFilePath); os.IsNotExist(err) {
				return nil
			}
			return errors.New("local_file did not get destroyed")
		},
	})

	defer os.Remove(destinationDirPath)

}

func testTempDir(t *testing.T) string {
	tmp, err := ioutil.TempDir("", "tf")
	if err != nil {
		t.Fatal(err)
	}
	return tmp
}
