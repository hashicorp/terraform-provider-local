package provider

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLocalFile() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocalFileRead,

		Schema: map[string]*schema.Schema{
			"filename": {
				Type:        schema.TypeString,
				Description: "Path to the output file",
				Required:    true,
				ForceNew:    true,
			},
			"sensitive": {
				Type:        schema.TypeBool,
				Description: "Treats content as sensitive and will not output it in plan",
				Optional:    true,
			},
			"content": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"sensitive_content": {
				Type:      schema.TypeString,
				Sensitive: true,
				Computed:  true,
			},
			"content_base64": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceLocalFileRead(d *schema.ResourceData, _ interface{}) error {
	path := d.Get("filename").(string)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}
	sensitive := d.Get("sensitive").(bool)
	if sensitive {
		d.Set("sensitive_content", string(content))
	} else {
		d.Set("content", string(content))
	}
	d.Set("content_base64", base64.StdEncoding.EncodeToString(content))

	checksum := sha1.Sum([]byte(content))
	d.SetId(hex.EncodeToString(checksum[:]))

	return nil
}
