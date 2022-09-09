package provider

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestLocalFileDataSource(t *testing.T) {
	content := "This is some content"
	md5Sum := md5.Sum([]byte(content))
	sha1Sum := sha1.Sum([]byte(content))
	sha256Sum := sha256.Sum256([]byte(content))
	sha512Sum := sha512.Sum512([]byte(content))

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
					resource.TestCheckResourceAttr("data.local_file.file", "content_md5", hex.EncodeToString(md5Sum[:])),
					resource.TestCheckResourceAttr("data.local_file.file", "content_sha1", hex.EncodeToString(sha1Sum[:])),
					resource.TestCheckResourceAttr("data.local_file.file", "content_sha256", hex.EncodeToString(sha256Sum[:])),
					resource.TestCheckResourceAttr("data.local_file.file", "content_base64sha256", base64.StdEncoding.EncodeToString(sha256Sum[:])),
					resource.TestCheckResourceAttr("data.local_file.file", "content_sha512", hex.EncodeToString(sha512Sum[:])),
					resource.TestCheckResourceAttr("data.local_file.file", "content_base64sha512", base64.StdEncoding.EncodeToString(sha512Sum[:])),
				),
			},
		},
	})
}
