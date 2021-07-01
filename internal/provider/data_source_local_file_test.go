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
					return nil
				},
			},
		},
	})
}
