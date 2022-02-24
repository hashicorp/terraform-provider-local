package provider

import (
	"encoding/base64"
	"fmt"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestLocalFileDataSource(t *testing.T) {
	content := "This is some content"

	config := `
data "local_file" "file" {
  filename = "./testdata/local_file"
}
`

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					m := s.RootModule()

					i := m.Resources["data.local_file.file"].Primary

					if got, want := i.Attributes["content"], content; got != want {
						return fmt.Errorf("wrong content %q; want %q", got, want)
					}
					if got, want := i.Attributes["content_base64"], base64.StdEncoding.EncodeToString([]byte(content)); got != want {
						return fmt.Errorf("wrong content_base64 %q; want %q", got, want)
					}
					if got, want := i.Attributes["sensitive_content"], ""; got != want {
						return fmt.Errorf("Content was not marked as sensitive and should be in 'content' attribute instead")
					}
					if got, want := i.Attributes["sensitive_content_base64"], ""; got != want {
						return fmt.Errorf("Content was marked as sensitive and should be in 'content_base64' attribute instead")
					}
					return nil
				},
			},
		},
	})
}

func TestLocalFileSensitiveDataSource(t *testing.T) {
	content := "This is some content"

	config := `
	data "local_file" "file" {
	  filename = "./testdata/local_file"
	  sensitive = true
	}
	`

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: func(s *terraform.State) error {
					m := s.RootModule()

					i := m.Resources["data.local_file.file"].Primary

					if got, want := i.Attributes["content"], ""; got != want {
						return fmt.Errorf("Content was marked as sensitive and should not be returned in 'content' attribute")
					}
					if got, want := i.Attributes["content_base64"], ""; got != want {
						return fmt.Errorf("Content was marked as sensitive and should not be returned in 'content_base64' attribute")
					}
					if got, want := i.Attributes["sensitive_content"], content; got != want {
						return fmt.Errorf("wrong content %q; want %q", got, want)
					}
					if got, want := i.Attributes["sensitive_content_base64"], base64.StdEncoding.EncodeToString([]byte(content)); got != want {
						return fmt.Errorf("wrong content_base64 %q; want %q", got, want)
					}
					return nil
				},
			},
		},
	})
}
