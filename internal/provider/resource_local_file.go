package provider

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/hex"
	"io/ioutil"
	"os"
	"path"
	"strconv"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceLocalFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceLocalFileCreate,
		Read:   resourceLocalFileRead,
		Delete: resourceLocalFileDelete,

		Description: "Generates a local file with the given content.",

		Schema: map[string]*schema.Schema{
			"content": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"content", "sensitive_content", "content_base64", "source"},
				Description:  "Content to store in the file, expected to be an UTF-8 encoded string.",
			},
			"sensitive_content": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Sensitive:    true,
				ExactlyOneOf: []string{"content", "sensitive_content", "content_base64", "source"},
				Description:  "Sensitive content to store in the file, expected to be an UTF-8 encoded string.",
				Deprecated:   "Use the `local_sensitive_file` resource instead.",
			},
			"content_base64": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"content", "sensitive_content", "content_base64", "source"},
				Description:  "Content to store in the file, expected to be binary encoded as base64 string.",
			},
			"source": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				ExactlyOneOf: []string{"content", "sensitive_content", "content_base64", "source"},
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
				Default:      "0777",
				ValidateFunc: validateModePermission,
				Description:  "Permissions to set for the output file (in numeric notation).",
			},
			"directory_permission": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Default:      "0777",
				ValidateFunc: validateModePermission,
				Description:  "Permissions to set for directories created (in numeric notation).",
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

func resourceLocalFileRead(d *schema.ResourceData, _ interface{}) error {
	// If the output file doesn't exist, mark the resource for creation.
	outputPath := d.Get("filename").(string)
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		d.SetId("")
		return nil
	}

	// Verify that the content of the destination file matches the content we
	// expect. Otherwise, the file might have been modified externally, and we
	// must reconcile.
	outputContent, err := ioutil.ReadFile(outputPath)
	if err != nil {
		return err
	}

	checksums := genFileChecksums(outputContent)
	d.Set("content_md5", checksums.md5Hex)
	d.Set("content_sha1", checksums.sha1Hex)
	d.Set("content_sha256", checksums.sha256Hex)
	d.Set("content_base64sha256", checksums.sha256Base64)
	d.Set("content_sha512", checksums.sha512Hex)
	d.Set("content_base64sha512", checksums.sha512Base64)

	if checksums.sha1Hex != d.Id() {
		d.SetId("")
		return nil
	}

	return nil
}

func resourceLocalFileContent(d *schema.ResourceData) ([]byte, error) {
	if sensitiveContent, ok := d.GetOk("sensitive_content"); ok {
		return []byte(sensitiveContent.(string)), nil
	}
	if contentBase64, ok := d.GetOk("content_base64"); ok {
		return base64.StdEncoding.DecodeString(contentBase64.(string))
	}

	if sourceFile, ok := d.GetOk("source"); ok {
		sourceFileContent := sourceFile.(string)
		return ioutil.ReadFile(sourceFileContent)
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

	if err := ioutil.WriteFile(destination, content, os.FileMode(fileMode)); err != nil {
		return err
	}

	checksums := genFileChecksums(content)
	d.Set("content_md5", checksums.md5Hex)
	d.Set("content_sha1", checksums.sha1Hex)
	d.Set("content_sha256", checksums.sha256Hex)
	d.Set("content_base64sha256", checksums.sha256Base64)
	d.Set("content_sha512", checksums.sha512Hex)
	d.Set("content_base64sha512", checksums.sha512Base64)

	d.SetId(checksums.sha1Hex)

	return nil
}

func resourceLocalFileDelete(d *schema.ResourceData, _ interface{}) error {
	os.Remove(d.Get("filename").(string))
	return nil
}

type fileChecksums struct {
	md5Hex       string
	sha1Hex      string
	sha256Hex    string
	sha256Base64 string
	sha512Hex    string
	sha512Base64 string
}

func genFileChecksums(data []byte) fileChecksums {
	checksums := fileChecksums{}

	md5Sum := md5.Sum(data)
	checksums.md5Hex = hex.EncodeToString(md5Sum[:])

	sha1Sum := sha1.Sum(data)
	checksums.sha1Hex = hex.EncodeToString(sha1Sum[:])

	sha256Sum := sha256.Sum256(data)
	checksums.sha256Hex = hex.EncodeToString(sha256Sum[:])
	checksums.sha256Base64 = base64.StdEncoding.EncodeToString(sha256Sum[:])

	sha512Sum := sha512.Sum512(data)
	checksums.sha512Hex = hex.EncodeToString(sha512Sum[:])
	checksums.sha512Base64 = base64.StdEncoding.EncodeToString(sha512Sum[:])

	return checksums
}
