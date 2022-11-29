package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

var (
	_ datasource.DataSource           = (*localSensitiveFileDataSource)(nil)
	_ datasource.DataSourceWithSchema = (*localSensitiveFileDataSource)(nil)
)

func NewLocalSensitiveFileDataSourceWithSchema() datasource.DataSourceWithSchema {
	return &localSensitiveFileDataSource{}
}

func NewLocalSensitiveFileDataSource() datasource.DataSource {
	return &localSensitiveFileDataSource{}
}

type localSensitiveFileDataSource struct{}

func (n *localSensitiveFileDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sensitive_file"
}

func (n *localSensitiveFileDataSource) Schema(ctx context.Context, req datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Reads a file from the local filesystem.",
		Attributes: map[string]schema.Attribute{
			"filename": schema.StringAttribute{
				Description: "Path to the file that will be read. The data source will return an error if the file does not exist.",
				Required:    true,
			},
			"content": schema.StringAttribute{
				Description: "Raw content of the file that was read, as UTF-8 encoded string.",
				Sensitive:   true,
				Computed:    true,
			},
			"content_base64": schema.StringAttribute{
				Description: "Base64 encoded version of the file content (use this when dealing with binary data).",
				Sensitive:   true,
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Description: "The hexadecimal encoding of the checksum of the file content",
				Computed:    true,
			},
		},
	}
}

func (n *localSensitiveFileDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// NOTE: We can use the read-method for the data source `local_file` as-is, because
	// all this data source does, is adding "Sensitive: true" to the schema of the property.
	//
	// The values and the property names are meant to be kept the same between data sources.
	NewLocalFileDataSource().Read(ctx, req, resp)
}
