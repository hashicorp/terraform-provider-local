package provider

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os/exec"
	"syscall"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceLocalExec() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceLocalExecRead,

		Schema: map[string]*schema.Schema{
			"command": {
				Type:        schema.TypeList,
				Description: "Command to execute",
				Required:    true,
				ForceNew:    true,
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
			"working_dir": {
				Type:        schema.TypeString,
				Description: "Directory to change into before executing provided command",
				Optional:    true,
				Default:     "",
				ForceNew:    true,
			},
			"ignore_failure": {
				Type:        schema.TypeBool,
				Description: "If set to true, command execution failures will be ignored",
				Optional:    true,
				Default:     false,
				ForceNew:    true,
			},
			"stdout": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"stderr": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"rc": {
				Type:     schema.TypeInt,
				Computed: true,
			},
		},
	}
}

func dataSourceLocalExecRead(d *schema.ResourceData, _ interface{}) error {
	exitCode := 0
	command, args, _ := expandCommand(d)
	ignoreFailure := d.Get("ignore_failure").(bool)

	cmd := exec.Command(command, args...)
	cmd.Dir = d.Get("working_dir").(string)

	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	stderr, _ := ioutil.ReadAll(stderrPipe)
	stdout, _ := ioutil.ReadAll(stdoutPipe)

	if err := cmd.Wait(); err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if status, ok := exitError.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
			} else {
				// unable to retrieve exit code from error, use default
				exitCode = -1
			}
		}

		if !ignoreFailure {
			return err
		}
	}

	d.Set("stderr", string(stderr))
	d.Set("stdout", string(stdout))
	d.Set("rc", exitCode)

	// use the checksum of (stdout, stderr, rc) to generate id
	checksum := sha1.Sum(
		append([]byte(stdout),
			append([]byte(stderr),
				[]byte(fmt.Sprintf("%d", exitCode))...)...))
	d.SetId(hex.EncodeToString(checksum[:]))

	return nil
}

func expandCommand(d *schema.ResourceData) (string, []string, error) {
	execCommand := d.Get("command").([]interface{})
	command := make([]string, 0)

	for _, commandRaw := range execCommand {
		command = append(command, commandRaw.(string))
	}

	return command[0], command[1:], nil
}
