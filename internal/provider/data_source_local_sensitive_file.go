package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLocalSensitiveFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocalSensitiveFileRead,

		Description: "Reads a file that contains sensitive data, from the local filesystem.",

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
				Sensitive:   true,
				Computed:    true,
			},
			"content_base64": {
				Type:        schema.TypeString,
				Description: "Base64 encoded version of the file raw content (use this when dealing with binary data).",
				Sensitive:   true,
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

func dataSourceLocalSensitiveFileRead(d *schema.ResourceData, m interface{}) error {
	// NOTE: We can use the read-method for the data source `local_file` as-is, because
	// all this data source does, is adding "Sensitive: true" to the schema of the property.
	//
	// The values and the property names are meant to be kept the same between data sources.
	return dataSourceLocalFileRead(d, m)
}
