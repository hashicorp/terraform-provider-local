package provider

import (
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
					resource.TestCheckResourceAttr("data.local_file.file", "content_base64", "VGhpcyBpcyBzb21lIGNvbnRlbnQ="),
					resource.TestCheckResourceAttr("data.local_file.file", "content_md5", "ee428920507e39e8d89c2cabe6641b67"),
					resource.TestCheckResourceAttr("data.local_file.file", "content_sha1", "f3705a38abd5d2bd1f4fecda606d216216c536b1"),
					resource.TestCheckResourceAttr("data.local_file.file", "content_sha256", "d68e560efbe6f20c31504b2fc1c6d3afa1f58b8ee293ad3311939a5fd5059a12"),
					resource.TestCheckResourceAttr("data.local_file.file", "content_base64sha256", "1o5WDvvm8gwxUEsvwcbTr6H1i47ik60zEZOaX9UFmhI="),
					resource.TestCheckResourceAttr("data.local_file.file", "content_sha512", "217150cec0dac8ba2d640eeb80f12407c3b9362650e716bc568fcd2cca0fd951db25fd4aa0aefa6454803697ecb74fd3dc8b36bd2c2e5a3a3ac2456e3017728d"),
					resource.TestCheckResourceAttr("data.local_file.file", "content_base64sha512", "IXFQzsDayLotZA7rgPEkB8O5NiZQ5xa8Vo/NLMoP2VHbJf1KoK76ZFSANpfst0/T3Is2vSwuWjo6wkVuMBdyjQ=="),
				),
			},
		},
	})
}
