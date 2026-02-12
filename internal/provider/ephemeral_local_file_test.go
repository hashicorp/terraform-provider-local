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
	"github.com/hashicorp/terraform-plugin-testing/tfversion"
)

func TestEphemeralLocalFile_Basic_SourceContent(t *testing.T) {
	sourceDirPath := t.TempDir()
	sourceFilePath := filepath.Join(sourceDirPath, "source_file")
	sourceFilePath = strings.ReplaceAll(sourceFilePath, `\`, `\\`)
	// create a local file that will be used as the "source" file
	if err := createSourceFile(sourceFilePath, "local file content"); err != nil {
		t.Fatal(err)
	}

	destinationDirPath := t.TempDir()
	destinationFilePath := filepath.Join(destinationDirPath, "new_file")
	destinationFilePath = strings.ReplaceAll(destinationFilePath, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccConfigEphemeralLocalSourceFile(sourceFilePath, destinationFilePath),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.local_file", tfjsonpath.New("data").AtMapKey("source"), knownvalue.StringExact(sourceFilePath)),
					// statecheck.ExpectKnownOutputValue("local_file", knownvalue.StringExact("local file content")),
				},
				Check: checkFileDeleted(destinationFilePath),
			},
		},
		CheckDestroy: checkFileDeleted(destinationFilePath),
	})
}

func TestEphemeralLocalFile_Basic_FileContent(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccConfigEphemeralLocalFileContent("This is some content", f),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.local_file", tfjsonpath.New("data").AtMapKey("content"), knownvalue.StringExact("This is some content")),
				},
				Check: checkFileDeleted(f),
			},
		},
		CheckDestroy: checkFileDeleted(f),
	})
}

func TestEphemeralLocalFile_Basic_EncodedBase64Content(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
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

func TestEphemeralLocalFile_Basic_DecodedBase64Content(t *testing.T) {
	f := filepath.Join(t.TempDir(), "local_file")
	f = strings.ReplaceAll(f, `\`, `\\`)

	r.UnitTest(t, r.TestCase{
		TerraformVersionChecks: []tfversion.TerraformVersionCheck{
			tfversion.SkipBelow(tfversion.Version1_10_0),
		},
		ProtoV5ProviderFactories: protoV5ProviderFactories(),
		ProtoV6ProviderFactories: protoV6ProviderFactories(),
		Steps: []r.TestStep{
			{
				Config: testAccConfigEphemeralLocalFileDecodedBase64Content("This is some base64 content", f),
				ConfigStateChecks: []statecheck.StateCheck{
					statecheck.ExpectKnownValue("echo.local_file", tfjsonpath.New("data").AtMapKey("content_base64"), knownvalue.StringExact("VGhpcyBpcyBzb21lIGJhc2U2NCBjb250ZW50")),
				},
				Check: checkFileDeleted(f),
			},
		},
		CheckDestroy: checkFileDeleted(f),
	})
}

func testAccConfigEphemeralLocalSourceFile(source, filename string) string {
	return fmt.Sprintf(`
ephemeral "local_file" "file" {
	source   = %[1]q
	filename = %[2]q
}

provider "echo" {
	data = ephemeral.local_file.file
}

resource "echo" "local_file" {}
`, source, filename)
}

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

func testAccConfigEphemeralLocalFileDecodedBase64Content(content, filename string) string {
	return fmt.Sprintf(`
ephemeral "local_file" "file" {
	content_base64 = base64encode(%[1]q)
	filename       = %[2]q
}

provider "echo" {
	data = ephemeral.local_file.file
}

resource "echo" "local_file" {}
`, content, filename)
}
