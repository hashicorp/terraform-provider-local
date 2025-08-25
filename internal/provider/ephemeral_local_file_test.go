// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	r "github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/knownvalue"
	"github.com/hashicorp/terraform-plugin-testing/statecheck"
	"github.com/hashicorp/terraform-plugin-testing/tfjsonpath"
)

func TestEphemeralLocalFile_Basic_FileContent(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccConfigEphemeralLocalFileContent("This is some content", f),
				ConfigStateChecks: []statecheck.StateCheck{
					// statecheck.ExpectKnownValue("echo.local_file", tfjsonpath.New("data").AtMapKey("value").AtMapKey("content"), knownvalue.StringExact("This is some content")),
					statecheck.ExpectKnownValue("echo.local_file", tfjsonpath.New("data").AtMapKey("content"), knownvalue.StringExact("This is some content")),
				},
				Check: checkFileDeleted(f),
			},
			// {
			// 	Config: testAccConfigEphemeralLocalFileEncodedBase64Content("VGhpcyBpcyBzb21lIG1vcmUgY29udGVudAo=", f),
			// 	ConfigStateChecks: []statecheck.StateCheck{
			// 		statecheck.ExpectKnownValue("echo.local_file", tfjsonpath.New("data").AtMapKey("content"), knownvalue.StringExact("This is some more content")),
			// 	},
			// 	Check: checkFileDeleted(f),
			// },
			// {
			// 	Config: testAccConfigEphemeralLocalFileDecodedBase64Content("This is some base64 content", f),
			// 	Check:  checkFileCreation("local_file_resource.test", f),
			// },
		},
		CheckDestroy: checkFileDeleted(f),
	})
}

func TestEphemeralLocalFile_Basic_EncodedBase64Content(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccConfigEphemeralLocalFileEncodedBase64Content("VGhpcyBpcyBzb21lIG1vcmUgY29udGVudAo=", f),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.local_file", tfjsonpath.New("data").AtMapKey("content_base64"), knownvalue.StringExact("VGhpcyBpcyBzb21lIG1vcmUgY29udGVudAo=")),
				},
				Check: checkFileDeleted(f),
			},
		},
		CheckDestroy: checkFileDeleted(f),
	})
}

// func testAccConfigEphemeralLocalSourceFile(source, filename string) string {
// 	return fmt.Sprintf(`
// 				ephemeral "local_file" "file" {
// 				  source  = %[1]q
// 				  filename = %[2]q
// 				}`, source, filename)
// }

func testAccConfigEphemeralLocalFileContent(content, filename string) string {
	return fmt.Sprintf(`
ephemeral "local_file" "file" {
	content  = %[1]q
	filename = %[2]q
}

provider "echo" {
	data = ephemeral.local_file.file
}

resource "echo" "local_file" {}
`, content, filename)
}

func testAccConfigEphemeralLocalFileEncodedBase64Content(content, filename string) string {
	return fmt.Sprintf(`
ephemeral "local_file" "file" {
	content_base64 = %[1]q
	filename       = %[2]q
}

provider "echo" {
	data = ephemeral.local_file.file
}

resource "echo" "local_file" {}
`, content, filename)
}

// func testAccConfigEphemeralLocalFileDecodedBase64Content(content, filename string) string {
// 	return fmt.Sprintf(`
// 				ephemeral "local_file" "file" {
// 				  content_base64  = base64encode(%[1]q)
// 				  filename = %[2]q
// 				}`, content, filename)
// }
