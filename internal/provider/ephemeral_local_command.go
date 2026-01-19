// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
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
			"All environment variables visible to the Terraform process are passed through to the child process. Both `stdout` and `stderr` returned by this data source " +
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
				MarkdownDescription: "The exit code returned by the command. By default, if the exit code is non-zero, the data source will return a diagnostic to Terraform. " +
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

func (e *localCommandEphemeral) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
}
