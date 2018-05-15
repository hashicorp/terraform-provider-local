package local

import (
	"github.com/hashicorp/terraform/helper/schema"
	"os"
)

func dataSourceLocalDirectory() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocalDirectoryRead,

		Schema: map[string]*schema.Schema{
			"directory": {
				Type:        schema.TypeString,
				Description: "Path to the directory",
				Required:    true,
				ForceNew:    true,
			},
			"directory_permission": {
				Type:        schema.TypeInt,
				Description: "Permissions to set for directories created",
				Optional:    true,
				ForceNew:    true,
				Default:     0777,
			},
		},
	}
}

func dataSourceLocalDirectoryRead(d *schema.ResourceData, _ interface{}) error {
	directory := d.Get("directory").(string)
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		return err
	}
	dirInfo, err := os.Stat(directory)
	if err != nil {
		return err
	}

	permission := int(dirInfo.Mode().Perm())
	d.Set("directory_permission", permission)

	d.SetId(directory)

	return nil
}
