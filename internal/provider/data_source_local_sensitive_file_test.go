// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

func TestLocalFileSensitiveDataSource(t *testing.T) {
	testFileContent := "This is some content"
	checkSums := genFileChecksums([]byte(testFileContent))

	config := `
		data "local_sensitive_file" "file" {
		  filename = "./testdata/local_file"
		}
	`

	resource.UnitTest(t, resource.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content", testFileContent),
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content_base64", base64.StdEncoding.EncodeToString([]byte(testFileContent))),
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content_md5", checkSums.md5Hex),
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content_sha1", checkSums.sha1Hex),
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content_sha256", checkSums.sha256Hex),
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content_base64sha256", checkSums.sha256Base64),
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content_sha512", checkSums.sha512Hex),
					resource.TestCheckResourceAttr("data.local_sensitive_file.file", "content_base64sha512", checkSums.sha512Base64),
				),
			},
		},
	})
}

func TestLocalFileSensitiveDataSourceCheckSensitiveAttributes(t *testing.T) {
	dataSource := NewLocalSensitiveFileDataSourceWithSchema()
	schemaResponse := datasource.SchemaResponse{}

	dataSource.Schema(context.Background(), datasource.SchemaRequest{}, &schemaResponse)
	if !schemaResponse.Schema.Attributes["content"].IsSensitive() {
		t.Errorf("attribute 'content' should be marked as 'Sensitive'")
	}

	if !schemaResponse.Schema.Attributes["content_base64"].IsSensitive() {
		t.Errorf("attribute 'content_base64' should be marked as 'Sensitive'")
	}
}
