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

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource = (*localCommandDataSource)(nil)
)

func NewLocalCommandDataSource() datasource.DataSource {
	return &localCommandDataSource{}
}

type localCommandDataSource struct{}

func (a *localCommandDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_command"
}

func (a *localCommandDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "", // TODO: document (mention no side-effects, similar to the caveats on the external data source)
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
			// TODO: naming (allow_non_zero_exit_code ?)
			"skip_error": schema.BoolAttribute{
				Description: "", // TODO: document what users can expect here and how to use it (when it will be populated, defaults)
				Optional:    true,
			},
			"exit_code": schema.Int64Attribute{
				Description: "The exit code returned after executing the given command.", // TODO: Describe it's relationship to diagnostics
				Computed:    true,
			},
			"stdout": schema.StringAttribute{
				Description: "", // TODO: document what users can expect here and how to use it
				Computed:    true,
			},
			"stderr": schema.StringAttribute{
				Description: "", // TODO: document what users can expect here and how to use it (when it will be populated, defaults)
				Computed:    true,
			},
		},
	}
}

type localCommandDataSourceModel struct {
	Command          types.String `tfsdk:"command"`
	Arguments        types.List   `tfsdk:"arguments"`
	Stdin            types.String `tfsdk:"stdin"`
	WorkingDirectory types.String `tfsdk:"working_directory"`
	SkipError        types.Bool   `tfsdk:"skip_error"`
	ExitCode         types.Int64  `tfsdk:"exit_code"`
	Stdout           types.String `tfsdk:"stdout"`
	Stderr           types.String `tfsdk:"stderr"`
}

func (a *localCommandDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state localCommandDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Prep the command
	command := state.Command.ValueString()
	if _, err := exec.LookPath(command); err != nil {
		resp.Diagnostics.AddAttributeError(
			path.Root("command"),
			"Command Lookup Failed",
			"The data source received an unexpected error while attempting to find the command."+
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
		return
	}

	arguments := make([]string, 0)
	for _, element := range state.Arguments.Elements() {
		strElement, ok := element.(types.String)
		// Mirroring the underlying os/exec Command support for args (no nil arguments, but does support empty strings)
		if element.IsNull() || !ok {
			continue
		}

		arguments = append(arguments, strElement.ValueString())
	}

	cmd := exec.CommandContext(ctx, command, arguments...)

	cmd.Dir = state.WorkingDirectory.ValueString()

	if !state.Stdin.IsNull() {
		cmd.Stdin = bytes.NewReader([]byte(state.Stdin.ValueString()))
	}

	var stderr strings.Builder
	cmd.Stderr = &stderr
	var stdout strings.Builder
	cmd.Stdout = &stdout

	tflog.Trace(ctx, "Executing local command", map[string]interface{}{"command": cmd.String()})

	// Run the command
	err := cmd.Run()
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if len(stderrStr) > 0 {
		// TODO: Should we raise an explicit error if this isn't utf8?
		// https://pkg.go.dev/unicode/utf8#example-Valid
		state.Stderr = types.StringValue(stderrStr)
	}

	if len(stdoutStr) > 0 {
		// TODO: Should we raise an explicit error if this isn't utf8?
		// https://pkg.go.dev/unicode/utf8#example-Valid
		state.Stdout = types.StringValue(stdoutStr)
	}

	// ProcessState will always be populated if the command has been was successfully started (regardless of exit code)
	if cmd.ProcessState != nil {
		exitCode := cmd.ProcessState.ExitCode()
		state.ExitCode = types.Int64Value(int64(exitCode))
	}

	tflog.Trace(ctx, "Executed local command", map[string]interface{}{"command": cmd.String(), "stdout": stdoutStr, "stderr": stderrStr})

	// Set all of the data to state
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)

	// If we received an error, we need to check and see if we should explicitly raise a diagnostic
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			// We won't return a diagnostic because the command was successfully started, it just
			// exited with a non-zero code (which the user has indicated they will handle in configuration).
			//
			// All data has already been saved to state, so we just return.
			if state.SkipError.ValueBool() {
				return
			}

			resp.Diagnostics.AddAttributeError(
				path.Root("command"),
				"Command Execution Failed",
				"The data source executed the command but received a non-zero exit code. If a non-zero exit code is expected and can be handled in configuration, set \"skip_error\" to true."+
					"\n\n"+
					fmt.Sprintf("Command: %s\n", cmd.String())+
					fmt.Sprintf("Command Error: %s\n", stderrStr)+
					fmt.Sprintf("State: %s", exitError),
			)
			return
		}

		// We can't skip this error because the command wasn't successfully started.
		resp.Diagnostics.AddAttributeError(
			path.Root("command"),
			"Command Execution Failed",
			"The data source received an unexpected error while attempting to execute the command."+
				"\n\n"+
				fmt.Sprintf("Command: %s\n", cmd.String())+
				fmt.Sprintf("State: %s", err),
		)
		return
	}
}
