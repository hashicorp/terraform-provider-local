// Copyright IBM Corp. 2017, 2025
// SPDX-License-Identifier: MPL-2.0

package provider

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/terraform-providers/terraform-provider-local/internal/localtypes"
)

var (
	_ resource.Resource = (*localSensitiveFileResource)(nil)
)

func NewLocalSensitiveFileResource() resource.Resource {
	return &localSensitiveFileResource{}
}

type localSensitiveFileResource struct{}

func (n *localSensitiveFileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Generates a local file with the given sensitive content.",
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Description: "The path to the file that will be created.\n " +
					"Missing parent directories will be created.\n " +
					"If the file already exists, it will be overridden with the given content.",
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Description: "Sensitive Content to store in the file, expected to be a UTF-8 encoded string.\n " +
					"Conflicts with `content_base64` and `source`.\n " +
					"Exactly one of these three arguments must be specified.",
				Sensitive: true,
				Optional:  true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("content_base64"),
						path.MatchRoot("source")),
				},
			},
			"content_base64": schema.StringAttribute{
				Description: "Sensitive Content to store in the file, expected to be binary encoded as base64 string.\n " +
					"Conflicts with `content` and `source`.\n " +
					"Exactly one of these three arguments must be specified.",
				Sensitive: true,
				Optional:  true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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
					"Default value is `\"0700\"`.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("0700"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"directory_permission": schema.StringAttribute{
				CustomType: localtypes.NewFilePermissionType(),
				Description: "Permissions to set for directories created (before umask), expressed as string in\n " +
					"[numeric notation](https://en.wikipedia.org/wiki/File-system_permissions#Numeric_notation).\n " +
					"Default value is `\"0700\"`.",
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString("0700"),
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
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

func (n *localSensitiveFileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive_file"
}

func (n *localSensitiveFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan localSensitiveFileResourceModelV0

	var filePerm, dirPerm string

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	content, err := parseLocalSensitiveFileContent(plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Create local sensitive file error",
			"An unexpected error occurred while parsing local file content\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	destination := plan.Filename.ValueString()

	destinationDir := filepath.Dir(destination)
	if _, err := os.Stat(destinationDir); err != nil {
		dirPerm = plan.DirectoryPermission.ValueString()
		dirMode, _ := strconv.ParseInt(dirPerm, 8, 64)
		if err := os.MkdirAll(destinationDir, os.FileMode(dirMode)); err != nil {
			resp.Diagnostics.AddError(
				"Create local sensitive file error",
				"An unexpected error occurred while creating file directory\n\n+"+
					fmt.Sprintf("Original Error: %s", err),
			)
			return
		}
	}

	filePerm = plan.FilePermission.ValueString()

	fileMode, _ := strconv.ParseInt(filePerm, 8, 64)

	if err := os.WriteFile(destination, content, os.FileMode(fileMode)); err != nil {
		resp.Diagnostics.AddError(
			"Create local sensitive file error",
			"An unexpected error occurred while writing the file\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	checksums := genFileChecksums(content)
	plan.ContentMd5 = types.StringValue(checksums.md5Hex)
	plan.ContentSha1 = types.StringValue(checksums.sha1Hex)
	plan.ContentSha256 = types.StringValue(checksums.sha256Hex)
	plan.ContentBase64sha256 = types.StringValue(checksums.sha256Base64)
	plan.ContentSha512 = types.StringValue(checksums.sha512Hex)
	plan.ContentBase64sha512 = types.StringValue(checksums.sha512Base64)

	plan.ID = types.StringValue(checksums.sha1Hex)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (n *localSensitiveFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state localSensitiveFileResourceModelV0

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// If the output file doesn't exist, mark the resource for creation.
	outputPath := state.Filename.ValueString()
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		resp.State.RemoveResource(ctx)
		return
	}

	// Verify that the content of the destination file matches the content we
	// expect. Otherwise, the file might have been modified externally, and we
	// must reconcile.
	outputContent, err := os.ReadFile(outputPath)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read local sensitive file error",
			"An unexpected error occurred while reading the file\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	outputChecksum := sha1.Sum(outputContent)
	if hex.EncodeToString(outputChecksum[:]) != state.ID.ValueString() {
		resp.State.RemoveResource(ctx)
		return
	}
}

func (n *localSensitiveFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan localSensitiveFileResourceModelV0

	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &plan)...)
}

func (n *localSensitiveFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var filename string
	req.State.GetAttribute(ctx, path.Root("filename"), &filename)
	os.Remove(filename)
}

func parseLocalSensitiveFileContent(plan localSensitiveFileResourceModelV0) ([]byte, error) {
	if !plan.ContentBase64.IsNull() && !plan.ContentBase64.IsUnknown() {
		return base64.StdEncoding.DecodeString(plan.ContentBase64.ValueString())
	}

	if !plan.Source.IsNull() && !plan.Source.IsUnknown() {
		sourceFileContent := plan.Source.ValueString()
		return os.ReadFile(sourceFileContent)
	}

	content := plan.Content.ValueString()
	return []byte(content), nil
}

type localSensitiveFileResourceModelV0 struct {
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
