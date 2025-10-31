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
		Description: "", // TODO: Describe action, mention that actions don't have output, this is meant to execute local commands, they can be non-idempotent as they are only executed during apply.
		// If the external command is idempotent/you need the output, use data source (coming soon).
		Attributes: map[string]schema.Attribute{
			"command": schema.StringAttribute{
				Description: "Executable name to be discovered on the PATH or absolute path to executable.",
				Required:    true,
			},
			"arguments": schema.ListAttribute{
				Description: "Arguments to be passed to the given command.",
				ElementType: types.StringType,
				Optional:    true,
			},
			"stdin": schema.StringAttribute{
				Description: "Data to be passed to the given command's standard input.",
				Optional:    true,
			},
			"working_directory": schema.StringAttribute{
				Description: "The directory where the command should be executed. Defaults to the current working directory.",
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

	// Prep the commmand
	command := config.Command.ValueString()
	resp.Diagnostics.Append(findCommand(command))
	if resp.Diagnostics.HasError() {
		return
	}

	arguments := make([]string, 0)
	resp.Diagnostics.Append(config.Arguments.ElementsAs(ctx, &arguments, true)...)
	if resp.Diagnostics.HasError() {
		return
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
					fmt.Sprintf("Commmand: %s\n", cmd.String())+
					fmt.Sprintf("Command Error: %s\n", stderrStr)+
					fmt.Sprintf("State: %s", err),
			)
			return
		}

		resp.Diagnostics.Append(genericCommandDiag(cmd, err))
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

func genericCommandDiag(cmd *exec.Cmd, err error) diag.Diagnostic {
	return diag.NewAttributeErrorDiagnostic(
		path.Root("command"),
		"Command Execution Failed",
		"The action received an unexpected error while attempting to execute the command."+
			"\n\n"+
			fmt.Sprintf("Commmand: %s\n", cmd.Path)+
			fmt.Sprintf("Error: %s", err),
	)
}
