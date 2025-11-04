// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ action.Action = (*localCommandAction)(nil)
)

func NewLocalCommandAction() action.Action {
	return &localCommandAction{}
}

type localCommandAction struct{}

func (a *localCommandAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_command"
}

func (a *localCommandAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.Schema = schema.Schema{
		// TODO: Once we have a local_command data source, reference that to be used if the user needs to the consume the output of the command (and it's idempotent)
		MarkdownDescription: "Invokes an executable on the local machine. All environment variables visible to the Terraform process are passed through " +
			"to the child process. After the child process successfully executes, the `stdout` will be returned for Terraform to display to the user.\n\n" +
			"Any non-zero exit code will be treated as an error and will return a diagnostic to Terraform containing the `stderr` message if available.",
		Attributes: map[string]schema.Attribute{
			"command": schema.StringAttribute{
				Description: "Executable name to be discovered on the PATH or absolute path to executable.",
				Required:    true,
			},
			"arguments": schema.ListAttribute{
				MarkdownDescription: "Arguments to be passed to the given command. Any `null` arguments will be removed from the list.",
				ElementType:         types.StringType,
				Optional:            true,
			},
			"stdin": schema.StringAttribute{
				Description: "Data to be passed to the given command's standard input.",
				Optional:    true,
			},
			"working_directory": schema.StringAttribute{
				Description: "The directory where the command should be executed. Defaults to the Terraform working directory.",
				Optional:    true,
			},
		},
	}
}

type localCommandActionModel struct {
	Command          types.String `tfsdk:"command"`
	Arguments        types.List   `tfsdk:"arguments"`
	Stdin            types.String `tfsdk:"stdin"`
	WorkingDirectory types.String `tfsdk:"working_directory"`
}

func (a *localCommandAction) ModifyPlan(ctx context.Context, req action.ModifyPlanRequest, resp *action.ModifyPlanResponse) {
	var command types.String
	resp.Diagnostics.Append(req.Config.GetAttribute(ctx, path.Root("command"), &command)...)
	if resp.Diagnostics.HasError() || command.IsUnknown() {
		return
	}

	resp.Diagnostics.Append(findCommand(command.ValueString()))
}

func (a *localCommandAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	var config localCommandActionModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prep the command
	command := config.Command.ValueString()
	resp.Diagnostics.Append(findCommand(command))
	if resp.Diagnostics.HasError() {
		return
	}

	arguments := make([]string, 0)
	for _, element := range config.Arguments.Elements() {
		strElement, ok := element.(types.String)
		// Mirroring the underlying os/exec Command support for args (no nil arguments, but does support empty strings)
		if element.IsNull() || !ok {
			continue
		}

		arguments = append(arguments, strElement.ValueString())
	}

	cmd := exec.CommandContext(ctx, command, arguments...)

	cmd.Dir = config.WorkingDirectory.ValueString()

	if !config.Stdin.IsNull() {
		cmd.Stdin = bytes.NewReader([]byte(config.Stdin.ValueString()))
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr

	tflog.Trace(ctx, "Executing local command", map[string]interface{}{"command": cmd.String()})

	// Run the command
	stdout, err := cmd.Output()
	stdoutStr := string(stdout)
	stderrStr := stderr.String()

	if err != nil {
		if len(stderrStr) > 0 {
			resp.Diagnostics.AddAttributeError(
				path.Root("command"),
				"Command Execution Failed",
				"The action received an unexpected error while attempting to execute the command."+
					"\n\n"+
					fmt.Sprintf("Command: %s\n", cmd.String())+
					fmt.Sprintf("Command Error: %s\n", stderrStr)+
					fmt.Sprintf("State: %s", err),
			)
			return
		}

		resp.Diagnostics.AddAttributeError(
			path.Root("command"),
			"Command Execution Failed",
			"The action received an unexpected error while attempting to execute the command."+
				"\n\n"+
				fmt.Sprintf("Command: %s\n", cmd.Path)+
				fmt.Sprintf("Error: %s", err),
		)
		return
	}

	tflog.Trace(ctx, "Executed local command", map[string]interface{}{"command": cmd.String(), "stdout": stdoutStr, "stderr": stderrStr})

	// Send the STDOUT to Terraform to display to the practitioner. The underlying action protocol supports streaming the
	// STDOUT line-by-line in real-time, although each progress message gets a prefix per line, so it'd be difficult
	// to read without batching lines together with an arbitrary time interval (this can be improved later if needed).
	resp.SendProgress(action.InvokeProgressEvent{
		Message: fmt.Sprintf("\n\n%s\n", stdoutStr),
	})
}

func findCommand(command string) diag.Diagnostic {
	if _, err := exec.LookPath(command); err != nil {
		return diag.NewAttributeErrorDiagnostic(
			path.Root("command"),
			"Command Lookup Failed",
			"The action received an unexpected error while attempting to find the command."+
				"\n\n"+
				"The command must be accessible according to the platform where Terraform is running."+
				"\n\n"+
				"If the expected command should be automatically found on the platform where Terraform is running, "+
				"ensure that the command is in an expected directory. On Unix-based platforms, these directories are "+
				"typically searched based on the '$PATH' environment variable. On Windows-based platforms, these directories "+
				"are typically searched based on the '%PATH%' environment variable."+
				"\n\n"+
				"If the expected command is relative to the Terraform configuration, it is recommended that the command name includes "+
				"the interpolated value of 'path.module' before the command name to ensure that it is compatible with varying module usage. For example: \"${path.module}/my-command\""+
				"\n\n"+
				"The command must also be executable according to the platform where Terraform is running. On Unix-based platforms, the file on the filesystem must have the executable bit set. "+
				"On Windows-based platforms, no action is typically necessary."+
				"\n\n"+
				fmt.Sprintf("Platform: %s\n", runtime.GOOS)+
				fmt.Sprintf("Command: %s\n", command)+
				fmt.Sprintf("Error: %s", err),
		)
	}

	return nil
}
