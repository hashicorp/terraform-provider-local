package local

import (
	"crypto/sha1"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"

	"github.com/hashicorp/terraform/helper/schema"
)

func resourceLocalFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceLocalFileCreate,
		Read:   resourceLocalFileRead,
		Delete: resourceLocalFileDelete,

		Schema: map[string]*schema.Schema{
			"content": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"mode": {
				Type:        schema.TypeInt,
				Description: "File mode of the output file",
				Optional:    true,
				ForceNew:    true,
				Default:     0777,
			},
			"dir_mode": {
				Type:        schema.TypeInt,
				Description: "File mode for parent directories if they are created",
				Optional:    true,
				ForceNew:    true,
				Default:     0777,
			},
			"filename": {
				Type:        schema.TypeString,
				Description: "Path to the output file",
				Required:    true,
				ForceNew:    true,
			},
		},
	}
}

func resourceLocalFileRead(d *schema.ResourceData, _ interface{}) error {
	// If the output file doesn't exist, mark the resource for creation.
	outputPath := d.Get("filename").(string)
	fi, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		d.SetId("")
		return nil
	}

	// Verify that the mode of the destination file mathes the mode we expect.
	// Otherwise, we must reconcile.
	if fi.Mode().Perm() != os.FileMode(d.Get("mode").(int)) {
		d.SetId("")
		return nil
	}

	// Verify that the content of the destination file matches the content we
	// expect. Otherwise, the file might have been modified externally and we
	// must reconcile.
	outputContent, err := ioutil.ReadFile(outputPath)
	if err != nil {
		return err
	}

	outputChecksum := sha1.Sum([]byte(outputContent))
	if hex.EncodeToString(outputChecksum[:]) != d.Id() {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceLocalFileCreate(d *schema.ResourceData, _ interface{}) error {
	content := d.Get("content").(string)
	destination := d.Get("filename").(string)
	dir_mode := os.FileMode(d.Get("dir_mode").(int))
	mode := os.FileMode(d.Get("mode").(int))

	destinationDir := path.Dir(destination)
	if _, err := os.Stat(destinationDir); err != nil {
		if err := os.MkdirAll(destinationDir, dir_mode); err != nil {
			return err
		}
	}

	if err := ioutil.WriteFile(destination, []byte(content), mode); err != nil {
		return err
	}

	checksum := sha1.Sum([]byte(content))
	d.SetId(hex.EncodeToString(checksum[:]))

	return nil
}

func resourceLocalFileDelete(d *schema.ResourceData, _ interface{}) error {
	os.Remove(d.Get("filename").(string))
	return nil
}
