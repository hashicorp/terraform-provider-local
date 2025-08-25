// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
	"github.com/hashicorp/terraform-plugin-framework/ephemeral/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/terraform-providers/terraform-provider-local/internal/localtypes"
)

var (
	_ ephemeral.EphemeralResource = (*localFileEphemeralResource)(nil)
)

func NewLocalFileEphemeralResource() ephemeral.EphemeralResource {
	return &localFileEphemeralResource{}
}

type localFileEphemeralResource struct{}

func (e *localFileEphemeralResource) Schema(_ context.Context, _ ephemeral.SchemaRequest, resp *ephemeral.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates an ephemeral local file with the given content.",
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Description: "The path to the file that will be created.\n " +
					"Missing parent directories will be created.\n " +
					"If the file already exists, it will be overridden with the given content.",
				Required: true,
			},
			"content": schema.StringAttribute{
				Description: "Content to store in the file, expected to be a UTF-8 encoded string.\n " +
					"Conflicts with `content_base64` and `source`.\n " +
					"Exactly one of these three arguments must be specified.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("content_base64"),
						path.MatchRoot("source")),
				},
			},
			"content_base64": schema.StringAttribute{
				Description: "Content to store in the file, expected to be binary encoded as base64 string.\n " +
					"Conflicts with `content` and `source`.\n " +
					"Exactly one of these three arguments must be specified.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("content"),
						path.MatchRoot("source")),
				},
			},
			"source": schema.StringAttribute{
				Description: "Path to file to use as source for the one we are creating.\n " +
					"Conflicts with `content` and `content_base64`.\n " +
					"Exactly one of these three arguments must be specified.",
				Optional: true,
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("content"),
						path.MatchRoot("content_base64")),
				},
			},
			"file_permission": schema.StringAttribute{
				CustomType: localtypes.NewFilePermissionType(),
				Description: "Permissions to set for the output file (before umask), expressed as string in\n " +
					"[numeric notation](https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation).\n " +
					"Default value is `\"0777\"`.",
				Optional: true,
				Computed: true,
				// Can't set a default value for ephemeral resources, this is here as a fingers-crossed placeholder.
				// Default:  stringdefault.StaticString("0777"),
			},
			"directory_permission": schema.StringAttribute{
				CustomType: localtypes.NewFilePermissionType(),
				Description: "Permissions to set for directories created (before umask), expressed as string in\n " +
					"[numeric notation](https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation).\n " +
					"Default value is `\"0777\"`.",
				Optional: true,
				Computed: true,
				// Can't set a default value for ephemeral resources, this is here as a fingers-crossed placeholder.
				// Default:  stringdefault.StaticString("0777"),
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

func (e *localFileEphemeralResource) Metadata(ctx context.Context, req ephemeral.MetadataRequest, resp *ephemeral.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file" // local_file
}

func (e *localFileEphemeralResource) Open(ctx context.Context, req ephemeral.OpenRequest, resp *ephemeral.OpenResponse) {
	var data localFileEphemeralResourceModelV0
	var filePerm, dirPerm string

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	content, err := parseEphemeralLocalFileContent(data)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create ephemeral local file error",
			"An unexpected error occurred while parsing ephemeral local file content\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	destination := data.Filename.ValueString()
	privateData, _ := json.Marshal(localFilePrivateData{Filename: destination})
	resp.Private.SetKey(ctx, "local_file_data", privateData)

	destinationDir := filepath.Dir(destination)
	if _, err := os.Stat(destinationDir); err != nil {
		dirPerm = data.DirectoryPermission.ValueString()
		if dirPerm == "" {
			dirPerm = "0777"
		}
		dirPermData := localtypes.FilePermissionValue{StringValue: basetypes.NewStringValue(dirPerm)}
		data.DirectoryPermission = dirPermData
		dirMode, _ := strconv.ParseInt(dirPerm, 8, 64)
		if err := os.MkdirAll(destinationDir, os.FileMode(dirMode)); err != nil {
			resp.Diagnostics.AddError(
				"Create local file error",
				"An unexpected error occurred while creating file directory\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
	}

	filePerm = data.FilePermission.ValueString()
	if filePerm == "" {
		filePerm = "0777"
	}
	filePermData := localtypes.FilePermissionValue{StringValue: basetypes.NewStringValue(filePerm)}
	data.FilePermission = filePermData

	fileMode, _ := strconv.ParseInt(filePerm, 8, 64)

	if err := os.WriteFile(destination, content, os.FileMode(fileMode)); err != nil {
		resp.Diagnostics.AddError(
			"Create local file error",
			"An unexpected error occurred while writing the file\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	tflog.Debug(ctx, fmt.Sprintf("Created ephemeral file with name: %s", destination))

	checksums := genFileChecksums(content)
	data.ContentMd5 = types.StringValue(checksums.md5Hex)
	data.ContentSha1 = types.StringValue(checksums.sha1Hex)
	data.ContentSha256 = types.StringValue(checksums.sha256Hex)
	data.ContentBase64sha256 = types.StringValue(checksums.sha256Base64)
	data.ContentSha512 = types.StringValue(checksums.sha512Hex)
	data.ContentBase64sha512 = types.StringValue(checksums.sha512Base64)

	data.ID = types.StringValue(checksums.sha1Hex)
	resp.Diagnostics.Append(resp.Result.Set(ctx, &data)...)
}

func (e *localFileEphemeralResource) Close(ctx context.Context, req ephemeral.CloseRequest, resp *ephemeral.CloseResponse) {
	// Destroy the file
	privateBytes, diags := req.Private.GetKey(ctx, "local_file_data")
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var privateData localFilePrivateData
	if err := json.Unmarshal(privateBytes, &privateData); err != nil {
		resp.Diagnostics.AddError(
			"Private data unmarshal error",
			"An unexpected error occurred while unmarshaling private data\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	if privateData.Filename != "" {
		tflog.Debug(ctx, fmt.Sprintf("Deleting ephemeral file: %s", privateData.Filename))
		os.Remove(privateData.Filename)
	}
}

func parseEphemeralLocalFileContent(data localFileEphemeralResourceModelV0) ([]byte, error) {
	if !data.ContentBase64.IsNull() && !data.ContentBase64.IsUnknown() {
		return base64.StdEncoding.DecodeString(data.ContentBase64.ValueString())
	}

	if !data.Source.IsNull() && !data.Source.IsUnknown() {
		sourceFileContent := data.Source.ValueString()
		return os.ReadFile(sourceFileContent)
	}

	content := data.Content.ValueString()
	return []byte(content), nil
}

type localFileEphemeralResourceModelV0 struct {
	Filename            types.String                   `tfsdk:"filename"`
	Content             types.String                   `tfsdk:"content"`
	ContentBase64       types.String                   `tfsdk:"content_base64"`
	Source              types.String                   `tfsdk:"source"`
	FilePermission      localtypes.FilePermissionValue `tfsdk:"file_permission"`
	DirectoryPermission localtypes.FilePermissionValue `tfsdk:"directory_permission"`
	ID                  types.String                   `tfsdk:"id"`
	ContentMd5          types.String                   `tfsdk:"content_md5"`
	ContentSha1         types.String                   `tfsdk:"content_sha1"`
	ContentSha256       types.String                   `tfsdk:"content_sha256"`
	ContentBase64sha256 types.String                   `tfsdk:"content_base64sha256"`
	ContentSha512       types.String                   `tfsdk:"content_sha512"`
	ContentBase64sha512 types.String                   `tfsdk:"content_base64sha512"`
}

type localFilePrivateData struct {
	Filename string `json:"filename"`
}
