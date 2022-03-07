package provider

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
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
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.local_file.file", "content", content),
					resource.TestCheckResourceAttr("data.local_file.file", "content_base64", base64.StdEncoding.EncodeToString([]byte(content))),
				),
			},
		},
	})
}
