package local

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestLocalDirectoryDataSource(t *testing.T) {
	var tests = []struct {
		directory_permission int
		config               string
	}{
		{
			0750,
			`
				resource "local_directory" "directory" {
					directory            = "local_directory"
					directory_permission = 0750
				}
				data "local_directory" "directory" {
					directory            = "${local_directory.directory.directory}"
				}
			`,
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
							m := s.RootModule()
							i := m.Resources["data.local_directory.directory"].Primary
							got, _ := strconv.Atoi(i.Attributes["directory_permission"])
							want := test.directory_permission
							if got != want {
								return fmt.Errorf("wrong directory_permission got %d; want %d", got, want)
							}
							return nil
						},
					},
				},
			})
		})
	}
}
