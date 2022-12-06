package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/terraform-providers/terraform-provider-local/internal/localtypes"
	"github.com/terraform-providers/terraform-provider-local/internal/modifiers/stringmodifier"
)

var (
	_ resource.Resource           = (*localSensitiveFileResource)(nil)
	_ resource.ResourceWithSchema = (*localSensitiveFileResource)(nil)
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
				Description: `
					The path to the file that will be created.
					Missing parent directories will be created.
					If the file already exists, it will be overridden with the given content.
				`,
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"content": schema.StringAttribute{
				Description: "Sensitive content to store in the file, expected to be an UTF-8 encoded string.",
				Sensitive:   true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("sensitive_content"),
						path.MatchRoot("content_base64"),
						path.MatchRoot("source")),
				},
			},
			"content_base64": schema.StringAttribute{
				Description: "Sensitive content to store in the file, expected to be binary encoded as base64 string.",
				Sensitive:   true,
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("content"),
						path.MatchRoot("sensitive_content"),
						path.MatchRoot("source")),
				},
			},
			"source": schema.StringAttribute{
				Description: "Path to file to use as source for the one we are creating.",
				Optional:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("content"),
						path.MatchRoot("sensitive_content"),
						path.MatchRoot("content_base64")),
				},
			},
			"file_permission": schema.StringAttribute{
				CustomType:  localtypes.NewFilePermissionType(),
				Description: "Permissions to set for the output file (in numeric notation).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringmodifier.StringDefault("0777"),
				},
			},
			"directory_permission": schema.StringAttribute{
				CustomType:  localtypes.NewFilePermissionType(),
				Description: "Permissions to set for directories created (in numeric notation).",
				Optional:    true,
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
					stringmodifier.StringDefault("0777"),
				},
			},
			"id": schema.StringAttribute{
				Description: "The hexadecimal encoding of the checksum of the file content",
				Computed:    true,
			},
			"sensitive_content": schema.StringAttribute{
				DeprecationMessage: "Use the `local_sensitive_file` resource instead",
				Description:        "Sensitive content to store in the file, expected to be an UTF-8 encoded string.",
				Sensitive:          true,
				Optional:           true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.ExactlyOneOf(
						path.MatchRoot("content"),
						path.MatchRoot("content_base64"),
						path.MatchRoot("source")),
				},
			},
		},
	}
}

func (n *localSensitiveFileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive_file"
}

func (n *localSensitiveFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	NewLocalFileResource().Create(ctx, req, resp)
}

func (n *localSensitiveFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	NewLocalFileResource().Read(ctx, req, resp)
}

func (n *localSensitiveFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	NewLocalFileResource().Update(ctx, req, resp)
}

func (n *localSensitiveFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	NewLocalFileResource().Delete(ctx, req, resp)
}
