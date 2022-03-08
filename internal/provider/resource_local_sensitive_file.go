package provider

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLocalSensitiveFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceLocalSensitiveFileCreate,
		Read:   resourceLocalSensitiveFileRead,
		Delete: resourceLocalSensitiveFileDelete,

		Description: "Generates a local file with the given sensitive content.",

		Schema: map[string]*schema.Schema{
			"content": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"content", "content_base64", "source"},
				Description:  "Sensitive content to store in the file, expected to be an UTF-8 encoded string.",
			},
			"content_base64": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"content", "content_base64", "source"},
				Description:  "Sensitive content to store in the file, expected to be binary encoded as base64 string.",
			},
			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"content", "content_base64", "source"},
				Description:  "Path to file to use as source for the one we are creating.",
			},
			"filename": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				Description: `
					The path to the file that will be created.
					Missing parent directories will be created.
					If the file already exists, it will be overridden with the given content.
				`,
			},
			"file_permission": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "0700",
				ValidateFunc: validateModePermission,
				Description:  "Permissions to set for the output file (in numeric notation).",
			},
			"directory_permission": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "0700",
				ValidateFunc: validateModePermission,
				Description:  "Permissions to set for directories created (in numeric notation).",
			},
		},
	}
}

func resourceLocalSensitiveFileRead(d *schema.ResourceData, m interface{}) error {
	return resourceLocalFileRead(d, m)
}

func resourceLocalSensitiveFileCreate(d *schema.ResourceData, m interface{}) error {
	return resourceLocalFileCreate(d, m)
}

func resourceLocalSensitiveFileDelete(d *schema.ResourceData, m interface{}) error {
	return resourceLocalFileDelete(d, m)
}
