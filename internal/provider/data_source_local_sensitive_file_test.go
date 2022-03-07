package provider

import (
	"encoding/base64"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestLocalFileSensitiveDataSource(t *testing.T) {
	content := "This is some content"

	config := `
		data "local_sensitive_file" "file" {
		  filename = "./testdata/local_file"
		}

		output "sensitive_file_content" {
			value = data.local_sensitive_file.file.content
			sensitive = true
		}

		output "sensitive_file_content_base64" {
			value = data.local_sensitive_file.file.content_base64
			sensitive = true
		}
	`

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content", content),
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content_base64", base64.StdEncoding.EncodeToString([]byte(content))),
					resource.TestMatchOutput("sensitive_file_content", regexp.MustCompile(content)),
				),
			},
		},
	})
}
