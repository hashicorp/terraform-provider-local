package local

import (
	"github.com/hashicorp/terraform/helper/schema"
	"os"
)

func resourceLocalSymlink() *schema.Resource {
	return &schema.Resource{
		Create: resourceLocalSymlinkCreate,
		Read:   resourceLocalSymlinkRead,
		Delete: resourceLocalSymlinkDelete,

		Schema: map[string]*schema.Schema{
			"symlink": {
				Type:        schema.TypeString,
				Description: "Path of the symlink",
				Required:    true,
				ForceNew:    true,
			},
			"destination": {
				Type:        schema.TypeString,
				Description: "Permissions to set for directories created",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceLocalSymlinkRead(d *schema.ResourceData, _ interface{}) error {
	// If the symlink doesn't exist, mark the resource for creation.
	wantedSymlink := d.Get("symlink").(string)
	if _, err := os.Stat(wantedSymlink); os.IsNotExist(err) {
		d.SetId("")
		return nil
	}

	// If the symlink is not targeting the good destination, mark the resource for creation.
	wantedDestination := d.Get("destination").(string)
	if destination, err := os.Readlink(wantedSymlink); err == nil && destination != wantedDestination {
		d.SetId("")
		return nil
	}

	d.SetId(wantedSymlink)

	return nil
}

func resourceLocalSymlinkCreate(d *schema.ResourceData, _ interface{}) error {
	wantedSymlink := d.Get("symlink").(string)
	wantedDestination := d.Get("destination").(string)
	if _, err := os.Stat(wantedSymlink); os.IsNotExist(err) {
		if err := os.Symlink(wantedDestination, wantedSymlink); err != nil {
			return err
		}
	}
	if destination, err := os.Readlink(wantedSymlink); err == nil && destination != wantedDestination {
		if err := os.Remove(d.Get("symlink").(string)); err != nil {
			return err
		}
		if err := os.Symlink(wantedDestination, wantedSymlink); err != nil {
			return err
		}
	}

	d.SetId(wantedSymlink)

	return nil
}

func resourceLocalSymlinkDelete(d *schema.ResourceData, _ interface{}) error {
	os.Remove(d.Get("symlink").(string))
	return nil
}
