package provider

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestLocalFileSensitiveDataSource(t *testing.T) {
	testFileContent := "This is some content"

	config := `
		data "local_sensitive_file" "file" {
		  filename = "./testdata/local_file"
		}
	`

	resource.UnitTest(t, resource.TestCase{
		Providers: testProviders,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content", testFileContent),
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content_base64", base64.StdEncoding.EncodeToString([]byte(testFileContent))),
				),
			},
		},
	})
}

func TestLocalFileSensitiveDataSourceCheckSensitiveAttributes(t *testing.T) {
	dsSchema := dataSourceLocalSensitiveFile()

	if !dsSchema.Schema["content"].Sensitive {
		t.Errorf("attribute 'content' should be marked as 'Sensitive'")
	}

	if !dsSchema.Schema["content_base64"].Sensitive {
		t.Errorf("attribute 'content_base64' should be marked as 'Sensitive'")
	}
}
