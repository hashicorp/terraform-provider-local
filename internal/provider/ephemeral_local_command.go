// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var _ ephemeral.EphemeralResource = (*localCommandEphemeral)(nil)

func NewLocalCommandEphemeral() ephemeral.EphemeralResource {
	return &localCommandEphemeral{}
}

type localCommandEphemeral struct{}

func (e *localCommandEphemeral) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_command"
}

func (e *localCommandEphemeral) Schema(ctx context.Context, req ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Runs an executable on the local machine and returns the exit code, standard output data (`stdout`), and standard error data (`stderr`). " +
			"All environment variables visible to the Terraform process are passed through to the child process. Both `stdout` and `stderr` returned by this ephemeral resource " +
			"are UTF-8 strings, which can be decoded into [Terraform values](https://developer.hashicorp.com/terraform/language/expressions/types) for use elsewhere in the Terraform configuration. " +
			"There are built-in decoding functions such as [`jsondecode`](https://developer.hashicorp.com/terraform/language/functions/jsondecode) or [`yamldecode`](https://developer.hashicorp.com/terraform/language/functions/yamldecode), " +
			"and more specialized [decoding functions](https://developer.hashicorp.com/terraform/plugin/framework/functions/concepts) can be built with a Terraform provider." +
			"\n\n" +
			"Any non-zero exit code returned by the command will be treated as an error and will return a diagnostic to Terraform containing the `stderr` message if available. " +
			"If a non-zero exit code is expected by the command, set `allow_non_zero_exit_code` to `true`." +
			"\n\n" +
			"~> **Warning** This mechanism is provided as an \"escape hatch\" for exceptional situations where a first-class Terraform provider is not more appropriate. " +
			"Its capabilities are limited in comparison to a true ephemeral resource, and implementing an ephemeral resource via a local executable is likely to hurt the " +
			"portability of your Terraform configuration by creating dependencies on external programs and libraries that may not be available (or may need to be used differently) " +
			"on different operating systems." +
			"\n\n" +
			"~> **Warning** HCP Terraform and Terraform Enterprise do not guarantee availability of any particular language runtimes or external programs beyond standard shell utilities, " +
			"so it is not recommended to use this ephemeral resource within configurations that are applied within either.",
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
				MarkdownDescription: "Data to be passed to the given command's standard input as a UTF-8 string. [Terraform values](https://developer.hashicorp.com/terraform/language/expressions/types) can be encoded " +
					"by any Terraform encode function, for example, [`jsonencode`](https://developer.hashicorp.com/terraform/language/functions/jsonencode).",
				Optional: true,
			},
			"working_directory": schema.StringAttribute{
				Description: "The directory path where the command should be executed, either an absolute path or relative to the Terraform working directory. If not provided, defaults to the Terraform working directory.",
				Optional:    true,
			},
			"allow_non_zero_exit_code": schema.BoolAttribute{
				MarkdownDescription: "Indicates that the command returning a non-zero exit code should be treated as a successful execution. " +
					"Further assertions can be made of the `exit_code` value with the [`check` block](https://developer.hashicorp.com/terraform/language/block/check). Defaults to false.",
				Optional: true,
			},
			"exit_code": schema.Int64Attribute{
				MarkdownDescription: "The exit code returned by the command. By default, if the exit code is non-zero, the ephemeral resource will return a diagnostic to Terraform. " +
					"If a non-zero exit code is expected by the command, set `allow_non_zero_exit_code` to `true`.",
				Computed: true,
			},
			"stdout": schema.StringAttribute{
				MarkdownDescription: "Data returned from the command's standard output stream. The data is returned directly from the command as a UTF-8 string, " +
					"which can then be decoded by any Terraform decode function, for example, [`jsondecode`](https://developer.hashicorp.com/terraform/language/functions/jsondecode).",
				Computed: true,
			},
			"stderr": schema.StringAttribute{
				Description: "Data returned from the command's standard error stream. The data is returned directly from the command as a UTF-8 string and will be " +
					"populated regardless of the exit code returned.",
				Computed: true,
			},
		},
	}
}

type localCommandEphemeralModel struct {
	Command              types.String `tfsdk:"command"`
	Arguments            types.List   `tfsdk:"arguments"`
	Stdin                types.String `tfsdk:"stdin"`
	WorkingDirectory     types.String `tfsdk:"working_directory"`
	AllowNonZeroExitCode types.Bool   `tfsdk:"allow_non_zero_exit_code"`
	ExitCode             types.Int64  `tfsdk:"exit_code"`
	Stdout               types.String `tfsdk:"stdout"`
	Stderr               types.String `tfsdk:"stderr"`
}

func (e *localCommandEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var state localCommandEphemeralModel
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
			"The ephemeral resource received an unexpected error while attempting to find the command."+
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
	commandErr := cmd.Run()
	stdoutStr := stdout.String()
	stderrStr := stderr.String()

	if len(stderrStr) > 0 {
		state.Stderr = types.StringValue(stderrStr)
	}

	if len(stdoutStr) > 0 {
		state.Stdout = types.StringValue(stdoutStr)
	}

	// ProcessState will always be populated if the command has been was successfully started (regardless of exit code)
	if cmd.ProcessState != nil {
		exitCode := cmd.ProcessState.ExitCode()
		state.ExitCode = types.Int64Value(int64(exitCode))
	}

	tflog.Trace(ctx, "Executed local command", map[string]interface{}{"command": cmd.String(), "stdout": stdoutStr, "stderr": stderrStr})

	// Set all of the data to result
	resp.Diagnostics.Append(resp.Result.Set(ctx, state)...)
	if commandErr == nil {
		return
	}

	// If running the command returned an exit error, we need to check and see if we should explicitly raise a diagnostic
	if exitError, ok := commandErr.(*exec.ExitError); ok {
		// We won't return a diagnostic because the command was successfully started and then exited
		// with a non-zero exit code (which the user has indicated they will handle in configuration).
		//
		// All data has already been saved to result, so we just return.
		if state.AllowNonZeroExitCode.ValueBool() {
			return
		}

		resp.Diagnostics.AddAttributeError(
			path.Root("command"),
			"Command Execution Failed",
			"The ephemeral resource executed the command but received a non-zero exit code. If a non-zero exit code is expected "+
				"and can be handled in configuration, set \"allow_non_zero_exit_code\" to true."+
				"\n\n"+
				fmt.Sprintf("Command: %s\n", cmd.String())+
				fmt.Sprintf("Command Error: %s\n", stderrStr)+
				fmt.Sprintf("State: %s", exitError),
		)
		return
	}

	// We need to raise a diagnostic because the command wasn't successfully started and we have no exit code.
	resp.Diagnostics.AddAttributeError(
		path.Root("command"),
		"Command Execution Failed",
		"The ephemeral resource received an unexpected error while attempting to execute the command."+
			"\n\n"+
			fmt.Sprintf("Command: %s\n", cmd.String())+
			fmt.Sprintf("State: %s", commandErr),
	)
}
