package provider

import (
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource           = (*localFileDataSource)(nil)
	_ datasource.DataSourceWithSchema = (*localFileDataSource)(nil)
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
				Description: "Raw content of the file that was read, as UTF-8 encoded string.",
				Computed:    true,
			},
			"content_base64": schema.StringAttribute{
				Description: "Base64 encoded version of the file content (use this when dealing with binary data).",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The hexadecimal encoding of the checksum of the file content",
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
	content, err := ioutil.ReadFile(filepath)

	if err != nil {
		resp.Diagnostics.AddError(
			"Read local file data source error",
			"The file at given path cannot be read.\n\n+"+
				fmt.Sprintf("Original Error: %s", err),
		)
		return
	}

	//calculate the checksum of file content
	checksum := sha1.Sum(content)

	state := localFileDataSourceModelV0{
		Filename:      config.Filename,
		Content:       types.StringValue(string(content)),
		ContentBase64: types.StringValue(base64.StdEncoding.EncodeToString(content)),
		ID:            types.StringValue(hex.EncodeToString(checksum[:])),
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

type localFileDataSourceModelV0 struct {
	Filename      types.String `tfsdk:"filename"`
	Content       types.String `tfsdk:"content"`
	ContentBase64 types.String `tfsdk:"content_base64"`
	ID            types.String `tfsdk:"id"`
}
