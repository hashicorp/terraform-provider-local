package local

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	r "github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"regexp"
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

func TestLocalFile_contentConfigThrowsError(t *testing.T) {
	configs := []string{`resource "local_file" "file" {
         content     = "This is some content"
         sensitive_content     = "This is some sensitive content"
         filename    = "local_file"
      }`, `resource "local_file" "file" {
         filename    = "local_file"
      }`,
	}
	for _, config := range configs {
		r.UnitTest(t, r.TestCase{
			Providers: testProviders,
			Steps: []r.TestStep{
				{
					Config:      config,
					ExpectError: regexp.MustCompile(regexp.QuoteMeta("Exactly one of `content` or `sensitive_content` must be specified")),
				},
			},
		})
	}
}
