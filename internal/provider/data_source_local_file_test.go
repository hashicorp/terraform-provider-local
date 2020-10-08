package provider

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
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

	for i, test := range tests {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
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
							if got, want := i.Attributes["content_base64"], base64.StdEncoding.EncodeToString([]byte(test.content)); got != want {
								return fmt.Errorf("wrong content_base64 %q; want %q", got, want)
							}
							return nil
						},
					},
				},
			})
		})
	}
}
