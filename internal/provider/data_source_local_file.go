package provider

import (
	"encoding/base64"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLocalFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocalFileRead,

		Description: "Reads a file from the local filesystem.",

		Schema: map[string]*schema.Schema{
			"filename": {
				Type:        schema.TypeString,
				Description: "Path to the file that will be read. The data source will return an error if the file does not exist.",
				Required:    true,
				ForceNew:    true,
			},
			"content": {
				Type:        schema.TypeString,
				Description: "Raw content of the file that was read, as UTF-8 encoded string.",
				Computed:    true,
			},
			"content_base64": {
				Type:        schema.TypeString,
				Description: "Base64 encoded version of the file content (use this when dealing with binary data).",
				Computed:    true,
			},
			"content_md5": {
				Type:        schema.TypeString,
				Description: "MD5 checksum of file content.",
				Computed:    true,
			},
			"content_sha1": {
				Type:        schema.TypeString,
				Description: "SHA1 checksum of file content.",
				Computed:    true,
			},
			"content_sha256": {
				Type:        schema.TypeString,
				Description: "SHA256 checksum of file content.",
				Computed:    true,
			},
			"content_base64sha256": {
				Type:        schema.TypeString,
				Description: "Base64 encoded SHA256 checksum of file content.",
				Computed:    true,
			},
			"content_sha512": {
				Type:        schema.TypeString,
				Description: "SHA512 checksum of file content.",
				Computed:    true,
			},
			"content_base64sha512": {
				Type:        schema.TypeString,
				Description: "Base64 encoded SHA512 checksum of file content.",
				Computed:    true,
			},
		},
	}
}

func dataSourceLocalFileRead(d *schema.ResourceData, _ interface{}) error {
	// Read the entire file content
	path := d.Get("filename").(string)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	// Set the content both as UTF-8 string, and as base64 encoded string
	d.Set("content", string(content))
	d.Set("content_base64", base64.StdEncoding.EncodeToString(content))

	checksums := genFileChecksums(content)
	d.Set("content_md5", checksums.md5Hex)
	d.Set("content_sha1", checksums.sha1Hex)
	d.Set("content_sha256", checksums.sha256Hex)
	d.Set("content_base64sha256", checksums.sha256Base64)
	d.Set("content_sha512", checksums.sha512Hex)
	d.Set("content_base64sha512", checksums.sha512Base64)

	// Use the hexadecimal encoding of the checksum of the file content as ID
	d.SetId(checksums.sha1Hex)

	return nil
}
