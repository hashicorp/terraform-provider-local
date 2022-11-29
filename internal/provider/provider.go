package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ provider.Provider             = (*localProvider)(nil)
	_ provider.ProviderWithSchema   = (*localProvider)(nil)
	_ provider.ProviderWithMetadata = (*localProvider)(nil)
)

func New() provider.Provider {
	return &localProvider{}
}

type localProvider struct{}

func (p *localProvider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "local"
}

func (p *localProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {

}

func (p *localProvider) DataSources(ctx context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewLocalFileDataSource,
		NewLocalSensitiveFileDataSource,
	}
}

func (p *localProvider) Resources(ctx context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewLocalFileResource,
		NewLocalSensitiveFileResource,
	}
}

func (p *localProvider) Schema(ctx context.Context, req provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{}
}
