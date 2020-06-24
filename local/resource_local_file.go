package local

import (
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/helper/schema"
)

func resourceLocalFile() *schema.Resource {
	return &schema.Resource{
		Read:   resourceLocalFileRead,
		Create: resourceLocalFileCreate,
		Delete: resourceLocalFileDelete,
		Update: nil,
		Exists: func(d *schema.ResourceData, meta interface{}) (bool, error) {
			if _, err := os.Stat(d.Get("filename").(string)); os.IsNotExist(err) {
				return false, nil
			}
			return true, nil
		},
		Schema: map[string]*schema.Schema{
			"content": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"sensitive_content", "content_base64", "source"},
			},
			"sensitive_content": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				Sensitive:     true,
				ConflictsWith: []string{"content", "content_base64", "source"},
			},
			"content_base64": {
				Type:          schema.TypeString,
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"sensitive_content", "content", "source"},
			},
			"filename": {
				Type:        schema.TypeString,
				Description: "Path to the output file",
				Required:    true,
				ForceNew:    true,
			},
			"file_permission": {
				Type:         schema.TypeString,
				Description:  "Permissions to set for the output file",
				Optional:     true,
				ForceNew:     true,
				Default:      "0777",
				ValidateFunc: validateMode,
			},
			"directory_permission": {
				Type:         schema.TypeString,
				Description:  "Permissions to set for directories created",
				Optional:     true,
				ForceNew:     true,
				Default:      "0777",
				ValidateFunc: validateMode,
			},
			"source": {
				Type:          schema.TypeString,
				Description:   "Path to file to use as source for content of output file",
				Optional:      true,
				ForceNew:      true,
				ConflictsWith: []string{"content", "sensitive_content", "content_base64"},
			},
		},
	}
}

func resourceLocalFileRead(d *schema.ResourceData, _ interface{}) error {
	// Get actual content from file.
	filePath := d.Get("filename").(string)
	byteContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return err
	}

	var setErr error

	// Set file_permission to match what is on disk.
	stat, _ := os.Stat(filePath)
	setErr = d.Set("file_permission", fmt.Sprintf("%04o", stat.Mode().Perm()))
	if setErr != nil {
		return err
	}

	// Set `content` or `content_base64` to match current value on disk.
	if _, exists := d.GetOkExists("content"); exists {
		setErr = d.Set("content", string(byteContent))
	} else if _, exists := d.GetOkExists("content_base64"); exists {
		setErr = d.Set("content_base64", base64.StdEncoding.EncodeToString(byteContent))
	}
	if setErr != nil {
		return err
	}

	return nil
}

func resourceLocalFileContent(d *schema.ResourceData) ([]byte, error) {
	if content, sensitiveSpecified := d.GetOk("sensitive_content"); sensitiveSpecified {
		return []byte(content.(string)), nil
	}
	if b64Content, b64Specified := d.GetOk("content_base64"); b64Specified {
		return base64.StdEncoding.DecodeString(b64Content.(string))
	}

	if v, ok := d.GetOk("source"); ok {
		source := v.(string)
		return ioutil.ReadFile(source)
	}

	content := d.Get("content")
	return []byte(content.(string)), nil
}

func resourceLocalFileCreate(d *schema.ResourceData, _ interface{}) error {
	content, err := resourceLocalFileContent(d)
	if err != nil {
		return err
	}

	destination := d.Get("filename").(string)

	destinationDir := path.Dir(destination)
	if _, err := os.Stat(destinationDir); err != nil {
		dirPerm := d.Get("directory_permission").(string)
		dirMode, _ := strconv.ParseInt(dirPerm, 8, 64)
		if err := os.MkdirAll(destinationDir, os.FileMode(dirMode)); err != nil {
			return err
		}
	}

	filePerm := d.Get("file_permission").(string)

	fileMode, _ := strconv.ParseInt(filePerm, 8, 64)

	if err := ioutil.WriteFile(destination, []byte(content), os.FileMode(fileMode)); err != nil {
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
