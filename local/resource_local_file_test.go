package local

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"path"
	"runtime"
)

func TestLocalFile_Basic(t *testing.T) {
	var cases = []struct {
		path    string
		content string
		config  string
	}{
		{
			"local_file",
			"This is some content",
			`resource "local_file" "file" {
         content     = "This is some content"
         filename    = "local_file"
      }`,
		},
		{
			"local_file",
			"This is some sensitive content",
			`resource "local_file" "file" {
         sensitive_content     = "This is some sensitive content"
         filename    = "local_file"
      }`,
		},
	}

	for _, tt := range cases {
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
	}
}

func TestLocalFile_Permissions(t *testing.T) {
	destination := "local_directory/local_file"
	filePermission := os.FileMode(0600)
	directoryPermission := os.FileMode(0700)
	skipDirCheck := false
	config := `resource "local_file" "file" {
		content              = "This is some content"
		filename             = "local_directory/local_file"
		file_permission      = 0600
		directory_permission = 0700
	}`

	r.UnitTest(t, r.TestCase{
		Providers: testProviders,
		Steps: []r.TestStep{
			{
				Config: config,
				PreConfig: func() {
					// if directory already existed prior to check, skip check
					if _, err := os.Stat(path.Dir(destination)); !os.IsNotExist(err) {
						skipDirCheck = true
					}
				},
				Check: func(s *terraform.State) error {
					if runtime.GOOS == "windows" {
						// skip all checks if windows
						return nil
					}

					fileInfo, err := os.Stat(destination)
					if err != nil {
						return fmt.Errorf("config:\n%s\ngot:%s\n", config, err)
					}

					if fileInfo.Mode() != filePermission {
						return fmt.Errorf(
							"File permission.\nconfig:\n%s\nexpected:%s\ngot: %s\n",
							config, filePermission, fileInfo.Mode())
					}

					if !skipDirCheck {
						dirInfo, _ := os.Stat(path.Dir(destination))
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
			if _, err := os.Stat(destination); os.IsNotExist(err) {
				return nil
			}
			return errors.New("local_file did not get destroyed")
		},
	})
}
