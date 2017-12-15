package local

import (
	"fmt"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestLocalFileDataSource(t *testing.T) {
	var tests = []struct {
		content string
		config  string
	}{
		{
			"This is some content",
			`
				resource "local_file" "file" {
					content  = "This is some content"
					filename = "local_file"
				}
				data "local_file" "file" {
					filename = "${local_file.file.filename}"
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
							i := m.Resources["data.local_file.file"].Primary
							if got, want := i.Attributes["content"], test.content; got != want {
								return fmt.Errorf("wrong content %q; want %q", got, want)
							}
							return nil
						},
					},
				},
			})
		})
	}
}
