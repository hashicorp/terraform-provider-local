// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource = (*localFileDataSource)(nil)
)

func NewLocalFileDataSource() datasource.DataSource {
	return &localFileDataSource{}
}

type localFileDataSource struct{}

func (n *localFileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a file from the local filesystem.",
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Description: "Path to the file that will be read. The data source will return an error if the file does not exist.",
				Required:    true,
			},
			"content": schema.StringAttribute{
				Description: "Raw content of the file that was read, as UTF-8 encoded string. " +
					"Files that do not contain UTF-8 text will have invalid UTF-8 sequences in `content`\n  replaced with the Unicode replacement character. ",
				Computed: true,
			},
			"content_base64": schema.StringAttribute{
				Description: "Base64 encoded version of the file content (use this when dealing with binary data).",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The hexadecimal encoding of the SHA1 checksum of the file content.",
				Computed:    true,
			},
			"content_md5": schema.StringAttribute{
				Description: "MD5 checksum of file content.",
				Computed:    true,
			},
			"content_sha1": schema.StringAttribute{
				Description: "SHA1 checksum of file content.",
				Computed:    true,
			},
			"content_sha256": schema.StringAttribute{
				Description: "SHA256 checksum of file content.",
				Computed:    true,
			},
			"content_base64sha256": schema.StringAttribute{
				Description: "Base64 encoded SHA256 checksum of file content.",
				Computed:    true,
			},
			"content_sha512": schema.StringAttribute{
				Description: "SHA512 checksum of file content.",
				Computed:    true,
			},
			"content_base64sha512": schema.StringAttribute{
				Description: "Base64 encoded SHA512 checksum of file content.",
				Computed:    true,
			},
		},
	}
}

func (n *localFileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (n *localFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var config localFileDataSourceModelV0

	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Read the entire file content
	filepath := config.Filename.ValueString()
	content, err := os.ReadFile(filepath)

	if err != nil {
		resp.Diagnostics.AddError(
			"Read local file data source error",
			"The file at given path cannot be read.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	//calculate the checksums of file content
	checksums := genFileChecksums(content)

	state := localFileDataSourceModelV0{
		Filename:            config.Filename,
		Content:             types.StringValue(string(content)),
		ContentBase64:       types.StringValue(base64.StdEncoding.EncodeToString(content)),
		ID:                  types.StringValue(checksums.sha1Hex),
		ContentMd5:          types.StringValue(checksums.md5Hex),
		ContentSha1:         types.StringValue(checksums.sha1Hex),
		ContentSha256:       types.StringValue(checksums.sha256Hex),
		ContentBase64sha256: types.StringValue(checksums.sha256Base64),
		ContentSha512:       types.StringValue(checksums.sha512Hex),
		ContentBase64sha512: types.StringValue(checksums.sha512Base64),
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

type localFileDataSourceModelV0 struct {
	Filename            types.String `tfsdk:"filename"`
	Content             types.String `tfsdk:"content"`
	ContentBase64       types.String `tfsdk:"content_base64"`
	ID                  types.String `tfsdk:"id"`
	ContentMd5          types.String `tfsdk:"content_md5"`
	ContentSha1         types.String `tfsdk:"content_sha1"`
	ContentSha256       types.String `tfsdk:"content_sha256"`
	ContentBase64sha256 types.String `tfsdk:"content_base64sha256"`
	ContentSha512       types.String `tfsdk:"content_sha512"`
	ContentBase64sha512 types.String `tfsdk:"content_base64sha512"`
}
