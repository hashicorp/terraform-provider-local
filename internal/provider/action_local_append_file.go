package provider

import (
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/hashicorp/terraform-plugin-framework/action"
	"github.com/hashicorp/terraform-plugin-framework/action/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

var (
	_ action.Action = (*localAppendFileAction)(nil)
)

func NewLocalAppendFile() action.Action {
	return &localAppendFileAction{}
}

type localAppendFileAction struct{}

func (l *localAppendFileAction) Metadata(ctx context.Context, req action.MetadataRequest, resp *action.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_append_file"

}

func (l *localAppendFileAction) Schema(ctx context.Context, req action.SchemaRequest, resp *action.SchemaResponse) {
	resp.LinkedResources = action.LinkedResources{
		"file": {
			ResourceTypeName: "local_file",
			AttributePath:    path.Root("file"),
		},
	}
	resp.Schema = schema.Schema{
		Description: "Generates a local file with the given content.",
		Attributes: map[string]schema.Attribute{
			"file": schema.ResourceAttribute{ // new ResourceAttribute which is backed by a Framework types.Object
				//TypeName:    "local_file",
				Resource:    NewLocalFileResource(),
				Description: "The file resource to be modified by this action.",
				//DriftableAttributes: []path.Path{ // Paths to attributes in 'local_file' schema that can cause action drift
				//	path.Root("content"),
				//	path.Root("content_md5"),
				//	path.Root("content_sha1"),
				//	path.Root("content_sha256"),
				//	path.Root("content_base64sha256"),
				//	path.Root("content_sha512"),
				//	path.Root("content_base64sha512"),
				//},
				Required: true,
			},
			"content": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (l *localAppendFileAction) Plan(ctx context.Context, req action.PlanRequest, resp *action.PlanResponse) {
	var plan actionModel
	var filePerm string

	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//content := plan.Content.ValueString()

	filePlan := localFileResourceModelV0{}

	diags = plan.File.As(ctx, filePlan, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destination := filePlan.Filename.ValueString()

	filePerm = filePlan.FilePermission.ValueString()

	fileMode, _ := strconv.ParseInt(filePerm, 8, 64)

	_, err := os.ReadFile(destination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read local file error",
			"An unexpected error occurred while reading the file\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
	}

	newContent := []byte("actionContent")
	if err := os.WriteFile(destination, newContent, os.FileMode(fileMode)); err != nil {
		resp.Diagnostics.AddError(
			"Create local file error",
			"An unexpected error occurred while writing the file\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	checksums := genFileChecksums(newContent)
	filePlan.Content = types.StringValue(string(newContent))
	filePlan.ContentMd5 = types.StringValue(checksums.md5Hex)
	filePlan.ContentSha1 = types.StringValue(checksums.sha1Hex)
	filePlan.ContentSha256 = types.StringValue(checksums.sha256Hex)
	filePlan.ContentBase64sha256 = types.StringValue(checksums.sha256Base64)
	filePlan.ContentSha512 = types.StringValue(checksums.sha512Hex)
	filePlan.ContentBase64sha512 = types.StringValue(checksums.sha512Base64)

	filePlan.ID = types.StringValue(checksums.sha1Hex)

	objectValue, diags := types.ObjectValueFrom(ctx, filePlan.AttributeTypes(), filePlan)
	plan.File = objectValue

	diags = resp.PlannedConfig.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (l *localAppendFileAction) Invoke(ctx context.Context, req action.InvokeRequest, resp *action.InvokeResponse) {
	err := resp.CallbackServer.Send(ctx, &action.StartedActionEvent{
		CancellationToken: "randomToken",
	})
	if err != nil {
		resp.Diagnostics.AddError("Error: sending started event", fmt.Sprintf("Original error %s", err.Error()))
		return
	}
	newConfig := &tfsdk.State{
		Raw:    req.Config.Raw.Copy(),
		Schema: req.Config.Schema,
	}

	var plan actionModel
	var filePerm string

	diags := req.Config.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	//content := plan.Content.ValueString()

	filePlan := localFileResourceModelV0{}

	diags = plan.File.As(ctx, filePlan, basetypes.ObjectAsOptions{})
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	destination := filePlan.Filename.ValueString()

	filePerm = filePlan.FilePermission.ValueString()

	fileMode, _ := strconv.ParseInt(filePerm, 8, 64)

	_, err = os.ReadFile(destination)
	if err != nil {
		resp.Diagnostics.AddError(
			"Read local file error",
			"An unexpected error occurred while reading the file\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
	}

	// newContent := append(readContent, []byte(content)...)
	newContent := []byte("actionContent")
	if err := os.WriteFile(destination, newContent, os.FileMode(fileMode)); err != nil {
		resp.Diagnostics.AddError(
			"Create local file error",
			"An unexpected error occurred while writing the file\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	checksums := genFileChecksums(newContent)
	filePlan.Content = types.StringValue(string(newContent))
	filePlan.ContentMd5 = types.StringValue(checksums.md5Hex)
	filePlan.ContentSha1 = types.StringValue(checksums.sha1Hex)
	filePlan.ContentSha256 = types.StringValue(checksums.sha256Hex)
	filePlan.ContentBase64sha256 = types.StringValue(checksums.sha256Base64)
	filePlan.ContentSha512 = types.StringValue(checksums.sha512Hex)
	filePlan.ContentBase64sha512 = types.StringValue(checksums.sha512Base64)

	filePlan.ID = types.StringValue(checksums.sha1Hex)

	objectValue, diags := types.ObjectValueFrom(ctx, filePlan.AttributeTypes(), filePlan)
	plan.File = objectValue

	diags = newConfig.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	err = resp.CallbackServer.Send(ctx, &action.FinishedActionEvent{
		NewConfig: newConfig,
	})
	if err != nil {
		resp.Diagnostics.AddError("Error: sending started event", fmt.Sprintf("Original error %s", err.Error()))
		return
	}
}

func (l *localAppendFileAction) Cancel(ctx context.Context, req action.CancelRequest, resp *action.CancelResponse) {
	//TODO implement me
	panic("implement me")
}

type actionModel struct {
	File    types.Object `tfsdk:"file"`
	Content types.String `tfsdk:"content"`
}
