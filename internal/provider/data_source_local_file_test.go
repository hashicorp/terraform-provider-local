// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLocalFileDataSource(t *testing.T) {
	content := "This is some content"
	checkSums := genFileChecksums([]byte(content))

	config := `
		data "local_file" "file" {
		  filename = "./testdata/local_file"
		}
	`

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.local_file.file", "content", content),
					resource.TestCheckResourceAttr("data.local_file.file", "content_base64", base64.StdEncoding.EncodeToString([]byte(content))),
					resource.TestCheckResourceAttr("data.local_file.file", "content_md5", checkSums.md5Hex),
					resource.TestCheckResourceAttr("data.local_file.file", "content_sha1", checkSums.sha1Hex),
					resource.TestCheckResourceAttr("data.local_file.file", "content_sha256", checkSums.sha256Hex),
					resource.TestCheckResourceAttr("data.local_file.file", "content_base64sha256", checkSums.sha256Base64),
					resource.TestCheckResourceAttr("data.local_file.file", "content_sha512", checkSums.sha512Hex),
					resource.TestCheckResourceAttr("data.local_file.file", "content_base64sha512", checkSums.sha512Base64),
				),
			},
		},
	})
}
