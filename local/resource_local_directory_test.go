package local

import (
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestLocalDirectory_Basic(t *testing.T) {
	var tests = []struct {
		directory            string
		directory_permission int
		config               string
	}{
		{
			"local_directory",
			0750,
			`resource "local_directory" "directory" {
         directory            = "local_directory"
         directory_permission = 0750
      }`,
		},
	}

	for _, test := range tests {
		t.Run("", func(t *testing.T) {
			resource.UnitTest(t, resource.TestCase{
				Providers: testProviders,
				Steps: []resource.TestStep{
					{
						Config: test.config,
						Check: func(s *terraform.State) error {
							dirInfo, err := os.Stat(test.directory)
							if err != nil {
								return fmt.Errorf("config:\n%s\n,got: %s\n", test.config, err)
							}
							if int(dirInfo.Mode().Perm()) != test.directory_permission {
								return fmt.Errorf("config:\n%s\ngot:\n%d\nwant:\n%d\n", test.config, int(dirInfo.Mode().Perm()), test.directory_permission)
							}
							return nil
						},
					},
				},
				CheckDestroy: func(*terraform.State) error {
					if _, err := os.Stat(test.directory); os.IsNotExist(err) {
						return nil
					}
					return errors.New("local_directory did not get destroyed")
				},
			})
		})
	}
}
