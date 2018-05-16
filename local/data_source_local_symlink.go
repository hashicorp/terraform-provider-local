package local

import (
	"github.com/hashicorp/terraform/helper/schema"
	"os"
)

func dataSourceLocalSymlink() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocalSymlinkRead,

		Schema: map[string]*schema.Schema{
			"symlink": {
				Type:        schema.TypeString,
				Description: "Path of the symlink",
				Required:    true,
				ForceNew:    false,
			},
			"destination": {
				Type:        schema.TypeString,
				Description: "Permissions to set for directories created",
				Optional:    true,
				ForceNew:    false,
			},
		},
	}
}

func dataSourceLocalSymlinkRead(d *schema.ResourceData, _ interface{}) error {
	wantedSymlink := d.Get("symlink").(string)
	if _, err := os.Stat(wantedSymlink); err != nil {
		return err
	}
	destination, err := os.Readlink(wantedSymlink)
	if err != nil {
		return err
	}

	d.SetId(wantedSymlink)
	d.Set("destination", destination)

	return nil
}
