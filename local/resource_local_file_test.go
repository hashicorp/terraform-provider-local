package local

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestLocalFile_Basic(t *testing.T) {
	var cases = []struct {
		path    string
		mode    os.FileMode
		content string
		config  string
	}{
		{
			"local_file",
			0777,
			"This is some content",
			`resource "local_file" "file" {
         content     = "This is some content"
         filename    = "local_file"
      }`,
		},
		{
			"other_local_file",
			0644,
			"some private content",
			`resource "local_file" "file2" {
         content  = "some private content"
         filename = "other_local_file"
         mode     = 0644
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

						fi, err := os.Stat(tt.path)
						if err != nil {
							return fmt.Errorf("config:\n%s\n,got: %s\n", tt.config, err)
						}
						mode := fi.Mode().Perm()
						if mode != tt.mode {
							return fmt.Errorf("config:\n%s\n, got: %s, want: %s\n", tt.config, mode, tt.mode)
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
